<template>
  <div class="admin-page-container">
    <div class="page-header">
      <h1 class="page-title">Claims Management</h1>
      <p class="page-subtitle">Review and manage all system claims</p>
    </div>

    <!-- Statistics Cards -->
    <div class="stats-grid">
      <Card class="stat-card">
        <template #content>
          <div class="stat-content">
            <div class="stat-icon pending">
              <i class="pi pi-clock"></i>
            </div>
            <div class="stat-details">
              <p class="stat-label">Pending</p>
              <p class="stat-value">{{ stats.pending }}</p>
            </div>
          </div>
        </template>
      </Card>

      <Card class="stat-card">
        <template #content>
          <div class="stat-content">
            <div class="stat-icon approved">
              <i class="pi pi-check-circle"></i>
            </div>
            <div class="stat-details">
              <p class="stat-label">Approved</p>
              <p class="stat-value">{{ stats.approved }}</p>
            </div>
          </div>
        </template>
      </Card>

      <Card class="stat-card">
        <template #content>
          <div class="stat-content">
            <div class="stat-icon rejected">
              <i class="pi pi-times-circle"></i>
            </div>
            <div class="stat-details">
              <p class="stat-label">Rejected</p>
              <p class="stat-value">{{ stats.rejected }}</p>
            </div>
          </div>
        </template>
      </Card>

      <Card class="stat-card">
        <template #content>
          <div class="stat-content">
            <div class="stat-icon total">
              <i class="pi pi-dollar"></i>
            </div>
            <div class="stat-details">
              <p class="stat-label">Total Amount</p>
              <p class="stat-value">${{ formatAmount(stats.totalAmount) }}</p>
            </div>
          </div>
        </template>
      </Card>
    </div>

    <!-- Filters and Search -->
    <div class="filters-section">
      <div class="filters-row">
        <IconField>
          <InputIcon class="pi pi-search" />
          <InputText v-model="filters.search" placeholder="Search claims..." />
        </IconField>

        <Dropdown
          v-model="filters.status"
          :options="statusOptions"
          optionLabel="label"
          optionValue="value"
          placeholder="All Statuses"
          showClear
        />

        <Dropdown
          v-model="filters.type"
          :options="typeOptions"
          optionLabel="label"
          optionValue="value"
          placeholder="All Types"
          showClear
        />

        <Calendar
          v-model="filters.dateRange"
          selectionMode="range"
          dateFormat="mm/dd/yy"
          placeholder="Date Range"
          showIcon
          showClear
        />

        <Button
          label="Export"
          icon="pi pi-download"
          severity="secondary"
          @click="exportClaims"
        />
      </div>
    </div>

    <!-- Claims Table -->
    <DataTable
      :value="filteredClaims"
      :loading="loading"
      paginator
      :rows="10"
      :rowsPerPageOptions="[10, 20, 50]"
      responsiveLayout="scroll"
      selectionMode="multiple"
      v-model:selection="selectedClaims"
      dataKey="id"
      :rowClass="getRowClass"
    >
      <Column selectionMode="multiple" headerStyle="width: 3rem"></Column>

      <!-- <Column field="id" header="ID" sortable style="width: 80px">
        <template #body="slotProps">
          <span class="claim-id">#{{ slotProps.data.id }}</span>
        </template>
      </Column> -->

      <Column field="employee" header="Employee" sortable>
        <template #body="slotProps">
          <div class="employee-info">
            <Avatar :label="getInitials(slotProps.data.employee)" size="small" />
            <div class="employee-details">
              <span class="employee-name">{{ slotProps.data.employee }}</span>
              <span class="employee-dept">{{ slotProps.data.department }}</span>
            </div>
          </div>
        </template>
      </Column>

      <!-- <Column field="type" header="Type" sortable>
        <template #body="slotProps">
          <Tag :value="slotProps.data.type" :severity="getTypeSeverity(slotProps.data.type)" />
        </template>
      </Column> -->

      <Column field="amount" header="Amount" sortable>
        <template #body="slotProps">
          <span class="amount">${{ formatAmount(slotProps.data.amount) }}</span>
        </template>
      </Column>

      <Column field="status" header="Status" sortable>
        <template #body="slotProps">
          <Tag :value="slotProps.data.status" :severity="getStatusSeverity(slotProps.data.status)" />
        </template>
      </Column>

      <Column field="submittedDate" header="Submitted" sortable>
        <template #body="slotProps">
          {{ formatDate(slotProps.data.submittedDate) }}
        </template>
      </Column>

      <Column field="approvalProgress" header="Approval Progress">
        <template #body="slotProps">
          <div class="approval-progress">
            <ProgressBar
              :value="getApprovalProgress(slotProps.data)"
              :showValue="false"
              style="height: 6px"
            />
            <span class="progress-text">
              {{ slotProps.data.approvalsReceived }}/{{ slotProps.data.approvalsRequired }}
            </span>
          </div>
        </template>
      </Column>

      <Column header="Review" :exportable="false" style="width: 120px">
        <template #body="slotProps">
          <div class="review-indicator">
            <Tag
              v-if="slotProps.data.canApprove && slotProps.data.allowedStatuses?.length > 0"
              value="Needs Review"
              severity="warning"
              icon="pi pi-exclamation-triangle"
              class="review-tag"
            />
            <span v-else class="no-action">-</span>
          </div>
        </template>
      </Column>

      <Column header="Actions" :exportable="false" style="width: 80px">
        <template #body="slotProps">
          <Button
            icon="pi pi-eye"
            severity="info"
            text
            rounded
            @click="viewClaim(slotProps.data)"
            v-tooltip="'View Details'"
          />
        </template>
      </Column>
    </DataTable>

    <!-- Bulk Actions -->
    <div v-if="selectedClaims.length > 0" class="bulk-actions">
      <span class="selection-count">{{ selectedClaims.length }} claims selected</span>
      <div class="actions">
        <Button
          label="Bulk Approve"
          icon="pi pi-check"
          severity="success"
          @click="bulkApprove"
        />
        <Button
          label="Bulk Reject"
          icon="pi pi-times"
          severity="danger"
          @click="bulkReject"
        />
        <Button
          label="Export Selected"
          icon="pi pi-download"
          severity="secondary"
          @click="exportSelected"
        />
      </div>
    </div>

    <!-- Claim Details Dialog -->
    <Dialog
      v-model:visible="showDetailsDialog"
      header="Claim Details"
      :style="{ width: '800px' }"
      modal
    >
      <div v-if="selectedClaim" class="claim-details">
        <div class="detail-section">
          <h3>Basic Information</h3>
          <div class="detail-grid">
            <div class="detail-item">
              <span class="label">Claim ID:</span>
              <span class="value">#{{ selectedClaim.id }}</span>
            </div>
            <div class="detail-item">
              <span class="label">Employee:</span>
              <span class="value">{{ selectedClaim.employee }}</span>
            </div>
            <div class="detail-item">
              <span class="label">Type:</span>
              <Tag :value="selectedClaim.type" :severity="getTypeSeverity(selectedClaim.type)" />
            </div>
            <div class="detail-item">
              <span class="label">Amount:</span>
              <span class="value amount">${{ formatAmount(selectedClaim.amount) }}</span>
            </div>
            <div class="detail-item">
              <span class="label">Status:</span>
              <Tag :value="selectedClaim.status" :severity="getStatusSeverity(selectedClaim.status)" />
            </div>
            <div class="detail-item">
              <span class="label">Submitted:</span>
              <span class="value">{{ formatDate(selectedClaim.submittedDate) }}</span>
            </div>
          </div>
        </div>

        <div class="detail-section">
          <h3>Description</h3>
          <p>{{ selectedClaim.description }}</p>
        </div>

        <div class="detail-section">
          <h3>Attachments</h3>
          <div class="attachments-list">
            <div v-for="attachment in selectedClaim.attachments" :key="attachment.id" class="attachment-item">
              <i class="pi pi-file"></i>
              <span>{{ attachment.name }}</span>
              <Button icon="pi pi-download" text rounded size="small" />
            </div>
          </div>
        </div>

        <div class="detail-section">
          <h3>Approval Workflow</h3>
          <div v-if="selectedClaim.approvalWorkflow && selectedClaim.approvalWorkflow.length > 0" class="workflow-container">
            <div class="workflow-summary">
              <div class="summary-item">
                <span class="summary-label">Progress:</span>
                <ProgressBar
                  :value="getWorkflowProgress(selectedClaim)"
                  :showValue="false"
                  style="width: 200px; height: 8px;"
                />
                <span class="summary-value">{{ selectedClaim.approvalsReceived }}/{{ selectedClaim.approvalsRequired }} completed</span>
              </div>
              <div v-if="selectedClaim.currentStep" class="summary-item">
                <span class="summary-label">Current Step:</span>
                <span class="summary-value">{{ selectedClaim.currentStep.name }}</span>
              </div>
            </div>

            <div class="workflow-steps">
              <div
                v-for="(step, index) in selectedClaim.approvalWorkflow"
                :key="step.id"
                class="workflow-step"
                :class="getStepClass(step)"
              >
                <div class="step-header">
                  <div class="step-number">{{ step.level }}</div>
                  <div class="step-info">
                    <h4 class="step-title">{{ step.name }}</h4>
                    <p class="step-approver">{{ step.approverName }} ({{ step.approverEmail }})</p>
                    <p class="step-group">Group: {{ step.userGroupName }}</p>
                  </div>
                  <div class="step-status">
                    <Tag
                      :value="getStepStatusLabel(step.status)"
                      :severity="getStepStatusSeverity(step.status)"
                      :icon="getStepStatusIcon(step.status)"
                    />
                  </div>
                </div>

                <div v-if="step.completedAt || step.comments" class="step-details">
                  <div v-if="step.completedAt" class="step-date">
                    <i class="pi pi-calendar"></i>
                    <span>{{ formatDateTime(step.completedAt) }}</span>
                  </div>
                  <div v-if="step.comments" class="step-comments">
                    <i class="pi pi-comment"></i>
                    <span>{{ step.comments }}</span>
                  </div>
                </div>

                <div v-if="step.status === 'pending'" class="step-permissions">
                  <div class="permissions-title">Available Actions:</div>
                  <div class="permissions-list">
                    <Tag v-if="step.permissions.canApprove" value="Approve" severity="success" size="small" />
                    <Tag v-if="step.permissions.canReject" value="Reject" severity="danger" size="small" />
                    <Tag v-if="step.permissions.canSetPaymentInProgress" value="Set Payment In Progress" severity="info" size="small" />
                    <Tag v-if="step.permissions.canSetPaid" value="Set Paid" severity="contrast" size="small" />
                  </div>
                </div>
              </div>
            </div>
          </div>
          <div v-else class="no-workflow">
            <p>No approval workflow configured for this claim.</p>
          </div>
        </div>
      </div>

      <template #footer>
        <Button label="Close" severity="secondary" @click="showDetailsDialog = false" />
        <div v-if="selectedClaim && selectedClaim.canApprove && selectedClaim.allowedStatuses?.length > 0" class="status-update-section">
          <Dropdown
            v-model="dialogSelectedStatus"
            :options="getStatusOptions(selectedClaim)"
            optionLabel="label"
            optionValue="value"
            placeholder="Change Status"
            class="status-dropdown"
          />
          <Button
            label="Update Status"
            icon="pi pi-check"
            severity="primary"
            @click="updateClaimStatusFromDialog"
            :disabled="!dialogSelectedStatus || dialogSelectedStatus === selectedClaim.status"
          />
        </div>
      </template>
    </Dialog>

    <!-- Approve/Reject Dialog -->
    <Dialog
      v-model:visible="showActionDialog"
      :header="actionType === 'approve' ? 'Approve Claim' : 'Reject Claim'"
      :style="{ width: '500px' }"
      modal
    >
      <div class="action-content">
        <p>
          Are you sure you want to {{ actionType }}
          <strong>{{ actionClaim?.id ? `claim #${actionClaim.id}` : 'selected claims' }}</strong>?
        </p>

        <div class="field">
          <label for="comments">Comments (optional)</label>
          <Textarea
            id="comments"
            v-model="actionComments"
            rows="4"
            class="w-full"
          />
        </div>
      </div>

      <template #footer>
        <Button label="Cancel" severity="secondary" @click="closeActionDialog" />
        <Button
          :label="actionType === 'approve' ? 'Approve' : 'Reject'"
          :severity="actionType === 'approve' ? 'success' : 'danger'"
          @click="performAction"
        />
      </template>
    </Dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useToast } from 'primevue/usetoast'
