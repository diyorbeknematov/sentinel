import { useState, useEffect } from 'react'
import { useNavigate } from 'react-router-dom'
import api from '../api/axios'

export default function Register() {
  const navigate = useNavigate()
  const [form, setForm] = useState({ username: '', email: '', password: '', confirmPassword: '' })
  const [showPass, setShowPass] = useState(false)
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState('')
  const [metrics, setMetrics] = useState({
    requests: 48392, errorRate: 2.7, alerts: 18, response: 142
  })

  useEffect(() => {
    const interval = setInterval(() => {
      setMetrics({
        requests: Math.floor(48000 + Math.random() * 2000),
        errorRate: parseFloat((2 + Math.random() * 2).toFixed(1)),
        alerts: Math.floor(15 + Math.random() * 8),
        response: Math.floor(120 + Math.random() * 80),
      })
    }, 3000)
    return () => clearInterval(interval)
  }, [])

  const handleSubmit = async (e) => {
    e.preventDefault()
    setError('')
    if (!form.username || !form.email || !form.password || !form.confirmPassword) {
      setError('Barcha maydonlarni toldiring')
      return
    }
    if (form.password !== form.confirmPassword) {
      setError('Parollar mos kelmadi')
      return
    }
    setLoading(true)
    try {
      const res = await api.post('/register', form)
      const token = res.data.token || res.data.access_token
      if (token) {
        localStorage.setItem('sentinel_token', token)
        navigate('/home')
      }
    } catch (err) {
      setError(err.response?.data?.message || "Username yoki parol noto'g'ri")
    } finally {
      setLoading(false)
    }
  }

  return (
    <div
      className="min-h-screen bg-[#080b14] flex items-center justify-center"
      style={{
        backgroundImage: `linear-gradient(rgba(99,102,241,0.05) 1px, transparent 1px),
                          linear-gradient(90deg, rgba(99,102,241,0.05) 1px, transparent 1px)`,
        backgroundSize: '40px 40px'
      }}
    >
      <div className="flex w-full max-w-3xl mx-6 rounded-2xl overflow-hidden border border-[#1e293b] shadow-2xl">

        {/* LEFT */}
        <div className="flex-1 bg-[#0d1120] flex flex-col justify-center px-10 py-12">

          <div className="flex items-center gap-3 mb-10">
            <div className="w-9 h-9 rounded-lg bg-gradient-to-br from-indigo-500 to-purple-600 flex items-center justify-center text-white font-bold text-sm">
              LM
            </div>
            <span className="text-white font-semibold text-base tracking-wide">LogMonitor</span>
          </div>

          <h1 className="text-2xl font-bold text-white mb-1">Ro'yhatdan o'tish</h1>
          <p className="text-slate-500 text-sm mb-7">Davom etish uchun tizimga kiring</p>

          {error && (
            <div className="mb-4 px-4 py-3 rounded-lg bg-red-500/10 border border-red-500/20 text-red-400 text-xs font-mono flex items-center gap-2">
              ⚠ {error}
            </div>
          )}

          <form onSubmit={handleSubmit} className="space-y-4">
            <div>
              <label className="block text-xs text-slate-500 uppercase tracking-widest mb-2">Username</label>
              <div className="relative">
                <svg className="absolute left-3 top-1/2 -translate-y-1/2 text-indigo-500 opacity-70" width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                  <circle cx="12" cy="8" r="4"/><path d="M4 20c0-4 3.6-7 8-7s8 3 8 7"/>
                </svg>
                <input
                  type="text"
                  value={form.username}
                  onChange={e => setForm({ ...form, username: e.target.value })}
                  placeholder="admin"
                  className="w-full bg-[#080b14] border border-[#1e293b] rounded-lg pl-10 pr-4 py-3 text-sm text-slate-200 font-mono placeholder-slate-700 outline-none focus:border-indigo-500 focus:ring-1 focus:ring-indigo-500/20 transition"
                />
              </div>
            </div>

            <div>
              <label className="block text-xs text-slate-500 uppercase tracking-widest mb-2">
                Email
              </label>

              <div className="relative">
                {/* icon */}
                <svg
                  className="absolute left-3 top-1/2 -translate-y-1/2 text-indigo-500 opacity-70"
                  width="15"
                  height="15"
                  viewBox="0 0 24 24"
                  fill="none"
                  stroke="currentColor"
                  strokeWidth="2"
                >
                  <path d="M4 4h16v16H4z" />
                  <path d="M22 6l-10 7L2 6" />
                </svg>

                <input
                  type="email"
                  value={form.email}
                  onChange={e => setForm({ ...form, email: e.target.value })}
                  placeholder="admin@sentinel.com"
                  className="w-full bg-[#080b14] border border-[#1e293b] rounded-lg pl-10 pr-4 py-3 text-sm text-slate-200 font-mono placeholder-slate-700 outline-none focus:border-indigo-500 focus:ring-1 focus:ring-indigo-500/20 transition"
                />
              </div>
            </div>

            <div>
              <label className="block text-xs text-slate-500 uppercase tracking-widest mb-2">Password</label>
              <div className="relative">
                <svg className="absolute left-3 top-1/2 -translate-y-1/2 text-indigo-500 opacity-70" width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                  <rect x="3" y="11" width="18" height="11" rx="2"/><path d="M7 11V7a5 5 0 0 1 10 0v4"/>
                </svg>
                <input
                  type={showPass ? 'text' : 'password'}
                  value={form.password}
                  onChange={e => setForm({ ...form, password: e.target.value })}
                  placeholder="••••••••"
                  className="w-full bg-[#080b14] border border-[#1e293b] rounded-lg pl-10 pr-12 py-3 text-sm text-slate-200 font-mono placeholder-slate-700 outline-none focus:border-indigo-500 focus:ring-1 focus:ring-indigo-500/20 transition"
                />
                <button
                  type="button"
                  onClick={() => setShowPass(!showPass)}
                  className="absolute right-3 top-1/2 -translate-y-1/2 text-slate-600 hover:text-slate-400 transition"
                >
                  {showPass ? (
                    <svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                      <path d="M17.94 17.94A10.07 10.07 0 0 1 12 20c-7 0-11-8-11-8a18.45 18.45 0 0 1 5.06-5.94"/>
                      <path d="M9.9 4.24A9.12 9.12 0 0 1 12 4c7 0 11 8 11 8a18.5 18.5 0 0 1-2.16 3.19"/>
                      <line x1="1" y1="1" x2="23" y2="23"/>
                    </svg>
                  ) : (
                    <svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                      <path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z"/>
                      <circle cx="12" cy="12" r="3"/>
                    </svg>
                  )}
                </button>
              </div>
            </div>

            <div>
              <label className="block text-xs text-slate-500 uppercase tracking-widest mb-2">Confirm Password</label>
              <div className="relative">
                <svg className="absolute left-3 top-1/2 -translate-y-1/2 text-indigo-500 opacity-70" width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                  <rect x="3" y="11" width="18" height="11" rx="2"/><path d="M7 11V7a5 5 0 0 1 10 0v4"/>
                </svg>
                <input
                  type={showPass ? 'text' : 'password'}
                  value={form.confirmPassword}
                  onChange={e => setForm({ ...form, confirmPassword: e.target.value })}
                  placeholder="••••••••"
                  className="w-full bg-[#080b14] border border-[#1e293b] rounded-lg pl-10 pr-12 py-3 text-sm text-slate-200 font-mono placeholder-slate-700 outline-none focus:border-indigo-500 focus:ring-1 focus:ring-indigo-500/20 transition"
                />
                <button
                  type="button"
                  onClick={() => setShowPass(!showPass)}
                  className="absolute right-3 top-1/2 -translate-y-1/2 text-slate-600 hover:text-slate-400 transition"
                >
                  {showPass ? (
                    <svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                      <path d="M17.94 17.94A10.07 10.07 0 0 1 12 20c-7 0-11-8-11-8a18.45 18.45 0 0 1 5.06-5.94"/>
                      <path d="M9.9 4.24A9.12 9.12 0 0 1 12 4c7 0 11 8 11 8a18.5 18.5 0 0 1-2.16 3.19"/>
                      <line x1="1" y1="1" x2="23" y2="23"/>
                    </svg>
                  ) : (
                    <svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                      <path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z"/>
                      <circle cx="12" cy="12" r="3"/>
                    </svg>
                  )}
                </button>
              </div>
            </div>

            <button
              type="submit"
              disabled={loading}
              className="w-full py-3 rounded-lg bg-gradient-to-r from-indigo-600 to-violet-600 text-white font-semibold text-sm hover:opacity-90 active:scale-[0.99] transition disabled:opacity-50 flex items-center justify-center gap-2"
            >
              {loading && (
                <svg className="animate-spin w-4 h-4" viewBox="0 0 24 24" fill="none">
                  <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4"/>
                  <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8v8z"/>
                </svg>
              )}
              {loading ? 'Tekshirilmoqda...' : 'Create Account'}
            </button>
            
            <p className="text-xs text-slate-500 text-center mt-4">
              Account bormi?{' '}
              <span
                onClick={() => navigate('/login')}
                className="text-indigo-400 hover:text-indigo-300 cursor-pointer"
              >
                Login
              </span>
            </p>
          </form>
        </div>

        {/* RIGHT */}
        <div className="w-64 bg-[#080b14] border-l border-[#1e293b] flex flex-col justify-center px-7 py-10">

          <div className="flex items-center gap-2 mb-6">
            <span className="w-2 h-2 rounded-full bg-emerald-400 animate-pulse"></span>
            <span className="text-xs text-emerald-400 font-mono">System online</span>
          </div>

          <p className="text-xs text-slate-600 uppercase tracking-widest mb-3">Live metrics</p>
          <div className="space-y-0 mb-7">
            {[
              { label: 'Nginx req/s',   value: metrics.requests.toLocaleString(), color: 'text-indigo-400' },
              { label: 'Error rate',    value: metrics.errorRate + '%',            color: 'text-red-400' },
              { label: 'Active alerts', value: metrics.alerts,                     color: 'text-amber-400' },
              { label: 'Avg response',  value: metrics.response + 'ms',            color: 'text-emerald-400' },
            ].map(m => (
              <div key={m.label} className="flex justify-between items-center py-2.5 border-b border-[#1e293b]/50">
                <span className="text-xs text-slate-500">{m.label}</span>
                <span className={`text-xs font-mono font-semibold ${m.color}`}>{m.value}</span>
              </div>
            ))}
          </div>

          <p className="text-xs text-slate-600 uppercase tracking-widest mb-3">Recent alerts</p>
          <div className="space-y-2">
            {[
              { text: 'SSH brute force',    time: '14:32', color: 'bg-red-500/10 border-red-500/15',        dot: 'bg-red-400' },
              { text: '500 Internal Error', time: '14:29', color: 'bg-amber-500/10 border-amber-500/15',    dot: 'bg-amber-400' },
              { text: 'Rate limit exceeded',time: '14:25', color: 'bg-indigo-500/10 border-indigo-500/15',  dot: 'bg-indigo-400' },
            ].map(a => (
              <div key={a.text} className={`flex items-center gap-2 px-3 py-2 rounded-lg border ${a.color}`}>
                <span className={`w-1.5 h-1.5 rounded-full flex-shrink-0 ${a.dot}`}></span>
                <span className="text-xs text-slate-500 font-mono flex-1 truncate">{a.text}</span>
                <span className="text-xs text-slate-700">{a.time}</span>
              </div>
            ))}
          </div>
        </div>

      </div>
    </div>
  )
}