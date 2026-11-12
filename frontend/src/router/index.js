import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '../stores/auth'

const routes = [
  {
    path: '/',
    name: 'landing',
    component: () => import('../views/Landing.vue')
  },
  // Auth (guest-only)
  {
    path: '/register',
    name: 'register',
    component: () => import('../views/Register.vue'),
    meta: { guestOnly: true }
  },
  {
    path: '/login',
    name: 'login',
    component: () => import('../views/Login.vue'),
    meta: { guestOnly: true }
  },

  // Authenticated (any role)
  {
    path: '/onboarding',
    name: 'onboarding',
    component: () => import('../views/Onboarding.vue'),
    meta: { requiresAuth: true }
  },
  {
    path: '/dashboard',
    name: 'dashboard',
    component: () => import('../views/Dashboard.vue'),
    meta: { requiresAuth: true }
  },
  {
    path: '/settings',
    name: 'settings',
    component: () => import('../views/Settings.vue'),
    meta: { requiresAuth: true }
  },
  {
    path: '/notifications',
    name: 'notifications',
    component: () => import('../views/Notifications.vue'),
    meta: { requiresAuth: true }
  },

  // Public browse
  {
    path: '/browse',
    name: 'browse',
    component: () => import('../views/BrowseJobs.vue')
  },
  {
    path: '/browse/:id',
    name: 'job-detail',
    component: () => import('../views/PublicJobDetail.vue')
  },
  {
    path: '/companies',
    name: 'companies',
    component: () => import('../views/CompaniesList.vue')
  },
  {
    path: '/companies/:id',
    name: 'company-public',
    component: () => import('../views/CompanyPublic.vue')
  },

  // Seeker-only
  {
    path: '/profile',
    name: 'seeker-profile',
    component: () => import('../views/SeekerProfile.vue'),
    meta: { requiresAuth: true, roles: ['job_seeker'] }
  },
  {
    path: '/applications',
    name: 'my-applications',
    component: () => import('../views/MyApplications.vue'),
    meta: { requiresAuth: true, roles: ['job_seeker'] }
  },
  {
    path: '/saved',
    name: 'saved-jobs',
    component: () => import('../views/SavedJobs.vue'),
    meta: { requiresAuth: true, roles: ['job_seeker'] }
  },
  {
    path: '/searches',
    name: 'saved-searches',
    component: () => import('../views/SavedSearches.vue'),
    meta: { requiresAuth: true, roles: ['job_seeker'] }
  },

  // Employer-only
  {
    path: '/company',
    name: 'company-profile',
    component: () => import('../views/CompanyProfile.vue'),
    meta: { requiresAuth: true, roles: ['employer'] }
  },
  {
    path: '/jobs',
    name: 'jobs',
    component: () => import('../views/JobsList.vue'),
    meta: { requiresAuth: true, roles: ['employer'] }
  },
  {
    path: '/jobs/new',
    name: 'job-create',
    component: () => import('../views/JobCreate.vue'),
    meta: { requiresAuth: true, roles: ['employer'] }
  },
  {
    path: '/jobs/:id/edit',
    name: 'job-edit',
    component: () => import('../views/JobEdit.vue'),
    meta: { requiresAuth: true, roles: ['employer'] }
  },
  {
    path: '/employer/applications',
    name: 'employer-applications',
    component: () => import('../views/EmployerApplications.vue'),
    meta: { requiresAuth: true, roles: ['employer'] }
  },
  {
    path: '/employer/applications/:id',
    name: 'employer-application-detail',
    component: () => import('../views/EmployerApplicationDetail.vue'),
    meta: { requiresAuth: true, roles: ['employer'] }
  },

  // Admin-only
  {
    path: '/admin',
    name: 'admin-dashboard',
    component: () => import('../views/AdminDashboard.vue'),
    meta: { requiresAuth: true, roles: ['admin'] }
  },
  {
    path: '/admin/jobs',
    name: 'admin-jobs',
    component: () => import('../views/AdminJobs.vue'),
    meta: { requiresAuth: true, roles: ['admin'] }
  },
  {
    path: '/admin/jobs/:id',
    name: 'admin-job-detail',
    component: () => import('../views/AdminJobDetail.vue'),
    meta: { requiresAuth: true, roles: ['admin'] }
  },
  {
    path: '/admin/users',
    name: 'admin-users',
    component: () => import('../views/AdminUsers.vue'),
    meta: { requiresAuth: true, roles: ['admin'] }
  },
  {
    path: '/admin/users/:id',
    name: 'admin-user-detail',
    component: () => import('../views/AdminUserDetail.vue'),
    meta: { requiresAuth: true, roles: ['admin'] }
  },
  {
    path: '/admin/applications',
    name: 'admin-applications',
    component: () => import('../views/AdminApplications.vue'),
    meta: { requiresAuth: true, roles: ['admin'] }
  },
  {
    path: '/admin/applications/:id',
    name: 'admin-application-detail',
    component: () => import('../views/AdminApplicationDetail.vue'),
    meta: { requiresAuth: true, roles: ['admin'] }
  },
    {
    path: '/forgot-password',
    name: 'forgot-password',
    component: () => import('../views/ForgotPassword.vue'),
    meta: { guestOnly: true }
  },
  {
    path: '/reset-password',
    name: 'reset-password',
    component: () => import('../views/ResetPassword.vue')
  },
  {
    path: '/verify-email',
    name: 'verify-email',
    component: () => import('../views/VerifyEmail.vue')
  },
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

router.beforeEach((to) => {
  const auth = useAuthStore()

  // Landing shows to anonymous users; logged-in users go to dashboard.
  if (to.name === 'landing' && auth.isAuthenticated) {
    return { name: 'dashboard' }
  }

  if (to.meta.requiresAuth && !auth.isAuthenticated) {
    return { name: 'login', query: { redirect: to.fullPath } }
  }

  if (
    to.meta.requiresAuth &&
    to.name !== 'onboarding' &&
    auth.user &&
    auth.user.onboardingCompleted === false
  ) {
    return { name: 'onboarding' }
  }

  if (
    to.name === 'onboarding' &&
    auth.user &&
    auth.user.onboardingCompleted === true
  ) {
    return { name: 'dashboard' }
  }

  if (to.meta.roles && auth.user && !to.meta.roles.includes(auth.user.role)) {
    return { name: 'dashboard' }
  }

  return true
})

export default router
