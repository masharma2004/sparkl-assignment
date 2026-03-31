import { useEffect, useState } from 'react'
import { Link, useParams } from 'react-router'
import { getCmsQuiz } from '../api/cms'
import { ApiError } from '../api/client'
import { useAuth } from '../auth/AuthContext'
import { LoadingBlock } from '../components/common/LoadingBlock'
import { normalizeCategory } from '../utils/category'
import type { GetQuizResponse } from '../types/api'
import { formatDate } from '../utils/format'

export function CMSQuizDetailPage() {
  const { quizId } = useParams()
  const { session } = useAuth()
  const [data, setData] = useState<GetQuizResponse | null>(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    if (!session || !quizId) {
      return
    }

    let isCancelled = false

    getCmsQuiz(quizId)
      .then((response) => {
        if (!isCancelled) {
          setData(response)
        }
      })
      .catch((loadError) => {
        if (!isCancelled) {
          setError(loadError instanceof ApiError ? loadError.message : 'Failed to load quiz.')
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
        title="Loading quiz details"
        message="Bringing in the quiz structure and question order."
      />
    )
  }

  if (!data) {
    return <p className="form-error">{error ?? 'Quiz not found.'}</p>
  }

  return (
    <div className="stack-lg">
      <section className="page-header">
        <div>
          <span className="eyebrow">Quiz Overview</span>
          <h1>{data.quiz.title}</h1>
          <p className="subtle-text">Category: {normalizeCategory(data.quiz.category)}</p>
          <p>
            {data.quiz.question_count} questions · {data.quiz.total_marks} marks ·{' '}
            {data.quiz.duration_minutes} minutes
          </p>
        </div>
        <div className="page-header-actions">
          <Link className="button secondary" to={`/cms/quizzes/${data.quiz.quiz_id}/participants`}>
            View participants
          </Link>
        </div>
      </section>

      <section className="metric-grid">
        <article className="metric-card">
          <span>Created</span>
          <strong>{formatDate(data.quiz.created_at)}</strong>
        </article>
        <article className="metric-card">
          <span>Last updated</span>
          <strong>{formatDate(data.quiz.updated_at)}</strong>
        </article>
        <article className="metric-card">
          <span>Question order</span>
          <strong>{data.questions.length}</strong>
        </article>
      </section>

      <section className="section-panel stack-md">
        <header className="section-heading">
          <div>
            <span className="eyebrow">Question order</span>
            <h2>{data.questions.length} questions in this quiz</h2>
          </div>
        </header>

        {data.questions.map((question) => (
          <article className="resource-card" key={question.quiz_question_id}>
            <div className="resource-card-header">
              <span className="eyebrow">Question {question.sequence_number}</span>
              <strong>{question.marks} marks</strong>
            </div>
            <h2>{question.prompt}</h2>
            <p className="subtle-text">Question ID: {question.question_id}</p>
          </article>
        ))}
      </section>
    </div>
  )
}

