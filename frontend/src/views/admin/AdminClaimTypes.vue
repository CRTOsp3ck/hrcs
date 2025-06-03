<template>
  <div class="page-container">
    <div class="page-header">
      <h1 class="page-title">Claim Types</h1>
      <p class="page-subtitle">Configure available claim types and their settings</p>
    </div>

    <!-- Toolbar -->
    <div class="toolbar">
      <span class="p-input-icon-left">
        <IconField>
          <InputIcon class="pi pi-search" />
          <InputText v-model="filters.global.value" placeholder="Search claim types..." />
        </IconField>
      </span>
      <Button label="Add Claim Type" icon="pi pi-plus" @click="showAddDialog = true" />
    </div>

    <!-- Claim Types Table -->
    <DataTable
      :value="claimTypes"
      :loading="loading"
      :filters="filters"
      filterDisplay="menu"
      :globalFilterFields="['name', 'code', 'category']"
      paginator
      :rows="10"
      :rowsPerPageOptions="[10, 20, 50]"
      responsiveLayout="scroll"
      class="claim-types-table"
    >
      <Column field="name" header="Name" sortable>
        <template #body="slotProps">
          <div class="claim-type-name">
            <i :class="getIconClass(slotProps.data.icon)" :style="{ color: slotProps.data.color }"></i>
            <span>{{ slotProps.data.name }}</span>
          </div>
        </template>
      </Column>

      <Column field="code" header="Code" sortable>
        <template #body="slotProps">
          <Tag :value="slotProps.data.code" severity="secondary" />
        </template>
      </Column>

      <Column field="category" header="Category" sortable>
        <template #body="slotProps">
          <Tag :value="slotProps.data.category" :severity="getCategorySeverity(slotProps.data.category)" />
        </template>
      </Column>

      <Column field="maxAmount" header="Max Amount" sortable>
        <template #body="slotProps">
          <span class="amount">${{ formatAmount(slotProps.data.maxAmount) }}</span>
        </template>
      </Column>

      <Column field="requiresReceipt" header="Receipt Required" sortable>
        <template #body="slotProps">
          <i
            :class="slotProps.data.requiresReceipt ? 'pi pi-check-circle' : 'pi pi-times-circle'"
            :style="{ color: slotProps.data.requiresReceipt ? 'var(--green-500)' : 'var(--surface-400)' }"
          ></i>
        </template>
      </Column>

      <!-- <Column field="approvalLevels" header="Approval Levels" sortable>
        <template #body="slotProps">
          <span class="approval-levels">{{ slotProps.data.approvalLevels }}</span>
        </template>
      </Column> -->

      <Column field="status" header="Status" sortable>
        <template #body="slotProps">
          <InputSwitch
            v-model="slotProps.data.active"
            @change="toggleStatus(slotProps.data)"
          />
        </template>
      </Column>

      <Column header="Actions" :exportable="false" style="width: 120px">
        <template #body="slotProps">
          <div class="flex row">
            <Button
              icon="pi pi-pencil"
              severity="secondary"
              text
              rounded
              @click="editClaimType(slotProps.data)"
              v-tooltip="'Edit'"
            />
            <Button
              icon="pi pi-trash"
              severity="danger"
              text
              rounded
              @click="confirmDelete(slotProps.data)"
              v-tooltip="'Delete'"
            />
          </div>
        </template>
      </Column>
    </DataTable>

    <!-- Add/Edit Dialog -->
    <Dialog
      v-model:visible="showAddDialog"
      :header="editingItem ? 'Edit Claim Type' : 'Add New Claim Type'"
      :style="{ width: '600px' }"
      modal
    >
      <div class="dialog-content">
        <div class="form-row">
          <div class="field">
            <label for="name">Name</label>
            <InputText id="name" v-model="form.name" class="w-full" />
          </div>

          <div class="field">
            <label for="code">Code</label>
            <InputText id="code" v-model="form.code" class="w-full" />
          </div>
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
            <label for="category">Category</label>
            <Dropdown
              id="category"
              v-model="form.category"
              :options="categories"
              optionLabel="label"
              optionValue="value"
              placeholder="Select category"
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
              class="w-full"
            />
          </div>
        </div>

        <div class="form-row">
          <div class="field">
            <label for="icon">Icon</label>
            <Dropdown
              id="icon"
              v-model="form.icon"
              :options="availableIcons"
              optionLabel="label"
              optionValue="value"
              placeholder="Select icon"
              class="w-full"
            >
              <template #value="slotProps">
                <div v-if="slotProps.value" class="icon-option">
                  <i :class="slotProps.value"></i>
                  <span>{{ getIconLabel(slotProps.value) }}</span>
                </div>
              </template>
              <template #option="slotProps">
                <div class="icon-option">
                  <i :class="slotProps.option.value"></i>
                  <span>{{ slotProps.option.label }}</span>
                </div>
              </template>
            </Dropdown>
          </div>

          <div class="field">
            <label for="color">Color</label>
            <ColorPicker id="color" v-model="form.color" />
          </div>
        </div>

        <div class="field">
          <label>Requirements</label>
          <div class="requirements-grid">
            <div class="requirement-item">
              <Checkbox v-model="form.requiresReceipt" inputId="requiresReceipt" binary />
              <label for="requiresReceipt">Requires Receipt</label>
            </div>
            <div class="requirement-item">
              <Checkbox v-model="form.requiresApproval" inputId="requiresApproval" binary />
              <label for="requiresApproval">Requires Approval</label>
            </div>
            <div class="requirement-item">
              <Checkbox v-model="form.requiresJustification" inputId="requiresJustification" binary />
              <label for="requiresJustification">Requires Justification</label>
            </div>
          </div>
        </div>

        <div class="field">
          <label for="approvalLevels">Number of Approval Levels</label>
          <InputNumber
            id="approvalLevels"
            v-model="form.approvalLevels"
            :min="0"
            :max="5"
            showButtons
            class="w-full"
          />
        </div>

        <div class="field">
          <label for="validityPeriod">Validity Period (days)</label>
          <InputNumber
            id="validityPeriod"
            v-model="form.validityPeriod"
            :min="0"
            suffix=" days"
            class="w-full"
          />
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
        <p>Are you sure you want to delete claim type <strong>{{ itemToDelete?.name }}</strong>?</p>
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
import { FilterMatchMode } from '@primevue/core/api'
import { useToast } from 'primevue/usetoast'
import { adminApi } from '@/api'

