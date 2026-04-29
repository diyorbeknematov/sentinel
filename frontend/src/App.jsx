import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom'
import Login from './pages/Login'
import Dashboard from './pages/Dashboard'
import NginxLogs from './pages/NginxLogs'
import AppLogs from './pages/AppLogs'
import Alerts from './pages/Alerts'
import Layout from './components/Layout'

const PrivateRoute = ({ children }) => {
  // const token = localStorage.getItem('sentinel_token')
  // return token ? children : <Navigate to="/login" />
}

export default function App() {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/login" element={<Login />} />
        <Route
          path="/"
          // element={
          //   <PrivateRoute>
          //     <Layout />
          //   </PrivateRoute>
          // }
          element={<Layout />} 
        >
          <Route index element={<Dashboard />} />
          <Route path="nginx-logs" element={<NginxLogs />} />
          <Route path="app-logs" element={<AppLogs />} />
          <Route path="alerts" element={<Alerts />} />
        </Route>
        <Route path="*" element={<Navigate to="/" />} />
      </Routes>
    </BrowserRouter>
  )
}