<template>
  <div id="app" :class="{ 'admin-mode': showAdminContext }">
    <Toast position="top-right" />
    <ConfirmDialog />

    <!-- Navigation Header -->
    <Navbar v-if="authStore.isAuthenticated" />
    
    <!-- Admin Context Bar -->
    <div v-if="showAdminContext" class="context-bar">
      <div class="context-indicator">
        <i class="pi pi-shield"></i>
        <span>Administrator Mode</span>
      </div>
      <Button 
        @click="exitAdminMode" 
        size="small"
        severity="secondary"
        outlined
        class="exit-admin-btn"
      >
        Exit Admin Mode
      </Button>
    </div>

    <!-- Main Content Area -->
    <main :class="{ 'main-content': authStore.isAuthenticated }">
      <router-view v-slot="{ Component }">
        <transition name="fade" mode="out-in">
          <component :is="Component" />
        </transition>
      </router-view>
    </main>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import Navbar from '@/components/Navbar.vue'
import Button from 'primevue/button'

const route = useRoute()
const router = useRouter()
const authStore = useAuthStore()

// Admin context detection
const isAdminRoute = computed(() => route.path.startsWith('/admin'))
const showAdminContext = computed(() => 
  isAdminRoute.value && authStore.isAuthenticated && authStore.user?.role === 'admin'
)

// Exit admin mode handler
const exitAdminMode = () => {
  router.push('/dashboard')
}

onMounted(async () => {
  await authStore.init()
})
</script>

<style scoped>
#app {
  min-height: 100vh;
  width: 100%;
}

.main-content {
  min-height: calc(100vh - 60px); /* Account for navbar height */
  background: linear-gradient(135deg, var(--surface-50) 0%, var(--surface-100) 100%);
}

.admin-mode .main-content {
  min-height: calc(100vh - 60px - 48px); /* Account for navbar + context bar */
}

.exit-admin-btn {
  color: white !important;
  border-color: rgba(255, 255, 255, 0.3) !important;
  font-size: var(--text-xs) !important;
  padding: var(--space-2) var(--space-5) !important;
}

.exit-admin-btn:hover {
  background: rgba(255, 255, 255, 0.1) !important;
  border-color: rgba(255, 255, 255, 0.5) !important;
}

/* Enhanced page transitions */
.fade-enter-active,
.fade-leave-active {
  transition: opacity var(--transition-normal) ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}

/* Responsive adjustments */
@media (max-width: var(--breakpoint-sm)) {
  .context-bar {
    padding: var(--space-3) var(--space-8) !important;
    font-size: var(--text-xs) !important;
  }
  
  .context-bar .context-indicator {
    gap: var(--space-2) !important;
  }
  
  .exit-admin-btn {
    padding: var(--space-1) var(--space-3) !important;
  }
}
</style>
