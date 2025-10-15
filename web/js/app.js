// Main application JavaScript
class KubernetesMonitorApp {
    constructor() {
        this.autoRefreshInterval = null;
        this.init();
    }

    init() {
        this.setupEventListeners();
        this.loadInitialData();
        this.loadVersion();
        this.toggleAutoRefresh();
    }

    setupEventListeners() {
        // Auto-refresh checkbox
        document.getElementById('autoRefresh').addEventListener('change', () => {
            this.toggleAutoRefresh();
        });

        // Initialize filter chips
        this.initializeFilterChips();
    }

    initializeFilterChips() {
        // Single-select chips (sort order and status)
        document.querySelectorAll('#sortOrderFilters .filter-chip, #statusFilters .filter-chip').forEach(chip => {
            chip.addEventListener('click', function() {
                // Remove selected from siblings and add to clicked
                const siblings = this.parentElement.querySelectorAll('.filter-chip');
                siblings.forEach(sibling => sibling.classList.remove('selected'));
                this.classList.add('selected');
                window.app.loadChanges();
            });
        });
        
        // Multi-select chips (event types)
        document.querySelectorAll('#eventTypeFilters .filter-chip').forEach(chip => {
            chip.addEventListener('click', function() {
                this.classList.toggle('selected');
                window.app.loadChanges();
            });
        });
        
        // Resource type chips (will be added dynamically)
        document.addEventListener('click', function(e) {
            if (e.target.closest('#resourceTypeFilters .filter-chip')) {
                e.target.closest('.filter-chip').classList.toggle('selected');
                window.app.loadChanges();
            }
        });
    }

    async loadInitialData() {
        await this.loadConfig();
        await this.loadChanges();
        await this.loadStats();
    }

    async loadVersion() {
        try {
            const response = await fetch('/api/debug');
            const data = await response.json();
            const versionElement = document.getElementById('appVersion');
            if (versionElement) {
                versionElement.textContent = data.version || 'dev';
            }
        } catch (error) {
            console.error('Error loading version:', error);
            const versionElement = document.getElementById('appVersion');
            if (versionElement) {
                versionElement.textContent = 'dev';
            }
        }
    }

    async loadConfig() {
        try {
            const response = await fetch('/api/config');
            const config = await response.json();
            
            this.populateResourceFilters(config.resources);
            this.updateConfigPanel(config);
        } catch (error) {
            console.error('Error loading config:', error);
        }
    }

