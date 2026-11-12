import { apiFetch } from './client'

const API_BASE = import.meta.env.VITE_API_BASE || 'http://localhost:8080/api'

async function uploadFile(path, file) {
  const token = localStorage.getItem('token')
  const form = new FormData()
  form.append('file', file)

  const res = await fetch(`${API_BASE}${path}`, {
    method: 'POST',
    headers: token ? { Authorization: `Bearer ${token}` } : {},
    body: form
  })

  const text = await res.text()
  let data = null
  if (text) {
    try {
      data = JSON.parse(text)
    } catch {
      data = { error: text }
    }
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
