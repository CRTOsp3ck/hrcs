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
            <div class="metric-icon" style="background: rgba(59, 130, 246, 0.1);">
              <i class="pi pi-users" style="color: #3b82f6;"></i>
            </div>
            <div class="metric-details">
              <p class="metric-label">Total Users</p>
              <p class="metric-value">{{ stats.totalUsers }}</p>
              <p class="metric-change positive">
                <i class="pi pi-arrow-up"></i>
                12% from last month
              </p>
            </div>
          </div>
        </template>
      </Card>

      <Card class="metric-card">
        <template #content>
          <div class="metric-content">
            <div class="metric-icon" style="background: rgba(34, 197, 94, 0.1);">
              <i class="pi pi-file" style="color: #22c55e;"></i>
            </div>
            <div class="metric-details">
              <p class="metric-label">Total Claims</p>
              <p class="metric-value">{{ stats.totalClaims }}</p>
              <p class="metric-change positive">
                <i class="pi pi-arrow-up"></i>
                8% from last month
              </p>
            </div>
          </div>
        </template>
      </Card>

      <Card class="metric-card">
        <template #content>
          <div class="metric-content">
            <div class="metric-icon" style="background: rgba(251, 146, 60, 0.1);">
              <i class="pi pi-clock" style="color: #fb923c;"></i>
            </div>
            <div class="metric-details">
              <p class="metric-label">Pending Approvals</p>
              <p class="metric-value">{{ stats.pendingClaims }}</p>
              <p class="metric-change negative">
                <i class="pi pi-arrow-down"></i>
                5% from last week
              </p>
            </div>
          </div>
        </template>
      </Card>

      <Card class="metric-card">
        <template #content>
          <div class="metric-content">
            <div class="metric-icon" style="background: rgba(168, 85, 247, 0.1);">
              <i class="pi pi-dollar" style="color: #a855f7;"></i>
            </div>
            <div class="metric-details">
              <p class="metric-label">Total Amount</p>
              <p class="metric-value">${{ formatAmount(stats.totalAmount) }}</p>
              <p class="metric-change positive">
                <i class="pi pi-arrow-up"></i>
                15% from last month
              </p>
            </div>
          </div>
        </template>
      </Card>
    </div>

    <div class="dashboard-grid">
      <!-- Claims by Status Chart -->
      <Card>
        <template #header>
          <h3 class="card-title">Claims by Status</h3>
        </template>
        <template #content>
          <Chart type="doughnut" :data="claimsByStatusData" :options="chartOptions" />
        </template>
      </Card>

      <!-- Claims by Type -->
      <Card>
        <template #header>
          <h3 class="card-title">Claims by Type</h3>
        </template>
        <template #content>
          <Chart type="bar" :data="claimsByTypeData" :options="barChartOptions" />
        </template>
      </Card>

      <!-- Recent Activity -->
      <Card class="recent-activity">
        <template #header>
          <h3 class="card-title">Recent Activity</h3>
        </template>
        <template #content>
          <Timeline :value="recentActivity" class="custom-timeline">
            <template #marker="slotProps">
              <span class="activity-marker" :class="getActivityClass(slotProps.item.type)">
                <i :class="getActivityIcon(slotProps.item.type)"></i>
              </span>
            </template>
            <template #content="slotProps">
              <div class="activity-content">
                <p class="activity-title">{{ slotProps.item.title }}</p>
                <p class="activity-time">{{ formatTimeAgo(slotProps.item.timestamp) }}</p>
              </div>
            </template>
          </Timeline>
        </template>
      </Card>
    </div>

    <!-- Quick Actions -->
    <Card class="quick-actions">
      <template #header>
        <h3 class="card-title">Quick Actions</h3>
      </template>
      <template #content>
        <div class="actions-grid">
          <Button 
            label="Manage Users" 
            icon="pi pi-users"
            @click="router.push('/admin/users')"
          />
          <Button 
            label="Configure Groups" 
            icon="pi pi-sitemap"
            severity="secondary"
            @click="router.push('/admin/groups')"
          />
          <Button 
            label="Claim Types" 
            icon="pi pi-tags"
            severity="info"
            @click="router.push('/admin/claim-types')"
          />
          <Button 
            label="Approval Workflow" 
            icon="pi pi-shield"
            severity="warning"
            @click="router.push('/admin/approval-levels')"
          />
        </div>
      </template>
    </Card>
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

const chartOptions = {
  responsive: true,
  maintainAspectRatio: false,
  plugins: {
    legend: {
      position: 'bottom'
    }
  }
}

const barChartOptions = {
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

.metric-label {
  color: var(--surface-600);
  font-size: 0.875rem;
  margin: 0;
}

.metric-value {
  font-size: 2rem;
  font-weight: 700;
  color: var(--surface-900);
  margin: 0.25rem 0;
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

.dashboard-grid {
  display: grid;
  grid-template-columns: 1fr 1fr 400px;
  gap: 1.5rem;
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

.quick-actions {
  margin-top: 2rem;
}

.actions-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 1rem;
}

@media (max-width: 1400px) {
  .dashboard-grid {
    grid-template-columns: 1fr 1fr;
  }
  
  .recent-activity {
    grid-column: span 2;
  }
}

@media (max-width: 768px) {
  .dashboard-grid {
    grid-template-columns: 1fr;
  }
  
  .metrics-grid {
    grid-template-columns: 1fr;
  }
}</style>