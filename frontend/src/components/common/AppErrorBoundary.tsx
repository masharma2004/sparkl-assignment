import { Component, type ErrorInfo, type ReactNode } from 'react'

interface AppErrorBoundaryProps {
  children: ReactNode
}

interface AppErrorBoundaryState {
  hasError: boolean
  errorMessage: string | null
}

const isDevelopment = import.meta.env.DEV

export class AppErrorBoundary extends Component<
  AppErrorBoundaryProps,
  AppErrorBoundaryState
> {
  constructor(props: AppErrorBoundaryProps) {
    super(props)
    this.state = {
      hasError: false,
      errorMessage: null,
    }
  }

  static getDerivedStateFromError(error: Error): AppErrorBoundaryState {
    return {
      hasError: true,
      errorMessage: error.message,
    }
  }

  componentDidCatch(error: Error, errorInfo: ErrorInfo) {
    console.error('Unhandled application error', error, errorInfo)
  }

  handleReload = () => {
    window.location.reload()
  }

  render() {
    if (this.state.hasError) {
      return (
        <div className="standalone-state">
          <section className="empty-state error-boundary-state">
            <span className="eyebrow">Something broke</span>
            <h1>The interface hit an unexpected error</h1>
            <p>
              Reload the page to try again, or head back to the home screen and
              restart the flow from a clean state.
            </p>

            <div className="hero-actions">
              <button className="button primary" onClick={this.handleReload} type="button">
                Reload page
              </button>
              <a className="button secondary" href="/">
                Back home
              </a>
            </div>

            {isDevelopment && this.state.errorMessage ? (
              <details className="error-details">
                <summary>Developer details</summary>
                <pre>{this.state.errorMessage}</pre>
              </details>
            ) : null}
          </section>
        </div>
      )
    }

    return this.props.children
  }
}
