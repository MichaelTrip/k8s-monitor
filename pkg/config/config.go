package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type ResourceConfig struct {
	Name        string `json:"name"`
	Enabled     bool   `json:"enabled"`
	Namespace   string `json:"namespace,omitempty"` // empty means all namespaces
	Description string `json:"description"`
}

type Config struct {
	WebPort     int              `json:"webPort"`
	Resources   []ResourceConfig `json:"resources"`
	Persistence PersistenceConfig `json:"persistence"`
}

type PersistenceConfig struct {
	Enabled    bool   `json:"enabled"`
	FilePath   string `json:"filePath"`
	AutoSave   bool   `json:"autoSave"`
	SaveInterval int  `json:"saveInterval"` // in seconds
}

func LoadConfig(configPath string) (*Config, error) {
	// Default configuration
	defaultConfig := &Config{
		WebPort: 8080,
		Persistence: PersistenceConfig{
			Enabled:      true,
			FilePath:     "changes.json",
			AutoSave:     true,
			SaveInterval: 30, // Save every 30 seconds
		},
		Resources: []ResourceConfig{
			{Name: "pods", Enabled: true, Description: "Kubernetes Pods"},
			{Name: "deployments", Enabled: true, Description: "Kubernetes Deployments"},
			{Name: "services", Enabled: true, Description: "Kubernetes Services"},
			{Name: "configmaps", Enabled: true, Description: "Kubernetes ConfigMaps"},
			{Name: "secrets", Enabled: true, Description: "Kubernetes Secrets"},
			{Name: "replicasets", Enabled: false, Description: "Kubernetes ReplicaSets"},
			{Name: "daemonsets", Enabled: false, Description: "Kubernetes DaemonSets"},
			{Name: "statefulsets", Enabled: false, Description: "Kubernetes StatefulSets"},
			{Name: "jobs", Enabled: false, Description: "Kubernetes Jobs"},
			{Name: "cronjobs", Enabled: false, Description: "Kubernetes CronJobs"},
			{Name: "persistentvolumes", Enabled: false, Description: "Kubernetes PersistentVolumes"},
			{Name: "persistentvolumeclaims", Enabled: false, Description: "Kubernetes PersistentVolumeClaims"},
			{Name: "ingresses", Enabled: false, Description: "Kubernetes Ingresses"},
			{Name: "networkpolicies", Enabled: false, Description: "Kubernetes NetworkPolicies"},
		},
	}

	// Check if config file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		// Create default config file
		if err := SaveConfig(configPath, defaultConfig); err != nil {
			return nil, fmt.Errorf("failed to create default config: %v", err)
		}
		fmt.Printf("Created default configuration file at %s\n", configPath)
		
		// Apply environment variable overrides to default config
		applyEnvironmentOverrides(defaultConfig)
		return defaultConfig, nil
	}

	// Load existing config file
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %v", err)
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %v", err)
	}

	// Apply environment variable overrides to loaded config
	applyEnvironmentOverrides(&config)

	return &config, nil
}

// applyEnvironmentOverrides applies environment variable overrides to the configuration
func applyEnvironmentOverrides(config *Config) {
	// Override persistence file path if environment variable is set
	if envFilePath := os.Getenv("PERSISTENCE_FILE_PATH"); envFilePath != "" {
		config.Persistence.FilePath = envFilePath
		fmt.Printf("Using persistence file path from environment: %s\n", envFilePath)
	}
	
	// Override web port if environment variable is set
	if envWebPort := os.Getenv("WEB_PORT"); envWebPort != "" {
		// Try to parse the port number
		if port := parsePort(envWebPort); port > 0 {
			config.WebPort = port
			fmt.Printf("Using web port from environment: %d\n", port)
		}
	}
}

// parsePort safely parses a port number from string
func parsePort(portStr string) int {
	var port int
	if _, err := fmt.Sscanf(portStr, "%d", &port); err == nil && port > 0 && port <= 65535 {
		return port
	}
	return 0
}

func SaveConfig(configPath string, config *Config) error {
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %v", err)
	}

	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %v", err)
	}

	return nil
}

func (c *Config) GetEnabledResources() []ResourceConfig {
	var enabled []ResourceConfig
	for _, resource := range c.Resources {
		if resource.Enabled {
			enabled = append(enabled, resource)
		}
	}
	return enabled
}

func (c *Config) IsResourceEnabled(resourceName string) bool {
	for _, resource := range c.Resources {
		if resource.Name == resourceName && resource.Enabled {
			return true
		}
	}
	return false
}
