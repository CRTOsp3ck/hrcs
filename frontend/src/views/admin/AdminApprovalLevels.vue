<template>
  <div class="page-container">
    <div class="page-header">
      <h1 class="page-title">Approval Workflow</h1>
      <p class="page-subtitle">Configure approval flow for each user group</p>
    </div>

    <!-- Info Banner -->
    <Message severity="info" :closable="false" class="info-banner">
      <template #default>
        <div class="info-content">
          <i class="pi pi-info-circle"></i>
          <span>Each user group has its own approval flow. Add approvers to create a sequential approval chain (1st → 2nd → 3rd approver, etc.)</span>
        </div>
      </template>
    </Message>

    <!-- Toolbar -->
    <div class="toolbar">
      <Button label="Add Approver" icon="pi pi-plus" @click="showAddDialog = true" />
      <Button label="Refresh" icon="pi pi-refresh" severity="secondary" @click="loadData" :loading="loading" />
    </div>

    <!-- User Groups with Approval Flows -->
    <div v-if="groupedApprovalLevels.length > 0" class="groups-container">
      <Card v-for="group in groupedApprovalLevels" :key="group.id" class="group-card">
        <template #header>
          <div class="group-header">
            <div class="group-info">
              <h3 class="group-name">{{ group.name }}</h3>
              <Tag :value="`${group.levels.length} approver${group.levels.length !== 1 ? 's' : ''}`" severity="info" />
            </div>
            <Button 
              label="Add Approver" 
              icon="pi pi-plus" 
              size="small"
              @click="addApproverToGroup(group)"
            />
          </div>
        </template>
        <template #content>
          <div v-if="group.levels.length === 0" class="empty-group">
            <i class="pi pi-users"></i>
            <p>No approvers configured for this group</p>
          </div>
          <div v-else class="approval-flow">
            <div v-for="(level, index) in group.levels" :key="level.id" class="approval-level">
              <div class="level-number">
                <Badge :value="`${index + 1}`" severity="success" size="large" />
                <div class="level-label">{{ getOrdinal(index + 1) }} Approver</div>
              </div>
              <div class="level-details">
                <div class="approver-info">
                  <Avatar 
                    :label="getInitials(level.approver)" 
                    size="large"
                  />
                  <div class="approver-name">
                    <strong>{{ level.approver?.first_name }} {{ level.approver?.last_name }}</strong>
                    <div class="approver-email">{{ level.approver?.email }}</div>
                  </div>
                </div>
                <div class="permissions">
                  <div class="permission-label">Can perform:</div>
                  <div class="permissions-grid">
                    <Tag v-if="level.can_draft" value="Draft" severity="secondary" />
                    <Tag v-if="level.can_submit" value="Submit" severity="info" />
                    <Tag v-if="level.can_approve" value="Approve" severity="success" />
                    <Tag v-if="level.can_reject" value="Reject" severity="danger" />
                    <Tag v-if="level.can_set_payment_in_progress" value="Payment in Progress" severity="warning" />
                    <Tag v-if="level.can_set_paid" value="Paid" severity="success" />
                  </div>
                </div>
                <div class="level-actions">
                  <Button 
                    icon="pi pi-pencil" 
                    severity="secondary" 
                    text 
                    rounded
                    @click="editLevel(level)"
                    v-tooltip.top="'Edit permissions'"
                  />
                  <Button 
                    icon="pi pi-trash" 
                    severity="danger" 
                    text 
                    rounded
                    @click="confirmDelete(level)"
                    v-tooltip.top="'Remove approver'"
                  />
                </div>
              </div>
              <div v-if="index < group.levels.length - 1" class="flow-arrow">
                <i class="pi pi-arrow-down"></i>
              </div>
            </div>
          </div>
        </template>
      </Card>
    </div>

    <!-- Empty State -->
    <div v-else-if="!loading" class="empty-state">
      <i class="pi pi-sitemap"></i>
      <h3>No User Groups Found</h3>
      <p>Create user groups first to configure approval workflows</p>
    </div>

    <!-- Add/Edit Dialog -->
    <Dialog 
      v-model:visible="showAddDialog" 
      :header="editingItem ? 'Edit Approver Permissions' : 'Add Approver'"
      :style="{ width: '600px' }"
      modal
    >
      <div class="form-container">
        <div v-if="!editingItem" class="field">
          <label for="userGroup">User Group *</label>
          <Dropdown
            id="userGroup"
            v-model="form.user_group_id"
            :options="availableGroups"
            optionLabel="name"
            optionValue="id"
            placeholder="Select a user group"
            class="w-full"
          />
        </div>

        <div v-if="!editingItem" class="field">
          <label for="approver">Approver *</label>
          <Dropdown
            id="approver"
            v-model="form.approver_id"
            :options="availableUsersForGroup"
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

        <div v-if="editingItem" class="edit-info">
          <div class="info-row">
            <strong>User Group:</strong> {{ editingItem.user_group?.name }}
          </div>
          <div class="info-row">
            <strong>Approver:</strong> {{ editingItem.approver?.first_name }} {{ editingItem.approver?.last_name }}
          </div>
          <div class="info-row">
            <strong>Level:</strong> {{ getOrdinal(editingItem.level) }} Approver
          </div>
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
        <Button label="Cancel" severity="secondary" @click="closeDialog" />
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
        <p>Are you sure you want to remove this approver?</p>
        <p class="text-secondary">This will update the approval flow for the group.</p>
      </div>
      
      <template #footer>
        <Button label="Cancel" severity="secondary" @click="showDeleteDialog = false" />
        <Button label="Delete" severity="danger" @click="deleteItem" />
      </template>
    </Dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
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
const selectedGroupId = ref<number | null>(null)

