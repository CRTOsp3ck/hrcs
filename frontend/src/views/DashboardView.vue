<template>
  <div>
    <div class="page-container">
      <div class="page-header">
        <h1 class="page-title">Welcome back, {{ authStore.user?.first_name }} {{ authStore.user?.last_name }}!</h1>
        <p class="page-subtitle">Here's an overview of your claims activity</p>
      </div>

      <div class="dashboard-grid">
        <!-- Stats Cards -->
        <Card class="stat-card">
          <template #content>
            <div class="stat-content">
              <div class="stat-icon-container" style="background: rgba(59, 130, 246, 0.1);">
                <i class="pi pi-file stat-icon" style="color: var(--primary-500);"></i>
              </div>
              <div class="stat-details">
                <p class="stat-label">Total Claims</p>
                <p class="stat-value">{{ stats.totalClaims }}</p>
              </div>
            </div>
          </template>
        </Card>

        <Card class="stat-card">
          <template #content>
            <div class="stat-content">
              <div class="stat-icon-container" style="background: rgba(251, 191, 36, 0.1);">
                <i class="pi pi-clock stat-icon" style="color: #fbbf24;"></i>
              </div>
              <div class="stat-details">
                <p class="stat-label">Pending Claims</p>
                <p class="stat-value">{{ stats.pendingClaims }}</p>
              </div>
            </div>
          </template>
        </Card>

        <Card class="stat-card">
          <template #content>
            <div class="stat-content">
              <div class="stat-icon-container" style="background: rgba(16, 185, 129, 0.1);">
                <i class="pi pi-check-circle stat-icon" style="color: #10b981;"></i>
              </div>
              <div class="stat-details">
                <p class="stat-label">Approved Claims</p>
                <p class="stat-value">{{ stats.approvedClaims }}</p>
              </div>
            </div>
          </template>
        </Card>

        <Card class="stat-card">
          <template #content>
            <div class="stat-content">
              <div class="stat-icon-container" style="background: rgba(99, 102, 241, 0.1);">
                <i class="pi pi-dollar stat-icon" style="color: #6366f1;"></i>
              </div>
              <div class="stat-details">
                <p class="stat-label">Total Amount</p>
                <p class="stat-value">${{ formatAmount(stats.totalAmount) }}</p>
              </div>
            </div>
          </template>
        </Card>
      </div>

      <div class="dashboard-content">
        <!-- Recent Claims -->
        <Card class="recent-claims-card">
          <template #header>
            <div class="card-header">
              <h2 class="card-title">Recent Claims</h2>
              <Button
                label="View All"
                icon="pi pi-arrow-right"
                text
                @click="router.push('/claims')"
              />
            </div>
          </template>
          <template #content>
            <DataTable
              :value="stats.recentClaims"
              :loading="loading"
              responsiveLayout="scroll"
              :paginator="false"
              class="p-datatable-sm"
            >
              <template #empty>
                <div class="empty-state">
                  <i class="pi pi-inbox empty-icon"></i>
                  <p>No claims found</p>
                  <Button
                    label="Create Your First Claim"
                    icon="pi pi-plus"
                    @click="router.push('/claims/new')"
                  />
                </div>
              </template>

              <Column field="title" header="Title" :sortable="true">
                <template #body="slotProps">
                  <router-link :to="`/claims/${slotProps.data.id}`" class="claim-link">
                    {{ slotProps.data.title }}
                  </router-link>
                </template>
              </Column>

              <Column field="claim_type.name" header="Type" :sortable="true">
                <template #body="slotProps">
                  <Tag :value="slotProps.data.claim_type?.name" severity="info" />
                </template>
              </Column>

              <Column field="amount" header="Amount" :sortable="true">
                <template #body="slotProps">
                  <span class="font-semibold">${{ formatAmount(slotProps.data.amount) }}</span>
                </template>
              </Column>

              <Column field="status" header="Status" :sortable="true">
                <template #body="slotProps">
                  <span :class="`status-badge status-${slotProps.data.status}`">
                    <i :class="getStatusIcon(slotProps.data.status)"></i>
                    {{ formatStatus(slotProps.data.status) }}
                  </span>
                </template>
              </Column>

              <Column field="created_at" header="Date" :sortable="true">
                <template #body="slotProps">
                  {{ formatDate(slotProps.data.created_at) }}
                </template>
              </Column>

              <Column header="Actions" :exportable="false" style="width: 100px">
                <template #body="slotProps">
                  <Button
                    icon="pi pi-eye"
                    text
                    rounded
                    @click="router.push(`/claims/${slotProps.data.id}`)"
                  />
                </template>
              </Column>
            </DataTable>
          </template>
        </Card>

        <!-- Quick Actions -->
        <Card class="quick-actions-card">
          <template #header>
            <h2 class="card-title">Quick Actions</h2>
          </template>
          <template #content>
            <div class="quick-actions-grid">
              <Button
                label="New Claim"
                icon="pi pi-plus"
                class="p-button-lg"
                @click="router.push('/claims/new')"
              />
              <Button
                label="View All Claims"
                icon="pi pi-list"
                severity="secondary"
                class="p-button-lg"
                @click="router.push('/claims')"
              />
              <Button
                label="Profile Settings"
                icon="pi pi-user"
                severity="help"
                class="p-button-lg"
                @click="showProfileDialog = true"
              />
            </div>

            <Panel header="Claim Guidelines" :toggleable="true" class="mt-4">
              <ul class="guidelines-list">
                <li>Submit claims within 30 days of expense</li>
                <li>Attach all required receipts and documentation</li>
                <li>Ensure claim amounts are accurate</li>
                <li>Select the appropriate claim type</li>
                <li>Provide clear descriptions for faster approval</li>
              </ul>
            </Panel>
          </template>
        </Card>
      </div>
    </div>

    <!-- Profile Dialog -->
    <Dialog
      v-model:visible="showProfileDialog"
      modal
      header="Profile Settings"
      :style="{ width: '450px' }"
    >
      <div class="profile-content">
        <div class="profile-avatar">
          <Avatar
            :label="userInitials"
            :style="{ backgroundColor: '#2563eb', color: '#ffffff' }"
            size="xlarge"
            shape="circle"
          />
        </div>
        <div class="profile-info">
          <p><strong>Name:</strong> {{ authStore.user?.name }}</p>
          <p><strong>Email:</strong> {{ authStore.user?.email }}</p>
          <p><strong>Role:</strong> <Tag :value="authStore.user?.role" :severity="authStore.isAdmin ? 'success' : 'info'" /></p>
          <p v-if="authStore.user?.user_group"><strong>Group:</strong> {{ authStore.user.user_group.name }}</p>
          <p><strong>Member Since:</strong> {{ formatDate(authStore.user?.created_at) }}</p>
        </div>
      </div>
    </Dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { dashboardApi } from '@/api'
