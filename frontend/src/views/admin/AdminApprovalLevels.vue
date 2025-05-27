<template>
  <div class="page-container">
    <div class="page-header">
      <h1 class="page-title">Approval Levels</h1>
      <p class="page-subtitle">Configure approval levels and permissions for each user group</p>
    </div>

    <!-- Info Banner -->
    <Message severity="info" :closable="false" class="info-banner">
      <template #default>
        <div class="info-content">
          <i class="pi pi-info-circle"></i>
          <span>Each user group can have multiple approval levels (1, 2, 3, etc.) with different status permissions.</span>
        </div>
      </template>
    </Message>

    <!-- Toolbar -->
    <div class="toolbar">
      <Button label="Add Approval Level" icon="pi pi-plus" @click="showAddDialog = true" />
      <Button label="Refresh" icon="pi pi-refresh" severity="secondary" @click="loadApprovalLevels" :loading="loading" />
    </div>

    <!-- Approval Levels Table -->
    <DataTable 
      :value="approvalLevels" 
      :loading="loading"
      stripedRows
      showGridlines
      class="approval-levels-table"
    >
      <template #empty>
        <div class="empty-state">
          <i class="pi pi-sitemap"></i>
          <p>No approval levels configured</p>
        </div>
      </template>

      <Column field="user_group.name" header="User Group" sortable>
        <template #body="{ data }">
          <Tag :value="data.user_group?.name || 'N/A'" severity="info" />
        </template>
      </Column>

      <Column field="level" header="Level" sortable>
        <template #body="{ data }">
          <Badge :value="`Level ${data.level}`" severity="success" />
        </template>
      </Column>

      <Column field="approver.name" header="Approver">
        <template #body="{ data }">
          <div class="approver-info">
            <Avatar 
              :label="getInitials(data.approver)" 
              size="small"
              style="margin-right: 0.5rem"
            />
            {{ data.approver?.first_name }} {{ data.approver?.last_name }}
          </div>
        </template>
      </Column>

      <Column header="Status Permissions">
        <template #body="{ data }">
          <div class="permissions-grid">
            <Tag v-if="data.can_draft" value="Draft" severity="secondary" />
            <Tag v-if="data.can_submit" value="Submit" severity="info" />
            <Tag v-if="data.can_approve" value="Approve" severity="success" />
            <Tag v-if="data.can_reject" value="Reject" severity="danger" />
            <Tag v-if="data.can_set_payment_in_progress" value="Payment in Progress" severity="warning" />
            <Tag v-if="data.can_set_paid" value="Paid" severity="success" />
          </div>
        </template>
      </Column>

      <Column header="Actions" style="width: 150px">
        <template #body="{ data }">
          <div class="action-buttons">
            <Button 
              icon="pi pi-pencil" 
              severity="secondary" 
              text 
              @click="editLevel(data)"
              v-tooltip.top="'Edit'"
            />
            <Button 
              icon="pi pi-trash" 
              severity="danger" 
              text 
              @click="confirmDelete(data)"
              v-tooltip.top="'Delete'"
            />
          </div>
        </template>
      </Column>
    </DataTable>

    <!-- Add/Edit Dialog -->
    <Dialog 
      v-model:visible="showAddDialog" 
      :header="editingItem ? 'Edit Approval Level' : 'Add Approval Level'"
      :style="{ width: '600px' }"
      modal
    >
      <div class="form-container">
        <div class="field">
          <label for="userGroup">User Group *</label>
          <Dropdown
            id="userGroup"
            v-model="form.user_group_id"
            :options="availableGroups"
            optionLabel="name"
            optionValue="id"
            placeholder="Select a user group"
            class="w-full"
            :disabled="editingItem"
          />
        </div>

        <div class="field">
          <label for="level">Level *</label>
          <InputNumber
            id="level"
            v-model="form.level"
            :min="1"
            :max="10"
            placeholder="Enter level (1, 2, 3...)"
            class="w-full"
            :disabled="editingItem"
          />
        </div>

        <div class="field">
          <label for="approver">Approver *</label>
          <Dropdown
            id="approver"
            v-model="form.approver_id"
            :options="availableUsers"
            optionLabel="name"
            optionValue="id"
            placeholder="Select an approver"
            class="w-full"
            filter
          >
            <template #option="{ option }">
              <div class="approver-option">
                <Avatar :label="getInitials(option)" size="small" style="margin-right: 0.5rem" />
                {{ option.name }}
              </div>
            </template>
          </Dropdown>
        </div>

        <Divider />

        <div class="field">
          <label>Status Permissions</label>
          <div class="permissions-checkboxes">
            <div class="checkbox-item">
              <Checkbox 
                v-model="form.can_draft" 
                inputId="canDraft" 
                binary 
              />
              <label for="canDraft">Can set to Draft</label>
            </div>
            <div class="checkbox-item">
              <Checkbox 
                v-model="form.can_submit" 
                inputId="canSubmit" 
                binary 
              />
              <label for="canSubmit">Can Submit claims</label>
            </div>
            <div class="checkbox-item">
              <Checkbox 
                v-model="form.can_approve" 
                inputId="canApprove" 
                binary 
              />
              <label for="canApprove">Can Approve claims</label>
            </div>
            <div class="checkbox-item">
              <Checkbox 
                v-model="form.can_reject" 
                inputId="canReject" 
                binary 
              />
              <label for="canReject">Can Reject claims</label>
            </div>
            <div class="checkbox-item">
              <Checkbox 
                v-model="form.can_set_payment_in_progress" 
                inputId="canSetPaymentInProgress" 
                binary 
              />
              <label for="canSetPaymentInProgress">Can set Payment in Progress</label>
            </div>
            <div class="checkbox-item">
              <Checkbox 
                v-model="form.can_set_paid" 
                inputId="canSetPaid" 
                binary 
              />
              <label for="canSetPaid">Can mark as Paid</label>
            </div>
          </div>
        </div>
      </div>
      
      <template #footer>
        <Button label="Cancel" severity="secondary" @click="showAddDialog = false" />
        <Button label="Save" @click="saveItem" :loading="saving" />
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
        <p>Are you sure you want to delete this approval level?</p>
        <p class="text-secondary">This action cannot be undone.</p>
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
import { approvalLevelsApi, userGroupsApi, usersApi } from '@/api'
import type { ApprovalLevel, UserGroup, User } from '@/types'

