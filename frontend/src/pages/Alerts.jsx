import { useState, useEffect, useRef } from 'react'
import api from '../api/axios'

// Och rangli dizayn uchun Severity (Xavflilik) stillari
const SEV_STYLE = {
  critical: { bg: '#FEE2E2', text: '#991B1B', row: '#FEF2F2', dot: '#EF4444' }, // Qizil tonlar
  warning:  { bg: '#FEF3C7', text: '#92400E', row: '#FFFBEB', dot: '#F59E0B' }, // Sariq/To'q sariq tonlar
  high:     { bg: '#F3E8FF', text: '#6B21A8', row: '#F9F5FF', dot: '#A855F7' }, // Binafsha tonlar
}

function formatTime(iso) {
  return new Date(iso).toLocaleString('en', {
    month: 'short', day: 'numeric',
    hour: '2-digit', minute: '2-digit', hour12: false,
  })
}

function Badge({ label, bg, text }) {
  return (
    <span style={{
      fontSize: '11px', padding: '2px 8px', borderRadius: '4px',
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
    const h = (e) => { if (ref.current && !ref.current.contains(e.target)) setOpen(false) }
    document.addEventListener('mousedown', h)
    return () => document.removeEventListener('mousedown', h)
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
        {label}{selected?.value ? `: ${selected.label}` : ''}
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
              {o.dot && <span style={{ width: '7px', height: '7px', borderRadius: '50%', background: o.dot, flexShrink: 0 }}/>}
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
    boxSizing: 'border-box',
    width: '100%',
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

export default function Alerts() {
  const [alerts, setAlerts]           = useState([])
  const [agents, setAgents]           = useState([])
  const [agentId, setAgentId]         = useState('')
  const [severity, setSeverity]       = useState('')
  const [isRead, setIsRead]           = useState('')
  const [from, setFrom]               = useState('')
  const [to, setTo]                   = useState('')
  const [loading, setLoading]         = useState(false)
  const [markingId, setMarkingId]     = useState(null)
  const [selectedAlert, setSelectedAlert] = useState(null)
  const [total, setTotal]             = useState(0)

  const [page, setPage] = useState(1)
  const limit = 30

  const counts = {
    critical: alerts.filter(a => a.severity === 'critical').length,
    warning:  alerts.filter(a => a.severity === 'warning').length,
    high:     alerts.filter(a => a.severity === 'high').length,
    unread:   alerts.filter(a => !a.is_read).length,
  }

  useEffect(() => {
    const fetchAgents = async () => {
      try {
        const res = await api.get('/agents')
        const agentsData = res.data?.data || []
        setAgents(agentsData)
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
      params.append('limit', limit)
      params.append('page', page)

      const res = await api.get(`/alerts?${params}`)
      const alertData = res.data?.data || []

      setAlerts(alertData)
      setTotal(res.data?.total || 0)
    } catch { /* mock */ }
    finally { setLoading(false) }
  }

  useEffect(() => { fetchAlerts() }, [agentId, severity, isRead, page])
  useEffect(() => { setPage(1) }, [agentId, severity, isRead, from, to])

  const markRead = async (id) => {
    setMarkingId(id)
    try {
      await api.put(`/alerts/${id}/markread`)
      setAlerts(prev => prev.map(a => a.id === id ? { ...a, is_read: true } : a))
    } catch {
      setAlerts(prev => prev.map(a => a.id === id ? { ...a, is_read: true } : a))
    }
    setMarkingId(null)
  }

  const agentOptions = [
    { value: '', label: 'All agents' },
    ...agents.map(a => ({ value: a.id, label: a.name || a.id })),
  ]

  const severityOptions = [
    { value: '',         label: 'All severity' },
    { value: 'critical', label: 'critical', dot: '#EF4444' },
    { value: 'warning',  label: 'warning',  dot: '#F59E0B' },
    { value: 'high',     label: 'high',     dot: '#A855F7' },
  ]

  const readOptions = [
    { value: '',      label: 'All' },
    { value: 'false', label: 'Unread', dot: '#3B82F6' },
    { value: 'true',  label: 'Read',   dot: '#64748B' },
  ]

  const totalPages = Math.ceil(total / limit)
  return (
    <div style={{ minHeight: '100vh', width: '100%', boxSizing: 'border-box' }}>

      {/* ── MODAL ── */}
      {selectedAlert && (
        <div
          onClick={() => setSelectedAlert(null)}
          style={{
            position: 'fixed', inset: 0, zIndex: 200,
            background: 'rgba(15, 23, 42, 0.3)',
            backdropFilter: 'blur(4px)',
            display: 'flex', alignItems: 'center', justifyContent: 'center',
            padding: '24px',
          }}
        >
          <div
            onClick={e => e.stopPropagation()}
            style={{
              background: '#FFFFFF',
              border: `1px solid ${SEV_STYLE[selectedAlert.severity]?.dot || '#E2E8F0'}`,
              borderRadius: '14px', padding: '24px',
              width: '100%', maxWidth: '500px',
              boxShadow: '0 20px 25px -5px rgba(0,0,0,0.1), 0 10px 10px -5px rgba(0,0,0,0.04)',
            }}
          >
            {/* Header */}
            <div style={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between', marginBottom: '20px' }}>
              <div style={{ display: 'flex', alignItems: 'center', gap: '10px' }}>
                <span style={{
                  width: '8px', height: '8px', borderRadius: '50%',
                  background: SEV_STYLE[selectedAlert.severity]?.dot,
                  display: 'inline-block',
                }}/>
                <span style={{ fontSize: '16px', fontWeight: '600', color: '#1E293B' }}>
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
                  background: 'none', border: 'none', cursor: 'pointer',
                  color: '#64748B', padding: '4px 8px', borderRadius: '6px',
                  fontSize: '16px', lineHeight: 1, fontWeight: '600',
                }}
                onMouseEnter={e => e.currentTarget.style.background = '#F1F5F9'}
                onMouseLeave={e => e.currentTarget.style.background = 'none'}
              >✕</button>
            </div>

            {/* Fields */}
            <div style={{ display: 'flex', flexDirection: 'column', gap: '0' }}>
              {[
                { label: 'ID',       value: selectedAlert.id },
                { label: 'Severity', value: selectedAlert.severity, color: SEV_STYLE[selectedAlert.severity]?.text },
                { label: 'Type',     value: selectedAlert.type },
                { label: 'Message',  value: selectedAlert.message, color: '#1E293B' },
                { label: 'Agent',    value: selectedAlert.agent_name, color: '#4F46E5' },
                { label: 'Created',  value: new Date(selectedAlert.created_at).toLocaleString('en') },
                { label: 'Status',   value: selectedAlert.is_read ? 'Read' : 'Unread', color: selectedAlert.is_read ? '#64748B' : '#3B82F6' },
              ].map((f, i, arr) => (
                <div key={f.label} style={{
                  display: 'flex', gap: '16px', alignItems: 'flex-start',
                  padding: '12px 0',
                  borderBottom: i < arr.length - 1 ? '1px solid #E2E8F0' : 'none',
                }}>
                  <span style={{
                    fontSize: '12px', color: '#64748B', fontFamily: 'monospace',
                    width: '75px', flexShrink: 0, paddingTop: '1px', fontWeight: '600',
                  }}>
                    {f.label}
                  </span>
                  <span style={{
                    fontSize: '13px', color: f.color || '#334155',
                    fontFamily: 'monospace', wordBreak: 'break-word', lineHeight: 1.5,
                  }}>
                    {f.value}
                  </span>
                </div>
              ))}
            </div>

            {/* Actions */}
            <div style={{ display: 'flex', justifyContent: 'flex-end', gap: '8px', marginTop: '24px' }}>
              {!selectedAlert.is_read && (
                <button
                  onClick={() => {
                    markRead(selectedAlert.id)
                    setSelectedAlert({ ...selectedAlert, is_read: true })
                  }}
                  style={{
                    fontSize: '12px', padding: '8px 18px', borderRadius: '8px',
                    border: '1px solid #10B981',
                    background: '#E8F5E9', color: '#15803D',
                    cursor: 'pointer', fontFamily: 'monospace', fontWeight: '600',
                  }}
                >
                  Mark as read
                </button>
              )}
              <button
                onClick={() => setSelectedAlert(null)}
                style={{
                  fontSize: '12px', padding: '8px 18px', borderRadius: '8px',
                  border: '1px solid #E2E8F0',
                  background: '#FFFFFF', color: '#64748B',
                  cursor: 'pointer', fontFamily: 'monospace', fontWeight: '500',
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
        position: 'sticky', top: 0, zIndex: 30,
        background: 'linear-gradient(90deg, rgba(79,70,229,0.08), rgba(16,185,129,0.05))',
        backdropFilter: 'blur(12px)',
        borderBottom: '1px solid rgba(148,163,184,0.2)',
        height: '56px', display: 'flex', alignItems: 'center', justifyContent: 'space-between',
        marginLeft: '-24px', marginRight: '-24px',
        paddingLeft: '24px', paddingRight: '24px', marginBottom: '20px',
      }}>
        <div style={{ display: 'flex', alignItems: 'center', gap: '8px' }}>
          <span style={{ fontSize: '16px', fontWeight: '600', color: '#1E293B' }}>Alerts</span>
          {counts.unread > 0 && (
            <span style={{
              fontSize: '11px', padding: '2px 8px', borderRadius: '10px',
              background: '#FEE2E2', color: '#991B1B',
              fontFamily: 'monospace', fontWeight: '600',
            }}>
              {counts.unread} unread
            </span>
          )}
        </div>
        {loading && (
          <span style={{ fontSize: '12px', color: '#64748B', fontFamily: 'monospace' }}>loading...</span>
        )}
      </div>

      {/* Summary cards */}
      <div style={{ display: 'grid', gridTemplateColumns: 'repeat(3,1fr)', gap: '16px', marginBottom: '20px', width: '100%' }}>
        {[
          { label: 'Critical', key: 'critical', color: '#B91C1C', border: '#EF4444' },
          { label: 'Warning',  key: 'warning',  color: '#C2410C', border: '#F59E0B' },
          { label: 'High',     key: 'high',     color: '#7E22CE', border: '#A855F7' },
        ].map(c => (
          <div key={c.label} style={{
            ...S.card,
            borderTop: `4px solid ${c.border}`,
            padding: '16px',
          }}>
            <p style={{ fontSize: '11px', color: '#64748B', textTransform: 'uppercase', letterSpacing: '0.7px', marginBottom: '8px', fontWeight: '500' }}>
              {c.label}
            </p>
            <p style={{ fontSize: '28px', fontWeight: '600', fontFamily: 'monospace', color: c.color, lineHeight: 1 }}>
              {counts[c.key]}
            </p>
            <p style={{ fontSize: '11px', color: '#94A3B8', marginTop: '6px', fontFamily: 'monospace' }}>
              {alerts.filter(a => a.severity === c.key && !a.is_read).length} unread
            </p>
          </div>
        ))}
      </div>

      {/* Filter bar */}
      <div style={{ ...S.card, padding: '16px', marginBottom: '16px' }}>
        <div style={{ display: 'flex', alignItems: 'center', gap: '16px', flexWrap: 'wrap' }}>

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

          <div style={{ display: 'flex', flexDirection: 'column', gap: '4px' }}>
            <span style={{ fontSize: '10px', color: '#64748B', textTransform: 'uppercase', letterSpacing: '0.6px', fontWeight: '600' }}>From → To</span>
            <div style={{
              display: 'flex', alignItems: 'center',
              border: '1px solid #E2E8F0', borderRadius: '6px', overflow: 'hidden',
              boxShadow: '0 1px 2px rgba(0,0,0,0.02)',
            }}>
              <input
                type="date" value={from} onChange={e => setFrom(e.target.value)}
                style={{
                  fontSize: '12px', padding: '6px 10px', border: 'none',
                  background: '#FFFFFF', color: from ? '#1E293B' : '#94A3B8',
                  fontFamily: 'monospace', outline: 'none', width: '135px',
                }}
              />
              <span style={{
                padding: '0 8px', color: '#94A3B8', fontSize: '12px',
                background: '#F8FAFC',
                borderLeft: '1px solid #E2E8F0', borderRight: '1px solid #E2E8F0',
              }}>→</span>
              <input
                type="date" value={to} onChange={e => setTo(e.target.value)}
                style={{
                  fontSize: '12px', padding: '6px 10px', border: 'none',
                  background: '#FFFFFF', color: to ? '#1E293B' : '#94A3B8',
                  fontFamily: 'monospace', outline: 'none', width: '135px',
                }}
              />
            </div>
          </div>

          <button
            onClick={fetchAlerts}
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

          {(agentId || severity || isRead || from || to) && (
            <button
              onClick={() => { setAgentId(''); setSeverity(''); setIsRead(''); setFrom(''); setTo('') }}
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
              <col style={{ width: '20px' }}/>
              <col style={{ width: '120px' }}/>
              <col style={{ width: '140px' }}/>
              <col style={{ width: '130px' }}/>
              <col style={{ minWidth: '200px' }}/>
              <col style={{ width: '110px' }}/>
              <col style={{ width: '120px' }}/>
              <col style={{ width: '110px' }}/>
            </colgroup>
            <thead>
              <tr style={{ borderBottom: '2px solid #E2E8F0', background: '#F8FAFC' }}>
                <th style={S.th}></th>
                <th style={S.th}>Time</th>
                <th style={{ ...S.th, padding: '6px 12px' }}>
                  <DropdownFilter label="Severity" options={severityOptions} value={severity} onChange={setSeverity}/>
                </th>
                <th style={S.th}>Type</th>
                <th style={S.th}>Message</th>
                <th style={S.th}>Agent</th>
                <th style={{ ...S.th, padding: '6px 12px' }}>
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
                      borderBottom: '1px solid #E2E8F0',
                      background: !alert.is_read ? s.row : 'transparent',
                      borderLeft: !alert.is_read ? `3px solid ${s.dot}` : '3px solid transparent',
                      cursor: 'pointer',
                    }}
                    onMouseEnter={e => e.currentTarget.style.background = !alert.is_read ? s.row : '#F8FAFC'}
                    onMouseLeave={e => e.currentTarget.style.background = !alert.is_read ? s.row : 'transparent'}
                  >
                    <td style={{ padding: '12px 0 12px 12px' }}>
                      <span style={{ width: '6px', height: '6px', borderRadius: '50%', background: s.dot, display: 'inline-block' }}/>
                    </td>
                    <td style={{ ...S.td, color: '#64748B' }}>{formatTime(alert.created_at)}</td>
                    <td style={S.td}>
                      <Badge label={alert.severity?.toUpperCase()} bg={s.bg} text={s.text}/>
                    </td>
                    <td style={{ ...S.td, color: '#1E293B', fontWeight: '500' }}>{alert.type}</td>
                    <td style={{ ...S.td, color: s.text, fontWeight: !alert.is_read ? '600' : '400' }}>{alert.message}</td>
                    <td style={S.td}>
                      <Badge label={alert.agent_name} bg='#F1F5F9' text='#475569'/>
                    </td>
                    <td style={S.td}>
                      {alert.is_read
                        ? <span style={{ fontSize: '11px', color: '#94A3B8', fontFamily: 'monospace' }}>read</span>
                        : <span style={{ fontSize: '11px', color: '#3B82F6', fontFamily: 'monospace', display: 'flex', alignItems: 'center', gap: '4px', fontWeight: '600' }}>
                            <span style={{ width: '5px', height: '5px', borderRadius: '50%', background: '#3B82F6', display: 'inline-block' }}/>
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
                            fontSize: '11px', padding: '4px 10px', borderRadius: '4px',
                            border: '1px solid #10B981',
                            background: '#E8F5E9',
                            color: markingId === alert.id ? '#94A3B8' : '#15803D',
                            cursor: markingId === alert.id ? 'not-allowed' : 'pointer',
                            fontFamily: 'monospace', fontWeight: '600',
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

        {alerts.length === 0 && (
          <div style={{ padding: '48px', textAlign: 'center', color: '#94A3B8', fontFamily: 'monospace', fontSize: '13px' }}>
            No alerts found
          </div>
        )}
      </div>
    </div>
  )
}