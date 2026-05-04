import { useState, useEffect, useRef } from 'react'
import api from '../api/axios'

const LEVEL_STYLE = {
  error: { bg:'rgba(239,68,68,0.15)',  text:'#fca5a5', row:'rgba(239,68,68,0.05)',  dot:'#ef4444' },
  warn:  { bg:'rgba(245,158,11,0.15)', text:'#fcd34d', row:'rgba(245,158,11,0.05)', dot:'#f59e0b' },
  info:  { bg:'rgba(99,102,241,0.15)', text:'#a5b4fc', row:'transparent',           dot:'#6366f1' },
  debug: { bg:'rgba(100,116,139,0.12)',text:'#94a3b8', row:'transparent',           dot:'#64748b' },
}

function formatTime(iso) {
  return new Date(iso).toLocaleTimeString('en', { hour:'2-digit', minute:'2-digit', second:'2-digit', hour12:false })
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
          zIndex:100, minWidth:'150px', overflow:'hidden',
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
              {o.dot && (
                <span style={{ width:'7px', height:'7px', borderRadius:'50%', background:o.dot, flexShrink:0 }}/>
              )}
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

export default function AppLogs() {
  const [logs, setLogs]             = useState([])
  const [agents, setAgents]         = useState([])
  const [agentId, setAgentId]       = useState('')
  const [level, setLevel]           = useState('')
  const [from, setFrom]             = useState('')
  const [to, setTo]                 = useState('')
  const [total, setTotal]           = useState(0)
  const [expandedId, setExpandedId] = useState(null)
  const [loading, setLoading]       = useState(false)

  const [page, setPage] = useState(1)
  const limit = 30

  useEffect(() => {
    const fetchAgents = async () => {
      try {
        const res = await api.get('/agents')
        const agentsData = res.data?.data || []
        if (Array.isArray(agentsData)) setAgents(agentsData)
      } catch { /* mock */ }
    }
    fetchAgents()
  }, [])

  const fetchLogs = async () => {
    setLoading(true)
    try {
      const params = new URLSearchParams()
      if (agentId) params.append('agent_id', agentId)
      if (level)   params.append('level', level)
      if (from)    params.append('from', new Date(from).toISOString())
      if (to)      params.append('to',   new Date(to).toISOString())
      params.append('limit', limit)
      params.append('page', page)

      const res = await api.get(`/applogs?${params}`)
      const appLogs = res.data?.data || []
      
      setLogs(appLogs)
      setTotal(res.data?.total || 0)
    } catch { /* mock */ }
    finally { setLoading(false) }
  }

  useEffect(() => { fetchLogs() }, [agentId, level, page])
  useEffect(() => { setPage(1) }, [agentId, level, from, to])

  const agentOptions = [
    { value:'', label:'All agents' },
    ...agents.map(a => ({ value:a.id, label:a.name || a.id })),
  ]

  const levelOptions = [
    { value:'',      label:'All levels' },
    { value:'error', label:'error', dot:'#ef4444' },
    { value:'warn',  label:'warn',  dot:'#f59e0b' },
    { value:'info',  label:'info',  dot:'#6366f1' },
    { value:'debug', label:'debug', dot:'#64748b' },
  ]

  const toggleExpand = (id) => setExpandedId(expandedId === id ? null : id)

  const totalPages = Math.ceil(total / limit)
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
          <span style={{ fontSize:'14px', fontWeight:'600', color:'#e2e8f0' }}>App Logs</span>
          <span style={{ fontSize:'11px', color:'#475569', fontFamily:'monospace' }}>
            · {total.toLocaleString()} records
          </span>
        </div>
        {loading && (
          <span style={{ fontSize:'10px', color:'#64748b', fontFamily:'monospace' }}>loading...</span>
        )}
      </div>

      {/* Filter bar */}
      <div style={{ ...S.card, padding:'14px 16px', marginBottom:'12px' }}>
        <div style={{ display:'flex', alignItems:'center', gap:'12px', flexWrap:'wrap' }}>

          {/* Agent */}
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

          {/* From → To */}
          <div style={{ display:'flex', flexDirection:'column', gap:'3px' }}>
            <span style={{ fontSize:'9px', color:'#475569', textTransform:'uppercase', letterSpacing:'0.6px' }}>From → To</span>
            <div style={{
              display:'flex', alignItems:'center',
              border:'1px solid #1e293b', borderRadius:'6px', overflow:'hidden',
            }}>
              <input
                type="datetime-local"
                value={from}
                onChange={e => setFrom(e.target.value)}
                style={{
                  fontSize:'11px', padding:'6px 10px', border:'none',
                  background:'#080b14', color: from ? '#e2e8f0' : '#475569',
                  fontFamily:'monospace', outline:'none', width:'170px',
                }}
              />
              <span style={{
                padding:'0 8px', color:'#334155', fontSize:'11px',
                background:'#0d1120',
                borderLeft:'1px solid #1e293b', borderRight:'1px solid #1e293b',
              }}>→</span>
              <input
                type="datetime-local"
                value={to}
                onChange={e => setTo(e.target.value)}
                style={{
                  fontSize:'11px', padding:'6px 10px', border:'none',
                  background:'#080b14', color: to ? '#e2e8f0' : '#475569',
                  fontFamily:'monospace', outline:'none', width:'170px',
                }}
              />
            </div>
          </div>

          <button
            onClick={fetchLogs}
            style={{
              marginTop:'14px', fontSize:'12px', padding:'7px 18px',
              borderRadius:'6px', border:'1px solid rgba(99,102,241,0.3)',
              background:'rgba(99,102,241,0.1)', color:'#a5b4fc',
              cursor:'pointer', fontFamily:'monospace',
            }}
          >
            Apply
          </button>

          {(agentId || level || from || to) && (
            <button
              onClick={() => { setAgentId(''); setLevel(''); setFrom(''); setTo('') }}
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
            <col style={{ width:'90px' }}/>
            <col style={{ width:'80px' }}/>
            <col style={{ width:'130px' }}/>
            <col/>
            <col style={{ width:'85px' }}/>
            <col style={{ width:'20px' }}/>
          </colgroup>
          <thead>
            <tr style={{ borderBottom:'1px solid #1e293b', background:'rgba(15,20,40,0.5)' }}>
              <th style={S.th}>Time</th>
              <th style={{ ...S.th }}>
                <DropdownFilter
                  label="Level"
                  options={levelOptions}
                  value={level}
                  onChange={setLevel}
                />
              </th>
              <th style={S.th}>Event</th>
              <th style={S.th}>Message</th>
              <th style={S.th}>Agent</th>
              <th style={S.th}></th>
            </tr>
          </thead>
          <tbody>
            {logs.map(log => {
              const lvl = log.level
              const ls = LEVEL_STYLE[lvl] || LEVEL_STYLE.debug
              const expanded = expandedId === log.id

              return (
                <>
                  <tr
                    key={log.id}
                    onClick={() => toggleExpand(log.id)}
                    style={{
                      borderBottom: expanded ? 'none' : '1px solid rgba(30,41,59,0.5)',
                      background: expanded ? 'rgba(30,41,59,0.2)' : ls.row,
                      cursor:'pointer',
                    }}
                    onMouseEnter={e => { if (!expanded) e.currentTarget.style.background='rgba(30,41,59,0.3)' }}
                    onMouseLeave={e => { if (!expanded) e.currentTarget.style.background=ls.row }}
                  >
                    <td style={{ ...S.td, color: ls.dot }}>
                      {formatTime(log.log_time)}
                    </td>
                    <td style={S.td}>
                      <Badge label={log.level?.toUpperCase()} bg={ls.bg} text={ls.text} />
                    </td>
                    <td style={{ ...S.td, color:'#cbd5e1' }}>{log.event}</td>
                    <td style={{ ...S.td, color: lvl === 'error' ? '#fca5a5' : lvl === 'warn' ? '#fcd34d' : '#94a3b8' }}>
                      {log.message}
                    </td>
                    <td style={S.td}>
                      <Badge label={log.agent_name} bg='rgba(100,116,139,0.12)' text='#94a3b8' />
                    </td>
                    <td style={{ ...S.td, color:'#475569', textAlign:'center' }}>
                      <span style={{ transition:'transform 0.2s', display:'inline-block', transform: expanded ? 'rotate(90deg)' : 'none' }}>›</span>
                    </td>
                  </tr>

                  {/* Expanded */}
                  {expanded && (
                    <tr key={`${log.id}-exp`} style={{ borderBottom:'1px solid rgba(30,41,59,0.5)' }}>
                      <td colSpan={6} style={{ padding:'0 16px 14px' }}>
                        <div style={{
                          background: lvl === 'error' ? '#1a0a0a' : lvl === 'warn' ? '#1a1400' : '#0a0f1a',
                          border: `1px solid ${ls.dot}33`,
                          borderRadius:'8px', padding:'12px 16px',
                          fontFamily:'monospace', fontSize:'11px', lineHeight:2,
                          color: ls.text,
                        }}>
                          <span style={{ color:'#475569' }}>log_time:  </span>{log.log_time}<br/>
                          <span style={{ color:'#475569' }}>level:     </span>
                          <span style={{ color: ls.text, fontWeight:'600' }}>{log.level}</span><br/>
                          <span style={{ color:'#475569' }}>event:     </span>
                          <span style={{ color:'#e2e8f0' }}>{log.event}</span><br/>
                          <span style={{ color:'#475569' }}>message:   </span>
                          <span style={{ color: ls.text }}>{log.message}</span><br/>
                          <span style={{ color:'#475569' }}>user_id:   </span>
                          <span style={{ color:'#a5b4fc' }}>{log.user_id}</span><br/>
                          <span style={{ color:'#475569' }}>agent_id:  </span>
                          <span style={{ color:'#818cf8' }}>{log.agent_name}</span>
                        </div>
                      </td>
                    </tr>
                  )}
                </>
              )
            })}
          </tbody>
        </table>

        {/* Pagination */}
        <div style={{
          display:'flex',
          justifyContent:'space-between',
          alignItems:'center',
          padding:'12px 16px',
          borderTop:'1px solid #1e293b'
        }}>
          <span style={{ fontSize:'11px', color:'#475569', fontFamily:'monospace' }}>
            Page {page} of {totalPages}
          </span>

          <div style={{ display:'flex', gap:'6px' }}>
            <button
              onClick={() => setPage(p => Math.max(1, p - 1))}
              disabled={page === 1}
              style={{
                padding:'6px 12px',
                fontSize:'11px',
                borderRadius:'6px',
                border:'1px solid #1e293b',
                background: page === 1 ? 'rgba(255,255,255,0.05)' : 'rgba(234,179,8,0.08)',
                color: page === 1 ? '#475569' : '#facc15',
                cursor: page === 1 ? 'not-allowed' : 'pointer',
                fontFamily:'monospace',
                transition:'all 0.2s ease',
              }}

              onMouseEnter={e => {
                if (page > 1) e.currentTarget.style.background = 'rgba(250,204,21,0.15)'
              }}

              onMouseLeave={e => {
                if (page > 1) e.currentTarget.style.background = 'rgba(234,179,8,0.08)'
              }}
            >
              Prev
            </button>

            <button
              onClick={() => setPage(p => Math.min(totalPages, p + 1))}
              disabled={page === totalPages}
              style={{
                padding:'6px 12px',
                fontSize:'11px',
                borderRadius:'6px',
                border:'1px solid #1e293b',
                background: page === totalPages ? 'rgba(255,255,255,0.05)' : 'rgba(234,179,8,0.08)',
                color: page === totalPages ? '#475569' : '#facc15',
                cursor: page === totalPages ? 'not-allowed' : 'pointer',
                fontFamily:'monospace',
                transition:'all 0.2s ease',
              }}   
              
              onMouseEnter={ e => {
                if (page != totalPages) e.currentTarget.style.background = 'rgba(250,204,21,0.15)'
              }}

              onMouseMove={ e => {
                if (page != totalPages) e.currentTarget.style.background = 'rgba(234,179,8,0.08)'
              }}
            >
              Next
            </button>
          </div>
        </div>

        {logs.length === 0 && (
          <div style={{ padding:'48px', textAlign:'center', color:'#334155', fontFamily:'monospace', fontSize:'13px' }}>
            No logs found
          </div>
        )}
      </div>
    </div>
  )
}