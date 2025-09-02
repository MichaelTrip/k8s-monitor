package monitor

import (
	"context"
	"crypto/rand"
	"fmt"
	"log"
	"sync"
	"time"

	v1 "k8s.io/api/core/v1"
	appsv1 "k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
	batchv1beta1 "k8s.io/api/batch/v1beta1"
	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	
	"k8s-monitor/pkg/config"
	"k8s-monitor/pkg/utils"
)

type Change struct {
	ID          string    `json:"id"`
	Timestamp   time.Time `json:"timestamp"`
	EventType   string    `json:"eventType"`
	ResourceType string   `json:"resourceType"`
	Namespace   string    `json:"namespace"`
	Name        string    `json:"name"`
	Details     string    `json:"details"`
	IsRead      bool      `json:"isRead"`
}

type K8sMonitor struct {
	clientset       *kubernetes.Clientset
	config          *config.Config
	changes         []Change
	changesMutex    sync.RWMutex
	startTime       time.Time
	stopChan        chan struct{}
	knownResources  map[string]map[string]string // resourceType -> namespace/name -> resourceVersion
	resourcesMutex  sync.RWMutex
}

func NewK8sMonitor(clientset *kubernetes.Clientset, cfg *config.Config) (*K8sMonitor, error) {
	monitor := &K8sMonitor{
		clientset:      clientset,
		config:         cfg,
		changes:        []Change{},
		startTime:      time.Now(),
		stopChan:       make(chan struct{}),
		knownResources: make(map[string]map[string]string),
	}

	// Initialize known resources map
	for _, resource := range cfg.Resources {
		if resource.Enabled {
			monitor.knownResources[resource.Name] = make(map[string]string)
		}
	}

	// Load existing changes from file if persistence is enabled
	if cfg.Persistence.Enabled {
		if loadedChanges, err := utils.LoadChangesFromFile(cfg.Persistence.FilePath); err == nil {
			// Convert loaded changes back to Change structs
			for _, changeData := range loadedChanges {
				if changeMap, ok := changeData.(map[string]interface{}); ok {
					change := Change{
						ID:           getString(changeMap, "id"),
						Timestamp:    getTime(changeMap, "timestamp"),
						EventType:    getString(changeMap, "eventType"),
						ResourceType: getString(changeMap, "resourceType"),
						Namespace:    getString(changeMap, "namespace"),
						Name:         getString(changeMap, "name"),
						Details:      getString(changeMap, "details"),
						IsRead:       getBool(changeMap, "isRead"),
					}
					monitor.changes = append(monitor.changes, change)
				}
			}
			log.Printf("Loaded %d changes from %s", len(monitor.changes), cfg.Persistence.FilePath)
		} else {
			log.Printf("Could not load changes from file: %v", err)
		}

		// Populate known resources from loaded changes to avoid duplicate ADDED events
		monitor.populateKnownResourcesFromChanges()
		
		// Create initial empty file if it doesn't exist
		monitor.saveToFile()
	}

	return monitor, nil
}

func (m *K8sMonitor) StartMonitoring() error {
	enabledResources := m.config.GetEnabledResources()
	
	// First, populate current state to avoid false ADDED events
	if err := m.populateCurrentState(); err != nil {
		log.Printf("Warning: Could not populate current state: %v", err)
	}

	for _, resource := range enabledResources {
		go m.startResourceWatcher(resource)
	}

	// Start auto-save goroutine if persistence is enabled
	if m.config.Persistence.Enabled && m.config.Persistence.AutoSave {
		go m.startAutoSave()
	}

	log.Printf("Started monitoring %d enabled Kubernetes resources...", len(enabledResources))
	return nil
}