import { adminApi } from '@/api'

const toast = useToast()

const loading = ref(false)
const claims = ref<any[]>([])
const selectedClaims = ref<any[]>([])
const showDetailsDialog = ref(false)
const showActionDialog = ref(false)
const selectedClaim = ref<any>(null)
const actionClaim = ref(null)
const actionType = ref('')
const actionComments = ref('')
const dialogSelectedStatus = ref('')

const stats = ref({
  pending: 0,
  approved: 0,
  rejected: 0,
  totalAmount: 0
})

const filters = ref({
  search: '',
  status: null,
  type: null,
  dateRange: null
})

const statusOptions = ref([
  { label: 'Pending', value: 'pending' },
  { label: 'In Progress', value: 'in-progress' },
  { label: 'Approved', value: 'approved' },
  { label: 'Rejected', value: 'rejected' },
  { label: 'Paid', value: 'paid' }
])

const typeOptions = ref([
  { label: 'Travel', value: 'travel' },
  { label: 'Medical', value: 'medical' },
  { label: 'Equipment', value: 'equipment' },
  { label: 'Training', value: 'training' },
  { label: 'Entertainment', value: 'entertainment' },
  { label: 'Other', value: 'other' }
])

const filteredClaims = computed(() => {
  let result = claims.value

  if (filters.value.search) {
    const search = filters.value.search.toLowerCase()
    result = result.filter((claim: any) =>
      claim.id.toString().includes(search) ||
      claim.employee.toLowerCase().includes(search) ||
      claim.description.toLowerCase().includes(search)
    )
  }

  if (filters.value.status) {
    result = result.filter((claim: any) => claim.status === filters.value.status)
  }

  if (filters.value.type) {
    result = result.filter((claim: any) => claim.type === filters.value.type)
  }

  if (filters.value.dateRange && filters.value.dateRange[0]) {
    const [start, end] = filters.value.dateRange
    result = result.filter((claim: any) => {
      const date = new Date(claim.submittedDate || claim.created_at)
      return date >= start && (!end || date <= end)
    })
  }

  return result
})

