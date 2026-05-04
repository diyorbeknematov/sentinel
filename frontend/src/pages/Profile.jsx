import { useState, useEffect } from 'react'
import api from '../api/axios'

export default function Profile() {
  const [user, setUser] = useState(null)

  useEffect(() => {
    const fetchUser = async () => {
      try {
        const res = await api.get('/me')
        setUser(res.data)
      } catch {
        setUser(null)
      }
    }

    fetchUser()
  }, [])

  if (!user) return <div style={{ color:'#94a3b8' }}>Loading...</div>

  return (
    <div style={{ color:'#e2e8f0' }}>
      <h2 style={{ fontSize:'20px', marginBottom:'16px' }}>Profile</h2>

      <div style={{
        background:'#0d1120',
        border:'1px solid #1e293b',
        borderRadius:'12px',
        padding:'20px',
        maxWidth:'400px'
      }}>
        <div style={{ marginBottom:'10px' }}>
          <strong>Username:</strong> {user.username}
        </div>

        <div>
          <strong>Email:</strong> {user.email}
        </div>
      </div>
    </div>
  )
}