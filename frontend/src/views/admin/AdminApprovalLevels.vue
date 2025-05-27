<template>
  <div class="page-container">
    <div class="page-header">
      <h1 class="page-title">Approval Levels</h1>
      <p class="page-subtitle">Configure approval hierarchy and permissions for each user group</p>
    </div>

    <!-- Info Banner -->
    <Message severity="info" :closable="false" class="info-banner">
      <template #default>
        <div class="info-content">
          <i class="pi pi-info-circle"></i>
          <span>Add approvers to each user group. They are automatically assigned levels (1st, 2nd, 3rd approver, etc.). You can drag to reorder the approval hierarchy.</span>
        </div>
      </template>
    </Message>

    <!-- Loading State -->
    <div v-if="loading" class="loading-container">
      <ProgressSpinner />
      <p>Loading approval levels...</p>
    </div>

    <!-- Debug Info -->
    <div v-else-if="!loading && groupsWithLevels.length === 0" class="debug-info">
      <Message severity="warn">
        <p>No approval level groups found. This could mean:</p>
        <ul>
          <li>No user groups exist in the system</li>
          <li>The API endpoint is not returning data correctly</li>
          <li>Check the browser console for errors</li>
        </ul>
      </Message>
      <div class="debug-data">
        <h4>Debug Data:</h4>
        <pre>{{ JSON.stringify(groupsWithLevels, null, 2) }}</pre>
      </div>
    </div>

    <!-- User Groups with Approval Levels -->
    <div v-else class="groups-container">
      <Panel v-for="group in groupsWithLevels" :key="group.groupId" class="group-panel">
        <template #header>
          <div class="group-header">
            <div class="group-info">
              <h3>{{ group.groupName }}</h3>
              <Tag :value="`${group.levels.length} Approver${group.levels.length !== 1 ? 's' : ''}`" severity="info" />
            </div>
            <Button 
              label="Add Approver" 
              icon="pi pi-plus" 
              size="small"
              @click="openAddApprover(group.groupId, group.groupName)"
            />
          </div>
        </template>

        <!-- Empty State -->
        <div v-if="group.levels.length === 0" class="empty-group">
          <i class="pi pi-users"></i>
          <p>No approvers configured for this group</p>
          <Button 
            label="Add First Approver" 
            icon="pi pi-plus" 
            severity="secondary"
            @click="openAddApprover(group.groupId, group.groupName)"
          />
        </div>

        <!-- Draggable Approval Levels -->
        <draggable 
          v-else
          v-model="group.levels"
          group="approvers"
          item-key="id"
          handle=".drag-handle"
          animation="200"
          @end="onDragEnd(group)"
          class="levels-list"
        >
          <template #item="{ element: level, index }">
            <div class="level-card">
              <div class="level-header">
                <div class="level-info">
                  <Button 
                    icon="pi pi-bars" 
                    class="drag-handle"
                    severity="secondary"
                    text
                    v-tooltip.top="'Drag to reorder'"
                  />
                  <Badge :value="`Level ${index + 1}`" :severity="getLevelSeverity(index)" />
                  <div class="approver-details">
                    <Avatar 
                      :label="getInitials(level.approver)" 
                      size="small"
                      :style="{ backgroundColor: getAvatarColor(level.approver.id) }"
                    />
                    <div>
                      <div class="approver-name">{{ level.approver.name }}</div>
                      <div class="approver-email">{{ level.approver.email }}</div>
                    </div>
                  </div>
                </div>
                <div class="level-actions">
                  <Button 
                    icon="pi pi-pencil" 
                    severity="secondary" 
                    text 
                    size="small"
                    @click="editLevel(level)"
                    v-tooltip.top="'Edit permissions'"
                  />
                  <Button 
                    icon="pi pi-trash" 
                    severity="danger" 
                    text 
                    size="small"
                    @click="confirmDelete(level, group)"
                    v-tooltip.top="'Remove approver'"
                  />
                </div>
              </div>
              <div class="permissions-display">
                <Tag v-if="level.permissions.canDraft" value="Draft" severity="secondary" />
                <Tag v-if="level.permissions.canSubmit" value="Submit" severity="info" />
                <Tag v-if="level.permissions.canApprove" value="Approve" severity="success" />
                <Tag v-if="level.permissions.canReject" value="Reject" severity="danger" />
                <Tag v-if="level.permissions.canSetPaymentInProgress" value="Payment in Progress" severity="warning" />
                <Tag v-if="level.permissions.canSetPaid" value="Paid" severity="success" />
              </div>
            </div>
          </template>
        </draggable>
      </Panel>
    </div>

    <!-- Add Approver Dialog -->
    <Dialog 
      v-model:visible="showAddDialog" 
      header="Add Approver"
      :style="{ width: '500px' }"
      modal
    >
      <div class="form-container">
        <div class="field">
          <label>User Group</label>
          <div class="group-display">
            <Tag :value="selectedGroupName" severity="info" />
          </div>
        </div>

        <div class="field">
          <label for="approver">Select Approver *</label>
          <Dropdown
            id="approver"
            v-model="addForm.approverId"
            :options="availableApprovers"
            optionLabel="name"
            optionValue="id"
            placeholder="Choose an approver"
            class="w-full"
            filter
            :loading="loadingUsers"
          >
            <template #option="{ option }">
              <div class="approver-option">
                <Avatar 
                  :label="getInitials(option)" 
                  size="small" 
                  :style="{ backgroundColor: getAvatarColor(option.id), marginRight: '0.5rem' }"
                />
                <div>
                  <div>{{ option.name }}</div>
                  <div class="option-subtitle">{{ option.email }}</div>
                </div>
              </div>
            </template>
          </Dropdown>
        </div>

        <Divider />

        <div class="field">
          <label>Permissions</label>
          <div class="permissions-info">
            <i class="pi pi-info-circle"></i>
            <span>This approver will be assigned as Level {{ nextLevel }} approver</span>
          </div>
          <div class="permissions-checkboxes">
            <div class="checkbox-item">
              <Checkbox 
                v-model="addForm.permissions.canDraft" 
                inputId="addCanDraft" 
                binary 
              />
              <label for="addCanDraft">Can set to Draft</label>
            </div>
            <div class="checkbox-item">
              <Checkbox 
                v-model="addForm.permissions.canSubmit" 
                inputId="addCanSubmit" 
                binary 
              />
              <label for="addCanSubmit">Can Submit claims</label>
            </div>
            <div class="checkbox-item">
              <Checkbox 
                v-model="addForm.permissions.canApprove" 
                inputId="addCanApprove" 
                binary 
              />
              <label for="addCanApprove">Can Approve claims</label>
            </div>
            <div class="checkbox-item">
              <Checkbox 
                v-model="addForm.permissions.canReject" 
                inputId="addCanReject" 
                binary 
              />
              <label for="addCanReject">Can Reject claims</label>
            </div>
            <div class="checkbox-item">
              <Checkbox 
                v-model="addForm.permissions.canSetPaymentInProgress" 
                inputId="addCanSetPaymentInProgress" 
                binary 
              />
              <label for="addCanSetPaymentInProgress">Can set Payment in Progress</label>
            </div>
            <div class="checkbox-item">
              <Checkbox 
                v-model="addForm.permissions.canSetPaid" 
                inputId="addCanSetPaid" 
                binary 
              />
              <label for="addCanSetPaid">Can mark as Paid</label>
            </div>
          </div>
        </div>
      </div>
      
      <template #footer>
        <Button label="Cancel" severity="secondary" @click="closeAddDialog" />
        <Button label="Add Approver" @click="addApprover" :loading="saving" />
      </template>
    </Dialog>

    <!-- Edit Permissions Dialog -->
    <Dialog 
      v-model:visible="showEditDialog" 
      header="Edit Permissions"
      :style="{ width: '500px' }"
      modal
    >
      <div class="form-container">
        <div class="field">
          <label>Approver</label>
          <div class="approver-display">
            <Avatar 
              :label="editingLevel ? getInitials(editingLevel.approver) : ''" 
              :style="{ backgroundColor: editingLevel ? getAvatarColor(editingLevel.approver.id) : '' }"
            />
            <div>
              <div class="approver-name">{{ editingLevel?.approver.name }}</div>
              <div class="approver-email">{{ editingLevel?.approver.email }}</div>
            </div>
          </div>
        </div>

        <Divider />

        <div class="field">
          <label>Status Permissions</label>
          <div class="permissions-checkboxes">
            <div class="checkbox-item">
              <Checkbox 
                v-model="editForm.permissions.canDraft" 
                inputId="editCanDraft" 
                binary 
              />
              <label for="editCanDraft">Can set to Draft</label>
            </div>
            <div class="checkbox-item">
              <Checkbox 
                v-model="editForm.permissions.canSubmit" 
                inputId="editCanSubmit" 
                binary 
              />
              <label for="editCanSubmit">Can Submit claims</label>
            </div>
            <div class="checkbox-item">
              <Checkbox 
                v-model="editForm.permissions.canApprove" 
                inputId="editCanApprove" 
                binary 
              />
              <label for="editCanApprove">Can Approve claims</label>
            </div>
            <div class="checkbox-item">
              <Checkbox 
                v-model="editForm.permissions.canReject" 
                inputId="editCanReject" 
                binary 
              />
              <label for="editCanReject">Can Reject claims</label>
            </div>
            <div class="checkbox-item">
              <Checkbox 
                v-model="editForm.permissions.canSetPaymentInProgress" 
                inputId="editCanSetPaymentInProgress" 
                binary 
              />
              <label for="editCanSetPaymentInProgress">Can set Payment in Progress</label>
            </div>
            <div class="checkbox-item">
              <Checkbox 
                v-model="editForm.permissions.canSetPaid" 
                inputId="editCanSetPaid" 
                binary 
              />
              <label for="editCanSetPaid">Can mark as Paid</label>
            </div>
          </div>
        </div>
      </div>
      
      <template #footer>
        <Button label="Cancel" severity="secondary" @click="showEditDialog = false" />
        <Button label="Save Changes" @click="updatePermissions" :loading="saving" />
      </template>
    </Dialog>

    <!-- Delete Confirmation -->
    <Dialog 
      v-model:visible="showDeleteDialog" 
      header="Confirm Remove"
      :style="{ width: '450px' }"
      modal
    >
      <div class="confirmation-content">
        <i class="pi pi-exclamation-triangle" style="font-size: 2rem; color: var(--red-500)"></i>
        <p>Are you sure you want to remove <strong>{{ levelToDelete?.approver.name }}</strong> as an approver?</p>
        <p class="text-secondary">This will automatically adjust the levels of remaining approvers in this group.</p>
      </div>
      
      <template #footer>
        <Button label="Cancel" severity="secondary" @click="showDeleteDialog = false" />
        <Button label="Remove Approver" severity="danger" @click="deleteLevel" :loading="deleting" />
      </template>
    </Dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useToast } from 'primevue/usetoast'
