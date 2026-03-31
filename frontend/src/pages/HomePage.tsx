import { Link } from 'react-router'
import { useAuth } from '../auth/AuthContext'

export function HomePage() {
  const { session } = useAuth()
  const dashboardPath = session?.user.role === 'cms_admin' ? '/cms' : '/student'

  return (
    <div className="landing-page">
      <section className="landing-hero">
        <div className="page-header">
          <div className="landing-copy">
            <span className="eyebrow">Sparkl Edventure</span>
            <h1>Quiz workspace for admins and students.</h1>
            <p>
              The CMS portal creates quizzes and reviews participation. The student
              portal starts, resumes, finishes, and reports on quiz attempts.
            </p>
            <div className="hero-actions">
              <Link className="button primary" to="/login/cms">
                CMS Login
              </Link>
              <Link className="button secondary" to="/login/student">
                Student Login
              </Link>
              {!session ? (
                <Link className="button ghost" to="/signup/student">
                  Student Signup
                </Link>
              ) : null}
              {session ? (
                <Link className="button ghost" to={dashboardPath}>
                  Continue Session
                </Link>
              ) : null}
            </div>
          </div>
        </div>

        <section className="section-panel">
          <div className="landing-grid">
            <article className="feature-card">
              <span className="eyebrow">CMS</span>
              <h2>Create quizzes</h2>
              <p>Build assessments, assign marks, and review participant reports.</p>
            </article>
            <article className="feature-card">
              <span className="eyebrow">Student</span>
              <h2>Take quizzes</h2>
              <p>Start, resume, finish, and review a score from the same workspace.</p>
            </article>
            <article className="feature-card">
              <span className="eyebrow">System</span>
              <h2>Connected API flow</h2>
              <p>The React frontend is backed directly by the Go quiz service.</p>
            </article>
          </div>
        </section>
      </section>
    </div>
  )
}
