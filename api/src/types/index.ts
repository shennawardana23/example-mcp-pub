export interface Service {
  id: number
  name: string
  description: string
  owner: string
  type: 'api' | 'library' | 'microservice' | 'frontend'
  status: 'active' | 'deprecated' | 'planning'
  version?: string
  repository?: string
  docs_url?: string
  tags?: string[]
  created_at: string
  updated_at: string
}

export interface User {
  id: number
  username: string
  email: string
  first_name: string
  last_name: string
  role: string
  created_at: string
}

export interface LoginRequest {
  username: string
  password: string
}

export interface LoginResponse {
  access_token: string
  refresh_token: string
  token_type: string
  expires_in: number
}

export interface PaginatedResponse<T> {
  data: T[]
  page: number
  limit: number
  total_count: number
  total_pages: number
}
