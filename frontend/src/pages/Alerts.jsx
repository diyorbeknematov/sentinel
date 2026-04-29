import { useState, useEffect, useRef } from 'react'
import api from '../api/axios'

const SEV_STYLE = {
  critical: { bg:'rgba(239,68,68,0.15)',  text:'#fca5a5', row:'rgba(239,68,68,0.05)',  dot:'#ef4444' },
  warning:  { bg:'rgba(245,158,11,0.15)', text:'#fcd34d', row:'rgba(245,158,11,0.05)', dot:'#f59e0b' },
  high:     { bg:'rgba(168,85,247,0.15)', text:'#d8b4fe', row:'rgba(168,85,247,0.05)', dot:'#a855f7' },
}

const mockAlerts = [
  { id:'1', severity:'critical', type:'auth_failure', message:'SSH brute force — 120 failed attempts', agent_id:'nginx-01', created_at:'2024-01-25T14:32:00Z', is_read:false },
  { id:'2', severity:'warning',  type:'high_memory',  message:'RAM usage exceeded 90% threshold',      agent_id:'app-01',   created_at:'2024-01-25T14:29:00Z', is_read:false },
  { id:'3', severity:'high',     type:'rate_limit',   message:'Rate limit exceeded: 1200 req/min',     agent_id:'nginx-01', created_at:'2024-01-25T14:25:00Z', is_read:false },
  { id:'4', severity:'critical', type:'auth_failure', message:'Failed login from unknown IP',          agent_id:'nginx-02', created_at:'2024-01-25T14:20:00Z', is_read:true  },
  { id:'5', severity:'warning',  type:'disk_space',   message:'Disk usage exceeded 80% threshold',     agent_id:'app-01',   created_at:'2024-01-25T14:15:00Z', is_read:true  },
  { id:'6', severity:'high',     type:'cpu_spike',    message:'CPU usage spiked to 95% for 5 minutes', agent_id:'app-02',   created_at:'2024-01-25T14:10:00Z', is_read:true  },
]

function formatTime(iso) {
  return new Date(iso).toLocaleString('en', {
    month:'short', day:'numeric',
    hour:'2-digit', minute:'2-digit', hour12:false,
  })
}

function Badge({ label, bg, text }) {
  return (
    <span style={{
      fontSize:'10px', padding:'2px 8px', borderRadius:'4px',
      fontFamily:'monospace', fontWeight:'500',
      background:bg, color:text, display:'inline-block',
    }}>
      {label}
    </span>
  )
}

function DropdownFilter({ label, options, value, onChange }) {
  const [open, setOpen] = useState(false)
  const ref = useRef()

  useEffect(() => {
    const h = (e) => { if (ref.current && !ref.current.contains(e.target)) setOpen(false) }
    document.addEventListener('mousedown', h)
    return () => document.removeEventListener('mousedown', h)
  }, [])

  const selected = options.find(o => o.value === value)

  return (
    <div ref={ref} style={{ position:'relative', display:'inline-block' }}>
      <button
        onClick={() => setOpen(!open)}
        style={{
          display:'flex', alignItems:'center', gap:'4px',
          fontSize:'10px', color: value ? '#a5b4fc' : '#64748b',
          fontWeight:'500', textTransform:'uppercase', letterSpacing:'0.5px',
          cursor:'pointer', padding:'2px 6px', borderRadius:'4px',
          border:'none', background:'none', fontFamily:'sans-serif',
        }}
        onMouseEnter={e => e.currentTarget.style.background='rgba(30,41,59,0.5)'}
        onMouseLeave={e => e.currentTarget.style.background='none'}
      >
        {label}{selected?.value ? `: ${selected.label}` : ''}
        <svg width="10" height="10" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2.5">
          <polyline points="6 9 12 15 18 9"/>
        </svg>
      </button>

      {open && (
        <div style={{
          position:'absolute', top:'calc(100% + 4px)', left:0,
          background:'#0d1120', border:'1px solid #1e293b',
          borderRadius:'8px', boxShadow:'0 8px 24px rgba(0,0,0,0.4)',
          zIndex:100, minWidth:'160px', overflow:'hidden',
        }}>
          {options.map(o => (
            <div
              key={o.value}
              onClick={() => { onChange(o.value === value ? '' : o.value); setOpen(false) }}
              style={{
                display:'flex', alignItems:'center', gap:'8px',
                padding:'8px 14px', fontSize:'12px',
                color: o.value === value ? '#a5b4fc' : '#94a3b8',
                background: o.value === value ? 'rgba(99,102,241,0.1)' : 'transparent',
                cursor:'pointer', fontFamily:'monospace',
              }}
              onMouseEnter={e => { if (o.value !== value) e.currentTarget.style.background='rgba(30,41,59,0.5)' }}
              onMouseLeave={e => { if (o.value !== value) e.currentTarget.style.background='transparent' }}
            >
              {o.dot && <span style={{ width:'7px', height:'7px', borderRadius:'50%', background:o.dot, flexShrink:0 }}/>}
              {o.label}
            </div>
          ))}
        </div>
      )}
    </div>
  )
}