    populateResourceFilters(resources) {
        const resourceFilters = document.getElementById('resourceTypeFilters');
        const currentSelected = Array.from(document.querySelectorAll('#resourceTypeFilters .filter-chip.selected'))
            .map(chip => chip.dataset.value);
        
        // Clear existing chips
        resourceFilters.innerHTML = '';
        
        // Add enabled resources as chips
        resources.forEach(resource => {
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
                    case 'replicasets': icon = '📊'; break;
                    case 'daemonsets': icon = '👹'; break;
                    case 'statefulsets': icon = '🏛️'; break;
                    case 'jobs': icon = '⚡'; break;
                    case 'cronjobs': icon = '⏰'; break;
                    case 'networkpolicies': icon = '🔒'; break;
                }
                
                chip.innerHTML = '<span>' + icon + ' ' + resource.name.charAt(0).toUpperCase() + resource.name.slice(1) + '</span>';
                resourceFilters.appendChild(chip);
            }
        });
    }

    updateConfigPanel(config) {
        // Update resource config panel
        const configContainer = document.getElementById('resourceConfig');
        configContainer.innerHTML = config.resources.map(resource => 
            '<div style="border: 1px solid #dee2e6; padding: 15px; border-radius: 8px; background: white;">' +
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
    }

    async loadChanges() {
        try {
            const response = await fetch('/api/changes');
            const changes = await response.json();
            
            const container = document.getElementById('changesList');
            if (!changes || changes.length === 0) {
                container.innerHTML = this.getEmptyState();
                return;
            }

            // Apply filters
            let filteredChanges = this.applyFilters(changes);

            // Sort changes
            const sortChip = document.querySelector('#sortOrderFilters .filter-chip.selected');
            const sortOrder = sortChip ? sortChip.dataset.value : 'newest';
            filteredChanges.sort((a, b) => {
                const dateA = new Date(a.timestamp);
                const dateB = new Date(b.timestamp);
                return sortOrder === 'newest' ? dateB - dateA : dateA - dateB;
            });

            if (filteredChanges.length === 0) {
                container.innerHTML = this.getEmptyState('No changes match the current filters');
                return;
            }

            container.innerHTML = filteredChanges.slice(0, 100).map(change => this.createChangeRow(change)).join('');
        } catch (error) {
            console.error('Error loading changes:', error);
            document.getElementById('changesList').innerHTML = this.getEmptyState('Error loading changes');
        }
    }

    createChangeRow(change) {
        return `
            <div class="change-item${change.isRead ? '' : ' unread'}">
                <div class="timestamp">${this.formatTimestamp(change.timestamp)}</div>
                <div class="event-type event-${change.eventType}">${change.eventType}</div>
                <div class="resource-type">${change.resourceType}</div>
                <div class="namespace">${change.namespace || 'default'}</div>
                <div class="name">${change.name}</div>
                <div class="details">${change.details}</div>
                <div>${change.isRead ? '✓' : `<button class="mark-read-btn" onclick="markAsRead('${change.id}')">Mark Read</button>`}</div>
            </div>
        `;
    }

    applyFilters(changes) {
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

    async loadStats() {
        try {
            const response = await fetch('/api/stats');
            const stats = await response.json();
            
            document.getElementById('totalChanges').textContent = stats.totalChanges || 0;
            document.getElementById('unreadChanges').textContent = stats.unreadChanges || 0;
            document.getElementById('currentSession').textContent = stats.currentSession || 0;
            document.getElementById('uptime').textContent = this.formatUptime(stats.uptime) || '-';
        } catch (error) {
            console.error('Error loading stats:', error);
        }
    }

    formatTimestamp(timestamp) {
        const date = new Date(timestamp);
        return date.toLocaleTimeString();
    }

    formatUptime(uptime) {
        if (!uptime) return '-';
        // Parse uptime string and make it more readable
        const parts = uptime.match(/(\d+h)?(\d+m)?(\d+\.?\d*s)?/);
        if (!parts) return uptime;
        
        let result = '';
        if (parts[1]) result += parts[1] + ' ';
        if (parts[2]) result += parts[2] + ' ';
        if (parts[3] && !parts[1] && !parts[2]) result += Math.round(parseFloat(parts[3])) + 's';
        
        return result.trim() || uptime;
    }

    getEmptyState(message = 'No changes detected yet...') {
        return `
            <div class="empty-state">
                <div class="empty-state-icon">📋</div>
                <div class="empty-state-title">Ready to Monitor</div>
                <div class="empty-state-subtitle">${message}</div>
            </div>
        `;
    }

    toggleAutoRefresh() {
        const checkbox = document.getElementById('autoRefresh');
        if (checkbox.checked) {
            this.autoRefreshInterval = setInterval(() => {
                this.loadChanges();
                this.loadStats();
            }, 2000);
        } else {
            clearInterval(this.autoRefreshInterval);
        }
    }

    async markAllAsRead() {
        try {
            const response = await fetch('/api/mark-all-read', { method: 'POST' });
            const data = await response.json();
            console.log('Marked', data.count, 'changes as read');
            await this.loadChanges();
            await this.loadStats();
            this.showNotification(`Marked ${data.count} changes as read`, 'success');
        } catch (error) {
            console.error('Error marking all as read:', error);
            this.showNotification('Error marking changes as read', 'error');
        }
    }

    async markAsRead(changeId) {
        try {
            const response = await fetch('/api/mark-read', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({ id: changeId })
            });
            const data = await response.json();
            if (data.success) {
                await this.loadChanges();
                await this.loadStats();
            }
        } catch (error) {
            console.error('Error marking as read:', error);
        }
    }

    async saveToFile() {
        try {
            const response = await fetch('/api/save-now', { method: 'POST' });
            const data = await response.json();
            if (data.success) {
                this.showNotification('Changes saved to file successfully!', 'success');
            } else {
                this.showNotification('Error saving to file: ' + (data.error || 'Unknown error'), 'error');
            }
        } catch (error) {
            console.error('Error saving to file:', error);
            this.showNotification('Error saving to file', 'error');
        }
    }

    toggleConfigPanel() {
        const panel = document.getElementById('configPanel');
        if (panel.style.display === 'none') {
            panel.style.display = 'block';
        } else {
            panel.style.display = 'none';
        }
    }

    toggleTheme() {
        document.body.classList.toggle('dark-theme');
        localStorage.setItem('darkTheme', document.body.classList.contains('dark-theme'));
    }

    refreshData() {
        this.loadChanges();
        this.loadStats();
        this.showNotification('Data refreshed', 'info');
    }

    showNotification(message, type = 'info') {
        // Create notification element
        const notification = document.createElement('div');
        notification.style.cssText = `
            position: fixed;
            top: 20px;
            right: 20px;
            padding: 15px 20px;
            border-radius: 8px;
            color: white;
            font-weight: 500;
            z-index: 1000;
            animation: slideIn 0.3s ease-out;
        `;

        // Set background color based on type
        switch (type) {
            case 'success':
                notification.style.backgroundColor = '#27ae60';
                break;
            case 'error':
                notification.style.backgroundColor = '#e74c3c';
                break;
            case 'warning':
                notification.style.backgroundColor = '#f39c12';
                break;
            default:
                notification.style.backgroundColor = '#3498db';
        }

        notification.textContent = message;
        document.body.appendChild(notification);

        // Add animation styles if not already present
        if (!document.querySelector('#notification-styles')) {
            const styles = document.createElement('style');
            styles.id = 'notification-styles';
            styles.textContent = `
                @keyframes slideIn {
                    from { transform: translateX(100%); opacity: 0; }
                    to { transform: translateX(0); opacity: 1; }
                }
                @keyframes slideOut {
                    from { transform: translateX(0); opacity: 1; }
                    to { transform: translateX(100%); opacity: 0; }
                }
            `;
            document.head.appendChild(styles);
        }

        // Auto-remove after 3 seconds
        setTimeout(() => {
            notification.style.animation = 'slideOut 0.3s ease-out';
            setTimeout(() => {
                if (notification.parentNode) {
                    notification.parentNode.removeChild(notification);
                }
            }, 300);
        }, 3000);
    }
}

// Global functions for onclick handlers
function toggleConfigPanel() {
    window.app.toggleConfigPanel();
}

function markAllAsRead() {
    window.app.markAllAsRead();
}

function markAsRead(changeId) {
    window.app.markAsRead(changeId);
}

function saveToFile() {
    window.app.saveToFile();
}

function toggleTheme() {
    window.app.toggleTheme();
}

function refreshData() {
    window.app.refreshData();
}

// Initialize app when DOM is loaded
document.addEventListener('DOMContentLoaded', () => {
    // Prevent multiple app instances
    if (window.app) {
        console.warn('App already initialized, skipping duplicate initialization');
        return;
    }
    
    console.log('Initializing Kubernetes Monitor App');
    window.app = new KubernetesMonitorApp();

    // Load theme preference
    if (localStorage.getItem('darkTheme') === 'true') {
        document.body.classList.add('dark-theme');
    }
});

// Export for module systems
if (typeof module !== 'undefined' && module.exports) {
    module.exports = KubernetesMonitorApp;
}