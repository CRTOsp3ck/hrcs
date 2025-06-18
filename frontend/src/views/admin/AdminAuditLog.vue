<template>
  <div class="admin-page-container">
    <div class="page-header">
      <h1 class="page-title">Audit Log</h1>
      <p class="page-subtitle">System activity and change tracking</p>
    </div>

    <!-- Filters Section -->
    <div class="filters-section">
      <div class="filters-row">
        <div class="filter-group">
          <label>Date Range:</label>
          <Calendar
            v-model="filters.dateRange"
            selectionMode="range"
            dateFormat="mm/dd/yy"
            placeholder="Select date range"
            showIcon
            :maxDate="maxDate"
            class="date-filter"
          />
        </div>
        
        <div class="filter-group">
          <label>Action:</label>
          <Dropdown
            v-model="filters.action"
            :options="actionOptions"
            optionLabel="label"
            optionValue="value"
            placeholder="All Actions"
            showClear
            class="action-filter"
          />
        </div>
        
        <div class="filter-group">
          <label>Entity Type:</label>
          <Dropdown
            v-model="filters.entityType"
            :options="entityTypeOptions"
            optionLabel="label"
            optionValue="value"
            placeholder="All Types"
            showClear
            class="entity-filter"
          />
        </div>
        
        <div class="filter-group">
          <label>User:</label>
          <Dropdown
            v-model="filters.userId"
            :options="userOptions"
            optionLabel="name"
            optionValue="id"
            placeholder="All Users"
            showClear
            filter
            class="user-filter"
          />
        </div>

        <div class="filter-group">
          <label>Search:</label>
          <InputText 
            v-model="filters.search"
            placeholder="Search activities..."
            class="search-filter"
          />
        </div>

        <div class="filter-actions">
          <Button 
            label="Apply Filters" 
            icon="pi pi-filter"
            @click="applyFilters"
            :loading="loading"
          />
          <Button 
            label="Clear" 
            icon="pi pi-times"
            severity="secondary"
            @click="clearFilters"
          />
        </div>
      </div>
    </div>

    <!-- Quick Stats -->
    <div class="quick-stats" v-if="stats">
      <div class="stat-card">
        <div class="stat-icon">
          <i class="pi pi-calendar"></i>
        </div>
        <div class="stat-content">
          <div class="stat-value">{{ stats.total_activity }}</div>
          <div class="stat-label">Total Activities</div>
        </div>
      </div>
      
      <div class="stat-card">
        <div class="stat-icon create">
          <i class="pi pi-plus"></i>
        </div>
        <div class="stat-content">
          <div class="stat-value">{{ getActionCount('create') }}</div>
          <div class="stat-label">Create Actions</div>
        </div>
      </div>
      
      <div class="stat-card">
        <div class="stat-icon update">
          <i class="pi pi-pencil"></i>
        </div>
        <div class="stat-content">
          <div class="stat-value">{{ getActionCount('update') }}</div>
          <div class="stat-label">Update Actions</div>
        </div>
      </div>
      
      <div class="stat-card">
        <div class="stat-icon delete">
          <i class="pi pi-trash"></i>
        </div>
        <div class="stat-content">
          <div class="stat-value">{{ getActionCount('delete') }}</div>
          <div class="stat-label">Delete Actions</div>
        </div>
      </div>
    </div>

    <!-- Audit Log Table -->
    <DataTable
      :value="auditLogs"
      :loading="loading"
      paginator
      :rows="pagination.per_page"
      :rowsPerPageOptions="[25, 50, 100]"
      :totalRecords="pagination.total_items"
      lazy
      @page="onPageChange"
      responsiveLayout="scroll"
      class="audit-table"
      scrollable
      scrollHeight="600px"
    >
      <Column field="created_at" header="Timestamp" sortable style="width: 160px" frozen>
        <template #body="slotProps">
          <div class="timestamp-cell">
            <div class="date">{{ formatDate(slotProps.data.created_at) }}</div>
            <div class="time">{{ formatTime(slotProps.data.created_at) }}</div>
          </div>
        </template>
      </Column>
      
      <Column field="user_name" header="User" sortable style="width: 200px">
        <template #body="slotProps">
          <div class="user-info">
            <Avatar :label="getInitials(slotProps.data.user_name)" size="small" />
            <div class="user-details">
              <div class="user-name">{{ slotProps.data.user_name }}</div>
              <div class="user-email">{{ slotProps.data.user_email }}</div>
            </div>
          </div>
        </template>
      </Column>
      
      <Column field="action" header="Action" sortable style="width: 120px">
        <template #body="slotProps">
          <Tag 
            :value="slotProps.data.action" 
            :severity="getActionSeverity(slotProps.data.action)"
            :icon="getActionIcon(slotProps.data.action)"
          />
        </template>
      </Column>
      
      <Column field="entity_type" header="Entity" sortable style="width: 140px">
        <template #body="slotProps">
          <div class="entity-info">
            <i :class="getEntityIcon(slotProps.data.entity_type)" class="entity-icon"></i>
            <span>{{ formatEntityType(slotProps.data.entity_type) }}</span>
          </div>
        </template>
      </Column>
      
      <Column field="description" header="Description" style="min-width: 300px">
        <template #body="slotProps">
          <div class="action-description">
            {{ slotProps.data.description }}
            <span v-if="slotProps.data.entity_id" class="entity-id">(ID: {{ slotProps.data.entity_id }})</span>
          </div>
        </template>
      </Column>

      <Column field="severity" header="Priority" sortable style="width: 100px">
        <template #body="slotProps">
          <Tag 
            :value="slotProps.data.severity" 
            :severity="getSeverityLevel(slotProps.data.severity)"
          />
        </template>
      </Column>

      <Column field="ip_address" header="IP Address" style="width: 140px">
        <template #body="slotProps">
          <code class="ip-address">{{ slotProps.data.ip_address || 'N/A' }}</code>
        </template>
      </Column>
      
      <Column header="Actions" style="width: 100px" frozen frozenPosition="right">
        <template #body="slotProps">
          <div class="action-buttons">
            <Button
              icon="pi pi-eye"
              text
              rounded
              size="small"
              @click="viewDetails(slotProps.data)"
              v-tooltip="'View Details'"
            />
            <Button
              icon="pi pi-info-circle"
              text
              rounded
              size="small"
              @click="viewChanges(slotProps.data)"
              v-tooltip="'View Changes'"
              :disabled="!hasChanges(slotProps.data)"
            />
          </div>
        </template>
      </Column>
    </DataTable>

    <!-- Loading Skeleton -->
    <div v-show="loading && auditLogs.length === 0" class="table-loading">
      <div class="loading-skeleton">
        <div class="skeleton-row" v-for="i in 10" :key="i">
          <div class="skeleton-cell timestamp"></div>
          <div class="skeleton-cell user"></div>
          <div class="skeleton-cell action"></div>
          <div class="skeleton-cell entity"></div>
          <div class="skeleton-cell description"></div>
          <div class="skeleton-cell priority"></div>
          <div class="skeleton-cell ip"></div>
          <div class="skeleton-cell actions"></div>
        </div>
      </div>
    </div>

    <!-- Details Dialog -->
    <Dialog 
      v-model:visible="showDetailsDialog" 
      header="Audit Log Details" 
      :modal="true" 
      :style="{ width: '700px' }"
      class="audit-details-dialog"
    >
      <div v-if="selectedLog" class="audit-details">
        <div class="detail-section">
          <h4>Activity Information</h4>
          <div class="detail-grid">
            <div class="detail-item">
              <label>Timestamp:</label>
              <span>{{ formatDateTime(selectedLog.created_at) }}</span>
            </div>
            <div class="detail-item">
              <label>User:</label>
              <span>{{ selectedLog.user_name }} ({{ selectedLog.user_email }})</span>
            </div>
            <div class="detail-item">
              <label>Action:</label>
              <Tag :value="selectedLog.action" :severity="getActionSeverity(selectedLog.action)" />
            </div>
            <div class="detail-item">
              <label>Entity:</label>
              <span>{{ formatEntityType(selectedLog.entity_type) }}</span>
            </div>
            <div class="detail-item">
              <label>Entity ID:</label>
              <span>{{ selectedLog.entity_id || 'N/A' }}</span>
            </div>
            <div class="detail-item">
              <label>IP Address:</label>
              <code>{{ selectedLog.ip_address || 'N/A' }}</code>
            </div>
          </div>
        </div>

        <div class="detail-section" v-if="selectedLog.user_agent">
          <h4>Browser Information</h4>
          <code class="user-agent">{{ selectedLog.user_agent }}</code>
        </div>

        <div class="detail-section">
          <h4>Description</h4>
          <p class="description">{{ selectedLog.description }}</p>
        </div>
      </div>
    </Dialog>

    <!-- Changes Dialog -->
    <Dialog 
      v-model:visible="showChangesDialog" 
      header="Data Changes" 
      :modal="true" 
      :style="{ width: '900px' }"
      class="audit-changes-dialog"
    >
      <div v-if="selectedLog" class="changes-content">
        <div class="changes-grid">
          <div class="changes-section" v-if="selectedLog.old_values">
            <h4>Before Changes</h4>
            <pre class="json-display">{{ formatJSON(selectedLog.old_values) }}</pre>
          </div>
          
          <div class="changes-section" v-if="selectedLog.new_values">
            <h4>After Changes</h4>
            <pre class="json-display">{{ formatJSON(selectedLog.new_values) }}</pre>
          </div>
        </div>
        
        <div v-if="!selectedLog.old_values && !selectedLog.new_values" class="no-changes">
          <i class="pi pi-info-circle"></i>
          <p>No detailed change information available for this activity.</p>
        </div>
      </div>
    </Dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { useToast } from 'primevue/usetoast'
