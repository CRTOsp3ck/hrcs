<template>
  <div class="page-container">
    <Breadcrumb :model="breadcrumbItems" class="mb-4" />

    <div class="content-layout">
      <Card class="claim-card">
        <template #header>
          <h1 class="card-title">Create New Claim</h1>
        </template>

        <template #content>
          <form @submit.prevent="handleSubmit" class="claim-form">
            <div class="form-grid">
              <div class="form-field span-2">
                <label for="title" class="form-label required">Claim Title</label>
                <InputText
                  id="title"
                  v-model="form.title"
                  placeholder="e.g., Business Trip to New York"
                  :invalid="!!errors.title"
                  class="w-full"
                />
                <small v-if="errors.title" class="p-error">{{ errors.title }}</small>
              </div>

              <div class="form-field">
                <label for="type" class="form-label required">Claim Type</label>
                <Dropdown
                  id="type"
                  v-model="form.claim_type_id"
                  :options="claimTypes"
                  optionLabel="name"
                  optionValue="id"
                  placeholder="Select claim type"
                  :invalid="!!errors.claim_type_id"
                  class="w-full"
                />
                <small v-if="errors.claim_type_id" class="p-error">{{ errors.claim_type_id }}</small>
              </div>

              <div class="form-field">
                <label for="amount" class="form-label required">Amount ($)</label>
                <InputNumber
                  id="amount"
                  v-model="form.amount"
                  mode="currency"
                  currency="USD"
                  locale="en-US"
                  :min="0"
                  :max="balanceInfo?.remaining_balance || 999999"
                  :invalid="!!errors.amount || isBalanceExceeded"
                  class="w-full"
                  @input="handleAmountChange"
                />
                <small v-if="errors.amount" class="p-error">{{ errors.amount }}</small>
                <small v-if="isBalanceExceeded && !errors.amount" class="p-error">
                  Amount exceeds remaining balance of ${{ balanceInfo?.remaining_balance.toFixed(2) }}
                </small>
                <small v-if="balanceInfo && form.amount > 0 && !isBalanceExceeded" class="p-help">
                  Remaining after this claim: ${{ (balanceInfo.remaining_balance - form.amount).toFixed(2) }}
                </small>
              </div>

              <!-- Enhanced Balance Information Card -->
              <Card v-if="selectedClaimType && balanceInfo" class="balance-info-card span-2">
                <template #header>
                  <div class="balance-header">
                    <i class="pi pi-wallet mr-2"></i>
                    <h3>Balance Information - {{ selectedClaimType.name }}</h3>
                    <ProgressSpinner v-if="loadingBalance" style="width: 20px; height: 20px" />
                  </div>
                </template>
                <template #content>
                  <div class="balance-summary">
                    <div class="balance-overview">
                      <div class="balance-stat">
                        <label>Limit Period:</label>
                        <Tag :value="selectedClaimType.limit_timespan" severity="info" />
                      </div>
                      <div class="balance-stat">
                        <label>Next Reset:</label>
                        <span>{{ getNextResetDate() }}</span>
                      </div>
                    </div>
                  </div>
                  
                  <div class="balance-grid">
                    <div class="balance-item">
                      <label>Total Limit ({{ selectedClaimType.limit_timespan }}):</label>
                      <span class="amount total-limit">${{ balanceInfo.total_limit.toFixed(2) }}</span>
                    </div>
                    <div class="balance-item">
                      <label>Already Spent:</label>
                      <span class="amount spent">${{ balanceInfo.current_spent.toFixed(2) }}</span>
                    </div>
                    <div class="balance-item">
                      <label>Remaining Balance:</label>
                      <span class="amount remaining" :class="{ 
                        'low-balance': balanceInfo.remaining_balance < (balanceInfo.total_limit * 0.1),
                        'zero-balance': balanceInfo.remaining_balance <= 0
                      }">
                        ${{ balanceInfo.remaining_balance.toFixed(2) }}
                      </span>
                    </div>
                  </div>
                  
                  <!-- Balance Warnings -->
                  <div v-if="balanceInfo.remaining_balance <= 0" class="balance-warning danger">
                    <i class="pi pi-exclamation-triangle mr-2"></i>
                    <span>No remaining balance for this claim type</span>
                  </div>
                  <div v-else-if="balanceInfo.remaining_balance < (balanceInfo.total_limit * 0.1)" class="balance-warning warning">
                    <i class="pi pi-info-circle mr-2"></i>
                    <span>Low balance remaining ({{ ((balanceInfo.remaining_balance / balanceInfo.total_limit) * 100).toFixed(1) }}% left)</span>
                  </div>
                </template>
              </Card>

              <div class="form-field span-2">
                <label for="description" class="form-label required">Description</label>
                <Textarea
                  id="description"
                  v-model="form.description"
                  rows="5"
                  placeholder="Provide detailed information about this claim..."
                  :invalid="!!errors.description"
                  class="w-full"
                />
                <small v-if="errors.description" class="p-error">{{ errors.description }}</small>
                <small class="text-color-secondary">
                  {{ form.description.length }}/1000 characters
                </small>
              </div>

              <div class="form-field span-2">
                <label class="form-label">Attachments</label>
                <div class="attachment-area">
                  <i class="pi pi-cloud-upload"></i>
                  <p>Drag and drop files here or click to browse</p>
                  <small class="text-color-secondary">
                    Supported formats: PDF, JPG, PNG (Max 10MB)
                  </small>
                </div>
              </div>
            </div>

            <div style="padding-bottom: 2rem;"></div>

            <div class="form-actions">
              <Button
                label="Cancel"
                severity="secondary"
                outlined
                @click="router.push('/claims')"
              />
              <div class="form-actions-right">
                <Button
                  label="Save as Draft"
                  icon="pi pi-save"
                  severity="secondary"
                  @click="saveDraft"
                  :loading="saving"
                />
                <Button
                  label="Submit Claim"
                  icon="pi pi-send"
                  type="submit"
                  :loading="submitting"
                />
              </div>
            </div>
          </form>
        </template>
      </Card>

      <!-- Tips Panel -->
      <Card class="tips-card">
        <template #header>
          <h3 class="tips-title">
            <i class="pi pi-info-circle mr-2"></i>
            Tips for Faster Approval
          </h3>
        </template>
        <template #content>
          <ul class="tips-list">
            <li>
              <i class="pi pi-check text-green-500 mr-2"></i>
              Provide clear and descriptive titles
            </li>
            <li>
              <i class="pi pi-check text-green-500 mr-2"></i>
              Include all relevant receipts and documentation
            </li>
            <li>
              <i class="pi pi-check text-green-500 mr-2"></i>
              Ensure amounts match your receipts exactly
            </li>
            <li>
              <i class="pi pi-check text-green-500 mr-2"></i>
              Submit claims within 30 days of expense
            </li>
            <li>
              <i class="pi pi-check text-green-500 mr-2"></i>
              Select the most appropriate claim type
            </li>
          </ul>
        </template>
      </Card>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, computed, watch } from 'vue'
