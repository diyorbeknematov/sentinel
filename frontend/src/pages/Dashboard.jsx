import { useState, useEffect, useCallback, useRef } from 'react'
import { AreaChart, Area, XAxis, YAxis, Tooltip, ResponsiveContainer } from 'recharts'
import api from '../api/axios'

// Och rangli dizayn (Light Mode) uchun xavfsizlik darajalari ranglari
const SEV = {
  critical: { bg:'rgba(254,226,226,0.4)',  border:'#fca5a5', badge:'#FEE2E2',  badgeText:'#DC2626', text:'#991B1B', dot:'#EF4444' },
  warning:  { bg:'rgba(254,243,199,0.4)', border:'#fcd34d', badge:'#FEF3C7', badgeText:'#D97706', text:'#92400E', dot:'#F59E0B' },
  high:     { bg:'rgba(243,232,255,0.4)', border:'#d8b4fe', badge:'#F3E8FF', badgeText:'#7C3AED', text:'#6B21A8', dot:'#A855F7' },
}

function formatTime(iso) {
  return new Date(iso).toLocaleTimeString('en', { hour:'2-digit', minute:'2-digit', hour12:false })
}

function getMetricColor(value) {
  if (value >= 80) return { stroke:'#EF4444', fill:'rgba(239,68,68,0.1)' }
  if (value >= 60) return { stroke:'#F59E0B', fill:'rgba(245,158,11,0.08)' }
  return { stroke:'#4F46E5', fill:'rgba(79,70,229,0.08)' } // Layout dagi Indigo rang
}

const CustomTooltip = ({ active, payload, label }) => {
  if (!active || !payload?.length) return null
  const val = payload[0].value
  const color = getMetricColor(val).stroke
  return (
    <div style={{
      background:'#FFFFFF', border:'1px solid #E2E8F0',
      borderRadius:'8px', padding:'8px 12px',
      fontSize:'11px', fontFamily:'monospace',
      boxShadow: '0 4px 12px rgba(0,0,0,0.05)',
    }}>
      <div style={{ color:'#94A3B8', marginBottom:'2px' }}>{label}</div>
      <div style={{ color, fontWeight:'600' }}>{val}%</div>
    </div>
  )
}

// Layout bilan bir xil panel (Card) stili
const S = {
  card: { 
    background:'#FFFFFF', 
    border:'1px solid #E2E8F0', 
    borderRadius:'12px', 
    padding:'20px',
    boxShadow: '0 1px 3px rgba(0,0,0,0.02), 0 1px 2px rgba(0,0,0,0.04)'
  },
}

const BAR_H = 160  // barlarning max pixel balandligi