import draggable from 'vuedraggable'
import api from '@/api'

const toast = useToast()

// State
const loading = ref(false)
const loadingUsers = ref(false)
const saving = ref(false)
const deleting = ref(false)
const groupsWithLevels = ref<any[]>([])
const availableUsers = ref<any[]>([])
const showAddDialog = ref(false)
const showEditDialog = ref(false)
const showDeleteDialog = ref(false)
const selectedGroupId = ref<number | null>(null)
const selectedGroupName = ref('')
const editingLevel = ref<any>(null)
const levelToDelete = ref<any>(null)
const groupToDelete = ref<any>(null)

// Forms
const addForm = ref({
  approverId: null as number | null,
  permissions: {
    canDraft: false,
    canSubmit: false,
    canApprove: true,
    canReject: true,
    canSetPaymentInProgress: false,
    canSetPaid: false
  }
})

const editForm = ref({
  permissions: {
    canDraft: false,
    canSubmit: false,
    canApprove: true,
    canReject: true,
    canSetPaymentInProgress: false,
    canSetPaid: false
  }
})

// Computed
const nextLevel = computed(() => {
  if (!selectedGroupId.value) return 1
  const group = groupsWithLevels.value.find(g => g.groupId === selectedGroupId.value)
  return group ? group.levels.length + 1 : 1
})

