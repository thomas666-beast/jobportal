import { apiFetch } from './client'

export function createApplication({ job_id, cover_letter = '' }) {
  return apiFetch('/seeker/applications', {
    method: 'POST',
    auth: true,
    body: JSON.stringify({ job_id, cover_letter })
  })
}

export function listMyApplications() {
  return apiFetch('/seeker/applications', { auth: true })
}

export function checkApplicationStatus(job_id) {
  return apiFetch(`/seeker/applications/status?job_id=${job_id}`, { auth: true })
}

export function withdrawApplication(id) {
  return apiFetch(`/seeker/applications/${id}/withdraw`, { method: 'POST', auth: true })
}
