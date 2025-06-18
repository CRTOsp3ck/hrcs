<template>
  <div class="admin-page-container">
    <div class="page-header">
      <h1 class="page-title">Integrations</h1>
      <p class="page-subtitle">Configure external system integrations and single sign-on</p>
    </div>

    <!-- Integration Overview -->
    <div class="integration-overview">
      <div class="overview-stats">
        <div class="stat-card">
          <div class="stat-icon configured">
            <i class="pi pi-check"></i>
          </div>
          <div class="stat-content">
            <div class="stat-value">{{ configuredIntegrations }}</div>
            <div class="stat-label">Configured</div>
          </div>
        </div>
        
        <div class="stat-card">
          <div class="stat-icon available">
            <i class="pi pi-cog"></i>
          </div>
          <div class="stat-content">
            <div class="stat-value">{{ availableIntegrations }}</div>
            <div class="stat-label">Available</div>
          </div>
        </div>
        
        <div class="stat-card">
          <div class="stat-icon coming-soon">
            <i class="pi pi-clock"></i>
          </div>
          <div class="stat-content">
            <div class="stat-value">{{ comingSoonIntegrations }}</div>
            <div class="stat-label">Coming Soon</div>
          </div>
        </div>
      </div>
    </div>

    <!-- Microsoft Integration Section -->
    <Card class="integration-card microsoft-card">
      <template #header>
        <div class="integration-header">
          <div class="integration-brand">
            <div class="brand-icon microsoft">
              <i class="pi pi-microsoft"></i>
            </div>
            <div class="brand-info">
              <h3 class="integration-title">Microsoft Integration</h3>
              <div class="integration-status">
                <Tag 
                  :value="microsoftIntegration.status" 
                  :severity="getStatusSeverity(microsoftIntegration.status)"
                  :icon="getStatusIcon(microsoftIntegration.status)"
                />
              </div>
            </div>
          </div>
          <div class="integration-actions">
            <Button 
              icon="pi pi-refresh"
              text
              rounded
              @click="refreshMicrosoftStatus"
              :loading="refreshing"
              v-tooltip="'Refresh Status'"
            />
          </div>
        </div>
      </template>
      
      <template #content>
        <div class="integration-content">
          <!-- Feature Overview -->
          <div class="feature-overview">
            <div class="feature-grid">
              <div 
                v-for="feature in microsoftFeatures" 
                :key="feature.id"
                class="feature-item"
                :class="{ 'feature-disabled': !feature.available }"
              >
                <div class="feature-icon-wrapper">
                  <i :class="feature.icon" class="feature-icon"></i>
                </div>
                <div class="feature-details">
                  <h4 class="feature-name">{{ feature.name }}</h4>
                  <p class="feature-description">{{ feature.description }}</p>
                  <div class="feature-status">
                    <Tag 
                      :value="feature.status" 
                      :severity="getFeatureStatusSeverity(feature.status)"
                      size="small"
                    />
                  </div>
                </div>
              </div>
            </div>
          </div>

          <!-- Configuration Status -->
          <div class="config-status">
            <h4>Configuration Status</h4>
            <div class="status-grid">
              <div class="status-item">
                <label>Application Registration:</label>
                <div class="status-value">
                  <Tag value="Pending" severity="warning" />
                  <span class="status-note">Azure AD app registration required</span>
                </div>
              </div>
              
              <div class="status-item">
                <label>Tenant Configuration:</label>
                <div class="status-value">
                  <Tag value="Not Set" severity="secondary" />
                  <span class="status-note">Microsoft 365 tenant ID needed</span>
                </div>
              </div>
              
              <div class="status-item">
                <label>API Permissions:</label>
                <div class="status-value">
                  <Tag value="Not Granted" severity="danger" />
                  <span class="status-note">Admin consent required</span>
                </div>
              </div>
              
              <div class="status-item">
                <label>Teams Webhook:</label>
                <div class="status-value">
                  <Tag value="Not Configured" severity="secondary" />
                  <span class="status-note">Teams integration endpoint needed</span>
                </div>
              </div>
            </div>
          </div>

          <!-- Setup Guide -->
          <div class="setup-guide">
            <h4>Setup Instructions</h4>
            <div class="setup-steps">
              <div 
                v-for="(step, index) in setupSteps" 
                :key="step.id"
                class="setup-step"
                :class="{ 'step-completed': step.completed, 'step-current': step.current }"
              >
                <div class="step-number">{{ index + 1 }}</div>
                <div class="step-content">
                  <h5 class="step-title">{{ step.title }}</h5>
                  <p class="step-description">{{ step.description }}</p>
                  <div v-if="step.actions" class="step-actions">
                    <Button 
                      v-for="action in step.actions"
                      :key="action.label"
                      :label="action.label"
                      :icon="action.icon"
                      :severity="action.severity || 'primary'"
                      size="small"
                      @click="executeStepAction(action.action)"
                      :disabled="!action.enabled"
                    />
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </template>

      <template #footer>
        <div class="integration-footer">
          <div class="footer-info">
            <Message 
              severity="info" 
              :closable="false"
              class="integration-message"
            >
              Microsoft integration is currently in development. Contact your system administrator for setup assistance.
            </Message>
          </div>
          <div class="footer-actions">
            <Button 
              label="Configuration Guide" 
              icon="pi pi-book"
              severity="secondary"
              @click="showConfigGuide = true"
            />
            <Button 
              label="Test Connection" 
              icon="pi pi-bolt"
              @click="testMicrosoftConnection"
              :loading="testing"
              disabled
            />
            <Button 
              label="Configure Integration" 
              icon="pi pi-cog"
              @click="showConfigDialog = true"
              disabled
            />
          </div>
        </div>
      </template>
    </Card>

    <!-- Other Integrations Section -->
    <Card class="integration-card other-integrations">
      <template #header>
        <div class="integration-header">
          <div class="integration-brand">
            <div class="brand-icon general">
              <i class="pi pi-cloud"></i>
            </div>
            <div class="brand-info">
              <h3 class="integration-title">Other Integrations</h3>
              <div class="integration-status">
                <Tag value="Coming Soon" severity="info" icon="pi pi-clock" />
              </div>
            </div>
          </div>
        </div>
      </template>
      
      <template #content>
        <div class="other-integrations-content">
          <p class="section-intro">
            Additional integrations will be available in future releases to enhance your workflow and connectivity.
          </p>
          
          <div class="coming-soon-grid">
            <div 
              v-for="integration in otherIntegrations" 
              :key="integration.id"
              class="coming-soon-item"
            >
              <div class="integration-preview">
                <div class="preview-icon">
                  <i :class="integration.icon"></i>
                </div>
                <div class="preview-content">
                  <h4 class="preview-title">{{ integration.name }}</h4>
                  <p class="preview-description">{{ integration.description }}</p>
                  <div class="preview-features">
                    <span 
                      v-for="feature in integration.features" 
                      :key="feature"
                      class="feature-tag"
                    >
                      {{ feature }}
                    </span>
                  </div>
                  <div class="preview-status">
                    <Tag 
                      :value="integration.eta" 
                      severity="info" 
                      size="small"
                    />
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </template>
    </Card>

    <!-- Configuration Guide Dialog -->
    <Dialog 
      v-model:visible="showConfigGuide" 
      header="Microsoft Integration Configuration Guide" 
      :modal="true" 
      :style="{ width: '800px' }"
      class="config-guide-dialog"
    >
      <div class="config-guide-content">
        <div class="guide-section">
          <h4>Prerequisites</h4>
          <ul class="prerequisites-list">
            <li>Microsoft 365 or Azure Active Directory tenant</li>
            <li>Global Administrator or Application Administrator role</li>
            <li>Access to Azure Portal</li>
            <li>Teams administrative permissions (for Teams integration)</li>
          </ul>
        </div>

        <div class="guide-section">
          <h4>Step-by-Step Configuration</h4>
          <div class="config-steps">
            <div class="config-step">
              <div class="step-header">
                <span class="step-badge">1</span>
                <h5>Azure AD App Registration</h5>
              </div>
              <div class="step-content">
                <p>Register your application in Azure Active Directory:</p>
                <ul>
                  <li>Navigate to Azure Portal → App Registrations</li>
                  <li>Click "New registration"</li>
                  <li>Enter application name: "HRCS Claims System"</li>
                  <li>Select "Accounts in this organizational directory only"</li>
                  <li>Set redirect URI: <code>{{ getRedirectUri() }}</code></li>
                </ul>
              </div>
            </div>

            <div class="config-step">
              <div class="step-header">
                <span class="step-badge">2</span>
                <h5>API Permissions</h5>
              </div>
              <div class="step-content">
                <p>Configure required permissions:</p>
                <ul>
                  <li>Microsoft Graph → User.Read (Delegated)</li>
                  <li>Microsoft Graph → User.ReadBasic.All (Application)</li>
                  <li>Microsoft Graph → Group.Read.All (Application)</li>
                  <li>Microsoft Graph → ChannelMessage.Send (for Teams)</li>
                </ul>
                <p><strong>Important:</strong> Grant admin consent for all permissions.</p>
              </div>
            </div>

            <div class="config-step">
              <div class="step-header">
                <span class="step-badge">3</span>
                <h5>Certificates & Secrets</h5>
              </div>
              <div class="step-content">
                <p>Create client secret:</p>
                <ul>
                  <li>Go to "Certificates & secrets"</li>
                  <li>Click "New client secret"</li>
                  <li>Set description: "HRCS Production"</li>
                  <li>Set expiration: "24 months"</li>
                  <li>Copy the secret value immediately</li>
                </ul>
              </div>
            </div>
          </div>
        </div>

        <div class="guide-section">
          <h4>Configuration Values</h4>
          <div class="config-values">
            <p>You will need these values for system configuration:</p>
            <div class="value-grid">
              <div class="value-item">
                <label>Tenant ID:</label>
                <code>xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx</code>
              </div>
              <div class="value-item">
                <label>Client ID:</label>
                <code>xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx</code>
              </div>
              <div class="value-item">
                <label>Client Secret:</label>
                <code>xxxxxxxxxxxxxxxxxxxxxxxxxxxxxx</code>
              </div>
            </div>
          </div>
        </div>
      </div>
      
      <template #footer>
        <div class="dialog-footer">
          <Button 
            label="Download PDF Guide" 
            icon="pi pi-download"
            severity="secondary"
            @click="downloadGuide"
          />
          <Button 
            label="Close" 
            @click="showConfigGuide = false"
          />
        </div>
      </template>
    </Dialog>

    <!-- Configuration Dialog -->
    <Dialog 
      v-model:visible="showConfigDialog" 
      header="Configure Microsoft Integration" 
      :modal="true" 
      :style="{ width: '600px' }"
      class="config-dialog"
    >
      <div class="config-form">
        <Message 
          severity="warn" 
          :closable="false"
          class="config-warning"
        >
          This feature is currently in development. Configuration will be available in a future release.
        </Message>
        
        <div class="form-section">
          <label class="form-label">Tenant ID</label>
          <InputText 
            v-model="configForm.tenantId"
            placeholder="Enter your Microsoft 365 Tenant ID"
            disabled
          />
        </div>

        <div class="form-section">
          <label class="form-label">Application (Client) ID</label>
          <InputText 
            v-model="configForm.clientId"
            placeholder="Enter your Azure AD Application ID"
            disabled
          />
        </div>

        <div class="form-section">
          <label class="form-label">Client Secret</label>
          <Password 
            v-model="configForm.clientSecret"
            placeholder="Enter your client secret"
            :feedback="false"
            toggleMask
            disabled
          />
        </div>

        <div class="form-section">
          <div class="checkbox-group">
            <Checkbox 
              v-model="configForm.enableSSO" 
              :binary="true" 
              inputId="enable-sso"
              disabled
            />
            <label for="enable-sso">Enable Single Sign-On (SSO)</label>
          </div>
          
          <div class="checkbox-group">
            <Checkbox 
              v-model="configForm.enableTeams" 
              :binary="true" 
              inputId="enable-teams"
              disabled
            />
            <label for="enable-teams">Enable Microsoft Teams notifications</label>
          </div>
          
          <div class="checkbox-group">
            <Checkbox 
              v-model="configForm.enableUserSync" 
              :binary="true" 
              inputId="enable-user-sync"
              disabled
            />
            <label for="enable-user-sync">Enable user synchronization</label>
          </div>
        </div>
      </div>
      
      <template #footer>
        <div class="dialog-footer">
          <Button 
            label="Cancel" 
            severity="secondary" 
            @click="showConfigDialog = false"
          />
          <Button 
            label="Save Configuration" 
            icon="pi pi-save"
            disabled
          />
        </div>
      </template>
    </Dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { useToast } from 'primevue/usetoast'

