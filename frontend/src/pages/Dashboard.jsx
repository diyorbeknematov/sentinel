import { useState, useEffect, useCallback, useRef } from 'react'
import { AreaChart, Area, XAxis, YAxis, Tooltip, ResponsiveContainer } from 'recharts'
import api from '../api/axios'

const SEV = {
  critical: { bg:'rgba(239,68,68,0.08)',  border:'#ef4444', badge:'rgba(239,68,68,0.15)',  badgeText:'#fca5a5', text:'#fca5a5', dot:'#ef4444' },
  warning:  { bg:'rgba(245,158,11,0.08)', border:'#f59e0b', badge:'rgba(245,158,11,0.15)', badgeText:'#fcd34d', text:'#fcd34d', dot:'#f59e0b' },
  high:     { bg:'rgba(168,85,247,0.08)', border:'#a855f7', badge:'rgba(168,85,247,0.15)', badgeText:'#d8b4fe', text:'#d8b4fe', dot:'#a855f7' },
}

function formatTime(iso) {
  return new Date(iso).toLocaleTimeString('en', { hour:'2-digit', minute:'2-digit', hour12:false })
}

function getMetricColor(value) {
  if (value >= 80) return { stroke:'#ef4444', fill:'rgba(239,68,68,0.15)' }
  if (value >= 60) return { stroke:'#f59e0b', fill:'rgba(245,158,11,0.12)' }
  return { stroke:'#6366f1', fill:'rgba(99,102,241,0.12)' }
}

const CustomTooltip = ({ active, payload, label }) => {
  if (!active || !payload?.length) return null
  const val = payload[0].value
  const color = getMetricColor(val).stroke
  return (
    <div style={{
      background:'#0d1120', border:'1px solid #1e293b',
      borderRadius:'6px', padding:'6px 10px',
      fontSize:'11px', fontFamily:'monospace',
    }}>
      <div style={{ color:'#475569', marginBottom:'2px' }}>{label}</div>
      <div style={{ color, fontWeight:'600' }}>{val}%</div>
    </div>
  )
}

const S = {
  card: { background:'#0d1120', border:'1px solid #1e293b', borderRadius:'12px', padding:'16px' },
}

const BAR_H = 160  // barlarning max pixel balandligi

