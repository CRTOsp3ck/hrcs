<template>
  <div class="admin-page-container">
    <div class="page-header">
      <h1 class="page-title">Reports</h1>
      <p class="page-subtitle">Generate and view comprehensive system reports</p>
    </div>

    <!-- Quick Actions Bar -->
    <div class="toolbar">
      <div class="toolbar-left">
        <Button 
          label="Refresh" 
          icon="pi pi-refresh"
          severity="secondary"
          @click="loadRecentReports"
          :loading="loading"
        />
      </div>
      <div class="toolbar-actions">
        <Button 
          label="Schedule Report" 
          icon="pi pi-calendar"
          severity="info"
          @click="showScheduleDialog = true"
          disabled
        />
        <Button 
          label="Export Settings" 
          icon="pi pi-cog"
          severity="secondary"
          @click="showSettingsDialog = true"
          disabled
        />
      </div>
    </div>

    <!-- Report Types Grid -->
    <div class="reports-grid">
      <Card v-for="report in reportTypes" :key="report.id" class="report-card">
        <template #header>
          <div class="report-header">
            <div class="report-icon-wrapper">
              <i :class="report.icon" class="report-icon"></i>
            </div>
            <div class="report-title-section">
              <h3 class="report-title">{{ report.name }}</h3>
              <span class="report-category">{{ report.category }}</span>
            </div>
          </div>
        </template>
        
        <template #content>
          <div class="report-content">
            <p class="report-description">{{ report.description }}</p>
            
            <div class="report-features">
              <h4>Includes:</h4>
              <ul class="feature-list">
                <li v-for="feature in report.features" :key="feature">
                  <i class="pi pi-check feature-check"></i>
                  <span>{{ feature }}</span>
                </li>
              </ul>
            </div>

            <div class="report-stats" v-if="report.stats">
              <div class="stat-item">
                <span class="stat-label">Last Generated:</span>
                <span class="stat-value">{{ report.stats.lastGenerated || 'Never' }}</span>
              </div>
              <div class="stat-item">
                <span class="stat-label">Generation Time:</span>
                <span class="stat-value">{{ report.stats.avgTime || 'N/A' }}</span>
              </div>
            </div>
          </div>
        </template>

        <template #footer>
          <div class="report-actions">
            <Button 
              :label="getGenerateButtonLabel(report)"
              :icon="getGenerateButtonIcon(report)"
              @click="generateReport(report)"
              :loading="generatingReports.includes(report.id)"
              :disabled="!report.available"
              class="generate-btn"
            />
            
            <div class="secondary-actions">
              <Button 
                icon="pi pi-calendar"
                v-tooltip="'Schedule Report'"
                text
                rounded
                @click="scheduleReport(report)"
                :disabled="!report.available"
              />
              <Button 
                icon="pi pi-info-circle"
                v-tooltip="'Report Details'"
                text
                rounded
                @click="showReportDetails(report)"
              />
              <Button 
                icon="pi pi-download"
                v-tooltip="'Download Sample'"
                text
                rounded
                @click="downloadSample(report)"
                :disabled="!report.sampleAvailable"
              />
            </div>
          </div>
        </template>
      </Card>
    </div>

    <!-- Recent Reports Section -->
    <Card class="recent-reports-section">
      <template #header>
        <div class="section-header">
          <h3 class="section-title">Recent Reports</h3>
          <div class="header-actions">
            <InputText 
              v-model="searchReports"
              placeholder="Search reports..."
              class="search-input"
            >
              <template #before>
                <i class="pi pi-search"></i>
              </template>
            </InputText>
          </div>
        </div>
      </template>
      
      <template #content>
        <DataTable
          :value="filteredRecentReports"
          :loading="loading"
          paginator
          :rows="10"
          :rowsPerPageOptions="[10, 25, 50]"
          responsiveLayout="scroll"
          class="reports-table"
        >
          <Column field="name" header="Report Name" sortable>
            <template #body="slotProps">
              <div class="report-name-cell">
                <i :class="getReportTypeIcon(slotProps.data.type)" class="report-type-icon"></i>
                <div class="report-info">
                  <div class="report-name">{{ slotProps.data.name }}</div>
                  <div class="report-type">{{ formatReportType(slotProps.data.type) }}</div>
                </div>
              </div>
            </template>
          </Column>

          <Column field="generated_at" header="Generated" sortable>
            <template #body="slotProps">
              <div class="date-cell">
                <div class="date">{{ formatDate(slotProps.data.generated_at) }}</div>
                <div class="time">{{ formatTime(slotProps.data.generated_at) }}</div>
              </div>
            </template>
          </Column>

          <Column field="generated_by" header="Generated By" sortable>
            <template #body="slotProps">
              <div class="user-info">
                <Avatar :label="getInitials(slotProps.data.generated_by)" size="small" />
                <span>{{ slotProps.data.generated_by }}</span>
              </div>
            </template>
          </Column>

          <Column field="status" header="Status" sortable>
            <template #body="slotProps">
              <Tag 
                :value="slotProps.data.status" 
                :severity="getStatusSeverity(slotProps.data.status)"
                :icon="getStatusIcon(slotProps.data.status)"
              />
            </template>
          </Column>

          <Column field="file_size" header="Size" sortable>
            <template #body="slotProps">
              <span class="file-size">{{ formatFileSize(slotProps.data.file_size) }}</span>
            </template>
          </Column>

          <Column field="download_count" header="Downloads" sortable>
            <template #body="slotProps">
              <div class="download-info">
                <i class="pi pi-download download-icon"></i>
                <span>{{ slotProps.data.download_count || 0 }}</span>
              </div>
            </template>
          </Column>

          <Column header="Actions" style="width: 150px">
            <template #body="slotProps">
              <div class="table-actions">
                <Button
                  icon="pi pi-download"
                  text
                  rounded
                  @click="downloadReport(slotProps.data)"
                  v-tooltip="'Download Report'"
                  :disabled="slotProps.data.status !== 'completed'"
                />
                <Button
                  icon="pi pi-eye"
                  text
                  rounded
                  @click="previewReport(slotProps.data)"
                  v-tooltip="'Preview Report'"
                  :disabled="slotProps.data.status !== 'completed'"
                />
                <Button
                  icon="pi pi-share-alt"
                  text
                  rounded
                  @click="shareReport(slotProps.data)"
                  v-tooltip="'Share Report'"
                  :disabled="slotProps.data.status !== 'completed'"
                />
                <Button
                  icon="pi pi-trash"
                  text
                  rounded
                  severity="danger"
                  @click="deleteReport(slotProps.data)"
                  v-tooltip="'Delete Report'"
                />
              </div>
            </template>
          </Column>

          <template #empty>
            <div class="empty-state">
              <i class="pi pi-chart-bar empty-icon"></i>
              <p>No reports generated yet</p>
              <Button
                label="Generate Your First Report"
                icon="pi pi-plus"
                @click="showFirstReportDialog = true"
              />
            </div>
          </template>
        </DataTable>
      </template>
    </Card>

    <!-- Report Details Dialog -->
    <Dialog 
      v-model:visible="showDetailsDialog" 
      :header="selectedReport?.name + ' Details'" 
      :modal="true" 
      :style="{ width: '600px' }"
      class="report-details-dialog"
    >
      <div v-if="selectedReport" class="report-details">
        <div class="detail-section">
          <h4>Report Information</h4>
          <div class="detail-grid">
            <div class="detail-item">
              <label>Type:</label>
              <span>{{ selectedReport.name }}</span>
            </div>
            <div class="detail-item">
              <label>Category:</label>
              <span>{{ selectedReport.category }}</span>
            </div>
            <div class="detail-item">
              <label>Data Sources:</label>
              <span>{{ selectedReport.dataSources?.join(', ') || 'Multiple' }}</span>
            </div>
            <div class="detail-item">
              <label>Estimated Time:</label>
              <span>{{ selectedReport.estimatedTime || '2-5 minutes' }}</span>
            </div>
          </div>
        </div>

        <div class="detail-section">
          <h4>Available Formats</h4>
          <div class="format-options">
            <Tag 
              v-for="format in selectedReport.formats" 
              :key="format"
              :value="format.toUpperCase()"
              severity="info"
              class="format-tag"
            />
          </div>
        </div>

        <div class="detail-section">
          <h4>Report Contents</h4>
          <ul class="contents-list">
            <li v-for="feature in selectedReport.features" :key="feature">
              <i class="pi pi-check content-check"></i>
              <span>{{ feature }}</span>
            </li>
          </ul>
        </div>
      </div>
    </Dialog>

    <!-- Generate Report Dialog -->
    <Dialog 
      v-model:visible="showGenerateDialog" 
      header="Generate Report" 
      :modal="true" 
      :style="{ width: '500px' }"
      class="generate-dialog"
    >
      <div v-if="selectedReport" class="generate-form">
        <div class="form-section">
          <h4>{{ selectedReport.name }}</h4>
          <p class="section-description">{{ selectedReport.description }}</p>
        </div>

        <div class="form-section">
          <label class="form-label">Date Range</label>
          <Calendar
            v-model="generateOptions.dateRange"
            selectionMode="range"
            dateFormat="mm/dd/yy"
            placeholder="Select date range"
            showIcon
            :maxDate="maxDate"
            class="date-range-input"
          />
        </div>

        <div class="form-section">
          <label class="form-label">Output Format</label>
          <Dropdown
            v-model="generateOptions.format"
            :options="formatOptions"
            optionLabel="label"
            optionValue="value"
            placeholder="Select format"
            class="format-dropdown"
          />
        </div>

        <div class="form-section" v-if="selectedReport.hasFilters">
          <label class="form-label">Additional Filters</label>
          <MultiSelect
            v-model="generateOptions.filters"
            :options="availableFilters"
            optionLabel="label"
            optionValue="value"
            placeholder="Select filters (optional)"
            class="filters-input"
          />
        </div>

        <div class="form-section">
          <div class="checkbox-group">
            <Checkbox 
              v-model="generateOptions.includeCharts" 
              :binary="true" 
              inputId="include-charts"
            />
            <label for="include-charts">Include charts and visualizations</label>
          </div>
          
          <div class="checkbox-group">
            <Checkbox 
              v-model="generateOptions.emailNotification" 
              :binary="true" 
              inputId="email-notification"
            />
            <label for="email-notification">Send email notification when ready</label>
          </div>
        </div>
      </div>

      <template #footer>
        <div class="dialog-footer">
          <Button 
            label="Cancel" 
            severity="secondary" 
            @click="showGenerateDialog = false"
          />
          <Button 
            label="Generate Report" 
            icon="pi pi-chart-bar"
            @click="confirmGenerate"
            :loading="generatingReports.includes(selectedReport?.id)"
          />
        </div>
      </template>
    </Dialog>

    <!-- Coming Soon Dialog -->
    <Dialog 
      v-model:visible="showComingSoonDialog" 
      header="Coming Soon" 
      :modal="true" 
      :style="{ width: '400px' }"
      class="coming-soon-dialog"
    >
      <div class="coming-soon-content">
        <i class="pi pi-clock coming-soon-icon"></i>
        <h4>Feature In Development</h4>
        <p>This reporting feature is currently being developed and will be available in a future release.</p>
        <p>For now, you can view placeholder data and explore the interface.</p>
      </div>
      
      <template #footer>
        <Button 
          label="Understood" 
          @click="showComingSoonDialog = false"
        />
      </template>
    </Dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { useToast } from 'primevue/usetoast'