// Toast
const toast = useToast()

// Loading states
const refreshing = ref(false)
const testing = ref(false)

// Dialog states
const showConfigGuide = ref(false)
const showConfigDialog = ref(false)

// Data
const microsoftIntegration = reactive({
  status: 'Not Configured',
  configured: false,
  lastChecked: null
})

const configForm = reactive({
  tenantId: '',
  clientId: '',
  clientSecret: '',
  enableSSO: true,
  enableTeams: true,
  enableUserSync: false
})

// Microsoft features
const microsoftFeatures = ref([
  {
    id: 'sso',
    name: 'Single Sign-On (SSO)',
    description: 'Allow users to sign in with their Microsoft 365 credentials',
    icon: 'pi pi-sign-in',
    status: 'Not Configured',
    available: true
  },
  {
    id: 'teams',
    name: 'Teams Notifications',
    description: 'Send claim notifications and updates to Microsoft Teams channels',
    icon: 'pi pi-comments',
    status: 'Not Configured',
    available: true
  },
  {
    id: 'ad-sync',
    name: 'Active Directory Sync',
    description: 'Synchronize user accounts and groups from Azure Active Directory',
    icon: 'pi pi-users',
    status: 'Coming Soon',
    available: false
  },
  {
    id: 'calendar',
    name: 'Calendar Integration',
    description: 'Schedule claim reviews and approvals in Outlook calendar',
    icon: 'pi pi-calendar',
    status: 'Coming Soon',
    available: false
  }
])