const toast = useToast()

const loading = ref(false)
const saving = ref(false)
const approvalLevels = ref<ApprovalLevel[]>([])
const availableUsers = ref<User[]>([])
const availableGroups = ref<UserGroup[]>([])
const showAddDialog = ref(false)
const showDeleteDialog = ref(false)
const editingItem = ref<ApprovalLevel | null>(null)
const itemToDelete = ref<ApprovalLevel | null>(null)

const form = ref({
  user_group_id: null as number | null,
  level: 1,
  approver_id: null as number | null,
  can_draft: false,
  can_submit: false,
  can_approve: true,
  can_reject: true,
  can_set_payment_in_progress: false,
  can_set_paid: false
})

const getInitials = (user: any) => {
  if (!user) return '?'
  const firstName = user.first_name || user.firstName || ''
  const lastName = user.last_name || user.lastName || ''
  return (firstName[0] || '') + (lastName[0] || '')
}

const loadApprovalLevels = async () => {
  loading.value = true
  try {
    const response = await approvalLevelsApi.getAll()
    approvalLevels.value = response.data.data || []
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
    const [usersRes, groupsRes] = await Promise.all([
      usersApi.getAll(),
      userGroupsApi.getAll()
    ])
    
    // Format users for display
    availableUsers.value = (usersRes.data.data || []).map(user => ({
      ...user,
      name: `${user.first_name || ''} ${user.last_name || ''}`
    }))
    
    availableGroups.value = groupsRes.data.data || []
  } catch (error) {
    console.error('Failed to load resources:', error)
  }
}