func (m *K8sMonitor) startResourceWatcher(resource config.ResourceConfig) {
	log.Printf("Starting watcher for %s (namespace: %s)", resource.Name, resource.Namespace)
	
	for {
		var watcher watch.Interface
		var err error

		namespace := resource.Namespace
		if namespace == "" {
			namespace = metav1.NamespaceAll
		}

		listOptions := metav1.ListOptions{
			ResourceVersion: "0", // Start from current version to only get new changes
		}

		switch resource.Name {
		case "pods":
			watcher, err = m.clientset.CoreV1().Pods(namespace).Watch(context.TODO(), listOptions)
		case "deployments":
			watcher, err = m.clientset.AppsV1().Deployments(namespace).Watch(context.TODO(), listOptions)
		case "services":
			watcher, err = m.clientset.CoreV1().Services(namespace).Watch(context.TODO(), listOptions)
		case "configmaps":
			watcher, err = m.clientset.CoreV1().ConfigMaps(namespace).Watch(context.TODO(), listOptions)
		case "secrets":
			watcher, err = m.clientset.CoreV1().Secrets(namespace).Watch(context.TODO(), listOptions)
		case "replicasets":
			watcher, err = m.clientset.AppsV1().ReplicaSets(namespace).Watch(context.TODO(), listOptions)
		case "daemonsets":
			watcher, err = m.clientset.AppsV1().DaemonSets(namespace).Watch(context.TODO(), listOptions)
		case "statefulsets":
			watcher, err = m.clientset.AppsV1().StatefulSets(namespace).Watch(context.TODO(), listOptions)
		case "jobs":
			watcher, err = m.clientset.BatchV1().Jobs(namespace).Watch(context.TODO(), listOptions)
		case "cronjobs":
			watcher, err = m.clientset.BatchV1beta1().CronJobs(namespace).Watch(context.TODO(), listOptions)
		case "persistentvolumes":
			watcher, err = m.clientset.CoreV1().PersistentVolumes().Watch(context.TODO(), listOptions)
		case "persistentvolumeclaims":
			watcher, err = m.clientset.CoreV1().PersistentVolumeClaims(namespace).Watch(context.TODO(), listOptions)
		case "ingresses":
			watcher, err = m.clientset.NetworkingV1().Ingresses(namespace).Watch(context.TODO(), listOptions)
		case "networkpolicies":
			watcher, err = m.clientset.NetworkingV1().NetworkPolicies(namespace).Watch(context.TODO(), listOptions)
		default:
			log.Printf("Unknown resource type: %s", resource.Name)
			return
		}

		if err != nil {
			log.Printf("Error watching %s: %v", resource.Name, err)
			time.Sleep(5 * time.Second)
			continue
		}

		for event := range watcher.ResultChan() {
			m.handleEvent(resource.Name, event)
		}

		// If we reach here, the watcher closed, restart it
		log.Printf("Watcher for %s closed, restarting...", resource.Name)
		time.Sleep(1 * time.Second)
	}
}