import api from '@/api'

// Toast
const toast = useToast()

// Loading states
const loading = ref(false)
const generatingReports = ref<string[]>([])

// Dialog states
const showDetailsDialog = ref(false)
const showGenerateDialog = ref(false)
const showScheduleDialog = ref(false)
const showSettingsDialog = ref(false)
const showFirstReportDialog = ref(false)
const showComingSoonDialog = ref(false)

// Data
const selectedReport = ref(null)
const recentReports = ref([])
const searchReports = ref('')

// Form data
const maxDate = new Date()
const generateOptions = reactive({
  dateRange: [new Date(Date.now() - 30 * 24 * 60 * 60 * 1000), new Date()],
  format: 'pdf',
  filters: [],
  includeCharts: true,
  emailNotification: false
})

// Options
const formatOptions = [
  { label: 'PDF Document', value: 'pdf' },
  { label: 'Excel Spreadsheet', value: 'xlsx' },
  { label: 'CSV Data', value: 'csv' },
  { label: 'JSON Data', value: 'json' }
]

const availableFilters = [
  { label: 'Include Draft Claims', value: 'include_drafts' },
  { label: 'Group by Department', value: 'group_by_dept' },
  { label: 'Show Detailed Breakdown', value: 'detailed_breakdown' },
  { label: 'Include User Profiles', value: 'include_profiles' }
]

