import { useEffect, useMemo, useRef, useState } from 'react'
import { useNavigate, useParams } from 'react-router'
import {
  finishStudentAttempt,
  getStudentAttempt,
  saveStudentAnswer,
} from '../api/student'
import { ApiError } from '../api/client'
import { useAuth } from '../auth/AuthContext'
import { LoadingBlock } from '../components/common/LoadingBlock'
import { formatCountdown } from '../utils/format'
import type { GetAttemptResponse } from '../types/api'

export function StudentAttemptPage() {
  const { attemptId } = useParams()
  const { session } = useAuth()
  const navigate = useNavigate()
  const [attempt, setAttempt] = useState<GetAttemptResponse | null>(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [currentIndex, setCurrentIndex] = useState(0)
  const [savingQuestionId, setSavingQuestionId] = useState<number | null>(null)
  const [finishing, setFinishing] = useState(false)
  const [remainingSeconds, setRemainingSeconds] = useState<number>(0)
  const autoFinishedRef = useRef(false)

  useEffect(() => {
    if (!session || !attemptId) {
      return
    }

    let isCancelled = false

    getStudentAttempt(attemptId)
      .then((response) => {
        if (isCancelled) {
          return
        }

        if (response.status === 'completed') {
          navigate(`/student/attempts/${response.attempt_id}/report`, { replace: true })
          return
        }

        setAttempt(response)
        setRemainingSeconds(response.remaining_seconds)
      })
      .catch((loadError) => {
        if (!isCancelled) {
          setError(loadError instanceof ApiError ? loadError.message : 'Failed to load attempt.')
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
  }, [attemptId, navigate, session])

  useEffect(() => {
    if (!attempt || attempt.status !== 'in_progress') {
      return
    }

    const timerId = window.setInterval(() => {
      setRemainingSeconds((current) => (current > 0 ? current - 1 : 0))
    }, 1000)

    return () => {
      window.clearInterval(timerId)
    }
  }, [attempt])

  useEffect(() => {
    if (
      !attempt ||
      attempt.status !== 'in_progress' ||
      remainingSeconds > 0 ||
      autoFinishedRef.current ||
      !session
    ) {
      return
    }

    autoFinishedRef.current = true
    setFinishing(true)

    finishStudentAttempt(attempt.attempt_id)
      .then(() => {
        navigate(`/student/attempts/${attempt.attempt_id}/report`, { replace: true })
      })
      .catch(() => {
        navigate(`/student/attempts/${attempt.attempt_id}/report`, { replace: true })
      })
  }, [attempt, navigate, remainingSeconds, session])

  const currentQuestion = useMemo(() => {
    if (!attempt) {
      return null
    }

    return attempt.questions[currentIndex] ?? null
  }, [attempt, currentIndex])

  if (loading) {
    return (
      <LoadingBlock
        title="Loading your attempt"
        message="Preparing the quiz, saved answers, and the live countdown timer."
      />
    )
  }

  if (!attempt || !currentQuestion) {
    return <p className="form-error">{error ?? 'Attempt not found.'}</p>
  }

  const activeAttempt = attempt
  const activeQuestion = currentQuestion

  async function handleToggle(option: string) {
    if (!session) {
      return
    }

    const nextChosenOptions = activeQuestion.chosen_options.includes(option)
      ? activeQuestion.chosen_options.filter((item) => item !== option)
      : [...activeQuestion.chosen_options, option]

    setAttempt((current) =>
      current
        ? {
            ...current,
            questions: current.questions.map((question) =>
              question.quiz_question_id === activeQuestion.quiz_question_id
                ? { ...question, chosen_options: nextChosenOptions }
                : question,
            ),
          }
        : current,
    )

    setSavingQuestionId(activeQuestion.quiz_question_id)
    setError(null)

    try {
      await saveStudentAnswer(activeAttempt.attempt_id, {
        quiz_question_id: activeQuestion.quiz_question_id,
        chosen_options: nextChosenOptions,
      })
    } catch (saveError) {
      if (saveError instanceof ApiError && saveError.status === 409) {
        navigate(`/student/attempts/${activeAttempt.attempt_id}/report`, { replace: true })
        return
      }

      setError(saveError instanceof ApiError ? saveError.message : 'Failed to save answer.')
    } finally {
      setSavingQuestionId(null)
    }
  }

  async function handleFinish() {
    if (!session || finishing) {
      return
    }

    setFinishing(true)
    setError(null)

    try {
      await finishStudentAttempt(activeAttempt.attempt_id)
      navigate(`/student/attempts/${activeAttempt.attempt_id}/report`, { replace: true })
    } catch (finishError) {
      setError(finishError instanceof ApiError ? finishError.message : 'Failed to finish attempt.')
      setFinishing(false)
    }
  }

  return (
    <div className="attempt-shell">
      <header className="attempt-topbar">
        <div>
          <span className="eyebrow">Quiz in progress</span>
          <h1>{activeAttempt.quiz.title}</h1>
        </div>
        <div className="attempt-topbar-actions">
          <div className="timer-card timer-card--inline">
            <span>Time Left</span>
            <strong>{formatCountdown(remainingSeconds)}</strong>
          </div>
          <button className="button danger" onClick={handleFinish} type="button">
            {finishing ? 'Finishing...' : 'Finish'}
          </button>
        </div>
      </header>

      {error ? <p className="form-error">{error}</p> : null}

      <section className="attempt-main attempt-main--wide">
        <div className="question-progress-strip">
          {attempt.questions.map((question, index) => (
            <button
              className={
                index === currentIndex
                  ? 'question-index-button question-index-button--active'
                  : question.chosen_options.length > 0
                    ? 'question-index-button question-index-button--answered'
                    : 'question-index-button'
              }
              key={question.quiz_question_id}
              onClick={() => setCurrentIndex(index)}
              type="button"
            >
              Q{question.sequence_number}
            </button>
          ))}
        </div>

        <div className="page-header compact-panel">
          <div>
            <span className="eyebrow">Question {activeQuestion.sequence_number}</span>
            <h2>{activeQuestion.prompt}</h2>
            <p>
              {activeAttempt.quiz.category} · {activeAttempt.quiz.question_count} questions ·{' '}
              {activeAttempt.quiz.total_marks} marks
            </p>
          </div>
          <div className="page-header-actions">
            <span className="score-badge">
              <span>Marks</span>
              <strong>{activeQuestion.marks}</strong>
            </span>
          </div>
        </div>

        <section className="option-stack">
          {activeQuestion.options.map((option) => {
            const checked = activeQuestion.chosen_options.includes(option)

            return (
              <button
                className={checked ? 'option-card option-card--selected' : 'option-card'}
                key={option}
                onClick={() => handleToggle(option)}
                type="button"
              >
                <span className="option-marker">{checked ? '✓' : ''}</span>
                <span>{option}</span>
              </button>
            )
          })}
        </section>

        <footer className="attempt-footer">
          <button
            className="button ghost"
            disabled={currentIndex === 0}
            onClick={() => setCurrentIndex((index) => Math.max(index - 1, 0))}
            type="button"
          >
            {'<< Prev'}
          </button>
          <p className="subtle-text">
            {savingQuestionId === activeQuestion.quiz_question_id
              ? 'Saving selection...'
              : 'Answers save while you move through the quiz.'}
          </p>
          <button
            className="button secondary"
            disabled={currentIndex === activeAttempt.questions.length - 1}
            onClick={() =>
              setCurrentIndex((index) =>
                Math.min(index + 1, activeAttempt.questions.length - 1),
              )
            }
            type="button"
          >
            {'Next >>'}
          </button>
        </footer>
      </section>
    </div>
  )
}

