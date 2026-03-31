import { useEffect, useState } from 'react'
import { Link, useParams } from 'react-router'
import { getCmsParticipantReport } from '../api/cms'
import { ApiError } from '../api/client'
import { useAuth } from '../auth/AuthContext'
import { LoadingBlock } from '../components/common/LoadingBlock'
import { ReportView } from '../components/common/ReportView'
import type { AttemptReportResponse } from '../types/api'

export function CMSParticipantReportPage() {
  const { quizId, studentId } = useParams()
  const { session } = useAuth()
  const [report, setReport] = useState<AttemptReportResponse | null>(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    if (!session || !quizId || !studentId) {
      return
    }

    let isCancelled = false

    getCmsParticipantReport(quizId, studentId)
      .then((response) => {
        if (!isCancelled) {
          setReport(response)
        }
      })
      .catch((loadError) => {
        if (!isCancelled) {
          setError(
            loadError instanceof ApiError
              ? loadError.message
              : 'Failed to load participant report.',
          )
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
  }, [quizId, session, studentId])

  if (loading) {
    return (
      <LoadingBlock
        title="Loading participant report"
        message="Bringing in answers, scoring, and explanations."
      />
    )
  }

  if (!report) {
    return <p className="form-error">{error ?? 'Participant report not found.'}</p>
  }

  return (
    <div className="stack-lg">
      <Link className="button ghost inline-button" to={`/cms/quizzes/${quizId}/participants`}>
        Back to participants
      </Link>
      <ReportView report={report} />
    </div>
  )
}

