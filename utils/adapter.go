package utils

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/axllent/semver"
	"github.com/go-git/go-git/v5"
	"github.com/panzouh/argocd-exporter/scheme"
	"gopkg.in/yaml.v2"
)

// GitRepositoryAdapter adapts Git repositories
type GitRepositoryAdapter struct{}

// HelmRepositoryAdapter adapts Helm repositories
type HelmRepositoryAdapter struct{}

// statusCodeMap = map[string]float64{
// 	"up-to-date":          0,
// 	"unknown":         1,
// 	"pending-upgrade":  2,
// 	"pending-rollback": 3,
// }

// Chart represents items in the index.yaml file in a Helm repository
type Chart struct {
	Name        string            `yaml:"name"`
	Version     string            `yaml:"version"`
	Digest      string            `yaml:"digest"`
	URLs        []string          `yaml:"urls"`
	Created     string            `yaml:"created"`
	Deprecated  bool              `yaml:"deprecated"`
	Annotations map[string]string `yaml:"annotations"`
}

// Index represents the index.yaml file in a Helm repository
type Index struct {
	APIVersion string             `yaml:"apiVersion"`
	Entries    map[string][]Chart `yaml:"entries"`
}

// RepositoryAnalyzer adapts and analyzes repositories
type RepositoryAnalyzer struct {
	GitAdapter  *GitRepositoryAdapter
	HelmAdapter *HelmRepositoryAdapter
}

// Local directory where git repositories are cloned
var tmpDir = "tmp"

// IsGitRepository checks if the URL is a Git repository
func (g *GitRepositoryAdapter) IsGitRepository(app scheme.Application) bool {
	return strings.HasSuffix(app.Spec.Source.RepoURL, ".git")
}

// ScanChartYAML searches for the Chart.yaml file in the Git repository
func (g *GitRepositoryAdapter) ScanChartYAML(app scheme.Application) (string, error) {
	_, err := git.PlainClone(tmpDir+app.Metadata.Name, false, &git.CloneOptions{
		URL:        app.Spec.Source.RepoURL,
		RemoteName: app.Spec.Source.TargetRevision,
	})
	if err != nil {
		return "", err
	}
	repo, err := git.PlainOpen(tmpDir + app.Metadata.Name)
	if err != nil {
		return "", err
	}
	ref, err := repo.Head()
	if err != nil {
		return "", err
	}
	commit, err := repo.CommitObject(ref.Hash())
	if err != nil {
		return "", err
	}
	tree, err := commit.Tree()
	if err != nil {
		return "", err
	}
	filePath := app.Spec.Source.Path + "/Chart.yaml"
	file, err := tree.File(filePath)
	if err != nil {
		return "", err
	}
	chartParser := Chart{}
	content, err := file.Contents()
	if err != nil {
		return "", err
	}
	os.RemoveAll(tmpDir + app.Metadata.Name)
	err = yaml.Unmarshal([]byte(content), &chartParser)
	if err != nil {
		return "", err
	}
	return chartParser.Version, nil
}

// IsHelmRepository checks if the URL is a Helm repository
func (h *HelmRepositoryAdapter) IsHelmRepository(app scheme.Application) bool {
	// return true if the http://<url>/index.yaml exists
	cli, err := http.Get(app.Spec.Source.RepoURL + "/index.yaml")
	if err != nil {
		return false
	}
	if cli.StatusCode == 200 {
		return true
	}
	return false
}

// ScanIndexYAML scans the index.yaml file in the Helm repository
func (h *HelmRepositoryAdapter) ScanIndexYAML(app scheme.Application) (string, error) {
	indexUrl := app.Spec.Source.RepoURL + "/index.yaml"
	resp, err := http.Get(indexUrl)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	index := Index{}
	err = yaml.Unmarshal(body, &index)
	if err != nil {
		return "", err
	}

	if charts, ok := index.Entries[app.Spec.Source.Chart]; ok {
		return charts[0].Version, nil

	} else {
		return "", errors.New("couldn't fetch the latest chart version")
	}
}

// AnalyzeRepository analyzes the repository based on the provided URL
func (r *RepositoryAnalyzer) AnalyzeRepository(app scheme.Application) (string, string, error) {
	if r.GitAdapter.IsGitRepository(app) {
		latestVersion, err := r.GitAdapter.ScanChartYAML(app)
		if err != nil {
			return "unknown", "", err
		}
		compare := semver.Compare(latestVersion, app.Spec.Source.TargetRevision)
		if compare == 0 {
			return "up-to-date", latestVersion, nil
		}
		if compare == 1 {
			return "pending-upgrade", latestVersion, nil
		}
		if compare == -1 {
			return "pending-rollback", latestVersion, nil
		}
		return "unknown", latestVersion, nil
	} else if r.HelmAdapter.IsHelmRepository(app) {
		latestVersion, err := r.HelmAdapter.ScanIndexYAML(app)
		if err != nil {
			return "unknown", "", err
		}
		compare := semver.Compare(latestVersion, app.Spec.Source.TargetRevision)
		if compare == 0 {
			return "up-to-date", latestVersion, nil
		}
		if compare == 1 {
			return "pending-upgrade", latestVersion, nil
		}
		if compare == -1 {
			return "pending-rollback", latestVersion, nil
		}
		return "unknown", latestVersion, nil
	}
	return "", "", fmt.Errorf("unsupported repository type")
}
