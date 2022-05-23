/*
Copyright 2021 The Kupenstack Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"flag"
	"os"

	// Import all Kubernetes client auth plugins (e.g. Azure, GCP, OIDC, etc.)
	// to ensure that exec-entrypoint and run can make use of them.
	_ "k8s.io/client-go/plugin/pkg/client/auth"

	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"

	clusterv1alpha1 "github.com/kupenstack/kupenstack/apis/cluster/v1alpha1"
	kupenstackiov1alpha1 "github.com/kupenstack/kupenstack/apis/v1alpha1"

	clustercontrollers "github.com/kupenstack/kupenstack/controllers/cluster"
	"github.com/kupenstack/kupenstack/controllers/flavor"
	"github.com/kupenstack/kupenstack/controllers/image"
	"github.com/kupenstack/kupenstack/controllers/keypair"
	"github.com/kupenstack/kupenstack/controllers/network"
	"github.com/kupenstack/kupenstack/controllers/project"
	"github.com/kupenstack/kupenstack/controllers/vm"
	"github.com/kupenstack/kupenstack/controllers/vn"
	"github.com/kupenstack/kupenstack/oskops"
	"github.com/kupenstack/kupenstack/pkg/openstack"
	//+kubebuilder:scaffold:imports
)

var (
	scheme   = runtime.NewScheme()
	setupLog = ctrl.Log.WithName("kupenstack.main")
)

func init() {
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))

	utilruntime.Must(kupenstackiov1alpha1.AddToScheme(scheme))
	utilruntime.Must(clusterv1alpha1.AddToScheme(scheme))
	//+kubebuilder:scaffold:scheme
}

func main() {
	var metricsAddr string
	var enableLeaderElection bool
	var probeAddr string
	var kupenstackConfigurationFile string
	flag.StringVar(&metricsAddr, "metrics-bind-address", ":8080", "The address the metric endpoint binds to.")
	flag.StringVar(&probeAddr, "health-probe-bind-address", ":8081", "The address the probe endpoint binds to.")
	flag.BoolVar(&enableLeaderElection, "leader-elect", false,
		"Enable leader election for controller manager. "+
			"Enabling this will ensure there is only one active controller manager.")
	flag.StringVar(&kupenstackConfigurationFile, "kupenstack-configuration-file", "config.yaml", "The filepath to KupenstackConfiguration.")
	opts := zap.Options{
		Development: true,
	}
	opts.BindFlags(flag.CommandLine)
	flag.Parse()

	ctrl.SetLogger(zap.New(zap.UseFlagOptions(&opts)))

	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme:                 scheme,
		MetricsBindAddress:     metricsAddr,
		Port:                   9443,
		HealthProbeBindAddress: probeAddr,
		LeaderElection:         enableLeaderElection,
		LeaderElectionID:       "e2bdf9e1.kupenstack.io",
	})
	if err != nil {
		setupLog.Error(err, "unable to start manager")
		os.Exit(1)
	}

	go oskops.ManageOskNodes(mgr.GetClient(), kupenstackConfigurationFile)
	go oskops.Start(mgr.GetClient(), kupenstackConfigurationFile)

	///// Temporary code

	OSclient := &openstack.Client{}
	go oskops.AuthenticateOpenstackClient(OSclient)

	//////////////////////////

	if err = (&project.Reconciler{
		Client:        mgr.GetClient(),
		OS:            OSclient,
		Log:           ctrl.Log.WithName("controllers").WithName("Project"),
		Scheme:        mgr.GetScheme(),
		EventRecorder: mgr.GetEventRecorderFor("kupenstack-controller"),
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "Project")
		os.Exit(1)
	}

	if err = (&keypair.Reconciler{
		Client:        mgr.GetClient(),
		OS:            OSclient,
		Log:           ctrl.Log.WithName("controllers").WithName("KeyPair"),
		Scheme:        mgr.GetScheme(),
		EventRecorder: mgr.GetEventRecorderFor("kupenstack-controller"),
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "KeyPair")
		os.Exit(1)
	}
	if err = (&image.Reconciler{
		Client:        mgr.GetClient(),
		OS:            OSclient,
		Log:           ctrl.Log.WithName("controllers").WithName("Image"),
		Scheme:        mgr.GetScheme(),
		EventRecorder: mgr.GetEventRecorderFor("kupenstack-controller"),
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "Image")
		os.Exit(1)
	}
	if err = (&flavor.Reconciler{
		Client:        mgr.GetClient(),
		OS:            OSclient,
		Log:           ctrl.Log.WithName("controllers").WithName("Flavor"),
		Scheme:        mgr.GetScheme(),
		EventRecorder: mgr.GetEventRecorderFor("kupenstack-controller"),
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "Flavor")
		os.Exit(1)
	}
	if err = (&network.Reconciler{
		Client:        mgr.GetClient(),
		OS:            OSclient,
		Log:           ctrl.Log.WithName("controllers").WithName("Network"),
		Scheme:        mgr.GetScheme(),
		EventRecorder: mgr.GetEventRecorderFor("kupenstack-controller"),
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "Network")
		os.Exit(1)
	}
	if err = (&vm.Reconciler{
		Client:        mgr.GetClient(),
		OS:            OSclient,
		Log:           ctrl.Log.WithName("controllers").WithName("VirtualMachine"),
		Scheme:        mgr.GetScheme(),
		EventRecorder: mgr.GetEventRecorderFor("kupenstack-controller"),
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "VirtualMachine")
		os.Exit(1)
	}
	if err = (&clustercontrollers.Reconciler{
		Client:        mgr.GetClient(),
		Log:           ctrl.Log.WithName("controllers").WithName("OpenstackNode"),
		Scheme:        mgr.GetScheme(),
		EventRecorder: mgr.GetEventRecorderFor("kupenstack-controller"),
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "OpenstackNode")
		os.Exit(1)
	}
	if err = (&vn.Reconciler{
		Client:        mgr.GetClient(),
		OS:            OSclient,
		Log:           ctrl.Log.WithName("controllers").WithName("VirtualNetwork"),
		Scheme:        mgr.GetScheme(),
		EventRecorder: mgr.GetEventRecorderFor("kupenstack-controller"),
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "VirtualNetwork")
		os.Exit(1)
	}
	//+kubebuilder:scaffold:builder

	if err := mgr.AddHealthzCheck("healthz", healthz.Ping); err != nil {
		setupLog.Error(err, "unable to set up health check")
		os.Exit(1)
	}
	if err := mgr.AddReadyzCheck("readyz", healthz.Ping); err != nil {
		setupLog.Error(err, "unable to set up ready check")
		os.Exit(1)
	}

	setupLog.Info("starting manager")
	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		setupLog.Error(err, "problem running manager")
		os.Exit(1)
	}
}
