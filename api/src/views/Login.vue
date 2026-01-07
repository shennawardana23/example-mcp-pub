<template>
  <div class="min-h-screen bg-gray-100">
    <div class="flex items-center justify-center min-h-screen">
      <div class="w-full max-w-md p-8 bg-white rounded-lg shadow-md">
        <h1 class="text-3xl font-bold text-center mb-8">Developer Portal</h1>
        <form @submit.prevent="handleLogin" class="space-y-6">
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-2">Username</label>
            <input
              v-model="credentials.username"
              type="text"
              required
              class="w-full px-4 py-2 border border-gray-300 rounded-md focus:ring-2 focus:ring-blue-500 focus:border-transparent"
              placeholder="Enter username"
            />
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-2">Password</label>
            <input
              v-model="credentials.password"
              type="password"
              required
              class="w-full px-4 py-2 border border-gray-300 rounded-md focus:ring-2 focus:ring-blue-500 focus:border-transparent"
              placeholder="Enter password"
            />
          </div>
          <button
            type="submit"
            :disabled="loading"
            class="w-full bg-blue-600 text-white py-2 px-4 rounded-md hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2 disabled:opacity-50"
          >
            {{ loading ? 'Logging in...' : 'Login' }}
          </button>
          <p v-if="error" class="text-red-500 text-sm text-center">{{ error }}</p>
          <p class="text-sm text-gray-600 text-center">Demo: username: admin, password: admin123</p>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'

const router = useRouter()
const authStore = useAuthStore()

const credentials = ref({ username: '', password: '' })
const loading = ref(false)
const error = ref('')

async function handleLogin() {
  loading.value = true
  error.value = ''
  
  const success = await authStore.login(credentials.value)
  
  if (success) {
    router.push('/')
  } else {
    error.value = 'Invalid username or password'
  }
  
  loading.value = false
}
</script>
