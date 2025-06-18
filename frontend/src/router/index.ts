import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/',
      redirect: '/dashboard'
    },
    {
      path: '/login',
      name: 'login',
      component: () => import('@/views/LoginView.vue'),
      meta: { requiresGuest: true }
    },
    {
      path: '/register',
      name: 'register',
      component: () => import('@/views/RegisterView.vue'),
      meta: { requiresGuest: true }
    },
    {
      path: '/dashboard',
      name: 'dashboard',
      component: () => import('@/views/DashboardView.vue'),
      meta: { requiresAuth: true }
    },
    {
      path: '/claims',
      name: 'claims',
      component: () => import('@/views/ClaimsView.vue'),
      meta: { requiresAuth: true }
    },
    {
      path: '/claims/new',
      name: 'new-claim',
      component: () => import('@/views/NewClaimView.vue'),
      meta: { requiresAuth: true }
    },
    {
      path: '/claims/:id',
      name: 'claim-detail',
      component: () => import('@/views/ClaimDetailView.vue'),
      meta: { requiresAuth: true }
    },
    {
      path: '/claims/:id/edit',
      name: 'edit-claim',
      component: () => import('@/views/EditClaimView.vue'),
      meta: { requiresAuth: true }
    },
    {
      path: '/admin',
      name: 'admin',
      component: () => import('@/views/AdminView.vue'),
      meta: { requiresAuth: true, requiresAdmin: true },
      children: [
        {
          path: '',
          redirect: '/admin/dashboard'
        },
        {
          path: 'dashboard',
          name: 'admin-dashboard',
          component: () => import('@/views/admin/AdminDashboard.vue')
        },
        {
          path: 'users',
          name: 'admin-users',
          component: () => import('@/views/admin/AdminUsers.vue')
        },
        {
          path: 'users/:id',
          name: 'admin-user-details',
          component: () => import('@/views/admin/UserDetailsView.vue')
        },
        {
          path: 'groups',
          name: 'admin-groups',
          component: () => import('@/views/admin/AdminGroups.vue')
        },
        {
          path: 'groups/:id',
          name: 'admin-group-details',
          component: () => import('@/views/admin/UserGroupDetailsView.vue')
        },
        {
          path: 'claim-types',
          name: 'admin-claim-types',
          component: () => import('@/views/admin/AdminClaimTypes.vue')
        },
        {
          path: 'claim-types/:id',
          name: 'admin-claim-type-details',
          component: () => import('@/views/admin/ClaimTypeDetailsView.vue')
        },
        {
          path: 'approval-levels',
          name: 'admin-approval-levels',
          component: () => import('@/views/admin/AdminApprovalLevels.vue')
        },
        {
          path: 'claims',
          name: 'admin-claims',
          component: () => import('@/views/admin/AdminClaims.vue')
        },
        // NEW: Phase 3 admin routes
        {
          path: 'audit-log',
          name: 'admin-audit-log',
          component: () => import('@/views/admin/AdminAuditLog.vue')
        },
        {
          path: 'reports',
          name: 'admin-reports',
          component: () => import('@/views/admin/AdminReports.vue')
        },
        {
          path: 'integrations',
          name: 'admin-integrations',
          component: () => import('@/views/admin/AdminIntegrations.vue')
        }
      ]
    }
  ]
})

// Navigation guards
router.beforeEach(async (to, from, next) => {
  const authStore = useAuthStore()

  // Initialize auth store if not already done
  if (!authStore.isAuthenticated && localStorage.getItem('token')) {
    await authStore.init()
  }

  const requiresAuth = to.matched.some(record => record.meta.requiresAuth)
  const requiresGuest = to.matched.some(record => record.meta.requiresGuest)
  const requiresAdmin = to.matched.some(record => record.meta.requiresAdmin)

  if (requiresAuth && !authStore.isAuthenticated) {
    next('/login')
  } else if (requiresGuest && authStore.isAuthenticated) {
    next('/dashboard')
  } else if (requiresAdmin && !authStore.isAdmin) {
    next('/dashboard')
  } else {
    next()
  }
})

export default router