// Report types configuration
const reportTypes = ref([
  {
    id: 'claims-summary',
    name: 'Claims Summary Report',
    category: 'Financial',
    description: 'Comprehensive overview of all claims by status, type, and period with financial summaries.',
    icon: 'pi pi-chart-pie',
    available: true,
    sampleAvailable: true,
    hasFilters: true,
    features: [
      'Claims breakdown by status and type',
      'Financial summaries and totals',
      'Approval workflow statistics',
      'Monthly and quarterly trends',
      'Export to multiple formats'
    ],
    formats: ['pdf', 'xlsx', 'csv'],
    dataSources: ['Claims', 'Claim Types', 'Users'],
    estimatedTime: '2-3 minutes',
    stats: {
      lastGenerated: '2 days ago',
      avgTime: '2.5 minutes'
    }
  },
  {
    id: 'user-activity',
    name: 'User Activity Report',
    category: 'Analytics',
    description: 'Detailed analysis of user engagement, claim submission patterns, and system usage.',
    icon: 'pi pi-users',
    available: true,
    sampleAvailable: true,
    hasFilters: true,
    features: [
      'User login and activity patterns',
      'Claim submission frequency',
      'Most active users and departments',
      'System usage statistics',
      'Productivity insights'
    ],
    formats: ['pdf', 'xlsx'],
    dataSources: ['Users', 'Claims', 'Audit Logs'],
    estimatedTime: '3-4 minutes',
    stats: {
      lastGenerated: '1 week ago',
      avgTime: '3.2 minutes'
    }
  },
  {
    id: 'financial-overview',
    name: 'Financial Overview',
    category: 'Financial',
    description: 'Complete financial analysis including total amounts, approvals, budget utilization, and forecasting.',
    icon: 'pi pi-dollar',
    available: true,
    sampleAvailable: false,
    hasFilters: true,
    features: [
      'Total approved and pending amounts',
      'Budget utilization by department',
      'Spending trends and forecasts',
      'Cost center analysis',
      'Variance reports'
    ],
    formats: ['pdf', 'xlsx'],
    dataSources: ['Claims', 'Budgets', 'Departments'],
    estimatedTime: '4-5 minutes',
    stats: {
      lastGenerated: 'Never',
      avgTime: 'N/A'
    }
  },
  {
    id: 'approval-efficiency',
    name: 'Approval Efficiency Report',
    category: 'Operations',
    description: 'Analysis of approval times, workflow performance, and bottleneck identification.',
    icon: 'pi pi-clock',
    available: true,
    sampleAvailable: true,
    hasFilters: false,
    features: [
      'Average approval times by level',
      'Workflow bottleneck analysis',
      'Approver workload distribution',
      'SLA compliance metrics',
      'Process improvement recommendations'
    ],
    formats: ['pdf', 'xlsx'],
    dataSources: ['Claims', 'Approvals', 'Workflow'],
    estimatedTime: '3-4 minutes',
    stats: {
      lastGenerated: '5 days ago',
      avgTime: '3.8 minutes'
    }
  },
  {
    id: 'compliance-audit',
    name: 'Compliance Audit Report',
    category: 'Compliance',
    description: 'Comprehensive compliance analysis and audit trail for regulatory requirements.',
    icon: 'pi pi-shield',
    available: false,
    sampleAvailable: false,
    hasFilters: true,
    features: [
      'Regulatory compliance status',
      'Audit trail documentation',
      'Policy adherence analysis',
      'Risk assessment metrics',
      'Compliance recommendations'
    ],
    formats: ['pdf'],
    dataSources: ['Audit Logs', 'Policies', 'Claims'],
    estimatedTime: '5-7 minutes',
    stats: {
      lastGenerated: 'Never',
      avgTime: 'N/A'
    }
  },
  {
    id: 'department-analysis',
    name: 'Department Analysis',
    category: 'Analytics',
    description: 'Department-wise breakdown of claims, spending, and performance metrics.',
    icon: 'pi pi-sitemap',
    available: false,
    sampleAvailable: false,
    hasFilters: true,
    features: [
      'Department spending analysis',
      'Claims volume by department',
      'Inter-department comparisons',
      'Budget allocation insights',
      'Performance benchmarking'
    ],
    formats: ['pdf', 'xlsx'],
    dataSources: ['Departments', 'Claims', 'Users'],
    estimatedTime: '3-4 minutes',
    stats: {
      lastGenerated: 'Never',
      avgTime: 'N/A'
    }
  }
])

