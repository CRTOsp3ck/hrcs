<template>
  <div class="admin-page-container">
    <PageHeader
      :title="pageTitle"
      :breadcrumbs="breadcrumbItems"
    >
      <template #actions>
        <Button 
          icon="pi pi-pencil" 
          label="Edit User" 
          @click="editUser" 
          :loading="loading"
        />
        <Button 
          icon="pi pi-key" 
          label="Reset Password" 
          severity="secondary" 
          @click="resetPassword"
          :loading="loading"
        />
        <Button 
          icon="pi pi-cog" 
          label="Manage Permissions" 
          severity="secondary" 
          @click="managePermissions"
          :loading="loading"
        />
        <Button 
          icon="pi pi-wallet" 
          label="Adjust Balances" 
          severity="secondary" 
          @click="adjustBalances"
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

    <div v-else-if="userDetails" class="content-grid">
      <!-- Basic Information -->
      <Card class="col-span-full lg:col-span-1">
        <template #title>Basic Information</template>
        <template #content>
          <div class="grid grid-cols-1 md:grid-cols-2" style="gap: var(--space-13);">
            <div class="info-item">
              <label>Full Name:</label>
              <p>{{ userDetails.user.first_name }} {{ userDetails.user.last_name }}</p>
            </div>
            <div class="info-item">
              <label>Email:</label>
              <p>{{ userDetails.user.email }}</p>
            </div>
            <div class="info-item">
              <label>Role:</label>
              <Tag 
                :value="userDetails.user.role" 
                :severity="userDetails.user.role === 'admin' ? 'danger' : 'info'"
                style="margin-top: var(--space-3);"
              />
            </div>
            <div class="info-item">
              <label>User Group:</label>
              <p>{{ userDetails.user.user_group?.name || 'No group assigned' }}</p>
            </div>
            <div class="info-item">
              <label>Member Since:</label>
              <p>{{ formatDate(userDetails.user.created_at) }}</p>
            </div>
          </div>
        </template>
      </Card>

      <!-- Claim Balances -->
      <Card class="col-span-full lg:col-span-1">
        <template #title>Claim Balances</template>
        <template #content>
          <DataTable 
            v-if="userDetails.balances.length > 0"
            :value="userDetails.balances" 
            :loading="loading"
            responsiveLayout="scroll"
            class="p-datatable-sm"
          >
            <Column field="claim_type.name" header="Claim Type" />
            <Column field="total_limit" header="Total Limit">
              <template #body="slotProps">
                ${{ slotProps.data.total_limit.toFixed(2) }}
              </template>
            </Column>
            <Column field="current_spent" header="Current Spent">
              <template #body="slotProps">
                ${{ slotProps.data.current_spent.toFixed(2) }}
              </template>
            </Column>
            <Column field="remaining_balance" header="Remaining">
              <template #body="slotProps">
                <span :class="getBalanceClass(slotProps.data.remaining_balance)">
                  ${{ slotProps.data.remaining_balance.toFixed(2) }}
                </span>
              </template>
            </Column>
            <Column field="reset_period" header="Reset Period">
              <template #body="slotProps">
                <Tag :value="slotProps.data.reset_period" severity="info" />
              </template>
            </Column>
          </DataTable>
          <div v-else class="no-data-message">
            No balance records found
          </div>
        </template>
      </Card>

      <!-- Claims History -->
      <Card class="col-span-full">
        <template #title>Recent Claims</template>
        <template #content>
          <DataTable 
            v-if="userDetails.claims.length > 0"
            :value="userDetails.claims" 
            :loading="loading"
            responsiveLayout="scroll"
            class="p-datatable-sm"
          >
            <Column field="title" header="Title" />
            <Column field="claim_type.name" header="Type" />
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
                />
              </template>
            </Column>
          </DataTable>
          <div v-else class="no-data-message">
            No claims found
          </div>
        </template>
      </Card>

      <!-- Permissions Override -->
      <Card class="col-span-full" v-if="userDetails.permissions.length > 0">
        <template #title>Claim Type Permissions</template>
        <template #content>
          <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3" style="gap: var(--space-13);">
            <div 
              v-for="permission in userDetails.permissions" 
              :key="permission.claim_type_id" 
              class="permission-card"
            >
              <div class="flex justify-between items-center" style="margin-bottom: var(--space-5);">
                <span style="font-weight: var(--font-semibold);">{{ permission.claim_type.name }}</span>
                <Tag 
                  :value="permission.is_allowed ? 'Allowed' : 'Denied'" 
                  :severity="permission.is_allowed ? 'success' : 'danger'" 
                />
              </div>
              <div v-if="permission.custom_limit_amount" style="font-size: var(--text-sm); color: var(--surface-600);">
                Custom Limit: ${{ permission.custom_limit_amount.toFixed(2) }}
              </div>
            </div>
          </div>
        </template>
      </Card>
    </div>

    <!-- Edit User Dialog -->
    <Dialog v-model:visible="showEditDialog" header="Edit User" modal class="w-full max-w-2xl">
      <div class="grid grid-cols-1 md:grid-cols-2" style="gap: var(--space-13); padding: var(--space-13) 0;">
        <div>
          <label class="block" style="font-size: var(--text-sm); font-weight: var(--font-medium); color: var(--surface-700); margin-bottom: var(--space-5);">First Name</label>
          <InputText v-model="editForm.first_name" class="w-full" />
        </div>
        <div>
          <label class="block" style="font-size: var(--text-sm); font-weight: var(--font-medium); color: var(--surface-700); margin-bottom: var(--space-5);">Last Name</label>
          <InputText v-model="editForm.last_name" class="w-full" />
        </div>
        <div>
          <label class="block" style="font-size: var(--text-sm); font-weight: var(--font-medium); color: var(--surface-700); margin-bottom: var(--space-5);">Email</label>
          <InputText v-model="editForm.email" class="w-full" />
        </div>
        <div>
          <label class="block" style="font-size: var(--text-sm); font-weight: var(--font-medium); color: var(--surface-700); margin-bottom: var(--space-5);">Role</label>
          <Dropdown 
            v-model="editForm.role" 
            :options="roleOptions" 
            optionLabel="label" 
            optionValue="value"
            class="w-full" 
          />
        </div>
      </div>
      <template #footer>
        <Button label="Cancel" severity="secondary" @click="showEditDialog = false" />
        <Button label="Save" @click="saveUser" :loading="saving" />
      </template>
    </Dialog>

    <!-- Manage Permissions Dialog -->
    <Dialog v-model:visible="showPermissionsDialog" header="Manage User Permission Overrides" modal class="w-full max-w-4xl">
      <div class="py-4">
        <p style="color: var(--surface-600); margin-bottom: var(--space-13);">Configure individual permission overrides for this user. These will override group permissions:</p>
        <div v-if="availableClaimTypes.length === 0" class="no-data-message">
          Loading claim types...
        </div>
        <div v-else class="permissions-grid">
          <div v-for="claimType in availableClaimTypes" :key="claimType.id" class="permission-item">
            <div class="permission-header">
              <h4>{{ claimType.name }}</h4>
              <p class="permission-description">{{ claimType.description || 'No description available' }}</p>
            </div>
            <div class="permission-controls">
              <div class="permission-toggle">
                <label class="toggle-label">Override Access:</label>
                <SelectButton 
                  v-model="ensurePermissionSetting(claimType.id).override" 
                  :options="overrideOptions" 
                  optionLabel="label" 
                  optionValue="value"
                  @change="updatePermissionSetting(claimType.id)" 
                />
                <small class="override-help">
                  {{ getOverrideText(ensurePermissionSetting(claimType.id).override) }}
                </small>
              </div>
              <div v-if="ensurePermissionSetting(claimType.id).override === 'allow'" class="custom-limit">
                <label class="limit-label">Custom Limit (optional):</label>
                <div class="limit-input-group">
                  <InputNumber 
                    v-model="ensurePermissionSetting(claimType.id).custom_limit_amount"
                    mode="currency" 
                    currency="USD" 
                    locale="en-US"
                    :min="0"
                    placeholder="Use default limit"
                    class="w-full"
                  />
                  <small class="default-limit-text">
                    Default: ${{ claimType.limit_amount.toFixed(2) }} ({{ claimType.limit_timespan }})
                  </small>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
      <template #footer>
        <Button label="Cancel" severity="secondary" @click="cancelPermissionChanges" />
        <Button label="Save Changes" @click="savePermissions" :loading="savingPermissions" />
      </template>
    </Dialog>

    <!-- Adjust Balances Dialog -->
    <Dialog v-model:visible="showBalancesDialog" header="Adjust User Balances" modal class="w-full max-w-4xl">
      <div class="py-4">
        <p style="color: var(--surface-600); margin-bottom: var(--space-13);">Adjust balance limits for this user. Changes will take effect immediately:</p>
        <div v-if="userDetails?.balances.length === 0" class="no-data-message">
          No balances found for this user
        </div>
        <div v-else class="balances-grid">
          <div v-for="balance in userDetails?.balances" :key="balance.claim_type_id" class="balance-item">
            <div class="balance-header">
              <h4>{{ balance.claim_type.name }}</h4>
              <p class="balance-description">{{ balance.claim_type.description || 'No description available' }}</p>
            </div>
            <div class="balance-info">
              <div class="info-row">
                <span class="info-label">Current Limit:</span>
                <span class="info-value">${{ balance.total_limit.toFixed(2) }}</span>
              </div>
              <div class="info-row">
                <span class="info-label">Current Spent:</span>
                <span class="info-value">${{ balance.current_spent.toFixed(2) }}</span>
              </div>
              <div class="info-row">
                <span class="info-label">Remaining:</span>
                <span class="info-value remaining">${{ balance.remaining_balance.toFixed(2) }}</span>
              </div>
              <div class="info-row">
                <span class="info-label">Reset Period:</span>
                <span class="info-value">{{ balance.reset_period }}</span>
              </div>
            </div>
            <div class="balance-adjustment">
              <label class="adjustment-label">New Limit:</label>
              <InputNumber 
                :model-value="ensureBalanceAdjustment(balance.claim_type_id)"
                @update:model-value="setBalanceAdjustment(balance.claim_type_id, $event)"
                mode="currency" 
                currency="USD" 
                locale="en-US"
                :min="0"
                class="w-full"
                :placeholder="balance.total_limit.toFixed(2)"
              />
              <small class="adjustment-help">
                Leave empty to keep current limit of ${{ balance.total_limit.toFixed(2) }}
              </small>
            </div>
          </div>
        </div>
      </div>
      <template #footer>
        <Button label="Cancel" severity="secondary" @click="cancelBalanceChanges" />
        <Button label="Apply Changes" @click="saveBalanceAdjustments" :loading="savingBalances" />
      </template>
    </Dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useToast } from 'primevue/usetoast'
