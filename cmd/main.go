package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gorilla/mux"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s-monitor/pkg/config"
	"k8s-monitor/pkg/monitor"
)

// Version information - set via ldflags at build time
var (
	Version   = "dev"
	BuildDate = "unknown"
	GitCommit = "unknown"
)

type Server struct {
	monitor *monitor.K8sMonitor
	config  *config.Config
}

func main() {
	// Print version information
	log.Printf("🚀 Kubernetes Monitor")
	log.Printf("   Version: %s", Version)
	log.Printf("   Build Date: %s", BuildDate)
	log.Printf("   Git Commit: %s", GitCommit)
	log.Println()

	// Load configuration
	cfg, err := config.LoadConfig("config.json")
	if err != nil {
		log.Fatalf("Error loading configuration: %s", err.Error())
	}

	// Try to create Kubernetes client config
	// First try in-cluster config (for pods running inside Kubernetes/OpenShift)
	k8sConfig, err := rest.InClusterConfig()
	if err != nil {
		// If in-cluster config fails, try kubeconfig file (for local development)
		log.Printf("In-cluster config not available, trying kubeconfig file: %v", err)
		
		kubeconfig := os.Getenv("KUBECONFIG")
		if kubeconfig == "" {
			kubeconfig = clientcmd.RecommendedHomeFile
		}

		k8sConfig, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			log.Printf("Warning: Failed to initialize Kubernetes client: %v", err)
			log.Printf("The application will start but Kubernetes features will be unavailable")
		} else {
			log.Printf("Using kubeconfig from: %s", kubeconfig)
		}
	} else {
		log.Printf("Using in-cluster configuration")
	}

	var clientset *kubernetes.Clientset
	if k8sConfig != nil {
		clientset, err = kubernetes.NewForConfig(k8sConfig)
		if err != nil {
			log.Printf("Warning: Error creating Kubernetes client: %s", err.Error())
		}
	}

	var m *monitor.K8sMonitor
	if clientset != nil {
		m, err = monitor.NewK8sMonitor(clientset, cfg)
		if err != nil {
			log.Fatalf("Error initializing monitor: %s", err.Error())
		}
		if err := m.StartMonitoring(); err != nil {
			log.Fatalf("Error starting monitoring: %s", err.Error())
		}
	}

	server := &Server{monitor: m, config: cfg}

	// Setup routes
	router := mux.NewRouter()

	// Static file serving setup first
	webDir := "web"
	if _, err := os.Stat(webDir); os.IsNotExist(err) {
		execPath, _ := os.Executable()
		execDir := filepath.Dir(execPath)
		webDir = filepath.Join(execDir, "..", "web")
		if _, err := os.Stat(webDir); os.IsNotExist(err) {
			webDir = filepath.Join(".", "web")
		}
	}

	// API routes (must be registered before static file handler)
	router.HandleFunc("/api/changes", server.handleAPIChanges).Methods("GET")
	router.HandleFunc("/api/stats", server.handleAPIStats).Methods("GET")
	router.HandleFunc("/api/config", server.handleAPIConfig).Methods("GET")
	router.HandleFunc("/api/mark-read", server.handleMarkRead).Methods("POST")
	router.HandleFunc("/api/mark-all-read", server.handleMarkAllRead).Methods("POST")
	router.HandleFunc("/api/save-now", server.handleSaveNow).Methods("POST")
	router.HandleFunc("/api/debug", server.debugStatus).Methods("GET")
	router.HandleFunc("/health", server.healthCheck).Methods("GET")

	// Serve static files (this must be last as it's a catch-all)
	router.PathPrefix("/").Handler(http.FileServer(http.Dir(webDir + "/")))

	// Start server
	port := fmt.Sprintf("%d", cfg.WebPort)
	if envPort := os.Getenv("PORT"); envPort != "" {
		port = envPort
	}

	fmt.Printf("🚀 Kubernetes Monitor starting on port %s\n", port)
	fmt.Printf("📂 Serving web files from: %s\n", webDir)
	if clientset != nil {
		fmt.Printf("🔗 Connected to Kubernetes cluster\n")
	}
	fmt.Printf("🌐 Open http://localhost:%s in your browser\n", port)
	
	log.Fatal(http.ListenAndServe(":"+port, router))
}

func (s *Server) debugStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	status := map[string]interface{}{
		"status":    "healthy",
		"version":   Version,
		"buildDate": BuildDate,
		"gitCommit": GitCommit,
	}
	
	if s.monitor != nil {
		stats := s.monitor.GetStats()
		status["monitoring"] = map[string]interface{}{
			"active":       true,
			"totalChanges": stats["totalChanges"],
			"uptime":      stats["uptime"],
		}
	} else {
		status["monitoring"] = map[string]interface{}{
			"active": false,
			"error":  "Kubernetes client not available",
		}
	}
	
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(status)
}

func (s *Server) healthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"healthy"}`))
}

func (s *Server) handleAPIChanges(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if s.monitor == nil {
		http.Error(w, "Monitor not available", http.StatusServiceUnavailable)
		return
	}
	changes := s.monitor.GetChanges()
	json.NewEncoder(w).Encode(changes)
}

func (s *Server) handleAPIStats(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if s.monitor == nil {
		http.Error(w, "Monitor not available", http.StatusServiceUnavailable)
		return
	}
	stats := s.monitor.GetStats()
	json.NewEncoder(w).Encode(stats)
}

func (s *Server) handleAPIConfig(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(s.config)
}

func (s *Server) handleMarkRead(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if s.monitor == nil {
		http.Error(w, "Monitor not available", http.StatusServiceUnavailable)
		return
	}

	var req struct {
		ID string `json:"id"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	success := s.monitor.MarkAsRead(req.ID)
	
	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"success": success,
	}
	json.NewEncoder(w).Encode(response)
}

func (s *Server) handleMarkAllRead(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if s.monitor == nil {
		http.Error(w, "Monitor not available", http.StatusServiceUnavailable)
		return
	}

	count := s.monitor.MarkAllAsRead()
	
	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"success": true,
		"count":   count,
	}
	json.NewEncoder(w).Encode(response)
}

func (s *Server) handleSaveNow(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if s.monitor == nil {
		http.Error(w, "Monitor not available", http.StatusServiceUnavailable)
		return
	}

	err := s.monitor.SaveToFileNow()
	
	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"success": err == nil,
	}
	
	if err != nil {
		response["error"] = err.Error()
	}
	
	json.NewEncoder(w).Encode(response)
}