func (m *K8sMonitor) handleEvent(resourceType string, event watch.Event) {
	if event.Object == nil {
		return
	}

	var namespace, name, details, resourceVersion string
	
	switch obj := event.Object.(type) {
	case *v1.Pod:
		namespace = obj.Namespace
		name = obj.Name
		resourceVersion = obj.ResourceVersion
		details = fmt.Sprintf("Phase: %s, Ready: %v", obj.Status.Phase, m.isPodReady(obj))
	case *appsv1.Deployment:
		namespace = obj.Namespace
		name = obj.Name
		resourceVersion = obj.ResourceVersion
		details = fmt.Sprintf("Replicas: %d/%d, Available: %d", 
			obj.Status.ReadyReplicas, obj.Status.Replicas, obj.Status.AvailableReplicas)
	case *appsv1.ReplicaSet:
		namespace = obj.Namespace
		name = obj.Name
		resourceVersion = obj.ResourceVersion
		details = fmt.Sprintf("Replicas: %d/%d", obj.Status.ReadyReplicas, obj.Status.Replicas)
	case *appsv1.DaemonSet:
		namespace = obj.Namespace
		name = obj.Name
		resourceVersion = obj.ResourceVersion
		details = fmt.Sprintf("Desired: %d, Ready: %d", obj.Status.DesiredNumberScheduled, obj.Status.NumberReady)
	case *appsv1.StatefulSet:
		namespace = obj.Namespace
		name = obj.Name
		resourceVersion = obj.ResourceVersion
		details = fmt.Sprintf("Replicas: %d/%d", obj.Status.ReadyReplicas, obj.Status.Replicas)
	case *v1.Service:
		namespace = obj.Namespace
		name = obj.Name
		resourceVersion = obj.ResourceVersion
		details = fmt.Sprintf("Type: %s, Ports: %d", obj.Spec.Type, len(obj.Spec.Ports))
	case *v1.ConfigMap:
		namespace = obj.Namespace
		name = obj.Name
		resourceVersion = obj.ResourceVersion
		details = fmt.Sprintf("Data keys: %d", len(obj.Data))
	case *v1.Secret:
		namespace = obj.Namespace
		name = obj.Name
		resourceVersion = obj.ResourceVersion
		details = fmt.Sprintf("Type: %s, Data keys: %d", obj.Type, len(obj.Data))
	case *batchv1.Job:
		namespace = obj.Namespace
		name = obj.Name
		resourceVersion = obj.ResourceVersion
		details = fmt.Sprintf("Active: %d, Succeeded: %d, Failed: %d", 
			obj.Status.Active, obj.Status.Succeeded, obj.Status.Failed)
	case *batchv1beta1.CronJob:
		namespace = obj.Namespace
		name = obj.Name
		resourceVersion = obj.ResourceVersion
		details = fmt.Sprintf("Schedule: %s, Suspend: %v", obj.Spec.Schedule, *obj.Spec.Suspend)
	case *v1.PersistentVolume:
		namespace = obj.Namespace
		name = obj.Name
		resourceVersion = obj.ResourceVersion
		details = fmt.Sprintf("Phase: %s, Capacity: %s", obj.Status.Phase, obj.Spec.Capacity.Storage().String())
	case *v1.PersistentVolumeClaim:
		namespace = obj.Namespace
		name = obj.Name
		resourceVersion = obj.ResourceVersion
		details = fmt.Sprintf("Phase: %s, Storage: %s", obj.Status.Phase, obj.Spec.Resources.Requests.Storage().String())
	case *networkingv1.Ingress:
		namespace = obj.Namespace
		name = obj.Name
		resourceVersion = obj.ResourceVersion
		details = fmt.Sprintf("Rules: %d", len(obj.Spec.Rules))
	case *networkingv1.NetworkPolicy:
		namespace = obj.Namespace
		name = obj.Name
		resourceVersion = obj.ResourceVersion
		details = fmt.Sprintf("Pod selector: %v", obj.Spec.PodSelector)
	default:
		// Handle unknown types
		if metaObj, ok := event.Object.(metav1.Object); ok {
			namespace = metaObj.GetNamespace()
			name = metaObj.GetName()
			resourceVersion = metaObj.GetResourceVersion()
			details = "Unknown resource type"
		} else {
			return
		}
	}

	resourceKey := fmt.Sprintf("%s/%s", namespace, name)
	
	// Check if this is a truly new resource or just a restart
	m.resourcesMutex.Lock()
	lastKnownVersion, existed := m.knownResources[resourceType][resourceKey]
	m.knownResources[resourceType][resourceKey] = resourceVersion
	m.resourcesMutex.Unlock()

	// Skip ADDED events for resources we already know about (from loaded state)
	if string(event.Type) == "ADDED" && existed && lastKnownVersion != "" {
		log.Printf("Skipping duplicate ADDED event for existing %s %s", resourceType, resourceKey)
		return
	}

	change := Change{
		ID:           generateID(),
		Timestamp:    time.Now(),
		EventType:    string(event.Type),
		ResourceType: resourceType,
		Namespace:    namespace,
		Name:         name,
		Details:      details,
		IsRead:       false,
	}

	m.changesMutex.Lock()
	m.changes = append(m.changes, change)
	// Keep only last 1000 changes to prevent memory issues
	if len(m.changes) > 1000 {
		m.changes = m.changes[len(m.changes)-1000:]
	}
	m.changesMutex.Unlock()

	log.Printf("Change detected: %s %s/%s in %s", 
		change.EventType, change.ResourceType, change.Name, change.Namespace)

	// Save to file immediately if persistence is enabled and not auto-saving
	if m.config.Persistence.Enabled && !m.config.Persistence.AutoSave {
		go m.saveToFile()
	}
}

