<template>
  <div class="admin-page-container">
    <div class="page-header">
      <h1 class="page-title">User Management</h1>
      <p class="page-subtitle">Manage system users and their permissions</p>
    </div>

    <!-- Search and Actions Bar -->
    <div class="toolbar">
      <span class="p-input-icon-left">
        <IconField>
          <InputIcon class="pi pi-search" />
          <InputText v-model="filters.global.value" placeholder="Search users..." />
        </IconField>
      </span>
      <div class="toolbar-actions">
        <Button label="Export" icon="pi pi-download" severity="secondary" @click="exportUsers" />
        <Button label="Add User" icon="pi pi-plus" @click="showAddUserDialog = true" />
      </div>
    </div>

    <!-- Users Table -->
    <DataTable
      :value="users"
      :loading="loading"
      :filters="filters"
      filterDisplay="menu"
      :globalFilterFields="['name', 'email', 'role', 'department']"
      paginator
      :rows="10"
      :rowsPerPageOptions="[10, 20, 50]"
      responsiveLayout="scroll"
      class="users-table"
    >
      <Column field="name" header="Name" sortable>
        <template #body="slotProps">
          <div class="user-info">
            <Avatar :label="getInitials(slotProps.data.name)" shape="circle" />
            <span>{{ slotProps.data.name }}</span>
          </div>
        </template>
      </Column>

      <Column field="email" header="Email" sortable />

      <Column field="role" header="Role" sortable>
        <template #body="slotProps">
          <Tag :value="slotProps.data.role" :severity="getRoleSeverity(slotProps.data.role)" />
        </template>
      </Column>

      <Column field="department" header="Department" sortable />

      <Column field="status" header="Status" sortable>
        <template #body="slotProps">
          <Tag
            :value="slotProps.data.status"
            :severity="slotProps.data.status === 'active' ? 'success' : 'danger'"
          />
        </template>
      </Column>

      <Column field="createdAt" header="Created" sortable>
        <template #body="slotProps">
          {{ formatDate(slotProps.data.createdAt) }}
        </template>
      </Column>

      <Column header="Actions" :exportable="false" style="width: 160px">
        <template #body="slotProps">
          <Button
            icon="pi pi-eye"
            severity="info"
            text
            rounded
            @click="viewUserDetails(slotProps.data)"
            v-tooltip="'View Details'"
          />
          <Button
            icon="pi pi-pencil"
            severity="secondary"
            text
            rounded
            @click="editUser(slotProps.data)"
            v-tooltip="'Edit'"
          />
          <Button
            icon="pi pi-trash"
            severity="danger"
            text
            rounded
            @click="confirmDeleteUser(slotProps.data)"
            v-tooltip="'Delete'"
          />
        </template>
      </Column>
    </DataTable>

    <!-- Add/Edit User Dialog -->
    <Dialog
      v-model:visible="showAddUserDialog"
      :header="editingUser ? 'Edit User' : 'Add New User'"
      :style="{ width: '500px' }"
      modal
    >
      <div class="dialog-content">
        <div class="field">
          <label for="name">Full Name</label>
          <InputText id="name" v-model="userForm.name" class="w-full" />
        </div>

        <div class="field">
          <label for="email">Email</label>
          <InputText id="email" v-model="userForm.email" type="email" class="w-full" />
        </div>

        <div class="field">
          <label for="role">Role</label>
          <Dropdown
            id="role"
            v-model="userForm.role"
            :options="roles"
            optionLabel="label"
            optionValue="value"
            placeholder="Select a role"
            class="w-full"
          />
        </div>

        <div class="field">
          <label for="department">Department</label>
          <InputText id="department" v-model="userForm.department" class="w-full" />
        </div>

        <div class="field">
          <label for="groups">Groups</label>
          <MultiSelect
            id="groups"
            v-model="userForm.groups"
            :options="groups"
            optionLabel="name"
            optionValue="id"
            placeholder="Select groups"
            class="w-full"
          />
        </div>
      </div>

      <template #footer>
        <Button label="Cancel" severity="secondary" @click="closeUserDialog" />
        <Button label="Save" @click="saveUser" />
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
        <p>Are you sure you want to delete user <strong>{{ userToDelete?.name }}</strong>?</p>
        <p class="text-secondary">This action cannot be undone.</p>
      </div>

      <template #footer>
        <Button label="Cancel" severity="secondary" @click="showDeleteDialog = false" />
        <Button label="Delete" severity="danger" @click="deleteUser" />
      </template>
    </Dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { FilterMatchMode } from '@primevue/core/api'
