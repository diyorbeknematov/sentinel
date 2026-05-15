import { useState, useEffect } from 'react'
import { useNavigate } from 'react-router-dom'
import api from '../api/axios'

const S = {
  page: { 
    minHeight: '100vh',
    padding: '0' 
  },
  navbar: {
    position: 'sticky', 
    top: 0, 
    zIndex: 30,
    background: 'linear-gradient(90deg, rgba(79,70,229,0.08), rgba(16,185,129,0.05))',
    backdropFilter: 'blur(12px)',
    borderBottom: '1px solid rgba(148,163,184,0.2)',
    height: '56px', 
    display: 'flex', 
    alignItems: 'center', 
    justifyContent: 'space-between',
    marginLeft: '-24px', 
    marginRight: '-24px',
    paddingLeft: '24px', 
    paddingRight: '24px', 
    marginBottom: '32px',
  },
  card: {
    background: '#FFFFFF', 
    border: '1px solid #E2E8F0',
    borderRadius: '12px', 
    padding: '20px 24px',
    display: 'flex', 
    gap: '18px', 
    alignItems: 'flex-start',
    boxShadow: '0 1px 3px rgba(0,0,0,0.02), 0 1px 2px rgba(0,0,0,0.04)',
  },
  stepNum: {
    width: '28px', 
    height: '28px', 
    borderRadius: '50%',
    background: '#EEF2FF', 
    border: '1px solid #C7D2FE',
    display: 'flex', 
    alignItems: 'center', 
    justifyContent: 'center',
    fontSize: '13px', 
    fontWeight: '700', 
    color: '#4F46E5',
    flexShrink: 0, 
    marginTop: '2px', 
    fontFamily: 'monospace',
  },
  codeBlock: {
    background: '#F1F5F9', // Kod bloki uchun och terminal foni
    border: '1px solid #E2E8F0',
    borderRadius: '8px',
    padding: '14px 16px',
    fontFamily: 'monospace',
    fontSize: '13px',
    color: '#334155',
    marginTop: '12px',
    position: 'relative',
    lineHeight: '1.8',
  },
  keyBox: {
    display: 'flex', 
    alignItems: 'center', 
    gap: '8px',
    background: '#F1F5F9', 
    border: '1px solid #E2E8F0',
    borderRadius: '8px', 
    padding: '10px 14px', 
    marginTop: '12px',
  },
  iconBtn: {
    background: 'none', 
    border: 'none',
    cursor: 'pointer', 
    padding: '4px',
    display: 'flex', 
    alignItems: 'center', 
    justifyContent: 'center',
    color: '#94A3B8', 
    flexShrink: 0, 
    borderRadius: '4px',
    transition: 'color 0.15s',
  },
}

const CopyIcon = ({ done }) => done ? (
  <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="#16A34A" strokeWidth="2.5" strokeLinecap="round" strokeLinejoin="round">
    <polyline points="20 6 9 17 4 12"/>
  </svg>
) : (
  <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
    <rect x="9" y="9" width="13" height="13" rx="2" ry="2"/>
    <path d="M5 15H4a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2h9a2 2 0 0 1 2 2v1"/>
  </svg>
)

function CopyButton({ getText, style = {} }) {
  const [done, setDone] = useState(false)

  const handle = () => {
    navigator.clipboard.writeText(getText()).then(() => {
      setDone(true)
      setTimeout(() => setDone(false), 1500)
    })
  }

  return (
    <button
      onClick={handle}
      title={done ? 'Nusxalandi!' : 'Nusxalash'}
      style={{
        position: 'absolute', 
        right: '10px', 
        top: '10px',
        background: done ? '#DCFCE7' : '#FFFFFF',
        border: `1px solid ${done ? '#BBF7D0' : '#E2E8F0'}`,
        borderRadius: '5px', 
        padding: '5px 8px',
        cursor: 'pointer', 
        display: 'flex', 
        alignItems: 'center',
        boxShadow: '0 1px 2px rgba(0,0,0,0.02)',
        transition: 'all 0.2s', 
        ...style,
      }}
    >
      <CopyIcon done={done} />
    </button>
  )
}

function InlineCopy({ getText }) {
  const [done, setDone] = useState(false)

  const handle = () => {
    navigator.clipboard.writeText(getText()).then(() => {
      setDone(true)
      setTimeout(() => setDone(false), 1500)
    })
  }

  return (
    <button
      onClick={handle}
      title={done ? 'Nusxalandi!' : 'Nusxalash'}
      style={{
        ...S.iconBtn,
        color: done ? '#16A34A' : '#94A3B8',
      }}
      onMouseEnter={e => !done && (e.currentTarget.style.color = '#475569')}
      onMouseLeave={e => !done && (e.currentTarget.style.color = '#94A3B8')}
    >
      <CopyIcon done={done} />
    </button>
  )
}