func (m *K8sMonitor) isPodReady(pod *v1.Pod) bool {
	for _, condition := range pod.Status.Conditions {
		if condition.Type == v1.PodReady {
			return condition.Status == v1.ConditionTrue
		}
	}
	return false
}

func (m *K8sMonitor) GetChanges() []Change {
	m.changesMutex.RLock()
	defer m.changesMutex.RUnlock()
	
	// Return a copy to avoid race conditions
	changes := make([]Change, len(m.changes))
	copy(changes, m.changes)
	return changes
}

func (m *K8sMonitor) GetStats() map[string]interface{} {
	m.changesMutex.RLock()
	defer m.changesMutex.RUnlock()
	
	unreadCount := 0
	loadedFromFile := 0
	
	for _, change := range m.changes {
		if !change.IsRead {
			unreadCount++
		}
		// Count changes that were loaded from file (before current session)
		if change.Timestamp.Before(m.startTime) {
			loadedFromFile++
		}
	}
	
	stats := map[string]interface{}{
		"totalChanges":   len(m.changes),
		"unreadChanges":  unreadCount,
		"loadedFromFile": loadedFromFile,
		"currentSession": len(m.changes) - loadedFromFile,
		"startTime":      m.startTime,
		"uptime":         time.Since(m.startTime).String(),
	}
	
	// Count by event type
	eventCounts := make(map[string]int)
	resourceCounts := make(map[string]int)
	
	for _, change := range m.changes {
		eventCounts[change.EventType]++
		resourceCounts[change.ResourceType]++
	}
	
	stats["eventCounts"] = eventCounts
	stats["resourceCounts"] = resourceCounts
	
	return stats
}

func (m *K8sMonitor) MarkAllAsRead() int {
	m.changesMutex.Lock()
	defer m.changesMutex.Unlock()
	
	count := 0
	for i := range m.changes {
		if !m.changes[i].IsRead {
			m.changes[i].IsRead = true
			count++
		}
	}
	
	log.Printf("Marked %d changes as read", count)
	return count
}

func (m *K8sMonitor) MarkAsRead(changeID string) bool {
	m.changesMutex.Lock()
	defer m.changesMutex.Unlock()
	
	for i := range m.changes {
		if m.changes[i].ID == changeID {
			m.changes[i].IsRead = true
			return true
		}
	}
	
	return false
}

func generateID() string {
	bytes := make([]byte, 8)
	rand.Read(bytes)
	return fmt.Sprintf("%x", bytes)
}

func (m *K8sMonitor) GetConfig() *config.Config {
	return m.config
}

func (m *K8sMonitor) startAutoSave() {
	interval := time.Duration(m.config.Persistence.SaveInterval) * time.Second
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			m.saveToFile()
		case <-m.stopChan:
			return
		}
	}
}

func (m *K8sMonitor) saveToFile() {
	if !m.config.Persistence.Enabled {
		return
	}

	m.changesMutex.RLock()
	// Convert changes to interface{} slice for utils function
	changes := make([]interface{}, len(m.changes))
	for i, change := range m.changes {
		changes[i] = change
	}
	m.changesMutex.RUnlock()

	if err := utils.SaveChangesToFile(m.config.Persistence.FilePath, changes); err != nil {
		log.Printf("Error saving changes to file: %v", err)
	} else {
		log.Printf("Saved %d changes to %s", len(changes), m.config.Persistence.FilePath)
	}
}

func (m *K8sMonitor) SaveToFileNow() error {
	if !m.config.Persistence.Enabled {
		return fmt.Errorf("persistence is not enabled")
	}

	m.saveToFile()
	return nil
}

