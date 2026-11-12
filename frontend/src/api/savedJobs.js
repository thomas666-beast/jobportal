import { apiFetch } from './client'

export function saveJob(jobId) {
  return apiFetch(`/seeker/saved-jobs/${jobId}`, { method: 'POST', auth: true })
}

export function unsaveJob(jobId) {
  return apiFetch(`/seeker/saved-jobs/${jobId}`, { method: 'DELETE', auth: true })
}

export function checkSaved(jobId) {
  return apiFetch(`/seeker/saved-jobs/${jobId}`, { auth: true })
}

export function listSavedJobs() {
  return apiFetch('/seeker/saved-jobs', { auth: true })
}