const form = ref({
  user_group_id: null as number | null,
  approver_id: null as number | null,
  can_draft: false,
  can_submit: true,
  can_approve: true,
  can_reject: true,
  can_set_payment_in_progress: false,
  can_set_paid: false
})

// Group approval levels by user group
const groupedApprovalLevels = computed(() => {
  const groups: any[] = []
  const groupMap = new Map()

  // First, create entries for all user groups
  availableGroups.value.forEach(group => {
    groupMap.set(group.id, {
      id: group.id,
      name: group.name,
      description: group.description,
      levels: []
    })
  })

  // Then, add approval levels to their respective groups
  approvalLevels.value.forEach(level => {
    if (groupMap.has(level.user_group_id)) {
      groupMap.get(level.user_group_id).levels.push(level)
    }
  })

  // Convert to array and sort levels within each group
  groupMap.forEach(group => {
    group.levels.sort((a: ApprovalLevel, b: ApprovalLevel) => a.level - b.level)
    groups.push(group)
  })

  return groups.sort((a, b) => a.name.localeCompare(b.name))
})

// Filter available users based on selected group
const availableUsersForGroup = computed(() => {
  if (!form.value.user_group_id) return availableUsers.value

  // Get users already assigned as approvers for this group
  const assignedUserIds = approvalLevels.value
    .filter(level => level.user_group_id === form.value.user_group_id)
    .map(level => level.approver_id)

  // Return users not already assigned
  return availableUsers.value.filter(user => !assignedUserIds.includes(user.id))
})

const getInitials = (user: any) => {
  if (!user) return '?'
  const firstName = user.first_name || user.firstName || ''
  const lastName = user.last_name || user.lastName || ''
  return (firstName[0] || '') + (lastName[0] || '')
}

const getOrdinal = (num: number) => {
  const suffixes = ['th', 'st', 'nd', 'rd']
  const v = num % 100
  return num + (suffixes[(v - 20) % 10] || suffixes[v] || suffixes[0])
}

const loadData = async () => {
  loading.value = true
  try {
    const [levelsRes, usersRes, groupsRes] = await Promise.all([
      approvalLevelsApi.getAll(),
      usersApi.getAll(),
      userGroupsApi.getAll()
    ])
    
    approvalLevels.value = levelsRes.data.data || []
    
    // Format users for display
    availableUsers.value = (usersRes.data.data || []).map(user => ({
      ...user,
      name: `${user.first_name || ''} ${user.last_name || ''}`
    }))
    
    availableGroups.value = groupsRes.data.data || []
  } catch (error) {
    toast.add({
      severity: 'error',
      summary: 'Error',
      detail: 'Failed to load data',
      life: 3000
    })
  } finally {
    loading.value = false
  }
}

const addApproverToGroup = (group: any) => {
  selectedGroupId.value = group.id
  form.value.user_group_id = group.id
  showAddDialog.value = true
}

