import { apiFetch } from './client'

export function createSavedSearch(payload) {
  return apiFetch('/seeker/saved-searches', {
    method: 'POST',
    auth: true,
    body: JSON.stringify(payload)
  })
}

export function listSavedSearches() {
  return apiFetch('/seeker/saved-searches', { auth: true })
}

export function deleteSavedSearch(id) {
  return apiFetch(`/seeker/saved-searches/${id}`, { method: 'DELETE', auth: true })
}