const getInitials = (name: string) => {
  return name
    .split(' ')
    .map(n => n[0])
    .join('')
    .toUpperCase()
    .slice(0, 2)
}

const formatAmount = (amount: number) => {
  return new Intl.NumberFormat('en-US', {
    minimumFractionDigits: 2,
    maximumFractionDigits: 2
  }).format(amount)
}

const formatDate = (date: string) => {
  return new Date(date).toLocaleDateString('en-US', {
    year: 'numeric',
    month: 'short',
    day: 'numeric'
  })
}

const formatDateTime = (date: string) => {
  return new Date(date).toLocaleString('en-US', {
    year: 'numeric',
    month: 'short',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  })
}

const getStatusSeverity = (status: string) => {
  const severities: Record<string, string> = {
    'pending': 'warning',
    'in-progress': 'info',
    'approved': 'success',
    'rejected': 'danger',
    'paid': 'contrast'
  }
  return severities[status] || 'secondary'
}

const getTypeSeverity = (type: string) => {
  const severities: Record<string, string> = {
    travel: 'info',
    medical: 'danger',
    equipment: 'warning',
    training: 'success',
    entertainment: 'secondary',
    other: 'contrast'
  }
  return severities[type] || 'secondary'
}

const getRowClass = (data: any) => {
  if (data.priority === 'high') return 'high-priority'
  if (data.status === 'rejected') return 'rejected-row'
  return ''
}

