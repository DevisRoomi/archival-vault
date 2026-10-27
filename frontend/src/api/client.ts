import type { ApiResponse } from '../types'

const getApiBase = (): string => {
  if (typeof window !== 'undefined') {
    // If running on Vite dev server (port 5173), Vite proxy routes /api to backend
    // If running on Herd (http://archival-vault.test), Herd /api or local backend
    if (window.location.port === '5173') {
      return '/api'
    }
    // If served via Herd or direct file, use relative /api or direct 8085
    return '/api'
  }
  return 'http://127.0.0.1:8085/api'
}

export const request = async <T = any>(
  endpoint: string,
  options: RequestInit = {}
): Promise<T> => {
  const token = localStorage.getItem('vault_token')
  const headers = new Headers(options.headers || {})

  if (!headers.has('Content-Type') && !(options.body instanceof FormData)) {
    headers.set('Content-Type', 'application/json')
  }

  if (token && !headers.has('Authorization')) {
    headers.set('Authorization', `Bearer ${token}`)
  }

  const base = getApiBase()
  const cleanEndpoint = endpoint.startsWith('/') ? endpoint : `/${endpoint}`
  let url = `${base}${cleanEndpoint}`

  let res: Response
  try {
    res = await fetch(url, { ...options, headers })
  } catch (err) {
    // If relative /api failed (e.g. running Herd without proxy), fallback directly to 8085
    if (base === '/api') {
      url = `http://127.0.0.1:8085/api${cleanEndpoint}`
      res = await fetch(url, { ...options, headers })
    } else {
      throw err
    }
  }

  const data: ApiResponse<T> = await res.json().catch(() => ({
    success: false,
    error: `HTTP ${res.status}: Failed to parse JSON response`
  }))

  if (!res.ok || !data.success) {
    if (res.status === 401) {
      localStorage.removeItem('vault_token')
      localStorage.removeItem('vault_user')
      window.dispatchEvent(new CustomEvent('auth:expired'))
    }
    throw new Error(data.error || data.message || `Request failed with status ${res.status}`)
  }

  return data.data as T
}
