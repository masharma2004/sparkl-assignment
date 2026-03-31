import { apiRequest } from './client'
import type {
  AttemptReportResponse,
  FinishAttemptResponse,
  GetAttemptResponse,
  SaveAnswerRequest,
  SaveAnswerResponse,
  StartQuizResponse,
  StudentQuizResponse,
} from '../types/api'

export function getStudentQuizzes() {
  return apiRequest<StudentQuizResponse>('/student/quizzes')
}

export function startStudentQuiz(quizId: string | number) {
  return apiRequest<StartQuizResponse>(`/student/quizzes/${quizId}/start`, {
    method: 'POST',
  })
}

export function getStudentAttempt(attemptId: string | number) {
  return apiRequest<GetAttemptResponse>(`/student/attempts/${attemptId}`)
}

export function saveStudentAnswer(
  attemptId: string | number,
  payload: SaveAnswerRequest,
) {
  return apiRequest<SaveAnswerResponse>(`/student/attempts/${attemptId}/answers`, {
    method: 'PATCH',
    body: payload,
  })
}

export function finishStudentAttempt(attemptId: string | number) {
  return apiRequest<FinishAttemptResponse>(`/student/attempts/${attemptId}/finish`, {
    method: 'POST',
  })
}

export function getStudentReport(attemptId: string | number) {
  return apiRequest<AttemptReportResponse>(`/student/attempts/${attemptId}/report`)
}
