<template>
  <div class="admin-page-container">
    <div class="page-header">
      <h1 class="page-title">Admin Dashboard</h1>
      <p class="page-subtitle">System overview and analytics</p>
    </div>

    <!-- Key Metrics -->
    <div class="metrics-grid">
      <Card class="metric-card">
        <template #content>
          <div class="metric-content">
            <div class="metric-top">
              <div class="metric-icon" style="background: rgba(59, 130, 246, 0.1);">
                <i class="pi pi-users" style="color: #3b82f6;"></i>
              </div>
              <div class="metric-details">
                <p class="metric-value">{{ stats.totalUsers }}</p>
                <p class="metric-label">Total Users</p>
              </div>
            </div>
            <p class="metric-change positive">
              <i class="pi pi-arrow-up"></i>
              12% from last month
            </p>
          </div>
        </template>
      </Card>

      <Card class="metric-card">
        <template #content>
          <div class="metric-content">
            <div class="metric-top">
              <div class="metric-icon" style="background: rgba(34, 197, 94, 0.1);">
                <i class="pi pi-file" style="color: #22c55e;"></i>
              </div>
              <div class="metric-details">
                <p class="metric-value">{{ stats.totalClaims }}</p>
                <p class="metric-label">Total Claims</p>
              </div>
            </div>
            <p class="metric-change positive">
              <i class="pi pi-arrow-up"></i>
              8% from last month
            </p>
          </div>
        </template>
      </Card>

      <Card class="metric-card">
        <template #content>
          <div class="metric-content">
            <div class="metric-top">
              <div class="metric-icon" style="background: rgba(251, 146, 60, 0.1);">
                <i class="pi pi-clock" style="color: #fb923c;"></i>
              </div>
              <div class="metric-details">
                <p class="metric-value">{{ stats.pendingClaims }}</p>
                <p class="metric-label">Pending Approvals</p>
              </div>
            </div>
            <p class="metric-change negative">
              <i class="pi pi-arrow-down"></i>
              5% from last week
            </p>
          </div>
        </template>
      </Card>

      <Card class="metric-card">
        <template #content>
          <div class="metric-content">
            <div class="metric-top">
              <div class="metric-icon" style="background: rgba(168, 85, 247, 0.1);">
                <i class="pi pi-dollar" style="color: #a855f7;"></i>
              </div>
              <div class="metric-details">
                <p class="metric-value">${{ formatAmount(stats.totalAmount) }}</p>
                <p class="metric-label">Total Amount</p>
              </div>
            </div>
            <p class="metric-change positive">
              <i class="pi pi-arrow-up"></i>
              15% from last month
            </p>
          </div>
        </template>
      </Card>
    </div>

    <div class="dashboard-grid-enhanced">
      <!-- Claims by Status Chart -->
      <Card class="chart-card">
        <template #header>
          <div class="chart-header">
            <h3 class="card-title">Claims by Status</h3>
            <div class="chart-stats">
              <span class="total-count">{{ getTotalClaimsCount() }} total</span>
            </div>
          </div>
        </template>
        <template #content>
          <Chart type="doughnut" :data="claimsByStatusData" :options="enhancedChartOptions" />
          <div class="chart-legend-details">
            <div v-for="item in stats.claimsByStatus" :key="item.status" class="legend-item">
              <div class="legend-color" :style="{ backgroundColor: getStatusColor(item.status) }"></div>
              <span class="legend-label">{{ formatStatus(item.status) }}</span>
              <span class="legend-count">{{ item.count }}</span>
              <span class="legend-percentage">{{ getStatusPercentage(item.count) }}%</span>
            </div>
          </div>
        </template>
      </Card>

      <!-- Claims by Type Chart -->
      <Card class="chart-card">
        <template #header>
          <div class="chart-header">
            <h3 class="card-title">Claims by Type</h3>
            <div class="chart-stats">
              <span class="highest-type">{{ getHighestClaimType() }}</span>
            </div>
          </div>
        </template>
        <template #content>
          <Chart type="bar" :data="claimsByTypeData" :options="enhancedBarChartOptions" />
          <div class="type-summary">
            <div v-for="item in stats.claimsByType" :key="item.type" class="type-item">
              <div class="type-info">
                <span class="type-name">{{ item.type }}</span>
                <span class="type-amount">${{ formatAmount(item.amount) }}</span>
              </div>
              <div class="type-count">{{ item.count }} claims</div>
            </div>
          </div>
        </template>
      </Card>

      <!-- Quick Actions -->
      <Card class="quick-actions-enhanced">
        <template #header>
          <h3 class="card-title">Quick Actions</h3>
        </template>
        <template #content>
          <div class="actions-grid-vertical">
            <Button 
              label="Manage Users" 
              icon="pi pi-users"
              @click="router.push('/admin/users')"
              class="action-button"
            />
            <Button 
              label="Configure Groups" 
              icon="pi pi-sitemap"
              severity="secondary"
              @click="router.push('/admin/groups')"
              class="action-button"
            />
            <Button 
              label="Claim Types" 
              icon="pi pi-tags"
              severity="info"
              @click="router.push('/admin/claim-types')"
              class="action-button"
            />
            <Button 
              label="Approval Workflow" 
              icon="pi pi-shield"
              severity="warning"
              @click="router.push('/admin/approval-levels')"
              class="action-button"
            />
          </div>
        </template>
      </Card>
    </div>

    <!-- Recent Claims Section -->
    <div class="dashboard-row">
      <Card class="recent-claims-admin">
        <template #header>
          <div class="section-header">
            <h3 class="card-title">Recent Claims</h3>
            <Button 
              label="View All Claims" 
              icon="pi pi-arrow-right"
              text
              @click="router.push('/admin/claims')"
            />
          </div>
        </template>
        <template #content>
          <DataTable
            :value="recentClaims"
            :loading="loading"
            responsiveLayout="scroll"
            :paginator="false"
            class="p-datatable-sm recent-claims-table"
          >
            <Column field="id" header="ID" style="width: 80px">
              <template #body="slotProps">
                #{{ slotProps.data.id }}
              </template>
            </Column>
            <Column field="title" header="Title" />
            <Column field="user.name" header="User" />
            <Column field="amount" header="Amount">
              <template #body="slotProps">
                ${{ formatAmount(slotProps.data.amount) }}
              </template>
            </Column>
            <Column field="status" header="Status">
              <template #body="slotProps">
                <Tag :value="formatStatus(slotProps.data.status)" :severity="getStatusSeverity(slotProps.data.status)" />
              </template>
            </Column>
            <Column field="created_at" header="Created">
              <template #body="slotProps">
                {{ formatTimeAgo(new Date(slotProps.data.created_at)) }}
              </template>
            </Column>
          </DataTable>
        </template>
      </Card>
    </div>

    <!-- Recent Activity (Audit Log) -->
    <div class="dashboard-row">
      <Card class="recent-activity-enhanced">
        <template #header>
          <div class="section-header">
            <h3 class="card-title">Recent Activity</h3>
            <Button 
              label="View Audit Log" 
              icon="pi pi-arrow-right"
              text
              @click="router.push('/admin/audit-log')"
            />
          </div>
        </template>
        <template #content>
          <Timeline :value="recentActivity" class="activity-timeline">
            <template #marker="slotProps">
              <span class="activity-marker" :class="getActivityClass(slotProps.item.type)">
                <i :class="getActivityIcon(slotProps.item.type)"></i>
              </span>
            </template>
            <template #content="slotProps">
              <div class="activity-content">
                <div class="activity-header">
                  <p class="activity-title">{{ slotProps.item.title }}</p>
                  <span class="activity-user">{{ slotProps.item.user || 'System' }}</span>
                </div>
                <p class="activity-time">{{ formatTimeAgo(slotProps.item.timestamp) }}</p>
                <p v-if="slotProps.item.details" class="activity-details">{{ slotProps.item.details }}</p>
              </div>
            </template>
          </Timeline>
        </template>
      </Card>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { dashboardApi } from '@/api'
