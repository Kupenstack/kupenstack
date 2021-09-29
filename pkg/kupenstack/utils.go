package kupenstack

import (
	"context"
	"strings"

	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"

	clusterv1alpha1 "github.com/kupenstack/kupenstack/apis/cluster/v1alpha1"
)

func OccpExists(c client.Client, profilename string) (bool, error) {
	profile := strings.Split(profilename, ".")
	profileDetail := &clusterv1alpha1.OpenStackCloudConfigurationProfile{}
	err := c.Get(context.Background(), types.NamespacedName{Name: profile[0], Namespace: profile[1]}, profileDetail)
	if err != nil {
		return false, err
	}
	return true, nil
}
