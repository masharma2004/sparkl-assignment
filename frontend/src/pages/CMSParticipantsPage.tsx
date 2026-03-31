import { useEffect, useState } from 'react'
import { Link, useParams } from 'react-router'
import { getCmsParticipants } from '../api/cms'
import { ApiError } from '../api/client'
import { useAuth } from '../auth/AuthContext'
import { LoadingBlock } from '../components/common/LoadingBlock'
import { StatusPill } from '../components/common/StatusPill'
import type { GetParticipantsResponse } from '../types/api'

export function CMSParticipantsPage() {
  const { quizId } = useParams()
  const { session } = useAuth()
  const [data, setData] = useState<GetParticipantsResponse | null>(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    if (!session || !quizId) {
      return
    }

    let isCancelled = false

    getCmsParticipants(quizId)
      .then((response) => {
        if (!isCancelled) {
          setData(response)
        }
      })
      .catch((loadError) => {
        if (!isCancelled) {
          setError(loadError instanceof ApiError ? loadError.message : 'Failed to load participants.')
        }
      })
      .finally(() => {
        if (!isCancelled) {
          setLoading(false)
        }
      })

    return () => {
      isCancelled = true
    }
  }, [quizId, session])

  if (loading) {
    return (
      <LoadingBlock
        title="Loading participants"
        message="Checking who has started, resumed, or completed this quiz."
      />
    )
  }

  if (!data) {
    return <p className="form-error">{error ?? 'Participants not found.'}</p>
  }

  return (
    <div className="stack-lg">
      <section className="page-header">
        <div>
          <span className="eyebrow">Participants</span>
          <h1>{data.quiz.title}</h1>
          <p>Track status and jump directly into any student report from this panel.</p>
        </div>
      </section>

      <section className="section-panel">
        <header className="section-heading">
          <div>
            <span className="eyebrow">Roster</span>
            <h2>{data.participants.length} students in view</h2>
          </div>
        </header>

        <div className="list-panel">
        {data.participants.map((participant) => (
          <article className="list-row" key={participant.student_id}>
            <div>
              <h2>{participant.full_name}</h2>
              <p>
                @{participant.username} · Student #{participant.student_id}
              </p>
            </div>
            <div className="list-row-actions">
              <StatusPill status={participant.status} />
              <strong>{participant.score} marks</strong>
              {participant.attempt_id ? (
                <Link
                  className="button ghost"
                  to={`/cms/quizzes/${data.quiz.quiz_id}/participants/${participant.student_id}/report`}
                >
                  Open report
                </Link>
              ) : (
                <span className="subtle-text">No attempt yet</span>
              )}
            </div>
          </article>
        ))}
        </div>
      </section>
    </div>
  )
}

