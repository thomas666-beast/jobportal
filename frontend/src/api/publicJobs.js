import { apiFetch } from './client'

export function browseJobs({ q = '', location = '', remote_type = '', limit = 20, offset = 0 } = {}) {
  const params = new URLSearchParams()
  if (q) params.set('q', q)
  if (location) params.set('location', location)
  if (remote_type) params.set('remote_type', remote_type)
  params.set('limit', String(limit))
  params.set('offset', String(offset))
  return apiFetch(`/jobs?${params.toString()}`)
}

export function getPublicJob(id) {
  return apiFetch(`/jobs/${id}`)
}