const availableApprovers = computed(() => {
  if (!selectedGroupId.value) return []
  
  const group = groupsWithLevels.value.find(g => g.groupId === selectedGroupId.value)
  if (!group) return availableUsers.value

  // Filter out users who are already approvers for this group
  const existingApproverIds = group.levels.map((l: any) => l.approverId)
  return availableUsers.value.filter(user => !existingApproverIds.includes(user.id))
})

// Methods
const getInitials = (user: any) => {
  if (!user) return '?'
  const names = user.name ? user.name.split(' ') : [user.firstName || '', user.lastName || '']
  return names.map((n: string) => n[0] || '').join('').toUpperCase() || '?'
}

const getAvatarColor = (userId: number) => {
  const colors = ['#2196F3', '#4CAF50', '#FF9800', '#9C27B0', '#F44336', '#00BCD4', '#795548']
  return colors[userId % colors.length]
}

const getLevelSeverity = (index: number) => {
  const severities = ['success', 'info', 'warning']
  return severities[index % severities.length]
}

const loadApprovalLevels = async () => {
  loading.value = true
  try {
    console.log('Loading approval levels...')
    const response = await api.get('/admin/approval-levels/by-group')
    console.log('Response:', response.data)
    groupsWithLevels.value = response.data.data || []
    console.log('Groups with levels:', groupsWithLevels.value)
  } catch (error: any) {
    console.error('Failed to load approval levels:', error)
    console.error('Error response:', error.response?.data)
    toast.add({
      severity: 'error',
      summary: 'Error',
      detail: error.response?.data?.message || 'Failed to load approval levels',
      life: 3000
    })
  } finally {
    loading.value = false
  }
}

