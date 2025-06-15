<template>
  <div class="admin-page-container">
    <PageHeader
      :title="pageTitle"
      :breadcrumbs="breadcrumbItems"
    >
      <template #actions>
        <Button 
          icon="pi pi-pencil" 
          label="Edit Group" 
          @click="editGroup" 
          :loading="loading"
        />
        <Button 
          icon="pi pi-users" 
          label="Manage Members" 
          severity="secondary" 
          @click="manageMembers"
          :loading="loading"
        />
      </template>
    </PageHeader>

    <div v-if="loading" class="loading-container">
      <ProgressSpinner />
    </div>

    <div v-else-if="error" class="error-container">
      {{ error }}
    </div>

    <div v-else-if="groupDetails" class="content-grid">
      <!-- Basic Information -->
      <Card class="col-span-full lg:col-span-1">
        <template #title>Basic Information</template>
        <template #content>
          <div class="grid grid-cols-1 md:grid-cols-2" style="gap: var(--space-13);">
            <div class="info-item">
              <label>Group Name:</label>
              <p>{{ groupDetails.group.name }}</p>
            </div>
            <div class="info-item">
              <label>Description:</label>
              <p>{{ groupDetails.group.description || 'No description' }}</p>
            </div>
            <div class="info-item">
              <label>Total Members:</label>
              <p>{{ groupDetails.members?.length || 0 }}</p>
            </div>
            <div class="info-item">
              <label>Created:</label>
              <p>{{ formatDate(groupDetails.group.created_at) }}</p>
            </div>
          </div>
        </template>
      </Card>

      <!-- Group Members -->
      <Card class="col-span-full lg:col-span-1">
        <template #title>Group Members</template>
        <template #content>
          <DataTable 
            v-if="groupDetails.members.length > 0"
            :value="groupDetails.members" 
            responsiveLayout="scroll"
            class="p-datatable-sm"
          >
            <Column header="Name">
              <template #body="slotProps">
                {{ slotProps.data.first_name }} {{ slotProps.data.last_name }}
              </template>
            </Column>
            <Column field="email" header="Email" />
            <Column field="role" header="Role">
              <template #body="slotProps">
                <Tag 
                  :value="slotProps.data.role" 
                  :severity="slotProps.data.role === 'admin' ? 'danger' : 'info'" 
                />
              </template>
            </Column>
            <Column header="Actions">
              <template #body="slotProps">
                <Button 
                  icon="pi pi-eye" 
                  size="small" 
                  text 
                  @click="viewUserDetails(slotProps.data.id)"
                  v-tooltip.top="'View Details'"
                />
              </template>
            </Column>
          </DataTable>
          <div v-else class="no-data-message">
            No members found
          </div>
        </template>
      </Card>

      <!-- Claim Type Permissions -->
      <Card class="col-span-full">
        <template #title>Claim Type Permissions</template>
        <template #content>
          <div class="section-actions">
            <Button 
              icon="pi pi-plus" 
              label="Manage Permissions" 
              @click="managePermissions"
              size="small"
            />
          </div>
          <DataTable 
            v-if="groupDetails.permissions.length > 0"
            :value="groupDetails.permissions" 
            responsiveLayout="scroll"
            class="p-datatable-sm"
          >
            <Column field="claim_type.name" header="Claim Type" />
            <Column field="is_allowed" header="Allowed">
              <template #body="slotProps">
                <Tag 
                  :value="slotProps.data.is_allowed ? 'Yes' : 'No'" 
                  :severity="slotProps.data.is_allowed ? 'success' : 'danger'" 
                />
              </template>
            </Column>
            <Column field="custom_limit_amount" header="Custom Limit">
              <template #body="slotProps">
                {{ slotProps.data.custom_limit_amount ? 
                   '$' + slotProps.data.custom_limit_amount.toFixed(2) : 
                   'Default' }}
              </template>
            </Column>
            <Column header="Actions">
              <template #body="slotProps">
                <Button 
                  icon="pi pi-pencil" 
                  size="small" 
                  text 
                  @click="editPermission(slotProps.data)"
                  v-tooltip.top="'Edit Permission'"
                />
              </template>
            </Column>
          </DataTable>
          <div v-else class="no-data-message">
            No claim type permissions configured
          </div>
        </template>
      </Card>

      <!-- Approval Levels -->
      <Card class="col-span-full">
        <template #title>Approval Levels</template>
        <template #content>
          <div class="section-actions">
            <Button 
              icon="pi pi-plus" 
              label="Add Approval Level" 
              @click="addApprovalLevel"
              size="small"
            />
          </div>
          <DataTable 
            v-if="groupDetails.approval_levels.length > 0"
            :value="groupDetails.approval_levels" 
            responsiveLayout="scroll"
            class="p-datatable-sm"
            :sortField="'level'"
            :sortOrder="1"
          >
            <Column field="level" header="Level" sortable />
            <Column header="Approver">
              <template #body="slotProps">
                {{ slotProps.data.approver?.first_name }} {{ slotProps.data.approver?.last_name }}
              </template>
            </Column>
            <Column header="Permissions">
              <template #body="slotProps">
                <div class="tag-container">
                  <Tag v-if="slotProps.data.can_approve" value="Approve" severity="success" />
                  <Tag v-if="slotProps.data.can_reject" value="Reject" severity="danger" />
                  <Tag v-if="slotProps.data.can_set_paid" value="Set Paid" severity="info" />
                </div>
              </template>
            </Column>
            <Column header="Actions">
              <template #body="slotProps">
                <Button 
                  icon="pi pi-pencil" 
                  size="small" 
                  text 
                  @click="editApprovalLevel(slotProps.data)"
                  v-tooltip.top="'Edit Level'"
                />
                <Button 
                  icon="pi pi-trash" 
                  size="small" 
                  text 
                  severity="danger"
                  @click="deleteApprovalLevel(slotProps.data.id)"
                  v-tooltip.top="'Delete Level'"
                />
              </template>
            </Column>
          </DataTable>
          <div v-else class="no-data-message">
            No approval levels configured
          </div>
        </template>
      </Card>
    </div>

    <!-- Edit Group Dialog -->
    <Dialog v-model:visible="showEditDialog" header="Edit Group" modal class="w-full max-w-2xl">
      <div class="grid grid-cols-1" style="gap: var(--space-13); padding: var(--space-13) 0;">
        <div>
          <label class="block" style="font-size: var(--text-sm); font-weight: var(--font-medium); color: var(--surface-700); margin-bottom: var(--space-5);">Group Name</label>
          <InputText v-model="editForm.name" class="w-full" />
        </div>
        <div>
          <label class="block" style="font-size: var(--text-sm); font-weight: var(--font-medium); color: var(--surface-700); margin-bottom: var(--space-5);">Description</label>
          <Textarea v-model="editForm.description" class="w-full" rows="3" />
        </div>
      </div>
      <template #footer>
        <Button label="Cancel" severity="secondary" @click="showEditDialog = false" />
        <Button label="Save" @click="saveGroup" :loading="saving" />
      </template>
    </Dialog>

    <!-- Manage Members Dialog -->
    <Dialog v-model:visible="showMembersDialog" header="Manage Members" modal class="w-full max-w-4xl">
      <div class="py-4">
        <p style="color: var(--surface-600); margin-bottom: var(--space-13);">Select users to add or remove from this group:</p>
        <!-- TODO: Implement member management interface -->
        <div class="no-data-message">
          Member management interface will be implemented here
        </div>
      </div>
      <template #footer>
        <Button label="Close" @click="showMembersDialog = false" />
      </template>
    </Dialog>

    <!-- Manage Permissions Dialog -->
    <Dialog v-model:visible="showPermissionsDialog" header="Manage Claim Type Permissions" modal class="w-full max-w-4xl">
      <div class="py-4">
        <p style="color: var(--surface-600); margin-bottom: var(--space-13);">Configure which claim types this group can access:</p>
        <!-- TODO: Implement permissions management interface -->
        <div class="no-data-message">
          Permissions management interface will be implemented here
        </div>
      </div>
      <template #footer>
        <Button label="Close" @click="showPermissionsDialog = false" />
      </template>
    </Dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useToast } from 'primevue/usetoast'
