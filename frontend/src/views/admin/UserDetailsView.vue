<template>
  <div class="user-details-view">
    <div class="flex justify-between items-center mb-6">
      <div>
        <h1 class="text-3xl font-bold text-gray-900">
          User Details: {{ userDetails?.user?.first_name }} {{ userDetails?.user?.last_name }}
        </h1>
        <nav class="breadcrumb mt-2">
          <router-link to="/admin/users" class="text-blue-600 hover:text-blue-800">Users</router-link>
          <span class="mx-2">/</span>
          <span class="text-gray-500">{{ userDetails?.user?.first_name }} {{ userDetails?.user?.last_name }}</span>
        </nav>
      </div>
      <div class="flex gap-3">
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
      </div>
    </div>

    <div v-if="loading" class="flex justify-center items-center h-64">
      <ProgressSpinner />
    </div>

    <div v-else-if="error" class="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded mb-4">
      {{ error }}
    </div>

    <div v-else-if="userDetails" class="grid grid-cols-1 lg:grid-cols-2 gap-6">
      <!-- Basic Information -->
      <Card class="col-span-full lg:col-span-1">
        <template #title>Basic Information</template>
        <template #content>
          <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div class="info-item">
              <label class="font-semibold text-gray-700">Full Name:</label>
              <p class="mt-1">{{ userDetails.user.first_name }} {{ userDetails.user.last_name }}</p>
            </div>
            <div class="info-item">
              <label class="font-semibold text-gray-700">Email:</label>
              <p class="mt-1">{{ userDetails.user.email }}</p>
            </div>
            <div class="info-item">
              <label class="font-semibold text-gray-700">Role:</label>
              <Tag 
                :value="userDetails.user.role" 
                :severity="userDetails.user.role === 'admin' ? 'danger' : 'info'"
                class="mt-1"
              />
            </div>
            <div class="info-item">
              <label class="font-semibold text-gray-700">User Group:</label>
              <p class="mt-1">{{ userDetails.user.user_group?.name || 'No group assigned' }}</p>
            </div>
            <div class="info-item">
              <label class="font-semibold text-gray-700">Member Since:</label>
              <p class="mt-1">{{ formatDate(userDetails.user.created_at) }}</p>
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
          <div v-else class="text-center py-8 text-gray-500">
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
          <div v-else class="text-center py-8 text-gray-500">
            No claims found
          </div>
        </template>
      </Card>

      <!-- Permissions Override -->
      <Card class="col-span-full" v-if="userDetails.permissions.length > 0">
        <template #title>Claim Type Permissions</template>
        <template #content>
          <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
            <div 
              v-for="permission in userDetails.permissions" 
              :key="permission.claim_type_id" 
              class="border rounded-lg p-4"
            >
              <div class="flex justify-between items-center mb-2">
                <span class="font-semibold">{{ permission.claim_type.name }}</span>
                <Tag 
                  :value="permission.is_allowed ? 'Allowed' : 'Denied'" 
                  :severity="permission.is_allowed ? 'success' : 'danger'" 
                />
              </div>
              <div v-if="permission.custom_limit_amount" class="text-sm text-gray-600">
                Custom Limit: ${{ permission.custom_limit_amount.toFixed(2) }}
              </div>
            </div>
          </div>
        </template>
      </Card>
    </div>

    <!-- Edit User Dialog -->
    <Dialog v-model:visible="showEditDialog" header="Edit User" modal class="w-full max-w-2xl">
      <div class="grid grid-cols-1 md:grid-cols-2 gap-4 py-4">
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-2">First Name</label>
          <InputText v-model="editForm.first_name" class="w-full" />
        </div>
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-2">Last Name</label>
          <InputText v-model="editForm.last_name" class="w-full" />
        </div>
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-2">Email</label>
          <InputText v-model="editForm.email" class="w-full" />
        </div>
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-2">Role</label>
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
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useToast } from 'primevue/usetoast'
import { adminApi } from '@/api'
import type { UserDetails } from '@/types'
import Button from 'primevue/button'
import Card from 'primevue/card'
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import Tag from 'primevue/tag'
import ProgressSpinner from 'primevue/progressspinner'
import Dialog from 'primevue/dialog'
import InputText from 'primevue/inputtext'
import Dropdown from 'primevue/dropdown'

const route = useRoute()
const router = useRouter()
const toast = useToast()

const userDetails = ref<UserDetails | null>(null)
const loading = ref(false)
const error = ref('')
const showEditDialog = ref(false)
const saving = ref(false)

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
.user-details-view {
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
</style>