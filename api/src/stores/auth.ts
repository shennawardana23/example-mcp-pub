import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import apiClient from '../utils/api'
import type { User, LoginRequest, LoginResponse } from '../types'

export const useAuthStore = defineStore('auth', () => {
  const user = ref<User | null>(null)
  const accessToken = ref<string | null>(localStorage.getItem('access_token'))
  const isAuthenticated = computed(() => !!accessToken.value)

  async function login(credentials: LoginRequest) {
    try {
      const response = await apiClient.post<LoginResponse>('/api/v1/auth/login', credentials)
      accessToken.value = response.data.access_token
      localStorage.setItem('access_token', response.data.access_token)
      await fetchCurrentUser()
      return true
    } catch (error) {
      console.error('Login failed:', error)
      return false
    }
  }

  async function fetchCurrentUser() {
    try {
      const response = await apiClient.get<User>('/api/v1/auth/me')
      user.value = response.data
    } catch (error) {
      console.error('Failed to fetch user:', error)
      logout()
    }
  }

  function logout() {
    user.value = null
    accessToken.value = null
    localStorage.removeItem('access_token')
  }

  return {
    user,
    accessToken,
    isAuthenticated,
    login,
    logout,
    fetchCurrentUser,
  }
})
