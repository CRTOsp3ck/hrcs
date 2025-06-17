<template>
  <div class="admin-page-container">
    <PageHeader
      :title="pageTitle"
      :breadcrumbs="breadcrumbItems"
    >
      <template #actions>
        <Button 
          icon="pi pi-pencil" 
          label="Edit Claim Type" 
          @click="editClaimType" 
          :loading="loading"
        />
        <Button 
          icon="pi pi-cog" 
          label="Manage Limits" 
          severity="secondary" 
          @click="manageLimits"
          :loading="loading"
        />
      </template>
    </PageHeader>

    <div v-if="loading" class="loading-container">
      <ProgressSpinner />
    </div>

    <div v-else-if="error" class="error-container">
      {{ error }}
    </div>

    <div v-else-if="claimTypeDetails" class="content-grid">
      <!-- Basic Information -->
      <Card class="col-span-full lg:col-span-1">
        <template #title>Basic Information</template>
        <template #content>
          <div class="grid grid-cols-1 md:grid-cols-2" style="gap: var(--space-13);">
            <div class="info-item">
              <label>Name:</label>
              <p>{{ claimTypeDetails.claim_type.name }}</p>
            </div>
            <div class="info-item">
              <label>Description:</label>
              <p>{{ claimTypeDetails.claim_type.description || 'No description' }}</p>
            </div>
            <div class="info-item">
              <label>Limit Amount:</label>
              <p>${{ claimTypeDetails.claim_type.limit_amount?.toFixed(2) || '0.00' }}</p>
            </div>
            <div class="info-item">
              <label>Limit Timespan:</label>
              <Tag :value="claimTypeDetails.claim_type.limit_timespan" severity="info" style="margin-top: var(--space-3);" />
            </div>
            <div class="info-item">
              <label>Created:</label>
              <p>{{ formatDate(claimTypeDetails.claim_type.created_at) }}</p>
            </div>
          </div>
        </template>
      </Card>

      <!-- Statistics -->
      <Card class="col-span-full lg:col-span-1">
        <template #title>Statistics</template>
        <template #content>
          <div class="stats-grid">
            <div class="metric-card">
              <div class="metric-content">
                <div class="metric-icon" style="background: rgba(59, 130, 246, 0.1);">
                  <i class="pi pi-file" style="color: #3b82f6;"></i>
                </div>
                <div class="metric-details">
                  <p class="metric-label">Total Claims</p>
                  <p class="metric-value">{{ claimTypeDetails.stats.total_claims }}</p>
                </div>
              </div>
            </div>
            <div class="metric-card">
              <div class="metric-content">
                <div class="metric-icon" style="background: rgba(34, 197, 94, 0.1);">
                  <i class="pi pi-check" style="color: #22c55e;"></i>
                </div>
                <div class="metric-details">
                  <p class="metric-label">Approved Claims</p>
                  <p class="metric-value">{{ claimTypeDetails.stats.approved_claims }}</p>
                </div>
              </div>
            </div>
            <div class="metric-card">
              <div class="metric-content">
                <div class="metric-icon" style="background: rgba(168, 85, 247, 0.1);">
                  <i class="pi pi-dollar" style="color: #a855f7;"></i>
                </div>
                <div class="metric-details">
                  <p class="metric-label">Total Amount</p>
                  <p class="metric-value">${{ claimTypeDetails.stats.total_amount.toFixed(2) }}</p>
                </div>
              </div>
            </div>
            <div class="metric-card">
              <div class="metric-content">
                <div class="metric-icon" style="background: rgba(249, 115, 22, 0.1);">
                  <i class="pi pi-chart-bar" style="color: #f97316;"></i>
                </div>
                <div class="metric-details">
                  <p class="metric-label">Average Amount</p>
                  <p class="metric-value">${{ claimTypeDetails.stats.average_amount.toFixed(2) }}</p>
                </div>
              </div>
            </div>
          </div>
        </template>
      </Card>

      <!-- User Group Permissions -->
      <Card class="col-span-full">
        <template #title>User Group Permissions</template>
        <template #content>
          <div class="section-actions">
            <Button 
              icon="pi pi-plus" 
              label="Configure Permissions" 
              @click="configurePermissions"
              size="small"
            />
          </div>
          <DataTable 
            v-if="claimTypeDetails.group_permissions.length > 0"
            :value="claimTypeDetails.group_permissions" 
            :loading="loading"
            responsiveLayout="scroll"
            class="p-datatable-sm"
          >
            <Column field="user_group.name" header="User Group" />
            <Column field="is_allowed" header="Allowed">
              <template #body="slotProps">
                <Tag 
                  :value="slotProps.data.is_allowed ? 'Yes' : 'No'" 
                  :severity="slotProps.data.is_allowed ? 'success' : 'danger'" 
                />
              </template>
            </Column>
            <Column field="custom_limit_amount" header="Custom Limit">
              <template #body="slotProps">
                {{ slotProps.data.custom_limit_amount ? 
                   '$' + slotProps.data.custom_limit_amount.toFixed(2) : 
                   'Default' }}
              </template>
            </Column>
            <Column header="Actions">
              <template #body="slotProps">
                <Button 
                  icon="pi pi-eye" 
                  size="small" 
                  text 
                  @click="viewGroupDetails(slotProps.data.user_group.id)"
                  v-tooltip.top="'View Group Details'"
                />
                <Button 
                  icon="pi pi-pencil" 
                  size="small" 
                  text 
                  @click="editGroupPermission(slotProps.data)"
                  v-tooltip.top="'Edit Permission'"
                />
              </template>
            </Column>
          </DataTable>
          <div v-else class="no-data-message">
            No group permissions configured
          </div>
        </template>
      </Card>

      <!-- Recent Claims -->
      <Card class="col-span-full">
        <template #title>Recent Claims</template>
        <template #content>
          <DataTable 
            v-if="claimTypeDetails.recent_claims.length > 0"
            :value="claimTypeDetails.recent_claims" 
            :loading="loading"
            responsiveLayout="scroll"
            class="p-datatable-sm"
          >
            <Column field="title" header="Title" />
            <Column header="User">
              <template #body="slotProps">
                {{ slotProps.data.user?.first_name }} {{ slotProps.data.user?.last_name }}
              </template>
            </Column>
            <Column field="amount" header="Amount">
              <template #body="slotProps">
                ${{ slotProps.data.amount.toFixed(2) }}
              </template>
            </Column>
            <Column field="status" header="Status">
              <template #body="slotProps">
                <Tag 
                  :value="slotProps.data.status" 
                  :severity="getStatusSeverity(slotProps.data.status)" 
                />
              </template>
            </Column>
            <Column field="created_at" header="Created">
              <template #body="slotProps">
                {{ formatDate(slotProps.data.created_at) }}
              </template>
            </Column>
            <Column header="Actions">
              <template #body="slotProps">
                <Button 
                  icon="pi pi-eye" 
                  size="small" 
                  text 
                  @click="viewClaim(slotProps.data.id)"
                  v-tooltip.top="'View Claim'"
                />
              </template>
            </Column>
          </DataTable>
          <div v-else class="no-data-message">
            No claims found
          </div>
        </template>
      </Card>
    </div>

    <!-- Edit Claim Type Dialog -->
    <Dialog v-model:visible="showEditDialog" header="Edit Claim Type" modal class="w-full max-w-2xl">
      <div class="grid grid-cols-1" style="gap: var(--space-13); padding: var(--space-13) 0;">
        <div>
          <label class="block" style="font-size: var(--text-sm); font-weight: var(--font-medium); color: var(--surface-700); margin-bottom: var(--space-5);">Name</label>
          <InputText v-model="editForm.name" class="w-full" />
        </div>
        <div>
          <label class="block" style="font-size: var(--text-sm); font-weight: var(--font-medium); color: var(--surface-700); margin-bottom: var(--space-5);">Description</label>
          <Textarea v-model="editForm.description" class="w-full" rows="3" />
        </div>
      </div>
      <template #footer>
        <Button label="Cancel" severity="secondary" @click="showEditDialog = false" />
        <Button label="Save" @click="saveClaimType" :loading="saving" />
      </template>
    </Dialog>

    <!-- Manage Limits Dialog -->
    <Dialog v-model:visible="showLimitsDialog" header="Manage Claim Limits" modal class="w-full max-w-2xl">
      <div class="grid grid-cols-1 md:grid-cols-2" style="gap: var(--space-13); padding: var(--space-13) 0;">
        <div>
          <label class="block" style="font-size: var(--text-sm); font-weight: var(--font-medium); color: var(--surface-700); margin-bottom: var(--space-5);">Limit Amount ($)</label>
          <InputNumber 
            v-model="limitsForm.limit_amount" 
            mode="currency" 
            currency="USD" 
            locale="en-US"
            class="w-full" 
          />
        </div>
        <div>
          <label class="block" style="font-size: var(--text-sm); font-weight: var(--font-medium); color: var(--surface-700); margin-bottom: var(--space-5);">Limit Timespan</label>
          <Dropdown 
            v-model="limitsForm.limit_timespan" 
            :options="timespanOptions" 
            optionLabel="label" 
            optionValue="value"
            class="w-full" 
          />
        </div>
      </div>
      <template #footer>
        <Button label="Cancel" severity="secondary" @click="showLimitsDialog = false" />
        <Button label="Save Limits" @click="saveLimits" :loading="saving" />
      </template>
    </Dialog>

    <!-- Configure Permissions Dialog -->
    <Dialog v-model:visible="showPermissionsDialog" header="Configure Group Permissions" modal class="w-full max-w-4xl">
      <div class="py-4">
        <p style="color: var(--surface-600); margin-bottom: var(--space-13);">Configure which user groups can access this claim type:</p>
        <!-- TODO: Implement permissions configuration interface -->
        <div class="no-data-message">
          Permissions configuration interface will be implemented here
        </div>
      </div>
      <template #footer>
        <Button label="Close" @click="showPermissionsDialog = false" />
      </template>
    </Dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useToast } from 'primevue/usetoast'