func (m *K8sMonitor) Stop() {
	close(m.stopChan)
	
	// Save changes one last time before stopping
	if m.config.Persistence.Enabled {
		m.saveToFile()
		log.Println("Final save completed")
	}
}

// Helper functions for loading changes from file
func getString(m map[string]interface{}, key string) string {
	if val, ok := m[key].(string); ok {
		return val
	}
	return ""
}

func getBool(m map[string]interface{}, key string) bool {
	if val, ok := m[key].(bool); ok {
		return val
	}
	return false
}

func getTime(m map[string]interface{}, key string) time.Time {
	if val, ok := m[key].(string); ok {
		if t, err := time.Parse(time.RFC3339, val); err == nil {
			return t
		}
	}
	return time.Now()
}

func (m *K8sMonitor) populateKnownResourcesFromChanges() {
	m.resourcesMutex.Lock()
	defer m.resourcesMutex.Unlock()
	
	// Track all resources that have been seen before from the loaded changes
	for _, change := range m.changes {
		resourceKey := fmt.Sprintf("%s/%s", change.Namespace, change.Name)
		if m.knownResources[change.ResourceType] == nil {
			m.knownResources[change.ResourceType] = make(map[string]string)
		}
		// Mark as known (with empty version since we don't have it from saved data)
		m.knownResources[change.ResourceType][resourceKey] = "loaded"
	}
	
	log.Printf("Populated known resources from %d loaded changes", len(m.changes))
}

func (m *K8sMonitor) populateCurrentState() error {
	enabledResources := m.config.GetEnabledResources()
	
	for _, resource := range enabledResources {
		namespace := resource.Namespace
		if namespace == "" {
			namespace = metav1.NamespaceAll
		}

		switch resource.Name {
		case "pods":
			if err := m.populatePods(namespace); err != nil {
				log.Printf("Error populating pods: %v", err)
			}
		case "deployments":
			if err := m.populateDeployments(namespace); err != nil {
				log.Printf("Error populating deployments: %v", err)
			}
		case "services":
			if err := m.populateServices(namespace); err != nil {
				log.Printf("Error populating services: %v", err)
			}
		// Add other resource types as needed
		}
	}
	
	return nil
}

func (m *K8sMonitor) populatePods(namespace string) error {
	pods, err := m.clientset.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return err
	}

	m.resourcesMutex.Lock()
	defer m.resourcesMutex.Unlock()
	
	for _, pod := range pods.Items {
		resourceKey := fmt.Sprintf("%s/%s", pod.Namespace, pod.Name)
		if m.knownResources["pods"] == nil {
			m.knownResources["pods"] = make(map[string]string)
		}
		m.knownResources["pods"][resourceKey] = pod.ResourceVersion
	}
	
	log.Printf("Populated %d existing pods", len(pods.Items))
	return nil
}

func (m *K8sMonitor) populateDeployments(namespace string) error {
	deployments, err := m.clientset.AppsV1().Deployments(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return err
	}

	m.resourcesMutex.Lock()
	defer m.resourcesMutex.Unlock()
	
	for _, deployment := range deployments.Items {
		resourceKey := fmt.Sprintf("%s/%s", deployment.Namespace, deployment.Name)
		if m.knownResources["deployments"] == nil {
			m.knownResources["deployments"] = make(map[string]string)
		}
		m.knownResources["deployments"][resourceKey] = deployment.ResourceVersion
	}
	
	log.Printf("Populated %d existing deployments", len(deployments.Items))
	return nil
}

func (m *K8sMonitor) populateServices(namespace string) error {
	services, err := m.clientset.CoreV1().Services(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return err
	}

	m.resourcesMutex.Lock()
	defer m.resourcesMutex.Unlock()
	
	for _, service := range services.Items {
		resourceKey := fmt.Sprintf("%s/%s", service.Namespace, service.Name)
		if m.knownResources["services"] == nil {
			m.knownResources["services"] = make(map[string]string)
		}
		m.knownResources["services"][resourceKey] = service.ResourceVersion
	}
	
	log.Printf("Populated %d existing services", len(services.Items))
	return nil
}