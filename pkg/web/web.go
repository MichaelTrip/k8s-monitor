package web

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"

	"k8s-monitor/pkg/monitor"
)

type WebServer struct {
	monitor *monitor.K8sMonitor
	port    int
}

func NewWebServer(monitor *monitor.K8sMonitor, port int) *WebServer {
	return &WebServer{
		monitor: monitor,
		port:    port,
	}
}

func (ws *WebServer) setupRoutes() {
	http.HandleFunc("/", ws.handleHome)
	http.HandleFunc("/api/changes", ws.handleAPIChanges)
	http.HandleFunc("/api/stats", ws.handleAPIStats)
	http.HandleFunc("/api/config", ws.handleAPIConfig)
	http.HandleFunc("/api/mark-read", ws.handleMarkRead)
	http.HandleFunc("/api/mark-all-read", ws.handleMarkAllRead)
	http.HandleFunc("/api/save-now", ws.handleSaveNow)
}

func (ws *WebServer) handleHome(w http.ResponseWriter, r *http.Request) {
	tmpl := `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Kubernetes Monitor</title>
    <style>
        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
            margin: 0;
            padding: 20px;
            background-color: #f5f5f5;
        }
        .container {
            max-width: 1200px;
            margin: 0 auto;
        }
        .header {
            background: white;
            padding: 20px;
            border-radius: 8px;
            box-shadow: 0 2px 4px rgba(0,0,0,0.1);
            margin-bottom: 20px;
        }
        .stats-grid {
            display: grid;
            grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
            gap: 20px;
            margin-bottom: 20px;
        }
        .stat-card {
            background: white;
            padding: 20px;
            border-radius: 8px;
            box-shadow: 0 2px 4px rgba(0,0,0,0.1);
            text-align: center;
        }
        .stat-card.unread {
            border-left: 4px solid #dc3545;
        }
        .stat-number {
            font-size: 2em;
            font-weight: bold;
            color: #007acc;
        }
        .stat-number.unread {
            color: #dc3545;
        }
        .changes-container {
            background: white;
            border-radius: 8px;
            box-shadow: 0 2px 4px rgba(0,0,0,0.1);
            overflow: hidden;
        }
        .changes-header {
            background: #007acc;
            color: white;
            padding: 15px 20px;
            display: flex;
            justify-content: space-between;
            align-items: center;
        }
        .changes-list {
            max-height: 600px;
            overflow-y: auto;
        }
        .change-item {
            padding: 15px 20px;
            border-bottom: 1px solid #eee;
            display: grid;
            grid-template-columns: 150px 100px 120px 150px 200px 1fr 80px;
            gap: 15px;
            align-items: center;
        }
        .change-item.unread {
            background-color: #fff3cd;
            border-left: 4px solid #ffc107;
        }
        .change-item:last-child {
            border-bottom: none;
        }
        .change-item:hover {
            background-color: #f8f9fa;
        }
        .timestamp {
            font-family: monospace;
            font-size: 0.9em;
            color: #666;
        }
        .event-type {
            padding: 4px 8px;
            border-radius: 4px;
            font-size: 0.85em;
            font-weight: bold;
        }
        .event-ADDED { background: #d4edda; color: #155724; }
        .event-MODIFIED { background: #fff3cd; color: #856404; }
        .event-DELETED { background: #f8d7da; color: #721c24; }
        .resource-type {
            background: #e9ecef;
            padding: 4px 8px;
            border-radius: 4px;
            font-size: 0.85em;
        }
        .namespace {
            color: #6f42c1;
            font-weight: 500;
        }
        .name {
            font-weight: bold;
        }
        .details {
            font-size: 0.9em;
            color: #666;
        }
        .loading {
            text-align: center;
            padding: 40px;
            color: #666;
        }
        .refresh-btn {
            background: #007acc;
            color: white;
            border: none;
            padding: 8px 16px;
            border-radius: 4px;
            cursor: pointer;
        }
        .refresh-btn:hover {
            background: #005a99;
        }
        .auto-refresh {
            display: flex;
            align-items: center;
            gap: 10px;
        }
        .controls {
            display: flex;
            align-items: center;
            gap: 15px;
        }
        .btn {
            background: #007acc;
            color: white;
            border: none;
            padding: 8px 16px;
            border-radius: 4px;
            cursor: pointer;
            font-size: 14px;
        }
        .btn:hover {
            background: #005a99;
        }
        .btn.btn-secondary {
            background: #6c757d;
        }
        .btn.btn-secondary:hover {
            background: #545b62;
        }
        .btn.btn-success {
            background: #28a745;
        }
        .btn.btn-success:hover {
            background: #1e7e34;
        }
        .mark-read-btn {
            background: #6c757d;
            color: white;
            border: none;
            padding: 4px 8px;
            border-radius: 3px;
            cursor: pointer;
            font-size: 12px;
        }
        .mark-read-btn:hover {
            background: #545b62;
        }
        .filter-controls {
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            padding: 20px;
            border-radius: 12px;
            margin-bottom: 20px;
            box-shadow: 0 4px 15px rgba(0,0,0,0.1);
        }
        .filter-header {
            color: white;
            font-size: 18px;
            font-weight: 600;
            margin-bottom: 15px;
            display: flex;
            align-items: center;
            gap: 8px;
        }
        .filters-grid {
            display: grid;
            grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
            gap: 20px;
        }
        .filter-section {
            background: rgba(255,255,255,0.95);
            border-radius: 8px;
            padding: 15px;
            backdrop-filter: blur(10px);
        }
        .filter-section h4 {
            margin: 0 0 10px 0;
            color: #333;
            font-size: 14px;
            font-weight: 600;
            text-transform: uppercase;
            letter-spacing: 0.5px;
        }
        .filter-options {
            display: flex;
            flex-wrap: wrap;
            gap: 8px;
        }
        .filter-chip {
            background: #f8f9fa;
            border: 2px solid #e9ecef;
            border-radius: 20px;
            padding: 6px 12px;
            font-size: 12px;
            font-weight: 500;
            cursor: pointer;
            transition: all 0.2s ease;
            user-select: none;
            position: relative;
        }
        .filter-chip:hover {
            background: #e9ecef;
            transform: translateY(-1px);
        }
        .filter-chip.selected {
            background: #007bff;
            border-color: #007bff;
            color: white;
        }
        .filter-chip.selected::after {
            content: "✓";
            position: absolute;
            right: 4px;
            top: 50%;
            transform: translateY(-50%);
            font-size: 10px;
        }
        .filter-chip.single-select {
            /* Special styling for single-select chips */
        }
        .classic-filter-group {
            display: flex;
            align-items: center;
            gap: 10px;
        }
        .classic-filter-group label {
            color: #333;
            font-weight: 500;
            font-size: 14px;
        }
        .classic-filter-group select {
            background: white;
            border: 2px solid #e9ecef;
            border-radius: 6px;
            padding: 8px 12px;
            font-size: 14px;
            min-width: 140px;
        }
        .filter-group {
            display: flex;
            align-items: center;
            gap: 5px;
        }
        select {
            padding: 6px 10px;
            border: 1px solid #ced4da;
            border-radius: 4px;
            font-size: 14px;
        }
        select[multiple] {
            min-height: 80px;
            width: 150px;
        }
        .header-grid {
            display: grid;
            grid-template-columns: 150px 100px 120px 150px 200px 1fr 80px;
            gap: 15px;
            font-weight: bold;
            padding: 15px 20px;
            background: #f8f9fa;
            border-bottom: 2px solid #dee2e6;
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>🚀 Kubernetes Monitor</h1>
            <p>Real-time monitoring of Kubernetes API objects</p>
            <div style="margin-top: 15px;">
                <button class="btn btn-secondary" onclick="toggleConfigPanel()">⚙️ Configuration</button>
                <button class="btn btn-success" onclick="saveToFile()">💾 Save to File</button>
            </div>
        </div>

        <div id="configPanel" style="display: none; background: white; padding: 20px; border-radius: 8px; box-shadow: 0 2px 4px rgba(0,0,0,0.1); margin-bottom: 20px;">
            <h3>📋 Monitored Resources Configuration</h3>
            <div id="resourceConfig" style="display: grid; grid-template-columns: repeat(auto-fit, minmax(300px, 1fr)); gap: 15px; margin-bottom: 20px;">
                <!-- Resource configuration will be loaded here -->
            </div>
            
            <h3>💾 Persistence Settings</h3>
            <div id="persistenceConfig" style="background: #f8f9fa; padding: 15px; border-radius: 8px;">
                <!-- Persistence configuration will be loaded here -->
            </div>
        </div>

        <div class="stats-grid" id="stats">
            <div class="stat-card">
                <div class="stat-number" id="totalChanges">-</div>
                <div>Total Changes</div>
            </div>
            <div class="stat-card unread">
                <div class="stat-number unread" id="unreadChanges">-</div>
                <div>Unread Changes</div>
            </div>
            <div class="stat-card">
                <div class="stat-number" id="uptime">-</div>
                <div>Uptime</div>
            </div>
            <div class="stat-card">
                <div class="stat-number" id="persistenceStatus">-</div>
                <div>Persistence</div>
            </div>
            <div class="stat-card">
                <div class="stat-number" id="currentSession">-</div>
                <div>This Session</div>
            </div>
        </div>

        <div class="changes-container">
            <div class="changes-header">
                <h2>📋 Recent Changes</h2>
                <div class="controls">
                    <button class="btn btn-success" onclick="markAllAsRead()">✅ Mark All Read</button>
                    <div class="auto-refresh">
                        <label>
                            <input type="checkbox" id="autoRefresh" checked> Auto-refresh
                        </label>
                        <button class="btn" onclick="loadChanges()">🔄 Refresh</button>
                    </div>
                </div>
            </div>
            <div class="filter-controls">
                <div class="filter-header">
                    🎛️ Filters & Sorting
                </div>
                <div class="filters-grid">
                    <div class="filter-section">
                        <h4>🔄 Sort Order</h4>
                        <div class="filter-options" id="sortOrderFilters">
                            <div class="filter-chip selected" data-value="newest">
                                <span>⬇️ Newest First</span>
                            </div>
                            <div class="filter-chip" data-value="oldest">
                                <span>⬆️ Oldest First</span>
                            </div>
                        </div>
                    </div>
                    <div class="filter-section">
                        <h4>👁️ Status Filter</h4>
                        <div class="filter-options" id="statusFilters">
                            <div class="filter-chip selected" data-value="all">
                                <span>📄 All Changes</span>
                            </div>
                            <div class="filter-chip" data-value="unread">
                                <span>🔴 Unread Only</span>
                            </div>
                            <div class="filter-chip" data-value="read">
                                <span>✅ Read Only</span>
                            </div>
                        </div>
                    </div>
                    <div class="filter-section">
                        <h4>📋 Event Types</h4>
                        <div class="filter-options" id="eventTypeFilters">
                            <div class="filter-chip" data-value="ADDED">
                                <span>➕ Added</span>
                            </div>
                            <div class="filter-chip" data-value="MODIFIED">
                                <span>✏️ Modified</span>
                            </div>
                            <div class="filter-chip" data-value="DELETED">
                                <span>🗑️ Deleted</span>
                            </div>
                        </div>
                    </div>
                    <div class="filter-section">
                        <h4>🎯 Resource Types</h4>
                        <div class="filter-options" id="resourceTypeFilters">
                            <!-- Will be populated dynamically -->
                        </div>
                    </div>
                </div>
            </div>
            <div class="header-grid">
                <div>Timestamp</div>
                <div>Event</div>
                <div>Resource</div>
                <div>Namespace</div>
                <div>Name</div>
                <div>Details</div>
                <div>Actions</div>
            </div>
            <div class="changes-list" id="changesList">
                <div class="loading">Loading changes...</div>
            </div>
        </div>
    </div>

    <script>
        let autoRefreshInterval;

        function formatTimestamp(timestamp) {
            const date = new Date(timestamp);
            return date.toLocaleTimeString();
        }

        function loadStats() {
            fetch('/api/stats')
                .then(response => response.json())
                .then(data => {
                    document.getElementById('totalChanges').textContent = data.totalChanges || 0;
                    document.getElementById('unreadChanges').textContent = data.unreadChanges || 0;
                    document.getElementById('uptime').textContent = data.uptime || '-';
                    document.getElementById('currentSession').textContent = data.currentSession || 0;
                })
                .catch(error => console.error('Error loading stats:', error));
            
            // Load persistence status
            fetch('/api/config')
                .then(response => response.json())
                .then(config => {
                    const persistenceEl = document.getElementById('persistenceStatus');
                    if (config.persistence && config.persistence.enabled) {
                        persistenceEl.textContent = config.persistence.autoSave ? '🟢 Auto' : '🟡 Manual';
                        persistenceEl.style.color = config.persistence.autoSave ? '#28a745' : '#ffc107';
                    } else {
                        persistenceEl.textContent = '🔴 Off';
                        persistenceEl.style.color = '#dc3545';
                    }
                })
                .catch(error => console.error('Error loading persistence status:', error));
        }

        function loadConfig() {
            fetch('/api/config')
                .then(response => response.json())
                .then(config => {
                    // Update resource filter chips
                    const resourceFilters = document.getElementById('resourceTypeFilters');
                    const currentSelected = Array.from(document.querySelectorAll('#resourceTypeFilters .filter-chip.selected'))
                        .map(chip => chip.dataset.value);
                    
                    // Clear existing chips
                    resourceFilters.innerHTML = '';
                    
                    // Add enabled resources as chips
                    config.resources.forEach(resource => {
                        if (resource.enabled) {
                            const chip = document.createElement('div');
                            chip.className = 'filter-chip';
                            chip.dataset.value = resource.name;
                            // Restore selection if it was previously selected
                            if (currentSelected.includes(resource.name)) {
                                chip.classList.add('selected');
                            }
                            
                            // Add appropriate icon based on resource type
                            let icon = '📦';
                            switch(resource.name) {
                                case 'pods': icon = '🎯'; break;
                                case 'services': icon = '🌐'; break;
                                case 'deployments': icon = '🚀'; break;
                                case 'configmaps': icon = '⚙️'; break;
                                case 'secrets': icon = '🔐'; break;
                                case 'ingresses': icon = '🌍'; break;
                                case 'persistentvolumes': icon = '💾'; break;
                                case 'persistentvolumeclaims': icon = '💽'; break;
                                case 'nodes': icon = '🖥️'; break;
                                case 'namespaces': icon = '📁'; break;
                                case 'replicasets': icon = '📊'; break;
                                case 'daemonsets': icon = '👹'; break;
                                case 'statefulsets': icon = '🏛️'; break;
                                case 'jobs': icon = '⚡'; break;
                                case 'cronjobs': icon = '⏰'; break;
                            }
                            
                            chip.innerHTML = '<span>' + icon + ' ' + resource.name.charAt(0).toUpperCase() + resource.name.slice(1) + '</span>';
                            resourceFilters.appendChild(chip);
                        }
                    });

                    // Update config panel
                    const configContainer = document.getElementById('resourceConfig');
                    configContainer.innerHTML = config.resources.map(resource => 
                        '<div style="border: 1px solid #dee2e6; padding: 15px; border-radius: 8px;">' +
                            '<div style="display: flex; justify-content: space-between; align-items: center; margin-bottom: 10px;">' +
                                '<strong>' + resource.name.charAt(0).toUpperCase() + resource.name.slice(1) + '</strong>' +
                                '<span style="color: ' + (resource.enabled ? '#28a745' : '#dc3545') + '; font-weight: bold;">' +
                                    (resource.enabled ? '✅ Enabled' : '❌ Disabled') +
                                '</span>' +
                            '</div>' +
                            '<div style="color: #666; font-size: 14px; margin-bottom: 8px;">' + resource.description + '</div>' +
                            (resource.namespace ? '<div style="font-size: 12px; color: #6f42c1;"><strong>Namespace:</strong> ' + resource.namespace + '</div>' : 
                             '<div style="font-size: 12px; color: #6f42c1;"><strong>Scope:</strong> All namespaces</div>') +
                        '</div>'
                    ).join('');

                    // Update persistence config panel
                    const persistenceContainer = document.getElementById('persistenceConfig');
                    persistenceContainer.innerHTML = 
                        '<div style="display: grid; grid-template-columns: repeat(auto-fit, minmax(250px, 1fr)); gap: 15px;">' +
                            '<div>' +
                                '<strong>Status:</strong> ' +
                                '<span style="color: ' + (config.persistence.enabled ? '#28a745' : '#dc3545') + ';">' +
                                    (config.persistence.enabled ? '✅ Enabled' : '❌ Disabled') +
                                '</span>' +
                            '</div>' +
                            '<div><strong>File Path:</strong> ' + config.persistence.filePath + '</div>' +
                            '<div><strong>Auto Save:</strong> ' + (config.persistence.autoSave ? 'Yes' : 'No') + '</div>' +
                            '<div><strong>Save Interval:</strong> ' + config.persistence.saveInterval + ' seconds</div>' +
                        '</div>';
                })
                .catch(error => console.error('Error loading config:', error));
        }

        function toggleConfigPanel() {
            const panel = document.getElementById('configPanel');
            if (panel.style.display === 'none') {
                panel.style.display = 'block';
                loadConfig();
            } else {
                panel.style.display = 'none';
            }
        }

        function loadChanges() {
            fetch('/api/changes')
                .then(response => response.json())
                .then(changes => {
                    const container = document.getElementById('changesList');
                    if (!changes || changes.length === 0) {
                        container.innerHTML = '<div class="loading">No changes detected yet...</div>';
                        return;
                    }

                    // Apply filters
                    let filteredChanges = applyFilters(changes);

                    // Sort changes
                    const sortChip = document.querySelector('#sortOrderFilters .filter-chip.selected');
                    const sortOrder = sortChip ? sortChip.dataset.value : 'newest';
                    filteredChanges.sort((a, b) => {
                        const dateA = new Date(a.timestamp);
                        const dateB = new Date(b.timestamp);
                        return sortOrder === 'newest' ? dateB - dateA : dateA - dateB;
                    });

                    if (filteredChanges.length === 0) {
                        container.innerHTML = '<div class="loading">No changes match the current filters...</div>';
                        return;
                    }

                    container.innerHTML = filteredChanges.slice(0, 100).map(change => 
                        '<div class="change-item' + (change.isRead ? '' : ' unread') + '">' +
                            '<div class="timestamp">' + formatTimestamp(change.timestamp) + '</div>' +
                            '<div class="event-type event-' + change.eventType + '">' + change.eventType + '</div>' +
                            '<div class="resource-type">' + change.resourceType + '</div>' +
                            '<div class="namespace">' + (change.namespace || 'default') + '</div>' +
                            '<div class="name">' + change.name + '</div>' +
                            '<div class="details">' + change.details + '</div>' +
                            '<div>' + (change.isRead ? '✓' : '<button class="mark-read-btn" onclick="markAsRead(\'' + change.id + '\')">Mark Read</button>') + '</div>' +
                        '</div>'
                    ).join('');
                })
                .catch(error => {
                    console.error('Error loading changes:', error);
                    document.getElementById('changesList').innerHTML = '<div class="loading">Error loading changes</div>';
                });
        }

        function applyFilters(changes) {
            // Get status filter from chips
            const statusChip = document.querySelector('#statusFilters .filter-chip.selected');
            const filterRead = statusChip ? statusChip.dataset.value : 'all';
            
            // Get selected event types from chips
            const selectedEventTypes = Array.from(document.querySelectorAll('#eventTypeFilters .filter-chip.selected'))
                .map(chip => chip.dataset.value);
            
            // Get selected resource types from chips
            const selectedResources = Array.from(document.querySelectorAll('#resourceTypeFilters .filter-chip.selected'))
                .map(chip => chip.dataset.value);

            return changes.filter(change => {
                // Filter by read status
                if (filterRead === 'unread' && change.isRead) return false;
                if (filterRead === 'read' && !change.isRead) return false;

                // Filter by event type (if any are selected)
                if (selectedEventTypes.length > 0 && !selectedEventTypes.includes(change.eventType)) return false;

                // Filter by resource type (if any are selected)
                if (selectedResources.length > 0 && !selectedResources.includes(change.resourceType)) return false;

                return true;
            });
        }

        function markAllAsRead() {
            fetch('/api/mark-all-read', { method: 'POST' })
                .then(response => response.json())
                .then(data => {
                    console.log('Marked', data.count, 'changes as read');
                    loadChanges();
                    loadStats();
                })
                .catch(error => console.error('Error marking all as read:', error));
        }

        function markAsRead(changeId) {
            fetch('/api/mark-read', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({ id: changeId })
            })
                .then(response => response.json())
                .then(data => {
                    if (data.success) {
                        loadChanges();
                        loadStats();
                    }
                })
                .catch(error => console.error('Error marking as read:', error));
        }

        function saveToFile() {
            fetch('/api/save-now', { method: 'POST' })
                .then(response => response.json())
                .then(data => {
                    if (data.success) {
                        alert('Changes saved to file successfully!');
                    } else {
                        alert('Error saving to file: ' + (data.error || 'Unknown error'));
                    }
                })
                .catch(error => {
                    console.error('Error saving to file:', error);
                    alert('Error saving to file');
                });
        }

        function toggleAutoRefresh() {
            const checkbox = document.getElementById('autoRefresh');
            if (checkbox.checked) {
                autoRefreshInterval = setInterval(() => {
                    loadChanges();
                    loadStats();
                }, 2000);
            } else {
                clearInterval(autoRefreshInterval);
            }
        }

        // Initialize chip functionality
        function initializeChips() {
            // Add click handlers for single-select chips (sort order and status)
            document.querySelectorAll('#sortOrderFilters .filter-chip, #statusFilters .filter-chip').forEach(chip => {
                chip.addEventListener('click', function() {
                    // For single-select, remove selected from siblings and add to clicked
                    const siblings = this.parentElement.querySelectorAll('.filter-chip');
                    siblings.forEach(sibling => sibling.classList.remove('selected'));
                    this.classList.add('selected');
                    loadChanges();
                });
            });
            
            // Add click handlers for multi-select chips (event types)
            document.querySelectorAll('#eventTypeFilters .filter-chip').forEach(chip => {
                chip.addEventListener('click', function() {
                    this.classList.toggle('selected');
                    loadChanges();
                });
            });
            
            // Add click handlers for resource type chips (will be added after loadConfig)
            document.addEventListener('click', function(e) {
                if (e.target.closest('#resourceTypeFilters .filter-chip')) {
                    e.target.closest('.filter-chip').classList.toggle('selected');
                    loadChanges();
                }
            });
        }

        // Initialize
        document.addEventListener('DOMContentLoaded', function() {
            loadChanges();
            loadStats();
            loadConfig();
            toggleAutoRefresh();
            initializeChips();
            
            document.getElementById('autoRefresh').addEventListener('change', toggleAutoRefresh);
            
            // No need for dropdown event listeners anymore - chips handle their own events
        });
    </script>
</body>
</html>
`

	t, _ := template.New("home").Parse(tmpl)
	t.Execute(w, nil)
}