import { adminApi } from '@/api'
import type { ClaimTypeDetails } from '@/types'
import Button from 'primevue/button'
import Card from 'primevue/card'
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import Tag from 'primevue/tag'
import ProgressSpinner from 'primevue/progressspinner'
import Dialog from 'primevue/dialog'
import InputText from 'primevue/inputtext'
import Textarea from 'primevue/textarea'
import InputNumber from 'primevue/inputnumber'
import Dropdown from 'primevue/dropdown'
import PageHeader from '@/components/base/PageHeader.vue'

const route = useRoute()
const router = useRouter()
const toast = useToast()

const claimTypeDetails = ref<ClaimTypeDetails | null>(null)
const loading = ref(false)
const error = ref('')
const showEditDialog = ref(false)
const showLimitsDialog = ref(false)
const showPermissionsDialog = ref(false)
const saving = ref(false)

const editForm = ref({
  name: '',
  description: ''
})

const limitsForm = ref({
  limit_amount: 0,
  limit_timespan: 'annual' as 'annual' | 'monthly' | 'weekly' | 'daily'
})

const timespanOptions = [
  { label: 'Annual', value: 'annual' },
  { label: 'Monthly', value: 'monthly' },
  { label: 'Weekly', value: 'weekly' },
  { label: 'Daily', value: 'daily' }
]

