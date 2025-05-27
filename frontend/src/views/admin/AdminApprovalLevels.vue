<template>
  <div class="page-container">
    <div class="page-header">
      <h1 class="page-title">Approval Workflow</h1>
      <p class="page-subtitle">Configure approval levels and workflow rules</p>
    </div>

    <!-- Info Banner -->
    <Message severity="info" :closable="false" class="info-banner">
      <template #default>
        <div class="info-content">
          <i class="pi pi-info-circle"></i>
          <span>Approval levels determine the authorization hierarchy for different claim amounts and types. Higher amounts require more approvals.</span>
        </div>
      </template>
    </Message>

    <!-- Toolbar -->
    <div class="toolbar">
      <Button label="Add Approval Level" icon="pi pi-plus" @click="showAddDialog = true" />
    </div>

    <!-- Approval Levels List -->
    <div class="levels-container">
      <Draggable 
        v-model="approvalLevels" 
        item-key="id"
        handle=".drag-handle"
        @end="updateOrder"
      >
        <template #item="{ element, index }">
          <Card class="level-card">
            <template #content>
              <div class="level-content">
                <div class="drag-handle">
                  <i class="pi pi-bars"></i>
                </div>
                
                <div class="level-info">
                  <div class="level-header">
                    <h3 class="level-name">
                      Level {{ index + 1 }}: {{ element.name }}
                    </h3>
                    <Tag :value="`Priority ${index + 1}`" severity="info" />
                  </div>
                  
                  <div class="level-details">
                    <div class="detail-row">
                      <span class="detail-label">Amount Range:</span>
                      <span class="detail-value">
                        ${{ formatAmount(element.minAmount) }} - 
                        {{ element.maxAmount ? `$${formatAmount(element.maxAmount)}` : 'No limit' }}
                      </span>
                    </div>
                    
                    <div class="detail-row">
                      <span class="detail-label">Claim Types:</span>
                      <span class="detail-value">
                        <Tag 
                          v-for="type in element.claimTypes" 
                          :key="type"
                          :value="type"
                          severity="secondary"
                          class="type-tag"
                        />
                        <span v-if="!element.claimTypes.length">All types</span>
                      </span>
                    </div>
                    
                    <div class="detail-row">
                      <span class="detail-label">Approvers:</span>
                      <span class="detail-value">
                        <div class="approvers-list">
                          <div v-for="approver in element.approvers" :key="approver.id" class="approver-item">
                            <Avatar 
                              v-if="approver.type === 'user'"
                              :label="getInitials(approver.name)" 
                              size="small"
                            />
                            <i v-else class="pi pi-users"></i>
                            <span>{{ approver.name }}</span>
                            <Tag :value="approver.type" severity="contrast" />
                          </div>
                        </div>
                      </span>
                    </div>
                    
                    <div class="detail-row">
                      <span class="detail-label">Options:</span>
                      <span class="detail-value">
                        <div class="options-list">
                          <Tag v-if="element.requiresAllApprovers" value="All must approve" severity="warning" />
                          <Tag v-if="element.autoApprove" value="Auto-approve enabled" severity="success" />
                          <Tag v-if="element.escalationDays" :value="`Escalates after ${element.escalationDays} days`" />
                        </div>
                      </span>
                    </div>
                  </div>
                </div>
                
                <div class="level-actions">
                  <Button 
                    icon="pi pi-pencil" 
                    severity="secondary" 
                    text 
                    rounded 
                    @click="editLevel(element)"
                    v-tooltip="'Edit'"
                  />
                  <Button 
                    icon="pi pi-trash" 
                    severity="danger" 
                    text 
                    rounded 
                    @click="confirmDelete(element)"
                    v-tooltip="'Delete'"
                  />
                </div>
              </div>
            </template>
          </Card>
        </template>
      </Draggable>
    </div>

    <!-- Add/Edit Dialog -->
    <Dialog 
      v-model:visible="showAddDialog" 
      :header="editingItem ? 'Edit Approval Level' : 'Add Approval Level'"
      :style="{ width: '700px' }"
      modal
    >
      <div class="dialog-content">
        <div class="field">
          <label for="name">Level Name</label>
          <InputText id="name" v-model="form.name" class="w-full" />
        </div>
        
        <div class="field">
          <label for="description">Description</label>
          <Textarea 
            id="description" 
            v-model="form.description" 
            rows="3"
            class="w-full" 
          />
        </div>
        
        <div class="form-row">
          <div class="field">
            <label for="minAmount">Minimum Amount</label>
            <InputNumber 
              id="minAmount" 
              v-model="form.minAmount" 
              mode="currency" 
              currency="USD"
              :min="0"
              class="w-full" 
            />
          </div>
          
          <div class="field">
            <label for="maxAmount">Maximum Amount</label>
            <InputNumber 
              id="maxAmount" 
              v-model="form.maxAmount" 
              mode="currency" 
              currency="USD"
              :min="0"
              placeholder="No limit"
              class="w-full" 
            />
          </div>
        </div>
        
        <div class="field">
          <label for="claimTypes">Applicable Claim Types</label>
          <MultiSelect 
            id="claimTypes"
            v-model="form.claimTypes" 
            :options="availableClaimTypes" 
            optionLabel="name"
            optionValue="code"
            placeholder="All types"
            class="w-full"
            display="chip"
          />
        </div>
        
        <div class="field">
          <label>Approvers</label>
          <div class="approvers-section">
            <div class="approvers-grid">
              <div v-for="(approver, index) in form.approvers" :key="index" class="approver-row">
                <Dropdown 
                  v-model="approver.type" 
                  :options="approverTypes" 
                  optionLabel="label"
                  optionValue="value"
                  placeholder="Type"
                  class="approver-type"
                />
                
                <Dropdown 
                  v-if="approver.type === 'user'"
                  v-model="approver.id" 
                  :options="availableUsers" 
                  optionLabel="name"
                  optionValue="id"
                  placeholder="Select user"
                  class="approver-select"
                  filter
                />
                
                <Dropdown 
                  v-else-if="approver.type === 'group'"
                  v-model="approver.id" 
                  :options="availableGroups" 
                  optionLabel="name"
                  optionValue="id"
                  placeholder="Select group"
                  class="approver-select"
                  filter
                />
                
                <InputText 
                  v-else
                  v-model="approver.value" 
                  placeholder="e.g., Department Manager"
                  class="approver-select"
                />
                
                <Button 
                  icon="pi pi-trash" 
                  severity="danger" 
                  text 
                  rounded
                  @click="removeApprover(index)"
                />
              </div>
            </div>
            <Button 
              label="Add Approver" 
              icon="pi pi-plus" 
              severity="secondary"
              size="small"
              @click="addApprover"
            />
          </div>
        </div>
        
        <div class="field">
          <label>Approval Options</label>
          <div class="options-grid">
            <div class="option-item">
              <Checkbox v-model="form.requiresAllApprovers" inputId="requiresAll" binary />
              <label for="requiresAll">All approvers must approve</label>
            </div>
            <div class="option-item">
              <Checkbox v-model="form.autoApprove" inputId="autoApprove" binary />
              <label for="autoApprove">Auto-approve if within budget</label>
            </div>
            <div class="option-item">
              <Checkbox v-model="form.notifyApprovers" inputId="notifyApprovers" binary />
              <label for="notifyApprovers">Send notifications to approvers</label>
            </div>
          </div>
        </div>
        
        <div class="form-row">
          <div class="field">
            <label for="escalationDays">Escalation After (days)</label>
            <InputNumber 
              id="escalationDays" 
              v-model="form.escalationDays" 
              :min="0"
              suffix=" days"
              placeholder="No escalation"
              class="w-full" 
            />
          </div>
          
          <div class="field">
            <label for="reminderDays">Send Reminder After (days)</label>
            <InputNumber 
              id="reminderDays" 
              v-model="form.reminderDays" 
              :min="0"
              suffix=" days"
              placeholder="No reminders"
              class="w-full" 
            />
          </div>
        </div>
      </div>
      
      <template #footer>
        <Button label="Cancel" severity="secondary" @click="closeDialog" />
        <Button label="Save" @click="save" />
      </template>
    </Dialog>

    <!-- Delete Confirmation -->
    <Dialog 
      v-model:visible="showDeleteDialog" 
      header="Confirm Delete"
      :style="{ width: '400px' }"
      modal
    >
      <div class="confirmation-content">
        <i class="pi pi-exclamation-triangle" style="font-size: 2rem; color: var(--red-500)"></i>
        <p>Are you sure you want to delete approval level <strong>{{ itemToDelete?.name }}</strong>?</p>
        <p class="text-secondary">This may affect existing claims in the approval process.</p>
      </div>
      
      <template #footer>
        <Button label="Cancel" severity="secondary" @click="showDeleteDialog = false" />
        <Button label="Delete" severity="danger" @click="deleteItem" />
      </template>
    </Dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useToast } from 'primevue/usetoast'