import { adminApi, claimTypesApi } from '@/api'
import type { UserDetails, ClaimType } from '@/types'
import Button from 'primevue/button'
import Card from 'primevue/card'
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import Tag from 'primevue/tag'
import ProgressSpinner from 'primevue/progressspinner'
import Dialog from 'primevue/dialog'
import InputText from 'primevue/inputtext'
import Dropdown from 'primevue/dropdown'
import SelectButton from 'primevue/selectbutton'
import InputNumber from 'primevue/inputnumber'
import PageHeader from '@/components/base/PageHeader.vue'

const route = useRoute()
const router = useRouter()
const toast = useToast()

const userDetails = ref<UserDetails | null>(null)
const loading = ref(false)
const error = ref('')
const showEditDialog = ref(false)
const saving = ref(false)
const showPermissionsDialog = ref(false)
const savingPermissions = ref(false)
const availableClaimTypes = ref<ClaimType[]>([])
const permissionSettings = ref<Record<number, { override: string; custom_limit_amount?: number }>>({})
const originalPermissionSettings = ref<Record<number, { override: string; custom_limit_amount?: number }>>({})

const overrideOptions = [
  { label: 'Group', value: 'group' },
  { label: 'Allow', value: 'allow' }, 
  { label: 'Deny', value: 'deny' }
]