import { useRouter } from 'vue-router'
import { claimsApi, claimTypesApi, balanceApi } from '@/api'
import { useToast } from 'primevue/usetoast'
import { useAuthStore } from '@/stores/auth'
import type { ClaimType, UserClaimBalance } from '@/types'

const router = useRouter()
const toast = useToast()
const authStore = useAuthStore()

const saving = ref(false)
const submitting = ref(false)
const claimTypes = ref<ClaimType[]>([])
const balanceInfo = ref<UserClaimBalance | null>(null)
const loadingBalance = ref(false)

const form = reactive({
  title: '',
  description: '',
  amount: 0,
  claim_type_id: null as number | null
})

const errors = reactive({
  title: '',
  description: '',
  amount: '',
  claim_type_id: ''
})

const selectedClaimType = computed(() => {
  return claimTypes.value.find(type => type.id === form.claim_type_id)
})

const isBalanceExceeded = computed(() => {
  if (!balanceInfo.value || !form.amount) return false
  return form.amount > (balanceInfo.value.remaining_balance || 0)
})

const formatBalance = (amount: number | undefined) => {
  return (amount || 0).toFixed(2)
}

const getNextResetDate = () => {
  if (!balanceInfo.value) return 'N/A'
  
  const lastReset = new Date(balanceInfo.value.last_reset_date)
  const resetPeriod = balanceInfo.value.reset_period
  
  let nextReset = new Date(lastReset)
  
  switch (resetPeriod) {
    case 'daily':
      nextReset.setDate(nextReset.getDate() + 1)
      break
    case 'weekly':
      nextReset.setDate(nextReset.getDate() + 7)
      break
    case 'monthly':
      nextReset.setMonth(nextReset.getMonth() + 1)
      break
    case 'annual':
      nextReset.setFullYear(nextReset.getFullYear() + 1)
      break
  }
  
  return nextReset.toLocaleDateString()
}

