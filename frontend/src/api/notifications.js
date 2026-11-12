import { apiFetch } from './client'

export function listNotifications({ limit = 20, offset = 0 } = {}) {
  const params = new URLSearchParams({ limit: String(limit), offset: String(offset) })
  return apiFetch(`/notifications?${params.toString()}`, { auth: true })
}

export function listRecentNotifications() {
  return apiFetch('/notifications/recent', { auth: true })
}

export function getUnreadCount() {
  return apiFetch('/notifications/unread-count', { auth: true })
}

export function markRead(id) {
  return apiFetch(`/notifications/${id}/read`, { method: 'POST', auth: true })
}

export function markAllRead() {
  return apiFetch('/notifications/read-all', { method: 'POST', auth: true })
}