// Setup steps
const setupSteps = ref([
  {
    id: 'register',
    title: 'Azure AD App Registration',
    description: 'Register your application in Azure Active Directory portal',
    completed: false,
    current: true,
    actions: [
      {
        label: 'Open Azure Portal',
        icon: 'pi pi-external-link',
        action: 'openAzurePortal',
        enabled: true,
        severity: 'info'
      }
    ]
  },
  {
    id: 'permissions',
    title: 'Configure API Permissions',
    description: 'Set up required Microsoft Graph API permissions and grant admin consent',
    completed: false,
    current: false,
    actions: [
      {
        label: 'View Required Permissions',
        icon: 'pi pi-list',
        action: 'showPermissions',
        enabled: true,
        severity: 'secondary'
      }
    ]
  },
  {
    id: 'credentials',
    title: 'Generate Credentials',
    description: 'Create client secret and gather tenant/application IDs',
    completed: false,
    current: false,
    actions: []
  },
  {
    id: 'configure',
    title: 'Configure HRCS',
    description: 'Enter Microsoft integration settings in HRCS admin panel',
    completed: false,
    current: false,
    actions: [
      {
        label: 'Configure Now',
        icon: 'pi pi-cog',
        action: 'configure',
        enabled: false,
        severity: 'primary'
      }
    ]
  },
  {
    id: 'test',
    title: 'Test Connection',
    description: 'Verify the integration is working correctly',
    completed: false,
    current: false,
    actions: [
      {
        label: 'Run Test',
        icon: 'pi pi-bolt',
        action: 'test',
        enabled: false,
        severity: 'success'
      }
    ]
  }
])