const ensurePermissionSetting = (claimTypeId: number) => {
  if (!permissionSettings.value[claimTypeId]) {
    permissionSettings.value[claimTypeId] = {
      override: 'group',
      custom_limit_amount: undefined
    }
  }
  return permissionSettings.value[claimTypeId]
}
const showBalancesDialog = ref(false)
const savingBalances = ref(false)
const balanceAdjustments = ref<Record<number, number | null>>({})

const ensureBalanceAdjustment = (claimTypeId: number) => {
  if (!(claimTypeId in balanceAdjustments.value)) {
    balanceAdjustments.value[claimTypeId] = null
  }
  return balanceAdjustments.value[claimTypeId]
}

const setBalanceAdjustment = (claimTypeId: number, value: number | null) => {
  balanceAdjustments.value[claimTypeId] = value
}

const editForm = ref({
  first_name: '',
  last_name: '',
  email: '',
  role: 'normal' as 'admin' | 'normal'
})

const roleOptions = [
  { label: 'Normal User', value: 'normal' },
  { label: 'Administrator', value: 'admin' }
]

const pageTitle = computed(() => {
  if (!userDetails.value?.user) return 'User Details'
  return `User Details: ${userDetails.value.user.first_name} ${userDetails.value.user.last_name}`
})

const breadcrumbItems = computed(() => {
  const items: Array<{ label: string; to?: string }> = [{ label: 'Users', to: '/admin/users' }]
  if (userDetails.value?.user) {
    items.push({ 
      label: `${userDetails.value.user.first_name} ${userDetails.value.user.last_name}`
    })
  }
  return items
})

