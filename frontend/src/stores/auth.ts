import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { User, LoginRequest, RegisterRequest } from '@/types'
import { authApi } from '@/api'

export const useAuthStore = defineStore('auth', () => {
  const user = ref<User | null>(null)
  const token = ref<string | null>(null)
  const loading = ref(false)
  const error = ref<string | null>(null)

  const isAuthenticated = computed(() => !!token.value && !!user.value)
  const isAdmin = computed(() => user.value?.role === 'admin')

  function isTokenExpired() {
    if (!token.value) return true
    
    try {
      // JWT tokens have 3 parts separated by dots
      const tokenParts = token.value.split('.')
      if (tokenParts.length !== 3) return true
      
      // Decode the payload (second part)
      const payload = JSON.parse(atob(tokenParts[1]))
      
      // Check if token has expired (exp is in seconds, Date.now() is in milliseconds)
      if (payload.exp && payload.exp * 1000 < Date.now()) {
        return true
      }
      
      return false
    } catch (error) {
      console.error('Error checking token expiration:', error)
      return true
    }
  }

  async function init() {
    const storedToken = localStorage.getItem('token')
    const storedUser = localStorage.getItem('user')

    if (storedToken && storedUser) {
      token.value = storedToken
      user.value = JSON.parse(storedUser)

      // Verify token is still valid
      try {
        const response = await authApi.getProfile()
        if (response.data.data) {
          user.value = response.data.data
          localStorage.setItem('user', JSON.stringify(response.data.data))
        }
      } catch (error) {
        logout()
      }
    }
  }

  async function login(credentials: LoginRequest) {
    loading.value = true
    error.value = null

    try {
      const response = await authApi.login(credentials)
      if (response.data.data) {
        token.value = response.data.data.token
        user.value = response.data.data.user

        localStorage.setItem('token', response.data.data.token)
        localStorage.setItem('user', JSON.stringify(response.data.data.user))
      }
    } catch (err: any) {
      error.value = err.response?.data?.message || 'Login failed'
      throw err
    } finally {
      loading.value = false
    }
  }

  async function register(data: RegisterRequest) {
    loading.value = true
    error.value = null

    try {
      const response = await authApi.register(data)
      if (response.data.data) {
        token.value = response.data.data.token
        user.value = response.data.data.user

        localStorage.setItem('token', response.data.data.token)
        localStorage.setItem('user', JSON.stringify(response.data.data.user))
      }
    } catch (err: any) {
      error.value = err.response?.data?.message || 'Registration failed'
      throw err
    } finally {
      loading.value = false
    }
  }

  function logout() {
    user.value = null
    token.value = null
    localStorage.removeItem('token')
    localStorage.removeItem('user')
  }

  return {
    user,
    token,
    loading,
    error,
    isAuthenticated,
    isAdmin,
    isTokenExpired,
    init,
    login,
    register,
    logout
  }
})