// Other integrations
const otherIntegrations = ref([
  {
    id: 'slack',
    name: 'Slack Integration',
    description: 'Send notifications and interact with claims through Slack',
    icon: 'pi pi-slack',
    features: ['Notifications', 'Bot Commands', 'Workflow Actions'],
    eta: 'Q2 2025'
  },
  {
    id: 'email',
    name: 'Email Providers',
    description: 'Integrate with SendGrid, Mailgun, or other email services',
    icon: 'pi pi-envelope',
    features: ['SendGrid', 'Mailgun', 'Custom SMTP'],
    eta: 'Q1 2025'
  },
  {
    id: 'storage',
    name: 'Document Storage',
    description: 'Connect to SharePoint, Google Drive, or cloud storage',
    icon: 'pi pi-cloud-upload',
    features: ['SharePoint', 'Google Drive', 'Box', 'AWS S3'],
    eta: 'Q3 2025'
  },
  {
    id: 'accounting',
    name: 'Accounting Systems',
    description: 'Sync with QuickBooks, SAP, or other financial systems',
    icon: 'pi pi-calculator',
    features: ['QuickBooks', 'SAP', 'Xero', 'NetSuite'],
    eta: 'Q4 2025'
  }
])

// Computed
const configuredIntegrations = computed(() => {
  return microsoftFeatures.value.filter(f => f.status === 'Configured').length
})