const getApprovalProgress = (claim: any) => {
  return (claim.approvalsReceived / claim.approvalsRequired) * 100
}

const getWorkflowProgress = (claim: any) => {
  if (!claim.approvalsRequired || claim.approvalsRequired === 0) return 0
  return (claim.approvalsReceived / claim.approvalsRequired) * 100
}

const getStepClass = (step: any) => {
  return {
    'step-completed': step.status === 'approved',
    'step-rejected': step.status === 'rejected',
    'step-pending': step.status === 'pending',
    'step-current': step.status === 'pending' // Could be enhanced to show actual current step
  }
}

const getStepStatusLabel = (status: string) => {
  const labels: Record<string, string> = {
    'pending': 'Pending',
    'approved': 'Approved',
    'rejected': 'Rejected',
    'skipped': 'Skipped'
  }
  return labels[status] || status
}

const getStepStatusSeverity = (status: string) => {
  const severities: Record<string, string> = {
    'pending': 'warning',
    'approved': 'success',
    'rejected': 'danger',
    'skipped': 'secondary'
  }
  return severities[status] || 'secondary'
}

const getStepStatusIcon = (status: string) => {
  const icons: Record<string, string> = {
    'pending': 'pi pi-clock',
    'approved': 'pi pi-check',
    'rejected': 'pi pi-times',
    'skipped': 'pi pi-forward'
  }
  return icons[status] || 'pi pi-circle'
}

