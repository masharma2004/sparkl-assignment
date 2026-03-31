type HttpMethod = 'GET' | 'POST' | 'PATCH' | 'PUT' | 'DELETE'

interface ApiRequestOptions {
  method?: HttpMethod
  body?: unknown
  signal?: AbortSignal
  retryOnUnauthorized?: boolean
}

const AUTH_REFRESHED_EVENT = 'sparklassignment:auth-refreshed'
const AUTH_EXPIRED_EVENT = 'sparklassignment:auth-expired'

const configuredBaseUrl = import.meta.env.VITE_API_BASE_URL?.trim()
const API_BASE_URL =
  configuredBaseUrl && configuredBaseUrl.length > 0
    ? configuredBaseUrl.replace(/\/$/, '')
    : 'http://localhost:8080/api/v1'

export class ApiError extends Error {
  status: number
  data: unknown

  constructor(message: string, status: number, data: unknown) {
    super(message)
    this.name = 'ApiError'
    this.status = status
    this.data = data
  }
}

function extractErrorMessage(data: unknown, fallback: string): string {
  if (typeof data === 'string' && data.trim().length > 0) {
    return data
  }

  if (typeof data === 'object' && data !== null) {
    const record = data as Record<string, unknown>
    if (typeof record.error === 'string' && record.error.trim().length > 0) {
      return record.error
    }
    if (typeof record.message === 'string' && record.message.trim().length > 0) {
      return record.message
    }
  }

  return fallback
}

function buildUrl(path: string): string {
  if (/^https?:\/\//i.test(path)) {
    return path
  }

  const normalizedPath = path.startsWith('/') ? path : `/${path}`
  return `${API_BASE_URL}${normalizedPath}`
}

async function parseResponseBody(response: Response): Promise<unknown> {
  const text = await response.text()
  if (!text) {
    return null
  }

  try {
    return JSON.parse(text) as unknown
  } catch {
    return text
  }
}

export async function apiRequest<T>(
  path: string,
  { method = 'GET', body, signal, retryOnUnauthorized = true }: ApiRequestOptions = {},
): Promise<T> {
  const headers = new Headers()

  if (body !== undefined) {
    headers.set('Content-Type', 'application/json')
  }

  const requestUrl = buildUrl(path)
  const requestInit: RequestInit = {
    method,
    headers,
    body: body !== undefined ? JSON.stringify(body) : undefined,
    credentials: 'include',
    signal,
  }

  const response = await fetch(requestUrl, requestInit)

  const data = await parseResponseBody(response)

  if (
    response.status === 401 &&
    retryOnUnauthorized &&
    shouldAttemptRefresh(path)
  ) {
    const refreshResponse = await refreshSession()
    if (refreshResponse) {
      return apiRequest<T>(path, {
        method,
        body,
        signal,
        retryOnUnauthorized: false,
      })
    }
  }

  if (!response.ok) {
    throw new ApiError(
      extractErrorMessage(data, 'Request failed'),
      response.status,
      data,
    )
  }

  return data as T
}

let refreshRequest: Promise<boolean> | null = null

function shouldAttemptRefresh(path: string): boolean {
  return ![
    '/auth/cms/login',
    '/auth/student/login',
    '/auth/student/signup',
    '/auth/refresh',
    '/auth/logout',
  ].some((authPath) => path.startsWith(authPath))
}

async function refreshSession(): Promise<boolean> {
  if (refreshRequest) {
    return refreshRequest
  }

  refreshRequest = (async () => {
    const response = await fetch(buildUrl('/auth/refresh'), {
      method: 'POST',
      credentials: 'include',
    })

    const data = await parseResponseBody(response)
    if (!response.ok) {
      if (typeof window !== 'undefined') {
        window.dispatchEvent(new Event(AUTH_EXPIRED_EVENT))
      }
      return false
    }

    if (typeof window !== 'undefined') {
      window.dispatchEvent(
        new CustomEvent(AUTH_REFRESHED_EVENT, {
          detail: data,
        }),
      )
    }

    return true
  })()

  try {
    return await refreshRequest
  } finally {
    refreshRequest = null
  }
}

export { AUTH_EXPIRED_EVENT, AUTH_REFRESHED_EVENT }