const availableIntegrations = computed(() => {
  return microsoftFeatures.value.filter(f => f.available).length
})

const comingSoonIntegrations = computed(() => {
  return microsoftFeatures.value.filter(f => !f.available).length + otherIntegrations.value.length
})

// Methods
const refreshMicrosoftStatus = async () => {
  refreshing.value = true
  try {
    // Simulate API call
    await new Promise(resolve => setTimeout(resolve, 1500))
    
    microsoftIntegration.lastChecked = new Date().toISOString()
    
    toast.add({
      severity: 'info',
      summary: 'Status Refreshed',
      detail: 'Microsoft integration status updated',
      life: 2000
    })
  } catch (error) {
    toast.add({
      severity: 'error',
      summary: 'Refresh Failed',
      detail: 'Could not refresh integration status',
      life: 3000
    })
  } finally {
    refreshing.value = false
  }
}

const testMicrosoftConnection = async () => {
  testing.value = true
  try {
    await new Promise(resolve => setTimeout(resolve, 2000))
    
    toast.add({
      severity: 'success',
      summary: 'Connection Test',
      detail: 'Microsoft integration test completed',
      life: 3000
    })
  } catch (error) {
    toast.add({
      severity: 'error',
      summary: 'Test Failed',
      detail: 'Connection test failed',
      life: 3000
    })
  } finally {
    testing.value = false
  }
}

const executeStepAction = (action: string) => {
  switch (action) {
    case 'openAzurePortal':
      window.open('https://portal.azure.com/#blade/Microsoft_AAD_IAM/ActiveDirectoryMenuBlade/RegisteredApps', '_blank')
      break
    case 'showPermissions':
      toast.add({
        severity: 'info',
        summary: 'Permissions Guide',
        detail: 'Opening configuration guide with required permissions',
        life: 3000
      })
      showConfigGuide.value = true
      break
    case 'configure':
      showConfigDialog.value = true
      break
    case 'test':
      testMicrosoftConnection()
      break
    default:
      toast.add({
        severity: 'info',
        summary: 'Coming Soon',
        detail: 'This feature will be available in a future release',
        life: 2000
      })
  }
}

const downloadGuide = () => {
  toast.add({
    severity: 'success',
    summary: 'Guide Downloaded',
    detail: 'Configuration guide has been downloaded',
    life: 2000
  })
}

const getRedirectUri = () => {
  return `${window.location.origin}/auth/microsoft/callback`
}

// Utility methods
const getStatusSeverity = (status: string) => {
  switch (status.toLowerCase()) {
    case 'configured': return 'success'
    case 'not configured': return 'warning'
    case 'error': return 'danger'
    default: return 'info'
  }
}

const getStatusIcon = (status: string) => {
  switch (status.toLowerCase()) {
    case 'configured': return 'pi pi-check'
    case 'not configured': return 'pi pi-exclamation-triangle'
    case 'error': return 'pi pi-times'
    default: return 'pi pi-info-circle'
  }
}

const getFeatureStatusSeverity = (status: string) => {
  switch (status.toLowerCase()) {
    case 'configured': return 'success'
    case 'not configured': return 'warning'
    case 'coming soon': return 'info'
    default: return 'secondary'
  }
}

