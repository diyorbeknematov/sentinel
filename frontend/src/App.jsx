import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom'
import Register from './pages/Register'
import Login from './pages/Login'
import ForgotPassword from './pages/ForgotPassword'
import ResetPassword from './pages/ResetPassword'
import Dashboard from './pages/Dashboard'
import Home from './pages/Home'
import Agents from './pages/Agents'
import NginxLogs from './pages/NginxLogs'
import AppLogs from './pages/AppLogs'
import Alerts from './pages/Alerts'
import Layout from './components/Layout'
import Profile from './pages/Profile'

const PrivateRoute = ({ children }) => {
  const token = localStorage.getItem('sentinel_token')
  return token ? children : <Navigate to="/login" />
}

export default function App() {
  return (
    <BrowserRouter>
      <Routes>

        <Route path="/register" element={<Register />} />
        <Route path="/forgot-password" element={<ForgotPassword />} /> 
        <Route path="/reset-password"  element={<ResetPassword />} />
        <Route
          path="/login"
          element={
            localStorage.getItem('sentinel_token')
              ? <Navigate to="/" replace />
              : <Login />
          }
        />

        <Route
          path="/"
          element={
            <PrivateRoute>
              <Layout />
            </PrivateRoute>
          }
        >
          <Route index element={<Dashboard />} />
          <Route path="home" element={<Home />} />
          <Route path="agents" element={<Agents />} />
          <Route path="nginx-logs" element={<NginxLogs />} />
          <Route path="app-logs" element={<AppLogs />} />
          <Route path="alerts" element={<Alerts />} />
          <Route path="profile" element={<Profile />} />
          
        </Route>

        <Route path="*" element={<Navigate to="/" />} />

      </Routes>
    </BrowserRouter>
  )
}