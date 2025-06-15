export interface User {
  id: number
  email: string
  first_name: string
  last_name: string
  name?: string // Optional display name
  role: 'admin' | 'normal'
  user_group_id?: number
  user_group?: UserGroup
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
  // NEW FIELDS FOR CLAIM LIMITS
  limit_amount: number
  limit_timespan: 'annual' | 'monthly' | 'weekly' | 'daily'
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
  // Enhanced fields for admin view
  employee?: string
  department?: string
  type?: string
  approvalsReceived?: number
  approvalsRequired?: number
  submittedDate?: string
  canApprove?: boolean
  allowedStatuses?: string[]
  approvalWorkflow?: ApprovalStep[]
  currentStep?: ApprovalStep
  nextSteps?: ApprovalStep[]
}

export interface ApprovalStep {
  id: number
  level: number
  name: string
  approverId: number
  approverName: string
  approverEmail: string
  userGroupId: number
  userGroupName: string
  status: 'pending' | 'approved' | 'rejected' | 'skipped'
  completedAt?: string
  comments: string
  permissions: ApprovalPermissions
}

export interface ApprovalPermissions {
  canDraft: boolean
  canSubmit: boolean
  canApprove: boolean
  canReject: boolean
  canSetPaymentInProgress: boolean
  canSetPaid: boolean
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
  level: number
  user_group_id: number
  user_group?: UserGroup
  approver_id: number
  approver?: User
  // Status permissions
  can_draft: boolean
  can_submit: boolean
  can_approve: boolean
  can_reject: boolean
  can_set_payment_in_progress: boolean
  can_set_paid: boolean
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

// NEW INTERFACES FOR BALANCE TRACKING

export interface UserClaimBalance {
  id: number
  user_id: number
  claim_type_id: number
  total_limit: number
  current_spent: number
  remaining_balance: number
  last_reset_date: string
  reset_period: 'annual' | 'monthly' | 'weekly' | 'daily'
  claim_type: ClaimType
  created_at: string
  updated_at: string
}

export interface BalanceInfo {
  claim_type_id: number
  claim_type_name: string
  total_limit: number
  current_spent: number
  remaining_balance: number
  can_claim_amount: number
}

export interface UserGroupClaimPermission {
  id: number
  user_group_id: number
  claim_type_id: number
  is_allowed: boolean
  custom_limit_amount?: number
  user_group: UserGroup
  claim_type: ClaimType
  created_at: string
  updated_at: string
}

export interface UserClaimPermission {
  id: number
  user_id: number
  claim_type_id: number
  is_allowed: boolean
  custom_limit_amount?: number
  user: User
  claim_type: ClaimType
  created_at: string
  updated_at: string
}

// DETAIL VIEW INTERFACES

export interface UserDetails {
  user: User
  balances: UserClaimBalance[]
  claims: Claim[]
  permissions: UserClaimPermission[]
}

export interface UserGroupDetails {
  group: UserGroup
  members: User[]
  permissions: UserGroupClaimPermission[]
  approval_levels: ApprovalLevel[]
}

export interface ClaimTypeDetails {
  claim_type: ClaimType
  stats: {
    total_claims: number
    approved_claims: number
    total_amount: number
    average_amount: number
  }
  group_permissions: UserGroupClaimPermission[]
  recent_claims: Claim[]
}