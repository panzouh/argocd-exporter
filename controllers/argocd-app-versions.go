package controllers

import (
	"context"
	"encoding/json"

	"github.com/panzouh/argocd-exporter/scheme"
	"github.com/panzouh/argocd-exporter/utils"
)

var (
	StatusCodes = []string{
		"up-to-date",
		"pending-upgrade",
		"pending-rollback",
		"unknown",
	}
	Api      = "apis/argoproj.io/v1alpha1"
	Resource = "applications"
)

func (c *Controller) Register() {
	c.MetricsServer.Register()
}

func (c *Controller) UpdateArgocdAppVersions() {
	appsList := scheme.ApplicationsList{}

	// Get applications from Kubernetes API
	apps, err := c.ClientSet.RESTClient().
		Get().AbsPath(Api).
		Resource(Resource).
		DoRaw(context.Background())

	if err != nil {
		c.Logger.Error().Msgf("Error getting applications: %v", err)
	}

	// Parse apps JSON response
	if err := json.Unmarshal(apps, &appsList); err != nil {
		c.Logger.Error().Msgf("Error parsing applications: %v", err)
	}

	// Update ArgoCD app versions
	for _, app := range appsList.Items {
		result, version, err := c.Analyze.AnalyzeRepository(app)
		if err != nil {
			c.Logger.Error().Msgf("Error analyzing repository (%v): %v", app.Spec.Source.RepoURL, err)
		}
		c.MetricsServer.ArgocdAppVersions.With(
			app.Metadata.Name,
			app.Spec.Source.RepoURL,
			app.Spec.Source.TargetRevision,
			app.Spec.Source.Chart,
			app.Spec.Source.Path,
			result,
			version,
		).Set(utils.IndexOf(StatusCodes, result))
		c.Logger.Info().Msgf("Updated ArgoCD app versions for %v", app.Metadata.Name)
	}
}
