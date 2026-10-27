export type Classification = 'PUBLIC' | 'INTERNAL' | 'RESTRICTED' | 'CONFIDENTIAL'
export type Category = 'SPEC' | 'INCIDENT' | 'SCHEMA' | 'DISPATCH'
export type Status = 'ACTIVE' | 'ARCHIVED' | 'DRAFT'
export type UserRole = 'ARCHIVIST' | 'OPERATOR' | 'VIEWER'

export interface User {
  id: number
  username: string
  email: string
  role: UserRole
  created_at: string
}

export interface RecordItem {
  id: number
  code: string
  title: string
  classification: Classification
  category: Category
  status: Status
  tags: string
  summary: string
  payload: string
  created_by: string
  created_at: string
  updated_at: string
}

export interface RecordPayload {
  title: string
  classification: Classification
  category: Category
  status: Status
  tags: string
  summary: string
  payload: string
}

export interface StatsData {
  total_records: number
  by_classification: Record<Classification, number>
  by_category: Record<Category, number>
  by_status: Record<Status, number>
  latest_code: string
  uptime_seconds: number
}

export interface ApiResponse<T = any> {
  success: boolean
  message?: string
  data?: T
  error?: string
}

export interface AuthResponse {
  token: string
  user: User
}

export interface FilterState {
  search: string
  category: string
  classification: string
  status: string
}
