import { useEffect, useMemo, useState } from 'react'
import { Link, useNavigate } from 'react-router'
import { createCmsQuiz, getCmsQuestions } from '../api/cms'
import { ApiError } from '../api/client'
import { useAuth } from '../auth/AuthContext'
import { LoadingBlock } from '../components/common/LoadingBlock'
import type { CreateQuizQuestionInput, QuestionResponse } from '../types/api'
import { matchesCategoryFilter, normalizeCategory } from '../utils/category'

interface QuizDraft {
  category: string
  title: string
  question_count: number
  total_marks: number
  duration_minutes: number
}

const initialDraft: QuizDraft = {
  category: '',
  title: '',
  question_count: 3,
  total_marks: 30,
  duration_minutes: 15,
}

function createQuestionRow(sequenceNumber: number): CreateQuizQuestionInput {
  return {
    question_id: 0,
    sequence_number: sequenceNumber,
    marks: 0,
  }
}

function syncQuestionRows(
  currentRows: CreateQuizQuestionInput[],
  questionCount: number,
): CreateQuizQuestionInput[] {
  return Array.from({ length: Math.max(questionCount, 0) }, (_, index) => {
    const existing = currentRows[index]
    if (!existing) {
      return createQuestionRow(index + 1)
    }

    return {
      ...existing,
      sequence_number: index + 1,
    }
  })
}

