<template>
  <Card class="stat-card" :interactive="interactive">
    <div class="stat-content">
      <div class="stat-icon" :class="`stat-icon--${color}`">
        <i :class="`pi pi-${icon}`"></i>
      </div>
      
      <div class="stat-details">
        <div class="stat-value">{{ value }}</div>
        <div class="stat-title">{{ title }}</div>
        
        <div v-if="change !== undefined" class="stat-change" :class="{
          'stat-change--positive': change > 0,
          'stat-change--negative': change < 0,
          'stat-change--neutral': change === 0
        }">
          <i :class="{
            'pi pi-arrow-up': change > 0,
            'pi pi-arrow-down': change < 0,
            'pi pi-minus': change === 0
          }"></i>
          <span>{{ formatChange(change) }}%</span>
          <span class="stat-change-period">vs last month</span>
        </div>
      </div>
    </div>
  </Card>
</template>

<script setup lang="ts">
import Card from './Card.vue'

interface Props {
  title: string
  value: string | number
  icon: string
  color?: 'primary' | 'success' | 'warning' | 'danger' | 'info'
  change?: number
  interactive?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  color: 'primary',
  interactive: false
})

const formatChange = (change: number): string => {
  return Math.abs(change).toFixed(1)
}
</script>

<style scoped>
.stat-card {
  position: relative;
  overflow: hidden;
}

.stat-card::before {
  content: '';
  position: absolute;
  top: -50%;
  right: -50%;
  width: 100%;
  height: 100%;
  background: radial-gradient(circle, rgba(59, 130, 246, 0.05) 0%, transparent 70%);
  pointer-events: none;
}

.stat-content {
  display: flex;
  align-items: flex-start;
  gap: var(--space-13);
}

.stat-icon {
  display: flex;
  align-items: center;
  justify-content: center;
  width: var(--space-21);
  height: var(--space-21);
  border-radius: var(--space-5);
  flex-shrink: 0;
}

.stat-icon--primary {
  background: var(--primary-100);
  color: var(--primary-600);
}

.stat-icon--success {
  background: #dcfce7;
  color: #16a34a;
}

.stat-icon--warning {
  background: #fef3c7;
  color: #d97706;
}

.stat-icon--danger {
  background: #fee2e2;
  color: #dc2626;
}

.stat-icon--info {
  background: #dbeafe;
  color: #2563eb;
}

.stat-icon i {
  font-size: var(--text-lg);
}

.stat-details {
  display: flex;
  flex-direction: column;
  gap: var(--space-2);
  flex: 1;
  min-width: 0;
}

.stat-value {
  font-size: var(--text-xl);
  font-weight: var(--font-bold);
  color: var(--surface-900);
  line-height: 1;
}

.stat-title {
  font-size: var(--text-sm);
  color: var(--surface-600);
  font-weight: var(--font-medium);
  line-height: 1.2;
}

.stat-change {
  display: flex;
  align-items: center;
  gap: var(--space-2);
  font-size: var(--text-xs);
  font-weight: var(--font-medium);
  margin-top: var(--space-3);
}

.stat-change--positive {
  color: #16a34a;
}

.stat-change--negative {
  color: #dc2626;
}

.stat-change--neutral {
  color: var(--surface-500);
}

.stat-change i {
  font-size: var(--text-xs);
}

.stat-change-period {
  color: var(--surface-500);
  font-weight: var(--font-normal);
}

/* Responsive Design */
@media (max-width: var(--breakpoint-sm)) {
  .stat-content {
    gap: var(--space-8);
  }
  
  .stat-icon {
    width: var(--space-13);
    height: var(--space-13);
  }
  
  .stat-icon i {
    font-size: var(--text-base);
  }
  
  .stat-value {
    font-size: var(--text-lg);
  }
}
</style>