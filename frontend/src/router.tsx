import { createBrowserRouter } from 'react-router'
import { ProtectedRoute, PublicOnlyRoute } from './auth/ProtectedRoute'
import { PortalLayout } from './components/layout/PortalLayout'
import { CMSCreateQuizPage } from './pages/CMSCreateQuizPage'
import { CMSDashboardPage } from './pages/CMSDashboardPage'
import { CMSQuestionBankPage } from './pages/CMSQuestionBankPage'
import { CMSLoginPage, StudentLoginPage, StudentSignupPage } from './pages/LoginPage'
import { CMSParticipantReportPage } from './pages/CMSParticipantReportPage'
import { CMSParticipantsPage } from './pages/CMSParticipantsPage'
import { CMSQuizDetailPage } from './pages/CMSQuizDetailPage'
import { HomePage } from './pages/HomePage'
import { NotFoundPage } from './pages/NotFoundPage'
import { RouteErrorPage } from './pages/RouteErrorPage'
import { StudentAttemptPage } from './pages/StudentAttemptPage'
import { StudentDashboardPage } from './pages/StudentDashboardPage'
import { StudentReportPage } from './pages/StudentReportPage'

export const router = createBrowserRouter([
  {
    path: '/',
    element: <HomePage />,
    errorElement: <RouteErrorPage />,
  },
  {
    path: '/login/cms',
    element: (
      <PublicOnlyRoute>
        <CMSLoginPage />
      </PublicOnlyRoute>
    ),
    errorElement: <RouteErrorPage />,
  },
  {
    path: '/login/student',
    element: (
      <PublicOnlyRoute>
        <StudentLoginPage />
      </PublicOnlyRoute>
    ),
    errorElement: <RouteErrorPage />,
  },
  {
    path: '/signup/student',
    element: (
      <PublicOnlyRoute>
        <StudentSignupPage />
      </PublicOnlyRoute>
    ),
    errorElement: <RouteErrorPage />,
  },
  {
    path: '/cms',
    element: (
      <ProtectedRoute role="cms_admin">
        <PortalLayout portal="cms" />
      </ProtectedRoute>
    ),
    errorElement: <RouteErrorPage />,
    children: [
      { index: true, element: <CMSDashboardPage /> },
      { path: 'questions', element: <CMSQuestionBankPage /> },
      { path: 'quizzes/new', element: <CMSCreateQuizPage /> },
      { path: 'quizzes/:quizId', element: <CMSQuizDetailPage /> },
      { path: 'quizzes/:quizId/participants', element: <CMSParticipantsPage /> },
      {
        path: 'quizzes/:quizId/participants/:studentId/report',
        element: <CMSParticipantReportPage />,
      },
    ],
  },
  {
    path: '/student',
    element: (
      <ProtectedRoute role="student">
        <PortalLayout portal="student" />
      </ProtectedRoute>
    ),
    errorElement: <RouteErrorPage />,
    children: [
      { index: true, element: <StudentDashboardPage /> },
      { path: 'attempts/:attemptId', element: <StudentAttemptPage /> },
      { path: 'attempts/:attemptId/report', element: <StudentReportPage /> },
    ],
  },
  {
    path: '*',
    element: <NotFoundPage />,
  },
])

