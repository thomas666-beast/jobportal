import { apiFetch } from './client'

export function listJobs() {
  return apiFetch('/employer/jobs', { auth: true })
}

export function getJob(id) {
  return apiFetch(`/employer/jobs/${id}`, { auth: true })
}

export function createJob(payload) {
  return apiFetch('/employer/jobs', {
    method: 'POST',
    auth: true,
    body: JSON.stringify(payload)
  })
}

export function updateJob(id, payload) {
  return apiFetch(`/employer/jobs/${id}`, {
    method: 'PUT',
    auth: true,
    body: JSON.stringify(payload)
  })
}

export function submitJob(id) {
  return apiFetch(`/employer/jobs/${id}/submit`, { method: 'POST', auth: true })
}

export function closeJob(id) {
  return apiFetch(`/employer/jobs/${id}/close`, { method: 'POST', auth: true })
}

export function archiveJob(id) {
  return apiFetch(`/employer/jobs/${id}/archive`, { method: 'POST', auth: true })
}
