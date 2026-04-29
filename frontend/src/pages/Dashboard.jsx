import { useState, useEffect } from 'react'
import api from '../api/axios'

const mockAlerts = [
  { id:'1', severity:'critical', type:'auth_failure', message:'SSH brute force — 120 failed attempts', agent_id:'nginx-01', created_at:'2024-01-25T14:32:00Z', is_read:false },
  { id:'2', severity:'warning',  type:'high_memory',  message:'RAM usage exceeded 90% threshold',      agent_id:'app-01',   created_at:'2024-01-25T14:29:00Z', is_read:false },
  { id:'3', severity:'high',     type:'rate_limit',   message:'Rate limit exceeded: 1200 req/min',     agent_id:'nginx-01', created_at:'2024-01-25T14:25:00Z', is_read:false },
  { id:'4', severity:'critical', type:'auth_failure', message:'Failed login from unknown IP',          agent_id:'nginx-02', created_at:'2024-01-25T14:20:00Z', is_read:true  },
]

const SEV = {
  critical: { bg:'rgba(239,68,68,0.08)',  border:'#ef4444', badge:'rgba(239,68,68,0.15)',  badgeText:'#fca5a5', text:'#fca5a5', dot:'#ef4444' },
  warning:  { bg:'rgba(245,158,11,0.08)', border:'#f59e0b', badge:'rgba(245,158,11,0.15)', badgeText:'#fcd34d', text:'#fcd34d', dot:'#f59e0b' },
  high:     { bg:'rgba(168,85,247,0.08)', border:'#a855f7', badge:'rgba(168,85,247,0.15)', badgeText:'#d8b4fe', text:'#d8b4fe', dot:'#a855f7' },
}

const TRACK = { red:'#ef4444', amber:'#f59e0b', blue:'#6366f1' }

function formatTime(iso) {
  return new Date(iso).toLocaleTimeString('en', { hour:'2-digit', minute:'2-digit', hour12:false })
}

function Gauge({ value, color }) {
  const circ = Math.PI * 23
  const offset = circ - (value / 100) * circ
  const stroke = TRACK[color]
  const fill = color === 'red' ? '#fca5a5' : color === 'amber' ? '#fcd34d' : '#a5b4fc'
  return (
    <svg width="64" height="38" viewBox="0 0 64 38">
      <path d="M6 34 A26 26 0 0 1 58 34" fill="none" stroke="#1e293b" strokeWidth="7" strokeLinecap="round"/>
      <path d="M6 34 A26 26 0 0 1 58 34" fill="none" stroke={stroke} strokeWidth="7"
        strokeLinecap="round" strokeDasharray={circ} strokeDashoffset={offset}/>
      <text x="32" y="32" textAnchor="middle" fontSize="10" fontWeight="600"
        fill={fill} fontFamily="monospace">{value}%</text>
    </svg>
  )
}

const S = {
  card: {
    background: '#0d1120',
    border: '1px solid #1e293b',
    borderRadius: '12px',
    padding: '16px',
  },
}

