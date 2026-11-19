import { apiFetch } from './client'

const API_BASE = import.meta.env.VITE_API_BASE || 'http://localhost:8080/api'

// Helper to read the CSRF token from the cookie
function getCsrfToken() {
  const match = document.cookie.match(new RegExp('(^| )csrf_token=([^;]+)'))
  if (match) return match[2]
  return ''
}

async function uploadFile(path, file) {
  const form = new FormData()
  form.append('file', file)

  const headers = {
    'X-CSRF-Token': getCsrfToken()
  }

  const res = await fetch(`${API_BASE}${path}`, {
    method: 'POST',
    headers,
    body: form,
    credentials: 'include'  // Send HttpOnly cookies automatically
    // NOTE: We do NOT set Content-Type — the browser sets multipart/form-data
  })

  const text = await res.text()
  let data = null
  if (text) {
    try { data = JSON.parse(text) } catch { data = { error: text } }
  }

  if (!res.ok) {
    const err = new Error(data?.error || `Upload failed (${res.status})`)
    err.status = res.status
    throw err
  }
  return data
}

export function uploadResume(file) {
  return uploadFile('/seeker/profile/resume', file)
}

export function deleteResume() {
  return apiFetch('/seeker/profile/resume', { method: 'DELETE', auth: true })
}

export function uploadLogo(file) {
  return uploadFile('/employer/company/logo', file)
}

export function deleteLogo() {
  return apiFetch('/employer/company/logo', { method: 'DELETE', auth: true })
}
