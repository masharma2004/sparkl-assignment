import {
  createContext,
  useContext,
  useEffect,
  useMemo,
  useState,
  type PropsWithChildren,
} from 'react'
import {
  AUTH_EXPIRED_EVENT,
  AUTH_REFRESHED_EVENT,
} from '../api/client'
import {
  getMe,
  loginCms as loginCmsRequest,
  loginStudent as loginStudentRequest,
  logout as logoutRequest,
  signupStudent as signupStudentRequest,
} from '../api/auth'
import type { LoginRequest, LoginResponse, StudentSignupRequest } from '../types/api'

interface AuthContextValue {
  session: LoginResponse | null
  isBootstrapping: boolean
  loginCms: (payload: LoginRequest) => Promise<LoginResponse>
  loginStudent: (payload: LoginRequest) => Promise<LoginResponse>
  signupStudent: (payload: StudentSignupRequest) => Promise<LoginResponse>
  logout: () => void
}

const AuthContext = createContext<AuthContextValue | undefined>(undefined)
const STORAGE_KEY = 'sparklassignment.session'

function isValidSession(value: unknown): value is LoginResponse {
  if (!value || typeof value !== 'object') {
    return false
  }

  const candidate = value as Record<string, unknown>
  const user = candidate.user

  if (!user || typeof user !== 'object') {
    return false
  }

  const sessionUser = user as Record<string, unknown>
  return typeof sessionUser.role === 'string'
}

function readSession(): LoginResponse | null {
  if (typeof window === 'undefined') {
    return null
  }

  const rawValue =
    window.sessionStorage.getItem(STORAGE_KEY) ??
    window.localStorage.getItem(STORAGE_KEY)
  if (!rawValue) {
    return null
  }

  try {
    const parsed = JSON.parse(rawValue)
    if (isValidSession(parsed)) {
      window.sessionStorage.setItem(STORAGE_KEY, rawValue)
      window.localStorage.removeItem(STORAGE_KEY)
      return parsed
    }
  } catch {
    window.sessionStorage.removeItem(STORAGE_KEY)
    window.localStorage.removeItem(STORAGE_KEY)
    return null
  }

  window.sessionStorage.removeItem(STORAGE_KEY)
  window.localStorage.removeItem(STORAGE_KEY)
  return null
}

function persistSession(session: LoginResponse | null) {
  if (typeof window === 'undefined') {
    return
  }

  if (!session) {
    window.sessionStorage.removeItem(STORAGE_KEY)
    window.localStorage.removeItem(STORAGE_KEY)
    return
  }

  const serialized = JSON.stringify(session)
  window.sessionStorage.setItem(STORAGE_KEY, serialized)
  window.localStorage.removeItem(STORAGE_KEY)
}

export function AuthProvider({ children }: PropsWithChildren) {
  const [session, setSession] = useState<LoginResponse | null>(() => readSession())
  const [isBootstrapping, setIsBootstrapping] = useState(true)

  useEffect(() => {
    let isCancelled = false

    getMe()
      .then((response) => {
        if (isCancelled) {
          return
        }

        const nextSession: LoginResponse = {
          user: response.user,
        }

        setSession(nextSession)
        persistSession(nextSession)
      })
      .catch(() => {
        if (isCancelled) {
          return
        }

        setSession(null)
        persistSession(null)
      })
      .finally(() => {
        if (!isCancelled) {
          setIsBootstrapping(false)
        }
      })

    return () => {
      isCancelled = true
    }
  }, [])

  useEffect(() => {
    if (typeof window === 'undefined') {
      return
    }

    const handleRefreshed = (event: Event) => {
      const detail = (event as CustomEvent<LoginResponse>).detail
      if (!detail?.user) {
        return
      }

      setSession(detail)
      persistSession(detail)
    }

    const handleExpired = () => {
      setSession(null)
      persistSession(null)
    }

    window.addEventListener(AUTH_REFRESHED_EVENT, handleRefreshed)
    window.addEventListener(AUTH_EXPIRED_EVENT, handleExpired)

    return () => {
      window.removeEventListener(AUTH_REFRESHED_EVENT, handleRefreshed)
      window.removeEventListener(AUTH_EXPIRED_EVENT, handleExpired)
    }
  }, [])

  const value = useMemo<AuthContextValue>(
    () => ({
      session,
      isBootstrapping,
      async loginCms(payload) {
        const response = await loginCmsRequest(payload)
        setSession(response)
        persistSession(response)
        return response
      },
      async loginStudent(payload) {
        const response = await loginStudentRequest(payload)
        setSession(response)
        persistSession(response)
        return response
      },
      async signupStudent(payload) {
        const response = await signupStudentRequest(payload)
        setSession(response)
        persistSession(response)
        return response
      },
      logout() {
        void logoutRequest().catch(() => undefined)
        setSession(null)
        persistSession(null)
      },
    }),
    [isBootstrapping, session],
  )

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>
}

export function useAuth() {
  const context = useContext(AuthContext)
  if (!context) {
    throw new Error('useAuth must be used inside AuthProvider')
  }

  return context
}