export function CMSCreateQuizPage() {
  const { session } = useAuth()
  const navigate = useNavigate()
  const [draft, setDraft] = useState<QuizDraft>(initialDraft)
  const [questionRows, setQuestionRows] = useState<CreateQuizQuestionInput[]>(() =>
    syncQuestionRows([], initialDraft.question_count),
  )
  const [questionBank, setQuestionBank] = useState<QuestionResponse[]>([])
  const [questionCategoryFilter, setQuestionCategoryFilter] = useState('all')
  const [activeRowIndex, setActiveRowIndex] = useState(0)
  const [loading, setLoading] = useState(true)
  const [step, setStep] = useState<1 | 2>(1)
  const [error, setError] = useState<string | null>(null)
  const [isSubmitting, setIsSubmitting] = useState(false)

  useEffect(() => {
    if (!session) {
      return
    }

    let isCancelled = false

    getCmsQuestions()
      .then((response) => {
        if (!isCancelled) {
          setQuestionBank(response.questions)
        }
      })
      .catch((loadError) => {
        if (!isCancelled) {
          setError(loadError instanceof ApiError ? loadError.message : 'Failed to load questions.')
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

  useEffect(() => {
    setQuestionRows((current) => syncQuestionRows(current, draft.question_count))
    setActiveRowIndex((current) => Math.min(current, Math.max(draft.question_count - 1, 0)))
  }, [draft.question_count])

  useEffect(() => {
    const trimmedCategory = draft.category.trim()
    setQuestionCategoryFilter(trimmedCategory.length > 0 ? trimmedCategory : 'all')
  }, [draft.category])

  const questionLookup = useMemo(
    () => new Map(questionBank.map((question) => [question.question_id, question])),
    [questionBank],
  )

  const availableQuestionCategories = useMemo(
    () =>
      Array.from(
        new Set(questionBank.map((question) => normalizeCategory(question.category))),
      ).sort((left, right) => left.localeCompare(right)),
    [questionBank],
  )

  const filteredQuestionBank = useMemo(() => {
    return questionBank.filter((question) =>
      matchesCategoryFilter(question.category, questionCategoryFilter),
    )
  }, [questionBank, questionCategoryFilter])

  const selectedQuestionIds = useMemo(
    () => questionRows.map((row) => row.question_id).filter((questionId) => questionId > 0),
    [questionRows],
  )

  const duplicateQuestionCount = useMemo(() => {
    const seen = new Set<number>()
    let duplicates = 0

    selectedQuestionIds.forEach((questionId) => {
      if (seen.has(questionId)) {
        duplicates += 1
        return
      }

      seen.add(questionId)
    })

    return duplicates
  }, [selectedQuestionIds])

  const marksAllocated = questionRows.reduce((sum, row) => sum + row.marks, 0)
  const remainingMarks = draft.total_marks - marksAllocated
  const rowsReady = questionRows.every((row) => row.question_id > 0 && row.marks > 0)
  const activeRow = questionRows[activeRowIndex] ?? null

  function updateRow(
    rowIndex: number,
    field: keyof Pick<CreateQuizQuestionInput, 'question_id' | 'marks'>,
    value: number,
  ) {
    setQuestionRows((current) =>
      current.map((row, index) =>
        index === rowIndex
          ? {
              ...row,
              [field]: Math.max(0, value),
            }
          : row,
      ),
    )
  }

  function assignQuestionToRow(questionId: number) {
    setQuestionRows((current) =>
      current.map((row, index) =>
        index === activeRowIndex
          ? {
              ...row,
              question_id: questionId,
            }
          : row,
      ),
    )
  }

  async function handleSubmit() {
    if (!session) {
      return
    }

    setError(null)
    setIsSubmitting(true)

    try {
      const payload = {
        ...draft,
        questions: questionRows,
      }

      const response = await createCmsQuiz(payload)
      navigate(`/cms/quizzes/${response.quiz.quiz_id}`, { replace: true })
    } catch (submitError) {
      if (submitError instanceof ApiError) {
        setError(submitError.message)
      } else {
        setError('Failed to create quiz.')
      }
    } finally {
      setIsSubmitting(false)
    }
  }

  if (loading) {
    return (
      <LoadingBlock
        title="Loading question bank"
        message="Pulling the CMS question library so you can map question IDs into the quiz rows."
      />
    )
  }

  return (
    <div className="stack-lg">
      <section className="page-header">
        <div>
          <span className="eyebrow">Create Quiz</span>
          <h1>Build the quiz in two steps</h1>
          <p>Enter the quiz details first, then fill each row with a question ID and marks.</p>
        </div>
        <div className="page-header-actions">
          <Link className="button secondary" to="/cms/questions">
            Open Question Bank
          </Link>
        </div>
      </section>

      {error ? <p className="form-error">{error}</p> : null}

      <section className="wizard-shell">
        <div className="wizard-steps">
          <button
            className={step === 1 ? 'wizard-step wizard-step--active' : 'wizard-step'}
            onClick={() => setStep(1)}
            type="button"
          >
            Step 1 · Quiz Details
          </button>
          <button
            className={step === 2 ? 'wizard-step wizard-step--active' : 'wizard-step'}
            onClick={() => setStep(2)}
            type="button"
          >
            Step 2 · Question Rows
          </button>
        </div>

        {step === 1 ? (
          <div className="form-grid">
            <label className="field">
              <span>Category</span>
              <input
                onChange={(event) =>
                  setDraft((current) => ({ ...current, category: event.target.value }))
                }
                placeholder="For example: Frontend"
                value={draft.category}
              />
            </label>

            <label className="field">
              <span>Title</span>
              <input
                onChange={(event) =>
                  setDraft((current) => ({ ...current, title: event.target.value }))
                }
                placeholder="For example: React Fundamentals"
                value={draft.title}
              />
            </label>

            <label className="field">
              <span>No. of Questions</span>
              <input
                min={1}
                onChange={(event) =>
                  setDraft((current) => ({
                    ...current,
                    question_count: Number(event.target.value) || 0,
                  }))
                }
                type="number"
                value={draft.question_count}
              />
            </label>

            <label className="field">
              <span>Total Score</span>
              <input
                min={1}
                onChange={(event) =>
                  setDraft((current) => ({
                    ...current,
                    total_marks: Number(event.target.value) || 0,
                  }))
                }
                type="number"
                value={draft.total_marks}
              />
            </label>

            <label className="field">
              <span>Duration (minutes)</span>
              <input
                min={1}
                onChange={(event) =>
                  setDraft((current) => ({
                    ...current,
                    duration_minutes: Number(event.target.value) || 0,
                  }))
                }
                type="number"
                value={draft.duration_minutes}
              />
            </label>

            <div className="hero-actions">
              <button
                className="button primary"
                disabled={
                  draft.category.trim().length === 0 ||
                  draft.title.trim().length === 0 ||
                  draft.question_count <= 0 ||
                  draft.total_marks <= 0 ||
                  draft.duration_minutes <= 0
                }
                onClick={() => setStep(2)}
                type="button"
              >
                Continue to Step 2
              </button>
            </div>
          </div>
        ) : (
          <div className="stack-lg">
            <section className="summary-strip">
              <div className="summary-item">
                <span>Rows generated</span>
                <strong>{questionRows.length}</strong>
              </div>
              <div className="summary-item">
                <span>Marks allocated</span>
                <strong>
                  {marksAllocated} / {draft.total_marks}
                </strong>
              </div>
              <div className="summary-item">
                <span>Question ID duplicates</span>
                <strong>{duplicateQuestionCount}</strong>
              </div>
            </section>

            <section className="quiz-builder-grid">
              <div className="stack-md">
                <header className="section-heading">
                  <div>
                    <span className="eyebrow">Step 2</span>
                    <h2>Fill the question rows</h2>
                    <p className="subtle-text">
                      Each row needs a question ID and marks. The order is locked to Q1, Q2, and so on.
                    </p>
                  </div>
                </header>

                <div className="quiz-row-list">
                  {questionRows.map((row, index) => {
                    const questionDetails = questionLookup.get(row.question_id)
                    const isActive = index === activeRowIndex

                    return (
                      <article
                        className={isActive ? 'quiz-row-card quiz-row-card--active' : 'quiz-row-card'}
                        key={row.sequence_number}
                      >
                        <button
                          className="quiz-row-card-trigger"
                          onClick={() => setActiveRowIndex(index)}
                          type="button"
                        >
                          <span className="quiz-row-tag">Q{row.sequence_number}</span>
                          <span className="subtle-text">
                            {questionDetails?.prompt ?? 'Select this row, then enter or pick a #QID.'}
                          </span>
                        </button>

                        <div className="quiz-row-fields">
                          <label className="field">
                            <span>#QID</span>
                            <input
                              min={0}
                              onChange={(event) =>
                                updateRow(index, 'question_id', Number(event.target.value) || 0)
                              }
                              type="number"
                              value={row.question_id || ''}
                            />
                          </label>

                          <label className="field">
                            <span>Marks</span>
                            <input
                              min={0}
                              onChange={(event) =>
                                updateRow(index, 'marks', Number(event.target.value) || 0)
                              }
                              type="number"
                              value={row.marks || ''}
                            />
                          </label>
                        </div>

                        {questionDetails ? (
                          <div className="quiz-row-preview">
                            <span className="category-pill">
                              {normalizeCategory(questionDetails.category)}
                            </span>
                            <span className="subtle-text">Question #{questionDetails.question_id}</span>
                          </div>
                        ) : null}
                      </article>
                    )
                  })}
                </div>
              </div>

              <aside className="stack-md question-reference-panel">
                <header className="section-heading">
                  <div>
                    <span className="eyebrow">Question bank</span>
                    <h2>Use a reference ID</h2>
                    <p className="subtle-text">
                      Select a question below to fill the active row: {activeRow ? `Q${activeRow.sequence_number}` : 'none'}.
                    </p>
                  </div>
                  <label className="field filter-field">
                    <span>Category filter</span>
                    <select
                      onChange={(event) => setQuestionCategoryFilter(event.target.value)}
                      value={questionCategoryFilter}
                    >
                      <option value="all">All categories</option>
                      {availableQuestionCategories.map((category) => (
                        <option key={category} value={category}>
                          {category}
                        </option>
                      ))}
                    </select>
                  </label>
                </header>

                <div className="question-reference-list">
                  {filteredQuestionBank.map((question) => {
                    const isActiveQuestion = activeRow?.question_id === question.question_id
                    const isUsedElsewhere =
                      selectedQuestionIds.includes(question.question_id) && !isActiveQuestion

                    return (
                      <article className="list-row question-reference-row" key={question.question_id}>
                        <div className="quiz-list-copy">
                          <div className="quiz-list-head">
                            <span className="quiz-row-tag">#{question.question_id}</span>
                            <span className="category-pill">{normalizeCategory(question.category)}</span>
                          </div>
                          <h3>{question.prompt}</h3>
                        </div>
                        <button
                          className={isActiveQuestion ? 'button ghost' : 'button secondary'}
                          onClick={() => assignQuestionToRow(question.question_id)}
                          type="button"
                        >
                          {isActiveQuestion
                            ? 'Selected'
                            : isUsedElsewhere
                              ? `Move to Q${(activeRowIndex ?? 0) + 1}`
                              : `Use for Q${(activeRowIndex ?? 0) + 1}`}
                        </button>
                      </article>
                    )
                  })}
                </div>
              </aside>
            </section>

            <div className="hero-actions">
              <button className="button ghost" onClick={() => setStep(1)} type="button">
                Back to Step 1
              </button>
              <button
                className="button primary"
                disabled={
                  isSubmitting ||
                  !rowsReady ||
                  remainingMarks !== 0 ||
                  duplicateQuestionCount > 0
                }
                onClick={handleSubmit}
                type="button"
              >
                {isSubmitting ? 'Creating quiz...' : 'Create Quiz'}
              </button>
            </div>
          </div>
        )}
      </section>
    </div>
  )
}
