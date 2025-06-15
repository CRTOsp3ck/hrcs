<template>
  <div class="page-header">
    <div class="page-header-content">
      <div class="page-header-text">
        <h1 class="page-title">{{ title }}</h1>
        <p v-if="subtitle" class="page-subtitle">{{ subtitle }}</p>
      </div>
      <div v-if="breadcrumbs && breadcrumbs.length" class="page-breadcrumbs">
        <nav class="breadcrumb-nav">
          <ol class="breadcrumb-list">
            <li v-for="(item, index) in breadcrumbs" :key="index" class="breadcrumb-item">
              <router-link 
                v-if="item.to && index < breadcrumbs.length - 1" 
                :to="item.to"
                class="breadcrumb-link"
              >
                {{ item.label }}
              </router-link>
              <span v-else class="breadcrumb-current">{{ item.label }}</span>
              <i v-if="index < breadcrumbs.length - 1" class="pi pi-chevron-right breadcrumb-separator"></i>
            </li>
          </ol>
        </nav>
      </div>
    </div>
    
    <div v-if="$slots.actions" class="page-header-actions">
      <slot name="actions" />
    </div>
  </div>
</template>

<script setup lang="ts">
interface BreadcrumbItem {
  label: string
  to?: string
}

interface Props {
  title: string
  subtitle?: string
  breadcrumbs?: BreadcrumbItem[]
}

defineProps<Props>()
</script>

<style scoped>
.page-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: var(--space-21);
  padding-bottom: var(--space-13);
  border-bottom: 1px solid var(--surface-200);
  gap: var(--space-21);
}

.page-header-content {
  display: flex;
  flex-direction: column;
  gap: var(--space-5);
  flex: 1;
}

.page-header-text {
  display: flex;
  flex-direction: column;
  gap: var(--space-3);
}

.page-title {
  font-size: var(--text-xl);
  font-weight: var(--font-bold);
  color: var(--surface-900);
  margin: 0;
  line-height: 1.2;
}

.page-subtitle {
  color: var(--surface-500);
  font-size: var(--text-base);
  margin: 0;
  line-height: 1.4;
}

.page-breadcrumbs {
  order: -1;
}

.breadcrumb-nav {
  margin-bottom: var(--space-5);
}

.breadcrumb-list {
  display: flex;
  align-items: center;
  list-style: none;
  margin: 0;
  padding: 0;
  gap: var(--space-3);
}

.breadcrumb-item {
  display: flex;
  align-items: center;
  gap: var(--space-3);
}

.breadcrumb-link {
  color: var(--surface-500);
  text-decoration: none;
  font-size: var(--text-sm);
  transition: color var(--transition-fast) ease;
}

.breadcrumb-link:hover {
  color: var(--primary-600);
}

.breadcrumb-current {
  color: var(--surface-900);
  font-size: var(--text-sm);
  font-weight: var(--font-medium);
}

.breadcrumb-separator {
  color: var(--surface-400);
  font-size: var(--text-xs);
}

.page-header-actions {
  display: flex;
  align-items: flex-start;
  gap: var(--space-8);
  flex-shrink: 0;
}

/* Responsive Design */
@media (max-width: var(--breakpoint-md)) {
  .page-header {
    flex-direction: column;
    align-items: stretch;
    gap: var(--space-13);
  }
  
  .page-header-actions {
    align-items: stretch;
    flex-direction: column;
  }
  
  .page-title {
    font-size: var(--text-lg);
  }
}

@media (max-width: var(--breakpoint-sm)) {
  .page-header {
    margin-bottom: var(--space-13);
    padding-bottom: var(--space-8);
  }
  
  .breadcrumb-list {
    flex-wrap: wrap;
  }
  
  .page-header-actions {
    gap: var(--space-5);
  }
}
</style>