// Computed
const filteredRecentReports = computed(() => {
  if (!searchReports.value) return recentReports.value
  
  return recentReports.value.filter((report: any) =>
    report.name.toLowerCase().includes(searchReports.value.toLowerCase()) ||
    report.type.toLowerCase().includes(searchReports.value.toLowerCase())
  )
})

// Methods
const loadRecentReports = async () => {
  loading.value = true
  try {
    // Simulate API call with mock data
    await new Promise(resolve => setTimeout(resolve, 1000))
    
    recentReports.value = [
      {
        id: 1,
        name: 'Monthly Claims Summary - November 2024',
        type: 'claims-summary',
        generated_at: '2024-11-15T10:30:00Z',
        generated_by: 'John Admin',
        status: 'completed',
        file_size: 2457600, // 2.4 MB
        download_count: 15
      },
      {
        id: 2,
        name: 'User Activity Report - Q4 2024',
        type: 'user-activity',
        generated_at: '2024-11-10T14:15:00Z',
        generated_by: 'Sarah Manager',
        status: 'completed',
        file_size: 1024000, // 1 MB
        download_count: 8
      },
      {
        id: 3,
        name: 'Approval Efficiency Analysis',
        type: 'approval-efficiency',
        generated_at: '2024-11-08T09:45:00Z',
        generated_by: 'Mike Director',
        status: 'failed',
        file_size: 0,
        download_count: 0
      },
      {
        id: 4,
        name: 'Financial Overview - October 2024',
        type: 'financial-overview',
        generated_at: '2024-11-05T16:20:00Z',
        generated_by: 'John Admin',
        status: 'processing',
        file_size: 0,
        download_count: 0
      }
    ]
  } catch (error) {
    console.error('Failed to load recent reports:', error)
    toast.add({
      severity: 'error',
      summary: 'Error',
      detail: 'Failed to load recent reports',
      life: 3000
    })
  } finally {
    loading.value = false
  }
}