function LogVolumeCard({ bars, scrollRef }) {
  const [tooltip, setTooltip] = useState(null) // { x, y, hour, nginx, app, error }

  return (
    <div style={{ ...S.card, display: 'flex', flexDirection: 'column', height:'100%' }}>
      {/* Header */}
      <div style={{ display:'flex', justifyContent:'space-between', alignItems:'center', marginBottom:'14px', flexShrink: 0 }}>
        <p style={{ fontSize:'10px', color:'#64748B', fontWeight:'600', textTransform:'uppercase', letterSpacing:'0.7px' }}>
          Log Volume — Last 24h
        </p>
        <div style={{ display:'flex', gap:'12px' }}>
          {[
            { label:'Nginx',  color:'#4F46E5' }, // Indigo ko'k
            { label:'App',    color:'#10B981' }, // Yashil
            { label:'Errors', color:'#EF4444' }, // Qizil
          ].map(l => (
            <span key={l.label} style={{ fontSize:'10px', color:l.color, fontFamily:'monospace', fontWeight:'500' }}>■ {l.label}</span>
          ))}
        </div>
      </div>

      {/* Chart area */}
      <div style={{ 
        flex: 1, 
        position: 'relative', 
        minHeight: 0,
        backgroundImage: `linear-gradient(to top, rgba(226,232,240,0.5) 1px, transparent 1px)`,
        backgroundSize: '100% 32px',
      }}>
        {bars.length === 0 ? (
          <div style={{ height:'100%', display:'flex', alignItems:'center', justifyContent:'center', color:'#94A3B8', fontSize:'12px', fontFamily:'monospace' }}>
            No log data
          </div>
        ) : (
          <>
            {/* Scroll container */}
            <div
              ref={scrollRef}
              style={{
                overflowX: 'auto',
                overflowY: 'hidden',
                height: '100%',
                scrollbarWidth: 'thin',
                scrollbarColor: '#CBD5E1 #F1F5F9',
              }}
            >
              <div style={{
                display: 'flex',
                alignItems: 'flex-end',
                gap: '8px',
                paddingBottom: '25px',  /* soat label uchun joy */
                paddingTop: '8px',
                width: 'max-content',
                height: '100%',
                boxSizing: 'border-box',
              }}>
                {bars.map((b) => (
                  <div
                    key={b.hour}
                    style={{ display: 'flex', flexDirection: 'column', alignItems: 'center', gap: '0', flexShrink: 0 }}
                    onMouseEnter={e => {
                      const rect = e.currentTarget.getBoundingClientRect()
                      const containerRect = scrollRef.current.getBoundingClientRect()
                      setTooltip({
                        x: rect.left - containerRect.left + rect.width / 2,
                        y: rect.top - containerRect.top - 8,
                        hour: b.hour,
                        nginx: b.nginx,
                        app: b.app,
                        error: b.error,
                      })
                    }}
                    onMouseLeave={() => setTooltip(null)}
                  >
                    {/* 3 ta bar yonma-yon, pastdan o'sadi */}
                    <div style={{ display: 'flex', alignItems: 'flex-end', gap: '3px' }}>
                      <div style={{ width: '9px', height: `${b.nginxH}px`, minHeight: b.nginxH > 0 ? '2px' : '0', background:'#4F46E5', borderRadius:'2px 2px 0 0', transition:'height 0.3s ease' }}/>
                      <div style={{ width: '9px', height: `${b.appH}px`,   minHeight: b.appH   > 0 ? '2px' : '0', background:'#10B981', borderRadius:'2px 2px 0 0', transition:'height 0.3s ease' }}/>
                      <div style={{ width: '9px', height: `${b.errorH}px`, minHeight: b.errorH > 0 ? '2px' : '0', background:'#EF4444', borderRadius:'2px 2px 0 0', transition:'height 0.3s ease' }}/>
                    </div>
                    {/* Soat label */}
                    <span style={{ fontSize:'10px', color:'#64748B', fontFamily:'monospace', marginTop:'6px', fontWeight:'500' }}>
                      {b.hour.toString().padStart(2,'0')}
                    </span>
                  </div>
                ))}
              </div>
            </div>

            {/* Tooltip */}
            {tooltip && (
              <div style={{
                position: 'absolute',
                left: tooltip.x,
                top: tooltip.y,
                transform: 'translate(-50%, -100%)',
                background: '#FFFFFF',
                border: '1px solid #E2E8F0',
                borderRadius: '8px',
                padding: '8px 12px',
                fontSize: '11px',
                fontFamily: 'monospace',
                pointerEvents: 'none',
                zIndex: 10,
                whiteSpace: 'nowrap',
                boxShadow: '0 4px 18px rgba(0,0,0,0.08)',
              }}>
                <div style={{ color:'#94A3B8', marginBottom:'6px', fontSize:'10px', fontWeight:'600' }}>{tooltip.hour.toString().padStart(2,'0')}:00</div>
                <div style={{ display:'flex', alignItems:'center', gap:'8px', marginBottom:'4px' }}>
                  <span style={{ width:'8px', height:'8px', background:'#4F46E5', borderRadius:'2px', display:'inline-block' }}/>
                  <span style={{ color:'#64748B' }}>Nginx</span>
                  <span style={{ color:'#4F46E5', fontWeight:'600', marginLeft:'auto', paddingLeft:'16px' }}>{tooltip.nginx.toLocaleString()}</span>
                </div>
                <div style={{ display:'flex', alignItems:'center', gap:'8px', marginBottom:'4px' }}>
                  <span style={{ width:'8px', height:'8px', background:'#10B981', borderRadius:'2px', display:'inline-block' }}/>
                  <span style={{ color:'#64748B' }}>App</span>
                  <span style={{ color:'#10B981', fontWeight:'600', marginLeft:'auto', paddingLeft:'16px' }}>{tooltip.app.toLocaleString()}</span>
                </div>
                <div style={{ display:'flex', alignItems:'center', gap:'8px' }}>
                  <span style={{ width:'8px', height:'8px', background:'#EF4444', borderRadius:'2px', display:'inline-block' }}/>
                  <span style={{ color:'#64748B' }}>Errors</span>
                  <span style={{ color:'#EF4444', fontWeight:'600', marginLeft:'auto', paddingLeft:'16px' }}>{tooltip.error.toLocaleString()}</span>
                </div>
              </div>
            )}
          </>
        )}
      </div>
    </div>
  )
}