import api from '@/api'

// Toast
const toast = useToast()

// Loading states
const loading = ref(false)
const statsLoading = ref(false)

// Data
const auditLogs = ref([])
const stats = ref(null)
const actionOptions = ref([])
const entityTypeOptions = ref([])
const userOptions = ref([])

// Dialogs
const showDetailsDialog = ref(false)
const showChangesDialog = ref(false)
const selectedLog = ref(null)

// Filters
const maxDate = new Date()
const filters = reactive({
  dateRange: null,
  action: null,
  entityType: null,
  userId: null,
  search: ''
})

// Pagination
const pagination = reactive({
  current_page: 1,
  total_pages: 1,
  total_items: 0,
  per_page: 25,
  has_next: false,
  has_prev: false
})

// Computed
const getActionCount = computed(() => {
  return (actionType: string) => {
    if (!stats.value?.by_action) return 0
    return stats.value.by_action
      .filter((item: any) => getActionType(item.action) === actionType)
      .reduce((sum: number, item: any) => sum + item.count, 0)
  }
})

// Methods
const loadAuditLogs = async (page = 1) => {
  loading.value = true
  try {
    const params: any = {
      page,
      limit: pagination.per_page
    }

    // Add filters
    if (filters.action) params.action = filters.action
    if (filters.entityType) params.entity_type = filters.entityType
    if (filters.userId) params.user_id = filters.userId
    if (filters.search) params.search = filters.search
    
    if (filters.dateRange && filters.dateRange.length === 2) {
      params.start_date = formatDateForAPI(filters.dateRange[0])
      params.end_date = formatDateForAPI(filters.dateRange[1])
    }

    const response = await api.get('/admin/audit-log', { params })
    
    if (response.data.success) {
      auditLogs.value = response.data.data.data
      Object.assign(pagination, response.data.data.pagination)
    }
  } catch (error) {
    console.error('Failed to load audit logs:', error)
    toast.add({
      severity: 'error',
      summary: 'Error',
      detail: 'Failed to load audit logs',
      life: 3000
    })
  } finally {
    loading.value = false
  }
}

