import { useState, useEffect, useRef } from 'react'
import { Outlet, NavLink, useNavigate } from 'react-router-dom'
import api from '../api/axios'

const SIDEBAR_WIDTH = 220

const navItems = [
  {
    to: '/home', label: 'Home',
    icon: (
      <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
        <polyline points="22 12 18 12 15 21 9 3 6 12 2 12"/>
      </svg>
    ),
  },
  {
    to: '/', label: 'Dashboard',
    icon: <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2"><rect x="3" y="3" width="7" height="7"/><rect x="14" y="3" width="7" height="7"/><rect x="3" y="14" width="7" height="7"/><rect x="14" y="14" width="7" height="7"/></svg>,
  },
  {
    to: '/agents', label: 'Agents',
    icon: <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
          <rect x="4" y="4" width="16" height="6" rx="2"/>
          <rect x="4" y="14" width="16" height="6" rx="2"/>
          <circle cx="8" cy="7" r="1"/>
          <circle cx="8" cy="17" r="1"/>
        </svg>
  },
  {
    to: '/nginx-logs', label: 'Nginx Logs',
    icon: <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2"><line x1="8" y1="6" x2="21" y2="6"/><line x1="8" y1="12" x2="21" y2="12"/><line x1="8" y1="18" x2="21" y2="18"/><line x1="3" y1="6" x2="3.01" y2="6"/><line x1="3" y1="12" x2="3.01" y2="12"/><line x1="3" y1="18" x2="3.01" y2="18"/></svg>,
  },
  {
    to: '/app-logs', label: 'App Logs',
    icon: <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2"><path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"/><polyline points="14 2 14 8 20 8"/></svg>,
  },
  {
    to: '/alerts', label: 'Alerts', badge: true,
    icon: <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2"><path d="M18 8A6 6 0 0 0 6 8c0 7-3 9-3 9h18s-3-2-3-9"/><path d="M13.73 21a2 2 0 0 1-3.46 0"/></svg>,
  },
]

