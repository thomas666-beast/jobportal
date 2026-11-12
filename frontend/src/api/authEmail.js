import { apiFetch } from './client'

export function verifyEmail(token) {
  return apiFetch('/auth/verify-email', {
    method: 'POST',
    body: JSON.stringify({ token })
  })
}

export function resendVerification() {
  return apiFetch('/auth/resend-verification', { method: 'POST', auth: true })
}

export function forgotPassword(email) {
  return apiFetch('/auth/forgot-password', {
    method: 'POST',
    body: JSON.stringify({ email })
  })
}

export function resetPassword(token, newPassword) {
  return apiFetch('/auth/reset-password', {
    method: 'POST',
    body: JSON.stringify({ token, new_password: newPassword })
  })
}
