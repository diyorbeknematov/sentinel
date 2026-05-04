import { useState, useEffect } from 'react'
import { useSearchParams, useNavigate } from 'react-router-dom'
import api from '../api/axios'

const EyeIcon = () => (
  <svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
    <path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z"/>
    <circle cx="12" cy="12" r="3"/>
  </svg>
)

const EyeOffIcon = () => (
  <svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
    <path d="M17.94 17.94A10.07 10.07 0 0 1 12 20c-7 0-11-8-11-8a18.45 18.45 0 0 1 5.06-5.94"/>
    <path d="M9.9 4.24A9.12 9.12 0 0 1 12 4c7 0 11 8 11 8a18.5 18.5 0 0 1-2.16 3.19"/>
    <line x1="1" y1="1" x2="23" y2="23"/>
  </svg>
)

const inputStyle = {
  width: '100%', fontSize: '12px', padding: '10px 40px 10px 14px',
  borderRadius: '8px', border: '1px solid #1e293b',
  background: '#080b14', color: '#e2e8f0',
  fontFamily: 'monospace', outline: 'none', boxSizing: 'border-box',
}

export default function ResetPassword() {
  const [searchParams] = useSearchParams()
  const navigate       = useNavigate()
  const token          = searchParams.get('token')

  const [password,  setPassword]  = useState('')
  const [confirm,   setConfirm]   = useState('')
  const [showPass,  setShowPass]  = useState(false)
  const [showConf,  setShowConf]  = useState(false)
  const [loading,   setLoading]   = useState(false)
  const [success,   setSuccess]   = useState(false)
  const [error,     setError]     = useState('')

  useEffect(() => {
    if (!token) setError("Reset link noto'g'ri yoki muddati tugagan")
  }, [token])

  const handleSubmit = async (e) => {
    e.preventDefault()
    setError('')

    if (!password || !confirm)       { setError('Barcha maydonlarni to\'ldiring'); return }
    if (password.length < 6)         { setError('Parol kamida 6 ta belgidan iborat bo\'lishi kerak'); return }
    if (password !== confirm)        { setError('Parollar mos kelmadi'); return }

    setLoading(true)
    try {
      await api.post('/reset-password', { token, password })
      setSuccess(true)
      setTimeout(() => navigate('/login'), 3000)
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
          Yangi parol o'rnatish
        </h1>
        <p style={{ fontSize: '12px', color: '#475569', margin: '0 0 24px', fontFamily: 'monospace' }}>
          Kamida 6 ta belgidan iborat parol kiriting
        </p>

        {/* Success */}
        {success ? (
          <div style={{
            padding: '16px', borderRadius: '8px',
            background: 'rgba(34,197,94,0.06)', border: '1px solid rgba(34,197,94,0.2)',
            color: '#4ade80', fontSize: '12px', fontFamily: 'monospace', lineHeight: '1.6',
          }}>
            ✓ Parol muvaffaqiyatli o'zgartirildi.
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

            {/* New password */}
            <div style={{ marginBottom: '14px' }}>
              <label style={{
                display: 'block', fontSize: '10px', color: '#475569',
                textTransform: 'uppercase', letterSpacing: '0.5px', marginBottom: '8px',
              }}>
                Yangi parol
              </label>
              <div style={{ position: 'relative' }}>
                <input
                  type={showPass ? 'text' : 'password'}
                  value={password}
                  onChange={e => setPassword(e.target.value)}
                  placeholder="••••••••"
                  style={inputStyle}
                  onFocus={e => e.target.style.borderColor = '#6366f1'}
                  onBlur={e  => e.target.style.borderColor = '#1e293b'}
                />
                <button type="button" onClick={() => setShowPass(v => !v)} style={{
                  position: 'absolute', right: '12px', top: '50%', transform: 'translateY(-50%)',
                  background: 'none', border: 'none', cursor: 'pointer', color: '#475569', padding: 0,
                }}>
                  {showPass ? <EyeOffIcon /> : <EyeIcon />}
                </button>
              </div>
            </div>

            {/* Confirm password */}
            <div style={{ marginBottom: '20px' }}>
              <label style={{
                display: 'block', fontSize: '10px', color: '#475569',
                textTransform: 'uppercase', letterSpacing: '0.5px', marginBottom: '8px',
              }}>
                Parolni tasdiqlang
              </label>
              <div style={{ position: 'relative' }}>
                <input
                  type={showConf ? 'text' : 'password'}
                  value={confirm}
                  onChange={e => setConfirm(e.target.value)}
                  placeholder="••••••••"
                  style={{
                    ...inputStyle,
                    borderColor: confirm && confirm !== password ? 'rgba(239,68,68,0.4)' : '#1e293b',
                  }}
                  onFocus={e => e.target.style.borderColor = '#6366f1'}
                  onBlur={e  => e.target.style.borderColor = confirm && confirm !== password
                    ? 'rgba(239,68,68,0.4)' : '#1e293b'
                  }
                />
                <button type="button" onClick={() => setShowConf(v => !v)} style={{
                  position: 'absolute', right: '12px', top: '50%', transform: 'translateY(-50%)',
                  background: 'none', border: 'none', cursor: 'pointer', color: '#475569', padding: 0,
                }}>
                  {showConf ? <EyeOffIcon /> : <EyeIcon />}
                </button>
              </div>
              {/* Inline match indicator */}
              {confirm && (
                <p style={{
                  fontSize: '10px', fontFamily: 'monospace', marginTop: '5px',
                  color: confirm === password ? '#4ade80' : '#f87171',
                }}>
                  {confirm === password ? '✓ Parollar mos' : '✗ Parollar mos emas'}
                </p>
              )}
            </div>

            {/* Submit */}
            <button
              type="submit"
              disabled={loading || !token}
              style={{
                width: '100%', padding: '10px', borderRadius: '8px', border: 'none',
                background: 'linear-gradient(135deg, #6366f1, #8b5cf6)',
                color: '#fff', fontSize: '13px', fontWeight: '600',
                cursor: loading || !token ? 'not-allowed' : 'pointer',
                opacity: loading || !token ? 0.7 : 1, fontFamily: 'monospace',
                display: 'flex', alignItems: 'center', justifyContent: 'center', gap: '8px',
              }}
            >
              {loading && (
                <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" style={{ animation: 'spin 1s linear infinite' }}>
                  <path d="M21 12a9 9 0 1 1-6.219-8.56"/>
                </svg>
              )}
              {loading ? 'Saqlanmoqda...' : 'Parolni o\'rnatish'}
            </button>

          </form>
        )}

        {/* Back */}
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