const S = {
  card: { background:'#0d1120', border:'1px solid #1e293b', borderRadius:'12px' },
  th: {
    textAlign:'left', fontSize:'10px', color:'#475569',
    fontWeight:'500', textTransform:'uppercase', letterSpacing:'0.5px',
    padding:'0 12px 12px',
  },
  td: {
    padding:'10px 12px', fontSize:'12px',
    color:'#94a3b8', fontFamily:'monospace',
    whiteSpace:'nowrap', overflow:'hidden', textOverflow:'ellipsis',
    verticalAlign:'middle',
  },
}

export default function Alerts() {
  const [alerts, setAlerts]           = useState(mockAlerts)
  const [agents, setAgents]           = useState([])
  const [agentId, setAgentId]         = useState('')
  const [severity, setSeverity]       = useState('')
  const [isRead, setIsRead]           = useState('')
  const [from, setFrom]               = useState('')
  const [to, setTo]                   = useState('')
  const [loading, setLoading]         = useState(false)
  const [markingId, setMarkingId]     = useState(null)
  const [selectedAlert, setSelectedAlert] = useState(null)

  const counts = {
    critical: alerts.filter(a => a.severity === 'critical').length,
    warning:  alerts.filter(a => a.severity === 'warning').length,
    high:     alerts.filter(a => a.severity === 'high').length,
    unread:   alerts.filter(a => !a.is_read).length,
  }

  useEffect(() => {
    const fetchAgents = async () => {
      try {
        const res = await api.get('/sentinel/api/agents')
        if (Array.isArray(res.data)) setAgents(res.data)
      } catch { /* mock */ }
    }
    fetchAgents()
  }, [])

  const fetchAlerts = async () => {
    setLoading(true)
    try {
      const params = new URLSearchParams()
      if (agentId)    params.append('agent_id', agentId)
      if (severity)   params.append('severity', severity)
      if (isRead !== '') params.append('is_read', isRead)
      if (from)       params.append('from', from)
      if (to)         params.append('to', to)
      params.append('limit', 50)

      const res = await api.get(`/sentinel/api/alerts?${params}`)
      if (Array.isArray(res.data)) setAlerts(res.data)
    } catch { /* mock */ }
    finally { setLoading(false) }
  }

  useEffect(() => { fetchAlerts() }, [agentId, severity, isRead])

  const markRead = async (id) => {
    setMarkingId(id)
    try {
      await api.put(`/sentinel/api/alerts/${id}/markread`)
      setAlerts(prev => prev.map(a => a.id === id ? { ...a, is_read:true } : a))
    } catch {
      setAlerts(prev => prev.map(a => a.id === id ? { ...a, is_read:true } : a))
    }
    setMarkingId(null)
  }

  const agentOptions = [
    { value:'', label:'All agents' },
    ...agents.map(a => ({ value:a.id, label:a.name || a.id })),
  ]

  const severityOptions = [
    { value:'',         label:'All severity' },
    { value:'critical', label:'critical', dot:'#ef4444' },
    { value:'warning',  label:'warning',  dot:'#f59e0b' },
    { value:'high',     label:'high',     dot:'#a855f7' },
  ]

  const readOptions = [
    { value:'',      label:'All' },
    { value:'false', label:'Unread', dot:'#60a5fa' },
    { value:'true',  label:'Read',   dot:'#475569' },
  ]

  return (
    <div style={{ minHeight:'100vh' }}>

      {/* ── MODAL ── */}
      {selectedAlert && (
        <div
          onClick={() => setSelectedAlert(null)}
          style={{
            position:'fixed', inset:0, zIndex:200,
            background:'rgba(0,0,0,0.75)',
            backdropFilter:'blur(6px)',
            display:'flex', alignItems:'center', justifyContent:'center',
            padding:'24px',
          }}
        >
          <div
            onClick={e => e.stopPropagation()}
            style={{
              background:'#0d1120',
              border:`1px solid ${SEV_STYLE[selectedAlert.severity]?.dot || '#475569'}55`,
              borderRadius:'14px', padding:'24px',
              width:'100%', maxWidth:'500px',
              boxShadow:'0 24px 64px rgba(0,0,0,0.6)',
            }}
          >
            {/* Header */}
            <div style={{ display:'flex', alignItems:'center', justifyContent:'space-between', marginBottom:'20px' }}>
              <div style={{ display:'flex', alignItems:'center', gap:'10px' }}>
                <span style={{
                  width:'8px', height:'8px', borderRadius:'50%',
                  background: SEV_STYLE[selectedAlert.severity]?.dot,
                  display:'inline-block',
                }}/>
                <span style={{ fontSize:'14px', fontWeight:'600', color:'#e2e8f0' }}>
                  Alert Details
                </span>
                <Badge
                  label={selectedAlert.severity}
                  bg={SEV_STYLE[selectedAlert.severity]?.bg}
                  text={SEV_STYLE[selectedAlert.severity]?.text}
                />
              </div>
              <button
                onClick={() => setSelectedAlert(null)}
                style={{
                  background:'none', border:'none', cursor:'pointer',
                  color:'#64748b', padding:'4px 8px', borderRadius:'6px',
                  fontSize:'16px', lineHeight:1,
                }}
                onMouseEnter={e => e.currentTarget.style.background='rgba(255,255,255,0.06)'}
                onMouseLeave={e => e.currentTarget.style.background='none'}
              >✕</button>
            </div>

            {/* Fields */}
            <div style={{ display:'flex', flexDirection:'column', gap:'0' }}>
              {[
                { label:'ID',       value: selectedAlert.id },
                { label:'Severity', value: selectedAlert.severity, color: SEV_STYLE[selectedAlert.severity]?.text },
                { label:'Type',     value: selectedAlert.type },
                { label:'Message',  value: selectedAlert.message },
                { label:'Agent',    value: selectedAlert.agent_id, color:'#a5b4fc' },
                { label:'Created',  value: new Date(selectedAlert.created_at).toLocaleString('en') },
                { label:'Status',   value: selectedAlert.is_read ? 'Read' : 'Unread', color: selectedAlert.is_read ? '#475569' : '#60a5fa' },
              ].map((f, i, arr) => (
                <div key={f.label} style={{
                  display:'flex', gap:'16px', alignItems:'flex-start',
                  padding:'10px 0',
                  borderBottom: i < arr.length - 1 ? '1px solid #1e293b' : 'none',
                }}>
                  <span style={{
                    fontSize:'11px', color:'#475569', fontFamily:'monospace',
                    width:'68px', flexShrink:0, paddingTop:'1px',
                  }}>
                    {f.label}
                  </span>
                  <span style={{
                    fontSize:'12px', color: f.color || '#94a3b8',
                    fontFamily:'monospace', wordBreak:'break-word', lineHeight:1.6,
                  }}>
                    {f.value}
                  </span>
                </div>
              ))}
            </div>

            {/* Actions */}
            <div style={{ display:'flex', justifyContent:'flex-end', gap:'8px', marginTop:'20px' }}>
              {!selectedAlert.is_read && (
                <button
                  onClick={() => {
                    markRead(selectedAlert.id)
                    setSelectedAlert({ ...selectedAlert, is_read:true })
                  }}
                  style={{
                    fontSize:'12px', padding:'8px 18px', borderRadius:'8px',
                    border:'1px solid rgba(34,197,94,0.25)',
                    background:'rgba(34,197,94,0.08)', color:'#86efac',
                    cursor:'pointer', fontFamily:'monospace',
                  }}
                >
                  Mark as read
                </button>
              )}
              <button
                onClick={() => setSelectedAlert(null)}
                style={{
                  fontSize:'12px', padding:'8px 18px', borderRadius:'8px',
                  border:'1px solid #1e293b',
                  background:'transparent', color:'#64748b',
                  cursor:'pointer', fontFamily:'monospace',
                }}
              >
                Close
              </button>
            </div>
          </div>
        </div>
      )}

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
          <span style={{ fontSize:'14px', fontWeight:'600', color:'#e2e8f0' }}>Alerts</span>
          {counts.unread > 0 && (
            <span style={{
              fontSize:'10px', padding:'1px 8px', borderRadius:'10px',
              background:'rgba(239,68,68,0.15)', color:'#fca5a5',
              fontFamily:'monospace',
            }}>
              {counts.unread} unread
            </span>
          )}
        </div>
        {loading && (
          <span style={{ fontSize:'10px', color:'#64748b', fontFamily:'monospace' }}>loading...</span>
        )}
      </div>

      {/* Summary cards */}
      <div style={{ display:'grid', gridTemplateColumns:'repeat(3,1fr)', gap:'12px', marginBottom:'16px' }}>
        {[
          { label:'Critical', key:'critical', color:'#ef4444', border:'#ef4444' },
          { label:'Warning',  key:'warning',  color:'#f59e0b', border:'#f59e0b' },
          { label:'High',     key:'high',     color:'#a855f7', border:'#a855f7' },
        ].map(c => (
          <div key={c.label} style={{
            background:'#0d1120', border:'1px solid #1e293b',
            borderTop:`2px solid ${c.border}`,
            borderRadius:'12px', padding:'16px',
          }}>
            <p style={{ fontSize:'9px', color:'#475569', textTransform:'uppercase', letterSpacing:'0.7px', marginBottom:'8px' }}>
              {c.label}
            </p>
            <p style={{ fontSize:'28px', fontWeight:'600', fontFamily:'monospace', color:c.color, lineHeight:1 }}>
              {counts[c.key]}
            </p>
            <p style={{ fontSize:'10px', color:'#475569', marginTop:'6px' }}>
              {alerts.filter(a => a.severity === c.key && !a.is_read).length} unread
            </p>
          </div>
        ))}
      </div>

      {/* Filter bar */}
      <div style={{ ...S.card, padding:'14px 16px', marginBottom:'12px' }}>
        <div style={{ display:'flex', alignItems:'center', gap:'12px', flexWrap:'wrap' }}>

          <div style={{ display:'flex', flexDirection:'column', gap:'3px' }}>
            <span style={{ fontSize:'9px', color:'#475569', textTransform:'uppercase', letterSpacing:'0.6px' }}>Agent</span>
            <select
              value={agentId}
              onChange={e => setAgentId(e.target.value)}
              style={{
                fontSize:'12px', padding:'6px 10px',
                borderRadius:'6px', border:'1px solid #1e293b',
                background:'#080b14', color:'#e2e8f0',
                fontFamily:'monospace', outline:'none', cursor:'pointer',
              }}
            >
              {agentOptions.map(o => (
                <option key={o.value} value={o.value}>{o.label}</option>
              ))}
            </select>
          </div>

          <div style={{ display:'flex', flexDirection:'column', gap:'3px' }}>
            <span style={{ fontSize:'9px', color:'#475569', textTransform:'uppercase', letterSpacing:'0.6px' }}>From → To</span>
            <div style={{
              display:'flex', alignItems:'center',
              border:'1px solid #1e293b', borderRadius:'6px', overflow:'hidden',
            }}>
              <input
                type="date" value={from} onChange={e => setFrom(e.target.value)}
                style={{
                  fontSize:'11px', padding:'6px 10px', border:'none',
                  background:'#080b14', color: from ? '#e2e8f0' : '#475569',
                  fontFamily:'monospace', outline:'none', width:'130px',
                }}
              />
              <span style={{
                padding:'0 8px', color:'#334155', fontSize:'11px',
                background:'#0d1120',
                borderLeft:'1px solid #1e293b', borderRight:'1px solid #1e293b',
              }}>→</span>
              <input
                type="date" value={to} onChange={e => setTo(e.target.value)}
                style={{
                  fontSize:'11px', padding:'6px 10px', border:'none',
                  background:'#080b14', color: to ? '#e2e8f0' : '#475569',
                  fontFamily:'monospace', outline:'none', width:'130px',
                }}
              />
            </div>
          </div>

          <button
            onClick={fetchAlerts}
            style={{
              marginTop:'14px', fontSize:'12px', padding:'7px 18px',
              borderRadius:'6px', border:'1px solid rgba(99,102,241,0.3)',
              background:'rgba(99,102,241,0.1)', color:'#a5b4fc',
              cursor:'pointer', fontFamily:'monospace',
            }}
          >
            Apply
          </button>

          {(agentId || severity || isRead || from || to) && (
            <button
              onClick={() => { setAgentId(''); setSeverity(''); setIsRead(''); setFrom(''); setTo('') }}
              style={{
                marginTop:'14px', fontSize:'12px', padding:'7px 12px',
                borderRadius:'6px', border:'1px solid #1e293b',
                background:'transparent', color:'#64748b',
                cursor:'pointer', fontFamily:'monospace',
              }}
            >
              Clear
            </button>
          )}
        </div>
      </div>

      {/* Table */}
      <div style={S.card}>
        <table style={{ width:'100%', borderCollapse:'collapse', tableLayout:'fixed' }}>
          <colgroup>
            <col style={{ width:'14px' }}/>
            <col style={{ width:'110px' }}/>
            <col style={{ width:'100px' }}/>
            <col style={{ width:'120px' }}/>
            <col/>
            <col style={{ width:'90px' }}/>
            <col style={{ width:'90px' }}/>
            <col style={{ width:'95px' }}/>
          </colgroup>
          <thead>
            <tr style={{ borderBottom:'1px solid #1e293b', background:'rgba(15,20,40,0.5)' }}>
              <th style={S.th}></th>
              <th style={S.th}>Time</th>
              <th style={{ ...S.th }}>
                <DropdownFilter label="Severity" options={severityOptions} value={severity} onChange={setSeverity}/>
              </th>
              <th style={S.th}>Type</th>
              <th style={S.th}>Message</th>
              <th style={S.th}>Agent</th>
              <th style={{ ...S.th }}>
                <DropdownFilter label="Status" options={readOptions} value={isRead} onChange={setIsRead}/>
              </th>
              <th style={S.th}>Action</th>
            </tr>
          </thead>
          <tbody>
            {alerts.map(alert => {
              const s = SEV_STYLE[alert.severity] || SEV_STYLE.warning
              return (
                <tr
                  key={alert.id}
                  onClick={() => setSelectedAlert(alert)}
                  style={{
                    borderBottom:'1px solid rgba(30,41,59,0.5)',
                    background: !alert.is_read ? s.row : 'transparent',
                    borderLeft: !alert.is_read ? `2px solid ${s.dot}` : '2px solid transparent',
                    cursor:'pointer',
                  }}
                  onMouseEnter={e => e.currentTarget.style.background='rgba(30,41,59,0.35)'}
                  onMouseLeave={e => e.currentTarget.style.background=!alert.is_read ? s.row : 'transparent'}
                >
                  <td style={{ padding:'10px 0 10px 10px' }}>
                    <span style={{ width:'6px', height:'6px', borderRadius:'50%', background:s.dot, display:'inline-block' }}/>
                  </td>
                  <td style={{ ...S.td, color:'#64748b' }}>{formatTime(alert.created_at)}</td>
                  <td style={S.td}>
                    <Badge label={alert.severity} bg={s.bg} text={s.text}/>
                  </td>
                  <td style={{ ...S.td, color:'#cbd5e1' }}>{alert.type}</td>
                  <td style={{ ...S.td, color:s.text }}>{alert.message}</td>
                  <td style={S.td}>
                    <Badge label={alert.agent_id} bg='rgba(100,116,139,0.12)' text='#94a3b8'/>
                  </td>
                  <td style={S.td}>
                    {alert.is_read
                      ? <span style={{ fontSize:'10px', color:'#334155', fontFamily:'monospace' }}>read</span>
                      : <span style={{ fontSize:'10px', color:'#60a5fa', fontFamily:'monospace', display:'flex', alignItems:'center', gap:'4px' }}>
                          <span style={{ width:'5px', height:'5px', borderRadius:'50%', background:'#60a5fa', display:'inline-block' }}/>
                          unread
                        </span>
                    }
                  </td>
                  <td style={S.td} onClick={e => e.stopPropagation()}>
                    {!alert.is_read && (
                      <button
                        onClick={() => markRead(alert.id)}
                        disabled={markingId === alert.id}
                        style={{
                          fontSize:'10px', padding:'3px 10px', borderRadius:'4px',
                          border:'1px solid rgba(34,197,94,0.25)',
                          background:'rgba(34,197,94,0.08)',
                          color: markingId === alert.id ? '#475569' : '#86efac',
                          cursor: markingId === alert.id ? 'not-allowed' : 'pointer',
                          fontFamily:'monospace',
                        }}
                      >
                        {markingId === alert.id ? '...' : 'Mark read'}
                      </button>
                    )}
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