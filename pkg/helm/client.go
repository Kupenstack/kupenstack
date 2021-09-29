package helm

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/gofrs/flock"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"

	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/downloader"
	"helm.sh/helm/v3/pkg/getter"
	"helm.sh/helm/v3/pkg/release"
	"helm.sh/helm/v3/pkg/repo"
	"helm.sh/helm/v3/pkg/storage/driver"
)

var settings *cli.EnvSettings = cli.New()

// AddRepoIfNotExist will add the repo if the repo doesn't exist
func AddRepoIfNotExist(repoName string, repoUrl string) error {
	repoFile := settings.RepositoryConfig

	//Ensure the file directory exists as it is required for file locking
	err := os.MkdirAll(filepath.Dir(repoFile), os.ModePerm)
	if err != nil && !os.IsExist(err) {
		return err
	}

	// Acquire a file lock for process synchronization
	fileLock := flock.New(strings.Replace(repoFile, filepath.Ext(repoFile), ".lock", 1))
	lockCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	locked, err := fileLock.TryLockContext(lockCtx, time.Second)
	if err == nil && locked {
		defer fileLock.Unlock()
	}
	if err != nil {
		return err
	}

	data, err := ioutil.ReadFile(repoFile)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	var file repo.File
	err = yaml.Unmarshal(data, &file)
	if err != nil {
		return err
	}

	if file.Has(repoName) {
		return nil
	}

	repoEntry := repo.Entry{
		Name: repoName,
		URL:  repoUrl,
	}

	repository, err := repo.NewChartRepository(&repoEntry, getter.All(settings))
	if err != nil {
		return err
	}

	if _, err := repository.DownloadIndexFile(); err != nil {
		err = errors.Wrap(err, "looks like "+repoUrl+" is not a valid chart repository or cannot be reached")
		return err
	}

	file.Update(&repoEntry)

	if err := file.WriteFile(repoFile, 0644); err != nil {
		return err
	}
	return nil
}

// ListReleases returns a list of all the releases for the given namespaces.
// Inputs: takes a list of namespaces to search releases in.
// If empty string is given as input then returns all releases from all namespaces.
func ListReleases(namespaces ...string) ([]*release.Release, error) {
	if namespaces == nil {
		// search across all namespaces.
		namespaces = append(namespaces, "")
	}

	var arr []*release.Release
	config := new(action.Configuration)

	for _, namespace := range namespaces {
		if err := config.Init(nil, namespace, os.Getenv("HELM_DRIVER"), debug); err != nil {
			return nil, err
		}

		client := action.NewList(config)
		client.Deployed = true

		releases, err := client.Run()
		if err != nil {
			return nil, err
		}

		arr = append(arr, releases...)
	}
	return arr, nil
}

func updateRepository(wg *sync.WaitGroup, repository *repo.ChartRepository, errchan chan error) {
	wg.Add(1)
	defer wg.Done()

	_, err := repository.DownloadIndexFile()
	if err != nil {
		errchan <- errors.Wrap(err, "Unable to get an update from the "+repository.Config.Name+"/"+repository.Config.URL)
	}

}

// UpdateHelmRepos will update all the repo
func UpdateHelmRepos() error {
	repoFile := settings.RepositoryConfig

	file, err := repo.LoadFile(repoFile)
	if os.IsNotExist(errors.Cause(err)) || len(file.Repositories) == 0 {
		return errors.New("No repositories found. You must add one before updating.")
	}

	// get all repositories
	var repositories []*repo.ChartRepository
	for _, cfg := range file.Repositories {
		repository, err := repo.NewChartRepository(cfg, getter.All(settings))
		if err != nil {
			return err
		}
		repositories = append(repositories, repository)
	}

	// update all repositories
	var wg sync.WaitGroup
	errchan := make(chan error)
	for _, repository := range repositories {
		go updateRepository(&wg, repository, errchan)
	}
	wg.Wait()
	close(errchan)
	// return first error
	return <-errchan
}

