<template>
  <div class="page-container">
    <div class="page-header">
      <div>
        <h1 class="page-title">My Claims</h1>
        <p class="page-subtitle">Manage and track all your expense claims</p>
      </div>
      <Button
        label="New Claim"
        icon="pi pi-plus"
        @click="router.push('/claims/new')"
      />
    </div>

    <!-- Filters -->
    <Card class="filters-card">
      <template #content>
        <div class="filters-grid">
          <div class="filter-field">
            <label class="filter-label">Search</label>
            <span class="w-full">
              <!-- <i class="pi pi-search" /> -->
              <InputText
                v-model="filters.search"
                placeholder="Search by title or description"
                class="w-full"
              />
            </span>
          </div>

          <div class="filter-field">
            <label class="filter-label">Status</label>
            <Dropdown
              v-model="filters.status"
              :options="statusOptions"
              optionLabel="label"
              optionValue="value"
              placeholder="All Statuses"
              class="w-full"
              showClear
            />
          </div>

          <div class="filter-field">
            <label class="filter-label">Type</label>
            <Dropdown
              v-model="filters.type"
              :options="claimTypes"
              optionLabel="name"
              optionValue="id"
              placeholder="All Types"
              class="w-full"
              showClear
            />
          </div>

          <div class="filter-field">
            <label class="filter-label">Date Range</label>
            <div class="date-range">
              <InputText
                v-model="filters.dateFrom"
                type="date"
                class="w-full"
              />
              <span class="date-separator">to</span>
              <InputText
                v-model="filters.dateTo"
                type="date"
                class="w-full"
              />
            </div>
          </div>
        </div>
      </template>
    </Card>

    <!-- Claims Table -->
    <Card>
      <template #content>
        <DataTable
          :value="filteredClaims"
          :loading="loading"
          paginator
          :rows="10"
          :rowsPerPageOptions="[10, 25, 50]"
          dataKey="id"
          responsiveLayout="scroll"
          :globalFilterFields="['title', 'description', 'claim_type.name']"
        >
          <template #header>
            <div class="table-header">
              <span class="table-title">
                <i class="pi pi-file mr-2"></i>
                Total Claims: {{ filteredClaims.length }}
              </span>
              <span class="table-summary">
                Total Amount: <strong>${{ totalAmount }}</strong>
              </span>
            </div>
          </template>

          <template #empty>
            <div class="empty-state">
              <i class="pi pi-inbox empty-icon"></i>
              <h3>No claims found</h3>
              <p>Start by creating your first claim</p>
              <Button
                label="Create Claim"
                icon="pi pi-plus"
                @click="router.push('/claims/new')"
              />
            </div>
          </template>

          <Column selectionMode="multiple" style="width: 3rem" :exportable="false"></Column>

          <Column field="id" header="ID" :sortable="true" style="width: 80px">
            <template #body="slotProps">
              <span class="font-mono">#{{ slotProps.data.id }}</span>
            </template>
          </Column>

          <Column field="title" header="Title" :sortable="true">
            <template #body="slotProps">
              <router-link :to="`/claims/${slotProps.data.id}`" class="claim-link">
                {{ slotProps.data.title }}
              </router-link>
            </template>
          </Column>

          <Column field="claim_type.name" header="Type" :sortable="true" style="width: 150px">
            <template #body="slotProps">
              <Tag :value="slotProps.data.claim_type?.name" severity="info" />
            </template>
          </Column>

          <Column field="amount" header="Amount" :sortable="true" style="width: 120px">
            <template #body="slotProps">
              <span class="amount">${{ formatAmount(slotProps.data.amount) }}</span>
            </template>
          </Column>

          <Column field="status" header="Status" :sortable="true" style="width: 180px">
            <template #body="slotProps">
              <span :class="`status-badge status-${slotProps.data.status}`">
                <i :class="getStatusIcon(slotProps.data.status)"></i>
                {{ formatStatus(slotProps.data.status) }}
              </span>
            </template>
          </Column>

          <Column field="created_at" header="Submitted" :sortable="true" style="width: 120px">
            <template #body="slotProps">
              {{ formatDate(slotProps.data.submitted_at || slotProps.data.created_at) }}
            </template>
          </Column>

          <Column header="Actions" :exportable="false" style="width: 120px">
            <template #body="slotProps">
              <div class="action-buttons">
                <Button
                  icon="pi pi-eye"
                  text
                  rounded
                  v-tooltip="'View'"
                  @click="router.push(`/claims/${slotProps.data.id}`)"
                />
                <Button
                  icon="pi pi-pencil"
                  text
                  rounded
                  v-tooltip="'Edit'"
                  :disabled="!canEdit(slotProps.data)"
                  @click="router.push(`/claims/${slotProps.data.id}/edit`)"
                />
                <Button
                  icon="pi pi-trash"
                  text
                  rounded
                  severity="danger"
                  v-tooltip="'Delete'"
                  :disabled="!canDelete(slotProps.data)"
                  @click="confirmDelete(slotProps.data)"
                />
              </div>
            </template>
          </Column>
        </DataTable>
      </template>
    </Card>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { claimsApi, claimTypesApi } from '@/api'