const canApprove = (claim: any) => {
  return claim.canApprove && claim.allowedStatuses?.length > 0
}

const canReject = (claim: any) => {
  return claim.canApprove && claim.allowedStatuses?.includes('rejected')
}

const getClaimStatus = (claim: any) => {
  return claim.status
}

const getStatusOptions = (claim: any) => {
  if (!claim.allowedStatuses) return []

  const statusLabels: Record<string, string> = {
    'draft': 'Draft',
    'submitted': 'Submitted',
    'approved': 'Approved',
    'rejected': 'Rejected',
    'payment-in-progress': 'Payment in Progress',
    'paid': 'Paid'
  }

  return claim.allowedStatuses.map((status: string) => ({
    label: statusLabels[status] || status,
    value: status
  }))
}

const updateClaimStatus = async (claim: any, newStatus: string) => {
  if (newStatus === claim.status) return

  try {
    await adminApi.updateClaimStatus(claim.id, {
      status: newStatus,
      comments: '' // TODO: Add comments dialog if needed
    })

    toast.add({
      severity: 'success',
      summary: 'Success',
      detail: `Claim status updated to ${newStatus}`,
      life: 3000
    })

    loadClaims()
  } catch (error) {
    toast.add({
      severity: 'error',
      summary: 'Error',
      detail: 'Failed to update claim status',
      life: 3000
    })
  }
}

const updateClaimStatusFromDialog = async () => {
  if (!selectedClaim.value || !dialogSelectedStatus.value) return

  try {
    await adminApi.updateClaimStatus(selectedClaim.value.id, {
      status: dialogSelectedStatus.value,
      comments: ''
    })

    toast.add({
      severity: 'success',
      summary: 'Success',
      detail: `Claim status updated to ${dialogSelectedStatus.value}`,
      life: 3000
    })

    showDetailsDialog.value = false
    dialogSelectedStatus.value = ''
    loadClaims()
  } catch (error) {
    toast.add({
      severity: 'error',
      summary: 'Error',
      detail: 'Failed to update claim status',
      life: 3000
    })
  }
}

const getApprovalClass = (action: string) => {
  const classes: Record<string, string> = {
    'submitted': 'timeline-primary',
    'approved': 'timeline-success',
    'rejected': 'timeline-danger',
    'commented': 'timeline-info'
  }
  return classes[action] || ''
}

const getApprovalIcon = (action: string) => {
  const icons: Record<string, string> = {
    'submitted': 'pi pi-send',
    'approved': 'pi pi-check',
    'rejected': 'pi pi-times',
    'commented': 'pi pi-comment'
  }
  return icons[action] || 'pi pi-circle'
}

const loadClaims = async () => {
  loading.value = true
  try {
    const response = await adminApi.getAllClaims()
    claims.value = response.data.data || []
    updateStats()
  } catch (error) {
    toast.add({
      severity: 'error',
      summary: 'Error',
      detail: 'Failed to load claims',
      life: 3000
    })
  } finally {
    loading.value = false
  }
}

const updateStats = () => {
  stats.value = {
    pending: claims.value.filter((c: any) => c.status === 'pending').length,
    approved: claims.value.filter((c: any) => c.status === 'approved').length,
    rejected: claims.value.filter((c: any) => c.status === 'rejected').length,
    totalAmount: claims.value.reduce((sum: number, c: any) => sum + c.amount, 0)
  }
}

const viewClaim = (claim: any) => {
  selectedClaim.value = claim
  dialogSelectedStatus.value = claim.status
  showDetailsDialog.value = true
}


const bulkApprove = () => {
  actionClaim.value = null
  actionType.value = 'approve'
  showActionDialog.value = true
}

const bulkReject = () => {
  actionClaim.value = null
  actionType.value = 'reject'
  showActionDialog.value = true
}

