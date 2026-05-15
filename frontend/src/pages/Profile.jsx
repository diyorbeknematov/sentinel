import { useState, useEffect } from 'react'
import api from '../api/axios'

export default function Profile() {
  const [user, setUser] = useState(null)
  const [loading, setLoading] = useState(true)
  
  // Tahrirlash holatlari
  const [usernameForm, setUsernameForm] = useState({ username: '' })
  const [passwordForm, setPasswordForm] = useState({ currentPassword: '', newPassword: '', confirmPassword: '' })
  
  // Panel holatlari (ochish/yopish)
  const [editUsername, setEditUsername] = useState(false)
  const [editPassword, setEditPassword] = useState(false)

  // Parolni ko'rsatish/yashirish holatlari (har bir maydon uchun alohida)
  const [showCurrentPass, setShowCurrentPass] = useState(false)
  const [showNewPass, setShowNewPass] = useState(false)
  const [showConfirmPass, setShowConfirmPass] = useState(false)
  
  // Status xabarlari
  const [msg, setMsg] = useState({ type: '', text: '' })
  const [actionLoading, setActionLoading] = useState(false)

  useEffect(() => {
    fetchUser()
  }, [])

  const fetchUser = async () => {
    try {
      const res = await api.get('/me')
      setUser(res.data)
      setUsernameForm({ username: res.data.username || '' })
    } catch {
      setUser(null)
    } finally {
      setLoading(false)
    }
  }

  const handleUpdateUsername = async (e) => {
    e.preventDefault()
    if (!usernameForm.username) return
    
    setActionLoading(true)
    setMsg({ type: '', text: '' })
    
    try {
      await api.put('/profile/update-username', usernameForm)
      setMsg({ type: 'success', text: 'Username muvaffaqiyatli yangilandi!' })
      setEditUsername(false)
      fetchUser()
    } catch (err) {
      setMsg({ type: 'error', text: err.response?.data?.message || 'Username oʻzgartirishda xatolik' })
    } finally {
      setActionLoading(false)
    }
  }

  const handleUpdatePassword = async (e) => {
    e.preventDefault()
    if (!passwordForm.currentPassword || !passwordForm.newPassword || !passwordForm.confirmPassword) {
      setMsg({ type: 'error', text: 'Barcha maydonlarni toʻldiring' })
      return
    }
    if (passwordForm.newPassword !== passwordForm.confirmPassword) {
      setMsg({ type: 'error', text: 'Yangi parollar mos kelmadi' })
      return
    }

    setActionLoading(true)
    setMsg({ type: '', text: '' })

    try {
      await api.put('/profile/update-password', {
        current_password: passwordForm.currentPassword,
        new_password: passwordForm.newPassword
      })
      setMsg({ type: 'success', text: 'Parol muvaffaqiyatli yangilandi!' })
      setPasswordForm({ currentPassword: '', newPassword: '', confirmPassword: '' })
      setEditPassword(false)
      // Holatlarni qayta yopish
      setShowCurrentPass(false)
      setShowNewPass(false)
      setShowConfirmPass(false)
    } catch (err) {
      setMsg({ type: 'error', text: err.response?.data?.message || 'Eski parol notoʻgʻri yoki xatolik yuz berdi' })
    } finally {
      setActionLoading(false)
    }
  }

  if (loading) {
    return (
      <div className="min-h-screen bg-[#F8FAFC] flex items-center justify-center text-slate-500 font-mono text-sm">
        Yuklanmoqda...
      </div>
    )
  }

  if (!user) {
    return (
      <div className="min-h-screen bg-[#F8FAFC] flex items-center justify-center text-red-500 font-mono text-sm">
        Foydalanuvchi ma'lumotlarini yuklashda xatolik yuz berdi.
      </div>
    )
  }

  // Ko'zcha Ikonkasi Komponenti (Toza va ixcham kod uchun)
  const EyeIcon = ({ visible }) => visible ? (
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
  )

  return (
    <div className="min-h-screen p-6 font-sans">
      
      {/* Sahifa sarlavhasi */}
      <div className="max-w-xl mx-auto mb-6">
        <h2 className="text-2xl font-bold text-[#1E293B]">Profil sozlamalari</h2>
        <p className="text-slate-500 text-sm font-mono mt-1">Hisob ma'lumotlarini boshqarish</p>
      </div>

      <div className="max-w-xl mx-auto space-y-4">
        
        {/* Status Alert xabarlari */}
        {msg.text && (
          <div className={`px-4 py-3 rounded-lg border text-xs font-mono flex items-center gap-2 shadow-sm ${
            msg.type === 'success' ? 'bg-emerald-50 border-emerald-200 text-emerald-700' : 'bg-red-50 border-red-200 text-red-600'
          }`}>
            {msg.type === 'success' ? '✓' : '⚠'} {msg.text}
          </div>
        )}

        {/* ASOSIY PROFIL KARTASI */}
        <div className="bg-white border border-[#E2E8F0] rounded-xl p-6 shadow-sm">
          <div className="space-y-4">
            
            {/* Email */}
            <div className="border-b border-[#F1F5F9] pb-4">
              <span className="block text-xs font-bold text-slate-400 uppercase tracking-widest mb-1">Email manzil</span>
              <span className="text-sm text-[#1E293B] font-mono font-medium">{user.email}</span>
            </div>

            {/* Username qismi */}
            <div className="border-b border-[#F1F5F9] pb-4">
              <div className="flex justify-between items-center">
                <div>
                  <span className="block text-xs font-bold text-slate-400 uppercase tracking-widest mb-1">Username</span>
                  <span className="text-sm text-[#1E293B] font-mono font-semibold">{user.username}</span>
                </div>
                <button
                  onClick={() => {
                    setEditUsername(!editUsername)
                    setEditPassword(false)
                    setMsg({ type: '', text: '' })
                  }}
                  className="text-xs text-indigo-600 hover:text-indigo-700 font-semibold border border-[#E2E8F0] px-3 py-1.5 rounded-lg bg-white shadow-sm hover:bg-[#F8FAFC] transition"
                >
                  {editUsername ? 'Yopish' : 'Oʻzgartirish'}
                </button>
              </div>

              {/* Username o'zgartirish formasi */}
              {editUsername && (
                <form onSubmit={handleUpdateUsername} className="mt-4 pt-4 border-t border-dashed border-[#E2E8F0] space-y-3">
                  <div>
                    <label className="block text-xs font-medium text-slate-500 mb-1">Yangi username</label>
                    <input
                      type="text"
                      value={usernameForm.username}
                      onChange={e => setUsernameForm({ username: e.target.value })}
                      className="w-full max-w-sm bg-[#F1F5F9] border border-[#E2E8F0] rounded-lg px-3 py-2 text-sm text-[#1E293B] font-mono outline-none focus:bg-white focus:border-indigo-500 transition"
                      required
                    />
                  </div>
                  <button
                    type="submit"
                    disabled={actionLoading}
                    className="px-4 py-2 bg-indigo-600 text-white rounded-lg text-xs font-semibold hover:bg-indigo-700 transition disabled:opacity-50"
                  >
                    {actionLoading ? 'Saqlanmoqda...' : 'Username-ni saqlash'}
                  </button>
                </form>
              )}
            </div>

            {/* Xavfsizlik / Parol qismi */}
            <div>
              <div className="flex justify-between items-center">
                <div>
                  <span className="block text-xs font-bold text-slate-400 uppercase tracking-widest mb-1">Xavfsizlik</span>
                  <span className="text-sm text-slate-500 font-mono">Parol oxirgi marta himoyalangan</span>
                </div>
                <button
                  onClick={() => {
                    setEditPassword(!editPassword)
                    setEditUsername(false)
                    setMsg({ type: '', text: '' })
                  }}
                  className="text-xs text-indigo-600 hover:text-indigo-700 font-semibold border border-[#E2E8F0] px-3 py-1.5 rounded-lg bg-white shadow-sm hover:bg-[#F8FAFC] transition"
                >
                  {editPassword ? 'Yopish' : 'Parolni yangilash'}
                </button>
              </div>

              {/* Parol o'zgartirish formasi */}
              {editPassword && (
                <form onSubmit={handleUpdatePassword} className="mt-4 pt-4 border-t border-dashed border-[#E2E8F0] space-y-3.5 max-w-sm">
                  
                  {/* Joriy Parol */}
                  <div>
                    <label className="block text-xs font-medium text-slate-500 mb-1">Joriy (eski) parol</label>
                    <div className="relative">
                      <input
                        type={showCurrentPass ? 'text' : 'password'}
                        placeholder="••••••••"
                        value={passwordForm.currentPassword}
                        onChange={e => setPasswordForm({ ...passwordForm, currentPassword: e.target.value })}
                        className="w-full bg-[#F1F5F9] border border-[#E2E8F0] rounded-lg pl-3 pr-10 py-2 text-sm text-[#1E293B] font-mono outline-none focus:bg-white focus:border-indigo-500 transition"
                        required
                      />
                      <button
                        type="button"
                        onClick={() => setShowCurrentPass(!showCurrentPass)}
                        className="absolute right-3 top-1/2 -translate-y-1/2 text-slate-400 hover:text-slate-600 transition"
                      >
                        <EyeIcon visible={showCurrentPass} />
                      </button>
                    </div>
                  </div>

                  {/* Yangi Parol */}
                  <div>
                    <label className="block text-xs font-medium text-slate-500 mb-1">Yangi parol</label>
                    <div className="relative">
                      <input
                        type={showNewPass ? 'text' : 'password'}
                        placeholder="••••••••"
                        value={passwordForm.newPassword}
                        onChange={e => setPasswordForm({ ...passwordForm, newPassword: e.target.value })}
                        className="w-full bg-[#F1F5F9] border border-[#E2E8F0] rounded-lg pl-3 pr-10 py-2 text-sm text-[#1E293B] font-mono outline-none focus:bg-white focus:border-indigo-500 transition"
                        required
                      />
                      <button
                        type="button"
                        onClick={() => setShowNewPass(!showNewPass)}
                        className="absolute right-3 top-1/2 -translate-y-1/2 text-slate-400 hover:text-slate-600 transition"
                      >
                        <EyeIcon visible={showNewPass} />
                      </button>
                    </div>
                  </div>

                  {/* Yangi Parolni Takrorlash */}
                  <div>
                    <label className="block text-xs font-medium text-slate-500 mb-1">Yangi parolni takrorlang</label>
                    <div className="relative">
                      <input
                        type={showConfirmPass ? 'text' : 'password'}
                        placeholder="••••••••"
                        value={passwordForm.confirmPassword}
                        onChange={e => setPasswordForm({ ...passwordForm, confirmPassword: e.target.value })}
                        className="w-full bg-[#F1F5F9] border border-[#E2E8F0] rounded-lg pl-3 pr-10 py-2 text-sm text-[#1E293B] font-mono outline-none focus:bg-white focus:border-indigo-500 transition"
                        required
                      />
                      <button
                        type="button"
                        onClick={() => setShowConfirmPass(!showConfirmPass)}
                        className="absolute right-3 top-1/2 -translate-y-1/2 text-slate-400 hover:text-slate-600 transition"
                      >
                        <EyeIcon visible={showConfirmPass} />
                      </button>
                    </div>
                  </div>

                  <button
                    type="submit"
                    disabled={actionLoading}
                    className="px-4 py-2 mt-2 bg-indigo-600 text-white rounded-lg text-xs font-semibold hover:bg-indigo-700 transition disabled:opacity-50 shadow-sm"
                  >
                    {actionLoading ? 'Yangilanmoqda...' : 'Parolni oʻzgartirish'}
                  </button>
                </form>
              )}
            </div>

          </div>
        </div>

      </div>
    </div>
  )
}