import { adminApi } from '@/api'
import type { UserGroupDetails } from '@/types'
import Button from 'primevue/button'
import Card from 'primevue/card'
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import Tag from 'primevue/tag'
import ProgressSpinner from 'primevue/progressspinner'
import Dialog from 'primevue/dialog'
import InputText from 'primevue/inputtext'
import Textarea from 'primevue/textarea'
import PageHeader from '@/components/base/PageHeader.vue'

const route = useRoute()
const router = useRouter()
const toast = useToast()

const groupDetails = ref<UserGroupDetails | null>(null)
const loading = ref(false)
const error = ref('')
const showEditDialog = ref(false)
const showMembersDialog = ref(false)
const showPermissionsDialog = ref(false)
const saving = ref(false)

const editForm = ref({
  name: '',
  description: ''
})

const pageTitle = computed(() => {
  if (!groupDetails.value?.group) return 'User Group Details'
  return `User Group Details: ${groupDetails.value.group.name}`
})

const breadcrumbItems = computed(() => {
  const items: Array<{ label: string; to?: string }> = [{ label: 'Groups', to: '/admin/groups' }]
  if (groupDetails.value?.group) {
    items.push({ label: groupDetails.value.group.name })
  }
  return items
})

const fetchGroupDetails = async () => {
  try {
    loading.value = true
    error.value = ''
    const groupId = parseInt(route.params.id as string)
    const response = await adminApi.getUserGroupDetails(groupId)
    if (response.data.success) {
      groupDetails.value = response.data.data as UserGroupDetails
    } else {
      error.value = response.data.message || 'Failed to load group details'
    }
  } catch (err: any) {
    error.value = err.message || 'Failed to load group details'
  } finally {
    loading.value = false
  }
}