const loadUsers = async () => {
  loadingUsers.value = true
  try {
    const response = await api.get('/admin/users')
    availableUsers.value = (response.data.data || []).map((user: any) => ({
      id: user.id,
      name: user.name || `${user.first_name} ${user.last_name}`,
      email: user.email,
      firstName: user.first_name,
      lastName: user.last_name
    }))
  } catch (error) {
    console.error('Failed to load users:', error)
  } finally {
    loadingUsers.value = false
  }
}

const openAddApprover = (groupId: number, groupName: string) => {
  selectedGroupId.value = groupId
  selectedGroupName.value = groupName
  resetAddForm()
  showAddDialog.value = true
}

const closeAddDialog = () => {
  showAddDialog.value = false
  selectedGroupId.value = null
  selectedGroupName.value = ''
  resetAddForm()
}

const resetAddForm = () => {
  addForm.value = {
    approverId: null,
    permissions: {
      canDraft: false,
      canSubmit: false,
      canApprove: true,
      canReject: true,
      canSetPaymentInProgress: false,
      canSetPaid: false
    }
  }
}

const addApprover = async () => {
  if (!addForm.value.approverId || !selectedGroupId.value) {
    toast.add({
      severity: 'warn',
      summary: 'Validation Error',
      detail: 'Please select an approver',
      life: 3000
    })
    return
  }

  saving.value = true
  try {
    await api.post('/admin/approval-levels', {
      userGroupId: selectedGroupId.value,
      approverId: addForm.value.approverId,
      canDraft: addForm.value.permissions.canDraft,
      canSubmit: addForm.value.permissions.canSubmit,
      canApprove: addForm.value.permissions.canApprove,
      canReject: addForm.value.permissions.canReject,
      canSetPaymentInProgress: addForm.value.permissions.canSetPaymentInProgress,
      canSetPaid: addForm.value.permissions.canSetPaid
    })

    toast.add({
      severity: 'success',
      summary: 'Success',
      detail: 'Approver added successfully',
      life: 3000
    })

    closeAddDialog()
    await loadApprovalLevels()
  } catch (error: any) {
    toast.add({
      severity: 'error',
      summary: 'Error',
      detail: error.response?.data?.message || 'Failed to add approver',
      life: 3000
    })
  } finally {
    saving.value = false
  }
}

const editLevel = (level: any) => {
  editingLevel.value = level
  editForm.value.permissions = { ...level.permissions }
  showEditDialog.value = true
}

const updatePermissions = async () => {
  if (!editingLevel.value) return

  saving.value = true
  try {
    await api.put(`/admin/approval-levels/${editingLevel.value.id}`, {
      canDraft: editForm.value.permissions.canDraft,
      canSubmit: editForm.value.permissions.canSubmit,
      canApprove: editForm.value.permissions.canApprove,
      canReject: editForm.value.permissions.canReject,
      canSetPaymentInProgress: editForm.value.permissions.canSetPaymentInProgress,
      canSetPaid: editForm.value.permissions.canSetPaid
    })

    toast.add({
      severity: 'success',
      summary: 'Success',
      detail: 'Permissions updated successfully',
      life: 3000
    })

    showEditDialog.value = false
    await loadApprovalLevels()
  } catch (error) {
    toast.add({
      severity: 'error',
      summary: 'Error',
      detail: 'Failed to update permissions',
      life: 3000
    })
  } finally {
    saving.value = false
  }
}

