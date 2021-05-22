package main

import (
	"embed"
	"fmt"
	"os"

	"github.com/go-logr/logr"
	"github.com/go-logr/zapr"
	flag "github.com/spf13/pflag"
	"go.uber.org/zap"
	"sigs.k8s.io/controller-runtime/pkg/manager/signals"

	k8sversion "github.com/phillebaba/which-k8s-version/pkg/k8s-version"
	"github.com/phillebaba/which-k8s-version/pkg/server"
)

//go:embed static/*
var staticFs embed.FS

var (
	port           string
	subscriptionID string
	projectID      string
)

func init() {
	flag.StringVar(&port, "port", ":8080", "port to bind proxy server to.")
	flag.StringVar(&subscriptionID, "subscription-id", "", "Azure subscription ID.")
	flag.StringVar(&projectID, "project-id", "", "GCP project id.")
	flag.Parse()
}

func main() {
	var logger logr.Logger
	zapLog, err := zap.NewProduction()
	if err != nil {
		panic(fmt.Sprintf("who watches the watchmen (%v)?", err))
	}
	logger = zapr.NewLogger(zapLog)
	mainLog := logger.WithName("main")

	ctx := signals.SetupSignalHandler()

	mainLog.Info("Getting version data")
	vss, err := k8sversion.GetVersions(ctx, subscriptionID, projectID)
	if err != nil {
		mainLog.Error(err, "Could not get version data")
		os.Exit(1)
	}

	mainLog.Info("Starting application")
	srv, err := server.NewFrontendServer(logger, staticFs, port, vss)
	if err != nil {
		mainLog.Error(err, "Could not initialize frontend server")
		os.Exit(1)
	}
	srv.ListenAndServe(ctx.Done())
	mainLog.Info("Stopped application")
}
