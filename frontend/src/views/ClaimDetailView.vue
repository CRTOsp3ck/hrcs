<template>
  <div class="page-container" v-if="claim">
    <Breadcrumb :model="breadcrumbItems" class="mb-4" />
    
    <div class="detail-header">
      <div>
        <h1 class="page-title">{{ claim.title }}</h1>
        <div class="claim-meta">
          <Tag :value="claim.claim_type?.name" severity="info" />
          <span :class="`status-badge status-${claim.status}`">
            <i :class="getStatusIcon(claim.status)"></i>
            {{ formatStatus(claim.status) }}
          </span>
          <span class="claim-id">ID: #{{ claim.id }}</span>
        </div>
      </div>
      
      <div class="header-actions" v-if="canPerformActions">
        <Button 
          v-if="canEdit"
          label="Edit" 
          icon="pi pi-pencil"
          severity="secondary"
          @click="router.push(`/claims/${claim.id}/edit`)"
        />
        <Button 
          v-if="canSubmit"
          label="Submit for Approval" 
          icon="pi pi-send"
          @click="submitClaim"
          :loading="submitting"
        />
        <Button 
          v-if="canCancel"
          label="Cancel Claim" 
          icon="pi pi-times"
          severity="danger"
          outlined
          @click="confirmCancel"
        />
      </div>
    </div>

    <div class="detail-grid">
      <!-- Main Details -->
      <Card class="detail-card">
        <template #header>
          <h3 class="card-title">Claim Details</h3>
        </template>
        <template #content>
          <div class="detail-section">
            <div class="detail-row">
              <span class="detail-label">Amount:</span>
              <span class="detail-value amount">${{ formatAmount(claim.amount) }}</span>
            </div>
            <div class="detail-row">
              <span class="detail-label">Submitted by:</span>
              <span class="detail-value">{{ claim.user?.name || 'N/A' }}</span>
            </div>
            <div class="detail-row">
              <span class="detail-label">Created:</span>
              <span class="detail-value">{{ formatDateTime(claim.created_at) }}</span>
            </div>
            <div class="detail-row" v-if="claim.submitted_at">
              <span class="detail-label">Submitted:</span>
              <span class="detail-value">{{ formatDateTime(claim.submitted_at) }}</span>
            </div>
            <div class="detail-row" v-if="claim.approved_at">
              <span class="detail-label">Approved:</span>
              <span class="detail-value">{{ formatDateTime(claim.approved_at) }}</span>
            </div>
            <div class="detail-row" v-if="claim.paid_at">
              <span class="detail-label">Paid:</span>
              <span class="detail-value">{{ formatDateTime(claim.paid_at) }}</span>
            </div>
          </div>
          
          <Divider />
          
          <div class="detail-section">
            <h4 class="section-title">Description</h4>
            <p class="description">{{ claim.description }}</p>
          </div>
          
          <Divider />
          
          <div class="detail-section">
            <h4 class="section-title">Attachments</h4>
            <div class="attachments">
              <p class="text-color-secondary">No attachments uploaded</p>
            </div>
          </div>
        </template>
      </Card>

      <!-- Approval History -->
      <Card class="history-card">
        <template #header>
          <h3 class="card-title">Approval History</h3>
        </template>
        <template #content>
          <Timeline :value="approvalHistory" class="custom-timeline">
            <template #marker="slotProps">
              <span class="timeline-marker" :class="getTimelineMarkerClass(slotProps.item)">
                <i :class="getTimelineIcon(slotProps.item)"></i>
              </span>
            </template>
            <template #content="slotProps">
              <div class="timeline-content">
                <div class="timeline-header">
                  <span class="timeline-title">{{ slotProps.item.title }}</span>
                  <span class="timeline-date">{{ formatDateTime(slotProps.item.date) }}</span>
                </div>
                <p class="timeline-description">{{ slotProps.item.description }}</p>
                <p v-if="slotProps.item.comments" class="timeline-comments">
                  <i class="pi pi-comment mr-1"></i>
                  {{ slotProps.item.comments }}
                </p>
              </div>
            </template>
          </Timeline>
        </template>
      </Card>
    </div>

  </div>

  <div v-else class="page-container">
    <div class="loading-state">
      <ProgressSpinner />
      <p>Loading claim details...</p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { claimsApi } from '@/api'
import { useToast } from 'primevue/usetoast'
import { useConfirm } from 'primevue/useconfirm'
import type { Claim } from '@/types'

const route = useRoute()
const router = useRouter()
const authStore = useAuthStore()
const toast = useToast()
const confirm = useConfirm()

const claim = ref<Claim | null>(null)
const loading = ref(false)
const submitting = ref(false)

const claimId = computed(() => Number(route.params.id))

const breadcrumbItems = computed(() => [
  { label: 'Claims', route: '/claims' },
  { label: claim.value?.title || 'Claim Detail' }
])

const canEdit = computed(() => claim.value?.status === 'draft')
const canSubmit = computed(() => claim.value?.status === 'draft')
const canCancel = computed(() => 
  claim.value?.status === 'draft' || claim.value?.status === 'submitted'
)
const canPerformActions = computed(() => canEdit.value || canSubmit.value || canCancel.value)