const confirmDelete = (level: any, group: any) => {
  levelToDelete.value = level
  groupToDelete.value = group
  showDeleteDialog.value = true
}

const deleteLevel = async () => {
  if (!levelToDelete.value) return

  deleting.value = true
  try {
    await api.delete(`/admin/approval-levels/${levelToDelete.value.id}`)

    toast.add({
      severity: 'success',
      summary: 'Success',
      detail: 'Approver removed successfully',
      life: 3000
    })

    showDeleteDialog.value = false
    await loadApprovalLevels()
  } catch (error) {
    toast.add({
      severity: 'error',
      summary: 'Error',
      detail: 'Failed to remove approver',
      life: 3000
    })
  } finally {
    deleting.value = false
    levelToDelete.value = null
    groupToDelete.value = null
  }
}

const onDragEnd = async (group: any) => {
  // Update the level numbers based on new order
  const orders = group.levels.map((level: any, index: number) => ({
    id: level.id,
    level: index + 1
  }))

  try {
    await api.put('/admin/approval-levels/order', {
      userGroupId: group.groupId,
      orders
    })

    toast.add({
      severity: 'success',
      summary: 'Success',
      detail: 'Approval order updated',
      life: 3000
    })

    // Reload to ensure consistency
    await loadApprovalLevels()
  } catch (error) {
    toast.add({
      severity: 'error',
      summary: 'Error',
      detail: 'Failed to update approval order',
      life: 3000
    })
    // Reload to revert changes
    await loadApprovalLevels()
  }
}

// Lifecycle
onMounted(() => {
  console.log('AdminApprovalLevels component mounted')
  loadApprovalLevels()
  loadUsers()
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

.loading-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 4rem;
  color: var(--text-color-secondary);
}

.groups-container {
  display: grid;
  gap: 1.5rem;
}

.group-panel {
  background: var(--surface-card);
}

.group-panel :deep(.p-panel-header) {
  background: var(--surface-b);
  border: none;
}

.group-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  width: 100%;
}

.group-info {
  display: flex;
  align-items: center;
  gap: 1rem;
}

.group-info h3 {
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

.empty-group p {
  margin: 0 0 1rem 0;
}

.levels-list {
  display: flex;
  flex-direction: column;
  gap: 0.75rem;
}

.level-card {
  background: var(--surface-a);
  border: 1px solid var(--surface-border);
  border-radius: var(--border-radius);
  padding: 1rem;
  transition: all 0.2s;
}

.level-card:hover {
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.level-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 0.75rem;
}

.level-info {
  display: flex;
  align-items: center;
  gap: 1rem;
}

.drag-handle {
  cursor: move;
}

.drag-handle:hover {
  background-color: var(--surface-100);
}

.approver-details {
  display: flex;
  align-items: center;
  gap: 0.75rem;
}

.approver-name {
  font-weight: 500;
  color: var(--text-color);
}

.approver-email {
  font-size: 0.875rem;
  color: var(--text-color-secondary);
}

.level-actions {
  display: flex;
  gap: 0.25rem;
}

.permissions-display {
  display: flex;
  flex-wrap: wrap;
  gap: 0.25rem;
  margin-left: 3rem;
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

.group-display {
  padding: 0.5rem;
  background: var(--surface-a);
  border-radius: var(--border-radius);
}

.permissions-info {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.75rem;
  background: var(--blue-50);
  color: var(--blue-700);
  border-radius: var(--border-radius);
  margin-bottom: 0.5rem;
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
  padding: 0.5rem 0;
}

.option-subtitle {
  font-size: 0.75rem;
  color: var(--text-color-secondary);
}

.approver-display {
  display: flex;
  align-items: center;
  gap: 1rem;
  padding: 1rem;
  background: var(--surface-a);
  border-radius: var(--border-radius);
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

/* Draggable styles */
.sortable-ghost {
  opacity: 0.5;
}

.sortable-drag {
  background: var(--surface-b);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
}

/* Debug styles */
.debug-info {
  margin-top: 2rem;
}

.debug-data {
  margin-top: 1rem;
  padding: 1rem;
  background: var(--surface-a);
  border-radius: var(--border-radius);
  border: 1px solid var(--surface-border);
}

.debug-data h4 {
  margin: 0 0 0.5rem 0;
  color: var(--text-color-secondary);
}

.debug-data pre {
  margin: 0;
  font-size: 0.875rem;
  color: var(--text-color);
  white-space: pre-wrap;
}
</style>