func (ws *WebServer) handleAPIChanges(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	changes := ws.monitor.GetChanges()
	json.NewEncoder(w).Encode(changes)
}

func (ws *WebServer) handleAPIStats(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	stats := ws.monitor.GetStats()
	json.NewEncoder(w).Encode(stats)
}

func (ws *WebServer) handleStatic(w http.ResponseWriter, r *http.Request) {
	// Serve static files if needed
	http.NotFound(w, r)
}

func (ws *WebServer) handleMarkAllRead(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	count := ws.monitor.MarkAllAsRead()
	
	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"success": true,
		"count":   count,
	}
	json.NewEncoder(w).Encode(response)
}

func (ws *WebServer) handleMarkRead(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		ID string `json:"id"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	success := ws.monitor.MarkAsRead(req.ID)
	
	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"success": success,
	}
	json.NewEncoder(w).Encode(response)
}

func (ws *WebServer) handleAPIConfig(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	config := ws.monitor.GetConfig()
	json.NewEncoder(w).Encode(config)
}

func (ws *WebServer) handleSaveNow(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	err := ws.monitor.SaveToFileNow()
	
	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"success": err == nil,
	}
	
	if err != nil {
		response["error"] = err.Error()
	}
	
	json.NewEncoder(w).Encode(response)
}

func (ws *WebServer) Start() error {
	ws.setupRoutes()
	
	addr := fmt.Sprintf(":%d", ws.port)
	fmt.Printf("Starting web server on %s\n", addr)
	return http.ListenAndServe(addr, nil)
}
