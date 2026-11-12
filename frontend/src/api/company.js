import { apiFetch } from './client'

export function getCompany() {
  return apiFetch('/employer/company', { auth: true })
}

export function saveCompany(payload) {
  return apiFetch('/employer/company', {
    method: 'PUT',
    auth: true,
    body: JSON.stringify(payload)
  })
}