const generateReport = (report: any) => {
  if (!report.available) {
    showComingSoonDialog.value = true
    return
  }
  
  selectedReport.value = report
  showGenerateDialog.value = true
}

const confirmGenerate = async () => {
  if (!selectedReport.value) return
  
  generatingReports.value.push(selectedReport.value.id)
  
  try {
    // Simulate report generation
    toast.add({
      severity: 'info',
      summary: 'Report Generation Started',
      detail: `Generating ${selectedReport.value.name}...`,
      life: 3000
    })
    
    // Simulate processing time
    await new Promise(resolve => setTimeout(resolve, 3000))
    
    toast.add({
      severity: 'success',
      summary: 'Report Generated',
      detail: `${selectedReport.value.name} has been generated successfully`,
      life: 5000
    })
    
    showGenerateDialog.value = false
    loadRecentReports()
    
  } catch (error) {
    toast.add({
      severity: 'error',
      summary: 'Generation Failed',
      detail: 'Failed to generate report. Please try again.',
      life: 3000
    })
  } finally {
    generatingReports.value = generatingReports.value.filter(id => id !== selectedReport.value?.id)
  }
}

const scheduleReport = (report: any) => {
  if (!report.available) {
    showComingSoonDialog.value = true
    return
  }
  
  toast.add({
    severity: 'info',
    summary: 'Coming Soon',
    detail: 'Report scheduling feature will be available soon',
    life: 3000
  })
}

const showReportDetails = (report: any) => {
  selectedReport.value = report
  showDetailsDialog.value = true
}

const downloadSample = (report: any) => {
  if (!report.sampleAvailable) {
    toast.add({
      severity: 'warn',
      summary: 'Sample Not Available',
      detail: 'Sample report is not available for this type',
      life: 3000
    })
    return
  }
  
  toast.add({
    severity: 'info',
    summary: 'Downloading Sample',
    detail: `Downloading sample ${report.name}...`,
    life: 2000
  })
}

const downloadReport = (report: any) => {
  toast.add({
    severity: 'success',
    summary: 'Download Started',
    detail: `Downloading ${report.name}...`,
    life: 2000
  })
}

const previewReport = (report: any) => {
  toast.add({
    severity: 'info',
    summary: 'Opening Preview',
    detail: `Opening preview for ${report.name}...`,
    life: 2000
  })
}

const shareReport = (report: any) => {
  toast.add({
    severity: 'info',
    summary: 'Share Link Generated',
    detail: 'Share link copied to clipboard',
    life: 2000
  })
}