const handleAmountChange = () => {
  // Real-time balance validation
  if (form.amount && balanceInfo.value) {
    if (form.amount > balanceInfo.value.remaining_balance) {
      form.amount = balanceInfo.value.remaining_balance
      toast.add({
        severity: 'warn',
        summary: 'Amount Adjusted',
        detail: `Amount automatically adjusted to remaining balance of $${balanceInfo.value.remaining_balance.toFixed(2)}`,
        life: 4000
      })
    }
  }
}

const breadcrumbItems = [
  { label: 'Claims', route: '/claims' },
  { label: 'New Claim' }
]

const loadUserBalance = async (claimTypeId: number) => {
  if (!claimTypeId) {
    balanceInfo.value = null
    return
  }

  loadingBalance.value = true
  try {
    const response = await balanceApi.getUserBalance(claimTypeId)
    if (response.data.data) {
      balanceInfo.value = response.data.data
      console.log('Successfully loaded balance info:', balanceInfo.value)
    } else {
      console.warn('No balance data returned, creating default balance')
      // Create default balance info if none exists
      const claimType = selectedClaimType.value
      if (claimType) {
        balanceInfo.value = {
          user_id: authStore.user?.id || 0,
          claim_type_id: claimTypeId,
          total_limit: claimType.limit_amount || 0,
          current_spent: 0,
          remaining_balance: claimType.limit_amount || 0,
          last_reset_date: new Date().toISOString(),
          reset_period: claimType.limit_timespan || 'annual'
        }
      }
    }
  } catch (error) {
    console.error('Failed to load balance:', error)

    // Don't logout for balance errors - provide graceful fallback
    if (error.response?.status === 401) {
      console.log('Auth error loading balance - checking token')
      if (authStore.isTokenExpired()) {
        toast.add({
          severity: 'warn',
          summary: 'Session Expired',
          detail: 'Please log in again to continue',
          life: 3000
        })
        authStore.logout()
        router.push('/login')
        return
      }
    }

    // For other errors, create default balance to allow user to continue
    const claimType = selectedClaimType.value
    if (claimType) {
      console.log('Creating default balance due to error')
      balanceInfo.value = {
        user_id: authStore.user?.id || 0,
        claim_type_id: claimTypeId,
        total_limit: claimType.limit_amount || 0,
        current_spent: 0,
        remaining_balance: claimType.limit_amount || 0,
        last_reset_date: new Date().toISOString(),
        reset_period: claimType.limit_timespan || 'annual'
      }

      toast.add({
        severity: 'info',
        summary: 'Balance Info',
        detail: 'Using default balance limits for this claim type',
        life: 4000
      })
    } else {
      balanceInfo.value = null
    }
  } finally {
    loadingBalance.value = false
  }
}

const validateAmount = async () => {
  if (!form.claim_type_id || !form.amount) return true

  try {
    const response = await balanceApi.checkClaimAmount({
      claim_type_id: form.claim_type_id,
      amount: form.amount
    })

    if (!response.data.data?.can_claim) {
      errors.amount = response.data.data?.message || `Amount exceeds remaining balance of $${(balanceInfo.value?.remaining_balance || 0).toFixed(2)}`
      return false
    }

    // Clear amount error if validation passes
    if (errors.amount.includes('balance') || errors.amount.includes('exceeds')) {
      errors.amount = ''
    }
    return true
  } catch (error) {
    errors.amount = 'Unable to validate balance'
    return false
  }
}

