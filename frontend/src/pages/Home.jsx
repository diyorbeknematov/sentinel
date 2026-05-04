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
    background: 'rgba(8,11,20,0.85)', 
    backdropFilter: 'blur(12px)',
    borderBottom: '1px solid #1e293b',
    height: '52px', 
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
    background: '#0d1120', 
    border: '1px solid #1e293b',
    borderRadius: '12px', 
    padding: '20px 24px',
    display: 'flex', 
    gap: '18px', 
    alignItems: 'flex-start',
  },
  stepNum: {
    width: '28px', 
    height: '28px', 
    borderRadius: '50%',
    background: 'rgba(99,102,241,0.1)', 
    border: '1px solid rgba(99,102,241,0.25)',
    display: 'flex', 
    alignItems: 'center', 
    justifyContent: 'center',
    fontSize: '12px', 
    fontWeight: '600', 
    color: '#818cf8',
    flexShrink: 0, 
    marginTop: '2px', 
    fontFamily: 'monospace',
  },
  codeBlock: {
    background: '#080b14',
    border: '1px solid #1e293b',
    borderRadius: '8px',
    padding: '12px 16px',
    fontFamily: 'monospace',
    fontSize: '12px',
    color: '#94a3b8',
    marginTop: '12px',
    position: 'relative',
    lineHeight: '1.8',
  },
  keyBox: {
    display: 'flex', 
    alignItems: 'center', 
    gap: '8px',
    background: '#080b14', 
    border: '1px solid #1e293b',
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
    color: '#475569', 
    flexShrink: 0, 
    borderRadius: '4px',
    transition: 'color 0.15s',
  },
}

const CopyIcon = ({ done }) => done ? (
  <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="#4ade80" strokeWidth="2.5" strokeLinecap="round" strokeLinejoin="round">
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
        background: done ? 'rgba(34,197,94,0.08)' : 'rgba(30,41,59,0.7)',
        border: `1px solid ${done ? 'rgba(34,197,94,0.25)' : '#1e293b'}`,
        borderRadius: '5px', 
        padding: '5px 7px',
        cursor: 'pointer', 
        display: 'flex', 
        alignItems: 'center',
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
        color: done ? '#4ade80' : '#475569',
      }}
      onMouseEnter={e => !done && (e.currentTarget.style.color = '#94a3b8')}
      onMouseLeave={e => !done && (e.currentTarget.style.color = '#475569')}
    >
      <CopyIcon done={done} />
    </button>
  )
}

