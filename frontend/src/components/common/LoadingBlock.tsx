interface LoadingBlockProps {
  title: string
  message?: string
}

export function LoadingBlock({ title, message }: LoadingBlockProps) {
  return (
    <section className="loading-block">
      <div className="loading-spinner" aria-hidden="true" />
      <div>
        <h2>{title}</h2>
        {message ? <p>{message}</p> : null}
      </div>
    </section>
  )
}
