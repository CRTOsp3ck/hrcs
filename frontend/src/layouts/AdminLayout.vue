<template>
  <div class="admin-layout">
    <!-- Admin Sidebar -->
    <nav class="admin-sidebar">
      <div class="nav-sections">
        <!-- Overview Section -->
        <div class="nav-section">
          <h3 class="section-title">Overview</h3>
          <ul class="nav-items">
            <li>
              <router-link to="/admin/dashboard" class="nav-item">
                <i class="pi pi-chart-bar"></i>
                <span>Dashboard</span>
              </router-link>
            </li>
          </ul>
        </div>
        
        <!-- User Management Section -->
        <div class="nav-section">
          <h3 class="section-title">User Management</h3>
          <ul class="nav-items">
            <li>
              <router-link to="/admin/users" class="nav-item">
                <i class="pi pi-users"></i>
                <span>Users</span>
              </router-link>
            </li>
            <li>
              <router-link to="/admin/groups" class="nav-item">
                <i class="pi pi-sitemap"></i>
                <span>Groups</span>
              </router-link>
            </li>
          </ul>
        </div>
        
        <!-- Claims Configuration Section -->
        <div class="nav-section">
          <h3 class="section-title">Claims Configuration</h3>
          <ul class="nav-items">
            <li>
              <router-link to="/admin/claim-types" class="nav-item">
                <i class="pi pi-tags"></i>
                <span>Claim Types</span>
              </router-link>
            </li>
            <li>
              <router-link to="/admin/approval-levels" class="nav-item">
                <i class="pi pi-shield"></i>
                <span>Approval Levels</span>
              </router-link>
            </li>
          </ul>
        </div>
        
        <!-- Claims Management Section -->
        <div class="nav-section">
          <h3 class="section-title">Claims Management</h3>
          <ul class="nav-items">
            <li>
              <router-link to="/admin/claims" class="nav-item">
                <i class="pi pi-file-text"></i>
                <span>All Claims</span>
                <Badge 
                  v-if="pendingClaimsCount > 0" 
                  :value="pendingClaimsCount" 
                  severity="warning" 
                  class="nav-badge"
                />
              </router-link>
            </li>
          </ul>
        </div>
        
        <!-- Analytics & Reporting Section (Phase 3) -->
        <div class="nav-section">
          <h3 class="section-title">Analytics & Reporting</h3>
          <ul class="nav-items">
            <li>
              <router-link to="/admin/reports" class="nav-item">
                <i class="pi pi-chart-bar"></i>
                <span>Reports</span>
              </router-link>
            </li>
            <li>
              <router-link to="/admin/audit-log" class="nav-item">
                <i class="pi pi-history"></i>
                <span>Audit Log</span>
              </router-link>
            </li>
          </ul>
        </div>
        
        <!-- System Configuration Section (Phase 3) -->
        <div class="nav-section">
          <h3 class="section-title">System Configuration</h3>
          <ul class="nav-items">
            <li>
              <router-link to="/admin/integrations" class="nav-item">
                <i class="pi pi-cloud"></i>
                <span>Integrations</span>
                <Badge 
                  value="NEW" 
                  severity="info" 
                  class="nav-badge"
                />
              </router-link>
            </li>
          </ul>
        </div>
      </div>
    </nav>
    
    <!-- Admin Content Area -->
    <main class="admin-content">
      <slot />
    </main>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import Badge from 'primevue/badge'

// Mock pending claims count - in real app this would come from a store/API
const pendingClaimsCount = ref(12)
</script>

<style scoped>
.admin-layout {
  display: grid;
  grid-template-columns: var(--space-55) 1fr;
  gap: var(--space-21);
  min-height: calc(100vh - 60px - 48px); /* Account for navbar + admin context bar */
  padding: 0 var(--space-21) var(--space-21) var(--space-21); /* No top padding - space already reserved */
  max-width: var(--container-xl);
  margin: 0 auto;
}

/* Admin Sidebar */
.admin-sidebar {
  background: var(--surface-0);
  border-radius: var(--space-8);
  padding: var(--space-21);
  box-shadow: 0 var(--space-1) var(--space-8) rgba(0, 0, 0, 0.1);
  height: fit-content;
  position: static;
}

/* Navigation Sections */
.nav-sections {
  display: flex;
  flex-direction: column;
  gap: var(--space-21);
}

.nav-section {
  display: flex;
  flex-direction: column;
  gap: var(--space-8);
}

.section-title {
  font-size: var(--text-xs);
  font-weight: var(--font-semibold);
  color: var(--surface-500);
  text-transform: uppercase;
  letter-spacing: 0.05em;
  margin: 0;
}

.nav-items {
  list-style: none;
  margin: 0;
  padding: 0;
  display: flex;
  flex-direction: column;
  gap: var(--space-2);
}

.nav-item {
  display: flex;
  align-items: center;
  gap: var(--space-8);
  padding: var(--space-8) var(--space-13);
  text-decoration: none;
  color: var(--surface-700);
  font-size: var(--text-sm);
  font-weight: var(--font-medium);
  border-radius: var(--space-5);
  transition: all var(--transition-fast) ease;
  position: relative;
}

.nav-item:hover {
  background: var(--admin-primary-light);
  color: var(--admin-primary-dark);
}

.nav-item.router-link-active {
  background: var(--admin-primary);
  color: white;
}

.nav-item i {
  font-size: var(--text-sm);
  width: var(--space-13);
  text-align: center;
}

.nav-badge {
  margin-left: auto;
  font-size: var(--text-xs) !important;
  min-width: var(--space-13) !important;
  height: var(--space-13) !important;
}

/* Admin Content */
.admin-content {
  background: var(--surface-0);
  border-radius: var(--space-8);
  padding: 3rem var(--space-21) var(--space-21);
  box-shadow: 0 var(--space-1) var(--space-5) rgba(0, 0, 0, 0.1);
}

/* Responsive Design */
@media (max-width: var(--breakpoint-lg)) {
  .admin-layout {
    grid-template-columns: 1fr;
    gap: var(--space-13);
    padding: 0 var(--space-13) var(--space-13) var(--space-13); /* No top padding - space already reserved */
  }
  
  .admin-sidebar {
    order: 2;
    position: static;
    padding: var(--space-13);
  }
  
  .nav-sections {
    gap: var(--space-13);
  }
  
  .admin-content {
    padding: var(--space-13);
  }
}

@media (max-width: var(--breakpoint-sm)) {
  .admin-layout {
    padding: 0 var(--space-8) var(--space-8) var(--space-8); /* No top padding - space already reserved */
    gap: var(--space-8);
  }
  
  .admin-sidebar {
    padding: var(--space-8);
  }
  
  .admin-content {
    padding: var(--space-8);
  }
  
  .nav-item {
    padding: var(--space-5) var(--space-8);
    gap: var(--space-5);
  }
  
  .nav-item i {
    width: var(--space-8);
  }
}

/* Mobile Navigation - Horizontal scroll for better UX */
@media (max-width: calc(var(--breakpoint-sm) * 0.8)) {
  .nav-sections {
    flex-direction: row;
    overflow-x: auto;
    padding-bottom: var(--space-5);
    gap: var(--space-21);
  }
  
  .nav-section {
    flex-shrink: 0;
    min-width: 200px;
  }
  
  .nav-items {
    flex-direction: row;
    gap: var(--space-3);
  }
  
  .nav-item {
    white-space: nowrap;
    min-width: fit-content;
  }
}
</style>