const loadStats = async () => {
  statsLoading.value = true
  try {
    const params: any = {}
    
    if (filters.dateRange && filters.dateRange.length === 2) {
      params.start_date = formatDateForAPI(filters.dateRange[0])
      params.end_date = formatDateForAPI(filters.dateRange[1])
    }

    const response = await api.get('/admin/audit-log/stats', { params })
    
    if (response.data.success) {
      stats.value = response.data.data
    }
  } catch (error) {
    console.error('Failed to load audit stats:', error)
  } finally {
    statsLoading.value = false
  }
}

const loadOptions = async () => {
  try {
    const response = await api.get('/admin/audit-log/options')
    
    if (response.data.success) {
      const data = response.data.data
      
      actionOptions.value = data.actions.map((action: string) => ({
        label: formatAction(action),
        value: action
      }))
      
      entityTypeOptions.value = data.entity_types.map((type: string) => ({
        label: formatEntityType(type),
        value: type
      }))
      
      userOptions.value = data.users
    }
  } catch (error) {
    console.error('Failed to load options:', error)
  }
}

const applyFilters = () => {
  pagination.current_page = 1
  loadAuditLogs(1)
  loadStats()
}

const clearFilters = () => {
  filters.dateRange = null
  filters.action = null
  filters.entityType = null
  filters.userId = null
  filters.search = ''
  applyFilters()
}

