<template>
  <div id="app">
    <Toast position="top-right" />
    <ConfirmDialog />

    <Navbar v-if="authStore.isAuthenticated" />

    <main :class="{ 'main-content': authStore.isAuthenticated }">
      <router-view v-slot="{ Component }">
        <transition name="fade" mode="out-in">
          <section>
            <component :is="Component" />
          </section>
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
  padding: 6rem 0rem;
  min-height: 100vh;
}

.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.3s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>