const editLevel = (level: ApprovalLevel) => {
  editingItem.value = level
  form.value = {
    user_group_id: level.user_group_id,
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

const closeDialog = () => {
  showAddDialog.value = false
  resetForm()
}

const saveItem = async () => {
  if (!editingItem.value && (!form.value.user_group_id || !form.value.approver_id)) {
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
      // Only update permissions for existing level
      await approvalLevelsApi.update(editingItem.value.id, {
        can_draft: form.value.can_draft,
        can_submit: form.value.can_submit,
        can_approve: form.value.can_approve,
        can_reject: form.value.can_reject,
        can_set_payment_in_progress: form.value.can_set_payment_in_progress,
        can_set_paid: form.value.can_set_paid
      })
      toast.add({
        severity: 'success',
        summary: 'Success',
        detail: 'Approver permissions updated successfully',
        life: 3000
      })
    } else {
      await approvalLevelsApi.create(form.value)
      toast.add({
        severity: 'success',
        summary: 'Success',
        detail: 'Approver added successfully',
        life: 3000
      })
    }
    closeDialog()
    loadData()
  } catch (error: any) {
    toast.add({
      severity: 'error',
      summary: 'Error',
      detail: error.response?.data?.message || 'Failed to save',
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
      detail: 'Approver removed successfully',
      life: 3000
    })
    showDeleteDialog.value = false
    loadData()
  } catch (error) {
    toast.add({
      severity: 'error',
      summary: 'Error',
      detail: 'Failed to remove approver',
      life: 3000
    })
  }
}

const resetForm = () => {
  editingItem.value = null
  selectedGroupId.value = null
  form.value = {
    user_group_id: null,
    approver_id: null,
    can_draft: false,
    can_submit: true,
    can_approve: true,
    can_reject: true,
    can_set_payment_in_progress: false,
    can_set_paid: false
  }
}

onMounted(() => {
  loadData()
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

.groups-container {
  display: grid;
  gap: 1.5rem;
}

.group-card {
  background: var(--surface-card);
  border-radius: var(--border-radius);
  overflow: hidden;
}

.group-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 1rem;
  background: var(--surface-ground);
  border-bottom: 1px solid var(--surface-border);
}

.group-info {
  display: flex;
  align-items: center;
  gap: 1rem;
}

.group-name {
  margin: 0;
  font-size: 1.25rem;
  font-weight: 600;
}

.empty-group {
  text-align: center;
  padding: 3rem;
  color: var(--text-color-secondary);
}

.empty-group i {
  font-size: 3rem;
  margin-bottom: 1rem;
  display: block;
  opacity: 0.3;
}

.approval-flow {
  padding: 1rem;
}

.approval-level {
  margin-bottom: 1rem;
}

.level-details {
  display: flex;
  align-items: center;
  gap: 1rem;
  background: var(--surface-ground);
  padding: 1rem;
  border-radius: var(--border-radius);
  border: 1px solid var(--surface-border);
}

.level-number {
  text-align: center;
}

.level-label {
  font-size: 0.875rem;
  color: var(--text-color-secondary);
  margin-top: 0.25rem;
}

.approver-info {
  display: flex;
  align-items: center;
  gap: 1rem;
  flex: 1;
}

.approver-name {
  flex: 1;
}

.approver-email {
  font-size: 0.875rem;
  color: var(--text-color-secondary);
}

.permissions {
  flex: 2;
}

.permission-label {
  font-size: 0.875rem;
  color: var(--text-color-secondary);
  margin-bottom: 0.5rem;
}

.permissions-grid {
  display: flex;
  flex-wrap: wrap;
  gap: 0.25rem;
}

.level-actions {
  display: flex;
  gap: 0.25rem;
}

.flow-arrow {
  text-align: center;
  color: var(--primary-color);
  font-size: 1.5rem;
  margin: 0.5rem 0;
}

.empty-state {
  text-align: center;
  padding: 4rem;
  color: var(--text-color-secondary);
}

.empty-state i {
  font-size: 4rem;
  margin-bottom: 1rem;
  display: block;
  opacity: 0.3;
}

.empty-state h3 {
  margin: 0 0 0.5rem 0;
  color: var(--text-color);
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

.edit-info {
  background: var(--surface-ground);
  padding: 1rem;
  border-radius: var(--border-radius);
}

.info-row {
  display: flex;
  gap: 0.5rem;
  margin-bottom: 0.5rem;
}

.info-row:last-child {
  margin-bottom: 0;
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