import { useState, useEffect, useRef } from 'react'
import api from '../api/axios'

const METHOD_COLORS = {
  GET:    { bg:'rgba(99,102,241,0.15)',  text:'#a5b4fc' },
  POST:   { bg:'rgba(245,158,11,0.15)', text:'#fcd34d' },
  PUT:    { bg:'rgba(168,85,247,0.15)', text:'#d8b4fe' },
  DELETE: { bg:'rgba(239,68,68,0.15)',  text:'#fca5a5' },
  PATCH:  { bg:'rgba(34,197,94,0.15)',  text:'#86efac' },
}

function getStatusStyle(status) {
  if (status >= 500) return { bg:'rgba(239,68,68,0.15)',  text:'#fca5a5' }
  if (status >= 400) return { bg:'rgba(245,158,11,0.15)', text:'#fcd34d' }
  return                    { bg:'rgba(34,197,94,0.15)',  text:'#86efac' }
}

function getRowBg(status) {
  if (status >= 500) return 'rgba(239,68,68,0.05)'
  if (status >= 400) return 'rgba(245,158,11,0.05)'
  return 'transparent'
}

function formatTime(iso) {
  return new Date(iso).toLocaleTimeString('en', { hour:'2-digit', minute:'2-digit', second:'2-digit', hour12:false })
}

function Badge({ label, bg, text }) {
  return (
    <span style={{
      fontSize:'10px', padding:'2px 8px', borderRadius:'4px',
      fontFamily:'monospace', fontWeight:'500',
      background: bg, color: text, display:'inline-block',
    }}>
      {label}
    </span>
  )
}

