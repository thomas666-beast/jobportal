import { apiFetch } from './client'

export function register({ username, email, password }) {
  return apiFetch('/register', {
    method: 'POST',
    body: JSON.stringify({ username, email, password })
  })
}

export function login({ email, password }) {
  return apiFetch('/login', {
    method: 'POST',
    body: JSON.stringify({ email, password })
  })
}

export function me() {
  return apiFetch('/me', { auth: true })
}

export function setRole(role) {
  return apiFetch('/onboarding/role', {
    method: 'POST',
    auth: true,
    body: JSON.stringify({ role })
  })
}
