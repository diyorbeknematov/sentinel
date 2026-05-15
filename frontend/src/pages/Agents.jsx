import { useState, useEffect, useCallback } from 'react'
import api from '../api/axios'

// Och rangli dizayn (Light Mode) uchun panel va karta stillari
const S = {
  card: {
    background: '#FFFFFF',
    border: '1px solid #E2E8F0',
    borderRadius: '12px',
    padding: '20px',
    boxShadow: '0 1px 3px rgba(0,0,0,0.02), 0 1px 2px rgba(0,0,0,0.04)',
    boxSizing: 'border-box',
    width: '100%',
  },
}

function formatDate(iso) {
  if (!iso) return '—'
  return new Date(iso).toLocaleDateString('en', {
    month: 'short', day: 'numeric', year: 'numeric',
  })
}

function timeAgo(iso) {
  if (!iso) return 'hech qachon'
  const diff = Math.floor((Date.now() - new Date(iso)) / 1000)
  if (diff < 60)    return `${diff}s oldin`
  if (diff < 3600)  return `${Math.floor(diff / 60)}m oldin`
  if (diff < 86400) return `${Math.floor(diff / 3600)}s oldin`
  return `${Math.floor(diff / 86400)}kun oldin`
}

function StatusDot({ status }) {
  const online = status === 'online'
  return (
    <div style={{ display: 'flex', alignItems: 'center', gap: '6px' }}>
      <span style={{
        width: '7px', height: '7px', borderRadius: '50%',
        background: online ? '#10B981' : '#64748B', // Yashil va Kulrang
        display: 'inline-block',
        boxShadow: online ? '0 0 0 2px rgba(16,185,129,0.2)' : 'none',
      }}/>
      <span style={{ fontSize: '11px', fontFamily: 'monospace', color: online ? '#059669' : '#64748B', fontWeight: '500' }}>
        {online ? 'online' : 'offline'}
      </span>
    </div>
  )
}

