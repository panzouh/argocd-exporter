package main

import (
	"flag"
	"net/http"
	"time"

	ginzerolog "github.com/dn365/gin-zerolog"
	"github.com/gin-gonic/gin"
	discovery "github.com/gkarthiks/k8s-discovery"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog/log"
	"github.com/panzouh/argocd-exporter/controllers"
)

func main() {
	// Set Gin release mode
	gin.SetMode(gin.ReleaseMode)

	address := flag.String("address", ":", "Address to listen on")
	port := flag.String("port", "8080", "Port to listen on")
	metricsPath := flag.String("metrics-path", "/metrics", "Path to metrics")
	verbosity := flag.String("verbosity", "info", "Verbosity level")
	updateInterval := flag.Duration("update-interval", 1, "Interval in minutes to update metrics")
	flag.Parse()

	// Create controller
	auth, err := discovery.NewK8s()
	if err != nil {
		panic(err)
	}
	controller, err := controllers.NewControllers(auth, *verbosity)
	if err != nil {
		panic(err)
	}

	// Register metrics
	controller.Register()

	// Start updating metrics
	go func() {
		for {
			controller.UpdateArgocdAppVersions()
			time.Sleep(*updateInterval * time.Minute)
		}
	}()

	// Create Gin router
	r := gin.New()

	r.Use(ginzerolog.Logger("gin"))

	// Add permanent redirect to metrics endpoint
	r.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, *metricsPath)
	})

	// Add /metrics endpoint
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// Add /liveness endpoint
	r.GET("/liveness", controller.Liveness)

	// Add /readiness endpoint
	r.GET("/readiness", controller.Readiness)

	// Run the server
	log.Info().Msgf("Application started at %v%v", *address, *port)
	log.Fatal().Err(r.Run(*address + *port)).Msg("Application has stopped")
}