import { adminApi } from '@/api'
import Draggable from 'vuedraggable'

const toast = useToast()

const loading = ref(false)
const approvalLevels = ref([])
const availableUsers = ref([])
const availableGroups = ref([])
const availableClaimTypes = ref([])
const showAddDialog = ref(false)
const showDeleteDialog = ref(false)
const editingItem = ref(null)
const itemToDelete = ref(null)

const form = ref({
  name: '',
  description: '',
  minAmount: 0,
  maxAmount: null,
  claimTypes: [],
  approvers: [{ type: 'user', id: null, value: '' }],
  requiresAllApprovers: false,
  autoApprove: false,
  notifyApprovers: true,
  escalationDays: null,
  reminderDays: null
})

const approverTypes = ref([
  { label: 'Specific User', value: 'user' },
  { label: 'User Group', value: 'group' },
  { label: 'Role', value: 'role' }
])

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
    minimumFractionDigits: 0,
    maximumFractionDigits: 0
  }).format(amount)
}

const loadApprovalLevels = async () => {
  loading.value = true
  try {
    const response = await adminApi.getApprovalLevels()
    approvalLevels.value = response.data.data
  } catch (error) {
    toast.add({
      severity: 'error',
      summary: 'Error',
      detail: 'Failed to load approval levels',
      life: 3000
    })
  } finally {
    loading.value = false
  }
}