import { useToast } from 'primevue/usetoast'
import { adminApi } from '@/api'

const toast = useToast()
const router = useRouter()

const loading = ref(false)
const users = ref([])
const groups = ref([])
const showAddUserDialog = ref(false)
const showDeleteDialog = ref(false)
const editingUser = ref(null)
const userToDelete = ref(null)

const userForm = ref({
  name: '',
  email: '',
  role: '',
  department: '',
  groups: []
})

const filters = ref({
  global: { value: null, matchMode: FilterMatchMode.CONTAINS }
})

const roles = ref([
  { label: 'Admin', value: 'admin' },
  { label: 'Manager', value: 'manager' },
  { label: 'Employee', value: 'employee' },
  { label: 'Viewer', value: 'viewer' }
])

const getInitials = (name: string) => {
  return name
    .split(' ')
    .map(n => n[0])
    .join('')
    .toUpperCase()
    .slice(0, 2)
}

const getRoleSeverity = (role: string) => {
  const severities: Record<string, string> = {
    admin: 'danger',
    manager: 'warning',
    employee: 'info',
    viewer: 'secondary'
  }
  return severities[role] || 'secondary'
}

const formatDate = (date: string) => {
  return new Date(date).toLocaleDateString('en-US', {
    year: 'numeric',
    month: 'short',
    day: 'numeric'
  })
}

const loadUsers = async () => {
  loading.value = true
  try {
    const response = await adminApi.getUsers()
    users.value = response.data.data
  } catch (error) {
    toast.add({
      severity: 'error',
      summary: 'Error',
      detail: 'Failed to load users',
      life: 3000
    })
  } finally {
    loading.value = false
  }
}

const loadGroups = async () => {
  try {
    const response = await adminApi.getGroups()
    groups.value = response.data.data
  } catch (error) {
    console.error('Failed to load groups:', error)
  }
}

const viewUserDetails = (user: any) => {
  router.push(`/admin/users/${user.id}`)
}

const editUser = (user: any) => {
  editingUser.value = user
  userForm.value = {
    name: user.name,
    email: user.email,
    role: user.role,
    department: user.department,
    groups: user.groups?.map((g: any) => g.id) || []
  }
  showAddUserDialog.value = true
}

const confirmDeleteUser = (user: any) => {
  userToDelete.value = user
  showDeleteDialog.value = true
}

const deleteUser = async () => {
  try {
    await adminApi.deleteUser(userToDelete.value.id)
    toast.add({
      severity: 'success',
      summary: 'Success',
      detail: 'User deleted successfully',
      life: 3000
    })
    loadUsers()
    showDeleteDialog.value = false
  } catch (error) {
    toast.add({
      severity: 'error',
      summary: 'Error',
      detail: 'Failed to delete user',
      life: 3000
    })
  }
}

const saveUser = async () => {
  try {
    if (editingUser.value) {
      await adminApi.updateUser(editingUser.value.id, userForm.value)
      toast.add({
        severity: 'success',
        summary: 'Success',
        detail: 'User updated successfully',
        life: 3000
      })
    } else {
      await adminApi.createUser(userForm.value)
      toast.add({
        severity: 'success',
        summary: 'Success',
        detail: 'User created successfully',
        life: 3000
      })
    }
    loadUsers()
    closeUserDialog()
  } catch (error) {
    toast.add({
      severity: 'error',
      summary: 'Error',
      detail: 'Failed to save user',
      life: 3000
    })
  }
}

const closeUserDialog = () => {
  showAddUserDialog.value = false
  editingUser.value = null
  userForm.value = {
    name: '',
    email: '',
    role: '',
    department: '',
    groups: []
  }
}

const exportUsers = () => {
  // Implement export functionality
  toast.add({
    severity: 'info',
    summary: 'Export',
    detail: 'Export functionality not implemented yet',
    life: 3000
  })
}

onMounted(() => {
  loadUsers()
  loadGroups()
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

.toolbar-actions {
  display: flex;
  gap: 0.5rem;
}

.user-info {
  display: flex;
  align-items: center;
  gap: 0.75rem;
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

  .toolbar-actions {
    justify-content: flex-end;
  }
}
</style>
