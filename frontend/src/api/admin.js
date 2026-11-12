import { apiFetch } from './client'

// Jobs moderation
export function listAdminJobs({ status = 'pending_review', limit = 50, offset = 0 } = {}) {
  const params = new URLSearchParams({ status, limit: String(limit), offset: String(offset) })
  return apiFetch(`/admin/jobs?${params.toString()}`, { auth: true })
}

export function getAdminJob(id) {
  return apiFetch(`/admin/jobs/${id}`, { auth: true })
}

export function approveJob(id) {
  return apiFetch(`/admin/jobs/${id}/approve`, { method: 'POST', auth: true })
}

export function rejectJob(id, reason) {
  return apiFetch(`/admin/jobs/${id}/reject`, {
    method: 'POST',
    auth: true,
    body: JSON.stringify({ reason })
  })
}

// Stats
export function getAdminStats() {
  return apiFetch('/admin/stats', { auth: true })
}

// Users
export function listAdminUsers({ role = '', limit = 50, offset = 0 } = {}) {
  const params = new URLSearchParams({ limit: String(limit), offset: String(offset) })
  if (role) params.set('role', role)
  return apiFetch(`/admin/users?${params.toString()}`, { auth: true })
}

export function getAdminUser(id) {
  return apiFetch(`/admin/users/${id}`, { auth: true })
}

export function listAdminApplications({ status = '', job_id = '', seeker_q = '', limit = 50, offset = 0 } = {}) {
  const params = new URLSearchParams({ limit: String(limit), offset: String(offset) })
  if (status) params.set('status', status)
  if (job_id) params.set('job_id', job_id)
  if (seeker_q) params.set('seeker_q', seeker_q)
  return apiFetch(`/admin/applications?${params.toString()}`, { auth: true })
}

export function getAdminApplication(id) {
  return apiFetch(`/admin/applications/${id}`, { auth: true })
}
