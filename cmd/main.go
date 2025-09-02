package main

import (
	"log"
	"os"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s-monitor/pkg/config"
	"k8s-monitor/pkg/monitor"
	"k8s-monitor/pkg/web"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig("config.json")
	if err != nil {
		log.Fatalf("Error loading configuration: %s", err.Error())
	}

	// Try to create Kubernetes client config
	// First try in-cluster config (for pods running inside Kubernetes/OpenShift)
	config, err := rest.InClusterConfig()
	if err != nil {
		// If in-cluster config fails, try kubeconfig file (for local development)
		log.Printf("In-cluster config not available, trying kubeconfig file: %v", err)
		
		kubeconfig := os.Getenv("KUBECONFIG")
		if kubeconfig == "" {
			kubeconfig = clientcmd.RecommendedHomeFile
		}

		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			log.Fatalf("Error building kubeconfig: %s", err.Error())
		}
		log.Printf("Using kubeconfig from: %s", kubeconfig)
	} else {
		log.Printf("Using in-cluster configuration")
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatalf("Error creating Kubernetes client: %s", err.Error())
	}

	m, err := monitor.NewK8sMonitor(clientset, cfg)
	if err != nil {
		log.Fatalf("Error initializing monitor: %s", err.Error())
	}
	if err := m.StartMonitoring(); err != nil {
		log.Fatalf("Error starting monitoring: %s", err.Error())
	}

	// Start web server
	webServer := web.NewWebServer(m, cfg.WebPort)
	log.Printf("Web UI available at http://localhost:%d", cfg.WebPort)
	
	if err := webServer.Start(); err != nil {
		log.Fatalf("Error starting web server: %s", err.Error())
	}
}