export default function Layout() {
  const navigate = useNavigate()
  const [unreadCount, setUnreadCount] = useState(0)
  const [openUser, setOpenUser] = useState(false)
  const [user, setUser] = useState(null)
  const userRef = useRef()

  useEffect(() => {
    const fetchUnread = async () => {
      try {
        const res = await api.get('/alerts?is_read=false')
        setUnreadCount(Number(res.data.total) || 0) // unread alertlar sonini total dan olamiz
      } catch {
        setUnreadCount(0)
      }
    }
    fetchUnread()
    const t = setInterval(fetchUnread, 10000)
    return () => clearInterval(t)
  }, [])

  useEffect(() => {
    const h = (e) => { if (userRef.current && !userRef.current.contains(e.target)) setOpenUser(false) }
    document.addEventListener('mousedown', h)
    return () => document.removeEventListener('mousedown', h)
  }, [])

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

  const logout = () => {
    localStorage.removeItem('sentinel_token')
    navigate('/login')
  }

  return (
    <div style={{
      minHeight: '100vh',
      display: 'flex',
      fontFamily: 'Inter, system-ui, sans-serif',
      background: '#080b14',
      backgroundImage: `
        linear-gradient(rgba(99,102,241,0.04) 1px, transparent 1px),
        linear-gradient(90deg, rgba(99,102,241,0.04) 1px, transparent 1px)
      `,
      backgroundSize: '40px 40px',
    }}>

      {/* SIDEBAR */}
      <aside style={{
        position: 'fixed', top: 0, left: 0,
        width: `${SIDEBAR_WIDTH}px`, height: '100vh',
        background: '#0d1120',
        borderRight: '1px solid #1e293b',
        display: 'flex', flexDirection: 'column',
        padding: '20px 12px',
        zIndex: 50,
      }}>

        {/* Logo */}
        <div style={{ display:'flex', alignItems:'center', gap:'10px', padding:'4px 8px', marginBottom:'28px' }}>
          <div style={{
            width:'32px', height:'32px', borderRadius:'8px',
            background:'rgba(99,102,241,0.15)',
            display:'flex', alignItems:'center', justifyContent:'center',
            fontSize:'12px', fontWeight:'700', color:'#a5b4fc',
          }}>LM</div>
          <div>
            <div style={{ fontSize:'14px', fontWeight:'600', color:'#e2e8f0' }}>LogMonitor</div>
            <div style={{ fontSize:'11px', color:'#475569' }}>Security Dashboard</div>
          </div>
        </div>

        {/* Nav */}
        <nav style={{ display:'flex', flexDirection:'column', gap:'2px', flex:1 }}>
          {navItems.map(item => (
            <NavLink
              key={item.to}
              to={item.to}
              end={item.to === '/'}
              style={({ isActive }) => ({
                display: 'flex', alignItems: 'center', gap: '10px',
                padding: '9px 10px', borderRadius: '8px',
                fontSize: '13px', textDecoration: 'none',
                transition: 'all 0.15s',
                background: isActive ? 'rgba(99,102,241,0.12)' : 'transparent',
                color: isActive ? '#a5b4fc' : '#64748b',
                border: isActive ? '1px solid rgba(99,102,241,0.2)' : '1px solid transparent',
              })}
            >
              {item.icon}
              <span style={{ flex:1 }}>{item.label}</span>
              {item.badge && unreadCount > 0 && (
                <span style={{
                  fontSize:'10px', padding:'1px 6px', borderRadius:'10px',
                  background:'rgba(239,68,68,0.15)', color:'#fca5a5',
                  fontFamily:'monospace',
                }}>{unreadCount}</span>
              )}
            </NavLink>
          ))}
        </nav>

        {/* User */}
        <div ref={userRef} style={{ position:'relative', borderTop:'1px solid #1e293b', paddingTop:'12px' }}>
          <div
            onClick={() => setOpenUser(!openUser)}
            style={{
              display:'flex', alignItems:'center', gap:'10px',
              padding:'8px 10px', borderRadius:'8px', cursor:'pointer',
              transition:'background 0.15s',
            }}
            onMouseEnter={e => e.currentTarget.style.background='rgba(255,255,255,0.04)'}
            onMouseLeave={e => e.currentTarget.style.background='transparent'}
          >
            <div style={{
              width:'30px', height:'30px', borderRadius:'50%', flexShrink:0,
              background:'rgba(168,85,247,0.15)',
              display:'flex', alignItems:'center', justifyContent:'center',
              fontSize:'11px', fontWeight:'600', color:'#c4b5fd',
            }}>A</div>
            <div style={{ flex:1, minWidth:0 }}>
              <div style={{ fontSize:'13px', fontWeight:'500', color:'#e2e8f0' }}>
                <div>{user?.username || 'Loading...' }</div>
              </div>
              <div style={{ fontSize:'11px', color:'#475569', overflow:'hidden', textOverflow:'ellipsis', whiteSpace:'nowrap' }}>
                <div>{user?.role || 'Admin'}</div>
              </div>
            </div>
            <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="#475569" strokeWidth="2">
              <polyline points="6 9 12 15 18 9"/>
            </svg>
          </div>

          {openUser && (
            <div style={{
              position:'absolute', bottom:'calc(100% + 6px)', left:0, right:0,
              background:'#0d1120', border:'1px solid #1e293b',
              borderRadius:'10px', padding:'6px', overflow:'hidden',
              boxShadow:'0 -8px 24px rgba(0,0,0,0.4)',
            }}>
              <div
                onClick={() => {
                  navigate('/profile')
                  setOpenUser(false)
                }}
                style={{
                  display:'flex', alignItems:'center', gap:'8px',
                  padding:'8px 12px', borderRadius:'6px',
                  fontSize:'13px', color:'#cbd5f5', cursor:'pointer',
                }}
                onMouseEnter={e => e.currentTarget.style.background='rgba(99,102,241,0.1)'}
                onMouseLeave={e => e.currentTarget.style.background='transparent'}
              >
                <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                  <circle cx="12" cy="7" r="4"/>
                  <path d="M5.5 21a6.5 6.5 0 0 1 13 0"/>
                </svg>
                Profile
              </div>

              <div
                onClick={logout}
                style={{
                  display:'flex', alignItems:'center', gap:'8px',
                  padding:'8px 12px', borderRadius:'6px',
                  fontSize:'13px', color:'#f87171', cursor:'pointer',
                }}
                onMouseEnter={e => e.currentTarget.style.background='rgba(239,68,68,0.1)'}
                onMouseLeave={e => e.currentTarget.style.background='transparent'}
              >
                <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                  <path d="M9 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h4"/>
                  <polyline points="16 17 21 12 16 7"/>
                  <line x1="21" y1="12" x2="9" y2="12"/>
                </svg>
                Log out
              </div>
            </div>
          )}
        </div>
      </aside>

      {/* MAIN */}
      <main style={{
        marginLeft: `${SIDEBAR_WIDTH}px`,
        flex: 1, minHeight: '100vh',
        display: 'flex', justifyContent: 'center',
      }}>
        <div style={{ width:'100%', maxWidth:'1200px', padding:'24px' }}>
          <Outlet />
        </div>
      </main>
    </div>
  )
}