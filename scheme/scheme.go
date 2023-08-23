package scheme

type ApplicationsList struct {
	Items []Application `json:"items"`
}

type Application struct {
	Metadata struct {
		Name string `json:"name"`
	} `json:"metadata"`
	Spec struct {
		Source struct {
			RepoURL        string `json:"repoURL"`
			TargetRevision string `json:"targetRevision"`
			Chart          string `json:"chart"`
			Path           string `json:"path"`
		} `json:"source"`
	} `json:"spec"`
	Status struct {
		Health struct {
			Status string `json:"status"`
		} `json:"health"`
		History []HistoryItem `json:"history"`
	} `json:"status"`
}

type HistoryItem struct {
	DeployedAt      string `json:"deployedAt"`
	DeployStartedAt string `json:"deployStartedAt"`
	Revision        string `json:"revision"`
}