// UpgradeRelease upgrades a existing release or creates it if not exists.
func UpgradeRelease(name, repo, chart, namespace string, vals map[string]interface{}) (*release.Release, error) {
	cfg := new(action.Configuration)
	if err := cfg.Init(settings.RESTClientGetter(), namespace, os.Getenv("HELM_DRIVER"), debug); err != nil {
		return nil, err
	}

	err := SetNamespace(cfg, namespace)
	if err != nil {
		return nil, err
	}

	upgradeClient := action.NewUpgrade(cfg)

	chartPath, err := upgradeClient.ChartPathOptions.LocateChart(fmt.Sprintf("%s/%s", repo, chart), settings)
	if err != nil {
		return nil, err
	}

	chartRequested, err := loader.Load(chartPath)
	if err != nil {
		return nil, err
	}

	validInstallableChart, err := isChartInstallable(chartRequested)
	if !validInstallableChart {
		return nil, err
	}

	err = checkDependencies(chartRequested, chartPath, upgradeClient)
	if err != nil {
		return nil, err
	}

	upgradeClient.Namespace = namespace
	upgradeClient.Install = true
	upgradeClient.DryRun = false

	result, err := upgradeClient.Run(name, chartRequested, vals)

	if isReleaseDoesNotExistsErrorWithName(name, err) {
		return installRelease(cfg, name, namespace, vals, chartRequested)
	}
	return result, err
}

func checkDependencies(helmChart *chart.Chart, chartPath string, client *action.Upgrade) error {
	req := helmChart.Metadata.Dependencies
	if req == nil {
		return nil
	}

	err := action.CheckDependencies(helmChart, req)
	if err != nil {

		if !client.DependencyUpdate {
			return err
		}

		manager := &downloader.Manager{
			Out:              os.Stdout,
			ChartPath:        chartPath,
			Keyring:          client.ChartPathOptions.Keyring,
			SkipUpdate:       false,
			Getters:          getter.All(settings),
			RepositoryConfig: settings.RepositoryConfig,
			RepositoryCache:  settings.RepositoryCache,
		}

		err := manager.Update()
		if err != nil {
			return err
		}
	}

	return nil
}

func isReleaseDoesNotExistsErrorWithName(name string, err error) bool {

	ErrNoDeployedReleasesMsg := fmt.Sprintf("%q %s", name, driver.ErrNoDeployedReleases.Error())

	if err != nil && err.Error() == ErrNoDeployedReleasesMsg {
		return true
	} else {
		return false
	}
}

func installRelease(cfg *action.Configuration, name, namespace string, vals map[string]interface{},
	helmChart *chart.Chart) (*release.Release, error) {
	client := action.NewInstall(cfg)
	client.ReleaseName = name
	client.Namespace = namespace
	client.DryRun = false
	client.CreateNamespace = true
	if client.Version == "" && client.Devel {
		client.Version = ">0.0.0-0"
	}
	return client.Run(helmChart, vals)

}

// DeleteRelease will delete the given release
func DeleteRelease(name, namespace string) error {

	actionConfig := new(action.Configuration)
	err := actionConfig.Init(settings.RESTClientGetter(), namespace, os.Getenv("HELM_DRIVER"), debug)
	if err != nil {
		return err
	}

	//Set namespace
	err = SetNamespace(actionConfig, namespace)
	if err != nil {
		return err
	}

	client := action.NewUninstall(actionConfig)
	_, err = client.Run(name)
	if err != nil {
		return err
	}
	return nil
}

func isChartInstallable(ch *chart.Chart) (bool, error) {
	switch ch.Metadata.Type {
	case "", "application":
		return true, nil
	}
	return false, errors.Errorf("%s charts are not installable", ch.Metadata.Type)
}

//TODO: for printing logs from helm client
func debug(format string, v ...interface{}) {
	_ = fmt.Sprintf("[debug] %s\n", format)
}

func GetRelease(name, namespace string) (*release.Release, error) {

	releases, err := ListReleases(namespace)
	if err != nil {
		return nil, err
	}

	for _, releaseDetail := range releases {

		if releaseDetail.Name == name {
			return releaseDetail, nil
		}
	}
	return nil, nil
}
