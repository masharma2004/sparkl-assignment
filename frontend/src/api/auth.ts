import { apiRequest } from './client'
import type {
  LoginRequest,
  LoginResponse,
  MeResponse,
  StudentSignupRequest,
} from '../types/api'

export function loginCms(payload: LoginRequest) {
  return apiRequest<LoginResponse>('/auth/cms/login', {
    method: 'POST',
    body: payload,
  })
}

export function loginStudent(payload: LoginRequest) {
  return apiRequest<LoginResponse>('/auth/student/login', {
    method: 'POST',
    body: payload,
  })
}

export function signupStudent(payload: StudentSignupRequest) {
  return apiRequest<LoginResponse>('/auth/student/signup', {
    method: 'POST',
    body: payload,
  })
}

export function logout() {
  return apiRequest<{ message: string }>('/auth/logout', {
    method: 'POST',
  })
}

export function getMe() {
  return apiRequest<MeResponse>('/auth/me')
}
