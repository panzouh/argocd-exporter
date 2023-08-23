package metrics

import "github.com/prometheus/client_golang/prometheus"

type Metric struct {
	Name   string
	Help   string
	Labels []string
	Gauge  prometheus.GaugeVec
}

type MetricsServer struct {
	ArgocdAppVersions *Metric
}

func NewMetricsController() *MetricsServer {
	return &MetricsServer{
		ArgocdAppVersions: NewMetric(
			"argocd_app_versions",
			"Versions of ArgoCD applications",
			[]string{
				"app_name",
				"repo_url",
				"version",
				"chart",
				"path",
				"status",
				"remote_version",
			},
		),
	}
}

func (c *MetricsServer) Register() {
	prometheus.MustRegister(c.ArgocdAppVersions.Get())
}

func NewMetric(name, help string, labels []string) *Metric {
	return &Metric{
		Name:   name,
		Help:   help,
		Labels: labels,
		Gauge: *prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: name,
			Help: help,
		}, labels),
	}
}

func (m *Metric) Get() prometheus.Collector {
	return m.Gauge
}

func (m *Metric) Set(n float64, labels ...string) {
	m.Gauge.WithLabelValues(labels...).Set(n)
}

func (m *Metric) With(labels ...string) prometheus.Gauge {
	return m.Gauge.WithLabelValues(labels...)
}

func (*MetricsServer) UpdateArgocdAppVersions() {

}