import { useToast } from 'primevue/usetoast'
import { useConfirm } from 'primevue/useconfirm'
import type { Claim, ClaimType } from '@/types'

const router = useRouter()
const toast = useToast()
const confirm = useConfirm()

const loading = ref(false)
const claims = ref<Claim[]>([])
const claimTypes = ref<ClaimType[]>([])

const filters = ref({
  search: '',
  status: null,
  type: null,
  dateFrom: '',
  dateTo: ''
})

const statusOptions = [
  { label: 'Draft', value: 'draft' },
  { label: 'Submitted', value: 'submitted' },
  { label: 'Approved', value: 'approved' },
  { label: 'Rejected', value: 'rejected' },
  { label: 'Payment in Progress', value: 'payment-in-progress' },
  { label: 'Paid', value: 'paid' }
]

const filteredClaims = computed(() => {
  return claims.value.filter(claim => {
    // Search filter
    if (filters.value.search) {
      const search = filters.value.search.toLowerCase()
      if (!claim.title.toLowerCase().includes(search) &&
          !claim.description.toLowerCase().includes(search)) {
        return false
      }
    }

    // Status filter
    if (filters.value.status && claim.status !== filters.value.status) {
      return false
    }

    // Type filter
    if (filters.value.type && claim.claim_type_id !== filters.value.type) {
      return false
    }

    // Date range filter
    if (filters.value.dateFrom || filters.value.dateTo) {
      const claimDate = new Date(claim.created_at)
      if (filters.value.dateFrom && claimDate < new Date(filters.value.dateFrom)) {
        return false
      }
      if (filters.value.dateTo && claimDate > new Date(filters.value.dateTo)) {
        return false
      }
    }

    return true
  })
})

const totalAmount = computed(() => {
  return formatAmount(
    filteredClaims.value.reduce((sum, claim) => sum + claim.amount, 0)
  )
})

const formatAmount = (amount: number) => {
  return new Intl.NumberFormat('en-US', {
    minimumFractionDigits: 2,
    maximumFractionDigits: 2
  }).format(amount)
}

const formatDate = (date: string) => {
  if (!date) return '-'
  return new Date(date).toLocaleDateString('en-US', {
    year: 'numeric',
    month: 'short',
    day: 'numeric'
  })
}

const formatStatus = (status: string) => {
  return status.split('-').map(word =>
    word.charAt(0).toUpperCase() + word.slice(1)
  ).join(' ')
}

const getStatusIcon = (status: string) => {
  const icons: Record<string, string> = {
    'draft': 'pi pi-pencil',
    'submitted': 'pi pi-send',
    'approved': 'pi pi-check-circle',
    'rejected': 'pi pi-times-circle',
    'payment-in-progress': 'pi pi-spin pi-spinner',
    'paid': 'pi pi-dollar'
  }
  return icons[status] || 'pi pi-circle'
}

const canEdit = (claim: Claim) => {
  return claim.status === 'draft'
}

const canDelete = (claim: Claim) => {
  return claim.status === 'draft' || claim.status === 'submitted'
}