const toast = useToast()

const loading = ref(false)
const claimTypes = ref([])
const showAddDialog = ref(false)
const showDeleteDialog = ref(false)
const editingItem = ref(null)
const itemToDelete = ref(null)

const form = ref({
  name: '',
  code: '',
  description: '',
  category: '',
  maxAmount: 0,
  icon: '',
  color: '#3b82f6',
  requiresReceipt: false,
  requiresApproval: false,
  requiresJustification: false,
  approvalLevels: 1,
  validityPeriod: 30
})

const filters = ref({
  global: { value: null, matchMode: FilterMatchMode.CONTAINS }
})

const categories = ref([
  { label: 'Travel', value: 'travel' },
  { label: 'Medical', value: 'medical' },
  { label: 'Equipment', value: 'equipment' },
  { label: 'Training', value: 'training' },
  { label: 'Entertainment', value: 'entertainment' },
  { label: 'Other', value: 'other' }
])

const availableIcons = ref([
  { label: 'Airplane', value: 'pi pi-send' },
  { label: 'Medical', value: 'pi pi-heart' },
  { label: 'Desktop', value: 'pi pi-desktop' },
  { label: 'Book', value: 'pi pi-book' },
  { label: 'Ticket', value: 'pi pi-ticket' },
  { label: 'Car', value: 'pi pi-car' },
  { label: 'Home', value: 'pi pi-home' },
  { label: 'Shopping', value: 'pi pi-shopping-cart' },
  { label: 'Gift', value: 'pi pi-gift' },
  { label: 'Phone', value: 'pi pi-phone' }
])

const getIconClass = (icon: string) => {
  return icon || 'pi pi-tag'
}

const getIconLabel = (icon: string) => {
  const found = availableIcons.value.find(i => i.value === icon)
  return found ? found.label : 'Icon'
}