const performAction = async () => {
  try {
    const claimsToProcess = actionClaim.value ? [actionClaim.value] : selectedClaims.value

    for (const claim of claimsToProcess) {
      if (actionType.value === 'approve') {
        await adminApi.approveClaim(claim.id, { comments: actionComments.value })
      } else {
        await adminApi.rejectClaim(claim.id, { comments: actionComments.value })
      }
    }

    toast.add({
      severity: 'success',
      summary: 'Success',
      detail: `${claimsToProcess.length} claim(s) ${actionType.value}d successfully`,
      life: 3000
    })

    loadClaims()
    closeActionDialog()
    showDetailsDialog.value = false
    selectedClaims.value = []
  } catch (error) {
    toast.add({
      severity: 'error',
      summary: 'Error',
      detail: `Failed to ${actionType.value} claims`,
      life: 3000
    })
  }
}

const closeActionDialog = () => {
  showActionDialog.value = false
  actionClaim.value = null
  actionType.value = ''
  actionComments.value = ''
}

const exportClaims = () => {
  toast.add({
    severity: 'info',
    summary: 'Export',
    detail: 'Export functionality not implemented yet',
    life: 3000
  })
}

const exportSelected = () => {
  toast.add({
    severity: 'info',
    summary: 'Export',
    detail: `Exporting ${selectedClaims.value.length} claims...`,
    life: 3000
  })
}

onMounted(() => {
  loadClaims()
})
</script>

<style scoped>
.stats-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 1rem;
  margin-bottom: 2rem;
}

.stat-card {
  position: relative;
  overflow: hidden;
}

.stat-content {
  display: flex;
  align-items: center;
  gap: 1rem;
}

.stat-icon {
  width: 48px;
  height: 48px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 1.5rem;
}

.stat-icon.pending {
  background: rgba(251, 146, 60, 0.1);
  color: #fb923c;
}

.stat-icon.approved {
  background: rgba(34, 197, 94, 0.1);
  color: #22c55e;
}

.stat-icon.rejected {
  background: rgba(239, 68, 68, 0.1);
  color: #ef4444;
}

.stat-icon.total {
  background: rgba(168, 85, 247, 0.1);
  color: #a855f7;
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
  font-size: 1.5rem;
  font-weight: 700;
  color: var(--surface-900);
  margin: 0;
}

.filters-section {
  margin-bottom: 1.5rem;
}

.filters-row {
  display: flex;
  gap: 1rem;
  align-items: center;
  flex-wrap: wrap;
}

.claim-id {
  font-weight: 600;
  color: var(--primary-color);
}

.employee-info {
  display: flex;
  align-items: center;
  gap: 0.75rem;
}

.employee-details {
  display: flex;
  flex-direction: column;
}

.employee-name {
  font-weight: 500;
}

.employee-dept {
  font-size: 0.875rem;
  color: var(--surface-600);
}

.amount {
  font-weight: 600;
}

.approval-progress {
  display: flex;
  flex-direction: column;
  gap: 0.25rem;
}

.progress-text {
  font-size: 0.75rem;
  color: var(--surface-600);
}

.bulk-actions {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 1rem;
  background: var(--surface-50);
  border-radius: 6px;
  margin-top: 1rem;
}

.selection-count {
  font-weight: 500;
}

.actions {
  display: flex;
  gap: 0.5rem;
}

.claim-details {
  display: flex;
  flex-direction: column;
  gap: 2rem;
}

.detail-section h3 {
  margin: 0 0 1rem;
  font-size: 1rem;
  font-weight: 600;
  color: var(--surface-800);
}

.detail-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 1rem;
}

.detail-item {
  display: flex;
  gap: 0.5rem;
}

.detail-item .label {
  font-weight: 500;
  color: var(--surface-600);
}

