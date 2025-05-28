<template>
  <Menubar :model="menuItems" class="navbar-fixed">
    <template #start>
      <router-link to="/dashboard" class="navbar-brand">
        <!-- <i class="pi pi-building mr-2"></i> -->
        <span class="font-bold">XR-Claimatix</span>
        <div class="mr-8"></div>
      </router-link>
    </template>

    <template #end>
      <div class="flex items-center gap-3">
        <Avatar
          :label="userInitials"
          :style="{ backgroundColor: '#2563eb', color: '#ffffff' }"
          shape="circle"
        />
        <div class="hide-mobile">
          <div class="font-semibold">{{ authStore.user?.name }}</div>
          <div class="text-sm text-color-secondary">{{ authStore.user?.email }}</div>
        </div>
        <Button
          icon="pi pi-sign-out"
          @click="handleLogout"
          severity="secondary"
          text
          rounded
          aria-label="Logout"
        />
      </div>
    </template>
  </Menubar>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { useConfirm } from 'primevue/useconfirm'

const router = useRouter()
const authStore = useAuthStore()
const confirm = useConfirm()

const userInitials = computed(() => {
  const name = authStore.user?.name || ''
  return name.split(' ').map(n => n[0]).join('').toUpperCase() || 'U'
})

const navigateToRoute = (routePath: string) => {
  router.push(routePath)
}

const menuItems = computed(() => [
  {
    label: 'Dashboard',
    icon: 'pi pi-home',
    command: () => navigateToRoute('/dashboard')
  },
  {
    label: 'My Claims',
    icon: 'pi pi-file',
    // command: () => navigateToRoute('/claims'),
    items: [
      {
        label: 'View All Claims',
        icon: 'pi pi-list',
        command: () => navigateToRoute('/claims')
      },
      {
        label: 'New Claim',
        icon: 'pi pi-plus',
        command: () => navigateToRoute('/claims/new')
      }
    ]
  },
  ...(authStore.isAdmin ? [{
    label: 'Admin',
    icon: 'pi pi-cog',
    items: [
      {
        label: 'Dashboard',
        icon: 'pi pi-chart-bar',
        command: () => navigateToRoute('/admin/dashboard')
      },
      {
        separator: true
      },
      {
        label: 'Users',
        icon: 'pi pi-users',
        command: () => navigateToRoute('/admin/users')
      },
      {
        label: 'User Groups',
        icon: 'pi pi-sitemap',
        command: () => navigateToRoute('/admin/groups')
      },
      {
        separator: true
      },
      {
        label: 'Claim Types',
        icon: 'pi pi-tags',
        command: () => navigateToRoute('/admin/claim-types')
      },
      {
        label: 'Approval Levels',
        icon: 'pi pi-shield',
        command: () => navigateToRoute('/admin/approval-levels')
      },
      {
        separator: true
      },
      {
        label: 'All Claims',
        icon: 'pi pi-file-o',
        command: () => navigateToRoute('/admin/claims')
      }
    ]
  }] : [])
])

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
</script>

<style scoped>
.navbar-fixed {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  z-index: var(--z-fixed);
  background: rgba(255, 255, 255, 0.9);
  backdrop-filter: blur(20px);
  border-bottom: 1px solid rgba(0, 0, 0, 0.05);
  transition: all var(--transition-base) ease;
}

.navbar-brand {
  display: flex;
  align-items: center;
  text-decoration: none;
  color: var(--primary-600);
  font-size: var(--text-xl);
  font-weight: var(--font-bold);
  letter-spacing: -0.025em;
  transition: color var(--transition-fast) ease;
}

.navbar-brand:hover {
  color: var(--primary-700);
}

.navbar-brand span {
  background: linear-gradient(135deg, var(--primary-600) 0%, var(--primary-700) 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

:deep(.p-menubar) {
  background: transparent;
  border: none;
  padding: var(--space-4) var(--space-6);
}

:deep(.p-menubar .p-menuitem-link) {
  font-size: var(--text-sm);
  font-weight: var(--font-medium);
  color: var(--surface-600);
  padding: var(--space-2) var(--space-3);
  border-radius: var(--radius-lg);
  transition: all var(--transition-fast) ease;
}

:deep(.p-menubar .p-menuitem-link:hover) {
  background: var(--surface-100);
  color: var(--surface-900);
}

:deep(.p-menubar .p-menuitem-link .p-menuitem-icon) {
  margin-right: var(--space-2);
}

:deep(.p-submenu-list) {
  margin-top: var(--space-2);
  border: 1px solid var(--surface-100);
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(10px);
}

:deep(.p-avatar) {
  transition: transform var(--transition-fast) ease;
}

:deep(.p-avatar:hover) {
  transform: scale(1.05);
}

.user-info {
  animation: slideIn var(--transition-slow) ease;
}

@media (max-width: 768px) {
  :deep(.p-menubar) {
    padding: var(--space-3) var(--space-4);
  }
}
</style>