const getCategorySeverity = (category: string) => {
  const severities: Record<string, string> = {
    travel: 'info',
    medical: 'danger',
    equipment: 'warning',
    training: 'success',
    entertainment: 'secondary',
    other: 'contrast'
  }
  return severities[category] || 'secondary'
}

const formatAmount = (amount: number) => {
  return new Intl.NumberFormat('en-US', {
    minimumFractionDigits: 0,
    maximumFractionDigits: 0
  }).format(amount)
}

const loadClaimTypes = async () => {
  loading.value = true
  try {
    const response = await adminApi.getClaimTypes()
    claimTypes.value = response.data.data
  } catch (error) {
    toast.add({
      severity: 'error',
      summary: 'Error',
      detail: 'Failed to load claim types',
      life: 3000
    })
  } finally {
    loading.value = false
  }
}

const toggleStatus = async (claimType: any) => {
  try {
    await adminApi.updateClaimType(claimType.id, { active: claimType.active })
    toast.add({
      severity: 'success',
      summary: 'Success',
      detail: `Claim type ${claimType.active ? 'activated' : 'deactivated'}`,
      life: 3000
    })
  } catch (error) {
    claimType.active = !claimType.active
    toast.add({
      severity: 'error',
      summary: 'Error',
      detail: 'Failed to update status',
      life: 3000
    })
  }
}

const editClaimType = (claimType: any) => {
  editingItem.value = claimType
  form.value = {
    name: claimType.name,
    code: claimType.code,
    description: claimType.description,
    category: claimType.category,
    maxAmount: claimType.maxAmount,
    icon: claimType.icon,
    color: claimType.color,
    requiresReceipt: claimType.requiresReceipt,
    requiresApproval: claimType.requiresApproval,
    requiresJustification: claimType.requiresJustification,
    approvalLevels: claimType.approvalLevels,
    validityPeriod: claimType.validityPeriod
  }
  showAddDialog.value = true
}

const confirmDelete = (claimType: any) => {
  itemToDelete.value = claimType
  showDeleteDialog.value = true
}

const deleteItem = async () => {
  try {
    await adminApi.deleteClaimType(itemToDelete.value.id)
    toast.add({
      severity: 'success',
      summary: 'Success',
      detail: 'Claim type deleted successfully',
      life: 3000
    })
    loadClaimTypes()
    showDeleteDialog.value = false
  } catch (error) {
    toast.add({
      severity: 'error',
      summary: 'Error',
      detail: 'Failed to delete claim type',
      life: 3000
    })
  }
}

const save = async () => {
  try {
    if (editingItem.value) {
      await adminApi.updateClaimType(editingItem.value.id, form.value)
      toast.add({
        severity: 'success',
        summary: 'Success',
        detail: 'Claim type updated successfully',
        life: 3000
      })
    } else {
      await adminApi.createClaimType(form.value)
      toast.add({
        severity: 'success',
        summary: 'Success',
        detail: 'Claim type created successfully',
        life: 3000
      })
    }
    loadClaimTypes()
    closeDialog()
  } catch (error) {
    toast.add({
      severity: 'error',
      summary: 'Error',
      detail: 'Failed to save claim type',
      life: 3000
    })
  }
}

const closeDialog = () => {
  showAddDialog.value = false
  editingItem.value = null
  form.value = {
    name: '',
    code: '',
    description: '',
    category: '',
    maxAmount: 0,
    icon: '',
    color: '#3b82f6',
    requiresReceipt: false,
    requiresApproval: false,
    requiresJustification: false,
    approvalLevels: 1,
    validityPeriod: 30
  }
}

onMounted(() => {
  loadClaimTypes()
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

.claim-type-name {
  display: flex;
  align-items: center;
  gap: 0.75rem;
}

.claim-type-name i {
  font-size: 1.25rem;
}

.amount {
  font-weight: 600;
}

.approval-levels {
  font-weight: 500;
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

.requirements-grid {
  display: grid;
  grid-template-columns: 1fr;
  gap: 0.75rem;
}

.requirement-item {
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.icon-option {
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

  .form-row {
    grid-template-columns: 1fr;
  }
}
</style>
