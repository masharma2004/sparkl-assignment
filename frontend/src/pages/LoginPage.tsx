import { useState } from 'react'
import { Link, useNavigate } from 'react-router'
import { ApiError } from '../api/client'
import { useAuth } from '../auth/AuthContext'

interface PortalAuthPageProps {
  portal: 'cms' | 'student'
  mode?: 'login' | 'signup'
}

const portalContent = {
  cms: {
    login: {
      title: 'Welcome Back to Sparkl',
      description:
        'Sign in to the CMS workspace to create quizzes, review participants, and open reports.',
      accent: 'CMS Login',
      submitLabel: 'Login',
      helperLinks: [
        { href: '/login/student', label: 'Use the student portal instead' },
      ],
      demoLabel: 'Demo credentials',
      demoValue: 'cmsadmin / password123',
    },
  },
  student: {
    login: {
      title: 'Welcome Back to Sparkl',
      description:
        'Sign in to start a quiz, resume your last attempt, or view a completed scorecard.',
      accent: 'Student Login',
      submitLabel: 'Login',
      helperLinks: [
        { href: '/signup/student', label: 'Need an account? Create one' },
        { href: '/login/cms', label: 'Looking for CMS login?' },
      ],
      demoLabel: 'Demo credentials',
      demoValue: 'student1 / password123',
    },
    signup: {
      title: 'Create your Sparkl account',
      description:
        'Register as a student to access the dashboard, start a quiz, and review your report later.',
      accent: 'Student Signup',
      submitLabel: 'Create account',
      helperLinks: [
        { href: '/login/student', label: 'Already have an account? Sign in' },
        { href: '/login/cms', label: 'Looking for CMS login?' },
      ],
      demoLabel: 'What happens next',
      demoValue: 'We sign you in immediately after the account is created.',
    },
  },
} as const

function PortalAuthPage({ portal, mode = 'login' }: PortalAuthPageProps) {
  const isStudentSignup = portal === 'student' && mode === 'signup'
  const content =
    portal === 'cms'
      ? portalContent.cms.login
      : isStudentSignup
        ? portalContent.student.signup
        : portalContent.student.login
  const navigate = useNavigate()
  const { loginCms, loginStudent, signupStudent } = useAuth()
  const [username, setUsername] = useState('')
  const [password, setPassword] = useState('')
  const [fullName, setFullName] = useState('')
  const [email, setEmail] = useState('')
  const [error, setError] = useState<string | null>(null)
  const [isSubmitting, setIsSubmitting] = useState(false)

  async function handleSubmit(event: React.FormEvent<HTMLFormElement>) {
    event.preventDefault()
    setError(null)
    setIsSubmitting(true)

    try {
      if (portal === 'cms') {
        await loginCms({ username, password })
        navigate('/cms', { replace: true })
      } else if (isStudentSignup) {
        await signupStudent({
          username,
          password,
          full_name: fullName,
          email,
        })
        navigate('/student', { replace: true })
      } else {
        await loginStudent({ username, password })
        navigate('/student', { replace: true })
      }
    } catch (submitError) {
      if (submitError instanceof ApiError) {
        setError(submitError.message)
      } else {
        setError(isStudentSignup ? 'Unable to create your account right now.' : 'Unable to sign in right now.')
      }
    } finally {
      setIsSubmitting(false)
    }
  }

  return (
    <div className="auth-page">
      <section className="auth-panel">
        <Link className="auth-back-link" to="/">
          <span aria-hidden="true">←</span>
        </Link>

        <div className="auth-copy">
          <span className="eyebrow">{content.accent}</span>
          <h1>{content.title}</h1>
          <p>{content.description}</p>
          <div className="auth-demo">
            <strong>{content.demoLabel}</strong>
            <p>{content.demoValue}</p>
          </div>
        </div>

        <form className="auth-form" onSubmit={handleSubmit}>
          {isStudentSignup ? (
            <>
              <label className="field">
                <span>Full name</span>
                <input
                  autoComplete="name"
                  onChange={(event) => setFullName(event.target.value)}
                  placeholder="Enter your full name"
                  required
                  value={fullName}
                />
              </label>

              <label className="field">
                <span>Email</span>
                <input
                  autoComplete="email"
                  onChange={(event) => setEmail(event.target.value)}
                  placeholder="Enter your email"
                  required
                  type="email"
                  value={email}
                />
              </label>
            </>
          ) : null}

          <label className="field">
            <span>Username</span>
            <input
              autoComplete="username"
              onChange={(event) => setUsername(event.target.value)}
              placeholder="Enter your username"
              required
              value={username}
            />
          </label>

          <label className="field">
            <span>Password</span>
            <input
              autoComplete={isStudentSignup ? 'new-password' : 'current-password'}
              minLength={8}
              onChange={(event) => setPassword(event.target.value)}
              placeholder={isStudentSignup ? 'Create a password' : 'Enter your password'}
              required
              type="password"
              value={password}
            />
          </label>

          {error ? <p className="form-error">{error}</p> : null}

          <button className="button primary" disabled={isSubmitting} type="submit">
            {isSubmitting
              ? isStudentSignup
                ? 'Creating account...'
                : 'Signing in...'
              : content.submitLabel}
          </button>

          <div className="auth-links">
            {content.helperLinks.map((link) => (
              <Link className="auth-secondary-link" key={link.href} to={link.href}>
                {link.label}
              </Link>
            ))}
          </div>
        </form>
      </section>
    </div>
  )
}

export function CMSLoginPage() {
  return <PortalAuthPage portal="cms" />
}

export function StudentLoginPage() {
  return <PortalAuthPage portal="student" mode="login" />
}

export function StudentSignupPage() {
  return <PortalAuthPage portal="student" mode="signup" />
}