const fetchUserDetails = async () => {
  try {
    loading.value = true
    error.value = ''
    const userId = parseInt(route.params.id as string)
    const response = await adminApi.getUserDetails(userId)
    if (response.data.success) {
      userDetails.value = response.data.data as UserDetails
    } else {
      error.value = response.data.message || 'Failed to load user details'
    }
  } catch (err: any) {
    error.value = err.message || 'Failed to load user details'
  } finally {
    loading.value = false
  }
}

const editUser = () => {
  if (userDetails.value) {
    editForm.value = {
      first_name: userDetails.value.user.first_name,
      last_name: userDetails.value.user.last_name,
      email: userDetails.value.user.email,
      role: userDetails.value.user.role
    }
    showEditDialog.value = true
  }
}

const saveUser = async () => {
  try {
    saving.value = true
    const userId = parseInt(route.params.id as string)
    await adminApi.updateUser(userId, editForm.value)
    showEditDialog.value = false
    toast.add({
      severity: 'success',
      summary: 'Success',
      detail: 'User updated successfully'
    })
    await fetchUserDetails() // Refresh data
  } catch (err: any) {
    toast.add({
      severity: 'error',
      summary: 'Error',
      detail: err.message || 'Failed to update user'
    })
  } finally {
    saving.value = false
  }
}

const resetPassword = () => {
  // TODO: Implement password reset functionality
  toast.add({
    severity: 'info',
    summary: 'Info',
    detail: 'Password reset functionality not yet implemented'
  })
}

const managePermissions = async () => {
  await loadClaimTypes()
  initializePermissionSettings()
  showPermissionsDialog.value = true
}

const loadClaimTypes = async () => {
  try {
    const response = await claimTypesApi.getAll()
    if (response.data.data) {
      availableClaimTypes.value = response.data.data
    }
  } catch (error) {
    toast.add({
      severity: 'error',
      summary: 'Error',
      detail: 'Failed to load claim types'
    })
  }
}

