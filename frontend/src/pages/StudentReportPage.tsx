import { useEffect, useState } from 'react'
import { Link, useNavigate, useParams } from 'react-router'
import { getStudentReport } from '../api/student'
import { ApiError } from '../api/client'
import { useAuth } from '../auth/AuthContext'
import { LoadingBlock } from '../components/common/LoadingBlock'
import { ReportView } from '../components/common/ReportView'
import type { AttemptReportResponse } from '../types/api'

export function StudentReportPage() {
  const { attemptId } = useParams()
  const { session } = useAuth()
  const navigate = useNavigate()
  const [report, setReport] = useState<AttemptReportResponse | null>(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    if (!session || !attemptId) {
      return
    }

    let isCancelled = false

    getStudentReport(attemptId)
      .then((response) => {
        if (!isCancelled) {
          setReport(response)
        }
      })
      .catch((loadError) => {
        if (isCancelled) {
          return
        }

        if (loadError instanceof ApiError && loadError.status === 409) {
          navigate(`/student/attempts/${attemptId}`, { replace: true })
          return
        }

        setError(loadError instanceof ApiError ? loadError.message : 'Failed to load report.')
      })
      .finally(() => {
        if (!isCancelled) {
          setLoading(false)
        }
      })

    return () => {
      isCancelled = true
    }
  }, [attemptId, navigate, session])

  if (loading) {
    return (
      <LoadingBlock
        title="Loading your report"
        message="Scoring, solutions, and explanations are on the way."
      />
    )
  }

  if (!report) {
    return <p className="form-error">{error ?? 'Report not found.'}</p>
  }

  return (
    <div className="stack-lg">
      <Link className="button ghost inline-button" to="/student">
        Back to dashboard
      </Link>
      <ReportView report={report} />
    </div>
  )
}

