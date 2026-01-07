<template>
  <div class="min-h-screen bg-gray-100">
    <nav class="bg-white shadow-sm">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="flex justify-between h-16">
          <div class="flex items-center">
            <h1 class="text-xl font-bold">Developer Portal</h1>
          </div>
          <div class="flex items-center space-x-4">
            <router-link to="/" class="text-gray-700 hover:text-gray-900">Dashboard</router-link>
            <span class="text-gray-700">{{ authStore.user?.username }}</span>
            <button @click="handleLogout" class="text-gray-700 hover:text-gray-900">Logout</button>
          </div>
        </div>
      </div>
    </nav>

    <div class="max-w-7xl mx-auto py-12 px-4 sm:px-6 lg:px-8">
      <div class="flex justify-between items-center mb-8">
        <h2 class="text-3xl font-bold text-gray-900">Service Catalog</h2>
      </div>

      <div class="mb-6 grid grid-cols-1 md:grid-cols-4 gap-4">
        <input
          v-model="filters.search"
          @input="handleSearch"
          type="text"
          placeholder="Search services..."
          class="px-4 py-2 border border-gray-300 rounded-md focus:ring-2 focus:ring-blue-500"
        />
        <select v-model="filters.type" @change="handleFilterChange" class="px-4 py-2 border border-gray-300 rounded-md">
          <option value="">All Types</option>
          <option value="api">API</option>
          <option value="microservice">Microservice</option>
          <option value="library">Library</option>
          <option value="frontend">Frontend</option>
        </select>
        <select v-model="filters.status" @change="handleFilterChange" class="px-4 py-2 border border-gray-300 rounded-md">
          <option value="">All Status</option>
          <option value="active">Active</option>
          <option value="deprecated">Deprecated</option>
          <option value="planning">Planning</option>
        </select>
      </div>

      <div v-if="serviceStore.loading" class="text-center py-12">
        <p class="text-gray-600">Loading services...</p>
      </div>

      <div v-else class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        <div
          v-for="service in serviceStore.services"
          :key="service.id"
          class="bg-white rounded-lg shadow p-6 hover:shadow-lg transition"
        >
          <div class="flex justify-between items-start mb-4">
            <h3 class="text-xl font-semibold">{{ service.name }}</h3>
            <span
              :class="{
                'bg-green-100 text-green-800': service.status === 'active',
                'bg-red-100 text-red-800': service.status === 'deprecated',
                'bg-yellow-100 text-yellow-800': service.status === 'planning',
              }"
              class="px-2 py-1 text-xs rounded-full"
            >
              {{ service.status }}
            </span>
          </div>
          <p class="text-gray-600 mb-4">{{ service.description }}</p>
          <div class="space-y-2 text-sm">
            <p><span class="font-semibold">Type:</span> {{ service.type }}</p>
            <p><span class="font-semibold">Owner:</span> {{ service.owner }}</p>
            <p v-if="service.version"><span class="font-semibold">Version:</span> {{ service.version }}</p>
          </div>
          <div v-if="service.tags && service.tags.length" class="mt-4 flex flex-wrap gap-2">
            <span
              v-for="tag in service.tags"
              :key="tag"
              class="px-2 py-1 bg-blue-100 text-blue-800 text-xs rounded"
            >
              {{ tag }}
            </span>
          </div>
          <div v-if="service.repository || service.docs_url" class="mt-4 flex gap-2">
            <a
              v-if="service.repository"
              :href="service.repository"
              target="_blank"
              class="text-blue-600 hover:text-blue-800 text-sm"
            >
              Repository →
            </a>
            <a
              v-if="service.docs_url"
              :href="service.docs_url"
              target="_blank"
              class="text-blue-600 hover:text-blue-800 text-sm"
            >
              Docs →
            </a>
          </div>
        </div>
      </div>

      <div v-if="!serviceStore.loading && serviceStore.services.length === 0" class="text-center py-12">
        <p class="text-gray-600">No services found</p>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'
import { useServiceStore } from '../stores/service'

const router = useRouter()
const authStore = useAuthStore()
const serviceStore = useServiceStore()

const filters = ref({
  search: '',
  type: '',
  status: '',
})

onMounted(() => {
  serviceStore.fetchServices()
})

function handleSearch() {
  serviceStore.fetchServices(filters.value)
}

function handleFilterChange() {
  serviceStore.fetchServices(filters.value)
}

function handleLogout() {
  authStore.logout()
  router.push('/login')
}
</script>