import { useToast } from 'primevue/usetoast'
import type { DashboardStats } from '@/types'

const router = useRouter()
const authStore = useAuthStore()
const toast = useToast()

const loading = ref(false)
const showProfileDialog = ref(false)
const stats = ref<DashboardStats>({
  totalClaims: 0,
  pendingClaims: 0,
  approvedClaims: 0,
  rejectedClaims: 0,
  totalAmount: 0,
  approvedAmount: 0,
  recentClaims: [],
  claimsByStatus: [],
  claimsByType: []
})

const userInitials = computed(() => {
  const name = authStore.user?.name || ''
  return name.split(' ').map(n => n[0]).join('').toUpperCase() || 'U'
})

const formatAmount = (amount: number) => {
  return new Intl.NumberFormat('en-US', {
    minimumFractionDigits: 2,
    maximumFractionDigits: 2
  }).format(amount)
}

const formatDate = (date: string | undefined) => {
  if (!date) return ''
  return new Date(date).toLocaleDateString('en-US', {
    year: 'numeric',
    month: 'short',
    day: 'numeric'
  })
}

const formatStatus = (status: string) => {
  return status.split('-').map(word =>
    word.charAt(0).toUpperCase() + word.slice(1)
  ).join(' ')
}

const getStatusIcon = (status: string) => {
  const icons: Record<string, string> = {
    'draft': 'pi pi-pencil',
    'submitted': 'pi pi-send',
    'approved': 'pi pi-check-circle',
    'rejected': 'pi pi-times-circle',
    'payment-in-progress': 'pi pi-spin pi-spinner',
    'paid': 'pi pi-dollar'
  }
  return icons[status] || 'pi pi-circle'
}

const loadDashboard = async () => {
  loading.value = true
  try {
    const response = await dashboardApi.getStats()
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
.dashboard-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  gap: 1.5rem;
  margin-bottom: 2rem;
}

.stat-content {
  display: flex;
  align-items: center;
  gap: 1rem;
}

.stat-icon-container {
  width: 60px;
  height: 60px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.stat-icon {
  font-size: 1.5rem;
}

.stat-details {
  flex: 1;
}

.stat-label {
  color: var(--surface-600);
  font-size: 0.875rem;
  margin: 0;
}

.stat-value {
  font-size: 1.75rem;
  font-weight: 700;
  color: var(--surface-900);
  margin: 0.25rem 0 0;
}

.dashboard-content {
  display: grid;
  grid-template-columns: 1fr 400px;
  gap: 1.5rem;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.card-title {
  font-size: 1.25rem;
  font-weight: 600;
  margin: 0;
}

.claim-link {
  color: var(--primary-600);
  text-decoration: none;
  font-weight: 500;
}

.claim-link:hover {
  text-decoration: underline;
}

.empty-state {
  text-align: center;
  padding: 3rem;
}

.empty-state p {
  margin-bottom: 1.5rem;
}

.empty-icon {
  font-size: 3rem;
  color: var(--surface-400);
  margin-bottom: 1rem;
}

.quick-actions-grid {
  display: grid;
  gap: 1rem;
}

.p-button-lg {
  width: 100%;
  padding: 0.75rem 1.25rem;
}

.guidelines-list {
  margin: 0;
  padding-left: 1.5rem;
}

.guidelines-list li {
  margin-bottom: 0.5rem;
  color: var(--surface-700);
}

.profile-content {
  text-align: center;
}

.profile-avatar {
  margin-bottom: 1.5rem;
}

.profile-info {
  text-align: left;
}

.profile-info p {
  margin-bottom: 0.75rem;
}

@media (max-width: 1200px) {
  .dashboard-content {
    grid-template-columns: 1fr;
  }
}

@media (max-width: 768px) {
  .dashboard-grid {
    grid-template-columns: 1fr;
  }
}
</style>
