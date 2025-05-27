<template>
  <div class="page-container">
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
        <span class="p-input-icon-left">
          <i class="pi pi-search" />
          <InputText v-model="filters.search" placeholder="Search claims..." />
        </span>
        
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
      
      <Column field="id" header="ID" sortable style="width: 80px">
        <template #body="slotProps">
          <span class="claim-id">#{{ slotProps.data.id }}</span>
        </template>
      </Column>
      
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
      
      <Column field="type" header="Type" sortable>
        <template #body="slotProps">
          <Tag :value="slotProps.data.type" :severity="getTypeSeverity(slotProps.data.type)" />
        </template>
      </Column>
      
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
      
      <Column header="Actions" :exportable="false" style="width: 200px">
        <template #body="slotProps">
          <div class="actions-row">
            <Button 
              icon="pi pi-eye" 
              severity="info" 
              text 
              rounded 
              @click="viewClaim(slotProps.data)"
              v-tooltip="'View Details'"
            />
            <Dropdown 
              v-if="slotProps.data.canApprove && slotProps.data.allowedStatuses?.length > 0"
              :modelValue="slotProps.data.status"
              :options="getStatusOptions(slotProps.data)"
              optionLabel="label"
              optionValue="value"
              placeholder="Change Status"
              @update:modelValue="updateClaimStatus(slotProps.data, $event)"
              class="status-dropdown"
            />
          </div>
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
          <h3>Approval History</h3>
          <Timeline :value="selectedClaim.approvalHistory" class="approval-timeline">
            <template #marker="slotProps">
              <span class="timeline-marker" :class="getApprovalClass(slotProps.item.action)">
                <i :class="getApprovalIcon(slotProps.item.action)"></i>
              </span>
            </template>
            <template #content="slotProps">
              <div class="timeline-content">
                <p class="timeline-title">{{ slotProps.item.title }}</p>
                <p class="timeline-subtitle">{{ slotProps.item.user }} - {{ formatDateTime(slotProps.item.timestamp) }}</p>
                <p v-if="slotProps.item.comments" class="timeline-comments">{{ slotProps.item.comments }}</p>
              </div>
            </template>
          </Timeline>
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

.actions-row {
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.status-dropdown {
  min-width: 120px;
}

.status-update-section {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  margin-left: auto;
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