const pageTitle = computed(() => {
  if (!claimTypeDetails.value?.claim_type) return 'Claim Type Details'
  return `Claim Type Details: ${claimTypeDetails.value.claim_type.name}`
})

const breadcrumbItems = computed(() => {
  const items: Array<{ label: string; to?: string }> = [{ label: 'Claim Types', to: '/admin/claim-types' }]
  if (claimTypeDetails.value?.claim_type) {
    items.push({ label: claimTypeDetails.value.claim_type.name })
  }
  return items
})

const fetchClaimTypeDetails = async () => {
  try {
    loading.value = true
    error.value = ''
    const claimTypeId = parseInt(route.params.id as string)
    const response = await adminApi.getClaimTypeDetails(claimTypeId)
    if (response.data.success) {
      claimTypeDetails.value = response.data.data as ClaimTypeDetails
    } else {
      error.value = response.data.message || 'Failed to load claim type details'
    }
  } catch (err: any) {
    error.value = err.message || 'Failed to load claim type details'
  } finally {
    loading.value = false
  }
}

const editClaimType = () => {
  if (claimTypeDetails.value) {
    editForm.value = {
      name: claimTypeDetails.value.claim_type.name,
      description: claimTypeDetails.value.claim_type.description || ''
    }
    showEditDialog.value = true
  }
}

const saveClaimType = async () => {
  try {
    saving.value = true
    const claimTypeId = parseInt(route.params.id as string)
    await adminApi.updateClaimType(claimTypeId, editForm.value)
    showEditDialog.value = false
    toast.add({
      severity: 'success',
      summary: 'Success',
      detail: 'Claim type updated successfully'
    })
    await fetchClaimTypeDetails() // Refresh data
  } catch (err: any) {
    toast.add({
      severity: 'error',
      summary: 'Error',
      detail: err.message || 'Failed to update claim type'
    })
  } finally {
    saving.value = false
  }
}

