<template>
  <div class="form-field" :class="{ 'form-field--error': !!error }">
    <label v-if="label" class="form-label" :class="{ required: required }" :for="fieldId">
      {{ label }}
    </label>
    
    <div class="form-input-wrapper">
      <slot :fieldId="fieldId" />
    </div>
    
    <small v-if="error" class="form-error">
      {{ error }}
    </small>
    
    <small v-else-if="helpText" class="form-help">
      {{ helpText }}
    </small>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'

interface Props {
  label?: string
  required?: boolean
  error?: string
  helpText?: string
  id?: string
}

const props = withDefaults(defineProps<Props>(), {
  required: false
})

const fieldId = computed(() => props.id || `field-${Math.random().toString(36).substr(2, 9)}`)
</script>

<style scoped>
.form-field {
  display: flex;
  flex-direction: column;
  gap: var(--space-3);
}

.form-label {
  font-weight: var(--font-medium);
  color: var(--surface-700);
  font-size: var(--text-sm);
  line-height: 1.4;
}

.form-label.required::after {
  content: " *";
  color: #ef4444;
}

.form-input-wrapper {
  position: relative;
}

.form-error {
  color: #ef4444;
  font-size: var(--text-xs);
  font-weight: var(--font-medium);
  line-height: 1.4;
}

.form-help {
  color: var(--surface-500);
  font-size: var(--text-xs);
  line-height: 1.4;
}

.form-field--error .form-label {
  color: #ef4444;
}

/* Focus styles for inputs */
.form-field:focus-within .form-label {
  color: var(--primary-600);
}

.form-field--error:focus-within .form-label {
  color: #ef4444;
}
</style>