import { useToast } from 'primevue/usetoast'

const router = useRouter()
const toast = useToast()

const loading = ref(false)
const stats = ref({
  totalUsers: 0,
  totalClaims: 0,
  pendingClaims: 0,
  totalAmount: 0,
  claimsByStatus: [] as any[],
  claimsByType: [] as any[]
})

const recentClaims = ref([])

const recentActivity = ref([
  {
    id: 1,
    type: 'claim_submitted',
    title: 'John Doe submitted a new claim for $1,200',
    timestamp: new Date(Date.now() - 1000 * 60 * 5)
  },
  {
    id: 2,
    type: 'claim_approved',
    title: 'Admin approved claim #123',
    timestamp: new Date(Date.now() - 1000 * 60 * 30)
  },
  {
    id: 3,
    type: 'user_added',
    title: 'New user Jane Smith registered',
    timestamp: new Date(Date.now() - 1000 * 60 * 60)
  },
  {
    id: 4,
    type: 'claim_rejected',
    title: 'Claim #456 was rejected',
    timestamp: new Date(Date.now() - 1000 * 60 * 120)
  }
])

const claimsByStatusData = computed(() => ({
  labels: stats.value.claimsByStatus.map(item => formatStatus(item.status)),
  datasets: [{
    data: stats.value.claimsByStatus.map(item => item.count),
    backgroundColor: [
      '#e2e8f0',
      '#3b82f6',
      '#22c55e',
      '#ef4444',
      '#fb923c',
      '#10b981'
    ]
  }]
}))