const deleteReport = (report: any) => {
  // Would typically show confirmation dialog
  toast.add({
    severity: 'success',
    summary: 'Report Deleted',
    detail: `${report.name} has been deleted`,
    life: 2000
  })
  loadRecentReports()
}

// Utility methods
const getGenerateButtonLabel = (report: any) => {
  return report.available ? 'Generate Report' : 'Coming Soon'
}

const getGenerateButtonIcon = (report: any) => {
  return report.available ? 'pi pi-chart-bar' : 'pi pi-clock'
}

const getReportTypeIcon = (type: string) => {
  const report = reportTypes.value.find(r => r.id === type)
  return report?.icon || 'pi pi-file'
}

const formatReportType = (type: string) => {
  return type.split('-').map(word => 
    word.charAt(0).toUpperCase() + word.slice(1)
  ).join(' ')
}

const formatDate = (dateString: string) => {
  return new Date(dateString).toLocaleDateString()
}

const formatTime = (dateString: string) => {
  return new Date(dateString).toLocaleTimeString()
}

const formatFileSize = (bytes: number) => {
  if (bytes === 0) return 'N/A'
  const k = 1024
  const sizes = ['Bytes', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(1)) + ' ' + sizes[i]
}

const getInitials = (name: string) => {
  return name
    .split(' ')
    .map(word => word.charAt(0))
    .join('')
    .toUpperCase()
    .slice(0, 2)
}

const getStatusSeverity = (status: string) => {
  switch (status) {
    case 'completed': return 'success'
    case 'processing': return 'warning'
    case 'failed': return 'danger'
    default: return 'info'
  }
}

const getStatusIcon = (status: string) => {
  switch (status) {
    case 'completed': return 'pi pi-check'
    case 'processing': return 'pi pi-spin pi-spinner'
    case 'failed': return 'pi pi-times'
    default: return 'pi pi-info-circle'
  }
}

// Lifecycle
onMounted(() => {
  loadRecentReports()
})
</script>

<style scoped>

.page-header {
  margin-bottom: 2rem;
}

.page-title {
  font-size: 2rem;
  font-weight: 700;
  color: var(--surface-900);
  margin: 0 0 0.5rem;
}

.page-subtitle {
  color: var(--surface-600);
  margin: 0;
}

/* Toolbar */
.toolbar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 2rem;
  padding: 1rem;
  background: var(--surface-0);
  border: 1px solid var(--surface-200);
  border-radius: var(--border-radius);
}

.toolbar-actions {
  display: flex;
  gap: 0.5rem;
}

/* Reports Grid */
.reports-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(400px, 1fr));
  gap: 1.5rem;
  margin-bottom: 2rem;
}

.report-card {
  border: 1px solid var(--surface-200);
  transition: all 0.3s ease;
}

.report-card:hover {
  border-color: var(--primary-200);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}

.report-header {
  display: flex;
  align-items: center;
  gap: 1rem;
  padding: 1.5rem 1.5rem 0;
}

.report-icon-wrapper {
  width: 3rem;
  height: 3rem;
  border-radius: 50%;
  background: var(--primary-100);
  display: flex;
  align-items: center;
  justify-content: center;
}

.report-icon {
  font-size: 1.5rem;
  color: var(--primary-600);
}

.report-title-section {
  flex: 1;
}

.report-title {
  font-size: 1.25rem;
  font-weight: 600;
  color: var(--surface-900);
  margin: 0 0 0.25rem;
}

