<template>
  <div class="admin-page-container">
    <div class="page-header">
      <h1 class="page-title">Group Management</h1>
      <p class="page-subtitle">Configure user groups and permissions</p>
    </div>

    <!-- Toolbar -->
    <div class="toolbar">
      <span class="p-input-icon-left">
        <IconField>
          <InputIcon class="pi pi-search" />
          <InputText v-model="filters.global.value" placeholder="Search groups..." />
        </IconField>
      </span>
      <Button label="Add Group" icon="pi pi-plus" @click="showAddGroupDialog = true" />
    </div>

    <!-- Groups Grid -->
    <div class="groups-grid">
      <Card v-for="group in filteredGroups" :key="group.id" class="group-card">
        <template #header>
          <div class="group-header">
            <div class="group-info">
              <h3 class="group-name">{{ group.name }}</h3>
              <Tag :value="`${group.members.length} members`" severity="secondary" />
            </div>
            <div class="group-actions">
              <Button
                icon="pi pi-eye"
                severity="info"
                text
                rounded
                @click="viewGroupDetails(group)"
                v-tooltip="'View Details'"
              />
              <Button
                icon="pi pi-pencil"
                severity="secondary"
                text
                rounded
                @click="editGroup(group)"
                v-tooltip="'Edit'"
              />
              <Button
                icon="pi pi-trash"
                severity="danger"
                text
                rounded
                @click="confirmDeleteGroup(group)"
                v-tooltip="'Delete'"
              />
            </div>
          </div>
        </template>
        <template #content>
          <p class="group-description">{{ group.description }}</p>

          <div class="group-details">
            <div class="detail-item">
              <span class="detail-label">Created:</span>
              <span class="detail-value">{{ formatDate(group.createdAt) }}</span>
            </div>
            <div class="detail-item">
              <span class="detail-label">Department:</span>
              <span class="detail-value">{{ group.department || 'All' }}</span>
            </div>
          </div>

          <div class="permissions-section">
            <h4 class="permissions-title">Permissions</h4>
            <div class="permissions-list">
              <Tag
                v-for="permission in group.permissions"
                :key="permission"
                :value="permission"
                severity="info"
              />
            </div>
          </div>
        </template>
      </Card>
    </div>

    <!-- Add/Edit Group Dialog -->
    <Dialog
      v-model:visible="showAddGroupDialog"
      :header="editingGroup ? 'Edit Group' : 'Add New Group'"
      :style="{ width: '600px' }"
      modal
    >
      <div class="dialog-content">
        <div class="field">
          <label for="name">Group Name</label>
          <InputText id="name" v-model="groupForm.name" class="w-full" />
        </div>

        <div class="field">
          <label for="description">Description</label>
          <Textarea
            id="description"
            v-model="groupForm.description"
            rows="3"
            class="w-full"
          />
        </div>

        <div class="field">
          <label for="department">Department</label>
          <Dropdown
            id="department"
            v-model="groupForm.department"
            :options="departments"
            optionLabel="label"
            optionValue="value"
            placeholder="Select department (optional)"
            class="w-full"
            showClear
          />
        </div>

        <div class="field">
          <label>Permissions</label>
          <div class="permissions-grid">
            <div v-for="permission in availablePermissions" :key="permission.value" class="permission-item">
              <Checkbox
                v-model="groupForm.permissions"
                :inputId="permission.value"
                :value="permission.value"
              />
              <label :for="permission.value">{{ permission.label }}</label>
            </div>
          </div>
        </div>

        <div class="field">
          <label for="members">Group Members</label>
          <MultiSelect
            id="members"
            v-model="groupForm.members"
            :options="availableUsers"
            optionLabel="name"
            optionValue="id"
            placeholder="Select members"
            class="w-full"
            filter
          />
        </div>
      </div>

      <template #footer>
        <Button label="Cancel" severity="secondary" @click="closeGroupDialog" />
        <Button label="Save" @click="saveGroup" />
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
        <p>Are you sure you want to delete group <strong>{{ groupToDelete?.name }}</strong>?</p>
        <p class="text-secondary">{{ groupToDelete?.members?.length || 0 }} users will be removed from this group.</p>
      </div>

      <template #footer>
        <Button label="Cancel" severity="secondary" @click="showDeleteDialog = false" />
        <Button label="Delete" severity="danger" @click="deleteGroup" />
      </template>
    </Dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { FilterMatchMode } from '@primevue/core/api'
import { useToast } from 'primevue/usetoast'
import { adminApi } from '@/api'

const toast = useToast()
const router = useRouter()

const loading = ref(false)
const groups = ref([])
const availableUsers = ref([])
const showAddGroupDialog = ref(false)
const showDeleteDialog = ref(false)
const editingGroup = ref(null)
const groupToDelete = ref(null)

const groupForm = ref({
  name: '',
  description: '',
  department: null,
  permissions: [],
  members: []
})

const filters = ref({
  global: { value: null, matchMode: FilterMatchMode.CONTAINS }
})

const departments = ref([
  { label: 'Finance', value: 'finance' },
  { label: 'HR', value: 'hr' },
  { label: 'IT', value: 'it' },
  { label: 'Operations', value: 'operations' },
  { label: 'Sales', value: 'sales' },
  { label: 'Marketing', value: 'marketing' }
])

