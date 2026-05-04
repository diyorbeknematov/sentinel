import { useState } from 'react'
import { useNavigate } from 'react-router-dom'
import api from '../api/axios'

export default function ForgotPassword() {
  const navigate = useNavigate()
  const [email,   setEmail]   = useState('')
  const [loading, setLoading] = useState(false)
  const [sent,    setSent]    = useState(false)
  const [error,   setError]   = useState('')

  const handleSubmit = async (e) => {
    e.preventDefault()
    setError('')
    if (!email) { setError('Email kiritish shart'); return }

    setLoading(true)
    try {
      await api.post('/forgot-password', { email })
      setSent(true)
      setTimeout(() => navigate('/login'), 4000)
    } catch (err) {
      setError(err.response?.data?.message || 'Xatolik yuz berdi')
    } finally {
      setLoading(false)
    }
  }

  return (
    <div style={{
      minHeight: '100vh', background: '#080b14',
      display: 'flex', alignItems: 'center', justifyContent: 'center',
      backgroundImage: `linear-gradient(rgba(99,102,241,0.04) 1px, transparent 1px),
                        linear-gradient(90deg, rgba(99,102,241,0.04) 1px, transparent 1px)`,
      backgroundSize: '40px 40px',
    }}>
      <div style={{
        width: '100%', maxWidth: '380px', margin: '0 16px',
        background: '#0d1120', border: '1px solid #1e293b',
        borderRadius: '16px', padding: '32px',
      }}>

        {/* Logo */}
        <div style={{ display: 'flex', alignItems: 'center', gap: '10px', marginBottom: '28px' }}>
          <div style={{
            width: '32px', height: '32px', borderRadius: '8px',
            background: 'linear-gradient(135deg, #6366f1, #8b5cf6)',
            display: 'flex', alignItems: 'center', justifyContent: 'center',
            fontSize: '12px', fontWeight: '700', color: '#fff',
          }}>S</div>
          <span style={{ fontSize: '14px', fontWeight: '600', color: '#e2e8f0' }}>Sentinel</span>
        </div>

        <h1 style={{ fontSize: '18px', fontWeight: '600', color: '#e2e8f0', margin: '0 0 6px' }}>
          Parolni tiklash
        </h1>
        <p style={{ fontSize: '12px', color: '#475569', margin: '0 0 24px', fontFamily: 'monospace' }}>
          Email manzilingizga reset link yuboramiz
        </p>

        {/* Sent state */}
        {sent ? (
          <div style={{
            padding: '16px', borderRadius: '8px',
            background: 'rgba(34,197,94,0.06)', border: '1px solid rgba(34,197,94,0.2)',
            color: '#4ade80', fontSize: '12px', fontFamily: 'monospace', lineHeight: '1.6',
          }}>
            ✓ Reset link yuborildi. Emailingizni tekshiring.
            <br/>
            <span style={{ color: '#334155' }}>Login sahifasiga yo'naltirilmoqdasiz...</span>
          </div>
        ) : (
          <form onSubmit={handleSubmit}>

            {/* Error */}
            {error && (
              <div style={{
                marginBottom: '16px', padding: '10px 14px', borderRadius: '7px',
                background: 'rgba(239,68,68,0.08)', border: '1px solid rgba(239,68,68,0.2)',
                color: '#fca5a5', fontSize: '12px', fontFamily: 'monospace',
              }}>
                {error}
              </div>
            )}

            {/* Email */}
            <div style={{ marginBottom: '16px' }}>
              <label style={{
                display: 'block', fontSize: '10px', color: '#475569',
                textTransform: 'uppercase', letterSpacing: '0.5px', marginBottom: '8px',
              }}>
                Email
              </label>
              <input
                type="email"
                value={email}
                onChange={e => setEmail(e.target.value)}
                placeholder="user@example.com"
                style={{
                  width: '100%', fontSize: '12px', padding: '10px 14px',
                  borderRadius: '8px', border: '1px solid #1e293b',
                  background: '#080b14', color: '#e2e8f0',
                  fontFamily: 'monospace', outline: 'none', boxSizing: 'border-box',
                }}
                onFocus={e => e.target.style.borderColor = '#6366f1'}
                onBlur={e  => e.target.style.borderColor = '#1e293b'}
              />
            </div>

            {/* Submit */}
            <button
              type="submit"
              disabled={loading}
              style={{
                width: '100%', padding: '10px', borderRadius: '8px', border: 'none',
                background: 'linear-gradient(135deg, #6366f1, #8b5cf6)',
                color: '#fff', fontSize: '13px', fontWeight: '600',
                cursor: loading ? 'not-allowed' : 'pointer',
                opacity: loading ? 0.7 : 1, fontFamily: 'monospace',
                display: 'flex', alignItems: 'center', justifyContent: 'center', gap: '8px',
              }}
            >
              {loading && (
                <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" style={{ animation: 'spin 1s linear infinite' }}>
                  <path d="M21 12a9 9 0 1 1-6.219-8.56"/>
                </svg>
              )}
              {loading ? 'Yuborilmoqda...' : 'Link yuborish'}
            </button>

          </form>
        )}

        {/* Back to login */}
        <div style={{ marginTop: '20px', textAlign: 'center' }}>
          <span
            onClick={() => navigate('/login')}
            style={{ fontSize: '12px', color: '#475569', cursor: 'pointer', fontFamily: 'monospace' }}
          >
            ← Loginga qaytish
          </span>
        </div>

      </div>

      <style>{`@keyframes spin { to { transform: rotate(360deg) } }`}</style>
    </div>
  )
}