import { useState, useEffect, useRef } from 'react'
import api from '../api/axios'

// Och rangli dizayn uchun HTTP Metodlar stillari
const METHOD_COLORS = {
  GET:    { bg: '#EEF2FF', text: '#4F46E5' }, // Ko'k/Binafsha tonlar
  POST:   { bg: '#FEF3C7', text: '#D97706' }, // Sariq/To'q sariq
  PUT:    { bg: '#F3E8FF', text: '#7E22CE' }, // Binafsha
  DELETE: { bg: '#FEE2E2', text: '#DC2626' }, // Qizil
  PATCH:  { bg: '#ECFDF5', text: '#059669' }, // Yashil
}

function getStatusStyle(status) {
  if (status >= 500) return { bg: '#FEE2E2', text: '#991B1B' } // Server xatoligi - Qizil
  if (status >= 400) return { bg: '#FEF3C7', text: '#92400E' } // Client xatoligi - Sariq
  return                    { bg: '#E6F4EA', text: '#137333' } // Muvaffaqiyatli - Yashil
}

function getRowBg(status) {
  if (status >= 500) return '#FEF2F2'
  if (status >= 400) return '#FFFBEB'
  return 'transparent'
}

function formatTime(iso) {
  return new Date(iso).toLocaleTimeString('en', { hour: '2-digit', minute: '2-digit', second: '2-digit', hour12: false })
}

