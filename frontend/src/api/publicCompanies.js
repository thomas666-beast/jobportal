import { apiFetch } from './client'

export function listCompanies({ q = '', limit = 20, offset = 0 } = {}) {
  const params = new URLSearchParams()
  if (q) params.set('q', q)
  params.set('limit', String(limit))
  params.set('offset', String(offset))
  return apiFetch(`/companies?${params.toString()}`)
}

export function getPublicCompany(id) {
  return apiFetch(`/companies/${id}`)
}