const loadResources = async () => {
  try {
    const [usersRes, groupsRes, claimTypesRes] = await Promise.all([
      adminApi.getUsers(),
      adminApi.getGroups(),
      adminApi.getClaimTypes()
    ])
    availableUsers.value = usersRes.data.data
    availableGroups.value = groupsRes.data.data
    availableClaimTypes.value = claimTypesRes.data.data
  } catch (error) {
    console.error('Failed to load resources:', error)
  }
}

const updateOrder = async () => {
  try {
    const orders = approvalLevels.value.map((level, index) => ({
      id: level.id,
      order: index + 1
    }))
    await adminApi.updateApprovalLevelOrder(orders)
    toast.add({
      severity: 'success',
      summary: 'Success',
      detail: 'Approval level order updated',
      life: 3000
    })
  } catch (error) {
    toast.add({
      severity: 'error',
      summary: 'Error',
      detail: 'Failed to update order',
      life: 3000
    })
    loadApprovalLevels()
  }
}

const addApprover = () => {
  form.value.approvers.push({ type: 'user', id: null, value: '' })
}

const removeApprover = (index: number) => {
  form.value.approvers.splice(index, 1)
}

const editLevel = (level: any) => {
  editingItem.value = level
  form.value = {
    name: level.name,
    description: level.description,
    minAmount: level.minAmount,
    maxAmount: level.maxAmount,
    claimTypes: level.claimTypes,
    approvers: level.approvers.map((a: any) => ({ ...a })),
    requiresAllApprovers: level.requiresAllApprovers,
    autoApprove: level.autoApprove,
    notifyApprovers: level.notifyApprovers,
    escalationDays: level.escalationDays,
    reminderDays: level.reminderDays
  }
  showAddDialog.value = true
}

const confirmDelete = (level: any) => {
  itemToDelete.value = level
  showDeleteDialog.value = true
}

