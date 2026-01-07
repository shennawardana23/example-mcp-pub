import { defineStore } from 'pinia'
import { ref } from 'vue'
import apiClient from '../utils/api'
import type { Service, PaginatedResponse } from '../types'

export const useServiceStore = defineStore('service', () => {
  const services = ref<Service[]>([])
  const currentService = ref<Service | null>(null)
  const loading = ref(false)
  const totalCount = ref(0)

  async function fetchServices(filters?: {
    search?: string
    type?: string
    status?: string
    page?: number
    limit?: number
  }) {
    loading.value = true
    try {
      const response = await apiClient.get<PaginatedResponse<Service>>('/api/v1/services', {
        params: filters,
      })
      services.value = response.data.data
      totalCount.value = response.data.total_count
    } catch (error) {
      console.error('Failed to fetch services:', error)
    } finally {
      loading.value = false
    }
  }

  async function fetchServiceById(id: number) {
    loading.value = true
    try {
      const response = await apiClient.get<Service>(`/api/v1/services/${id}`)
      currentService.value = response.data
      return response.data
    } catch (error) {
      console.error('Failed to fetch service:', error)
      return null
    } finally {
      loading.value = false
    }
  }

  async function createService(service: Partial<Service>) {
    try {
      const response = await apiClient.post<Service>('/api/v1/services', service)
      services.value.unshift(response.data)
      return response.data
    } catch (error) {
      console.error('Failed to create service:', error)
      throw error
    }
  }

  async function updateService(id: number, service: Partial<Service>) {
    try {
      const response = await apiClient.put<Service>(`/api/v1/services/${id}`, service)
      const index = services.value.findIndex((s) => s.id === id)
      if (index !== -1) {
        services.value[index] = response.data
      }
      return response.data
    } catch (error) {
      console.error('Failed to update service:', error)
      throw error
    }
  }

  async function deleteService(id: number) {
    try {
      await apiClient.delete(`/api/v1/services/${id}`)
      services.value = services.value.filter((s) => s.id !== id)
    } catch (error) {
      console.error('Failed to delete service:', error)
      throw error
    }
  }

  return {
    services,
    currentService,
    loading,
    totalCount,
    fetchServices,
    fetchServiceById,
    createService,
    updateService,
    deleteService,
  }
})
