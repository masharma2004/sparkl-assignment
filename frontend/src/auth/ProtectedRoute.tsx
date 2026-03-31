import type { ReactNode } from 'react'
import { Navigate } from 'react-router'
import { LoadingBlock } from '../components/common/LoadingBlock'
import { useAuth } from './AuthContext'
import type { Role } from '../types/api'

interface ProtectedRouteProps {
  role: Role
  children: ReactNode
}

function dashboardPathForRole(role: Role) {
  return role === 'cms_admin' ? '/cms' : '/student'
}

function loginPathForRole(role: Role) {
  return role === 'cms_admin' ? '/login/cms' : '/login/student'
}

export function ProtectedRoute({ role, children }: ProtectedRouteProps) {
  const { session, isBootstrapping } = useAuth()

  if (isBootstrapping) {
    return (
      <div className="standalone-state">
        <LoadingBlock title="Restoring your workspace" message="Checking your session and preparing the right dashboard." />
      </div>
    )
  }

  if (!session) {
    return <Navigate to={loginPathForRole(role)} replace />
  }

  if (session.user.role !== role) {
    return <Navigate to={dashboardPathForRole(session.user.role)} replace />
  }

  return <>{children}</>
}

export function PublicOnlyRoute({ children }: { children: ReactNode }) {
  const { session, isBootstrapping } = useAuth()

  if (isBootstrapping) {
    return (
      <div className="standalone-state">
        <LoadingBlock title="Checking your session" message="If you already signed in, I’ll take you straight back to your dashboard." />
      </div>
    )
  }

  if (!session) {
    return <>{children}</>
  }

  return <Navigate to={dashboardPathForRole(session.user.role)} replace />
}