const deleteItem = async () => {
  try {
    await adminApi.deleteApprovalLevel(itemToDelete.value.id)
    toast.add({
      severity: 'success',
      summary: 'Success',
      detail: 'Approval level deleted successfully',
      life: 3000
    })
    loadApprovalLevels()
    showDeleteDialog.value = false
  } catch (error) {
    toast.add({
      severity: 'error',
      summary: 'Error',
      detail: 'Failed to delete approval level',
      life: 3000
    })
  }
}

const save = async () => {
  try {
    const validApprovers = form.value.approvers.filter(a => 
      (a.type === 'user' && a.id) || 
      (a.type === 'group' && a.id) || 
      (a.type === 'role' && a.value)
    )
    
    const data = {
      ...form.value,
      approvers: validApprovers
    }
    
    if (editingItem.value) {
      await adminApi.updateApprovalLevel(editingItem.value.id, data)
      toast.add({
        severity: 'success',
        summary: 'Success',
        detail: 'Approval level updated successfully',
        life: 3000
      })
    } else {
      await adminApi.createApprovalLevel(data)
      toast.add({
        severity: 'success',
        summary: 'Success',
        detail: 'Approval level created successfully',
        life: 3000
      })
    }
    loadApprovalLevels()
    closeDialog()
  } catch (error) {
    toast.add({
      severity: 'error',
      summary: 'Error',
      detail: 'Failed to save approval level',
      life: 3000
    })
  }
}

const closeDialog = () => {
  showAddDialog.value = false
  editingItem.value = null
  form.value = {
    name: '',
    description: '',
    minAmount: 0,
    maxAmount: null,
    claimTypes: [],
    approvers: [{ type: 'user', id: null, value: '' }],
    requiresAllApprovers: false,
    autoApprove: false,
    notifyApprovers: true,
    escalationDays: null,
    reminderDays: null
  }
}

onMounted(() => {
  loadApprovalLevels()
  loadResources()
})
</script>

<style scoped>
.info-banner {
  margin-bottom: 1.5rem;
}

.info-content {
  display: flex;
  align-items: center;
  gap: 0.75rem;
}

.toolbar {
  display: flex;
  justify-content: flex-end;
  margin-bottom: 1.5rem;
}

.levels-container {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.level-card {
  cursor: move;
}

.level-content {
  display: flex;
  align-items: flex-start;
  gap: 1rem;
}

.drag-handle {
  padding: 0.5rem;
  cursor: grab;
  color: var(--surface-500);
}

.drag-handle:active {
  cursor: grabbing;
}

.level-info {
  flex: 1;
}

.level-header {
  display: flex;
  align-items: center;
  gap: 1rem;
  margin-bottom: 1rem;
}

.level-name {
  margin: 0;
  font-size: 1.125rem;
  font-weight: 600;
}

.level-details {
  display: flex;
  flex-direction: column;
  gap: 0.75rem;
}

.detail-row {
  display: flex;
  gap: 0.5rem;
}

.detail-label {
  color: var(--surface-600);
  font-weight: 500;
  min-width: 120px;
}

.detail-value {
  flex: 1;
}

.type-tag {
  margin-right: 0.25rem;
}

.approvers-list {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.approver-item {
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.options-list {
  display: flex;
  flex-wrap: wrap;
  gap: 0.5rem;
}

.level-actions {
  display: flex;
  gap: 0.25rem;
}

.dialog-content {
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
}

.form-row {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 1rem;
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

.approvers-section {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.approvers-grid {
  display: flex;
  flex-direction: column;
  gap: 0.75rem;
}

.approver-row {
  display: flex;
  gap: 0.5rem;
  align-items: center;
}

.approver-type {
  width: 150px;
}

.approver-select {
  flex: 1;
}

.options-grid {
  display: grid;
  grid-template-columns: 1fr;
  gap: 0.75rem;
}

.option-item {
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.confirmation-content {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 1rem;
  text-align: center;
  padding: 1rem 0;
}

.text-secondary {
  color: var(--surface-600);
  font-size: 0.875rem;
}

@media (max-width: 768px) {
  .level-content {
    flex-direction: column;
  }
  
  .form-row {
    grid-template-columns: 1fr;
  }
  
  .approver-row {
    flex-wrap: wrap;
  }
  
  .approver-type,
  .approver-select {
    width: 100%;
  }
}
</style>