const initializePermissionSettings = () => {
  const settings: Record<number, { override: boolean | null; custom_limit_amount?: number }> = {}
  
  // First, initialize all claim types with default values
  availableClaimTypes.value.forEach(claimType => {
    settings[claimType.id] = {
      override: 'group',
      custom_limit_amount: undefined
    }
  })
  
  // Then override with existing permissions
  if (userDetails.value?.permissions) {
    userDetails.value.permissions.forEach(permission => {
      if (settings[permission.claim_type_id]) {
        settings[permission.claim_type_id] = {
          override: permission.is_allowed ? 'allow' : 'deny',
          custom_limit_amount: permission.custom_limit_amount
        }
      }
    })
  }
  
  permissionSettings.value = settings
  originalPermissionSettings.value = JSON.parse(JSON.stringify(settings))
}

const updatePermissionSetting = (claimTypeId: number) => {
  if (permissionSettings.value[claimTypeId] && permissionSettings.value[claimTypeId].override !== 'allow') {
    // Clear custom limit when not allowing access
    permissionSettings.value[claimTypeId].custom_limit_amount = undefined
  }
}

const getOverrideText = (override: string) => {
  if (override === 'group') return 'Use group permissions'
  if (override === 'allow') return 'Allow access (override)'
  return 'Deny access (override)'
}

const cancelPermissionChanges = () => {
  permissionSettings.value = JSON.parse(JSON.stringify(originalPermissionSettings.value))
  showPermissionsDialog.value = false
}

const savePermissions = async () => {
  if (!userDetails.value?.user.id) return
  
  savingPermissions.value = true
  try {
    const overrides = Object.entries(permissionSettings.value)
      .filter(([, setting]) => setting.override !== 'group')
      .map(([claimTypeId, setting]) => ({
        claim_type_id: parseInt(claimTypeId),
        is_allowed: setting.override === 'allow',
        custom_limit_amount: setting.custom_limit_amount
      }))
    
    await adminApi.setUserClaimOverrides(userDetails.value.user.id, { overrides })
    
    toast.add({
      severity: 'success',
      summary: 'Success',
      detail: 'Permission overrides updated successfully'
    })
    
    showPermissionsDialog.value = false
    // Refresh the user details to show updated permissions
    await fetchUserDetails()
  } catch (error) {
    toast.add({
      severity: 'error',
      summary: 'Error',
      detail: 'Failed to update permission overrides'
    })
  } finally {
    savingPermissions.value = false
  }
}

const adjustBalances = () => {
  initializeBalanceAdjustments()
  showBalancesDialog.value = true
}

const initializeBalanceAdjustments = () => {
  const adjustments: Record<number, number | null> = {}
  
  // Initialize all claim type balances as null (no change)
  if (userDetails.value?.balances) {
    userDetails.value.balances.forEach(balance => {
      adjustments[balance.claim_type_id] = null
    })
  }
  
  balanceAdjustments.value = adjustments
}

const cancelBalanceChanges = () => {
  balanceAdjustments.value = {}
  showBalancesDialog.value = false
}

const saveBalanceAdjustments = async () => {
  if (!userDetails.value?.user.id) return
  
  // Filter out null values (no changes)
  const changes = Object.entries(balanceAdjustments.value)
    .filter(([, newLimit]) => newLimit !== null && newLimit !== undefined)
    .map(([claimTypeId, newLimit]) => ({
      user_id: userDetails.value!.user.id,
      claim_type_id: parseInt(claimTypeId),
      new_limit: newLimit as number
    }))
  
  if (changes.length === 0) {
    toast.add({
      severity: 'info',
      summary: 'No Changes',
      detail: 'No balance adjustments to apply'
    })
    showBalancesDialog.value = false
    return
  }
  
  savingBalances.value = true
  try {
    // Apply each balance adjustment
    for (const change of changes) {
      await adminApi.adjustBalance(change)
    }
    
    toast.add({
      severity: 'success',
      summary: 'Success',
      detail: `Updated ${changes.length} balance limit(s) successfully`
    })
    
    showBalancesDialog.value = false
    // Refresh the user details to show updated balances
    await fetchUserDetails()
  } catch (error) {
    toast.add({
      severity: 'error',
      summary: 'Error',
      detail: 'Failed to update balance limits'
    })
  } finally {
    savingBalances.value = false
  }
}