function Badge({ label, bg, text }) {
  return (
    <span style={{
      fontSize: '10px', padding: '2px 8px', borderRadius: '4px',
      fontFamily: 'monospace', fontWeight: '600',
      background: bg, color: text, display: 'inline-block',
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
    <div ref={ref} style={{ position: 'relative', display: 'inline-block' }}>
      <button
        onClick={() => setOpen(!open)}
        style={{
          display: 'flex', alignItems: 'center', gap: '4px',
          fontSize: '11px', color: value ? '#4F46E5' : '#64748B',
          fontWeight: '600', textTransform: 'uppercase', letterSpacing: '0.5px',
          cursor: 'pointer', padding: '4px 8px', borderRadius: '4px',
          border: '1px solid #E2E8F0', background: '#FFFFFF', fontFamily: 'sans-serif',
          boxShadow: '0 1px 2px rgba(0,0,0,0.02)',
        }}
        onMouseEnter={e => e.currentTarget.style.background = '#F1F5F9'}
        onMouseLeave={e => e.currentTarget.style.background = '#FFFFFF'}
      >
        {label}{selected ? `: ${selected.label}` : ''}
        <svg width="10" height="10" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2.5">
          <polyline points="6 9 12 15 18 9"/>
        </svg>
      </button>

      {open && (
        <div style={{
          position: 'absolute', top: 'calc(100% + 4px)', left: 0,
          background: '#FFFFFF', border: '1px solid #E2E8F0',
          borderRadius: '8px', boxShadow: '0 10px 15px -3px rgba(0,0,0,0.05), 0 4px 6px -4px rgba(0,0,0,0.05)',
          zIndex: 100, minWidth: '160px', overflow: 'hidden',
        }}>
          {options.map(o => (
            <div
              key={o.value}
              onClick={() => { onChange(o.value === value ? '' : o.value); setOpen(false) }}
              style={{
                display: 'flex', alignItems: 'center', gap: '8px',
                padding: '8px 14px', fontSize: '13px',
                color: o.value === value ? '#4F46E5' : '#475569',
                background: o.value === value ? '#EEF2FF' : 'transparent',
                cursor: 'pointer', fontFamily: 'monospace',
                fontWeight: o.value === value ? '600' : '500',
              }}
              onMouseEnter={e => { if (o.value !== value) e.currentTarget.style.background = '#F8FAFC' }}
              onMouseLeave={e => { if (o.value !== value) e.currentTarget.style.background = 'transparent' }}
            >
              {o.dot && (
                <span style={{ width: '7px', height: '7px', borderRadius: '50%', background: o.dot, flexShrink: 0 }}/>
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
  card: { 
    background: '#FFFFFF', 
    border: '1px solid #E2E8F0', 
    borderRadius: '12px',
    boxShadow: '0 1px 3px rgba(0,0,0,0.02), 0 1px 2px rgba(0,0,0,0.04)',
  },
  th: {
    textAlign: 'left', fontSize: '11px', color: '#64748B',
    fontWeight: '600', textTransform: 'uppercase', letterSpacing: '0.5px',
    padding: '12px',
  },
  td: {
    padding: '12px', fontSize: '13px',
    color: '#334155', fontFamily: 'monospace',
    whiteSpace: 'nowrap', overflow: 'hidden', textOverflow: 'ellipsis',
    verticalAlign: 'middle',
  },
}

export default function NginxLogs() {
  const [logs, setLogs]           = useState([])
  const [agents, setAgents]       = useState([])
  const [agentId, setAgentId]     = useState('')
  const [method, setMethod]       = useState('')
  const [status, setStatus]       = useState('')
  const [from, setFrom]           = useState('')
  const [to, setTo]               = useState('')
  const [total, setTotal]         = useState(0)
  const [expandedId, setExpandedId] = useState(null)
  const [loading, setLoading]     = useState(false)

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
      if (method)  params.append('method', method)
      if (status) {
        if (status === '2xx') params.append('status', 200)
        if (status === '4xx') params.append('status', 400)
        if (status === '5xx') params.append('status', 500)
      }
      if (from) params.append('from', new Date(from).toISOString())
      if (to)   params.append('to',   new Date(to).toISOString())
      params.append('limit', limit)
      params.append('page', page)

      const res = await api.get(`/nginxlogs?${params}`)
      const nginxLogs = res.data?.data || []

      setLogs(nginxLogs)
      setTotal(res.data?.total || 0)
    } catch { /* mock */ }
    finally { setLoading(false) }
  }

  useEffect(() => { fetchLogs() }, [agentId, method, status, page])
  useEffect(() => { setPage(1) }, [agentId, method, status, from, to])

  const agentOptions = [
    { value: '', label: 'All agents' },
    ...agents.map(a => ({ value: a.id, label: a.name || a.id })),
  ]

  const methodOptions = [
    { value: '', label: 'All methods' },
    { value: 'GET',    label: 'GET',    dot: '#4F46E5' },
    { value: 'POST',   label: 'POST',   dot: '#D97706' },
    { value: 'PUT',    label: 'PUT',    dot: '#7E22CE' },
    { value: 'DELETE', label: 'DELETE', dot: '#DC2626' },
    { value: 'PATCH',  label: 'PATCH',  dot: '#059669' },
  ]

  const statusOptions = [
    { value: '',    label: 'All status' },
    { value: '2xx', label: '2xx — Success',     dot: '#137333' },
    { value: '4xx', label: '4xx — Client err',  dot: '#D97706' },
    { value: '5xx', label: '5xx — Server err',  dot: '#DC2626' },
  ]

  const toggleExpand = (id) => setExpandedId(expandedId === id ? null : id)

  const totalPages = Math.ceil(total / limit)
  return (
    <div style={{ minHeight: '100vh', width: '100%', boxSizing: 'border-box' }}>

      {/* Top navbar */}
      <div style={{
        position: 'sticky', top: 0, zIndex: 30,
        background: 'linear-gradient(90deg, rgba(79,70,229,0.08), rgba(16,185,129,0.05))',
        backdropFilter: 'blur(12px)',
        borderBottom: '1px solid rgba(148,163,184,0.2)',
        height: '56px', display: 'flex', alignItems: 'center', justifyContent: 'space-between',
        marginLeft: '-24px', marginRight: '-24px',
        paddingLeft: '24px', paddingRight: '24px', marginBottom: '20px',
      }}>
        <div style={{ display: 'flex', alignItems: 'center', gap: '8px' }}>
          <span style={{ fontSize: '16px', fontWeight: '600', color: '#1E293B' }}>Nginx Logs</span>
          <span style={{ fontSize: '12px', color: '#64748B', fontFamily: 'monospace', fontWeight: '500' }}>
            · {total.toLocaleString()} records
          </span>
        </div>
        <div style={{ display: 'flex', alignItems: 'center', gap: '6px' }}>
          {loading && (
            <span style={{ fontSize: '12px', color: '#64748B', fontFamily: 'monospace' }}>loading...</span>
          )}
        </div>
      </div>

      {/* Filter bar */}
      <div style={{ ...S.card, padding: '16px', marginBottom: '16px' }}>
        <div style={{ display: 'flex', alignItems: 'center', gap: '16px', flexWrap: 'wrap' }}>

          {/* Agent select */}
          <div style={{ display: 'flex', flexDirection: 'column', gap: '4px' }}>
            <span style={{ fontSize: '10px', color: '#64748B', textTransform: 'uppercase', letterSpacing: '0.6px', fontWeight: '600' }}>Agent</span>
            <select
              value={agentId}
              onChange={e => setAgentId(e.target.value)}
              style={{
                fontSize: '13px', padding: '6px 12px',
                borderRadius: '6px', border: '1px solid #E2E8F0',
                background: '#FFFFFF', color: '#1E293B',
                fontFamily: 'monospace', outline: 'none', cursor: 'pointer',
                boxShadow: '0 1px 2px rgba(0,0,0,0.02)',
              }}
            >
              {agentOptions.map(o => (
                <option key={o.value} value={o.value}>{o.label}</option>
              ))}
            </select>
          </div>

          {/* From - To */}
          <div style={{ display: 'flex', flexDirection: 'column', gap: '4px' }}>
            <span style={{ fontSize: '10px', color: '#64748B', textTransform: 'uppercase', letterSpacing: '0.6px', fontWeight: '600' }}>From → To</span>
            <div style={{
              display: 'flex', alignItems: 'center',
              border: '1px solid #E2E8F0', borderRadius: '6px', overflow: 'hidden',
              boxShadow: '0 1px 2px rgba(0,0,0,0.02)',
            }}>
              <input
                type="datetime-local"
                value={from}
                onChange={e => setFrom(e.target.value)}
                style={{
                  fontSize: '12px', padding: '6px 10px', border: 'none',
                  background: '#FFFFFF', color: from ? '#1E293B' : '#94A3B8',
                  fontFamily: 'monospace', outline: 'none', width: '180px',
                }}
              />
              <span style={{
                padding: '0 8px', color: '#94A3B8', fontSize: '12px',
                background: '#F8FAFC', borderLeft: '1px solid #E2E8F0', borderRight: '1px solid #E2E8F0',
              }}>→</span>
              <input
                type="datetime-local"
                value={to}
                onChange={e => setTo(e.target.value)}
                style={{
                  fontSize: '12px', padding: '6px 10px', border: 'none',
                  background: '#FFFFFF', color: to ? '#1E293B' : '#94A3B8',
                  fontFamily: 'monospace', outline: 'none', width: '180px',
                }}
              />
            </div>
          </div>

          <button
            onClick={fetchLogs}
            style={{
              marginTop: '16px', fontSize: '12px', padding: '8px 20px',
              borderRadius: '6px', border: '1px solid #4F46E5',
              background: '#EEF2FF', color: '#4F46E5',
              cursor: 'pointer', fontFamily: 'monospace', fontWeight: '600',
              boxShadow: '0 1px 2px rgba(0,0,0,0.02)',
            }}
          >
            Apply
          </button>

          {(agentId || method || status || from || to) && (
            <button
              onClick={() => { setAgentId(''); setMethod(''); setStatus(''); setFrom(''); setTo('') }}
              style={{
                marginTop: '16px', fontSize: '12px', padding: '8px 14px',
                borderRadius: '6px', border: '1px solid #E2E8F0',
                background: '#FFFFFF', color: '#64748B',
                cursor: 'pointer', fontFamily: 'monospace', fontWeight: '500',
              }}
            >
              Clear
            </button>
          )}
        </div>
      </div>

      {/* Table */}
      <div style={S.card}>
        <div style={{ overflowX: 'auto', width: '100%' }}>
          <table style={{ width: '100%', borderCollapse: 'collapse', tableLayout: 'fixed' }}>
            <colgroup>
              <col style={{ width: '100px' }}/>
              <col style={{ width: '120px' }}/>
              <col style={{ width: '120px' }}/>
              <col style={{ minWidth: '220px' }}/>
              <col style={{ width: '140px' }}/>
              <col style={{ width: '100px' }}/>
              <col style={{ width: '120px' }}/>
              <col style={{ width: '30px' }}/>
            </colgroup>
            <thead>
              <tr style={{ borderBottom: '2px solid #E2E8F0', background: '#F8FAFC' }}>
                <th style={S.th}>Time</th>
                <th style={{ ...S.th, padding: '6px 12px' }}>
                  <DropdownFilter label="Method" options={methodOptions} value={method} onChange={setMethod}/>
                </th>
                <th style={{ ...S.th, padding: '6px 12px' }}>
                  <DropdownFilter label="Status" options={statusOptions} value={status} onChange={setStatus}/>
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
                const currentMethod = log.method?.toUpperCase()
                const currentStatus = Number(log.status)
                const mStyle = METHOD_COLORS[currentMethod] || { bg: '#F1F5F9', text: '#64748B' }
                const sStyle = getStatusStyle(currentStatus)
                const rowBg  = getRowBg(currentStatus)
                const expanded = expandedId === log.id

                return (
                  <>
                    <tr
                      key={log.id}
                      onClick={() => toggleExpand(log.id)}
                      style={{
                        borderBottom: expanded ? 'none' : '1px solid #E2E8F0',
                        background: expanded ? '#F1F5F9' : rowBg,
                        cursor: 'pointer',
                      }}
                      onMouseEnter={e => { if (!expanded) e.currentTarget.style.background = expanded ? '#F1F5F9' : (rowBg !== 'transparent' ? rowBg : '#F8FAFC') }}
                      onMouseLeave={e => { if (!expanded) e.currentTarget.style.background = expanded ? '#F1F5F9' : rowBg }}
                    >
                      <td style={{ ...S.td, color: currentStatus >= 400 ? sStyle.text : '#64748B' }}>
                        {formatTime(log.log_time)}
                      </td>
                      <td style={S.td}>
                        <Badge label={currentMethod} bg={mStyle.bg} text={mStyle.text} />
                      </td>
                      <td style={S.td}>
                        <Badge label={log.status} bg={sStyle.bg} text={sStyle.text} />
                      </td>
                      <td style={{ ...S.td, color: currentStatus >= 400 ? sStyle.text : '#1E293B', fontWeight: currentStatus >= 400 ? '600' : '400' }}>
                        {log.path}
                      </td>
                      <td style={S.td}>{log.ip_address}</td>
                      <td style={S.td}>{log.bytes.toLocaleString()}</td>
                      <td style={S.td}>
                        <Badge label={log.agent_name} bg='#F1F5F9' text='#475569' />
                      </td>
                      <td style={{ ...S.td, color: '#64748B', textAlign: 'center', fontSize: '14px' }}>
                        <span style={{ transition: 'transform 0.2s', display: 'inline-block', transform: expanded ? 'rotate(90deg)' : 'none' }}>›</span>
                      </td>
                    </tr>

                    {/* Expanded detail */}
                    {expanded && (
                      <tr key={`${log.id}-detail`} style={{ borderBottom: '1px solid #E2E8F0' }}>
                        <td colSpan={8} style={{ padding: '0 16px 14px', background: '#F1F5F9' }}>
                          <div style={{
                            background: '#FFFFFF',
                            border: `1px solid ${currentStatus >= 500 ? '#EF4444' : currentStatus >= 400 ? '#F59E0B' : '#E2E8F0'}`,
                            borderRadius: '8px', padding: '14px 18px',
                            fontFamily: 'monospace', fontSize: '12px', lineHeight: 1.8,
                            boxShadow: '0 4px 6px -1px rgba(0,0,0,0.02)',
                          }}>
                            <span style={{ color: '#64748B', fontWeight: '600' }}>ip_address: </span><span style={{ color: '#1E293B' }}>{log.ip_address}</span><br/>
                            <span style={{ color: '#64748B', fontWeight: '600' }}>method: </span>
                            <span style={{ color: mStyle.text, fontWeight: '600' }}>{currentMethod}</span><br/>
                            <span style={{ color: '#64748B', fontWeight: '600' }}>path: </span>
                            <span style={{ color: '#1E293B', wordBreak: 'break-all' }}>{log.path}</span><br/>
                            <span style={{ color: '#64748B', fontWeight: '600' }}>status: </span>
                            <span style={{ color: sStyle.text, fontWeight: '700' }}>{log.status}</span><br/>
                            <span style={{ color: '#64748B', fontWeight: '600' }}>bytes: </span><span style={{ color: '#1E293B' }}>{log.bytes}</span><br/>
                            <span style={{ color: '#64748B', fontWeight: '600' }}>user_agent: </span>
                            <span style={{ color: '#475569', wordBreak: 'break-all' }}>{log.user_agent}</span><br/>
                            <span style={{ color: '#64748B', fontWeight: '600' }}>log_time: </span><span style={{ color: '#1E293B' }}>{log.log_time}</span><br/>
                            <span style={{ color: '#64748B', fontWeight: '600' }}>agent_id: </span>
                            <span style={{ color: '#4F46E5', fontWeight: '600' }}>{log.agent_name}</span>
                          </div>
                        </td>
                      </tr>
                    )}
                  </>
                )
              })}
            </tbody>
          </table>
        </div>

        {/* Pagination */}
        <div style={{
          display: 'flex',
          justifyContent: 'space-between',
          alignItems: 'center',
          padding: '14px 16px',
          borderTop: '1px solid #E2E8F0',
          background: '#F8FAFC',
          borderBottomLeftRadius: '12px',
          borderBottomRightRadius: '12px',
        }}>
          <span style={{ fontSize: '12px', color: '#64748B', fontFamily: 'monospace', fontWeight: '500' }}>
            Page {page} of {totalPages}
          </span>

          <div style={{ display: 'flex', gap: '6px' }}>
            <button
              onClick={() => setPage(p => Math.max(1, p - 1))}
              disabled={page === 1}
              style={{
                padding: '6px 14px',
                fontSize: '12px',
                borderRadius: '6px',
                border: '1px solid #E2E8F0',
                background: page === 1 ? '#F1F5F9' : '#FFFFFF',
                color: page === 1 ? '#94A3B8' : '#1E293B',
                cursor: page === 1 ? 'not-allowed' : 'pointer',
                fontFamily: 'monospace',
                transition: 'all 0.15s ease',
                fontWeight: '500',
                boxShadow: '0 1px 2px rgba(0,0,0,0.02)',
              }}
              onMouseEnter={e => {
                if (page > 1) e.currentTarget.style.background = '#F1F5F9'
              }}
              onMouseLeave={e => {
                if (page > 1) e.currentTarget.style.background = '#FFFFFF'
              }}
            >
              Prev
            </button>

            <button
              onClick={() => setPage(p => Math.min(totalPages, p + 1))}
              disabled={page === totalPages}
              style={{
                padding: '6px 14px',
                fontSize: '12px',
                borderRadius: '6px',
                border: '1px solid #E2E8F0',
                background: page === totalPages ? '#F1F5F9' : '#FFFFFF',
                color: page === totalPages ? '#94A3B8' : '#1E293B',
                cursor: page === totalPages ? 'not-allowed' : 'pointer',
                fontFamily: 'monospace',
                transition: 'all 0.15s ease',
                fontWeight: '500',
                boxShadow: '0 1px 2px rgba(0,0,0,0.02)',
              }}   
              onMouseEnter={e => {
                if (page !== totalPages) e.currentTarget.style.background = '#F1F5F9'
              }}
              onMouseLeave={e => {
                if (page !== totalPages) e.currentTarget.style.background = '#FFFFFF'
              }}
            >
              Next
            </button>
          </div>
        </div>

        {logs.length === 0 && (
          <div style={{ padding: '48px', textAlign: 'center', color: '#94A3B8', fontFamily: 'monospace', fontSize: '13px' }}>
            No logs found
          </div>
        )}
      </div>
    </div>
  )
}