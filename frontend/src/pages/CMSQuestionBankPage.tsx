import { useEffect, useMemo, useState } from 'react'
import { Link } from 'react-router'
import { createCmsQuestion, getCmsQuestions } from '../api/cms'
import { ApiError } from '../api/client'
import { useAuth } from '../auth/AuthContext'
import { LoadingBlock } from '../components/common/LoadingBlock'
import type { QuestionResponse } from '../types/api'
import { normalizeCategory, matchesCategoryFilter } from '../utils/category'

interface QuestionOptionDraft {
  id: number
  value: string
  isCorrect: boolean
}

function createInitialOptionDrafts(): QuestionOptionDraft[] {
  return Array.from({ length: 4 }, (_, index) => ({
    id: index + 1,
    value: '',
    isCorrect: false,
  }))
}

export function CMSQuestionBankPage() {
  const { session } = useAuth()
  const [questions, setQuestions] = useState<QuestionResponse[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [successMessage, setSuccessMessage] = useState<string | null>(null)
  const [category, setCategory] = useState('')
  const [categoryFilter, setCategoryFilter] = useState('all')
  const [prompt, setPrompt] = useState('')
  const [solution, setSolution] = useState('')
  const [optionDrafts, setOptionDrafts] = useState<QuestionOptionDraft[]>(() =>
    createInitialOptionDrafts(),
  )
  const [isSubmitting, setIsSubmitting] = useState(false)

  useEffect(() => {
    if (!session) {
      return
    }

    let isCancelled = false

    getCmsQuestions()
      .then((response) => {
        if (!isCancelled) {
          setQuestions(response.questions)
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

  const usableOptions = useMemo(
    () => optionDrafts.map((option) => option.value.trim()).filter((option) => option.length > 0),
    [optionDrafts],
  )

  const selectedCorrectOptions = useMemo(
    () =>
      optionDrafts
        .filter((option) => option.isCorrect && option.value.trim().length > 0)
        .map((option) => option.value.trim()),
    [optionDrafts],
  )

  const trimmedCategory = category.trim()
  const trimmedPrompt = prompt.trim()
  const trimmedSolution = solution.trim()

  const duplicateOptionCount = useMemo(() => {
    const seen = new Set<string>()
    let duplicates = 0

    usableOptions.forEach((option) => {
      const normalized = option.toLocaleLowerCase()
      if (seen.has(normalized)) {
        duplicates += 1
        return
      }

      seen.add(normalized)
    })

    return duplicates
  }, [usableOptions])

  const canSubmit =
    trimmedCategory.length > 0 &&
    trimmedPrompt.length > 0 &&
    usableOptions.length >= 2 &&
    selectedCorrectOptions.length > 0 &&
    duplicateOptionCount === 0

  const availableCategories = useMemo(
    () =>
      Array.from(
        new Set(questions.map((question) => normalizeCategory(question.category))),
      ).sort((left, right) => left.localeCompare(right)),
    [questions],
  )

  const filteredQuestions = useMemo(() => {
    return questions.filter((question) => matchesCategoryFilter(question.category, categoryFilter))
  }, [categoryFilter, questions])

  function updateOptionDraft(id: number, value: string) {
    setOptionDrafts((current) =>
      current.map((option) => (option.id === id ? { ...option, value } : option)),
    )
  }

  function toggleCorrectOption(id: number) {
    setOptionDrafts((current) =>
      current.map((option) =>
        option.id === id ? { ...option, isCorrect: !option.isCorrect } : option,
      ),
    )
  }

  function addOptionField() {
    setOptionDrafts((current) => [
      ...current,
      {
        id: current.length === 0 ? 1 : Math.max(...current.map((option) => option.id)) + 1,
        value: '',
        isCorrect: false,
      },
    ])
  }

  function removeOptionField(id: number) {
    setOptionDrafts((current) => {
      if (current.length <= 2) {
        return current
      }

      return current.filter((option) => option.id !== id)
    })
  }

  async function handleSubmit(event: React.FormEvent<HTMLFormElement>) {
    event.preventDefault()
    if (!session) {
      return
    }

    setError(null)
    setSuccessMessage(null)
    setIsSubmitting(true)

    try {
      const response = await createCmsQuestion({
        category: trimmedCategory,
        prompt: trimmedPrompt,
        options: usableOptions,
        correct_options: selectedCorrectOptions,
        solution: trimmedSolution,
      })

      setQuestions((current) => [response.question, ...current])
      setCategory('')
      setPrompt('')
      setSolution('')
      setOptionDrafts(createInitialOptionDrafts())
      setSuccessMessage('Question added to the bank.')
    } catch (submitError) {
      if (submitError instanceof ApiError) {
        setError(submitError.message)
      } else {
        setError('Failed to create question.')
      }
    } finally {
      setIsSubmitting(false)
    }
  }

  if (loading) {
    return (
      <LoadingBlock
        title="Loading question bank"
        message="Preparing the question catalogue so you can add new items and review what already exists."
      />
    )
  }

  return (
    <div className="stack-lg">
      <section className="page-header">
        <div>
          <span className="eyebrow">Question Bank</span>
          <h1>Author questions before you assemble quizzes</h1>
          <p>
            Create reusable questions with clear options, mark the correct answers, and keep
            the bank ready for new quizzes.
          </p>
        </div>
        <div className="page-header-actions">
          <Link className="button secondary" to="/cms/quizzes/new">
            Back to quiz builder
          </Link>
        </div>
      </section>

      {error ? <p className="form-error">{error}</p> : null}
      {successMessage ? <p className="success-banner">{successMessage}</p> : null}

      <section className="question-bank-layout">
        <section className="wizard-shell question-composer-shell">
          <header className="section-heading question-composer-header">
            <div>
              <span className="eyebrow">Create question</span>
              <h2>Build one reusable question at a time</h2>
              <p className="subtle-text">
                Start with the prompt and category, then shape the answer set and mark the
                correct choices.
              </p>
            </div>
            <span
              className={
                canSubmit
                  ? 'status-pill status-pill--success'
                  : 'status-pill status-pill--muted'
              }
            >
              {canSubmit ? 'Ready to add' : 'Needs input'}
            </span>
          </header>

          <form className="question-composer-grid" onSubmit={handleSubmit}>
            <div className="stack-md">
              <label className="field">
                <span>Category</span>
                <input
                  onChange={(event) => setCategory(event.target.value)}
                  placeholder="For example: Frontend"
                  required
                  value={category}
                />
              </label>

              <label className="field">
                <span>Prompt</span>
                <textarea
                  onChange={(event) => setPrompt(event.target.value)}
                  placeholder="Write the question prompt"
                  required
                  rows={5}
                  value={prompt}
                />
              </label>

              <label className="field">
                <span>Solution</span>
                <textarea
                  onChange={(event) => setSolution(event.target.value)}
                  placeholder="Optional explanation shown in reports"
                  rows={5}
                  value={solution}
                />
              </label>

              <article className="composer-checklist">
                <div className="composer-checklist-item">
                  <span
                    className={
                      usableOptions.length >= 2
                        ? 'status-pill status-pill--success'
                        : 'status-pill status-pill--muted'
                    }
                  >
                    {usableOptions.length >= 2 ? 'Done' : 'Needed'}
                  </span>
                  <p>Keep at least two filled answer choices.</p>
                </div>
                <div className="composer-checklist-item">
                  <span
                    className={
                      selectedCorrectOptions.length > 0
                        ? 'status-pill status-pill--success'
                        : 'status-pill status-pill--muted'
                    }
                  >
                    {selectedCorrectOptions.length > 0 ? 'Done' : 'Needed'}
                  </span>
                  <p>Mark one or more correct answers.</p>
                </div>
                <div className="composer-checklist-item">
                  <span
                    className={
                      duplicateOptionCount === 0
                        ? 'status-pill status-pill--success'
                        : 'status-pill status-pill--warning'
                    }
                  >
                    {duplicateOptionCount === 0 ? 'Clear' : 'Fix'}
                  </span>
                  <p>Avoid duplicate answer labels so scoring stays unambiguous.</p>
                </div>
              </article>
            </div>

            <div className="stack-md">
              <div className="section-heading">
                <div>
                  <span className="eyebrow">Answer choices</span>
                  <h3>Shape the answer set</h3>
                  <p className="subtle-text">
                    Fill the choice text first, then mark the correct answer cards.
                  </p>
                </div>
                <button className="button ghost inline-button" onClick={addOptionField} type="button">
                  Add option
                </button>
              </div>

              <div className="question-option-grid">
                {optionDrafts.map((option, index) => (
                  <article
                    className={
                      option.isCorrect
                        ? 'question-option-card question-option-card--correct'
                        : 'question-option-card'
                    }
                    key={option.id}
                  >
                    <div className="question-option-card-top">
                      <span className="option-letter">
                        {String.fromCharCode(65 + index)}
                      </span>
                      <button
                        className={
                          option.isCorrect
                            ? 'button secondary option-toggle option-toggle--active'
                            : 'button ghost option-toggle'
                        }
                        onClick={() => toggleCorrectOption(option.id)}
                        type="button"
                      >
                        {option.isCorrect ? 'Correct answer' : 'Mark correct'}
                      </button>
                    </div>

                    <label className="field">
                      <span>Choice text</span>
                      <input
                        onChange={(event) => updateOptionDraft(option.id, event.target.value)}
                        placeholder={`Option ${index + 1}`}
                        value={option.value}
                      />
                    </label>

                    <div className="question-option-card-actions">
                      <button
                        className="button ghost inline-button"
                        disabled={optionDrafts.length <= 2}
                        onClick={() => removeOptionField(option.id)}
                        type="button"
                      >
                        Remove option
                      </button>
                    </div>
                  </article>
                ))}
              </div>

              <div className="composer-summary-grid">
                <article className="metric-card">
                  <span>Filled options</span>
                  <strong>{usableOptions.length}</strong>
                </article>
                <article className="metric-card">
                  <span>Correct choices</span>
                  <strong>{selectedCorrectOptions.length}</strong>
                </article>
                <article className="metric-card">
                  <span>Duplicates</span>
                  <strong>{duplicateOptionCount}</strong>
                </article>
              </div>

              <article className="question-preview-card">
                <div className="resource-card-header">
                  <span className="eyebrow">Preview</span>
                  <span className="category-pill">
                    {trimmedCategory.length > 0 ? trimmedCategory : 'No category yet'}
                  </span>
                </div>
                <h3>
                  {trimmedPrompt.length > 0
                    ? trimmedPrompt
                    : 'Your question prompt will appear here as you type.'}
                </h3>
                <ul className="detail-list">
                  {usableOptions.length > 0 ? (
                    usableOptions.map((option, index) => (
                      <li key={`${index + 1}-${option}`}>
                        {String.fromCharCode(65 + index)}. {option}
                      </li>
                    ))
                  ) : (
                    <li>Add at least two options to preview the answer set.</li>
                  )}
                </ul>
                <p className="subtle-text">
                  Correct choices: {selectedCorrectOptions.length > 0 ? selectedCorrectOptions.join(', ') : 'None selected yet'}
                </p>
                {trimmedSolution.length > 0 ? (
                  <p className="subtle-text">Solution: {trimmedSolution}</p>
                ) : null}
              </article>

              <div className="hero-actions">
                <button className="button primary" disabled={!canSubmit || isSubmitting} type="submit">
                  {isSubmitting ? 'Adding question...' : 'Add question to bank'}
                </button>
              </div>
            </div>
          </form>
        </section>

        <section className="section-panel stack-md">
          <header className="section-heading">
            <div>
              <span className="eyebrow">Available questions</span>
              <h2>{filteredQuestions.length} questions in view</h2>
            </div>
            <label className="field filter-field">
              <span>Category filter</span>
              <select onChange={(event) => setCategoryFilter(event.target.value)} value={categoryFilter}>
                <option value="all">All categories</option>
                {availableCategories.map((item) => (
                  <option key={item} value={item}>
                    {item}
                  </option>
                ))}
              </select>
            </label>
          </header>

          <div className="list-panel">
            {filteredQuestions.map((question) => (
              <article className="list-row question-bank-row" key={question.question_id}>
                <div className="quiz-list-copy">
                  <div className="quiz-list-head">
                    <span className="quiz-row-tag">#{question.question_id}</span>
                    <span className="category-pill">{normalizeCategory(question.category)}</span>
                  </div>
                  <h3>{question.prompt}</h3>
                  <p className="quiz-list-meta">
                    Correct: {question.correct_options?.join(', ') ?? 'Not available'}
                  </p>
                  {question.solution ? <p className="subtle-text">Solution: {question.solution}</p> : null}
                </div>
                <ul className="detail-list question-bank-options">
                  {question.options.map((option) => (
                    <li key={`${question.question_id}-${option}`}>{option}</li>
                  ))}
                </ul>
              </article>
            ))}
          </div>
        </section>
      </section>
    </div>
  )
}
