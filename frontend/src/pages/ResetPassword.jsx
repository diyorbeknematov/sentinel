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
    <div 
      className="min-h-screen bg-[#F8FAFC] flex items-center justify-center font-sans"
      style={{
        backgroundImage: `linear-gradient(rgba(79,70,229,0.03) 1px, transparent 1px),
                          linear-gradient(90deg, rgba(79,70,229,0.03) 1px, transparent 1px)`,
        backgroundSize: '40px 40px'
      }}
    >
      <div className="w-full max-w-sm mx-6 bg-white border border-[#E2E8F0] rounded-2xl p-8 shadow-xl">

        {/* Logo / Brending */}
        <div className="flex items-center gap-3 mb-6">
          <div className="w-9 h-9 rounded-lg bg-gradient-to-br from-indigo-500 to-purple-600 flex items-center justify-center text-white font-bold text-sm shadow-md shadow-indigo-500/20">
            SM
          </div>
          <span className="text-[#1E293B] font-bold text-base tracking-wide">Sentinel</span>
        </div>

        <h1 className="text-xl font-bold text-[#1E293B] mb-1">
          Yangi parol o'rnatish
        </h1>
        <p className="text-slate-500 text-xs font-mono mb-6">
          Kamida 6 ta belgidan iborat parol kiriting
        </p>

        {/* Muvaffaqiyat paneli */}
        {success ? (
          <div className="p-4 rounded-lg bg-emerald-50 border border-emerald-200 text-emerald-700 text-xs font-mono mb-4 leading-relaxed shadow-sm">
            <div className="font-bold flex items-center gap-1.5 mb-1">
              <span>✓</span> Parol muvaffaqiyatli o'zgartirildi.
            </div>
            <div className="text-slate-400 text-[11px] mt-2 pt-2 border-t border-dashed border-emerald-200/60">
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

            {/* Yangi parol */}
            <div>
              <label className="block text-xs font-semibold text-slate-500 uppercase tracking-widest mb-1.5">
                Yangi parol
              </label>
              <div className="relative">
                <input
                  type={showPass ? 'text' : 'password'}
                  value={password}
                  onChange={e => setPassword(e.target.value)}
                  placeholder="••••••••"
                  className="w-full bg-[#F1F5F9] border border-[#E2E8F0] rounded-lg pl-3 pr-10 py-2.5 text-sm text-[#1E293B] font-mono placeholder-slate-400 outline-none focus:bg-white focus:border-indigo-500 focus:ring-1 focus:ring-indigo-500/20 transition"
                  required
                />
                <button 
                  type="button" 
                  onClick={() => setShowPass(v => !v)} 
                  className="absolute right-3 top-1/2 -translate-y-1/2 text-slate-400 hover:text-slate-600 transition"
                >
                  {showPass ? <EyeOffIcon /> : <EyeIcon />}
                </button>
              </div>
            </div>

            {/* Parolni tasdiqlash */}
            <div>
              <label className="block text-xs font-semibold text-slate-500 uppercase tracking-widest mb-1.5">
                Parolni tasdiqlang
              </label>
              <div className="relative">
                <input
                  type={showConf ? 'text' : 'password'}
                  value={confirm}
                  onChange={e => setConfirm(e.target.value)}
                  placeholder="••••••••"
                  className={`w-full bg-[#F1F5F9] border rounded-lg pl-3 pr-10 py-2.5 text-sm text-[#1E293B] font-mono placeholder-slate-400 outline-none focus:bg-white transition ${
                    confirm && confirm !== password 
                      ? 'border-red-300 focus:border-red-500 focus:ring-1 focus:ring-red-500/20' 
                      : 'border-[#E2E8F0] focus:border-indigo-500 focus:ring-1 focus:ring-indigo-500/20'
                  }`}
                  required
                />
                <button 
                  type="button" 
                  onClick={() => setShowConf(v => !v)} 
                  className="absolute right-3 top-1/2 -translate-y-1/2 text-slate-400 hover:text-slate-600 transition"
                >
                  {showConf ? <EyeOffIcon /> : <EyeIcon />}
                </button>
              </div>
              
              {/* Parollar moslik indikatori */}
              {confirm && (
                <p className={`text-[11px] font-mono mt-1.5 flex items-center gap-1 ${
                  confirm === password ? 'text-emerald-600' : 'text-red-500'
                }`}>
                  {confirm === password ? '✓ Parollar mos keldi' : '✗ Parollar mos emas'}
                </p>
              )}
            </div>

            {/* Tasdiqlash tugmasi */}
            <button
              type="submit"
              disabled={loading || !token}
              className="w-full py-2.5 rounded-lg bg-gradient-to-r from-indigo-600 to-purple-600 text-white font-semibold text-sm hover:opacity-95 active:scale-[0.99] transition disabled:opacity-50 flex items-center justify-center gap-2 shadow-sm font-mono"
            >
              {loading && (
                <svg className="animate-spin w-4 h-4" viewBox="0 0 24 24" fill="none">
                  <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4"/>
                  <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8v8z"/>
                </svg>
              )}
              {loading ? 'Saqlanmoqda...' : 'Parolni o\'rnatish'}
            </button>

          </form>
        )}

        {/* Orqaga qaytish */}
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