const validateForm = async () => {
  let isValid = true
  Object.keys(errors).forEach(key => errors[key as keyof typeof errors] = '')

  if (!form.title) {
    errors.title = 'Title is required'
    isValid = false
  } else if (form.title.length < 5) {
    errors.title = 'Title must be at least 5 characters'
    isValid = false
  }

  if (!form.claim_type_id) {
    errors.claim_type_id = 'Please select a claim type'
    isValid = false
  }

  if (!form.amount || form.amount <= 0) {
    errors.amount = 'Amount must be greater than 0'
    isValid = false
  } else {
    // Validate balance if amount is valid
    const balanceValid = await validateAmount()
    if (!balanceValid) {
      isValid = false
    }
  }

  if (!form.description) {
    errors.description = 'Description is required'
    isValid = false
  } else if (form.description.length < 10) {
    errors.description = 'Description must be at least 10 characters'
    isValid = false
  } else if (form.description.length > 1000) {
    errors.description = 'Description cannot exceed 1000 characters'
    isValid = false
  }

  return isValid
}

const saveDraft = async () => {
  if (!(await validateForm())) return

  saving.value = true
  try {
    const response = await claimsApi.create({
      ...form,
      status: 'draft'
    })

    toast.add({
      severity: 'success',
      summary: 'Success',
      detail: 'Claim saved as draft',
      life: 3000
    })

    router.push(`/claims/${response.data.data?.id}`)
  } catch (error) {
    toast.add({
      severity: 'error',
      summary: 'Error',
      detail: 'Failed to save claim',
      life: 3000
    })
  } finally {
    saving.value = false
  }
}

const handleSubmit = async () => {
  if (!(await validateForm())) return

  submitting.value = true
  try {
    const response = await claimsApi.create({
      ...form,
      status: 'draft'
    })

    if (response.data.data?.id) {
      await claimsApi.submit(response.data.data.id)

      toast.add({
        severity: 'success',
        summary: 'Success',
        detail: 'Claim submitted successfully',
        life: 3000
      })

      router.push(`/claims/${response.data.data.id}`)
    }
  } catch (error) {
    toast.add({
      severity: 'error',
      summary: 'Error',
      detail: 'Failed to submit claim',
      life: 3000
    })
  } finally {
    submitting.value = false
  }
}

const loadClaimTypes = async () => {
  // Check if token is expired before making API call
  if (authStore.isTokenExpired()) {
    console.log('Token expired, redirecting to login')
    toast.add({
      severity: 'warn',
      summary: 'Session Expired',
      detail: 'Please log in again to continue',
      life: 3000
    })
    authStore.logout()
    router.push('/login')
    return
  }

  try {
    console.log('Loading claim types...')
    const response = await claimTypesApi.getAll()
    console.log('Claim types response:', response.data)

    if (response.data.data) {
      claimTypes.value = response.data.data.filter(type => {
        return type.is_active !== false
      })
      console.log('Successfully loaded claim types:', claimTypes.value)
    } else {
      console.warn('No claim types data in response')
    }
  } catch (error) {
    console.error('Failed to load claim types:', error)

    // Check if it's an auth error
    if (error.response?.status === 401) {
      console.log('Auth error when loading claim types - token may have expired during request')
      toast.add({
        severity: 'warn',
        summary: 'Session Expired',
        detail: 'Please log in again to continue',
        life: 3000
      })
      authStore.logout()
      router.push('/login')
      return
    }

    toast.add({
      severity: 'error',
      summary: 'Error',
      detail: 'Failed to load claim types',
      life: 3000
    })
  }
}

// Watch for claim type changes to load balance
watch(() => form.claim_type_id, (newClaimTypeId) => {
  if (newClaimTypeId) {
    loadUserBalance(newClaimTypeId)
  } else {
    balanceInfo.value = null
  }
})

// Watch amount changes for real-time validation
watch(() => form.amount, async (newAmount) => {
  if (newAmount && form.claim_type_id && balanceInfo.value) {
    // Clear previous balance-related errors
    if (errors.amount.includes('balance') || errors.amount.includes('exceeds')) {
      errors.amount = ''
    }
    // Debounce validation to avoid too many API calls
    setTimeout(() => {
      if (form.amount === newAmount) { // Only validate if amount hasn't changed
        validateAmount()
      }
    }, 500)
  }
})

onMounted(() => {
  loadClaimTypes()
})
</script>

<style scoped>
.content-layout {
  display: grid;
  grid-template-columns: 1fr auto;
  gap: 2rem;
  align-items: start;
}

.claim-card {
  min-width: 0;
}

.tips-card {
  width: 320px;
  position: sticky;
  top: 1rem;
}