const editLevel = (level: ApprovalLevel) => {
  editingItem.value = level
  form.value = {
    user_group_id: level.user_group_id,
    level: level.level,
    approver_id: level.approver_id,
    can_draft: level.can_draft,
    can_submit: level.can_submit,
    can_approve: level.can_approve,
    can_reject: level.can_reject,
    can_set_payment_in_progress: level.can_set_payment_in_progress,
    can_set_paid: level.can_set_paid
  }
  showAddDialog.value = true
}

const confirmDelete = (level: ApprovalLevel) => {
  itemToDelete.value = level
  showDeleteDialog.value = true
}

const saveItem = async () => {
  if (!form.value.user_group_id || !form.value.approver_id) {
    toast.add({
      severity: 'warn',
      summary: 'Validation Error',
      detail: 'Please fill all required fields',
      life: 3000
    })
    return
  }

  saving.value = true
  try {
    if (editingItem.value) {
      await approvalLevelsApi.update(editingItem.value.id, form.value)
      toast.add({
        severity: 'success',
        summary: 'Success',
        detail: 'Approval level updated successfully',
        life: 3000
      })
    } else {
      await approvalLevelsApi.create(form.value)
      toast.add({
        severity: 'success',
        summary: 'Success',
        detail: 'Approval level created successfully',
        life: 3000
      })
    }
    showAddDialog.value = false
    loadApprovalLevels()
    resetForm()
  } catch (error: any) {
    toast.add({
      severity: 'error',
      summary: 'Error',
      detail: error.response?.data?.message || 'Failed to save approval level',
      life: 3000
    })
  } finally {
    saving.value = false
  }
}

const deleteItem = async () => {
  if (!itemToDelete.value) return

  try {
    await approvalLevelsApi.delete(itemToDelete.value.id)
    toast.add({
      severity: 'success',
      summary: 'Success',
      detail: 'Approval level deleted successfully',
      life: 3000
    })
    showDeleteDialog.value = false
    loadApprovalLevels()
  } catch (error) {
    toast.add({
      severity: 'error',
      summary: 'Error',
      detail: 'Failed to delete approval level',
      life: 3000
    })
  }
}

const resetForm = () => {
  editingItem.value = null
  form.value = {
    user_group_id: null,
    level: 1,
    approver_id: null,
    can_draft: false,
    can_submit: false,
    can_approve: true,
    can_reject: true,
    can_set_payment_in_progress: false,
    can_set_paid: false
  }
}

onMounted(() => {
  loadApprovalLevels()
  loadResources()
})
</script>

<style scoped>
.page-container {
  padding: 1.5rem;
  max-width: 1400px;
  margin: 0 auto;
}

.page-header {
  margin-bottom: 2rem;
}

.page-title {
  font-size: 2rem;
  font-weight: 600;
  margin: 0 0 0.5rem 0;
}

.page-subtitle {
  color: var(--text-color-secondary);
  margin: 0;
}

.info-banner {
  margin-bottom: 1.5rem;
}

.info-content {
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.toolbar {
  display: flex;
  gap: 0.5rem;
  margin-bottom: 1.5rem;
}

.approval-levels-table {
  background: var(--surface-card);
  border-radius: var(--border-radius);
  overflow: hidden;
}

.empty-state {
  text-align: center;
  padding: 3rem;
  color: var(--text-color-secondary);
}

.empty-state i {
  font-size: 3rem;
  margin-bottom: 1rem;
  display: block;
  opacity: 0.3;
}

.approver-info {
  display: flex;
  align-items: center;
}

.permissions-grid {
  display: flex;
  flex-wrap: wrap;
  gap: 0.25rem;
}

.action-buttons {
  display: flex;
  gap: 0.25rem;
}

.form-container {
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
  font-weight: 500;
  color: var(--text-color);
}

.permissions-checkboxes {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 1rem;
  margin-top: 0.5rem;
}

.checkbox-item {
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.checkbox-item label {
  font-weight: normal;
  margin: 0;
  cursor: pointer;
}

.approver-option {
  display: flex;
  align-items: center;
}

.confirmation-content {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 1rem;
  text-align: center;
  padding: 1rem;
}

.text-secondary {
  color: var(--text-color-secondary);
  font-size: 0.875rem;
}

.w-full {
  width: 100%;
}
</style>