// Lifecycle
onMounted(() => {
  // Initialize integration status
  microsoftIntegration.lastChecked = new Date().toISOString()
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

/* Integration Overview */
.integration-overview {
  margin-bottom: 2rem;
}

.overview-stats {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 1rem;
}

.stat-card {
  background: var(--surface-0);
  border: 1px solid var(--surface-200);
  border-radius: var(--border-radius);
  padding: 1.5rem;
  display: flex;
  align-items: center;
  gap: 1rem;
}

.stat-icon {
  width: 3rem;
  height: 3rem;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 1.25rem;
}

.stat-icon.configured {
  background: var(--green-100);
  color: var(--green-600);
}

.stat-icon.available {
  background: var(--blue-100);
  color: var(--blue-600);
}

.stat-icon.coming-soon {
  background: var(--orange-100);
  color: var(--orange-600);
}

.stat-value {
  font-size: 1.5rem;
  font-weight: 700;
  color: var(--surface-900);
}

.stat-label {
  color: var(--surface-600);
  font-size: 0.875rem;
}

/* Integration Cards */
.integration-card {
  margin-bottom: 2rem;
  border: 1px solid var(--surface-200);
}

.microsoft-card {
  border-left: 4px solid #0078d4;
}

.integration-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 1.5rem 1.5rem 0;
}

.integration-brand {
  display: flex;
  align-items: center;
  gap: 1rem;
}

.brand-icon {
  width: 3rem;
  height: 3rem;
  border-radius: 0.5rem;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 1.5rem;
}

.brand-icon.microsoft {
  background: #0078d4;
  color: white;
}

.brand-icon.general {
  background: var(--surface-100);
  color: var(--surface-600);
}

.integration-title {
  font-size: 1.5rem;
  font-weight: 600;
  color: var(--surface-900);
  margin: 0 0 0.25rem;
}

.integration-status {
  margin-top: 0.5rem;
}

/* Feature Overview */
.feature-overview {
  margin-bottom: 2rem;
}

.feature-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
  gap: 1rem;
}

.feature-item {
  background: var(--surface-50);
  border: 1px solid var(--surface-200);
  border-radius: var(--border-radius);
  padding: 1rem;
  transition: all 0.3s ease;
}

.feature-item:hover:not(.feature-disabled) {
  border-color: var(--primary-200);
  background: var(--surface-0);
}

.feature-disabled {
  opacity: 0.6;
}

.feature-icon-wrapper {
  width: 2.5rem;
  height: 2.5rem;
  border-radius: 50%;
  background: var(--primary-100);
  display: flex;
  align-items: center;
  justify-content: center;
  margin-bottom: 1rem;
}

.feature-icon {
  color: var(--primary-600);
  font-size: 1.25rem;
}

.feature-name {
  font-size: 1rem;
  font-weight: 600;
  color: var(--surface-900);
  margin: 0 0 0.5rem;
}

.feature-description {
  color: var(--surface-600);
  font-size: 0.875rem;
  line-height: 1.4;
  margin: 0 0 0.75rem;
}

/* Configuration Status */
.config-status {
  margin-bottom: 2rem;
}

.config-status h4 {
  margin: 0 0 1rem;
  color: var(--surface-900);
  font-weight: 600;
}

.status-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
  gap: 1rem;
}

.status-item {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
  padding: 1rem;
  background: var(--surface-50);
  border: 1px solid var(--surface-200);
  border-radius: var(--border-radius);
}

.status-item label {
  font-weight: 600;
  color: var(--surface-700);
  font-size: 0.875rem;
}