// Main 
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
  const connectCmd  = `sentinel start \\\n  --api-key ${apiKey} \\\n  --server  ${serverURL}`

  const steps = [
    {
      title: 'API kalitingizni oling',
      desc:  'Ushbu kalit serveringizni identifikatsiya qiladi. Hech kim bilan ulashmang.',
      content: (
        <div style={S.keyBox}>
          <span style={{
            fontFamily: 'monospace',
            fontSize: '13px',
            color: '#e2e8f0',
            flex: 1,
            letterSpacing: '0.5px',
            wordBreak: 'break-all',
          }}>
            {keyVisible ? apiKey : maskedKey}
          </span>
          <button
            onClick={() => setKeyVisible(v => !v)}
            style={{
              fontSize: '10px',
              color: '#475569',
              background: 'rgba(30,41,59,0.6)',
              border: '1px solid #1e293b',
              padding: '3px 10px',
              borderRadius: '4px',
              cursor: 'pointer',
              fontFamily: 'monospace',
              flexShrink: 0,
            }}
          >
            {keyVisible ? 'Hide key' : 'Show key'}
          </button>

          {/* Nusxalash */}
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
          <span style={{ color: '#475569' }}>$ </span>
          <span style={{ color: '#c4b5fd' }}>curl</span>
          <span> -sSL {serverURL}/install.sh | </span>
          <span style={{ color: '#c4b5fd' }}>bash</span>
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
            <span style={{ color: '#475569' }}>$ </span>
            <span style={{ color: '#c4b5fd' }}>sentinel</span>
            <span style={{ color: '#94a3b8' }}> start \</span>
          </div>
          <div style={{ paddingLeft: '16px' }}>
            <span style={{ color: '#34d399' }}>--api-key </span>
            <span style={{ color: '#fcd34d' }}>
              {keyVisible ? apiKey : maskedKey}
            </span>
            <span style={{ color: '#94a3b8' }}> \</span>
          </div>
          <div style={{ paddingLeft: '16px' }}>
            <span style={{ color: '#34d399' }}>--server </span>
            <span style={{ color: '#a5b4fc' }}>{serverURL}</span>
          </div>
        </div>
      ),
    },
    {
      title: 'Dashboard da kuzating',
      desc:  "Agent ulangandan so'ng 1–2 daqiqa ichida ma'lumotlar chiqadi.",
      content: (
        <div style={{ marginTop: '12px', display: 'flex', gap: '10px' }}>
          <button
            onClick={() => navigate('/dashboard')}
            style={{
              fontSize: '12px', 
              padding: '7px 16px', 
              borderRadius: '6px',
              border: '1px solid rgba(99,102,241,0.35)',
              background: 'rgba(99,102,241,0.1)', 
              color: '#a5b4fc',
              cursor: 'pointer', 
              fontFamily: 'monospace',
            }}
          >
            Dashboard →
          </button>
          <button
            onClick={() => navigate('/agents')}
            style={{
              fontSize: '12px', 
              padding: '7px 16px', 
              borderRadius: '6px',
              border: '1px solid #1e293b', 
              background: 'transparent',
              color: '#64748b', 
              cursor: 'pointer', 
              fontFamily: 'monospace',
            }}
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
          <span style={{ fontSize: '14px', fontWeight: '600', color: '#e2e8f0' }}>Sentinel</span>
          <span style={{ fontSize: '11px', color: '#475569', fontFamily: 'monospace' }}>· Boshlash</span>
        </div>
        <div style={{ display: 'flex', alignItems: 'center', gap: '12px' }}>
          {agentCount > 0 && (
            <span style={{
              fontSize: '11px', 
              padding: '3px 10px', 
              borderRadius: '20px',
              background: 'rgba(34,197,94,0.08)', 
              border: '1px solid rgba(34,197,94,0.2)',
              color: '#4ade80', 
              fontFamily: 'monospace',
            }}>
              {agentCount} agent ulangan
            </span>
          )}
          {user?.username && (
            <span style={{ fontSize: '12px', color: '#475569', fontFamily: 'monospace' }}>
              {user.username}
            </span>
          )}
        </div>
      </div>

      {/* Content */}
      <div style={{ maxWidth: '640px', margin: '0 auto', padding: '0 24px 48px' }}>

        <div style={{ marginBottom: '28px' }}>
          <h1 style={{ 
            fontSize: '22px', 
            fontWeight: '600', 
            color: '#e2e8f0', 
            margin: '0 0 6px' 
          }}>
            Xush kelibsiz{user?.username ? `, ${user.username}` : ''} 👋
          </h1>
          <p style={{ fontSize: '13px', color: '#475569', margin: 0, fontFamily: 'monospace' }}>
            Monitoring boshlash uchun serveringizga agent o'rnating
          </p>
        </div>

        <div style={{ display: 'flex', flexDirection: 'column', gap: '10px' }}>
          {steps.map((step, i) => (
            <div key={i} style={S.card}>
              <div style={S.stepNum}>{i + 1}</div>
              <div style={{ flex: 1, minWidth: 0 }}>
                <p style={{ 
                  fontSize: '14px', 
                  fontWeight: '600', 
                  color: '#e2e8f0', 
                  margin: '0 0 4px' }}>
                  {step.title}
                </p>
                <p style={{ 
                  fontSize: '12px', 
                  color: '#475569', 
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
          background: 'rgba(99,102,241,0.04)', 
          border: '1px solid rgba(99,102,241,0.12)',
          borderRadius: '10px', 
          display: 'flex', gap: '24px', 
          flexWrap: 'wrap',
        }}>
          <div>
            <p style={{ fontSize: '10px', color: '#475569', textTransform: 'uppercase', letterSpacing: '0.5px', margin: '0 0 6px' }}>
              Qo'llab-quvvatlanadigan OS
            </p>
            <div style={{ display: 'flex', gap: '6px', flexWrap: 'wrap' }}>
              {['Ubuntu 20+', 'Debian 11+', 'CentOS 8+'].map(os => (
                <span key={os} style={{
                  fontSize: '10px', 
                  padding: '2px 8px', 
                  borderRadius: '4px',
                  background: 'rgba(99,102,241,0.12)', 
                  color: '#818cf8', 
                  fontFamily: 'monospace',
                }}>
                  {os}
                </span>
              ))}
            </div>
          </div>
          <div>
            <p style={{ fontSize: '10px', color: '#475569', textTransform: 'uppercase', letterSpacing: '0.5px', margin: '0 0 6px' }}>
              Agent holati
            </p>
            <span style={{
              fontSize: '10px', 
              padding: '2px 10px', 
              borderRadius: '4px',
              background: agentCount > 0 ? 'rgba(34,197,94,0.1)' : 'rgba(30,41,59,0.5)',
              color: agentCount > 0 ? '#4ade80' : '#475569', fontFamily: 'monospace',
              border: `1px solid ${agentCount > 0 ? 'rgba(34,197,94,0.2)' : '#1e293b'}`,
            }}>
              {agentCount > 0 ? `${agentCount} ta ulangan` : "Hali agent yo'q"}
            </span>
          </div>
        </div>

      </div>
    </div>
  )
}