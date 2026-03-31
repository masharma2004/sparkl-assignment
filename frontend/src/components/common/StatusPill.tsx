import { formatStatusLabel } from '../../utils/format'

interface StatusPillProps {
  status: string
}

const statusClassByValue: Record<string, string> = {
  cms_admin: 'status-pill--accent',
  student: 'status-pill--muted',
  not_started: 'status-pill--muted',
  in_progress: 'status-pill--warning',
  completed: 'status-pill--success',
  resume: 'status-pill--warning',
  start: 'status-pill--accent',
  view_score: 'status-pill--success',
}

export function StatusPill({ status }: StatusPillProps) {
  const className = statusClassByValue[status] ?? 'status-pill--muted'

  return (
    <span className={`status-pill ${className}`}>
      {formatStatusLabel(status)}
    </span>
  )
}
