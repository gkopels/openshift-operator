// SPDX-License-Identifier: Apache-2.0
// Copyright (c) 2020-2021 Intel Corporation

package main

import (
	"flag"
	"os"

	sriovv1 "github.com/open-ness/openshift-operator/sriov-fec/api/v1"
	"github.com/open-ness/openshift-operator/sriov-fec/pkg/daemon"

	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientset "k8s.io/client-go/kubernetes"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
)

var (
	scheme   = runtime.NewScheme()
	setupLog = ctrl.Log.WithName("setup")
)

func init() {
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))
	utilruntime.Must(sriovv1.AddToScheme(scheme))
}

func main() {
	opts := zap.Options{}
	opts.BindFlags(flag.CommandLine)
	flag.Parse()

	ctrl.SetLogger(zap.New(zap.UseFlagOptions(&opts)))

	nodeName := os.Getenv("NODENAME")
	if nodeName == "" {
		setupLog.Error(nil, "NODENAME environment variable is empty")
		os.Exit(1)
	}

	ns := os.Getenv("SRIOV_FEC_NAMESPACE")
	if ns == "" {
		setupLog.Error(nil, "SRIOV_FEC_NAMESPACE environment variable is empty")
		os.Exit(1)
	}

	config := ctrl.GetConfigOrDie()
	directClient, err := client.New(config, client.Options{Scheme: scheme})
	if err != nil {
		setupLog.Error(err, "failed to create direct client")
		os.Exit(1)
	}

	cset, err := clientset.NewForConfig(config)
	if err != nil {
		setupLog.Error(err, "failed to create clientset")
		os.Exit(1)
	}

	mgr, err := ctrl.NewManager(config, ctrl.Options{
		Scheme:             scheme,
		MetricsBindAddress: "0",
		LeaderElection:     false,
		Namespace:          ns,
	})
	if err != nil {
		setupLog.Error(err, "unable to start manager")
		os.Exit(1)
	}

	daemon, err := daemon.NewNodeConfigReconciler(mgr.GetClient(), cset, ctrl.Log.WithName("daemon"), nodeName, ns)
	if err != nil {
		setupLog.Error(err, "unable to create reconciler")
		os.Exit(1)
	}

	if err := daemon.SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "NodeConfig")
		os.Exit(1)
	}

	if err := daemon.CreateEmptyNodeConfigIfNeeded(directClient); err != nil {
		setupLog.Error(err, "failed to create initial NodeConfig CR")
		os.Exit(1)
	}

	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		setupLog.Error(err, "problem running manager")
		os.Exit(1)
	}
}
