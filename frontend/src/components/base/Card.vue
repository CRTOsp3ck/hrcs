<template>
  <div class="card" :class="[`card--${variant}`, { 'card--interactive': interactive }]">
    <header v-if="$slots.header" class="card-header">
      <slot name="header" />
    </header>
    
    <div class="card-content">
      <slot />
    </div>
    
    <footer v-if="$slots.footer" class="card-footer">
      <slot name="footer" />
    </footer>
  </div>
</template>

<script setup lang="ts">
interface Props {
  variant?: 'default' | 'elevated' | 'outlined'
  interactive?: boolean
}

withDefaults(defineProps<Props>(), {
  variant: 'default',
  interactive: false
})
</script>

<style scoped>
.card {
  background: var(--surface-0);
  border-radius: var(--space-8);
  overflow: hidden;
  transition: all var(--transition-normal) ease;
}

.card--default {
  box-shadow: 0 var(--space-1) var(--space-5) rgba(0, 0, 0, 0.1);
}

.card--elevated {
  box-shadow: 0 var(--space-5) var(--space-21) rgba(0, 0, 0, 0.15);
}

.card--outlined {
  border: 1px solid var(--surface-200);
  box-shadow: none;
}

.card--interactive {
  cursor: pointer;
}

.card--interactive:hover {
  box-shadow: 0 var(--space-3) var(--space-13) rgba(0, 0, 0, 0.15);
  transform: translateY(-var(--space-1));
}

.card-header {
  padding: var(--space-13) var(--space-21) var(--space-8);
  border-bottom: 1px solid var(--surface-100);
  background: var(--surface-0);
}

.card-content {
  padding: var(--space-21);
}

.card-footer {
  padding: var(--space-8) var(--space-21) var(--space-13);
  background: var(--surface-50);
  border-top: 1px solid var(--surface-100);
}

/* Responsive adjustments */
@media (max-width: var(--breakpoint-sm)) {
  .card-header {
    padding: var(--space-8) var(--space-13) var(--space-5);
  }
  
  .card-content {
    padding: var(--space-13);
  }
  
  .card-footer {
    padding: var(--space-5) var(--space-13) var(--space-8);
  }
}
</style>