export default function Home() {
  const navigate = useNavigate()
  const [user,       setUser]       = useState(null)
  const [agentCount, setAgentCount] = useState(0)
  const [keyVisible, setKeyVisible] = useState(false)

  useEffect(() => {
    const fetchData = async () => {
      try {
        const [uRes, aRes] = await Promise.allSettled([
          api.get('/me'),
          api.get('/agents'),
        ])
        if (uRes.status === 'fulfilled') setUser(uRes.value?.data)
        if (aRes.status === 'fulfilled') setAgentCount((aRes.value?.data?.total || 0))
      } catch { /* ignore */ }
    }
    fetchData()
  }, [])

  const apiKey    = user?.api_key || 'sk-xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx'
  const maskedKey = apiKey.slice(0, 5) + '•'.repeat(24)
  const serverURL = window.location.origin

  const installCmd  = `curl -sSL ${serverURL}/install.sh | bash`

  const steps = [
    {
      title: 'API kalitingizni oling',
      desc:  'Ushbu kalit serveringizni identifikatsiya qiladi. Hech kim bilan ulashmang.',
      content: (
        <div style={S.keyBox}>
          <span style={{
            fontFamily: 'monospace',
            fontSize: '13px',
            color: '#1E293B',
            flex: 1,
            letterSpacing: '0.5px',
            wordBreak: 'break-all',
            fontWeight: '500',
          }}>
            {keyVisible ? apiKey : maskedKey}
          </span>
          <button
            onClick={() => setKeyVisible(v => !v)}
            style={{
              fontSize: '11px',
              color: '#475569',
              background: '#FFFFFF',
              border: '1px solid #E2E8F0',
              padding: '4px 10px',
              borderRadius: '5px',
              cursor: 'pointer',
              fontFamily: 'monospace',
              fontWeight: '500',
              flexShrink: 0,
              boxShadow: '0 1px 2px rgba(0,0,0,0.02)',
            }}
            onMouseEnter={e => e.currentTarget.style.background = '#F8FAFC'}
            onMouseLeave={e => e.currentTarget.style.background = '#FFFFFF'}
          >
            {keyVisible ? 'Hide key' : 'Show key'}
          </button>

          <InlineCopy getText={() => apiKey} />
        </div>
      ),
    },
    {
      title: "Agentni serveringizga o'rnating",
      desc:  'Linux serveringizda terminalga quyidagi buyruqni kiriting:',
      content: (
        <div style={S.codeBlock}>
          <CopyButton getText={() => installCmd} />
          <span style={{ color: '#94A3B8', fontWeight: '600' }}>$ </span>
          <span style={{ color: '#7C3AED', fontWeight: '600' }}>curl</span>
          <span> -sSL {serverURL}/install.sh | </span>
          <span style={{ color: '#7C3AED', fontWeight: '600' }}>bash</span>
        </div>
      ),
    },
    {
      title: 'Agentni ulang',
      desc:  "O'rnatishdan so'ng API kalit bilan ishga tushiring:",
      content: (
        <div style={S.codeBlock}>
          <CopyButton getText={() =>
            `sentinel start --api-key ${apiKey} --server ${serverURL}`
          } />
          <div>
            <span style={{ color: '#94A3B8', fontWeight: '600' }}>$ </span>
            <span style={{ color: '#7C3AED', fontWeight: '600' }}>sentinel</span>
            <span style={{ color: '#475569' }}> start \</span>
          </div>
          <div style={{ paddingLeft: '16px' }}>
            <span style={{ color: '#059669', fontWeight: '600' }}>--api-key </span>
            <span style={{ color: '#D97706', fontWeight: '500' }}>
              {keyVisible ? apiKey : maskedKey}
            </span>
            <span style={{ color: '#475569' }}> \</span>
          </div>
          <div style={{ paddingLeft: '16px' }}>
            <span style={{ color: '#059669', fontWeight: '600' }}>--server </span>
            <span style={{ color: '#2563EB' }}>{serverURL}</span>
          </div>
        </div>
      ),
    },
    {
      title: 'Dashboard da kuzating',
      desc:  "Agent ulangandan so'ng 1–2 daqiqa ichida ma'lumotlar chiqadi.",
      content: (
        <div style={{ marginTop: '14px', display: 'flex', gap: '10px' }}>
          <button
            onClick={() => navigate('/dashboard')}
            style={{
              fontSize: '12px', 
              padding: '8px 18px', 
              borderRadius: '6px',
              border: '1px solid #4F46E5',
              background: '#EEF2FF', 
              color: '#4F46E5',
              cursor: 'pointer', 
              fontFamily: 'monospace',
              fontWeight: '600',
              boxShadow: '0 1px 2px rgba(0,0,0,0.02)',
            }}
          >
            Dashboard →
          </button>
          <button
            onClick={() => navigate('/agents')}
            style={{
              fontSize: '12px', 
              padding: '8px 18px', 
              borderRadius: '6px',
              border: '1px solid #E2E8F0', 
              background: '#FFFFFF',
              color: '#475569', 
              cursor: 'pointer', 
              fontFamily: 'monospace',
              fontWeight: '500',
              boxShadow: '0 1px 2px rgba(0,0,0,0.02)',
            }}
            onMouseEnter={e => e.currentTarget.style.background = '#F8FAFC'}
            onMouseLeave={e => e.currentTarget.style.background = '#FFFFFF'}
          >
            Agentlarni ko'rish
          </button>
        </div>
      ),
    },
  ]

  return (
    <div style={S.page}>

      {/* Navbar */}
      <div style={S.navbar}>
        <div style={{ display: 'flex', alignItems: 'center', gap: '8px' }}>
          <span style={{ fontSize: '16px', fontWeight: '600', color: '#1E293B' }}>Sentinel</span>
          <span style={{ fontSize: '12px', color: '#64748B', fontFamily: 'monospace', fontWeight: '500' }}>· Boshlash</span>
        </div>
        <div style={{ display: 'flex', alignItems: 'center', gap: '12px' }}>
          {agentCount > 0 && (
            <span style={{
              fontSize: '12px', 
              padding: '4px 12px', 
              borderRadius: '20px',
              background: '#DCFCE7', 
              border: '1px solid #BBF7D0',
              color: '#15803D', 
              fontFamily: 'monospace',
              fontWeight: '600',
            }}>
              {agentCount} agent ulangan
            </span>
          )}
          {user?.username && (
            <span style={{ fontSize: '13px', color: '#475569', fontFamily: 'monospace', fontWeight: '500' }}>
              {user.username}
            </span>
          )}
        </div>
      </div>

      {/* Content */}
      <div style={{ maxWidth: '820px', margin: '0 auto', padding: '0 24px 48px' }}>

        <div style={{ marginBottom: '28px' }}>
          <h1 style={{ 
            fontSize: '24px', 
            fontWeight: '700', 
            color: '#1E293B', 
            margin: '0 0 6px' 
          }}>
            Xush kelibsiz{user?.username ? `, ${user.username}` : ''} 👋
          </h1>
          <p style={{ fontSize: '13px', color: '#64748B', margin: 0, fontFamily: 'monospace', fontWeight: '500' }}>
            Monitoring boshlash uchun serveringizga agent o'rnating
          </p>
        </div>

        <div style={{ display: 'flex', flexDirection: 'column', gap: '14px' }}>
          {steps.map((step, i) => (
            <div key={i} style={S.card}>
              <div style={S.stepNum}>{i + 1}</div>
              <div style={{ flex: 1, minWidth: 0 }}>
                <p style={{ 
                  fontSize: '15px', 
                  fontWeight: '600', 
                  color: '#1E293B', 
                  margin: '0 0 4px' }}>
                  {step.title}
                </p>
                <p style={{ 
                  fontSize: '13px', 
                  color: '#64748B', 
                  margin: 0, 
                  fontFamily: 'monospace' }}>
                  {step.desc}
                </p>
                {step.content}
              </div>
            </div>
          ))}
        </div>

        {/* Footer */}
        <div style={{
          marginTop: '24px', 
          padding: '16px 20px',
          background: '#FFFFFF', 
          border: '1px solid #E2E8F0',
          borderRadius: '12px', 
          display: 'flex', gap: '24px', 
          flexWrap: 'wrap',
          boxShadow: '0 1px 3px rgba(0,0,0,0.02)',
        }}>
          <div>
            <p style={{ fontSize: '10px', color: '#64748B', textTransform: 'uppercase', letterSpacing: '0.5px', margin: '0 0 6px', fontWeight: '700' }}>
              Qo'llab-quvvatlanadigan OS
            </p>
            <div style={{ display: 'flex', gap: '6px', flexWrap: 'wrap' }}>
              {['Ubuntu 20+', 'Debian 11+', 'CentOS 8+'].map(os => (
                <span key={os} style={{
                  fontSize: '11px', 
                  padding: '3px 8px', 
                  borderRadius: '4px',
                  background: '#EEF2FF', 
                  color: '#4F46E5', 
                  fontFamily: 'monospace',
                  fontWeight: '600',
                }}>
                  {os}
                </span>
              ))}
            </div>
          </div>
          <div>
            <p style={{ fontSize: '10px', color: '#64748B', textTransform: 'uppercase', letterSpacing: '0.5px', margin: '0 0 6px', fontWeight: '700' }}>
              Agent holati
            </p>
            <span style={{
              fontSize: '11px', 
              padding: '3px 10px', 
              borderRadius: '4px',
              background: agentCount > 0 ? '#DCFCE7' : '#F1F5F9',
              color: agentCount > 0 ? '#15803D' : '#64748B', 
              fontFamily: 'monospace',
              fontWeight: '600',
              border: `1px solid ${agentCount > 0 ? '#BBF7D0' : '#E2E8F0'}`,
            }}>
              {agentCount > 0 ? `${agentCount} ta ulangan` : "Hali agent yo'q"}
            </span>
          </div>
        </div>

      </div>
    </div>
  )
}