const onPageChange = (event: any) => {
  pagination.current_page = event.page + 1
  loadAuditLogs(pagination.current_page)
}

const viewDetails = (log: any) => {
  selectedLog.value = log
  showDetailsDialog.value = true
}

const viewChanges = (log: any) => {
  selectedLog.value = log
  showChangesDialog.value = true
}

const hasChanges = (log: any) => {
  return log.old_values || log.new_values
}

// Formatting methods
const formatDate = (dateString: string) => {
  return new Date(dateString).toLocaleDateString()
}

const formatTime = (dateString: string) => {
  return new Date(dateString).toLocaleTimeString()
}

const formatDateTime = (dateString: string) => {
  return new Date(dateString).toLocaleString()
}

const formatDateForAPI = (date: Date) => {
  return date.toISOString().split('T')[0]
}

const formatAction = (action: string) => {
  return action.charAt(0).toUpperCase() + action.slice(1)
}

const formatEntityType = (entityType: string) => {
  return entityType
    .split('_')
    .map(word => word.charAt(0).toUpperCase() + word.slice(1))
    .join(' ')
}

const formatJSON = (jsonString: string) => {
  try {
    return JSON.stringify(JSON.parse(jsonString), null, 2)
  } catch {
    return jsonString
  }
}

const getInitials = (name: string) => {
  return name
    .split(' ')
    .map(word => word.charAt(0))
    .join('')
    .toUpperCase()
    .slice(0, 2)
}

const getActionSeverity = (action: string) => {
  switch (action) {
    case 'delete': return 'danger'
    case 'create': return 'success'
    case 'update': case 'approve': case 'reject': return 'warning'
    case 'login': case 'logout': case 'view': return 'info'
    default: return 'secondary'
  }
}

const getActionIcon = (action: string) => {
  switch (action) {
    case 'create': return 'pi pi-plus'
    case 'update': return 'pi pi-pencil'
    case 'delete': return 'pi pi-trash'
    case 'approve': return 'pi pi-check'
    case 'reject': return 'pi pi-times'
    case 'submit': return 'pi pi-send'
    case 'login': return 'pi pi-sign-in'
    case 'logout': return 'pi pi-sign-out'
    case 'view': return 'pi pi-eye'
    case 'export': return 'pi pi-download'
    default: return 'pi pi-info-circle'
  }
}

const getEntityIcon = (entityType: string) => {
  switch (entityType) {
    case 'claim': return 'pi pi-file'
    case 'user': return 'pi pi-user'
    case 'user_group': return 'pi pi-users'
    case 'claim_type': return 'pi pi-tags'
    case 'approval_level': return 'pi pi-shield'
    case 'system': return 'pi pi-cog'
    default: return 'pi pi-box'
  }
}

const getSeverityLevel = (severity: string) => {
  switch (severity) {
    case 'high': return 'danger'
    case 'medium': return 'warning'
    case 'low': return 'info'
    default: return 'secondary'
  }
}

const getActionType = (action: string) => {
  switch (action) {
    case 'create': case 'import': return 'create'
    case 'update': case 'approve': case 'reject': case 'submit': return 'update'
    case 'delete': return 'delete'
    case 'login': case 'logout': return 'auth'
    case 'view': case 'export': return 'read'
    default: return 'other'
  }
}

// Lifecycle
onMounted(() => {
  loadOptions()
  loadAuditLogs()
  loadStats()
})
</script>

<style scoped>

.page-header {
  margin-bottom: 2rem;
}

.page-title {
  font-size: 2rem;
  font-weight: 700;
  color: var(--surface-900);
  margin: 0 0 0.5rem;
}

.page-subtitle {
  color: var(--surface-600);
  margin: 0;
}

/* Filters Section */
.filters-section {
  background: var(--surface-0);
  border: 1px solid var(--surface-200);
  border-radius: var(--border-radius);
  padding: 1.5rem;
  margin-bottom: 1.5rem;
}

.filters-row {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 1rem;
  align-items: end;
}

.filter-group {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.filter-group label {
  font-weight: 600;
  color: var(--surface-700);
  font-size: 0.875rem;
}

.filter-actions {
  display: flex;
  gap: 0.5rem;
  grid-column: span 2;
}

/* Quick Stats */
.quick-stats {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 1rem;
  margin-bottom: 1.5rem;
}

.stat-card {
  background: var(--surface-0);
  border: 1px solid var(--surface-200);
  border-radius: var(--border-radius);
  padding: 1.5rem;
  display: flex;
  align-items: center;
  gap: 1rem;
}

.stat-icon {
  width: 3rem;
  height: 3rem;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--primary-100);
  color: var(--primary-600);
  font-size: 1.25rem;
}

