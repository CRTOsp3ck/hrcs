<template>
  <div class="claim-type-details-view">
    <div class="flex justify-between items-center mb-6">
      <div>
        <h1 class="text-3xl font-bold text-gray-900">
          Claim Type Details: {{ claimTypeDetails?.claim_type?.name }}
        </h1>
        <nav class="breadcrumb mt-2">
          <router-link to="/admin/claim-types" class="text-blue-600 hover:text-blue-800">Claim Types</router-link>
          <span class="mx-2">/</span>
          <span class="text-gray-500">{{ claimTypeDetails?.claim_type?.name }}</span>
        </nav>
      </div>
      <div class="flex gap-3">
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
      </div>
    </div>

    <div v-if="loading" class="flex justify-center items-center h-64">
      <ProgressSpinner />
    </div>

    <div v-else-if="error" class="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded mb-4">
      {{ error }}
    </div>

    <div v-else-if="claimTypeDetails" class="grid grid-cols-1 lg:grid-cols-2 gap-6">
      <!-- Basic Information -->
      <Card class="col-span-full lg:col-span-1">
        <template #title>Basic Information</template>
        <template #content>
          <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div class="info-item">
              <label class="font-semibold text-gray-700">Name:</label>
              <p class="mt-1">{{ claimTypeDetails.claim_type.name }}</p>
            </div>
            <div class="info-item">
              <label class="font-semibold text-gray-700">Description:</label>
              <p class="mt-1">{{ claimTypeDetails.claim_type.description || 'No description' }}</p>
            </div>
            <div class="info-item">
              <label class="font-semibold text-gray-700">Limit Amount:</label>
              <p class="mt-1">${{ claimTypeDetails.claim_type.limit_amount?.toFixed(2) || '0.00' }}</p>
            </div>
            <div class="info-item">
              <label class="font-semibold text-gray-700">Limit Timespan:</label>
              <Tag :value="claimTypeDetails.claim_type.limit_timespan" severity="info" />
            </div>
            <div class="info-item">
              <label class="font-semibold text-gray-700">Created:</label>
              <p class="mt-1">{{ formatDate(claimTypeDetails.claim_type.created_at) }}</p>
            </div>
          </div>
        </template>
      </Card>

      <!-- Statistics -->
      <Card class="col-span-full lg:col-span-1">
        <template #title>Statistics</template>
        <template #content>
          <div class="grid grid-cols-2 gap-4">
            <div class="stat-card bg-blue-50 border border-blue-200 p-4 rounded-lg">
              <div class="flex items-center">
                <i class="pi pi-file text-blue-600 text-xl mr-3"></i>
                <div>
                  <p class="text-sm text-gray-600">Total Claims</p>
                  <p class="text-2xl font-bold text-blue-900">{{ claimTypeDetails.stats.total_claims }}</p>
                </div>
              </div>
            </div>
            <div class="stat-card bg-green-50 border border-green-200 p-4 rounded-lg">
              <div class="flex items-center">
                <i class="pi pi-check text-green-600 text-xl mr-3"></i>
                <div>
                  <p class="text-sm text-gray-600">Approved Claims</p>
                  <p class="text-2xl font-bold text-green-900">{{ claimTypeDetails.stats.approved_claims }}</p>
                </div>
              </div>
            </div>
            <div class="stat-card bg-purple-50 border border-purple-200 p-4 rounded-lg">
              <div class="flex items-center">
                <i class="pi pi-dollar text-purple-600 text-xl mr-3"></i>
                <div>
                  <p class="text-sm text-gray-600">Total Amount</p>
                  <p class="text-2xl font-bold text-purple-900">${{ claimTypeDetails.stats.total_amount.toFixed(2) }}</p>
                </div>
              </div>
            </div>
            <div class="stat-card bg-orange-50 border border-orange-200 p-4 rounded-lg">
              <div class="flex items-center">
                <i class="pi pi-chart-bar text-orange-600 text-xl mr-3"></i>
                <div>
                  <p class="text-sm text-gray-600">Average Amount</p>
                  <p class="text-2xl font-bold text-orange-900">${{ claimTypeDetails.stats.average_amount.toFixed(2) }}</p>
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
          <div class="mb-4">
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
          <div v-else class="text-center py-8 text-gray-500">
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
          <div v-else class="text-center py-8 text-gray-500">
            No claims found
          </div>
        </template>
      </Card>
    </div>

    <!-- Edit Claim Type Dialog -->
    <Dialog v-model:visible="showEditDialog" header="Edit Claim Type" modal class="w-full max-w-2xl">
      <div class="grid grid-cols-1 gap-4 py-4">
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-2">Name</label>
          <InputText v-model="editForm.name" class="w-full" />
        </div>
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-2">Description</label>
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
      <div class="grid grid-cols-1 md:grid-cols-2 gap-4 py-4">
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-2">Limit Amount ($)</label>
          <InputNumber 
            v-model="limitsForm.limit_amount" 
            mode="currency" 
            currency="USD" 
            locale="en-US"
            class="w-full" 
          />
        </div>
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-2">Limit Timespan</label>
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
        <p class="text-gray-600 mb-4">Configure which user groups can access this claim type:</p>
        <!-- TODO: Implement permissions configuration interface -->
        <div class="text-center py-8 text-gray-500">
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
import { ref, onMounted } from 'vue'
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
.claim-type-details-view {
  max-width: 80rem;
  margin: 0 auto;
  padding: 1.5rem;
}

.breadcrumb {
  font-size: 0.875rem;
}

.info-item {
  background-color: #f9fafb;
  padding: 0.75rem;
  border-radius: 0.25rem;
}

.stat-card {
  transition: all 0.2s;
}

.stat-card:hover {
  box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1);
}
</style>