const approvalHistory = computed(() => {
  if (!claim.value) return []
  
  const history = [
    {
      id: 1,
      type: 'created',
      title: 'Claim Created',
      description: `Created by ${claim.value.user?.name}`,
      date: claim.value.created_at
    }
  ]
  
  if (claim.value.submitted_at) {
    history.push({
      id: 2,
      type: 'submitted',
      title: 'Submitted for Approval',
      description: `Submitted by ${claim.value.user?.name}`,
      date: claim.value.submitted_at
    })
  }
  
  if (claim.value.approvals) {
    claim.value.approvals.forEach((approval, index) => {
      history.push({
        id: 3 + index,
        type: approval.action,
        title: approval.action === 'approve' ? 'Approved' : 'Rejected',
        description: `${approval.action === 'approve' ? 'Approved' : 'Rejected'} by ${approval.user?.name}`,
        date: approval.created_at,
        comments: approval.comments
      })
    })
  }
  
  return history
})

const formatAmount = (amount: number) => {
  return new Intl.NumberFormat('en-US', {
    minimumFractionDigits: 2,
    maximumFractionDigits: 2
  }).format(amount)
}

const formatDateTime = (date: string | undefined) => {
  if (!date) return ''
  return new Date(date).toLocaleString('en-US', {
    year: 'numeric',
    month: 'short',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
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

const getTimelineMarkerClass = (item: any) => {
  const classes: Record<string, string> = {
    'created': 'timeline-marker-info',
    'submitted': 'timeline-marker-primary',
    'approve': 'timeline-marker-success',
    'reject': 'timeline-marker-danger'
  }
  return classes[item.type] || ''
}

const getTimelineIcon = (item: any) => {
  const icons: Record<string, string> = {
    'created': 'pi pi-plus',
    'submitted': 'pi pi-send',
    'approve': 'pi pi-check',
    'reject': 'pi pi-times'
  }
  return icons[item.type] || 'pi pi-circle'
}

const loadClaim = async () => {
  loading.value = true
  try {
    const response = await claimsApi.getById(claimId.value)
    if (response.data.data) {
      claim.value = response.data.data
    }
  } catch (error) {
    toast.add({
      severity: 'error',
      summary: 'Error',
      detail: 'Failed to load claim details',
      life: 3000
    })
    router.push('/claims')
  } finally {
    loading.value = false
  }
}

const submitClaim = async () => {
  submitting.value = true
  try {
    await claimsApi.submit(claimId.value)
    toast.add({
      severity: 'success',
      summary: 'Success',
      detail: 'Claim submitted for approval',
      life: 3000
    })
    loadClaim()
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

const confirmCancel = () => {
  confirm.require({
    message: 'Are you sure you want to cancel this claim?',
    header: 'Cancel Confirmation',
    icon: 'pi pi-exclamation-triangle',
    acceptClass: 'p-button-danger',
    accept: cancelClaim
  })
}

const cancelClaim = async () => {
  try {
    await claimsApi.delete(claimId.value)
    toast.add({
      severity: 'success',
      summary: 'Success',
      detail: 'Claim cancelled successfully',
      life: 3000
    })
    router.push('/claims')
  } catch (error) {
    toast.add({
      severity: 'error',
      summary: 'Error',
      detail: 'Failed to cancel claim',
      life: 3000
    })
  }
}


onMounted(() => {
  loadClaim()
})
</script>

<style scoped>
.detail-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 2rem;
}

.claim-meta {
  display: flex;
  align-items: center;
  gap: 1rem;
  margin-top: 0.5rem;
}

.claim-id {
  color: var(--surface-500);
  font-family: monospace;
}

.header-actions {
  display: flex;
  gap: 1rem;
}

.detail-grid {
  display: grid;
  grid-template-columns: 1fr 400px;
  gap: 1.5rem;
}

.card-title {
  font-size: 1.125rem;
  font-weight: 600;
  margin: 0;
}

.detail-section {
  margin-bottom: 1.5rem;
}

.detail-section:last-child {
  margin-bottom: 0;
}

.section-title {
  font-size: 1rem;
  font-weight: 600;
  margin: 0 0 1rem;
  color: var(--surface-700);
}

.detail-row {
  display: flex;
  justify-content: space-between;
  padding: 0.75rem 0;
  border-bottom: 1px solid var(--surface-100);
}

.detail-row:last-child {
  border-bottom: none;
}

.detail-label {
  font-weight: 500;
  color: var(--surface-600);
}

.detail-value {
  color: var(--surface-900);
}

.amount {
  font-size: 1.25rem;
  font-weight: 600;
  color: var(--primary-600);
}

.description {
  line-height: 1.6;
  color: var(--surface-700);
}

.attachments {
  padding: 1rem;
  background: var(--surface-50);
  border-radius: 8px;
}

.custom-timeline .timeline-marker {
  width: 2rem;
  height: 2rem;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
}

.timeline-marker-info {
  background: var(--surface-500);
}

.timeline-marker-primary {
  background: var(--primary-500);
}

.timeline-marker-success {
  background: #10b981;
}

.timeline-marker-danger {
  background: #ef4444;
}

.timeline-content {
  padding-bottom: 2rem;
}

.timeline-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 0.5rem;
}

.timeline-title {
  font-weight: 600;
  color: var(--surface-900);
}

.timeline-date {
  font-size: 0.875rem;
  color: var(--surface-500);
}

.timeline-description {
  color: var(--surface-700);
  margin: 0.5rem 0;
}

.timeline-comments {
  background: var(--surface-100);
  padding: 0.75rem;
  border-radius: 6px;
  margin-top: 0.5rem;
  font-style: italic;
  color: var(--surface-700);
}

.admin-actions-card {
  margin-top: 1.5rem;
}

.admin-form {
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
}

.admin-buttons {
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

@media (max-width: 1200px) {
  .detail-grid {
    grid-template-columns: 1fr;
  }
  
  .detail-header {
    flex-direction: column;
    gap: 1rem;
  }
  
  .header-actions {
    width: 100%;
  }
  
  .header-actions button {
    flex: 1;
  }
}
</style>