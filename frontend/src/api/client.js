const API_BASE = import.meta.env.VITE_API_BASE || 'http://localhost:8080/api'

export class ApiError extends Error {
  constructor(message, status) {
    super(message)
    this.status = status
  }
}

// Helper to read the CSRF token from the cookie
function getCsrfToken() {
  const match = document.cookie.match(new RegExp('(^| )csrf_token=([^;]+)'))
  if (match) return match[2]
  return ''
}

// Prevent multiple concurrent refreshes
let isRefreshing = false
let refreshPromise = null

async function refreshToken() {
  if (isRefreshing) return refreshPromise

  isRefreshing = true
  refreshPromise = fetch(`${API_BASE}/auth/refresh`, {
    method: 'POST',
    credentials: 'include'
  }).then((res) => {
    isRefreshing = false
    if (!res.ok) throw new Error('Refresh failed')
    return true
  }).catch((err) => {
    isRefreshing = false
    throw err
  })

  return refreshPromise
}

export async function apiFetch(path, options = {}) {
  const { auth = false, skipRefresh = false, ...rest } = options

  const headers = {
    'Content-Type': 'application/json',
    ...(rest.headers || {})
  }

  // Attach CSRF token for state-changing requests
  const method = (rest.method || 'GET').toUpperCase()
  if (auth && method !== 'GET' && method !== 'HEAD') {
    headers['X-CSRF-Token'] = getCsrfToken()
  }

  const config = {
    ...rest,
    headers,
    // THIS IS THE MAGIC: Browsers send HttpOnly cookies automatically if this is set
    credentials: 'include'
  }

  let res = await fetch(`${API_BASE}${path}`, config)

  // If 401 and auth is required, try to refresh and retry
  if (res.status === 401 && auth && !skipRefresh) {
    try {
      await refreshToken()
      // Retry the original request with the new cookie
      res = await fetch(`${API_BASE}${path}`, config)
    } catch {
      // If refresh fails, dispatch event so app can logout
      window.dispatchEvent(new CustomEvent('auth:unauthorized'))
      throw new ApiError('Session expired. Please log in again.', 401)
    }
  }

  if (res.status === 401 && auth) {
    window.dispatchEvent(new CustomEvent('auth:unauthorized'))
    throw new ApiError('Session expired. Please log in again.', 401)
  }

  let data = null
  const text = await res.text()
  if (text) {
    try { data = JSON.parse(text) } catch { data = { error: text } }
  }

  if (!res.ok) {
    const message = data?.error || `Request failed (${res.status})`
    throw new ApiError(message, res.status)
  }

  return data
}
