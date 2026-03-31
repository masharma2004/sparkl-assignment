import { Link, isRouteErrorResponse, useRouteError } from 'react-router'

function getErrorMessage(error: unknown) {
  if (isRouteErrorResponse(error)) {
    return error.statusText || 'The route could not finish loading.'
  }

  if (error instanceof Error) {
    return error.message
  }

  return 'The route could not finish loading.'
}

export function RouteErrorPage() {
  const error = useRouteError()

  return (
    <div className="standalone-state">
      <section className="empty-state error-boundary-state">
        <span className="eyebrow">Route error</span>
        <h1>This page could not be rendered</h1>
        <p>{getErrorMessage(error)}</p>
        <div className="hero-actions">
          <Link className="button primary" to="/">
            Back home
          </Link>
          <button className="button secondary" onClick={() => window.location.reload()} type="button">
            Reload page
          </button>
        </div>
      </section>
    </div>
  )
}