const getBalanceClass = (remaining: number) => {
  if (remaining <= 0) return 'text-red-600 font-semibold'
  if (remaining < 100) return 'text-orange-600'
  return 'text-green-600'
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

const viewClaim = (claimId: number) => {
  router.push(`/claims/${claimId}`)
}

onMounted(() => {
  fetchUserDetails()
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

.permission-card {
  border: 1px solid var(--surface-200);
  border-radius: var(--space-8);
  padding: var(--space-13);
  transition: all var(--transition-fast) ease;
}

.permission-card:hover {
  border-color: var(--surface-300);
  box-shadow: 0 var(--space-2) var(--space-8) rgba(0, 0, 0, 0.1);
}

.no-data-message {
  text-align: center;
  padding: var(--space-34) 0;
  color: var(--surface-500);
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
}

/* Permissions Management Styles */
.permissions-grid {
  display: grid;
  gap: 1rem;
  max-height: 400px;
  overflow-y: auto;
}

.permission-item {
  border: 1px solid var(--surface-200);
  border-radius: 8px;
  padding: 1rem;
  background: white;
}

.permission-header h4 {
  margin: 0 0 0.5rem 0;
  color: var(--surface-800);
  font-size: 1.1rem;
  font-weight: 600;
}

.permission-description {
  margin: 0 0 1rem 0;
  color: var(--surface-600);
  font-size: 0.875rem;
}

.permission-controls {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.permission-toggle {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  flex-wrap: wrap;
}

.toggle-label {
  font-weight: 500;
  color: var(--surface-700);
  min-width: 120px;
}

.override-help {
  color: var(--surface-500);
  font-size: 0.75rem;
  flex-basis: 100%;
  margin-top: 0.25rem;
}

.custom-limit {
  padding-left: 1rem;
  border-left: 2px solid var(--primary-200);
}

.limit-label {
  display: block;
  font-size: 0.875rem;
  font-weight: 500;
  color: var(--surface-700);
  margin-bottom: 0.5rem;
}

.limit-input-group {
  display: flex;
  flex-direction: column;
  gap: 0.25rem;
}

.default-limit-text {
  color: var(--surface-500);
  font-size: 0.75rem;
}

/* Balance Adjustment Styles */
.balances-grid {
  display: grid;
  gap: 1rem;
  max-height: 400px;
  overflow-y: auto;
}

.balance-item {
  border: 1px solid var(--surface-200);
  border-radius: 8px;
  padding: 1rem;
  background: white;
}

.balance-header h4 {
  margin: 0 0 0.5rem 0;
  color: var(--surface-800);
  font-size: 1.1rem;
  font-weight: 600;
}

.balance-description {
  margin: 0 0 1rem 0;
  color: var(--surface-600);
  font-size: 0.875rem;
}

.balance-info {
  margin-bottom: 1rem;
  padding: 0.75rem;
  background: var(--surface-50);
  border-radius: 6px;
  border: 1px solid var(--surface-100);
}

.info-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 0.5rem;
}

.info-row:last-child {
  margin-bottom: 0;
}

.info-label {
  font-size: 0.875rem;
  color: var(--surface-600);
  font-weight: 500;
}

.info-value {
  font-size: 0.875rem;
  color: var(--surface-800);
  font-weight: 600;
}

.info-value.remaining {
  color: var(--green-600);
}

.balance-adjustment {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.adjustment-label {
  font-size: 0.875rem;
  font-weight: 500;
  color: var(--surface-700);
}

.adjustment-help {
  color: var(--surface-500);
  font-size: 0.75rem;
}
</style>