const claimsByTypeData = computed(() => ({
  labels: stats.value.claimsByType.map(item => item.type),
  datasets: [{
    label: 'Count',
    data: stats.value.claimsByType.map(item => item.count),
    backgroundColor: '#3b82f6'
  }]
}))

const enhancedChartOptions = {
  responsive: true,
  maintainAspectRatio: false,
  plugins: {
    legend: {
      display: false
    }
  }
}

const enhancedBarChartOptions = {
  responsive: true,
  maintainAspectRatio: false,
  plugins: {
    legend: {
      display: false
    }
  },
  scales: {
    y: {
      beginAtZero: true
    }
  }
}

const formatAmount = (amount: number) => {
  return new Intl.NumberFormat('en-US', {
    minimumFractionDigits: 0,
    maximumFractionDigits: 0
  }).format(amount)
}

const formatStatus = (status: string) => {
  return status.split('-').map(word => 
    word.charAt(0).toUpperCase() + word.slice(1)
  ).join(' ')
}

const formatTimeAgo = (date: Date) => {
  const seconds = Math.floor((new Date().getTime() - date.getTime()) / 1000)
  
  if (seconds < 60) return 'just now'
  if (seconds < 3600) return `${Math.floor(seconds / 60)} minutes ago`
  if (seconds < 86400) return `${Math.floor(seconds / 3600)} hours ago`
  return `${Math.floor(seconds / 86400)} days ago`
}

const getActivityClass = (type: string) => {
  const classes: Record<string, string> = {
    'claim_submitted': 'activity-primary',
    'claim_approved': 'activity-success',
    'claim_rejected': 'activity-danger',
    'user_added': 'activity-info'
  }
  return classes[type] || ''
}

const getActivityIcon = (type: string) => {
  const icons: Record<string, string> = {
    'claim_submitted': 'pi pi-send',
    'claim_approved': 'pi pi-check',
    'claim_rejected': 'pi pi-times',
    'user_added': 'pi pi-user-plus'
  }
  return icons[type] || 'pi pi-circle'
}

const getTotalClaimsCount = () => {
  return stats.value.claimsByStatus.reduce((total, item) => total + item.count, 0)
}

const getStatusColor = (status: string) => {
  const colors: Record<string, string> = {
    'draft': '#e2e8f0',
    'submitted': '#3b82f6',
    'approved': '#22c55e',
    'rejected': '#ef4444',
    'payment-in-progress': '#fb923c',
    'paid': '#10b981'
  }
  return colors[status] || '#e2e8f0'
}

const getStatusPercentage = (count: number) => {
  const total = getTotalClaimsCount()
  return total > 0 ? Math.round((count / total) * 100) : 0
}

const getHighestClaimType = () => {
  if (!stats.value.claimsByType.length) return 'N/A'
  const highest = stats.value.claimsByType.reduce((prev, current) => 
    current.count > prev.count ? current : prev
  )
  return `${highest.type} (${highest.count})`
}

const getStatusSeverity = (status: string) => {
  const severities: Record<string, string> = {
    'draft': 'secondary',
    'submitted': 'info',
    'approved': 'success',
    'rejected': 'danger',
    'payment-in-progress': 'warning',
    'paid': 'success'
  }
  return severities[status] || 'secondary'
}

