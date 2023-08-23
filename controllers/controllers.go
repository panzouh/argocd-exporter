package controllers

import (
	discovery "github.com/gkarthiks/k8s-discovery"
	"github.com/rs/zerolog"
	"github.com/panzouh/argocd-exporter/metrics"
	"github.com/panzouh/argocd-exporter/utils"
	"k8s.io/client-go/kubernetes"
)

type Controller struct {
	MetricsServer *metrics.MetricsServer
	ClientSet     *kubernetes.Clientset
	Analyze       *utils.RepositoryAnalyzer
	Logger        *zerolog.Logger
}

func NewControllers(auth *discovery.K8s, verbosity string) (*Controller, error) {
	ClientSet, err := kubernetes.NewForConfig(auth.RestConfig)
	if err != nil {
		return nil, err
	}
	return &Controller{
		MetricsServer: metrics.NewMetricsController(),
		ClientSet:     ClientSet,
		Analyze:       &utils.RepositoryAnalyzer{},
		Logger:        utils.SetupLogger(verbosity),
	}, nil
}
