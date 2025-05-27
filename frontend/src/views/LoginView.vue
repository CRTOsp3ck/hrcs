<template>
  <div class="auth-container">
    <Card class="auth-card">
      <template #header>
        <div class="auth-header">
          <i class="pi pi-building auth-logo"></i>
          <h1 class="auth-title">HR Claims Management</h1>
          <p class="auth-subtitle">Sign in to manage your claims</p>
        </div>
      </template>
      
      <template #content>
        <form @submit.prevent="handleLogin" class="auth-form">
          <div class="form-field">
            <label for="email" class="form-label">Email Address</label>
            <InputText 
              id="email"
              v-model="form.email"
              type="email"
              placeholder="john.doe@company.com"
              :invalid="!!errors.email"
              required
              class="w-full"
            />
            <small v-if="errors.email" class="p-error">{{ errors.email }}</small>
          </div>
          
          <div class="form-field">
            <label for="password" class="form-label">Password</label>
            <Password 
              id="password"
              v-model="form.password"
              placeholder="Enter your password"
              :feedback="false"
              :invalid="!!errors.password"
              required
              toggleMask
              class="w-full"
            />
            <small v-if="errors.password" class="p-error">{{ errors.password }}</small>
          </div>
          
          <Message v-if="error" severity="error" :closable="false" class="mb-3">
            {{ error }}
          </Message>
          
          <Button 
            type="submit"
            label="Sign In"
            icon="pi pi-sign-in"
            :loading="authStore.loading"
            class="w-full"
            size="large"
          />
        </form>
        
        <Divider align="center" class="my-4">
          <span class="text-sm text-color-secondary">New to the platform?</span>
        </Divider>
        
        <router-link to="/register" class="no-underline">
          <Button 
            label="Create Account"
            icon="pi pi-user-plus"
            severity="secondary"
            outlined
            class="w-full"
            size="large"
          />
        </router-link>
        
        <div class="demo-credentials">
          <Fieldset legend="Demo Credentials" :toggleable="true" :collapsed="true">
            <div class="demo-grid">
              <div>
                <h4>Admin User</h4>
                <p class="mb-1"><strong>Email:</strong> admin@hrcs.com</p>
                <p><strong>Password:</strong> password123</p>
              </div>
              <div>
                <h4>Normal User</h4>
                <p class="mb-1"><strong>Email:</strong> john.doe@hrcs.com</p>
                <p><strong>Password:</strong> password123</p>
              </div>
            </div>
          </Fieldset>
        </div>
      </template>
    </Card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { useToast } from 'primevue/usetoast'
import type { LoginRequest } from '@/types'

const router = useRouter()
const authStore = useAuthStore()
const toast = useToast()

const form = reactive<LoginRequest>({
  email: '',
  password: ''
})

const errors = reactive({
  email: '',
  password: ''
})

const error = ref('')

const validateForm = () => {
  let isValid = true
  errors.email = ''
  errors.password = ''
  
  if (!form.email) {
    errors.email = 'Email is required'
    isValid = false
  } else if (!/\S+@\S+\.\S+/.test(form.email)) {
    errors.email = 'Invalid email format'
    isValid = false
  }
  
  if (!form.password) {
    errors.password = 'Password is required'
    isValid = false
  }
  
  return isValid
}

const handleLogin = async () => {
  if (!validateForm()) return
  
  error.value = ''
  
  try {
    await authStore.login(form)
    toast.add({
      severity: 'success',
      summary: 'Welcome back!',
      detail: 'Login successful',
      life: 3000
    })
    router.push('/dashboard')
  } catch (err: any) {
    error.value = err.response?.data?.message || 'Invalid credentials'
  }
}
</script>

<style scoped>
.auth-container {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 2rem;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

.auth-card {
  width: 100%;
  max-width: 480px;
  box-shadow: 0 20px 25px -5px rgba(0, 0, 0, 0.1), 0 10px 10px -5px rgba(0, 0, 0, 0.04);
}

.auth-header {
  text-align: center;
  padding: 2rem 2rem 0;
}

.auth-logo {
  font-size: 4rem;
  color: var(--primary-500);
  margin-bottom: 1rem;
}

.auth-title {
  font-size: 1.75rem;
  font-weight: 700;
  color: var(--surface-900);
  margin: 0 0 0.5rem;
}

.auth-subtitle {
  color: var(--surface-600);
  margin: 0;
}

.auth-form {
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
}

.form-field {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.w-full {
  width: 100%;
}

.no-underline {
  text-decoration: none;
}

.demo-credentials {
  margin-top: 2rem;
}

.demo-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 2rem;
}

.demo-grid h4 {
  margin: 0 0 0.5rem;
  color: var(--primary-600);
}

.demo-grid p {
  margin: 0;
  font-size: 0.875rem;
}

:deep(.p-fieldset) {
  background: var(--surface-50);
}

:deep(.p-password-input) {
  width: 100%;
}

@media (max-width: 640px) {
  .demo-grid {
    grid-template-columns: 1fr;
    gap: 1rem;
  }
}
</style>