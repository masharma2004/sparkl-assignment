import { apiRequest } from './client'
import type {
  CreateQuestionRequest,
  CreateQuestionResponse,
  CreateQuizRequest,
  CreateQuizResponse,
  GetParticipantsResponse,
  GetQuestionsResponse,
  GetQuizResponse,
  GetQuizzesResponse,
  AttemptReportResponse,
} from '../types/api'

export function getCmsQuestions() {
  return apiRequest<GetQuestionsResponse>('/cms/questions')
}

export function createCmsQuestion(payload: CreateQuestionRequest) {
  return apiRequest<CreateQuestionResponse>('/cms/questions', {
    method: 'POST',
    body: payload,
  })
}

export function getCmsQuizzes() {
  return apiRequest<GetQuizzesResponse>('/cms/quizzes')
}

export function getCmsQuiz(quizId: string | number) {
  return apiRequest<GetQuizResponse>(`/cms/quizzes/${quizId}`)
}

export function createCmsQuiz(payload: CreateQuizRequest) {
  return apiRequest<CreateQuizResponse>('/cms/quizzes', {
    method: 'POST',
    body: payload,
  })
}

export function getCmsParticipants(quizId: string | number) {
  return apiRequest<GetParticipantsResponse>(`/cms/quizzes/${quizId}/participants`)
}

export function getCmsParticipantReport(
  quizId: string | number,
  studentId: string | number,
) {
  return apiRequest<AttemptReportResponse>(
    `/cms/quizzes/${quizId}/participants/${studentId}/report`,
  )
}
