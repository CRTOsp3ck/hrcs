<template>
  <header class="app-header">
    <!-- Primary Navigation -->
    <nav class="navbar">
      <!-- Brand Section -->
      <div class="navbar-brand">
        <router-link to="/dashboard" class="brand-link">
          <i class="pi pi-building brand-icon"></i>
          <span class="brand-text">Claimatic</span>
        </router-link>
      </div>

      <!-- Navigation Menu (Hidden on Admin Routes) -->
      <ul v-if="!isAdminRoute" class="nav-menu">
        <li class="nav-item">
          <router-link to="/dashboard" class="nav-link">
            <i class="pi pi-home"></i>
            <span class="hide-mobile">Dashboard</span>
          </router-link>
        </li>
        <li class="nav-item nav-dropdown">
          <div class="nav-link dropdown-trigger" @click="toggleClaimsDropdown">
            <i class="pi pi-file"></i>
            <span class="hide-mobile">Claims</span>
            <i class="pi pi-chevron-down dropdown-icon" :class="{ 'rotated': showClaimsDropdown }"></i>
          </div>
          <ul v-if="showClaimsDropdown" class="dropdown-menu">
            <li>
              <router-link to="/claims" class="dropdown-link" @click="closeDropdowns">
                <i class="pi pi-list"></i>
                View All
              </router-link>
            </li>
            <li>
              <router-link to="/claims/new" class="dropdown-link" @click="closeDropdowns">
                <i class="pi pi-plus"></i>
                New Claim
              </router-link>
            </li>
          </ul>
        </li>
      </ul>

      <!-- Right Section: Admin Button + User Profile -->
      <div class="navbar-right">
        <!-- Admin Button -->
        <Button
          v-if="authStore.isAdmin && !isAdminRoute"
          @click="goToAdmin"
          size="small"
          class="admin-button"
        >
          <i class="pi pi-shield"></i>
          <span class="hide-mobile">Admin Mode</span>
        </Button>

        <!-- User Profile Section -->
        <div class="user-profile">
        <div class="user-avatar">
          <Avatar
            :label="userInitials"
            :style="{ backgroundColor: 'var(--primary-600)', color: '#ffffff' }"
            shape="circle"
            size="normal"
          />
        </div>
        <div class="user-info hide-mobile">
          <span class="user-name">{{ authStore.user?.name }}</span>
          <span class="user-role">{{ authStore.user?.role }}</span>
        </div>
        <div class="user-menu">
          <Button
            icon="pi pi-sign-out"
            @click="handleLogout"
            severity="secondary"
            text
            rounded
            aria-label="Logout"
            class="logout-btn"
          />
        </div>
      </div>
      </div>
    </nav>
  </header>
</template>

