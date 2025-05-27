<template>
  <div class="page-container">
    <Breadcrumb :model="breadcrumbItems" class="mb-4" />
    
    <Card v-if="claim">
      <template #header>
        <h1 class="card-title">Edit Claim</h1>
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
                :max="999999"
                :invalid="!!errors.amount"
                class="w-full"
              />
              <small v-if="errors.amount" class="p-error">{{ errors.amount }}</small>
            </div>
            
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
          </div>
          
          <Divider />
          
          <div class="form-actions">
            <Button 
              label="Cancel" 
              severity="secondary" 
              outlined
              @click="router.push(`/claims/${claimId}`)"
            />
            <div class="form-actions-right">
              <Button 
                label="Save Changes" 
                icon="pi pi-save"
                type="submit"
                :loading="saving"
              />
              <Button 
                label="Save & Submit" 
                icon="pi pi-send"
                severity="success"
                @click="saveAndSubmit"
                :loading="submitting"
              />
            </div>
          </div>
        </form>
      </template>
    </Card>
    
    <div v-else class="loading-state">
      <ProgressSpinner />
      <p>Loading claim details...</p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { claimsApi, claimTypesApi } from '@/api'
import { useToast } from 'primevue/usetoast'
import type { Claim, ClaimType } from '@/types'

const route = useRoute()
const router = useRouter()
const toast = useToast()

const claim = ref<Claim | null>(null)
const saving = ref(false)
const submitting = ref(false)
const claimTypes = ref<ClaimType[]>([])

const claimId = computed(() => Number(route.params.id))

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

const breadcrumbItems = computed(() => [
  { label: 'Claims', route: '/claims' },
  { label: claim.value?.title || 'Edit Claim' }
])

const validateForm = () => {
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

const handleSubmit = async () => {
  if (!validateForm()) return
  
  saving.value = true
  try {
    await claimsApi.update(claimId.value, form)
    
    toast.add({
      severity: 'success',
      summary: 'Success',
      detail: 'Claim updated successfully',
      life: 3000
    })
    
    router.push(`/claims/${claimId.value}`)
  } catch (error) {
    toast.add({
      severity: 'error',
      summary: 'Error',
      detail: 'Failed to update claim',
      life: 3000
    })
  } finally {
    saving.value = false
  }
}

const saveAndSubmit = async () => {
  if (!validateForm()) return
  
  submitting.value = true
  try {
    await claimsApi.update(claimId.value, form)
    await claimsApi.submit(claimId.value)
    
    toast.add({
      severity: 'success',
      summary: 'Success',
      detail: 'Claim updated and submitted for approval',
      life: 3000
    })
    
    router.push(`/claims/${claimId.value}`)
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

const loadClaim = async () => {
  try {
    const response = await claimsApi.getById(claimId.value)
    if (response.data.data) {
      claim.value = response.data.data
      
      // Populate form
      form.title = claim.value.title
      form.description = claim.value.description
      form.amount = claim.value.amount
      form.claim_type_id = claim.value.claim_type_id
      
      // Check if claim can be edited
      if (claim.value.status !== 'draft') {
        toast.add({
          severity: 'warn',
          summary: 'Cannot Edit',
          detail: 'Only draft claims can be edited',
          life: 3000
        })
        router.push(`/claims/${claimId.value}`)
      }
    }
  } catch (error) {
    toast.add({
      severity: 'error',
      summary: 'Error',
      detail: 'Failed to load claim',
      life: 3000
    })
    router.push('/claims')
  }
}

const loadClaimTypes = async () => {
  try {
    const response = await claimTypesApi.getAll()
    if (response.data.data) {
      claimTypes.value = response.data.data.filter(type => type.is_active)
    }
  } catch (error) {
    toast.add({
      severity: 'error',
      summary: 'Error',
      detail: 'Failed to load claim types',
      life: 3000
    })
  }
}

onMounted(() => {
  loadClaim()
  loadClaimTypes()
})
</script>

<style scoped>
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

.form-actions {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.form-actions-right {
  display: flex;
  gap: 1rem;
}

.loading-state {
  text-align: center;
  padding: 4rem;
}

.loading-state p {
  margin-top: 1rem;
  color: var(--surface-600);
}

.w-full {
  width: 100%;
}

.text-color-secondary {
  color: var(--surface-500);
}

@media (max-width: 768px) {
  .form-grid {
    grid-template-columns: 1fr;
  }
  
  .span-2 {
    grid-column: span 1;
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