.status-value {
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.status-note {
  font-size: 0.75rem;
  color: var(--surface-500);
}

/* Setup Guide */
.setup-guide h4 {
  margin: 0 0 1.5rem;
  color: var(--surface-900);
  font-weight: 600;
}

.setup-steps {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.setup-step {
  display: flex;
  gap: 1rem;
  padding: 1rem;
  border: 1px solid var(--surface-200);
  border-radius: var(--border-radius);
  transition: border-color 0.3s ease;
}

.step-current {
  border-color: var(--primary-300);
  background: var(--primary-50);
}

.step-completed {
  border-color: var(--green-300);
  background: var(--green-50);
}

.step-number {
  width: 2rem;
  height: 2rem;
  border-radius: 50%;
  background: var(--surface-200);
  color: var(--surface-700);
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: 600;
  font-size: 0.875rem;
  flex-shrink: 0;
}

.step-current .step-number {
  background: var(--primary-500);
  color: white;
}

.step-completed .step-number {
  background: var(--green-500);
  color: white;
}

.step-content {
  flex: 1;
}

.step-title {
  font-size: 1rem;
  font-weight: 600;
  color: var(--surface-900);
  margin: 0 0 0.5rem;
}

.step-description {
  color: var(--surface-600);
  font-size: 0.875rem;
  line-height: 1.4;
  margin: 0 0 0.75rem;
}

.step-actions {
  display: flex;
  gap: 0.5rem;
  flex-wrap: wrap;
}

/* Integration Footer */
.integration-footer {
  display: flex;
  flex-direction: column;
  gap: 1rem;
  padding: 1rem 1.5rem;
  border-top: 1px solid var(--surface-200);
}

.footer-info .integration-message {
  margin: 0;
}

.footer-actions {
  display: flex;
  gap: 0.5rem;
  flex-wrap: wrap;
}

/* Other Integrations */
.other-integrations-content {
  padding: 0 1.5rem 1.5rem;
}

.section-intro {
  color: var(--surface-600);
  line-height: 1.5;
  margin-bottom: 1.5rem;
}

.coming-soon-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
  gap: 1rem;
}

.coming-soon-item {
  border: 1px solid var(--surface-200);
  border-radius: var(--border-radius);
  overflow: hidden;
  transition: all 0.3s ease;
}

.coming-soon-item:hover {
  border-color: var(--surface-300);
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.integration-preview {
  padding: 1.5rem;
}

.preview-icon {
  width: 3rem;
  height: 3rem;
  border-radius: 0.5rem;
  background: var(--surface-100);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 1.5rem;
  color: var(--surface-600);
  margin-bottom: 1rem;
}

.preview-title {
  font-size: 1.125rem;
  font-weight: 600;
  color: var(--surface-900);
  margin: 0 0 0.5rem;
}

.preview-description {
  color: var(--surface-600);
  font-size: 0.875rem;
  line-height: 1.4;
  margin: 0 0 1rem;
}

.preview-features {
  display: flex;
  flex-wrap: wrap;
  gap: 0.5rem;
  margin-bottom: 1rem;
}

.feature-tag {
  background: var(--surface-100);
  color: var(--surface-700);
  padding: 0.25rem 0.5rem;
  border-radius: 0.25rem;
  font-size: 0.75rem;
  font-weight: 500;
}

/* Dialog Styles */
.config-guide-content {
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
}

.guide-section h4 {
  margin: 0 0 1rem;
  color: var(--surface-900);
  font-weight: 600;
}

.prerequisites-list {
  margin: 0;
  padding-left: 1.5rem;
}

.prerequisites-list li {
  margin-bottom: 0.5rem;
  color: var(--surface-700);
}

.config-steps {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.config-step {
  border: 1px solid var(--surface-200);
  border-radius: var(--border-radius);
  overflow: hidden;
}

.step-header {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  padding: 1rem;
  background: var(--surface-50);
  border-bottom: 1px solid var(--surface-200);
}

.step-badge {
  width: 1.5rem;
  height: 1.5rem;
  border-radius: 50%;
  background: var(--primary-500);
  color: white;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 0.75rem;
  font-weight: 600;
}

.step-header h5 {
  margin: 0;
  color: var(--surface-900);
}

.step-content {
  padding: 1rem;
}

.step-content p {
  margin: 0 0 0.75rem;
  color: var(--surface-700);
}

.step-content ul {
  margin: 0;
  padding-left: 1.5rem;
}

.step-content li {
  margin-bottom: 0.25rem;
  color: var(--surface-600);
}

.config-values {
  background: var(--surface-50);
  border: 1px solid var(--surface-200);
  border-radius: var(--border-radius);
  padding: 1rem;
}

.value-grid {
  display: flex;
  flex-direction: column;
  gap: 0.75rem;
  margin-top: 0.75rem;
}

.value-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.value-item label {
  font-weight: 500;
  color: var(--surface-700);
}

.value-item code {
  background: var(--surface-100);
  padding: 0.25rem 0.5rem;
  border-radius: 0.25rem;
  font-family: monospace;
  font-size: 0.875rem;
}

/* Config Form */
.config-form {
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
}

.config-warning {
  margin: 0 0 1rem;
}

.form-section {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
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

/* Responsive Design */
@media (max-width: 1200px) {
  .feature-grid {
    grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  }
  
  .coming-soon-grid {
    grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
  }
}

@media (max-width: 768px) {
  .overview-stats {
    grid-template-columns: 1fr;
  }
  
  .feature-grid {
    grid-template-columns: 1fr;
  }
  
  .status-grid {
    grid-template-columns: 1fr;
  }
  
  .coming-soon-grid {
    grid-template-columns: 1fr;
  }
  
  .integration-header {
    flex-direction: column;
    gap: 1rem;
    align-items: flex-start;
  }
  
  .footer-actions {
    flex-direction: column;
  }
  
  .value-item {
    flex-direction: column;
    align-items: flex-start;
    gap: 0.5rem;
  }
}
</style>