export default function Dashboard() {
  const [alerts,      setAlerts]      = useState([])
  const [agents,      setAgents]      = useState([])
  const [agentId,     setAgentId]     = useState('')
  const [metrics,     setMetrics]     = useState({ cpu:0, ram:0, disk:0 })
  const [stats,       setStats]       = useState({
    requests:     { value:0, percent:0 },
    errorRate:    { value:0, percent:0 },
    authFail:     { value:0, percent:0 },
    activeAlerts: { total:0, critical:0 },
  })
  const [cpuHistory,  setCpuHistory]  = useState([])
  const [ramHistory,  setRamHistory]  = useState([])
  const [diskHistory, setDiskHistory] = useState([])
  const [bars, setBars] = useState([])
  const scrollRef = useRef(null)

  // Agentlar
  useEffect(() => {
    const fetchAgents = async () => {
      try {
        const res = await api.get('/agents')
        const list = (res.data?.data || []).map(a => ({ id:a.id, name:a.name }))
        setAgents(list)
        if (list.length > 0) setAgentId(list[0].id)
      } catch { /* ignore */ }
    }
    fetchAgents()
  }, [])

  const load = useCallback(async () => {
    try {
      const p   = agentId ? `agent_id=${agentId}&limit=30` : `limit=30`
      const agP = agentId ? `agent_id=${agentId}&` : ''

      const [mRes, aRes, sRes, lRes] = await Promise.allSettled([
        api.get(`/metrics?${p}`),
        api.get('/alerts?is_read=false&limit=10'),
        api.get(`/stats?${agP}period=7d`),
        api.get(`/logvolume?${agP}`),
      ])
      localStorage.getItem('sentinel_token')
      
      // Metrics
      if (mRes.status === 'fulfilled') {
        const data = mRes.value.data?.data || []
        if (Array.isArray(data) && data.length) {
          const sorted = [...data].reverse()
          const mapH = key => sorted.map((m, i) => ({
            t: i === sorted.length - 1 ? 'now' : `-${sorted.length - 1 - i}m`,
            v: Math.round(m[key]),
          }))
          setCpuHistory(mapH('cpu'))
          setRamHistory(mapH('ram'))
          setDiskHistory(mapH('disk'))
          const latest = sorted[sorted.length - 1]
          setMetrics({
            cpu:  Math.round(latest.cpu),
            ram:  Math.round(latest.ram),
            disk: Math.round(latest.disk),
          })
        }
      }

      // Alerts
      if (aRes.status === 'fulfilled') {
        const data = aRes.value.data?.data || []
        if (Array.isArray(data)) setAlerts(data)
      }

      // Stats
      if (sRes.status === 'fulfilled') {
        const d = sRes.value.data?.data || {}
        setStats({
          requests:     d.requests     || { value:0, percent:0 },
          errorRate:    d.errorRate    || { value:0, percent:0 },
          authFail:     d.authFail     || { value:0, percent:0 },
          activeAlerts: d.activeAlerts || { total:0, critical:0 },
        })
      }

      // Log volume
      if (lRes.status === 'fulfilled') {
        const raw = lRes.value.data?.data || lRes.value.data || []
        const grouped = Array.from({ length: 24 }, (_, i) => ({ hour: i, nginx: 0, app: 0, error: 0 }))
        if (Array.isArray(raw)) {
          raw.forEach(l => {
            const hour = new Date(l.time).getHours()
            if (hour >= 0 && hour < 24) {
              grouped[hour].nginx += l.nginx_logs  || 0
              grouped[hour].app   += l.app_logs    || 0
              grouped[hour].error += l.error_count || 0
            }
          })
        }
        const maxVal = Math.max(...grouped.flatMap(g => [g.nginx, g.app, g.error]), 1)
        setBars(grouped.map(g => ({
          hour:   g.hour,
          nginx:  g.nginx,
          app:    g.app,
          error:  g.error,
          nginxH: Math.round((g.nginx / maxVal) * BAR_H),
          appH:   Math.round((g.app   / maxVal) * BAR_H),
          errorH: Math.round((g.error / maxVal) * BAR_H),
        })))
      }

    } catch {
      // Mock fallback
      setMetrics(prev => {
        const cpu  = Math.max(5, Math.min(99, prev.cpu  + (Math.random()-0.5)*8))
        const ram  = Math.max(5, Math.min(99, prev.ram  + (Math.random()-0.5)*4))
        const disk = Math.max(5, Math.min(99, prev.disk + (Math.random()-0.5)*2))
        const c = Math.round(cpu), r = Math.round(ram), d = Math.round(disk)
        setCpuHistory(h  => [...h.slice(1), { t:'now', v:c }])
        setRamHistory(h  => [...h.slice(1), { t:'now', v:r }])
        setDiskHistory(h => [...h.slice(1), { t:'now', v:d }])
        return { cpu:c, ram:r, disk:d }
      })
    }
  }, [agentId])

  // Bars o'ng tomonga scroll
  useEffect(() => {
    if (scrollRef.current) {
      scrollRef.current.scrollLeft = scrollRef.current.scrollWidth
    }
  }, [bars])

  useEffect(() => {
    load()
    const t = setInterval(load, 5000)
    return () => clearInterval(t)
  }, [load])

  const statCards = [
    {
      label:'TOTAL REQUESTS',
      value: Math.round(stats.requests.value).toLocaleString(),
      sub: `${stats.requests.percent >= 0 ? '▲' : '▼'} ${Math.abs(stats.requests.percent).toFixed(1)}%`,
      subOk: stats.requests.percent >= 0,
      note:'vs last period',
      color:'#4F46E5', valColor:'#1E293B',
    },
    {
      label:'ERROR RATE',
      value: stats.errorRate.value.toFixed(1)+'%',
      sub: `${stats.errorRate.percent >= 0 ? '▲' : '▼'} ${Math.abs(stats.errorRate.percent).toFixed(1)}%`,
      subOk: stats.errorRate.percent <= 0,
      note:'vs last period',
      color:'#EF4444', valColor:'#1E293B',
    },
    {
      label:'AUTH FAILURES',
      value: Math.round(stats.authFail.value),
      sub: `${stats.authFail.percent >= 0 ? '▲' : '▼'} ${Math.abs(stats.authFail.percent).toFixed(1)}%`,
      subOk: stats.authFail.percent <= 0,
      note:'vs last period',
      color:'#F59E0B', valColor:'#1E293B',
    },
    {
      label:'ACTIVE ALERTS',
      value: stats.activeAlerts.total,
      sub: `${stats.activeAlerts.critical} critical`,
      subOk: false,
      note:`${stats.activeAlerts.total - stats.activeAlerts.critical} warning`,
      color:'#EF4444', valColor:'#DC2626',
    },
  ]

  const metricCharts = [
    { label:'CPU',  value:metrics.cpu,  history:cpuHistory,  note:'4 cores' },
    { label:'RAM',  value:metrics.ram,  history:ramHistory,  note:`${(metrics.ram * 0.09).toFixed(1)} / 9 GB` },
    { label:'Disk', value:metrics.disk, history:diskHistory, note:`${Math.round(metrics.disk * 5)} / 500 GB` },
  ]

  return (
    <div style={{ minHeight:'100vh' }}>

      {/* Top navbar */}
      <div style={{
        position:'sticky', top:0, zIndex:30,
        background:'linear-gradient(90deg, rgba(79,70,229,0.08), rgba(16,185,129,0.05))', 
        backdropFilter:'blur(12px)',
        borderBottom:'1px solid rgba(148,163,184,0.2)',
        height:'52px', display:'flex', alignItems:'center', justifyContent:'space-between',
        marginLeft:'-24px', marginRight:'-24px',
        paddingLeft:'24px', paddingRight:'24px', marginBottom:'20px',
      }}>
        <div style={{ display:'flex', alignItems:'center', gap:'8px' }}>
          <span style={{ fontSize:'15px', fontWeight:'600', color:'#1E293B' }}>Dashboard</span>
          <span style={{ fontSize:'12px', color:'#64748B', fontFamily:'monospace' }}>
            · {new Date().toLocaleDateString('en', { month:'short', day:'numeric', year:'numeric' })}
          </span>
        </div>
        <div style={{ display:'flex', alignItems:'center', gap:'12px' }}>
          <div style={{ display:'flex', alignItems:'center', gap:'6px' }}>
            <span style={{ fontSize:'11px', color:'#64748B', fontWeight:'500', textTransform:'uppercase', letterSpacing:'0.5px' }}>Agent:</span>
            <select
              value={agentId}
              onChange={e => setAgentId(e.target.value)}
              style={{
                fontSize:'12px', padding:'5px 10px',
                borderRadius:'6px', border:'1px solid #E2E8F0',
                background:'#FFFFFF', color:'#1E293B',
                fontFamily:'monospace', outline:'none', cursor:'pointer',
                boxShadow: '0 1px 2px rgba(0,0,0,0.05)',
              }}
            >
              {agents.map(a => (
                <option key={a.id} value={a.id}>{a.name || a.id}</option>
              ))}
            </select>
          </div>
          <div style={{
            display:'flex', alignItems:'center', gap:'6px',
            padding:'4px 12px', borderRadius:'20px',
            background:'#E8F5E9', border:'1px solid #A5D6A7',
          }}>
            <span style={{ width:'6px', height:'6px', borderRadius:'50%', background:'#2E7D32', display:'inline-block' }}/>
            <span style={{ fontSize:'11px', color:'#2E7D32', fontFamily:'monospace', fontWeight:'600' }}>Live</span>
          </div>
        </div>
      </div>

      {/* Stat cards */}
      <div style={{ display:'grid', gridTemplateColumns:'repeat(4,1fr)', gap:'12px', marginBottom:'16px' }}>
        {statCards.map(c => (
          <div key={c.label} style={{ ...S.card, borderTop:`3px solid ${c.color}` }}>
            <p style={{ fontSize:'10px', color:'#64748B', fontWeight:'500', letterSpacing:'0.7px', textTransform:'uppercase', marginBottom:'8px' }}>
              {c.label}
            </p>
            <p style={{ fontSize:'24px', fontWeight:'600', fontFamily:'monospace', color:'#1E293B', lineHeight:1 }}>
              {c.value}
            </p>
            <div style={{ display:'flex', justifyContent:'space-between', marginTop:'8px' }}>
              <span style={{ fontSize:'11px', fontFamily:'monospace', fontWeight:'600', color: c.subOk ? '#10B981' : '#EF4444' }}>{c.sub}</span>
              <span style={{ fontSize:'10px', color:'#94A3B8' }}>{c.note}</span>
            </div>
          </div>
        ))}
      </div>

      {/* Metrics + Log volume */}
      <div style={{ display:'grid', gridTemplateColumns:'1fr 1fr', gap:'12px', marginBottom:'16px', alignItems:'stretch' }}>

        {/* Metric charts */}
        <div style={{ ...S.card, minWidth: 0, height:'100%' }}>
          <div style={{ display:'flex', justifyContent:'space-between', alignItems:'center', marginBottom:'16px' }}>
            <p style={{ fontSize:'10px', color:'#64748B', fontWeight:'600', textTransform:'uppercase', letterSpacing:'0.7px' }}>
              Server Metrics — Live
            </p>
            {agentId && (
              <span style={{
                fontSize:'11px', padding:'2px 8px', borderRadius:'4px',
                background:'#EEF2FF', color:'#4F46E5', fontFamily:'monospace', fontWeight:'500',
              }}>
                {agents.find(a => a.id === agentId)?.name || agentId}
              </span>
            )}
          </div>
          <div style={{ display:'flex', flexDirection:'column', gap:'16px' }}>
            {metricCharts.map(m => {
              const c = getMetricColor(m.value)
              return (
                <div key={m.label}>
                  <div style={{ display:'flex', justifyContent:'space-between', alignItems:'center', marginBottom:'4px' }}>
                    <div style={{ display:'flex', alignItems:'center', gap:'8px' }}>
                      <span style={{ fontSize:'13px', color:'#475569', fontWeight:'500' }}>{m.label}</span>
                      <span style={{ fontSize:'11px', color:'#94A3B8', fontFamily:'monospace' }}>{m.note}</span>
                    </div>
                    <span style={{
                      fontSize:'13px', fontWeight:'600', fontFamily:'monospace',
                      color:c.stroke, padding:'2px 8px', borderRadius:'4px', background:c.fill,
                    }}>
                      {m.value}%
                    </span>
                  </div>
                  <ResponsiveContainer width="100%" height={52}>
                    <AreaChart data={m.history} margin={{ top:2, right:0, bottom:0, left:0 }}>
                      <defs>
                        <linearGradient id={`grad-${m.label}`} x1="0" y1="0" x2="0" y2="1">
                          <stop offset="0%"   stopColor={c.stroke} stopOpacity={0.25}/>
                          <stop offset="100%" stopColor={c.stroke} stopOpacity={0.01}/>
                        </linearGradient>
                      </defs>
                      <XAxis dataKey="t" hide />
                      <YAxis domain={[0,100]} hide />
                      <Tooltip content={<CustomTooltip />} />
                      <Area
                        type="monotone" dataKey="v"
                        stroke={c.stroke} strokeWidth={2}
                        fill={`url(#grad-${m.label})`}
                        dot={false}
                        activeDot={{ r:3, fill:c.stroke, strokeWidth:0 }}
                        isAnimationActive={false}
                      />
                    </AreaChart>
                  </ResponsiveContainer>
                  <div style={{ display:'flex', gap:'12px', marginTop:'4px' }}>
                    {[
                      { label:'0–60%',  color:'#4F46E5' },
                      { label:'60–80%', color:'#F59E0B' },
                      { label:'80%+',   color:'#EF4444' },
                    ].map(r => (
                      <span key={r.label} style={{ display:'flex', alignItems:'center', gap:'4px', fontSize:'10px', color:'#64748B' }}>
                        <span style={{ width:'6px', height:'6px', borderRadius:'50%', background:r.color, display:'inline-block' }}/>
                        {r.label}
                      </span>
                    ))}
                  </div>
                </div>
              )
            })}
          </div>
        </div>

        {/* Log volume */}
        <div style={{ minWidth: 0, height:'100%' }}>
          <LogVolumeCard bars={bars} scrollRef={scrollRef} />
        </div>
      </div>

      {/* Recent alerts */}
      <div style={S.card}>
        <div style={{ display:'flex', justifyContent:'space-between', alignItems:'center', marginBottom:'16px' }}>
          <p style={{ fontSize:'10px', color:'#64748B', fontWeight:'600', textTransform:'uppercase', letterSpacing:'0.7px' }}>Recent Alerts</p>
          <a href="/alerts" style={{ fontSize:'12px', color:'#4F46E5', fontWeight:'500', textDecoration:'none' }}>View all →</a>
        </div>
        <table style={{ width:'100%', borderCollapse:'collapse' }}>
          <thead>
            <tr style={{ borderBottom:'2px solid #E2E8F0' }}>
              {['','Severity','Type','Message','Agent','Time',''].map((h,i) => (
                <th key={i} style={{
                  textAlign:'left', fontSize:'11px', color:'#64748B',
                  fontWeight:'600', textTransform:'uppercase', letterSpacing:'0.5px',
                  padding:'0 10px 10px',
                }}>{h}</th>
              ))}
            </tr>
          </thead>
          <tbody>
            {alerts.map((a, i) => {
              const s = SEV[a.severity] || SEV.warning
              return (
                <tr
                  key={i}
                  style={{
                    borderBottom:'1px solid #E2E8F0',
                    background: !a.is_read ? s.bg : 'transparent',
                    borderLeft: !a.is_read ? `3px solid ${s.dot}` : '3px solid transparent',
                    transition: 'background 0.15s',
                  }}
                  onMouseEnter={e => e.currentTarget.style.background='#F8FAFC'}
                  onMouseLeave={e => e.currentTarget.style.background=!a.is_read ? s.bg : 'transparent'}
                >
                  <td style={{ padding:'12px 0 12px 10px' }}>
                    <span style={{ width:'6px', height:'6px', borderRadius:'50%', background:s.dot, display:'inline-block' }}/>
                  </td>
                  <td style={{ padding:'12px' }}>
                    <span style={{ fontSize:'10px', padding:'3px 8px', borderRadius:'4px', background:s.badge, color:s.badgeText, fontFamily:'monospace', fontWeight:'600' }}>
                      {a.severity?.toUpperCase()}
                    </span>
                  </td>
                  <td style={{ padding:'12px', fontSize:'12px', color:'#475569', fontFamily:'monospace', fontWeight:'500' }}>{a.type}</td>
                  <td style={{ padding:'12px', fontSize:'12px', color:s.text, fontFamily:'monospace', maxWidth:'260px', overflow:'hidden', textOverflow:'ellipsis', whiteSpace:'nowrap', fontWeight:'500' }}>
                    {a.message}
                  </td>
                  <td style={{ padding:'12px', fontSize:'12px', color:'#64748B', fontFamily:'monospace' }}>{a.agent_name}</td>
                  <td style={{ padding:'12px', fontSize:'12px', color:'#64748B', fontFamily:'monospace' }}>{formatTime(a.created_at)}</td>
                  <td style={{ padding:'12px' }}>
                    {!a.is_read && <span style={{ width:'6px', height:'6px', borderRadius:'50%', background:'#3B82F6', display:'inline-block' }}/>}
                  </td>
                </tr>
              )
            })}
          </tbody>
        </table>
        {alerts.length === 0 && (
          <div style={{ padding:'48px', textAlign:'center', color:'#94A3B8', fontFamily:'monospace', fontSize:'13px' }}>
            No alerts found
          </div>
        )}
      </div>
    </div>
  )
}