<template>
  <div id="app">
    <Toast position="top-right" />
    <ConfirmDialog />

    <Navbar v-if="authStore.isAuthenticated" />

    <main :class="{ 'main-content': authStore.isAuthenticated }">
      <router-view v-slot="{ Component }">
        <transition name="fade" mode="out-in">
          <!-- <section> -->
            <component :is="Component" />
          <!-- </section> -->
        </transition>
      </router-view>
    </main>
  </div>
</template>

<script setup lang="ts">
import { onMounted } from 'vue'
import { useAuthStore } from '@/stores/auth'
import Navbar from '@/components/Navbar.vue'

const authStore = useAuthStore()

onMounted(async () => {
  await authStore.init()
})
</script>

<style scoped>
.main-content {
  flex: 1;
  margin-top: 3rem;
  min-height: calc(100vh - 72px);
  background: var(--surface-50);
}

.fade-enter-active,
.fade-leave-active {
  transition: all var(--transition-slow) ease;
}

.fade-enter-from {
  opacity: 0;
  /* transform: translateY(20px); */
}

.fade-leave-to {
  opacity: 0;
  /* transform: translateY(-20px); */
}

.fade-enter-to,
.fade-leave-from {
  opacity: 1;
  transform: translateY(0);
}
</style>