.report-category {
  font-size: 0.875rem;
  color: var(--surface-500);
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.report-content {
  padding: 0 1.5rem 1.5rem;
}

.report-description {
  color: var(--surface-700);
  line-height: 1.5;
  margin-bottom: 1.5rem;
}

.report-features h4 {
  font-size: 1rem;
  color: var(--surface-800);
  margin: 0 0 0.75rem;
}

.feature-list {
  list-style: none;
  padding: 0;
  margin: 0 0 1.5rem;
}

.feature-list li {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.25rem 0;
  font-size: 0.875rem;
  color: var(--surface-700);
}

.feature-check {
  color: var(--green-600);
  font-size: 0.875rem;
}

.report-stats {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 1rem;
  margin-top: 1rem;
  padding-top: 1rem;
  border-top: 1px solid var(--surface-200);
}

.stat-item {
  display: flex;
  flex-direction: column;
  gap: 0.25rem;
}

.stat-label {
  font-size: 0.75rem;
  color: var(--surface-500);
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.stat-value {
  font-size: 0.875rem;
  color: var(--surface-700);
  font-weight: 500;
}

.report-actions {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 1rem 1.5rem;
  border-top: 1px solid var(--surface-200);
}

.generate-btn {
  flex: 1;
  max-width: 200px;
}

.secondary-actions {
  display: flex;
  gap: 0.25rem;
}

/* Recent Reports Section */
.recent-reports-section {
  margin-top: 2rem;
}

.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 1.5rem 1.5rem 0;
}

.section-title {
  font-size: 1.5rem;
  font-weight: 600;
  color: var(--surface-900);
  margin: 0;
}

.header-actions {
  display: flex;
  gap: 1rem;
  align-items: center;
}

.search-input {
  width: 300px;
}

/* Table Styles */
.reports-table {
  margin: 0;
}

.report-name-cell {
  display: flex;
  align-items: center;
  gap: 0.75rem;
}

.report-type-icon {
  font-size: 1.25rem;
  color: var(--primary-600);
}

.report-info {
  display: flex;
  flex-direction: column;
  gap: 0.25rem;
}

.report-name {
  font-weight: 600;
  color: var(--surface-900);
}

.report-type {
  font-size: 0.75rem;
  color: var(--surface-500);
}

.date-cell {
  display: flex;
  flex-direction: column;
  gap: 0.25rem;
}

.date-cell .date {
  font-weight: 500;
  color: var(--surface-900);
}

.date-cell .time {
  font-size: 0.75rem;
  color: var(--surface-500);
}

.user-info {
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.file-size {
  font-family: monospace;
  font-size: 0.875rem;
}

.download-info {
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.download-icon {
  color: var(--surface-500);
}

.table-actions {
  display: flex;
  gap: 0.25rem;
}

.empty-state {
  text-align: center;
  padding: 3rem;
  color: var(--surface-500);
}

.empty-icon {
  font-size: 3rem;
  margin-bottom: 1rem;
  color: var(--surface-400);
}

/* Dialog Styles */
.report-details {
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
}

.detail-section h4 {
  margin: 0 0 1rem;
  color: var(--surface-900);
  font-weight: 600;
}

.detail-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 1rem;
}

.detail-item {
  display: flex;
  flex-direction: column;
  gap: 0.25rem;
}

.detail-item label {
  font-weight: 600;
  color: var(--surface-700);
  font-size: 0.875rem;
}

.format-options {
  display: flex;
  gap: 0.5rem;
  flex-wrap: wrap;
}

.format-tag {
  font-size: 0.75rem;
}

.contents-list {
  list-style: none;
  padding: 0;
  margin: 0;
}

.contents-list li {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.5rem 0;
  color: var(--surface-700);
}

.content-check {
  color: var(--green-600);
}

/* Generate Dialog */
.generate-form {
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
}

.form-section {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.form-section h4 {
  margin: 0;
  color: var(--surface-900);
}

.section-description {
  color: var(--surface-600);
  font-size: 0.875rem;
  margin: 0;
}

.form-label {
  font-weight: 600;
  color: var(--surface-700);
  font-size: 0.875rem;
}

.checkbox-group {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  margin: 0.5rem 0;
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 0.5rem;
}

/* Coming Soon Dialog */
.coming-soon-content {
  text-align: center;
  padding: 1rem;
}

.coming-soon-icon {
  font-size: 3rem;
  color: var(--surface-400);
  margin-bottom: 1rem;
}

.coming-soon-content h4 {
  color: var(--surface-900);
  margin: 0 0 1rem;
}

.coming-soon-content p {
  color: var(--surface-600);
  line-height: 1.5;
  margin: 0 0 1rem;
}

/* Responsive Design */
@media (max-width: 1200px) {
  .reports-grid {
    grid-template-columns: repeat(auto-fit, minmax(350px, 1fr));
  }
  
  .detail-grid {
    grid-template-columns: 1fr;
  }
}

@media (max-width: 768px) {
  .reports-grid {
    grid-template-columns: 1fr;
  }
  
  .toolbar {
    flex-direction: column;
    gap: 1rem;
    align-items: stretch;
  }
  
  .section-header {
    flex-direction: column;
    gap: 1rem;
    align-items: stretch;
  }
  
  .search-input {
    width: 100%;
  }
  
  .report-actions {
    flex-direction: column;
    gap: 1rem;
  }
  
  .generate-btn {
    max-width: none;
  }
}
</style>