const manageLimits = () => {
  if (claimTypeDetails.value) {
    limitsForm.value = {
      limit_amount: claimTypeDetails.value.claim_type.limit_amount,
      limit_timespan: claimTypeDetails.value.claim_type.limit_timespan
    }
    showLimitsDialog.value = true
  }
}

const saveLimits = async () => {
  try {
    saving.value = true
    const claimTypeId = parseInt(route.params.id as string)
    await adminApi.updateClaimTypeLimits(claimTypeId, limitsForm.value)
    showLimitsDialog.value = false
    toast.add({
      severity: 'success',
      summary: 'Success',
      detail: 'Claim limits updated successfully'
    })
    await fetchClaimTypeDetails() // Refresh data
  } catch (err: any) {
    toast.add({
      severity: 'error',
      summary: 'Error',
      detail: err.message || 'Failed to update claim limits'
    })
  } finally {
    saving.value = false
  }
}

const configurePermissions = () => {
  showPermissionsDialog.value = true
}

const viewGroupDetails = (groupId: number) => {
  router.push(`/admin/groups/${groupId}`)
}

const editGroupPermission = (permission: any) => {
  // TODO: Implement group permission editing
  toast.add({
    severity: 'info',
    summary: 'Info',
    detail: 'Group permission editing not yet implemented'
  })
}

const viewClaim = (claimId: number) => {
  router.push(`/claims/${claimId}`)
}

const getStatusSeverity = (status: string) => {
  switch (status) {
    case 'approved':
    case 'paid':
      return 'success'
    case 'rejected':
      return 'danger'
    case 'submitted':
    case 'payment-in-progress':
      return 'warning'
    default:
      return 'info'
  }
}

const formatDate = (dateString: string) => {
  return new Date(dateString).toLocaleDateString()
}

onMounted(() => {
  fetchClaimTypeDetails()
})
</script>

<style scoped>
.admin-page-container {
  max-width: var(--container-lg);
  margin: 0 auto;
  padding: var(--space-21);
}

.breadcrumb {
  font-size: var(--text-sm);
}

.info-item {
  background-color: var(--surface-50);
  padding: var(--space-8);
  border-radius: var(--space-5);
  margin-bottom: var(--space-8);
  transition: all var(--transition-fast) ease;
}

.info-item:hover {
  background-color: var(--surface-100);
}

.info-item label {
  color: var(--surface-700);
  font-weight: var(--font-semibold);
  font-size: var(--text-sm);
}

.info-item p {
  color: var(--surface-900);
  margin-top: var(--space-3);
}

.metric-card {
  background: var(--surface-0);
  border: 1px solid var(--surface-200);
  border-radius: var(--space-8);
  padding: var(--space-13);
  transition: all var(--transition-fast) ease;
}

.metric-card:hover {
  border-color: var(--surface-300);
  box-shadow: 0 var(--space-2) var(--space-8) rgba(0, 0, 0, 0.1);
}

.metric-content {
  display: flex;
  align-items: center;
  gap: var(--space-8);
}

.metric-icon {
  width: var(--space-34);
  height: var(--space-34);
  border-radius: var(--space-8);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: var(--text-lg);
}

.metric-details {
  flex: 1;
}

.metric-label {
  font-size: var(--text-sm);
  color: var(--surface-600);
  margin: 0 0 var(--space-3) 0;
}

.metric-value {
  font-size: var(--text-xl);
  font-weight: var(--font-bold);
  color: var(--surface-900);
  margin: 0;
  line-height: 1.2;
}

.error-container {
  background: var(--red-50);
  border: 1px solid var(--red-200);
  color: var(--red-700);
  padding: var(--space-13);
  border-radius: var(--space-8);
  margin-bottom: var(--space-13);
}

.loading-container {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 16rem;
}

.content-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: var(--space-21);
}

.content-grid .col-span-full {
  grid-column: 1 / -1;
}

.no-data-message {
  text-align: center;
  padding: var(--space-34) 0;
  color: var(--surface-500);
}

.section-actions {
  margin-bottom: var(--space-13);
}

.stats-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: var(--space-13);
}

@media (max-width: var(--breakpoint-lg)) {
  .content-grid {
    grid-template-columns: 1fr;
    gap: var(--space-13);
  }
}

@media (max-width: var(--breakpoint-md)) {
  .admin-page-container {
    padding: var(--space-13);
  }
  
  .content-grid {
    gap: var(--space-8);
  }
  
  .stats-grid {
    grid-template-columns: 1fr;
    gap: var(--space-8);
  }
}
</style>