const loadClaims = async () => {
  loading.value = true
  try {
    const [claimsResponse, typesResponse] = await Promise.all([
      claimsApi.getAll(),
      claimTypesApi.getAll()
    ])

    if (claimsResponse.data.data) {
      claims.value = claimsResponse.data.data
    }

    if (typesResponse.data.data) {
      claimTypes.value = typesResponse.data.data
    }
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

const confirmDelete = (claim: Claim) => {
  confirm.require({
    message: `Are you sure you want to delete claim "${claim.title}"?`,
    header: 'Delete Confirmation',
    icon: 'pi pi-exclamation-triangle',
    acceptClass: 'p-button-danger',
    accept: () => deleteClaim(claim.id)
  })
}

const deleteClaim = async (id: number) => {
  try {
    await claimsApi.delete(id)
    toast.add({
      severity: 'success',
      summary: 'Success',
      detail: 'Claim deleted successfully',
      life: 3000
    })
    loadClaims()
  } catch (error) {
    toast.add({
      severity: 'error',
      summary: 'Error',
      detail: 'Failed to delete claim',
      life: 3000
    })
  }
}

onMounted(() => {
  loadClaims()
})
</script>

<style scoped>
.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: var(--space-8);
  animation: fadeIn var(--transition-slow) ease-out;
}

.filters-card {
  margin-bottom: var(--space-6);
  animation: slideIn var(--transition-slow) ease-out;
  animation-delay: 100ms;
  animation-fill-mode: both;
}

.filters-grid {
  display: grid;
  grid-template-columns: 2fr 1fr 1fr 2fr;
  gap: var(--space-6);
  align-items: end;
}

.filter-field {
  display: flex;
  flex-direction: column;
  gap: var(--space-2);
}

.filter-label {
  font-weight: var(--font-medium);
  color: var(--surface-600);
  font-size: var(--text-sm);
  text-transform: uppercase;
  letter-spacing: 0.025em;
}

.date-range {
  display: grid;
  grid-template-columns: 1fr auto 1fr;
  gap: var(--space-2);
  align-items: center;
}

.date-separator {
  color: var(--surface-400);
  font-size: var(--text-sm);
  font-weight: var(--font-medium);
  padding: 0 var(--space-1);
}

.table-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding-bottom: var(--space-4);
  border-bottom: 1px solid var(--surface-100);
}

.table-title {
  font-weight: var(--font-semibold);
  display: flex;
  align-items: center;
  color: var(--surface-700);
  font-size: var(--text-base);
}

.table-title i {
  color: var(--primary-500);
  margin-right: var(--space-2);
}

.table-summary {
  color: var(--surface-600);
  font-size: var(--text-sm);
}

.table-summary strong {
  color: var(--primary-600);
  font-size: var(--text-lg);
  font-weight: var(--font-bold);
}

.claim-link {
  color: var(--primary-600);
  text-decoration: none;
  font-weight: var(--font-medium);
  transition: all var(--transition-fast) ease;
  display: inline-block;
}

.claim-link:hover {
  color: var(--primary-700);
  transform: translateX(2px);
}

.amount {
  font-weight: var(--font-semibold);
  color: var(--surface-900);
  font-size: var(--text-base);
}

.action-buttons {
  display: flex;
  gap: var(--space-1);
  opacity: 0.6;
  transition: opacity var(--transition-fast) ease;
}

tr:hover .action-buttons {
  opacity: 1;
}

.empty-state {
  text-align: center;
  padding: var(--space-16) var(--space-8);
}

.empty-icon {
  font-size: 4rem;
  color: var(--surface-300);
  margin-bottom: var(--space-4);
  animation: scaleIn var(--transition-slow) ease-out;
}

.empty-state h3 {
  font-size: var(--text-2xl);
  margin-bottom: var(--space-2);
  color: var(--surface-700);
  font-weight: var(--font-semibold);
}

.empty-state p {
  color: var(--surface-500);
  margin-bottom: var(--space-6);
  font-size: var(--text-base);
}

.w-full {
  width: 100%;
}

.font-mono {
  font-family: var(--font-family-mono);
  font-size: var(--text-sm);
  color: var(--surface-600);
  background: var(--surface-100);
  padding: var(--space-1) var(--space-2);
  border-radius: var(--radius-sm);
}

:deep(.p-datatable) {
  animation: fadeIn var(--transition-slow) ease-out;
  animation-delay: 200ms;
  animation-fill-mode: both;
}

:deep(.p-dropdown-panel),
:deep(.p-tooltip) {
  animation: scaleIn var(--transition-fast) ease-out;
}

:deep(.p-input-icon-left > .pi) {
  color: var(--surface-400);
}

:deep(.p-input-icon-left > input) {
  padding-left: 2.5rem;
}

:deep(.p-button-text:not(:hover)) {
  background: transparent;
}

:deep(.p-button-text:hover) {
  background: var(--surface-100);
}

:deep(.p-column-header-content) {
  display: flex;
  align-items: center;
  gap: var(--space-2);
}

@media (max-width: 1200px) {
  .filters-grid {
    grid-template-columns: 1fr 1fr;
  }
}

@media (max-width: 768px) {
  .page-header {
    flex-direction: column;
    align-items: flex-start;
    gap: var(--space-4);
  }

  .filters-grid {
    grid-template-columns: 1fr;
    gap: var(--space-4);
  }

  .table-header {
    flex-direction: column;
    gap: var(--space-2);
    align-items: flex-start;
  }

  .date-range {
    grid-template-columns: 1fr;
    gap: var(--space-2);
  }

  .date-separator {
    display: none;
  }
}
</style>
