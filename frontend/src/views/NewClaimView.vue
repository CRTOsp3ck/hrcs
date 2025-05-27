<template>
  <div class="page-container">
    <Breadcrumb :model="breadcrumbItems" class="mb-4" />

    <Card>
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

          <Divider />

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
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { claimsApi, claimTypesApi } from '@/api'
import { useToast } from 'primevue/usetoast'
import type { ClaimType } from '@/types'

const router = useRouter()
const toast = useToast()

const saving = ref(false)
const submitting = ref(false)
const claimTypes = ref<ClaimType[]>([])

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

const breadcrumbItems = [
  { label: 'Claims', route: '/claims' },
  { label: 'New Claim' }
]

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

const saveDraft = async () => {
  if (!validateForm()) return

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
  if (!validateForm()) return

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
  try {
    const response = await claimTypesApi.getAll()
    // console.log('Claim types response', response.data.data)
    if (response.data.data) {
      // claimTypes.value = response.data.data.filter(type => type.is_active) - will implement is active later
      claimTypes.value = response.data.data
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

.tips-card {
  margin-top: 2rem;
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
