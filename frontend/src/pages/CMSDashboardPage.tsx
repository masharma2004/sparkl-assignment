import { useEffect, useMemo, useState } from 'react'
import { Link } from 'react-router'
import { getCmsQuizzes } from '../api/cms'
import { ApiError } from '../api/client'
import { useAuth } from '../auth/AuthContext'
import { LoadingBlock } from '../components/common/LoadingBlock'
import type { QuizSummaryResponse } from '../types/api'
import { normalizeCategory, matchesCategoryFilter } from '../utils/category'
import { formatDate } from '../utils/format'

export function CMSDashboardPage() {
  const { session } = useAuth()
  const [quizzes, setQuizzes] = useState<QuizSummaryResponse[]>([])
  const [categoryFilter, setCategoryFilter] = useState('all')
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    if (!session) {
      return
    }

    let isCancelled = false

    getCmsQuizzes()
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

  if (loading) {
    return <LoadingBlock title="Loading quiz inventory" message="Pulling the latest quiz catalogue from the CMS API." />
  }

  return (
    <div className="stack-lg">
      <section className="page-header">
        <div>
          <span className="eyebrow">CMS Dashboard</span>
          <h1>Quiz control room</h1>
          <p>Open quizzes, review participants, and start the next assessment build.</p>
        </div>
        <div className="page-header-actions">
          <Link className="button secondary" to="/cms/questions">
            Question Bank
          </Link>
          <Link className="button primary" to="/cms/quizzes/new">
            Create Quiz
          </Link>
        </div>
      </section>

      {error ? <p className="form-error">{error}</p> : null}

      <section className="section-panel stack-md">
        <div className="filter-bar">
          <div>
            <span className="eyebrow">Overview</span>
            <h2>Quiz inventory</h2>
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

        <div className="summary-strip">
          <div className="summary-item">
            <span>Quizzes in view</span>
            <strong>{filteredQuizzes.length}</strong>
          </div>
          <div className="summary-item">
            <span>Questions in view</span>
            <strong>{filteredQuizzes.reduce((sum, quiz) => sum + quiz.question_count, 0)}</strong>
          </div>
        </div>
      </section>

      <section className="section-panel stack-md">
        <header className="section-heading">
          <div>
            <span className="eyebrow">Quiz library</span>
            <h2>{filteredQuizzes.length} quizzes ready</h2>
          </div>
        </header>

        <section className="list-panel">
          {filteredQuizzes.map((quiz) => (
            <article className="list-row quiz-list-row" key={quiz.quiz_id}>
              <div className="quiz-list-copy">
                <div className="quiz-list-head">
                  <span className="category-pill">{normalizeCategory(quiz.category)}</span>
                  <span className="subtle-text">{formatDate(quiz.updated_at)}</span>
                </div>
                <h2>{quiz.title}</h2>
                <p className="quiz-list-meta">
                  {quiz.question_count} questions · {quiz.total_marks} marks ·{' '}
                  {quiz.duration_minutes} minutes
                </p>
              </div>
              <div className="resource-actions">
                <Link className="button secondary" to={`/cms/quizzes/${quiz.quiz_id}`}>
                  View quiz
                </Link>
                <Link className="button ghost" to={`/cms/quizzes/${quiz.quiz_id}/participants`}>
                  Participants
                </Link>
              </div>
            </article>
          ))}
        </section>
      </section>
    </div>
  )
}