export default function Agents() {
  const [agents,  setAgents]  = useState([])
  const [loading, setLoading] = useState(true)
  const [search,  setSearch]  = useState('')
  const [filter,  setFilter]  = useState('all') // all | online | offline

  const load = useCallback(async () => {
    try {
      const res = await api.get('/agents')
      setAgents(res.data?.data || [])
    } catch {
      setAgents([])
    } finally {
      setLoading(false)
    }
  }, [])

  useEffect(() => {
    load()
    const t = setInterval(load, 10000)
    return () => clearInterval(t)
  }, [load])

  const filtered = agents.filter(a => {
    const matchSearch = a.name?.toLowerCase().includes(search.toLowerCase())
      || a.ip_address?.includes(search)
    const matchFilter = filter === 'all'
      || (filter === 'online'  && a.status === 'online')
      || (filter === 'offline' && a.status !== 'online')
    return matchSearch && matchFilter
  })

  const onlineCount  = agents.filter(a => a.status === 'online').length
  const offlineCount = agents.length - onlineCount

  return (
    <div style={{ minHeight: '100vh', width: '100%', boxSizing: 'border-box' }}>

      {/* Navbar */}
      <div style={{
        position: 'sticky', top: 0, zIndex: 30,
        background: 'linear-gradient(90deg, rgba(79,70,229,0.08), rgba(16,185,129,0.05))',
        backdropFilter: 'blur(12px)',
        borderBottom: '1px solid rgba(148,163,184,0.2)',
        height: '56px', display: 'flex', alignItems: 'center', justifyContent: 'space-between',
        marginLeft: '-24px', marginRight: '-24px',
        paddingLeft: '24px', paddingRight: '24px', marginBottom: '24px',
      }}>
        <div style={{ display: 'flex', alignItems: 'center', gap: '8px' }}>
          <span style={{ fontSize: '16px', fontWeight: '600', color: '#1E293B' }}>Agents</span>
          <span style={{ fontSize: '12px', color: '#64748B', fontFamily: 'monospace' }}>
            · {agents.length} ta
          </span>
        </div>
        <div style={{
          display: 'flex', alignItems: 'center', gap: '6px',
          padding: '5px 14px', borderRadius: '20px',
          background: '#E8F5E9', border: '1px solid #A5D6A7',
        }}>
          <span style={{ width: '6px', height: '6px', borderRadius: '50%', background: '#2E7D32', display: 'inline-block' }}/>
          <span style={{ fontSize: '12px', color: '#2E7D32', fontFamily: 'monospace', fontWeight: '600' }}>Live</span>
        </div>
      </div>

      {/* Stat cards */}
      <div style={{ display: 'grid', gridTemplateColumns: 'repeat(3, 1fr)', gap: '16px', marginBottom: '20px', width: '100%' }}>
        {[
          { label: 'JAMI',    value: agents.length, color: '#4F46E5', badgeBg: '#EEF2FF', valColor: '#1E293B' },
          { label: 'ONLINE',  value: onlineCount,   color: '#10B981', badgeBg: '#E8F5E9', valColor: '#059669' },
          { label: 'OFFLINE', value: offlineCount,  color: '#64748B', badgeBg: '#F1F5F9', valColor: '#475569' },
        ].map(c => (
          <div key={c.label} style={{ ...S.card, borderTop: `4px solid ${c.color}` }}>
            <p style={{ fontSize: '11px', color: '#64748B', letterSpacing: '0.7px', textTransform: 'uppercase', marginBottom: '10px', fontWeight: '500' }}>
              {c.label}
            </p>
            <p style={{ fontSize: '28px', fontWeight: '600', fontFamily: 'monospace', color: c.valColor, lineHeight: 1 }}>
              {c.value}
            </p>
          </div>
        ))}
      </div>

      {/* Toolbar */}
      <div style={{ display: 'flex', gap: '12px', marginBottom: '16px', alignItems: 'center', width: '100%' }}>
        <input
          value={search}
          onChange={e => setSearch(e.target.value)}
          placeholder="Qidirish — nom yoki IP..."
          style={{
            flex: 1, fontSize: '13px', padding: '10px 16px',
            borderRadius: '8px', border: '1px solid #E2E8F0',
            background: '#FFFFFF', color: '#1E293B',
            fontFamily: 'monospace', outline: 'none',
            boxShadow: '0 1px 2px rgba(0,0,0,0.02)',
          }}
        />
        <div style={{ display: 'flex', gap: '6px' }}>
          {['all', 'online', 'offline'].map(f => (
            <button
              key={f}
              onClick={() => setFilter(f)}
              style={{
                fontSize: '12px', padding: '8px 16px', borderRadius: '6px',
                border: filter === f ? '1px solid #4F46E5' : '1px solid #E2E8F0',
                background: filter === f ? '#EEF2FF' : '#FFFFFF',
                color: filter === f ? '#4F46E5' : '#475569',
                cursor: 'pointer', fontFamily: 'monospace', transition: 'all 0.15s',
                fontWeight: filter === f ? '600' : '500',
                boxShadow: '0 1px 2px rgba(0,0,0,0.02)',
              }}
            >
              {f === 'all' ? 'Barchasi' : f === 'online' ? 'Online' : 'Offline'}
            </button>
          ))}
        </div>
      </div>

      {/* Table */}
      <div style={{ ...S.card, width: '100%' }}>
        <div style={{ overflowX: 'auto', width: '100%' }}>
          <table style={{ width: '100%', borderCollapse: 'collapse' }}>
            <thead>
              <tr style={{ borderBottom: '2px solid #E2E8F0' }}>
                {['', 'Nomi', 'IP manzil', 'Oxirgi faollik', "Qo'shilgan", ''].map((h, i) => (
                  <th key={i} style={{
                    textAlign: 'left', fontSize: '12px', color: '#64748B',
                    fontWeight: '600', textTransform: 'uppercase', letterSpacing: '0.5px',
                    padding: '0 12px 12px',
                  }}>
                    {h}
                  </th>
                ))}
              </tr>
            </thead>
            <tbody>
              {loading ? (
                <tr>
                  <td colSpan={6} style={{ padding: '48px', textAlign: 'center', color: '#94A3B8', fontFamily: 'monospace', fontSize: '13px' }}>
                    Yuklanmoqda...
                  </td>
                </tr>
              ) : filtered.length === 0 ? (
                <tr>
                  <td colSpan={6} style={{ padding: '48px', textAlign: 'center', color: '#94A3B8', fontFamily: 'monospace', fontSize: '13px' }}>
                    {agents.length === 0 ? "No Agents Found" : 'Hech narsa topilmadi'}
                  </td>
                </tr>
              ) : filtered.map((a, i) => {
                const online = a.status === 'online'
                return (
                  <tr
                    key={a.id}
                    style={{
                      borderBottom: '1px solid #E2E8F0',
                      transition: 'background 0.15s',
                    }}
                    onMouseEnter={e => e.currentTarget.style.background = '#F8FAFC'}
                    onMouseLeave={e => e.currentTarget.style.background = 'transparent'}
                  >
                    {/* Status dot */}
                    <td style={{ padding: '14px 12px', width: '32px' }}>
                      <span style={{
                        width: '7px', height: '7px', borderRadius: '50%',
                        background: online ? '#10B981' : '#94A3B8',
                        display: 'inline-block',
                        boxShadow: online ? '0 0 0 2px rgba(16,185,129,0.15)' : 'none',
                      }}/>
                    </td>

                    {/* Name + status */}
                    <td style={{ padding: '14px 12px' }}>
                      <span style={{ fontSize: '14px', color: '#1E293B', display: 'block', fontWeight: '500' }}>
                        {a.name || '—'}
                      </span>
                      <div style={{ marginTop: '4px' }}>
                        <StatusDot status={a.status} />
                      </div>
                    </td>

                    {/* IP */}
                    <td style={{ padding: '14px 12px' }}>
                      <span style={{ fontSize: '13px', color: '#475569', fontFamily: 'monospace' }}>
                        {a.ip_address || '—'}
                      </span>
                    </td>

                    {/* Last seen */}
                    <td style={{ padding: '14px 12px' }}>
                      <span style={{ fontSize: '13px', fontFamily: 'monospace', color: online ? '#059669' : '#64748B', fontWeight: '500' }}>
                        {timeAgo(a.last_seen)}
                      </span>
                    </td>

                    {/* Created at */}
                    <td style={{ padding: '14px 12px' }}>
                      <span style={{ fontSize: '13px', color: '#64748B', fontFamily: 'monospace' }}>
                        {formatDate(a.created_at)}
                      </span>
                    </td>

                    {/* Action */}
                    <td style={{ padding: '14px 12px', textAlign: 'right' }}>
                      <a
                        href={`/agents/${a.id}`}
                        style={{
                          fontSize: '12px', color: '#475569', textDecoration: 'none',
                          padding: '6px 12px', borderRadius: '6px',
                          border: '1px solid #E2E8F0', fontFamily: 'monospace',
                          background: '#FFFFFF', transition: 'all 0.15s',
                          display: 'inline-block',
                        }}
                        onMouseEnter={e => { 
                          e.target.style.color = '#4F46E5'; 
                          e.target.style.borderColor = '#4F46E5';
                          e.target.style.background = '#EEF2FF';
                        }}
                        onMouseLeave={e => { 
                          e.target.style.color = '#475569'; 
                          e.target.style.borderColor = '#E2E8F0';
                          e.target.style.background = '#FFFFFF';
                        }}
                      >
                        ko'rish →
                      </a>
                    </td>
                  </tr>
                )
              })}
            </tbody>
          </table>
        </div>
      </div>
    </div>
  )
}