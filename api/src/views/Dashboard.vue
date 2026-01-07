<template>
  <div class="min-h-screen bg-gray-100">
    <nav class="bg-white shadow-sm">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="flex justify-between h-16">
          <div class="flex items-center">
            <h1 class="text-xl font-bold">Developer Portal</h1>
          </div>
          <div class="flex items-center space-x-4">
            <router-link to="/services" class="text-gray-700 hover:text-gray-900">Services</router-link>
            <span class="text-gray-700">{{ authStore.user?.username }}</span>
            <button @click="handleLogout" class="text-gray-700 hover:text-gray-900">Logout</button>
          </div>
        </div>
      </div>
    </nav>

    <div class="max-w-7xl mx-auto py-12 px-4 sm:px-6 lg:px-8">
      <h2 class="text-3xl font-bold text-gray-900 mb-8">Dashboard</h2>
      
      <div class="grid grid-cols-1 md:grid-cols-3 gap-6 mb-12">
        <div class="bg-white p-6 rounded-lg shadow">
          <h3 class="text-lg font-semibold mb-2">Total Services</h3>
          <p class="text-3xl font-bold text-blue-600">{{ serviceStore.totalCount }}</p>
        </div>
        <div class="bg-white p-6 rounded-lg shadow">
          <h3 class="text-lg font-semibold mb-2">Active Services</h3>
          <p class="text-3xl font-bold text-green-600">{{ activeServices }}</p>
        </div>
        <div class="bg-white p-6 rounded-lg shadow">
          <h3 class="text-lg font-semibold mb-2">Deprecated</h3>
          <p class="text-3xl font-bold text-red-600">{{ deprecatedServices }}</p>
        </div>
      </div>

      <div class="bg-white rounded-lg shadow p-6">
        <h3 class="text-xl font-semibold mb-4">Quick Actions</h3>
        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
          <router-link
            to="/services"
            class="p-4 border border-gray-200 rounded-lg hover:bg-gray-50 transition"
          >
            <h4 class="font-semibold mb-2">Service Catalog</h4>
            <p class="text-sm text-gray-600">Browse and manage services</p>
          </router-link>
          <div class="p-4 border border-gray-200 rounded-lg opacity-50 cursor-not-allowed">
            <h4 class="font-semibold mb-2">API Documentation</h4>
            <p class="text-sm text-gray-600">View API docs (Coming soon)</p>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'
import { useServiceStore } from '../stores/service'

const router = useRouter()
const authStore = useAuthStore()
const serviceStore = useServiceStore()

const activeServices = computed(() => 
  serviceStore.services.filter(s => s.status === 'active').length
)

const deprecatedServices = computed(() => 
  serviceStore.services.filter(s => s.status === 'deprecated').length
)

onMounted(() => {
  serviceStore.fetchServices()
})

function handleLogout() {
  authStore.logout()
  router.push('/login')
}
</script>
