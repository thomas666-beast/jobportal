import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import * as authApi from '../api/auth'
import { apiFetch } from '../api/client'

export const useAuthStore = defineStore('auth', () => {
  const user = ref(null)
  const loading = ref(true) // true on initial load while we verify session

  const isAuthenticated = computed(() => !!user.value)
  const role = computed(() => user.value?.role || null)
  const emailVerified = computed(() => user.value?.emailVerified === true)
  const needsOnboarding = computed(
    () => isAuthenticated.value && user.value?.onboardingCompleted === false
  )

  async function init() {
    try {
      const res = await apiFetch('/me', { auth: true })
      user.value = {
        id: res.user_id,
        username: res.username,
        role: res.role,
        onboardingCompleted: res.onboarding_completed,
        emailVerified: res.email_verified
      }
    } catch {
      user.value = null
    } finally {
      loading.value = false
    }
  }

  async function login(credentials) {
    const res = await authApi.login(credentials)
    user.value = {
      id: res.user_id,
      username: res.username,
      role: res.role,
      onboardingCompleted: res.onboarding_completed,
      emailVerified: res.email_verified
    }
    return res
  }

  async function register(payload) {
    const res = await authApi.register(payload)
    user.value = {
      id: res.user_id,
      username: res.username,
      role: res.role,
      onboardingCompleted: res.onboarding_completed,
      emailVerified: res.email_verified
    }
    return res
  }

  async function fetchMe() {
    const res = await apiFetch('/me', { auth: true, skipRefresh: true })
    user.value = {
      id: res.user_id,
      username: res.username,
      role: res.role,
      onboardingCompleted: res.onboarding_completed,
      emailVerified: res.email_verified
    }
    return res
  }

  async function completeOnboarding(selectedRole) {
    const res = await authApi.setRole(selectedRole)
    user.value = {
      id: res.user_id,
      username: res.username,
      role: res.role,
      onboardingCompleted: res.onboarding_completed,
      emailVerified: res.email_verified
    }
    return res
  }

  async function logout() {
    try {
      await authApi.logout()
    } catch {
      // Ignore network errors on logout
    }
    user.value = null
  }

  return {
    user,
    loading,
    role,
    emailVerified,
    isAuthenticated,
    needsOnboarding,
    init,
    login,
    register,
    fetchMe,
    completeOnboarding,
    logout
  }
})
