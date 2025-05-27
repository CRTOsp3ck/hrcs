<template>
  <div class="auth-container">
    <Card class="auth-card">
      <template #header>
        <div class="auth-header">
          <i class="pi pi-user-plus auth-logo"></i>
          <h1 class="auth-title">Create Account</h1>
          <p class="auth-subtitle">Join HR Claims Management System</p>
        </div>
      </template>
      
      <template #content>
        <form @submit.prevent="handleRegister" class="auth-form">
          <div class="form-field">
            <label for="name" class="form-label">Full Name</label>
            <InputText 
              id="name"
              v-model="form.name"
              placeholder="John Doe"
              :invalid="!!errors.name"
              required
              class="w-full"
            />
            <small v-if="errors.name" class="p-error">{{ errors.name }}</small>
          </div>
          
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
              placeholder="Create a strong password"
              :invalid="!!errors.password"
              required
              toggleMask
              class="w-full"
              :feedback="true"
              :strongRegex="'^(?=.*[a-z])(?=.*[A-Z])(?=.*[0-9])(?=.*[!@#\$%]).{8,}'"
              weakLabel="Weak password"
              mediumLabel="Medium password" 
              strongLabel="Strong password"
            />
            <small v-if="errors.password" class="p-error">{{ errors.password }}</small>
          </div>
          
          <div class="form-field">
            <label for="confirmPassword" class="form-label">Confirm Password</label>
            <Password 
              id="confirmPassword"
              v-model="form.confirmPassword"
              placeholder="Re-enter your password"
              :invalid="!!errors.confirmPassword"
              required
              :feedback="false"
              toggleMask
              class="w-full"
            />
            <small v-if="errors.confirmPassword" class="p-error">{{ errors.confirmPassword }}</small>
          </div>
          
          <Message v-if="error" severity="error" :closable="false" class="mb-3">
            {{ error }}
          </Message>
          
          <Button 
            type="submit"
            label="Create Account"
            icon="pi pi-user-plus"
            :loading="authStore.loading"
            class="w-full"
            size="large"
          />
        </form>
        
        <Divider align="center" class="my-4">
          <span class="text-sm text-color-secondary">Already have an account?</span>
        </Divider>
        
        <router-link to="/login" class="no-underline">
          <Button 
            label="Sign In"
            icon="pi pi-sign-in"
            severity="secondary"
            outlined
            class="w-full"
            size="large"
          />
        </router-link>
      </template>
    </Card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { useToast } from 'primevue/usetoast'

const router = useRouter()
const authStore = useAuthStore()
const toast = useToast()

const form = reactive({
  name: '',
  email: '',
  password: '',
  confirmPassword: ''
})

const errors = reactive({
  name: '',
  email: '',
  password: '',
  confirmPassword: ''
})

const error = ref('')

const validateForm = () => {
  let isValid = true
  Object.keys(errors).forEach(key => errors[key as keyof typeof errors] = '')
  
  if (!form.name) {
    errors.name = 'Name is required'
    isValid = false
  } else if (form.name.length < 2) {
    errors.name = 'Name must be at least 2 characters'
    isValid = false
  }
  
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
  } else if (form.password.length < 8) {
    errors.password = 'Password must be at least 8 characters'
    isValid = false
  }
  
  if (!form.confirmPassword) {
    errors.confirmPassword = 'Please confirm your password'
    isValid = false
  } else if (form.password !== form.confirmPassword) {
    errors.confirmPassword = 'Passwords do not match'
    isValid = false
  }
  
  return isValid
}

const handleRegister = async () => {
  if (!validateForm()) return
  
  error.value = ''
  
  try {
    await authStore.register({
      name: form.name,
      email: form.email,
      password: form.password
    })
    toast.add({
      severity: 'success',
      summary: 'Welcome!',
      detail: 'Account created successfully',
      life: 3000
    })
    router.push('/dashboard')
  } catch (err: any) {
    error.value = err.response?.data?.message || 'Registration failed'
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
  gap: 1.25rem;
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

:deep(.p-password-input) {
  width: 100%;
}
</style>