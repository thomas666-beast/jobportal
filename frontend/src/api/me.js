import { apiFetch } from './client'

export function changePassword({ current_password, new_password }) {
  return apiFetch('/me/password', {
    method: 'POST',
    auth: true,
    body: JSON.stringify({ current_password, new_password })
  })
}