.card-title {
  font-size: 1.5rem;
  font-weight: 600;
  margin: 0;
}

.claim-form {
  max-width: 800px;
}

.form-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 1.5rem;
}

.span-2 {
  grid-column: span 2;
}

.form-label {
  font-weight: 500;
  color: var(--surface-700);
  margin-bottom: 0.5rem;
  display: block;
}

.form-label.required::after {
  content: ' *';
  color: var(--danger-500);
}

.attachment-area {
  border: 2px dashed var(--surface-300);
  border-radius: 8px;
  padding: 3rem;
  text-align: center;
  cursor: pointer;
  transition: all 0.3s;
}

.attachment-area:hover {
  border-color: var(--primary-500);
  background: var(--surface-50);
}

.attachment-area i {
  font-size: 3rem;
  color: var(--surface-400);
  margin-bottom: 1rem;
}

.attachment-area p {
  margin: 0.5rem 0;
  color: var(--surface-600);
}

.form-actions {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.form-actions-right {
  display: flex;
  gap: 1rem;
}

.tips-title {
  font-size: 1.125rem;
  font-weight: 600;
  margin: 0;
  display: flex;
  align-items: center;
}

.tips-list {
  list-style: none;
  padding: 0;
  margin: 0;
}

.tips-list li {
  padding: 0.75rem 0;
  border-bottom: 1px solid var(--surface-100);
  display: flex;
  align-items: center;
}

.tips-list li:last-child {
  border-bottom: none;
}

.w-full {
  width: 100%;
}

.text-green-500 {
  color: #10b981;
}

.text-color-secondary {
  color: var(--surface-500);
}

/* Enhanced Balance Information Styles */
.balance-info-card {
  border: 1px solid var(--primary-200);
  background: var(--primary-50);
}

.balance-header {
  display: flex;
  align-items: center;
  font-size: 1.125rem;
  font-weight: 600;
  color: var(--primary-700);
}

.balance-header h3 {
  margin: 0;
  flex: 1;
}

.balance-summary {
  margin-bottom: 1.5rem;
}

.balance-overview {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 1rem;
}

.balance-stat {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  font-size: 0.875rem;
}

.balance-stat label {
  color: var(--surface-600);
  font-weight: 500;
  margin: 0;
}

.balance-grid {
  display: grid;
  grid-template-columns: 1fr 1fr 1fr;
  gap: 1.5rem;
  margin-bottom: 1rem;
}

.balance-item {
  text-align: center;
  padding: 1rem;
  background: white;
  border-radius: 8px;
  border: 1px solid var(--surface-200);
}

.balance-item label {
  display: block;
  font-size: 0.875rem;
  color: var(--surface-600);
  margin-bottom: 0.5rem;
  font-weight: 500;
}

.balance-item .amount {
  display: block;
  font-size: 1.25rem;
  font-weight: 600;
  color: var(--surface-800);
}

.balance-item .amount.total-limit {
  color: var(--primary-600);
}

.balance-item .amount.spent {
  color: var(--surface-700);
}

.balance-item .amount.remaining {
  color: var(--green-600);
}

.balance-item .amount.low-balance {
  color: var(--orange-600);
}

.balance-item .amount.zero-balance {
  color: var(--red-600);
}

.balance-warning {
  display: flex;
  align-items: center;
  padding: 0.75rem 1rem;
  border-radius: 6px;
  font-weight: 500;
  margin-top: 1rem;
}

.balance-warning.danger {
  background: var(--red-100);
  border: 1px solid var(--red-300);
  color: var(--red-800);
}

.balance-warning.warning {
  background: var(--orange-100);
  border: 1px solid var(--orange-300);
  color: var(--orange-800);
}

@media (max-width: 1200px) {
  .content-layout {
    grid-template-columns: 1fr;
  }

  .tips-card {
    width: 100%;
    position: static;
    margin-top: 2rem;
  }
}

@media (max-width: 768px) {
  .form-grid {
    grid-template-columns: 1fr;
  }

  .span-2 {
    grid-column: span 1;
  }

  .balance-grid {
    grid-template-columns: 1fr;
    gap: 1rem;
  }

  .form-actions {
    flex-direction: column;
    gap: 1rem;
  }

  .form-actions-right {
    width: 100%;
  }

  .form-actions-right button {
    flex: 1;
  }
}
</style>