.stat-icon.create {
  background: var(--green-100);
  color: var(--green-600);
}

.stat-icon.update {
  background: var(--orange-100);
  color: var(--orange-600);
}

.stat-icon.delete {
  background: var(--red-100);
  color: var(--red-600);
}

.stat-value {
  font-size: 1.5rem;
  font-weight: 700;
  color: var(--surface-900);
}

.stat-label {
  color: var(--surface-600);
  font-size: 0.875rem;
}

/* Table Styles */
.audit-table {
  background: var(--surface-0);
  border-radius: var(--border-radius);
  overflow: hidden;
}

.timestamp-cell {
  display: flex;
  flex-direction: column;
  gap: 0.25rem;
}

.timestamp-cell .date {
  font-weight: 600;
  color: var(--surface-900);
}

.timestamp-cell .time {
  font-size: 0.75rem;
  color: var(--surface-500);
}

.user-info {
  display: flex;
  align-items: center;
  gap: 0.75rem;
}

.user-details {
  display: flex;
  flex-direction: column;
  gap: 0.25rem;
}

.user-name {
  font-weight: 600;
  color: var(--surface-900);
}

.user-email {
  font-size: 0.75rem;
  color: var(--surface-500);
}

.entity-info {
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.entity-icon {
  color: var(--surface-500);
}

.action-description {
  color: var(--surface-700);
}

.entity-id {
  font-size: 0.75rem;
  color: var(--surface-500);
  font-style: italic;
}

.ip-address {
  font-family: monospace;
  font-size: 0.875rem;
  background: var(--surface-100);
  padding: 0.25rem 0.5rem;
  border-radius: 0.25rem;
}

.action-buttons {
  display: flex;
  gap: 0.25rem;
}

/* Loading Skeleton */
.table-loading {
  padding: 1rem;
}

.loading-skeleton {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.skeleton-row {
  display: grid;
  grid-template-columns: 160px 200px 120px 140px 1fr 100px 140px 100px;
  gap: 1rem;
  align-items: center;
}

.skeleton-cell {
  height: 1.25rem;
  background: var(--surface-200);
  border-radius: 0.25rem;
  animation: pulse 1.5s ease-in-out infinite;
}

@keyframes pulse {
  0%, 100% { opacity: 1; }
  50% { opacity: 0.5; }
}

/* Dialog Styles */
.audit-details {
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
}

.detail-section h4 {
  margin: 0 0 1rem;
  color: var(--surface-900);
  font-weight: 600;
}

.detail-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 1rem;
}

.detail-item {
  display: flex;
  flex-direction: column;
  gap: 0.25rem;
}

.detail-item label {
  font-weight: 600;
  color: var(--surface-700);
  font-size: 0.875rem;
}

.user-agent {
  background: var(--surface-100);
  padding: 0.75rem;
  border-radius: 0.25rem;
  font-size: 0.875rem;
  word-break: break-all;
  white-space: pre-wrap;
}

.description {
  color: var(--surface-700);
  line-height: 1.5;
}

.changes-content {
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
}

.changes-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 1.5rem;
}

.changes-section h4 {
  margin: 0 0 0.75rem;
  color: var(--surface-900);
  font-weight: 600;
}

.json-display {
  background: var(--surface-100);
  border: 1px solid var(--surface-200);
  border-radius: 0.25rem;
  padding: 1rem;
  font-family: monospace;
  font-size: 0.875rem;
  overflow-x: auto;
  max-height: 300px;
  overflow-y: auto;
}

.no-changes {
  text-align: center;
  padding: 2rem;
  color: var(--surface-500);
}

.no-changes i {
  font-size: 2rem;
  margin-bottom: 1rem;
  display: block;
}

/* Responsive Design */
@media (max-width: 1200px) {
  .filters-row {
    grid-template-columns: repeat(auto-fit, minmax(180px, 1fr));
  }
  
  .detail-grid {
    grid-template-columns: 1fr;
  }
  
  .changes-grid {
    grid-template-columns: 1fr;
  }
}

@media (max-width: 768px) {
  .quick-stats {
    grid-template-columns: 1fr;
  }
  
  .filters-row {
    grid-template-columns: 1fr;
  }
  
  .filter-actions {
    grid-column: span 1;
  }
}
</style>