const editGroup = () => {
  if (groupDetails.value) {
    editForm.value = {
      name: groupDetails.value.group.name,
      description: groupDetails.value.group.description || ''
    }
    showEditDialog.value = true
  }
}

const saveGroup = async () => {
  try {
    saving.value = true
    const groupId = parseInt(route.params.id as string)
    await adminApi.updateGroup(groupId, editForm.value)
    showEditDialog.value = false
    toast.add({
      severity: 'success',
      summary: 'Success',
      detail: 'Group updated successfully'
    })
    await fetchGroupDetails() // Refresh data
  } catch (err: any) {
    toast.add({
      severity: 'error',
      summary: 'Error',
      detail: err.message || 'Failed to update group'
    })
  } finally {
    saving.value = false
  }
}

const manageMembers = () => {
  showMembersDialog.value = true
}

const managePermissions = () => {
  showPermissionsDialog.value = true
}

const viewUserDetails = (userId: number) => {
  router.push(`/admin/users/${userId}`)
}

const editPermission = (permission: any) => {
  // TODO: Implement permission editing
  toast.add({
    severity: 'info',
    summary: 'Info',
    detail: 'Permission editing not yet implemented'
  })
}

const addApprovalLevel = () => {
  // TODO: Implement approval level creation
  toast.add({
    severity: 'info',
    summary: 'Info',
    detail: 'Approval level creation not yet implemented'
  })
}

const editApprovalLevel = (level: any) => {
  // TODO: Implement approval level editing
  toast.add({
    severity: 'info',
    summary: 'Info',
    detail: 'Approval level editing not yet implemented'
  })
}

const deleteApprovalLevel = async (levelId: number) => {
  // TODO: Implement approval level deletion
  toast.add({
    severity: 'info',
    summary: 'Info',
    detail: 'Approval level deletion not yet implemented'
  })
}

const formatDate = (dateString: string) => {
  return new Date(dateString).toLocaleDateString()
}

onMounted(() => {
  fetchGroupDetails()
})
</script>

<style scoped>
.admin-page-container {
  max-width: var(--container-lg);
  margin: 0 auto;
  padding: var(--space-21);
}

.breadcrumb {
  font-size: var(--text-sm);
}

.info-item {
  background-color: var(--surface-50);
  padding: var(--space-8);
  border-radius: var(--space-5);
  margin-bottom: var(--space-8);
  transition: all var(--transition-fast) ease;
}

.info-item:hover {
  background-color: var(--surface-100);
}

.info-item label {
  color: var(--surface-700);
  font-weight: var(--font-semibold);
  font-size: var(--text-sm);
}

.info-item p {
  color: var(--surface-900);
  margin-top: var(--space-3);
}

.error-container {
  background: var(--red-50);
  border: 1px solid var(--red-200);
  color: var(--red-700);
  padding: var(--space-13);
  border-radius: var(--space-8);
  margin-bottom: var(--space-13);
}

.loading-container {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 16rem;
}

.content-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: var(--space-21);
}

.content-grid .col-span-full {
  grid-column: 1 / -1;
}

.no-data-message {
  text-align: center;
  padding: var(--space-34) 0;
  color: var(--surface-500);
}

.section-actions {
  margin-bottom: var(--space-13);
}

.tag-container {
  display: flex;
  flex-wrap: wrap;
  gap: var(--space-3);
}

@media (max-width: var(--breakpoint-lg)) {
  .content-grid {
    grid-template-columns: 1fr;
    gap: var(--space-13);
  }
}

@media (max-width: var(--breakpoint-md)) {
  .admin-page-container {
    padding: var(--space-13);
  }
  
  .content-grid {
    gap: var(--space-8);
  }
}
</style>