const loadDashboard = async () => {
  loading.value = true
  try {
    const response = await dashboardApi.getAdminStats()
    if (response.data.data) {
      stats.value = response.data.data
    }
  } catch (error) {
    toast.add({
      severity: 'error',
      summary: 'Error',
      detail: 'Failed to load dashboard data',
      life: 3000
    })
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  loadDashboard()
})
</script>

<style scoped>
.metrics-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
  gap: 1.5rem;
  margin-bottom: 2rem;
}

.metric-card {
  position: relative;
  overflow: hidden;
}

.metric-content {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.metric-top {
  display: flex;
  align-items: center;
  gap: 1.5rem;
}

.metric-icon {
  width: 64px;
  height: 64px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.metric-icon i {
  font-size: 1.75rem;
}

.metric-details {
  flex: 1;
}

.metric-value {
  font-size: 2rem;
  font-weight: 700;
  color: var(--surface-900);
  margin: 0 0 0.25rem;
  line-height: 1;
}

.metric-label {
  color: var(--surface-600);
  font-size: 0.875rem;
  margin: 0 0 0.5rem;
}

.metric-change {
  font-size: 0.875rem;
  margin: 0;
  display: flex;
  align-items: center;
  gap: 0.25rem;
}

.metric-change.positive {
  color: #22c55e;
}

.metric-change.negative {
  color: #ef4444;
}

/* Enhanced Dashboard Layout */
.dashboard-grid-enhanced {
  display: grid;
  grid-template-columns: 1fr 1fr 300px;
  gap: 1.5rem;
  margin-bottom: 2rem;
}

.dashboard-row {
  margin-bottom: 2rem;
}

.card-title {
  font-size: 1.125rem;
  font-weight: 600;
  margin: 0;
}

.custom-timeline .activity-marker {
  width: 2rem;
  height: 2rem;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
}

.activity-primary {
  background: #3b82f6;
}

.activity-success {
  background: #22c55e;
}

.activity-danger {
  background: #ef4444;
}

.activity-info {
  background: #6366f1;
}

.activity-content {
  padding-bottom: 1.5rem;
}

.activity-title {
  margin: 0;
  color: var(--surface-900);
}

.activity-time {
  font-size: 0.875rem;
  color: var(--surface-500);
  margin: 0.25rem 0 0;
}

/* Chart Enhancements */
.chart-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.chart-stats {
  font-size: 0.875rem;
  color: var(--surface-600);
}

.chart-legend-details {
  margin-top: 1rem;
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.legend-item {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.25rem 0;
}

.legend-color {
  width: 12px;
  height: 12px;
  border-radius: 50%;
}

.legend-label {
  flex: 1;
  font-size: 0.875rem;
}

.legend-count {
  font-weight: 600;
  font-size: 0.875rem;
}

.legend-percentage {
  font-size: 0.75rem;
  color: var(--surface-500);
  min-width: 35px;
  text-align: right;
}

.type-summary {
  margin-top: 1rem;
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.type-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0.5rem;
  background: var(--surface-50);
  border-radius: 6px;
}

.type-info {
  display: flex;
  flex-direction: column;
  gap: 0.25rem;
}

.type-name {
  font-weight: 500;
  color: var(--surface-900);
}

.type-amount {
  font-size: 0.875rem;
  color: var(--surface-600);
}

.type-count {
  font-size: 0.875rem;
  color: var(--surface-500);
}

/* Quick Actions */
.actions-grid-vertical {
  display: flex;
  flex-direction: column;
  gap: 0.75rem;
}

.action-button {
  width: 100%;
  justify-content: flex-start;
}

/* Section Headers */
.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

/* Activity Timeline */
.activity-timeline .activity-marker {
  width: 1.75rem;
  height: 1.75rem;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  font-size: 0.875rem;
}

.activity-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 0.25rem;
}

.activity-user {
  font-size: 0.75rem;
  color: var(--surface-500);
  font-weight: 500;
}

.activity-details {
  font-size: 0.875rem;
  color: var(--surface-600);
  margin: 0.25rem 0 0;
}

@media (max-width: 1400px) {
  .dashboard-grid-enhanced {
    grid-template-columns: 1fr 1fr;
  }
  
  .quick-actions-enhanced {
    grid-column: span 2;
  }
}

@media (max-width: 768px) {
  .dashboard-grid-enhanced {
    grid-template-columns: 1fr;
  }
  
  .metrics-grid {
    grid-template-columns: 1fr;
  }
}</style>