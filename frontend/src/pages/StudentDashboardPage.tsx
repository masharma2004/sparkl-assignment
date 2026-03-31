import { useEffect, useMemo, useState } from 'react'
import { useNavigate } from 'react-router'
import { getStudentQuizzes, startStudentQuiz } from '../api/student'
import { ApiError } from '../api/client'
import { useAuth } from '../auth/AuthContext'
import { LoadingBlock } from '../components/common/LoadingBlock'
import { StatusPill } from '../components/common/StatusPill'
import type { StudentQuizItemResponse } from '../types/api'
import { normalizeCategory, matchesCategoryFilter } from '../utils/category'

function extractAttemptId(error: ApiError): number | null {
  if (typeof error.data !== 'object' || error.data === null) {
    return null
  }

  const record = error.data as Record<string, unknown>
  return typeof record.attempt_id === 'number' ? record.attempt_id : null
}

export function StudentDashboardPage() {
  const { session } = useAuth()
  const navigate = useNavigate()
  const [quizzes, setQuizzes] = useState<StudentQuizItemResponse[]>([])
  const [categoryFilter, setCategoryFilter] = useState('all')
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [busyQuizId, setBusyQuizId] = useState<number | null>(null)

  useEffect(() => {
    if (!session) {
      return
    }

    let isCancelled = false

    getStudentQuizzes()
      .then((response) => {
        if (!isCancelled) {
          setQuizzes(response.quizzes)
        }
      })
      .catch((loadError) => {
        if (!isCancelled) {
          setError(loadError instanceof ApiError ? loadError.message : 'Failed to load quizzes.')
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
  }, [session])

  const availableCategories = useMemo(
    () =>
      Array.from(
        new Set(quizzes.map((quiz) => normalizeCategory(quiz.category))),
      ).sort((left, right) => left.localeCompare(right)),
    [quizzes],
  )

  const filteredQuizzes = useMemo(() => {
    return quizzes.filter((quiz) => matchesCategoryFilter(quiz.category, categoryFilter))
  }, [categoryFilter, quizzes])

  async function handleAction(quiz: StudentQuizItemResponse) {
    if (!session) {
      return
    }

    if (quiz.action === 'resume' && quiz.attempt_id) {
      navigate(`/student/attempts/${quiz.attempt_id}`)
      return
    }

    if (quiz.action === 'view_score' && quiz.attempt_id) {
      navigate(`/student/attempts/${quiz.attempt_id}/report`)
      return
    }

    setBusyQuizId(quiz.id)
    setError(null)

    try {
      const response = await startStudentQuiz(quiz.id)
      navigate(`/student/attempts/${response.attempt_id}`)
    } catch (startError) {
      if (startError instanceof ApiError && startError.status === 409) {
        const attemptId = extractAttemptId(startError)
        if (attemptId) {
          navigate(`/student/attempts/${attemptId}/report`)
          return
        }
      }

      setError(startError instanceof ApiError ? startError.message : 'Failed to start quiz.')
    } finally {
      setBusyQuizId(null)
    }
  }

  if (loading) {
    return <LoadingBlock title="Loading your quiz desk" message="Checking what you can start, resume, or review." />
  }

  return (
    <div className="stack-lg">
      <section className="page-header">
        <div>
          <span className="eyebrow">Student Dashboard</span>
          <h1>Your quiz desk</h1>
          <p>Start a new quiz, resume an in-progress attempt, or open a completed score.</p>
        </div>
      </section>

      {error ? <p className="form-error">{error}</p> : null}

      <section className="section-panel stack-md">
        <div className="filter-bar">
          <div>
            <span className="eyebrow">Available quizzes</span>
            <h2>{filteredQuizzes.length} quizzes match this view</h2>
          </div>
          <label className="field filter-field">
            <span>Filter by category</span>
            <select onChange={(event) => setCategoryFilter(event.target.value)} value={categoryFilter}>
              <option value="all">All categories</option>
              {availableCategories.map((category) => (
                <option key={category} value={category}>
                  {category}
                </option>
              ))}
            </select>
          </label>
        </div>

        <section className="list-panel">
          {filteredQuizzes.map((quiz) => (
            <article className="list-row quiz-list-row" key={quiz.id}>
              <div className="quiz-list-copy">
                <div className="quiz-list-head">
                  <StatusPill status={quiz.status} />
                  <span className="category-pill">{normalizeCategory(quiz.category)}</span>
                </div>
                <h2>{quiz.title}</h2>
                <p className="quiz-list-meta">
                  {quiz.question_count} questions · {quiz.total_marks} marks · {quiz.duration_minutes} min
                </p>
              </div>
              <button
                className="button primary"
                disabled={busyQuizId === quiz.id}
                onClick={() => handleAction(quiz)}
                type="button"
              >
                {busyQuizId === quiz.id
                  ? 'Opening...'
                  : quiz.action === 'view_score'
                    ? 'View Score'
                    : quiz.action === 'resume'
                      ? 'Resume'
                      : 'Start'}
              </button>
            </article>
          ))}
        </section>
      </section>
    </div>
  )
}

