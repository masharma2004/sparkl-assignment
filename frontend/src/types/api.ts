export type Role = 'cms_admin' | 'student'

export interface LoginRequest {
  username: string
  password: string
}

export interface StudentSignupRequest {
  username: string
  password: string
  full_name: string
  email: string
}

export interface UserResponse {
  user_id: number
  username: string
  role: Role
  full_name: string
  email: string
}

export interface LoginResponse {
  user: UserResponse
}

export interface MeResponse {
  user: UserResponse
}

export interface QuizSummaryResponse {
  quiz_id: number
  category: string
  title: string
  question_count: number
  total_marks: number
  duration_minutes: number
  created_by: number
  created_at: string
  updated_at: string
}

export interface GetQuizzesResponse {
  quizzes: QuizSummaryResponse[]
}

export interface QuestionResponse {
  question_id: number
  category: string
  prompt: string
  options: string[]
  correct_options?: string[]
  solution?: string
}

export interface GetQuestionsResponse {
  questions: QuestionResponse[]
}

export interface CreateQuestionRequest {
  category: string
  prompt: string
  options: string[]
  correct_options: string[]
  solution: string
}

export interface CreateQuestionResponse {
  message: string
  question: QuestionResponse
}

export interface CreateQuizQuestionInput {
  question_id: number
  sequence_number: number
  marks: number
}

export interface CreateQuizRequest {
  category: string
  title: string
  question_count: number
  total_marks: number
  duration_minutes: number
  questions: CreateQuizQuestionInput[]
}

export interface CreateQuizResponse {
  message: string
  quiz: QuizSummaryResponse
}

export interface QuizQuestionDetailResponse {
  quiz_question_id: number
  question_id: number
  sequence_number: number
  marks: number
  prompt: string
}

export interface GetQuizResponse {
  quiz: QuizSummaryResponse
  questions: QuizQuestionDetailResponse[]
}

export interface ParticipantResponse {
  student_id: number
  username: string
  full_name: string
  status: string
  score: number
  attempt_id?: number
}

export interface GetParticipantsResponse {
  quiz: QuizSummaryResponse
  participants: ParticipantResponse[]
}

export interface StudentQuizItemResponse {
  id: number
  category: string
  title: string
  question_count: number
  total_marks: number
  duration_minutes: number
  status: string
  action: string
  attempt_id?: number
}

export interface StudentQuizResponse {
  quizzes: StudentQuizItemResponse[]
}

export interface StartQuizResponse {
  message: string
  attempt_id: number
  status: string
}

export interface AttemptQuestionResponse {
  quiz_question_id: number
  question_id: number
  sequence_number: number
  marks: number
  prompt: string
  options: string[]
  chosen_options: string[]
}

export interface AttemptQuizResponse {
  id: number
  category: string
  title: string
  duration_minutes: number
  question_count: number
  total_marks: number
}

export interface GetAttemptResponse {
  attempt_id: number
  status: string
  remaining_seconds: number
  quiz: AttemptQuizResponse
  questions: AttemptQuestionResponse[]
}

export interface SaveAnswerRequest {
  quiz_question_id: number
  chosen_options: string[]
}

export interface SaveAnswerResponse {
  message: string
  attempt_id: number
  quiz_question_id: number
}

export interface FinishAttemptResponse {
  message: string
  attempt_id: number
  status: string
  score: number
  submitted_at?: string | null
}

export interface AttemptReportQuestionResponse {
  quiz_question_id: number
  question_id: number
  sequence_number: number
  prompt: string
  options: string[]
  chosen_options: string[]
  correct_options: string[]
  solution: string
  awarded_marks: number
  max_marks: number
  is_correct: boolean
}

export interface AttemptReportResponse {
  attempt_id: number
  status: string
  score: number
  submitted_at?: string | null
  quiz: AttemptQuizResponse
  questions: AttemptReportQuestionResponse[]
}
