import { apiFetch } from './client'

export function getProfile() {
  return apiFetch('/seeker/profile', { auth: true })
}

export function saveProfile(payload) {
  return apiFetch('/seeker/profile', {
    method: 'PUT',
    auth: true,
    body: JSON.stringify(payload)
  })
}