.attachments-list {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.attachment-item {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.5rem;
  background: var(--surface-50);
  border-radius: 4px;
}

.timeline-marker {
  width: 2rem;
  height: 2rem;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
}

.timeline-primary {
  background: #3b82f6;
}

.timeline-success {
  background: #22c55e;
}

.timeline-danger {
  background: #ef4444;
}

.timeline-info {
  background: #6366f1;
}

.timeline-content {
  padding-bottom: 1.5rem;
}

.timeline-title {
  margin: 0;
  font-weight: 500;
}

.timeline-subtitle {
  font-size: 0.875rem;
  color: var(--surface-600);
  margin: 0.25rem 0;
}

.timeline-comments {
  font-size: 0.875rem;
  margin: 0.5rem 0 0;
  padding: 0.5rem;
  background: var(--surface-50);
  border-radius: 4px;
}

.action-content {
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
}

.field {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.field label {
  font-weight: 600;
  color: var(--surface-700);
}

.high-priority {
  background-color: rgba(251, 146, 60, 0.05);
}

.rejected-row {
  opacity: 0.7;
}

.review-indicator {
  display: flex;
  align-items: center;
  justify-content: center;
}

.review-tag {
  font-size: 0.75rem;
  font-weight: 600;
}

.no-action {
  color: var(--surface-400);
  font-style: italic;
}

.status-update-section {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  margin-left: auto;
}

/* Workflow Styles */
.workflow-container {
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
}

.workflow-summary {
  display: flex;
  flex-direction: column;
  gap: 0.75rem;
  padding: 1rem;
  background: var(--surface-50);
  border-radius: 6px;
  border-left: 4px solid var(--primary-color);
}

.summary-item {
  display: flex;
  align-items: center;
  gap: 0.75rem;
}

.summary-label {
  font-weight: 600;
  color: var(--surface-700);
  min-width: 80px;
}

.summary-value {
  color: var(--surface-900);
  font-weight: 500;
}

.workflow-steps {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.workflow-step {
  border: 1px solid var(--surface-200);
  border-radius: 8px;
  padding: 1rem;
  background: var(--surface-0);
  transition: all 0.2s ease;
}

.workflow-step.step-completed {
  border-color: var(--green-300);
  background: var(--green-50);
}

.workflow-step.step-rejected {
  border-color: var(--red-300);
  background: var(--red-50);
}

.workflow-step.step-pending {
  border-color: var(--orange-300);
  background: var(--orange-50);
}

.workflow-step.step-current {
  border-color: var(--primary-color);
  background: var(--primary-50);
  box-shadow: 0 0 0 2px rgba(var(--primary-color-rgb), 0.1);
}

.step-header {
  display: flex;
  align-items: flex-start;
  gap: 1rem;
}

.step-number {
  width: 32px;
  height: 32px;
  border-radius: 50%;
  background: var(--primary-color);
  color: white;
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: 600;
  flex-shrink: 0;
}

.step-completed .step-number {
  background: var(--green-500);
}

.step-rejected .step-number {
  background: var(--red-500);
}

.step-pending .step-number {
  background: var(--orange-500);
}

.step-info {
  flex: 1;
}

.step-title {
  margin: 0 0 0.25rem;
  font-size: 1rem;
  font-weight: 600;
  color: var(--surface-900);
}

.step-approver {
  margin: 0 0 0.25rem;
  font-size: 0.875rem;
  color: var(--surface-700);
}

.step-group {
  margin: 0;
  font-size: 0.75rem;
  color: var(--surface-600);
}

.step-status {
  flex-shrink: 0;
}

.step-details {
  margin-top: 0.75rem;
  padding-top: 0.75rem;
  border-top: 1px solid var(--surface-200);
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.step-date, .step-comments {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  font-size: 0.875rem;
  color: var(--surface-700);
}

.step-permissions {
  margin-top: 0.75rem;
  padding-top: 0.75rem;
  border-top: 1px solid var(--surface-200);
}

.permissions-title {
  font-size: 0.875rem;
  font-weight: 600;
  color: var(--surface-700);
  margin-bottom: 0.5rem;
}

.permissions-list {
  display: flex;
  flex-wrap: wrap;
  gap: 0.5rem;
}

.no-workflow {
  padding: 2rem;
  text-align: center;
  color: var(--surface-600);
  background: var(--surface-50);
  border-radius: 6px;
  border: 2px dashed var(--surface-200);
}

@media (max-width: 1200px) {
  .detail-grid {
    grid-template-columns: 1fr;
  }
}

@media (max-width: 768px) {
  .stats-grid {
    grid-template-columns: 1fr;
  }

  .filters-row {
    flex-direction: column;
    align-items: stretch;
  }

  .bulk-actions {
    flex-direction: column;
    gap: 1rem;
    align-items: stretch;
  }

  .actions {
    flex-direction: column;
  }
}
</style>