const availablePermissions = ref([
  { label: 'View Claims', value: 'claims.view' },
  { label: 'Create Claims', value: 'claims.create' },
  { label: 'Edit Claims', value: 'claims.edit' },
  { label: 'Delete Claims', value: 'claims.delete' },
  { label: 'Approve Claims', value: 'claims.approve' },
  { label: 'View Reports', value: 'reports.view' },
  { label: 'Export Reports', value: 'reports.export' },
  { label: 'Manage Users', value: 'users.manage' },
  { label: 'System Settings', value: 'system.settings' }
])

const filteredGroups = computed(() => {
  if (!filters.value.global.value) {
    return groups.value
  }

  const searchTerm = filters.value.global.value.toLowerCase()
  return groups.value.filter((group: any) => {
    return group.name.toLowerCase().includes(searchTerm) ||
           group.description.toLowerCase().includes(searchTerm) ||
           (group.department && group.department.toLowerCase().includes(searchTerm))
  })
})

const formatDate = (date: string) => {
  return new Date(date).toLocaleDateString('en-US', {
    year: 'numeric',
    month: 'short',
    day: 'numeric'
  })
}

const loadGroups = async () => {
  loading.value = true
  try {
    const response = await adminApi.getGroups()
    groups.value = response.data.data
  } catch (error) {
    toast.add({
      severity: 'error',
      summary: 'Error',
      detail: 'Failed to load groups',
      life: 3000
    })
  } finally {
    loading.value = false
  }
}

const loadUsers = async () => {
  try {
    const response = await adminApi.getUsers()
    availableUsers.value = response.data.data
  } catch (error) {
    console.error('Failed to load users:', error)
  }
}

const viewGroupDetails = (group: any) => {
  router.push(`/admin/groups/${group.id}`)
}

const editGroup = (group: any) => {
  editingGroup.value = group
  groupForm.value = {
    name: group.name,
    description: group.description,
    department: group.department,
    permissions: [...group.permissions],
    members: group.members.map((m: any) => m.id)
  }
  showAddGroupDialog.value = true
}

const confirmDeleteGroup = (group: any) => {
  groupToDelete.value = group
  showDeleteDialog.value = true
}

const deleteGroup = async () => {
  try {
    await adminApi.deleteGroup(groupToDelete.value.id)
    toast.add({
      severity: 'success',
      summary: 'Success',
      detail: 'Group deleted successfully',
      life: 3000
    })
    loadGroups()
    showDeleteDialog.value = false
  } catch (error) {
    toast.add({
      severity: 'error',
      summary: 'Error',
      detail: 'Failed to delete group',
      life: 3000
    })
  }
}

const saveGroup = async () => {
  try {
    if (editingGroup.value) {
      await adminApi.updateGroup(editingGroup.value.id, groupForm.value)
      toast.add({
        severity: 'success',
        summary: 'Success',
        detail: 'Group updated successfully',
        life: 3000
      })
    } else {
      await adminApi.createGroup(groupForm.value)
      toast.add({
        severity: 'success',
        summary: 'Success',
        detail: 'Group created successfully',
        life: 3000
      })
    }
    loadGroups()
    closeGroupDialog()
  } catch (error) {
    toast.add({
      severity: 'error',
      summary: 'Error',
      detail: 'Failed to save group',
      life: 3000
    })
  }
}

const closeGroupDialog = () => {
  showAddGroupDialog.value = false
  editingGroup.value = null
  groupForm.value = {
    name: '',
    description: '',
    department: null,
    permissions: [],
    members: []
  }
}

onMounted(() => {
  loadGroups()
  loadUsers()
})
</script>

<style scoped>
.toolbar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 1.5rem;
  gap: 1rem;
}

.groups-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(350px, 1fr));
  gap: 1.5rem;
}

.group-card {
  height: 100%;
}

.group-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.group-info {
  display: flex;
  align-items: center;
  gap: 0.75rem;
}

.group-name {
  margin: 0;
  font-size: 1.125rem;
  font-weight: 600;
}

.group-actions {
  display: flex;
  gap: 0.25rem;
}

.group-description {
  color: var(--surface-600);
  margin: 0 0 1rem;
}

.group-details {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
  margin-bottom: 1.5rem;
}

.detail-item {
  display: flex;
  gap: 0.5rem;
  font-size: 0.875rem;
}

.detail-label {
  color: var(--surface-600);
  font-weight: 500;
}

.detail-value {
  color: var(--surface-700);
}

.permissions-section {
  border-top: 1px solid var(--surface-200);
  padding-top: 1rem;
}

.permissions-title {
  font-size: 0.875rem;
  font-weight: 600;
  margin: 0 0 0.75rem;
  color: var(--surface-700);
}

.permissions-list {
  display: flex;
  flex-wrap: wrap;
  gap: 0.5rem;
}

.dialog-content {
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

.permissions-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 0.75rem;
}

.permission-item {
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
  .toolbar {
    flex-direction: column;
    align-items: stretch;
  }

  .groups-grid {
    grid-template-columns: 1fr;
  }

  .permissions-grid {
    grid-template-columns: 1fr;
  }
}
</style>