<script setup lang="ts">
import { computed, ref, onMounted, onUnmounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { useConfirm } from 'primevue/useconfirm'
import Avatar from 'primevue/avatar'
import Button from 'primevue/button'

const route = useRoute()
const router = useRouter()
const authStore = useAuthStore()
const confirm = useConfirm()

// Dropdown state
const showClaimsDropdown = ref(false)

// Admin route detection
const isAdminRoute = computed(() => route.path.startsWith('/admin'))

const userInitials = computed(() => {
  const name = authStore.user?.name || ''
  return name.split(' ').map(n => n[0]).join('').toUpperCase() || 'U'
})

// Dropdown handlers
const toggleClaimsDropdown = () => {
  showClaimsDropdown.value = !showClaimsDropdown.value
}

const closeDropdowns = () => {
  showClaimsDropdown.value = false
}

const goToAdmin = () => {
  router.push('/admin/dashboard')
}

// Close dropdowns when clicking outside
const handleClickOutside = (event: Event) => {
  const target = event.target as Element
  if (!target.closest('.nav-dropdown')) {
    closeDropdowns()
  }
}

const handleLogout = () => {
  confirm.require({
    message: 'Are you sure you want to logout?',
    header: 'Logout Confirmation',
    icon: 'pi pi-sign-out',
    acceptClass: 'p-button-danger',
    accept: () => {
      authStore.logout()
      router.push('/login')
    }
  })
}

onMounted(() => {
  document.addEventListener('click', handleClickOutside)
})

onUnmounted(() => {
  document.removeEventListener('click', handleClickOutside)
})
</script>

<style scoped>
.app-header {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  z-index: 1000;
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(10px);
  border-bottom: 1px solid var(--surface-100);
  box-shadow: 0 var(--space-1) var(--space-5) rgba(0, 0, 0, 0.1);
}

.navbar {
  display: flex;
  align-items: center;
  padding: var(--space-8) var(--space-21);
  max-width: var(--container-xl);
  margin: 0 auto;
  height: 60px;
  position: relative;
}

/* Brand Section */
.navbar-brand {
  display: flex;
  align-items: center;
  flex: 0 0 auto;
}

.brand-link {
  display: flex;
  align-items: center;
  gap: var(--space-5);
  text-decoration: none;
  color: var(--primary-600);
  font-size: var(--text-lg);
  font-weight: var(--font-bold);
  transition: color var(--transition-fast) ease;
}

.brand-link:hover {
  color: var(--primary-700);
}

.brand-icon {
  font-size: var(--text-lg);
}

.brand-text {
  font-size: var(--text-lg);
}

/* Navigation Menu */
.nav-menu {
  display: flex;
  align-items: center;
  list-style: none;
  margin: 0;
  padding: 0;
  gap: var(--space-8);
  position: absolute;
  left: 50%;
  transform: translateX(-50%);
}

.nav-item {
  position: relative;
}

.nav-link {
  display: flex;
  align-items: center;
  gap: var(--space-3);
  padding: var(--space-5) var(--space-8);
  text-decoration: none;
  color: var(--surface-700);
  font-size: var(--text-sm);
  font-weight: var(--font-medium);
  border-radius: var(--space-5);
  transition: all var(--transition-fast) ease;
  cursor: pointer;
}

.nav-link:hover {
  background: var(--surface-100);
  color: var(--primary-600);
}

.nav-link.router-link-active {
  background: var(--primary-100);
  color: var(--primary-700);
}

/* Dropdown */
.nav-dropdown {
  position: relative;
}

.dropdown-trigger {
  position: relative;
}

.dropdown-icon {
  font-size: var(--text-xs);
  transition: transform var(--transition-fast) ease;
}

.dropdown-icon.rotated {
  transform: rotate(180deg);
}

.dropdown-menu {
  position: absolute;
  top: 100%;
  left: 0;
  background: var(--surface-0);
  border: 1px solid var(--surface-200);
  border-radius: var(--space-5);
  box-shadow: 0 var(--space-5) var(--space-21) rgba(0, 0, 0, 0.15);
  min-width: 180px;
  padding: var(--space-5) 0;
  margin-top: var(--space-2);
  list-style: none;
  z-index: 1001;
}

.dropdown-link {
  display: flex;
  align-items: center;
  gap: var(--space-5);
  padding: var(--space-5) var(--space-13);
  text-decoration: none;
  color: var(--surface-700);
  font-size: var(--text-sm);
  transition: background var(--transition-fast) ease;
}

.dropdown-link:hover {
  background: var(--surface-50);
  color: var(--primary-600);
}

.dropdown-separator {
  height: 1px;
  background: var(--surface-200);
  margin: var(--space-3) 0;
}

/* Admin Button Special Styling */
.admin-button {
  gap: var(--space-2) !important;
  font-weight: var(--font-medium) !important;
  background: var(--primary-600) !important;
  color: white !important;
  border: 1px solid var(--primary-600) !important;
  border-radius: var(--space-5) !important;
  padding: var(--space-3) var(--space-8) !important;
  transition: all var(--transition-fast) ease !important;
}

.admin-button:hover {
  background: var(--primary-700) !important;
  border-color: var(--primary-700) !important;
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(59, 130, 246, 0.3) !important;
}

/* Right Section */
.navbar-right {
  display: flex;
  align-items: center;
  gap: var(--space-21);
  margin-left: auto;
}

/* User Profile */
.user-profile {
  display: flex;
  align-items: center;
  gap: var(--space-8);
}

.user-avatar {
  display: flex;
  align-items: center;
}

.user-info {
  display: flex;
  flex-direction: column;
  gap: var(--space-1);
}

.user-name {
  font-size: var(--text-sm);
  font-weight: var(--font-semibold);
  color: var(--surface-900);
  line-height: 1;
}

.user-role {
  font-size: var(--text-xs);
  color: var(--surface-500);
  text-transform: capitalize;
  line-height: 1;
}

.logout-btn {
  color: var(--surface-600) !important;
  padding: var(--space-5) !important;
}

.logout-btn:hover {
  color: var(--primary-600) !important;
  background: var(--surface-100) !important;
}

/* Responsive Design */
@media (max-width: var(--breakpoint-sm)) {
  .navbar {
    padding: var(--space-5) var(--space-8);
  }

  .nav-menu {
    gap: var(--space-3);
  }

  .nav-link {
    padding: var(--space-3) var(--space-5);
  }

  .brand-text {
    font-size: var(--text-base);
  }

  .user-profile {
    gap: var(--space-5);
  }

  .dropdown-menu {
    right: 0;
    left: auto;
    min-width: 160px;
  }
}

@media (max-width: calc(var(--breakpoint-sm) * 0.8)) {
  .nav-menu {
    display: none;
  }

  .navbar {
    justify-content: space-between;
  }
}
</style>