function LogVolumeCard({ bars, scrollRef }) {
  const [tooltip, setTooltip] = useState(null) // { x, y, hour, nginx, app, error }

  return (
    <div style={{ ...S.card, display: 'flex', flexDirection: 'column', height:'100%' }}>
      {/* Header */}
      <div style={{ display:'flex', justifyContent:'space-between', alignItems:'center', marginBottom:'14px', flexShrink: 0 }}>
        <p style={{ fontSize:'9px', color:'#475569', textTransform:'uppercase', letterSpacing:'0.7px' }}>
          Log Volume — Last 24h
        </p>
        <div style={{ display:'flex', gap:'10px' }}>
          {[
            { label:'Nginx',  color:'#818cf8' },
            { label:'App',    color:'#34d399' },
            { label:'Errors', color:'#f87171' },
          ].map(l => (
            <span key={l.label} style={{ fontSize:'9px', color:l.color, fontFamily:'monospace' }}>■ {l.label}</span>
          ))}
        </div>
      </div>

      {/* Chart area — flex: 1 → header qolganini to'ldiradi */}
      <div style={{ 
        flex: 1, 
        position: 'relative', 
        minHeight: 0,
        backgroundImage: ` linear-gradient(to top, rgba(148,163,184,0.08) 1px, transparent 1px)`,
        backgroundSize: '100% 32px',
      }}>
        {bars.length === 0 ? (
          <div style={{ height:'100%', display:'flex', alignItems:'center', justifyContent:'center', color:'#334155', fontSize:'12px', fontFamily:'monospace' }}>
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
                scrollbarColor: '#1e293b #080b14',
              }}
            >
              {/* width: max-content → ichidan kengayadi, tashqi card o'zgarmaydi */}
              <div style={{
                display: 'flex',
                alignItems: 'flex-end',
                gap: '8px',
                paddingBottom: '20px',  /* soat label uchun joy */
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
                    <div style={{ display: 'flex', alignItems: 'flex-end', gap: '2px' }}>
                      <div style={{ width: '9px', height: `${b.nginxH}px`, minHeight: b.nginxH > 0 ? '2px' : '0', background:'#818cf8', borderRadius:'2px 2px 0 0', transition:'height 0.3s ease' }}/>
                      <div style={{ width: '9px', height: `${b.appH}px`,   minHeight: b.appH   > 0 ? '2px' : '0', background:'#34d399', borderRadius:'2px 2px 0 0', transition:'height 0.3s ease' }}/>
                      <div style={{ width: '9px', height: `${b.errorH}px`, minHeight: b.errorH > 0 ? '2px' : '0', background:'#f87171', borderRadius:'2px 2px 0 0', transition:'height 0.3s ease' }}/>
                    </div>
                    {/* Soat label */}
                    <span style={{ fontSize:'9px', color:'#334155', fontFamily:'monospace', marginTop:'4px' }}>
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
                background: '#0d1120',
                border: '1px solid #1e293b',
                borderRadius: '6px',
                padding: '7px 10px',
                fontSize: '11px',
                fontFamily: 'monospace',
                pointerEvents: 'none',
                zIndex: 10,
                whiteSpace: 'nowrap',
              }}>
                <div style={{ color:'#64748b', marginBottom:'5px', fontSize:'10px' }}>{tooltip.hour.toString().padStart(2,'0')}:00</div>
                <div style={{ display:'flex', alignItems:'center', gap:'6px', marginBottom:'3px' }}>
                  <span style={{ width:'8px', height:'8px', background:'#818cf8', borderRadius:'2px', display:'inline-block' }}/>
                  <span style={{ color:'#94a3b8' }}>Nginx</span>
                  <span style={{ color:'#818cf8', fontWeight:'600', marginLeft:'auto', paddingLeft:'12px' }}>{tooltip.nginx.toLocaleString()}</span>
                </div>
                <div style={{ display:'flex', alignItems:'center', gap:'6px', marginBottom:'3px' }}>
                  <span style={{ width:'8px', height:'8px', background:'#34d399', borderRadius:'2px', display:'inline-block' }}/>
                  <span style={{ color:'#94a3b8' }}>App</span>
                  <span style={{ color:'#34d399', fontWeight:'600', marginLeft:'auto', paddingLeft:'12px' }}>{tooltip.app.toLocaleString()}</span>
                </div>
                <div style={{ display:'flex', alignItems:'center', gap:'6px' }}>
                  <span style={{ width:'8px', height:'8px', background:'#f87171', borderRadius:'2px', display:'inline-block' }}/>
                  <span style={{ color:'#94a3b8' }}>Errors</span>
                  <span style={{ color:'#f87171', fontWeight:'600', marginLeft:'auto', paddingLeft:'12px' }}>{tooltip.error.toLocaleString()}</span>
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

  // Agentlar — bir marta
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

  // Bars yangilanganda o'ng tomonga scroll
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
      color:'#6366f1', valColor:'#a5b4fc',
    },
    {
      label:'ERROR RATE',
      value: stats.errorRate.value.toFixed(1)+'%',
      sub: `${stats.errorRate.percent >= 0 ? '▲' : '▼'} ${Math.abs(stats.errorRate.percent).toFixed(1)}%`,
      subOk: stats.errorRate.percent <= 0,
      note:'vs last period',
      color:'#ef4444', valColor:'#fca5a5',
    },
    {
      label:'AUTH FAILURES',
      value: Math.round(stats.authFail.value),
      sub: `${stats.authFail.percent >= 0 ? '▲' : '▼'} ${Math.abs(stats.authFail.percent).toFixed(1)}%`,
      subOk: stats.authFail.percent <= 0,
      note:'vs last period',
      color:'#f59e0b', valColor:'#fcd34d',
    },
    {
      label:'ACTIVE ALERTS',
      value: stats.activeAlerts.total,
      sub: `${stats.activeAlerts.critical} critical`,
      subOk: false,
      note:`${stats.activeAlerts.total - stats.activeAlerts.critical} warning`,
      color:'#ef4444', valColor:'#fca5a5',
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
        background:'rgba(8,11,20,0.85)', backdropFilter:'blur(12px)',
        borderBottom:'1px solid #1e293b',
        height:'52px', display:'flex', alignItems:'center', justifyContent:'space-between',
        marginLeft:'-24px', marginRight:'-24px',
        paddingLeft:'24px', paddingRight:'24px', marginBottom:'20px',
      }}>
        <div style={{ display:'flex', alignItems:'center', gap:'8px' }}>
          <span style={{ fontSize:'14px', fontWeight:'600', color:'#e2e8f0' }}>Dashboard</span>
          <span style={{ fontSize:'11px', color:'#475569', fontFamily:'monospace' }}>
            · {new Date().toLocaleDateString('en', { month:'short', day:'numeric', year:'numeric' })}
          </span>
        </div>
        <div style={{ display:'flex', alignItems:'center', gap:'10px' }}>
          <div style={{ display:'flex', alignItems:'center', gap:'6px' }}>
            <span style={{ fontSize:'10px', color:'#475569', textTransform:'uppercase', letterSpacing:'0.5px' }}>Agent:</span>
            <select
              value={agentId}
              onChange={e => setAgentId(e.target.value)}
              style={{
                fontSize:'11px', padding:'5px 10px',
                borderRadius:'6px', border:'1px solid #1e293b',
                background:'#080b14', color:'#e2e8f0',
                fontFamily:'monospace', outline:'none', cursor:'pointer',
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
            background:'rgba(34,197,94,0.08)', border:'1px solid rgba(34,197,94,0.2)',
          }}>
            <span style={{ width:'6px', height:'6px', borderRadius:'50%', background:'#22c55e', display:'inline-block' }}/>
            <span style={{ fontSize:'10px', color:'#4ade80', fontFamily:'monospace' }}>Live</span>
          </div>
        </div>
      </div>

      {/* Stat cards */}
      <div style={{ display:'grid', gridTemplateColumns:'repeat(4,1fr)', gap:'12px', marginBottom:'16px' }}>
        {statCards.map(c => (
          <div key={c.label} style={{ ...S.card, borderTop:`2px solid ${c.color}` }}>
            <p style={{ fontSize:'9px', color:'#475569', letterSpacing:'0.7px', textTransform:'uppercase', marginBottom:'8px' }}>
              {c.label}
            </p>
            <p style={{ fontSize:'22px', fontWeight:'600', fontFamily:'monospace', color:c.valColor, lineHeight:1 }}>
              {c.value}
            </p>
            <div style={{ display:'flex', justifyContent:'space-between', marginTop:'8px' }}>
              <span style={{ fontSize:'10px', fontFamily:'monospace', color: c.subOk ? '#4ade80' : '#f87171' }}>{c.sub}</span>
              <span style={{ fontSize:'9px', color:'#475569' }}>{c.note}</span>
            </div>
          </div>
        ))}
      </div>

      {/* Metrics + Log volume — teng ikkiga */}
      <div style={{ display:'grid', gridTemplateColumns:'1fr 1fr', gap:'12px', marginBottom:'16px', alignItems:'stretch', gridAutoRows:'1fx'}}>

        {/* Metric charts */}
        <div style={{ ...S.card, minWidth: 0, height:'100%' }}>
          <div style={{ display:'flex', justifyContent:'space-between', alignItems:'center', marginBottom:'14px' }}>
            <p style={{ fontSize:'9px', color:'#475569', textTransform:'uppercase', letterSpacing:'0.7px' }}>
              Server Metrics — Live
            </p>
            {agentId && (
              <span style={{
                fontSize:'10px', padding:'2px 8px', borderRadius:'4px',
                background:'rgba(99,102,241,0.12)', color:'#a5b4fc', fontFamily:'monospace',
              }}>
                {agents.find(a => a.id === agentId)?.name || agentId}
              </span>
            )}
          </div>
          <div style={{ display:'flex', flexDirection:'column', gap:'14px' }}>
            {metricCharts.map(m => {
              const c = getMetricColor(m.value)
              return (
                <div key={m.label}>
                  <div style={{ display:'flex', justifyContent:'space-between', alignItems:'center', marginBottom:'4px' }}>
                    <div style={{ display:'flex', alignItems:'center', gap:'8px' }}>
                      <span style={{ fontSize:'12px', color:'#94a3b8' }}>{m.label}</span>
                      <span style={{ fontSize:'10px', color:'#475569', fontFamily:'monospace' }}>{m.note}</span>
                    </div>
                    <span style={{
                      fontSize:'13px', fontWeight:'600', fontFamily:'monospace',
                      color:c.stroke, padding:'1px 8px', borderRadius:'4px', background:c.fill,
                    }}>
                      {m.value}%
                    </span>
                  </div>
                  <ResponsiveContainer width="100%" height={52}>
                    <AreaChart data={m.history} margin={{ top:2, right:0, bottom:0, left:0 }}>
                      <defs>
                        <linearGradient id={`grad-${m.label}`} x1="0" y1="0" x2="0" y2="1">
                          <stop offset="0%"   stopColor={c.stroke} stopOpacity={0.35}/>
                          <stop offset="100%" stopColor={c.stroke} stopOpacity={0.02}/>
                        </linearGradient>
                      </defs>
                      <XAxis dataKey="t" hide />
                      <YAxis domain={[0,100]} hide />
                      <Tooltip content={<CustomTooltip />} />
                      <Area
                        type="monotone" dataKey="v"
                        stroke={c.stroke} strokeWidth={1.5}
                        fill={`url(#grad-${m.label})`}
                        dot={false}
                        activeDot={{ r:3, fill:c.stroke, strokeWidth:0 }}
                        isAnimationActive={false}
                      />
                    </AreaChart>
                  </ResponsiveContainer>
                  <div style={{ display:'flex', gap:'10px', marginTop:'3px' }}>
                    {[
                      { label:'0–60%',  color:'#6366f1' },
                      { label:'60–80%', color:'#f59e0b' },
                      { label:'80%+',   color:'#ef4444' },
                    ].map(r => (
                      <span key={r.label} style={{ display:'flex', alignItems:'center', gap:'3px', fontSize:'9px', color:'#334155' }}>
                        <span style={{ width:'6px', height:'2px', borderRadius:'1px', background:r.color, display:'inline-block' }}/>
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
        <div style={{ display:'flex', justifyContent:'space-between', alignItems:'center', marginBottom:'14px' }}>
          <p style={{ fontSize:'9px', color:'#475569', textTransform:'uppercase', letterSpacing:'0.7px' }}>Recent Alerts</p>
          <a href="/alerts" style={{ fontSize:'11px', color:'#818cf8', textDecoration:'none' }}>View all →</a>
        </div>
        <table style={{ width:'100%', borderCollapse:'collapse' }}>
          <thead>
            <tr style={{ borderBottom:'1px solid #1e293b' }}>
              {['','Severity','Type','Message','Agent','Time',''].map((h,i) => (
                <th key={i} style={{
                  textAlign:'left', fontSize:'10px', color:'#475569',
                  fontWeight:'500', textTransform:'uppercase', letterSpacing:'0.5px',
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
                    borderBottom:'1px solid rgba(30,41,59,0.5)',
                    background: !a.is_read ? s.bg : 'transparent',
                    borderLeft: !a.is_read ? `2px solid ${s.dot}` : '2px solid transparent',
                  }}
                  onMouseEnter={e => e.currentTarget.style.background='rgba(30,41,59,0.3)'}
                  onMouseLeave={e => e.currentTarget.style.background=!a.is_read ? s.bg : 'transparent'}
                >
                  <td style={{ padding:'10px 0 10px 10px' }}>
                    <span style={{ width:'6px', height:'6px', borderRadius:'50%', background:s.dot, display:'inline-block' }}/>
                  </td>
                  <td style={{ padding:'10px' }}>
                    <span style={{ fontSize:'10px', padding:'2px 8px', borderRadius:'4px', background:s.badge, color:s.badgeText, fontFamily:'monospace', fontWeight:'500' }}>
                      {a.severity?.toUpperCase()}
                    </span>
                  </td>
                  <td style={{ padding:'10px', fontSize:'11px', color:'#94a3b8', fontFamily:'monospace' }}>{a.type}</td>
                  <td style={{ padding:'10px', fontSize:'11px', color:s.text, fontFamily:'monospace', maxWidth:'260px', overflow:'hidden', textOverflow:'ellipsis', whiteSpace:'nowrap' }}>
                    {a.message}
                  </td>
                  <td style={{ padding:'10px', fontSize:'11px', color:'#64748b', fontFamily:'monospace' }}>{a.agent_name}</td>
                  <td style={{ padding:'10px', fontSize:'11px', color:'#475569', fontFamily:'monospace' }}>{formatTime(a.created_at)}</td>
                  <td style={{ padding:'10px' }}>
                    {!a.is_read && <span style={{ width:'6px', height:'6px', borderRadius:'50%', background:'#60a5fa', display:'inline-block' }}/>}
                  </td>
                </tr>
              )
            })}
          </tbody>
        </table>
        {alerts.length === 0 && (
          <div style={{ padding:'48px', textAlign:'center', color:'#334155', fontFamily:'monospace', fontSize:'13px' }}>
            No alerts found
          </div>
        )}
      </div>
    </div>
  )
}