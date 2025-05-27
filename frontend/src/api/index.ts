import axios, { AxiosError } from 'axios'
import type { ApiResponse, User, LoginRequest, RegisterRequest, Claim, ClaimType, UserGroup, ApprovalLevel, DashboardStats } from '@/types'

// const API_BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:8000/api'
const API_BASE_URL = 'http://localhost:8000/api'

const api = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    'Content-Type': 'application/json'
  }
})

// Request interceptor to add auth token
api.interceptors.request.use((config) => {
  const token = localStorage.getItem('token')
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})

// Response interceptor to handle errors
api.interceptors.response.use(
  (response) => response,
  (error: AxiosError) => {
    if (error.response?.status === 401) {
      localStorage.removeItem('token')
      localStorage.removeItem('user')
      window.location.href = '/login'
    }
    return Promise.reject(error)
  }
)

// Auth API
export const authApi = {
  login: (data: LoginRequest) => api.post<ApiResponse<{ token: string; user: User }>>('/auth/login', data),
  register: (data: RegisterRequest) => api.post<ApiResponse<{ token: string; user: User }>>('/auth/register', data),
  getProfile: () => api.get<ApiResponse<User>>('/profile')
}

// Claims API
export const claimsApi = {
  getAll: () => api.get<ApiResponse<Claim[]>>('/claims'),
  getById: (id: number) => api.get<ApiResponse<Claim>>(`/claims/${id}`),
  create: (data: Partial<Claim>) => api.post<ApiResponse<Claim>>('/claims', data),
  update: (id: number, data: Partial<Claim>) => api.put<ApiResponse<Claim>>(`/claims/${id}`, data),
  delete: (id: number) => api.delete<ApiResponse>(`/claims/${id}`),
  submit: (id: number) => api.post<ApiResponse<Claim>>(`/claims/${id}/submit`),
  approve: (id: number, data: { comments?: string }) => api.post<ApiResponse<Claim>>(`/claims/${id}/approve`, data),
  reject: (id: number, data: { comments?: string }) => api.post<ApiResponse<Claim>>(`/claims/${id}/reject`, data)
}

// Users API (Admin)
export const usersApi = {
  getAll: () => api.get<ApiResponse<User[]>>('/users'),
  updateRole: (id: number, role: 'admin' | 'normal') => api.put<ApiResponse<User>>(`/users/${id}/role`, { role }),
  updateGroup: (id: number, groupId: number) => api.put<ApiResponse<User>>(`/users/${id}/group`, { group_id: groupId })
}

// Claim Types API
export const claimTypesApi = {
  getAll: () => api.get<ApiResponse<ClaimType[]>>('/claim-types'),
  create: (data: Partial<ClaimType>) => api.post<ApiResponse<ClaimType>>('/claim-types', data),
  update: (id: number, data: Partial<ClaimType>) => api.put<ApiResponse<ClaimType>>(`/claim-types/${id}`, data),
  delete: (id: number) => api.delete<ApiResponse>(`/claim-types/${id}`)
}

// User Groups API
export const userGroupsApi = {
  getAll: () => api.get<ApiResponse<UserGroup[]>>('/user-groups'),
  create: (data: Partial<UserGroup>) => api.post<ApiResponse<UserGroup>>('/user-groups', data),
  update: (id: number, data: Partial<UserGroup>) => api.put<ApiResponse<UserGroup>>(`/user-groups/${id}`, data),
  delete: (id: number) => api.delete<ApiResponse>(`/user-groups/${id}`)
}

// Approval Levels API
export const approvalLevelsApi = {
  getAll: () => api.get<ApiResponse<ApprovalLevel[]>>('/approval-levels'),
  getByGroup: (groupId: number) => api.get<ApiResponse<ApprovalLevel[]>>(`/approval-levels/group/${groupId}`),
  create: (data: Partial<ApprovalLevel>) => api.post<ApiResponse<ApprovalLevel>>('/approval-levels', data),
  update: (id: number, data: Partial<ApprovalLevel>) => api.put<ApiResponse<ApprovalLevel>>(`/approval-levels/${id}`, data),
  delete: (id: number) => api.delete<ApiResponse>(`/approval-levels/${id}`),
  addApprover: (levelId: number, userId: number) => api.post<ApiResponse>(`/approval-levels/${levelId}/approvers/${userId}`),
  removeApprover: (levelId: number, userId: number) => api.delete<ApiResponse>(`/approval-levels/${levelId}/approvers/${userId}`)
}

// Dashboard API
export const dashboardApi = {
  getStats: () => api.get<ApiResponse<DashboardStats>>('/dashboard/stats'),
  getAdminStats: () => api.get<ApiResponse<DashboardStats>>('/dashboard/admin-stats')
}

// Admin API - Consolidated admin functions
export const adminApi = {
  // Claims management
  getAllClaims: () => api.get<ApiResponse<Claim[]>>('/admin/claims'),
  approveClaim: (id: number, data: { comments?: string }) => api.post<ApiResponse<Claim>>(`/admin/claims/${id}/approve`, data),
  rejectClaim: (id: number, data: { comments?: string }) => api.post<ApiResponse<Claim>>(`/admin/claims/${id}/reject`, data),
  updateClaimStatus: (id: number, data: { status: string; comments?: string }) => api.put<ApiResponse<Claim>>(`/admin/claims/${id}/status`, data),

  // Users management
  getUsers: () => api.get<ApiResponse<User[]>>('/admin/users'),
  createUser: (data: Partial<User>) => api.post<ApiResponse<User>>('/admin/users', data),
  updateUser: (id: number, data: Partial<User>) => api.put<ApiResponse<User>>(`/admin/users/${id}`, data),
  deleteUser: (id: number) => api.delete<ApiResponse>(`/admin/users/${id}`),

  // Groups management
  getGroups: () => api.get<ApiResponse<UserGroup[]>>('/admin/groups'),
  createGroup: (data: Partial<UserGroup>) => api.post<ApiResponse<UserGroup>>('/admin/groups', data),
  updateGroup: (id: number, data: Partial<UserGroup>) => api.put<ApiResponse<UserGroup>>(`/admin/groups/${id}`, data),
  deleteGroup: (id: number) => api.delete<ApiResponse>(`/admin/groups/${id}`),

  // Claim Types management
  getClaimTypes: () => api.get<ApiResponse<ClaimType[]>>('/admin/claim-types'),
  createClaimType: (data: Partial<ClaimType>) => api.post<ApiResponse<ClaimType>>('/admin/claim-types', data),
  updateClaimType: (id: number, data: Partial<ClaimType>) => api.put<ApiResponse<ClaimType>>(`/admin/claim-types/${id}`, data),
  deleteClaimType: (id: number) => api.delete<ApiResponse>(`/admin/claim-types/${id}`),

  // Approval Levels management
  getApprovalLevels: () => api.get<ApiResponse<ApprovalLevel[]>>('/admin/approval-levels'),
  createApprovalLevel: (data: Partial<ApprovalLevel>) => api.post<ApiResponse<ApprovalLevel>>('/admin/approval-levels', data),
  updateApprovalLevel: (id: number, data: Partial<ApprovalLevel>) => api.put<ApiResponse<ApprovalLevel>>(`/admin/approval-levels/${id}`, data),
  deleteApprovalLevel: (id: number) => api.delete<ApiResponse>(`/admin/approval-levels/${id}`),
  updateApprovalLevelOrder: (orders: Array<{ id: number; order: number }>) => api.put<ApiResponse>('/admin/approval-levels/order', { orders })
}

export default api
