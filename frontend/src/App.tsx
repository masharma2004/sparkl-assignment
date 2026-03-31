import { RouterProvider } from 'react-router/dom'
import { AuthProvider } from './auth/AuthContext'
import { AppErrorBoundary } from './components/common/AppErrorBoundary'
import { router } from './router'

function App() {
  return (
    <AppErrorBoundary>
      <AuthProvider>
        <RouterProvider router={router} />
      </AuthProvider>
    </AppErrorBoundary>
  )
}

export default App

