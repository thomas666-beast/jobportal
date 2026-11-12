import { apiFetch } from './client'

export function listCompanyApplications() {
  return apiFetch('/employer/applications', { auth: true })
}

export function listJobApplications(jobId) {
  return apiFetch(`/employer/jobs/${jobId}/applications`, { auth: true })
}

export function getApplicationDetail(id) {
  return apiFetch(`/employer/applications/${id}`, { auth: true })
}

export function updateApplicationStatus(id, status) {
  return apiFetch(`/employer/applications/${id}/status`, {
    method: 'POST',
    auth: true,
    body: JSON.stringify({ status })
  })
}
