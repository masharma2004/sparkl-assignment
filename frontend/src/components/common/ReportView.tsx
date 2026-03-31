import type { AttemptReportResponse } from '../../types/api'
import { formatDate, formatStatusLabel } from '../../utils/format'
import { StatusPill } from './StatusPill'

interface ReportViewProps {
  report: AttemptReportResponse
}

export function ReportView({ report }: ReportViewProps) {
  return (
    <div className="stack-lg">
      <section className="page-header compact-panel">
        <div>
          <span className="eyebrow">Assessment Report</span>
          <h1>{report.quiz.title}</h1>
          <p>
            Attempt #{report.attempt_id} · Submitted {formatDate(report.submitted_at)}
          </p>
        </div>
        <div className="page-header-actions">
          <StatusPill status={report.status} />
          <div className="score-badge">
            <span>Score</span>
            <strong>
              {report.score} / {report.quiz.total_marks}
            </strong>
          </div>
        </div>
      </section>

      <section className="section-panel">
        <div className="metric-grid">
          <article className="metric-card">
            <span>Questions</span>
            <strong>{report.quiz.question_count}</strong>
          </article>
          <article className="metric-card">
            <span>Duration</span>
            <strong>{report.quiz.duration_minutes} min</strong>
          </article>
          <article className="metric-card">
            <span>Status</span>
            <strong>{formatStatusLabel(report.status)}</strong>
          </article>
        </div>
      </section>

      <section className="section-panel stack-md">
        {report.questions.map((question) => (
          <article className="question-report-card" key={question.quiz_question_id}>
            <header className="question-report-header">
              <div>
                <span className="eyebrow">Question {question.sequence_number}</span>
                <h2>{question.prompt}</h2>
              </div>
              <div className="question-marks">
                <span>{question.is_correct ? 'Correct' : 'Needs review'}</span>
                <strong>
                  {question.awarded_marks} / {question.max_marks}
                </strong>
              </div>
            </header>

            <div className="option-columns">
              <div className="option-column">
                <h3>Your answer</h3>
                <ul className="detail-list">
                  {(question.chosen_options.length > 0
                    ? question.chosen_options
                    : ['Not answered']
                  ).map((option) => (
                    <li key={`${question.quiz_question_id}-chosen-${option}`}>{option}</li>
                  ))}
                </ul>
              </div>
              <div className="option-column">
                <h3>Correct answer</h3>
                <ul className="detail-list">
                  {question.correct_options.map((option) => (
                    <li key={`${question.quiz_question_id}-correct-${option}`}>{option}</li>
                  ))}
                </ul>
              </div>
            </div>

            <div className="solution-box">
              <h3>Solution</h3>
              <p>{question.solution}</p>
            </div>
          </article>
        ))}
      </section>
    </div>
  )
}
