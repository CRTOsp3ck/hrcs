export interface User {
  id: number
  email: string
  name: string
  role: 'admin' | 'normal'
  group_id?: number
  group?: UserGroup
  created_at: string
  updated_at: string
}

export interface UserGroup {
  id: number
  name: string
  description?: string
  created_at: string
  updated_at: string
}

export interface ClaimType {
  id: number
  name: string
  description?: string
  is_active: boolean
  created_at: string
  updated_at: string
}

export interface Claim {
  id: number
  user_id: number
  user?: User
  claim_type_id: number
  claim_type?: ClaimType
  title: string
  description: string
  amount: number
  status: ClaimStatus
  submitted_at?: string
  approved_at?: string
  rejected_at?: string
  paid_at?: string
  approvals?: Approval[]
  created_at: string
  updated_at: string
}

export type ClaimStatus = 
  | 'draft'
  | 'submitted'
  | 'approved'
  | 'rejected'
  | 'payment-in-progress'
  | 'paid'

export interface Approval {
  id: number
  claim_id: number
  user_id: number
  user?: User
  level_id: number
  level?: ApprovalLevel
  action: 'approve' | 'reject'
  comments?: string
  created_at: string
}

export interface ApprovalLevel {
  id: number
  group_id: number
  group?: UserGroup
  level: number
  name: string
  approvers: User[]
  permissions: ('approve' | 'reject')[]
  created_at: string
  updated_at: string
}

export interface LoginRequest {
  email: string
  password: string
}

export interface RegisterRequest {
  email: string
  password: string
  name: string
}

export interface ApiResponse<T = any> {
  success: boolean
  data?: T
  message?: string
  error?: string
}

export interface DashboardStats {
  totalClaims: number
  pendingClaims: number
  approvedClaims: number
  rejectedClaims: number
  totalAmount: number
  approvedAmount: number
  recentClaims: Claim[]
  claimsByStatus: { status: string; count: number }[]
  claimsByType: { type: string; count: number; amount: number }[]
}