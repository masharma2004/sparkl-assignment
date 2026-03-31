import { Link } from 'react-router'

export function NotFoundPage() {
  return (
    <div className="standalone-state">
      <section className="empty-state">
        <span className="eyebrow">404</span>
        <h1>That page is off the map</h1>
        <p>Use the links below to get back to the right workspace.</p>
        <div className="hero-actions">
          <Link className="button primary" to="/">
            Back home
          </Link>
          <Link className="button secondary" to="/student">
            Student dashboard
          </Link>
          <Link className="button ghost" to="/cms">
            CMS dashboard
          </Link>
        </div>
      </section>
    </div>
  )
}