export default function Dashboard() {
  const [alerts, setAlerts]   = useState(mockAlerts)
  const [metrics, setMetrics] = useState({ cpu:67, ram:91, disk:50 })
  const [stats, setStats]     = useState({ requests:48392, errorRate:2.7, authFail:247, activeAlerts:18 })
  const [bars] = useState(() =>
    Array.from({ length: 20 }, () => ({
      nginx: Math.floor(35 + Math.random() * 55),
      app:   Math.floor(20 + Math.random() * 40),
      error: Math.random() > 0.75,
    }))
  )

  useEffect(() => {
    const load = async () => {
      try {
        const [mRes, aRes] = await Promise.all([
          api.get('/sentinel/api/metrics?limit=1'),
          api.get('/sentinel/api/alerts?limit=4'),
        ])
        if (mRes.data?.length) {
          const m = mRes.data[0]
          setMetrics({ cpu: Math.round(m.cpu), ram: Math.round(m.ram), disk: Math.round(m.disk) })
        }
        if (Array.isArray(aRes.data)) {
          setAlerts(aRes.data)
          setStats(p => ({ ...p, activeAlerts: aRes.data.length }))
        }
      } catch { /* mock */ }
    }
    load()
    const t = setInterval(load, 5000)
    return () => clearInterval(t)
  }, [])

  const metricRows = [
    { label:'CPU',  value:metrics.cpu,  color: metrics.cpu  >= 85 ? 'red' : metrics.cpu  >= 70 ? 'amber' : 'blue', note:'4 cores' },
    { label:'RAM',  value:metrics.ram,  color: metrics.ram  >= 85 ? 'red' : metrics.ram  >= 70 ? 'amber' : 'blue', note:`${(metrics.ram * 0.09).toFixed(1)} / 9 GB` },
    { label:'Disk', value:metrics.disk, color: metrics.disk >= 85 ? 'red' : metrics.disk >= 70 ? 'amber' : 'amber', note:`${Math.round(metrics.disk * 5)} / 500 GB` },
  ]

  const statCards = [
    { label:'TOTAL REQUESTS', value:stats.requests.toLocaleString(), sub:'▲ 2.4%', subOk:true,  note:'vs last hour', color:'#6366f1', valColor:'#a5b4fc' },
    { label:'ERROR RATE',     value:stats.errorRate.toFixed(1)+'%',  sub:'▼ 0.3%', subOk:true,  note:'improved',     color:'#ef4444', valColor:'#fca5a5' },
    { label:'AUTH FAILURES',  value:stats.authFail,                  sub:'▲ 12 new',subOk:false, note:'today',        color:'#f59e0b', valColor:'#fcd34d' },
    { label:'ACTIVE ALERTS',  value:stats.activeAlerts,              sub:'3 critical',subOk:false,note:'5 unread',    color:'#ef4444', valColor:'#fca5a5' },
  ]

  return (
    <div style={{ minHeight:'100vh' }}>

      {/* Top navbar */}
      <div style={{
        position: 'sticky', top: 0, zIndex: 30,
        background: 'rgba(8,11,20,0.85)',
        backdropFilter: 'blur(12px)',
        borderBottom: '1px solid #1e293b',
        padding: '0 24px',
        height: '52px',
        display: 'flex', alignItems: 'center', justifyContent: 'space-between',
        marginLeft: '-24px', marginRight: '-24px', marginBottom: '20px',
        paddingLeft: '24px', paddingRight: '24px',
      }}>
        <div style={{ display:'flex', alignItems:'center', gap:'8px' }}>
          <span style={{ fontSize:'14px', fontWeight:'600', color:'#e2e8f0' }}>Dashboard</span>
          <span style={{ fontSize:'11px', color:'#475569', fontFamily:'monospace' }}>
            · {new Date().toLocaleDateString('en', { month:'short', day:'numeric', year:'numeric' })}
          </span>
        </div>
        <div style={{
          display:'flex', alignItems:'center', gap:'6px',
          padding:'4px 12px', borderRadius:'20px',
          background:'rgba(34,197,94,0.08)', border:'1px solid rgba(34,197,94,0.2)',
        }}>
          <span style={{
            width:'6px', height:'6px', borderRadius:'50%',
            background:'#22c55e', display:'inline-block',
            animation:'pulse 2s infinite',
          }}></span>
          <span style={{ fontSize:'10px', color:'#4ade80', fontFamily:'monospace' }}>
            All systems operational
          </span>
        </div>
      </div>

      {/* Stat cards */}
      <div style={{ display:'grid', gridTemplateColumns:'repeat(4, 1fr)', gap:'12px', marginBottom:'16px' }}>
        {statCards.map(c => (
          <div key={c.label} style={{
            ...S.card,
            borderTop: `2px solid ${c.color}`,
          }}>
            <p style={{ fontSize:'9px', color:'#475569', letterSpacing:'0.7px', textTransform:'uppercase', marginBottom:'8px' }}>
              {c.label}
            </p>
            <p style={{ fontSize:'22px', fontWeight:'600', fontFamily:'monospace', color:c.valColor, lineHeight:1 }}>
              {c.value}
            </p>
            <div style={{ display:'flex', justifyContent:'space-between', marginTop:'8px' }}>
              <span style={{ fontSize:'10px', fontFamily:'monospace', color: c.subOk ? '#4ade80' : '#f87171' }}>
                {c.sub}
              </span>
              <span style={{ fontSize:'9px', color:'#475569' }}>{c.note}</span>
            </div>
          </div>
        ))}
      </div>

      {/* Metrics + Log volume */}
      <div style={{ display:'grid', gridTemplateColumns:'1fr 1fr', gap:'12px', marginBottom:'16px' }}>

        {/* Metrics — horizontal row */}
        <div style={S.card}>
          <p style={{ fontSize:'9px', color:'#475569', textTransform:'uppercase', letterSpacing:'0.7px', marginBottom:'16px' }}>
            Server Metrics
          </p>
          <div style={{ display:'grid', gridTemplateColumns:'repeat(3, 1fr)', gap:'8px' }}>
            {metricRows.map(m => (
              <div key={m.label} style={{
                background:'#111827', border:'1px solid #1e293b',
                borderRadius:'8px', padding:'12px 10px',
                display:'flex', flexDirection:'column', alignItems:'center', gap:'6px',
              }}>
                <Gauge value={m.value} color={m.color} />
                <p style={{ fontSize:'11px', fontWeight:'600', color:'#e2e8f0' }}>{m.label}</p>
                <div style={{
                  width:'100%', height:'4px', background:'#1e293b',
                  borderRadius:'4px', overflow:'hidden',
                }}>
                  <div style={{
                    width:`${m.value}%`, height:'100%',
                    background: TRACK[m.color], borderRadius:'4px',
                    transition:'width 0.5s',
                  }}/>
                </div>
                <p style={{ fontSize:'9px', color:'#475569', fontFamily:'monospace' }}>{m.note}</p>
              </div>
            ))}
          </div>
        </div>

        {/* Log volume chart */}
        <div style={S.card}>
          <div style={{ display:'flex', justifyContent:'space-between', alignItems:'center', marginBottom:'16px' }}>
            <p style={{ fontSize:'9px', color:'#475569', textTransform:'uppercase', letterSpacing:'0.7px' }}>
              Log Volume — Last 30 min
            </p>
            <div style={{ display:'flex', gap:'10px' }}>
              {[
                { label:'Nginx', color:'#818cf8' },
                { label:'App',   color:'#34d399' },
                { label:'Errors',color:'#f87171' },
              ].map(l => (
                <span key={l.label} style={{ fontSize:'9px', color:l.color, fontFamily:'monospace' }}>
                  ■ {l.label}
                </span>
              ))}
            </div>
          </div>
          <div style={{ display:'flex', alignItems:'flex-end', gap:'2px', height:'120px' }}>
            {bars.map((b, i) => (
              <div key={i} style={{ flex:1, display:'flex', flexDirection:'column', justifyContent:'flex-end', gap:'1px' }}>
                {b.error
                  ? <div style={{ borderRadius:'2px 2px 0 0', height:`${60 + Math.floor(Math.random()*25)}%`, background:'rgba(239,68,68,0.55)' }}/>
                  : <>
                      <div style={{ borderRadius:'2px 2px 0 0', height:`${b.nginx}%`, background:'rgba(99,102,241,0.6)' }}/>
                      <div style={{ borderRadius:'2px 2px 0 0', height:`${b.app}%`, background:'rgba(52,211,153,0.5)' }}/>
                    </>
                }
              </div>
            ))}
          </div>
          <div style={{ display:'flex', justifyContent:'space-between', marginTop:'8px' }}>
            {['-30min','-20min','-10min','now'].map(t => (
              <span key={t} style={{ fontSize:'9px', color:'#475569', fontFamily:'monospace' }}>{t}</span>
            ))}
          </div>
        </div>
      </div>

      {/* Recent alerts table */}
      <div style={S.card}>
        <div style={{ display:'flex', justifyContent:'space-between', alignItems:'center', marginBottom:'14px' }}>
          <p style={{ fontSize:'9px', color:'#475569', textTransform:'uppercase', letterSpacing:'0.7px' }}>
            Recent Alerts
          </p>
          <a href="/alerts" style={{ fontSize:'11px', color:'#818cf8', textDecoration:'none' }}>
            View all →
          </a>
        </div>

        <table style={{ width:'100%', borderCollapse:'collapse' }}>
          <thead>
            <tr style={{ borderBottom:'1px solid #1e293b' }}>
              {['', 'Severity', 'Type', 'Message', 'Agent', 'Time', ''].map((h, i) => (
                <th key={i} style={{
                  textAlign:'left', fontSize:'10px', color:'#475569',
                  fontWeight:'500', textTransform:'uppercase', letterSpacing:'0.5px',
                  padding:'0 10px 10px',
                }}>
                  {h}
                </th>
              ))}
            </tr>
          </thead>
          <tbody>
            {alerts.map((a, i) => {
              const s = SEV[a.severity] || SEV.warning
              return (
                <tr
                  key={i}
                  style={{ borderBottom:'1px solid rgba(30,41,59,0.5)', cursor:'default' }}
                  onMouseEnter={e => e.currentTarget.style.background='rgba(30,41,59,0.3)'}
                  onMouseLeave={e => e.currentTarget.style.background='transparent'}
                >
                  <td style={{ padding:'10px 10px' }}>
                    <span style={{ width:'6px', height:'6px', borderRadius:'50%', background:s.dot, display:'inline-block' }}/>
                  </td>
                  <td style={{ padding:'10px 10px' }}>
                    <span style={{
                      fontSize:'10px', padding:'2px 8px', borderRadius:'4px',
                      fontFamily:'monospace', fontWeight:'500',
                      background:s.badge, color:s.badgeText,
                    }}>
                      {a.severity}
                    </span>
                  </td>
                  <td style={{ padding:'10px 10px', fontSize:'11px', color:'#94a3b8', fontFamily:'monospace' }}>
                    {a.type}
                  </td>
                  <td style={{ padding:'10px 10px', fontSize:'11px', color:'#94a3b8', fontFamily:'monospace', maxWidth:'260px' }}>
                    <span style={{ display:'block', overflow:'hidden', textOverflow:'ellipsis', whiteSpace:'nowrap' }}>
                      {a.message}
                    </span>
                  </td>
                  <td style={{ padding:'10px 10px', fontSize:'11px', color:'#64748b', fontFamily:'monospace' }}>
                    {a.agent_id}
                  </td>
                  <td style={{ padding:'10px 10px', fontSize:'11px', color:'#475569', fontFamily:'monospace' }}>
                    {formatTime(a.created_at)}
                  </td>
                  <td style={{ padding:'10px 10px' }}>
                    {!a.is_read && (
                      <span style={{
                        width:'6px', height:'6px', borderRadius:'50%',
                        background:'#60a5fa', display:'inline-block',
                      }}/>
                    )}
                  </td>
                </tr>
              )
            })}
          </tbody>
        </table>
      </div>

    </div>
  )
}