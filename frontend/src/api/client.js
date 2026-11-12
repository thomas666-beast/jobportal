const API_BASE = import.meta.env.VITE_API_BASE || 'http://localhost:8080/api'

export class ApiError extends Error {
  constructor(message, status) {
    super(message)
    this.status = status
  }
}

export async function apiFetch(path, options = {}) {
  const { auth = false, ...rest } = options

  const headers = {
    'Content-Type': 'application/json',
    ...(rest.headers || {})
  }

  if (auth) {
    const token = localStorage.getItem('token')
    if (token) {
      headers['Authorization'] = `Bearer ${token}`
    }
  }

  const res = await fetch(`${API_BASE}${path}`, {
    ...rest,
    headers
  })

  let data = null
  const text = await res.text()
  if (text) {
    try {
      data = JSON.parse(text)
    } catch {
      data = { error: text }
    }
  }

  if (!res.ok) {
    const message = data?.error || `Request failed (${res.status})`
    throw new ApiError(message, res.status)
  }

  return data
}
