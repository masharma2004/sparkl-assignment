import { Link, NavLink, Outlet } from 'react-router'
import { useAuth } from '../../auth/AuthContext'

interface PortalLayoutProps {
  portal: 'cms' | 'student'
}

const portalConfig = {
  cms: {
    eyebrow: 'CMS Workspace',
    title: 'Sparkl Edventure',
    nav: [
      { label: 'Dashboard', to: '/cms' },
      { label: 'Question Bank', to: '/cms/questions' },
      { label: 'Create Quiz', to: '/cms/quizzes/new' },
    ],
  },
  student: {
    eyebrow: 'Student Workspace',
    title: 'Sparkl Edventure',
    nav: [{ label: 'Dashboard', to: '/student' }],
  },
} as const

export function PortalLayout({ portal }: PortalLayoutProps) {
  const { session, logout } = useAuth()
  const config = portalConfig[portal]
  const firstName = session?.user.full_name.split(' ')[0] ?? 'User'
  const greeting = portal === 'cms' ? 'Hi, Admin' : `Hi, ${firstName}`

  return (
    <div className="portal-shell">
      <header className="portal-header">
        <Link className="portal-brand" to="/">
          <span className="brand-mark">SE</span>
          <div>
            <span className="eyebrow">{config.eyebrow}</span>
            <strong>{config.title}</strong>
          </div>
        </Link>

        <nav className="portal-nav">
          {config.nav.map((item) => (
            <NavLink
              className={({ isActive }) =>
                isActive ? 'portal-tab portal-tab--active' : 'portal-tab'
              }
              end={item.to === '/cms' || item.to === '/student'}
              key={item.to}
              to={item.to}
            >
              {item.label}
            </NavLink>
          ))}
        </nav>

        <div className="portal-user">
          <span className="portal-user-label">{greeting}</span>
          <button className="button ghost" onClick={logout} type="button">
            Logout
          </button>
        </div>
      </header>

      <main className="page-content portal-main">
        <Outlet />
      </main>
    </div>
  )
}
