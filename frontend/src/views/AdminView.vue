<template>
  <div class="admin-layout">
    <div class="admin-sidebar">
      <div class="sidebar-header">
        <i class="pi pi-shield"></i>
        <span>Admin Panel</span>
      </div>
      
      <nav class="sidebar-nav">
        <router-link 
          v-for="item in menuItems" 
          :key="item.route"
          :to="item.route"
          class="nav-item"
          :class="{ active: $route.path === item.route }"
        >
          <i :class="item.icon"></i>
          <span>{{ item.label }}</span>
        </router-link>
      </nav>
    </div>
    
    <div class="admin-content">
      <router-view />
    </div>
  </div>
</template>

<script setup lang="ts">
const menuItems = [
  {
    label: 'Dashboard',
    icon: 'pi pi-chart-bar',
    route: '/admin/dashboard'
  },
  {
    label: 'Users',
    icon: 'pi pi-users',
    route: '/admin/users'
  },
  {
    label: 'User Groups',
    icon: 'pi pi-sitemap',
    route: '/admin/groups'
  },
  {
    label: 'Claim Types',
    icon: 'pi pi-tags',
    route: '/admin/claim-types'
  },
  {
    label: 'Approval Levels',
    icon: 'pi pi-shield',
    route: '/admin/approval-levels'
  },
  {
    label: 'All Claims',
    icon: 'pi pi-file-o',
    route: '/admin/claims'
  }
]
</script>

<style scoped>
.admin-layout {
  display: flex;
  min-height: calc(100vh - 60px);
}

.admin-sidebar {
  width: 250px;
  background: var(--surface-0);
  border-right: 1px solid var(--surface-200);
  flex-shrink: 0;
}

.sidebar-header {
  padding: 1.5rem;
  display: flex;
  align-items: center;
  gap: 0.75rem;
  font-size: 1.25rem;
  font-weight: 600;
  color: var(--surface-900);
  border-bottom: 1px solid var(--surface-200);
}

.sidebar-header i {
  color: var(--primary-500);
}

.sidebar-nav {
  padding: 1rem 0;
}

.nav-item {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  padding: 0.75rem 1.5rem;
  color: var(--surface-600);
  text-decoration: none;
  transition: all 0.2s;
}

.nav-item:hover {
  background: var(--surface-50);
  color: var(--surface-900);
}

.nav-item.active {
  background: var(--primary-50);
  color: var(--primary-600);
  font-weight: 500;
  border-left: 3px solid var(--primary-600);
}

.admin-content {
  flex: 1;
  background: var(--surface-50);
  overflow-y: auto;
}

@media (max-width: 1024px) {
  .admin-sidebar {
    width: 200px;
  }
}

@media (max-width: 768px) {
  .admin-layout {
    flex-direction: column;
  }
  
  .admin-sidebar {
    width: 100%;
    border-right: none;
    border-bottom: 1px solid var(--surface-200);
  }
  
  .sidebar-nav {
    display: flex;
    overflow-x: auto;
    padding: 0;
  }
  
  .nav-item {
    white-space: nowrap;
    border-left: none !important;
    border-bottom: 3px solid transparent;
  }
  
  .nav-item.active {
    border-bottom-color: var(--primary-600);
  }
}
</style>