function DropdownFilter({ label, options, value, onChange }) {
  const [open, setOpen] = useState(false)
  const ref = useRef()

  useEffect(() => {
    const handler = (e) => { if (ref.current && !ref.current.contains(e.target)) setOpen(false) }
    document.addEventListener('mousedown', handler)
    return () => document.removeEventListener('mousedown', handler)
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
        {label}{selected ? `: ${selected.label}` : ''}
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

const mockLogs = [
  { id:'1', method:'GET',    path:'/api/users',        status:200, ip_address:'192.168.1.12', bytes:2142, user_agent:'Mozilla/5.0 (Mac)', log_time:'2024-01-25T14:32:01Z', agent_id:'nginx-01' },
  { id:'2', method:'POST',   path:'/sentinel/login',   status:401, ip_address:'203.0.113.44', bytes:412,  user_agent:'Mozilla/5.0 (Windows NT 10.0; Win64)', log_time:'2024-01-25T14:31:55Z', agent_id:'nginx-01' },
  { id:'3', method:'GET',    path:'/api/metrics',      status:500, ip_address:'10.0.0.5',     bytes:102,  user_agent:'axios/1.4.0', log_time:'2024-01-25T14:31:48Z', agent_id:'nginx-02' },
  { id:'4', method:'GET',    path:'/api/agents',       status:200, ip_address:'10.0.0.3',     bytes:5842, user_agent:'axios/1.4.0', log_time:'2024-01-25T14:31:40Z', agent_id:'nginx-01' },
  { id:'5', method:'DELETE', path:'/api/users/9982',   status:404, ip_address:'10.0.0.8',     bytes:201,  user_agent:'insomnia/2023.5', log_time:'2024-01-25T14:31:22Z', agent_id:'nginx-01' },
  { id:'6', method:'PUT',    path:'/api/users/12/role',status:200, ip_address:'192.168.1.55', bytes:344,  user_agent:'Mozilla/5.0 (Linux)', log_time:'2024-01-25T14:31:10Z', agent_id:'nginx-02' },
  { id:'7', method:'GET',    path:'/api/applogs',      status:200, ip_address:'10.0.0.3',     bytes:8921, user_agent:'axios/1.4.0', log_time:'2024-01-25T14:30:58Z', agent_id:'nginx-01' },
  { id:'8', method:'POST',   path:'/sentinel/login',   status:401, ip_address:'203.0.113.44', bytes:412,  user_agent:'Python-urllib/3.9', log_time:'2024-01-25T14:30:44Z', agent_id:'nginx-02' },
]

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

export default function NginxLogs() {
  const [logs, setLogs]           = useState(mockLogs)
  const [agents, setAgents]       = useState([])
  const [agentId, setAgentId]     = useState('')
  const [method, setMethod]       = useState('')
  const [status, setStatus]       = useState('')
  const [from, setFrom]           = useState('')
  const [to, setTo]               = useState('')
  const [total, setTotal]         = useState(mockLogs.length)
  const [expandedId, setExpandedId] = useState(null)
  const [loading, setLoading]     = useState(false)

  // Agentlar ro'yxatini olish
  useEffect(() => {
    const fetchAgents = async () => {
      try {
        const res = await api.get('/sentinel/api/agents')
        if (Array.isArray(res.data)) setAgents(res.data)
      } catch { /* mock */ }
    }
    fetchAgents()
  }, [])

  const fetchLogs = async () => {
    setLoading(true)
    try {
      const params = new URLSearchParams()
      if (agentId) params.append('agent_id', agentId)
      if (method)  params.append('method', method)
      if (status) {
        if (status === '2xx') params.append('status', 200)
        if (status === '4xx') params.append('status', 400)
        if (status === '5xx') params.append('status', 500)
      }
      if (from) params.append('from', new Date(from).toISOString())
      if (to)   params.append('to',   new Date(to).toISOString())
      params.append('limit', 50)

      const res = await api.get(`/sentinel/api/nginxlogs?${params}`)
      if (Array.isArray(res.data)) {
        setLogs(res.data)
        setTotal(res.data.length)
      }
    } catch { /* mock */ }
    finally { setLoading(false) }
  }

  // Agent yoki filter o'zgarganda darhol fetch
  useEffect(() => { fetchLogs() }, [agentId, method, status])

  const agentOptions = [
    { value:'', label:'All agents' },
    ...agents.map(a => ({ value: a.id, label: a.name || a.id })),
  ]

  const methodOptions = [
    { value:'', label:'All methods' },
    { value:'GET',    label:'GET',    dot:'#a5b4fc' },
    { value:'POST',   label:'POST',   dot:'#fcd34d' },
    { value:'PUT',    label:'PUT',    dot:'#d8b4fe' },
    { value:'DELETE', label:'DELETE', dot:'#fca5a5' },
    { value:'PATCH',  label:'PATCH',  dot:'#86efac' },
  ]

  const statusOptions = [
    { value:'',    label:'All status' },
    { value:'2xx', label:'2xx — Success',     dot:'#86efac' },
    { value:'4xx', label:'4xx — Client err',  dot:'#fcd34d' },
    { value:'5xx', label:'5xx — Server err',  dot:'#fca5a5' },
  ]

  const toggleExpand = (id) => setExpandedId(expandedId === id ? null : id)

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
          <span style={{ fontSize:'14px', fontWeight:'600', color:'#e2e8f0' }}>Nginx Logs</span>
          <span style={{ fontSize:'11px', color:'#475569', fontFamily:'monospace' }}>
            · {total.toLocaleString()} records
          </span>
        </div>
        <div style={{ display:'flex', alignItems:'center', gap:'6px' }}>
          {loading && (
            <span style={{ fontSize:'10px', color:'#64748b', fontFamily:'monospace' }}>loading...</span>
          )}
        </div>
      </div>

      {/* Filter bar */}
      <div style={{ ...S.card, padding:'14px 16px', marginBottom:'12px' }}>
        <div style={{ display:'flex', alignItems:'center', gap:'12px', flexWrap:'wrap' }}>

          {/* Agent select */}
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

          {/* From - To */}
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
                background:'#0d1120', borderLeft:'1px solid #1e293b', borderRight:'1px solid #1e293b',
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

          {(agentId || method || status || from || to) && (
            <button
              onClick={() => { setAgentId(''); setMethod(''); setStatus(''); setFrom(''); setTo('') }}
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
            <col style={{ width:'70px' }}/>
            <col/>
            <col style={{ width:'124px' }}/>
            <col style={{ width:'70px' }}/>
            <col style={{ width:'85px' }}/>
            <col style={{ width:'20px' }}/>
          </colgroup>
          <thead>
            <tr style={{ borderBottom:'1px solid #1e293b', background:'rgba(15,20,40,0.5)' }}>
              <th style={S.th}>Time</th>
              <th style={{ ...S.th }}>
                <DropdownFilter
                  label="Method"
                  options={methodOptions}
                  value={method}
                  onChange={setMethod}
                />
              </th>
              <th style={{ ...S.th }}>
                <DropdownFilter
                  label="Status"
                  options={statusOptions}
                  value={status}
                  onChange={setStatus}
                />
              </th>
              <th style={S.th}>Path</th>
              <th style={S.th}>IP Address</th>
              <th style={S.th}>Bytes</th>
              <th style={S.th}>Agent</th>
              <th style={S.th}></th>
            </tr>
          </thead>
          <tbody>
            {logs.map(log => {
              const mStyle = METHOD_COLORS[log.method] || { bg:'rgba(100,116,139,0.15)', text:'#94a3b8' }
              const sStyle = getStatusStyle(log.status)
              const rowBg  = getRowBg(log.status)
              const expanded = expandedId === log.id

              return (
                <>
                  <tr
                    key={log.id}
                    onClick={() => toggleExpand(log.id)}
                    style={{
                      borderBottom: expanded ? 'none' : '1px solid rgba(30,41,59,0.5)',
                      background: expanded ? 'rgba(30,41,59,0.2)' : rowBg,
                      cursor:'pointer',
                    }}
                    onMouseEnter={e => { if (!expanded) e.currentTarget.style.background='rgba(30,41,59,0.3)' }}
                    onMouseLeave={e => { if (!expanded) e.currentTarget.style.background=rowBg }}
                  >
                    <td style={{ ...S.td, color: log.status >= 400 ? (log.status >= 500 ? '#fca5a5' : '#fcd34d') : '#64748b' }}>
                      {formatTime(log.log_time)}
                    </td>
                    <td style={S.td}>
                      <Badge label={log.method} bg={mStyle.bg} text={mStyle.text} />
                    </td>
                    <td style={S.td}>
                      <Badge label={log.status} bg={sStyle.bg} text={sStyle.text} />
                    </td>
                    <td style={{ ...S.td, color: log.status >= 400 ? (log.status >= 500 ? '#fca5a5' : '#fcd34d') : '#e2e8f0' }}>
                      {log.path}
                    </td>
                    <td style={S.td}>{log.ip_address}</td>
                    <td style={S.td}>{log.bytes.toLocaleString()}</td>
                    <td style={S.td}>
                      <Badge label={log.agent_id} bg='rgba(100,116,139,0.12)' text='#94a3b8' />
                    </td>
                    <td style={{ ...S.td, color:'#475569', textAlign:'center' }}>
                      <span style={{ transition:'transform 0.2s', display:'inline-block', transform: expanded ? 'rotate(90deg)' : 'none' }}>›</span>
                    </td>
                  </tr>

                  {/* Expanded detail */}
                  {expanded && (
                    <tr key={`${log.id}-detail`} style={{ borderBottom:'1px solid rgba(30,41,59,0.5)' }}>
                      <td colSpan={8} style={{ padding:'0 16px 14px' }}>
                        <div style={{
                          background: log.status >= 500 ? '#1a0a0a' : log.status >= 400 ? '#1a1400' : '#0a1020',
                          border: `1px solid ${log.status >= 500 ? 'rgba(239,68,68,0.2)' : log.status >= 400 ? 'rgba(245,158,11,0.2)' : 'rgba(99,102,241,0.2)'}`,
                          borderRadius:'8px', padding:'12px 16px',
                          fontFamily:'monospace', fontSize:'11px', lineHeight:2,
                          color: log.status >= 500 ? '#fca5a5' : log.status >= 400 ? '#fcd34d' : '#94a3b8',
                        }}>
                          <span style={{ color:'#64748b' }}>ip_address: </span>{log.ip_address}<br/>
                          <span style={{ color:'#64748b' }}>method: </span>
                          <span style={{ color: METHOD_COLORS[log.method]?.text }}>{log.method}</span><br/>
                          <span style={{ color:'#64748b' }}>path: </span>
                          <span style={{ color:'#e2e8f0' }}>{log.path}</span><br/>
                          <span style={{ color:'#64748b' }}>status: </span>
                          <span style={{ color: sStyle.text, fontWeight:'600' }}>{log.status}</span><br/>
                          <span style={{ color:'#64748b' }}>bytes: </span>{log.bytes}<br/>
                          <span style={{ color:'#64748b' }}>user_agent: </span>
                          <span style={{ color:'#a5b4fc' }}>{log.user_agent}</span><br/>
                          <span style={{ color:'#64748b' }}>log_time: </span>{log.log_time}<br/>
                          <span style={{ color:'#64748b' }}>agent_id: </span>
                          <span style={{ color:'#818cf8' }}>{log.agent_id}</span>
                        </div>
                      </td>
                    </tr>
                  )}
                </>
              )
            })}
          </tbody>
        </table>

        {logs.length === 0 && (
          <div style={{ padding:'48px', textAlign:'center', color:'#334155', fontFamily:'monospace', fontSize:'13px' }}>
            No logs found
          </div>
        )}
      </div>
    </div>
  )
}