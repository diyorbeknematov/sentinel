import { useState } from 'react'
import { useNavigate } from 'react-router-dom'
import api from '../api/axios'

export default function ForgotPassword() {
  const navigate = useNavigate()
  const [email, setEmail] = useState('')
  const [loading, setLoading] = useState(false)
  const [sent, setSent] = useState(false)
  const [error, setError] = useState('')

  const handleSubmit = async (e) => {
    e.preventDefault()
    setError('')
    if (!email) {
      setError('Email kiritish shart')
      return
    }

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
    <div
      className="min-h-screen bg-[#F8FAFC] flex items-center justify-center font-sans"
      style={{
        backgroundImage: `linear-gradient(rgba(79,70,229,0.03) 1px, transparent 1px),
                          linear-gradient(90deg, rgba(79,70,229,0.03) 1px, transparent 1px)`,
        backgroundSize: '40px 40px'
      }}
    >
      <div className="w-full max-w-sm mx-6 bg-white border border-[#E2E8F0] rounded-2xl p-8 shadow-xl">
        
        {/* Logo / Brend qismi */}
        <div className="flex items-center gap-3 mb-6">
          <div className="w-9 h-9 rounded-lg bg-gradient-to-br from-indigo-500 to-purple-600 flex items-center justify-center text-white font-bold text-sm shadow-md shadow-indigo-500/20">
            SM
          </div>
          <span className="text-[#1E293B] font-bold text-base tracking-wide">Sentinel</span>
        </div>

        <h1 className="text-xl font-bold text-[#1E293B] mb-1">Parolni tiklash</h1>
        <p className="text-slate-500 text-xs font-mono mb-6">Email manzilingizga reset link yuboramiz</p>

        {/* Muvaffaqiyatli yuborilgandagi holat */}
        {sent ? (
          <div className="p-4 rounded-lg bg-emerald-50 border border-emerald-200 text-emerald-700 text-xs font-mono mb-4 leading-relaxed shadow-sm">
            <div className="font-bold flex items-center gap-1.5 mb-1">
              <span>✓</span> Reset link yuborildi.
            </div>
            Emailingizni tekshiring.
            <div className="text-slate-400 text-[11px] mt-3 pt-2 border-t border-dashed border-emerald-200/60">
              Login sahifasiga yo'naltirilmoqdasiz...
            </div>
          </div>
        ) : (
          <form onSubmit={handleSubmit} className="space-y-4">
            
            {/* Xatolik xabari */}
            {error && (
              <div className="px-4 py-3 rounded-lg bg-red-50 border border-red-200 text-red-600 text-xs font-mono flex items-center gap-2 shadow-sm">
                ⚠ {error}
              </div>
            )}

            {/* Email input maydoni */}
            <div>
              <label className="block text-xs font-semibold text-slate-500 uppercase tracking-widest mb-1.5">
                Email manzil
              </label>
              <div className="relative">
                <svg className="absolute left-3 top-1/2 -translate-y-1/2 text-indigo-500" width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2.5">
                  <path d="M4 4h16v16H4z" />
                  <path d="M22 6l-10 7L2 6" />
                </svg>
                <input
                  type="email"
                  value={email}
                  onChange={e => setEmail(e.target.value)}
                  placeholder="admin@sentinel.com"
                  className="w-full bg-[#F1F5F9] border border-[#E2E8F0] rounded-lg pl-10 pr-4 py-2.5 text-sm text-[#1E293B] font-mono placeholder-slate-400 outline-none focus:bg-white focus:border-indigo-500 focus:ring-1 focus:ring-indigo-500/20 transition"
                  required
                />
              </div>
            </div>

            {/* Tasdiqlash tugmasi */}
            <button
              type="submit"
              disabled={loading}
              className="w-full py-2.5 rounded-lg bg-gradient-to-r from-indigo-600 to-purple-600 text-white font-semibold text-sm hover:opacity-95 active:scale-[0.99] transition disabled:opacity-50 flex items-center justify-center gap-2 shadow-sm font-mono"
            >
              {loading && (
                <svg className="animate-spin w-4 h-4" viewBox="0 0 24 24" fill="none">
                  <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4"/>
                  <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8v8z"/>
                </svg>
              )}
              {loading ? 'Yuborilmoqda...' : 'Link yuborish'}
            </button>

          </form>
        )}

        {/* Loginga qaytish havolasi */}
        <div className="mt-6 text-center">
          <span
            onClick={() => navigate('/login')}
            className="text-xs text-slate-500 hover:text-indigo-600 font-semibold font-mono cursor-pointer transition flex items-center justify-center gap-1"
          >
            ← Loginga qaytish
          </span>
        </div>

      </div>
    </div>
  )
}