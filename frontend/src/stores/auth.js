import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import * as authApi from '../api/auth'

export const useAuthStore = defineStore('auth', () => {
  const token = ref(localStorage.getItem('token') || '')
  const user = ref(JSON.parse(localStorage.getItem('user') || 'null'))

  const isAuthenticated = computed(() => !!token.value)
  const role = computed(() => user.value?.role || null)
  const emailVerified = computed(() => user.value?.emailVerified === true)
  const needsOnboarding = computed(
    () => isAuthenticated.value && user.value?.onboardingCompleted === false
  )

  function setSession({ token: newToken, user: newUser }) {
    token.value = newToken
    user.value = newUser
    localStorage.setItem('token', newToken)
    localStorage.setItem('user', JSON.stringify(newUser))
  }

  function clearSession() {
    token.value = ''
    user.value = null
    localStorage.removeItem('token')
    localStorage.removeItem('user')
  }

  async function login(credentials) {
    const res = await authApi.login(credentials)
    setSession({
      token: res.token,
      user: {
        id: res.user_id,
        username: res.username,
        role: res.role,
        onboardingCompleted: res.onboarding_completed,
        emailVerified: false // refreshed by fetchMe on dashboard mount
      }
    })
    return res
  }

  async function register(payload) {
    return authApi.register(payload)
  }

  async function fetchMe() {
    const res = await authApi.me()
    user.value = {
      ...user.value,
      id: res.user_id,
      username: res.username,
      role: res.role,
      onboardingCompleted: res.onboarding_completed,
      emailVerified: res.email_verified === true
    }
    localStorage.setItem('user', JSON.stringify(user.value))
    return res
  }

  async function completeOnboarding(selectedRole) {
    const res = await authApi.setRole(selectedRole)
    setSession({
      token: res.token,
      user: {
        id: res.user_id,
        username: res.username,
        role: res.role,
        onboardingCompleted: res.onboarding_completed,
        emailVerified: user.value?.emailVerified === true
      }
    })
    return res
  }

  function logout() {
    clearSession()
  }

  return {
    token,
    user,
    role,
    emailVerified,
    isAuthenticated,
    needsOnboarding,
    login,
    register,
    fetchMe,
    completeOnboarding,
    logout
  }
})
