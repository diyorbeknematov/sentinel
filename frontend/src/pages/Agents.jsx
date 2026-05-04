import { useState, useEffect, useCallback } from 'react'
import api from '../api/axios'

const S = {
  card: {
    background: '#0d1120',
    border: '1px solid #1e293b',
    borderRadius: '12px',
    padding: '16px',
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

// ✅ isOnline yo'q — to'g'ridan status ishlatiladi
function StatusDot({ status }) {
  const online = status === 'online'
  return (
    <div style={{ display: 'flex', alignItems: 'center', gap: '6px' }}>
      <span style={{
        width: '7px', height: '7px', borderRadius: '50%',
        background: online ? '#22c55e' : '#475569',
        display: 'inline-block',
        boxShadow: online ? '0 0 0 2px rgba(34,197,94,0.2)' : 'none',
      }}/>
      <span style={{ fontSize: '11px', fontFamily: 'monospace', color: online ? '#4ade80' : '#475569' }}>
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
    // ✅ a.status dan foydalanamiz
    const matchFilter = filter === 'all'
      || (filter === 'online'  && a.status === 'online')
      || (filter === 'offline' && a.status !== 'online')
    return matchSearch && matchFilter
  })

  // ✅ a.status dan foydalanamiz
  const onlineCount  = agents.filter(a => a.status === 'online').length
  const offlineCount = agents.length - onlineCount

  return (
    <div style={{ minHeight: '100vh' }}>

      {/* Navbar */}
      <div style={{
        position: 'sticky', top: 0, zIndex: 30,
        background: 'rgba(8,11,20,0.85)', backdropFilter: 'blur(12px)',
        borderBottom: '1px solid #1e293b',
        height: '52px', display: 'flex', alignItems: 'center', justifyContent: 'space-between',
        marginLeft: '-24px', marginRight: '-24px',
        paddingLeft: '24px', paddingRight: '24px', marginBottom: '24px',
      }}>
        <div style={{ display: 'flex', alignItems: 'center', gap: '8px' }}>
          <span style={{ fontSize: '14px', fontWeight: '600', color: '#e2e8f0' }}>Agents</span>
          <span style={{ fontSize: '11px', color: '#475569', fontFamily: 'monospace' }}>
            · {agents.length} ta
          </span>
        </div>
        <div style={{ display: 'flex', alignItems: 'center', gap: '6px' }}>
          <span style={{ width: '6px', height: '6px', borderRadius: '50%', background: '#22c55e', display: 'inline-block' }}/>
          <span style={{ fontSize: '10px', color: '#4ade80', fontFamily: 'monospace' }}>Live</span>
        </div>
      </div>

      {/* Stat cards */}
      <div style={{ display: 'grid', gridTemplateColumns: 'repeat(3, 1fr)', gap: '12px', marginBottom: '20px' }}>
        {[
          { label: 'JAMI',    value: agents.length, color: '#6366f1', valColor: '#a5b4fc' },
          { label: 'ONLINE',  value: onlineCount,   color: '#22c55e', valColor: '#4ade80' },
          { label: 'OFFLINE', value: offlineCount,  color: '#475569', valColor: '#94a3b8' },
        ].map(c => (
          <div key={c.label} style={{ ...S.card, borderTop: `2px solid ${c.color}` }}>
            <p style={{ fontSize: '9px', color: '#475569', letterSpacing: '0.7px', textTransform: 'uppercase', marginBottom: '8px' }}>
              {c.label}
            </p>
            <p style={{ fontSize: '26px', fontWeight: '600', fontFamily: 'monospace', color: c.valColor, lineHeight: 1 }}>
              {c.value}
            </p>
          </div>
        ))}
      </div>

      {/* Toolbar */}
      <div style={{ display: 'flex', gap: '10px', marginBottom: '14px', alignItems: 'center' }}>
        <input
          value={search}
          onChange={e => setSearch(e.target.value)}
          placeholder="Qidirish — nom yoki IP..."
          style={{
            flex: 1, fontSize: '12px', padding: '8px 14px',
            borderRadius: '8px', border: '1px solid #1e293b',
            background: '#080b14', color: '#e2e8f0',
            fontFamily: 'monospace', outline: 'none',
          }}
        />
        {['all', 'online', 'offline'].map(f => (
          <button
            key={f}
            onClick={() => setFilter(f)}
            style={{
              fontSize: '11px', padding: '7px 14px', borderRadius: '6px',
              border: filter === f ? '1px solid rgba(99,102,241,0.4)' : '1px solid #1e293b',
              background: filter === f ? 'rgba(99,102,241,0.1)' : 'transparent',
              color: filter === f ? '#a5b4fc' : '#475569',
              cursor: 'pointer', fontFamily: 'monospace', transition: 'all 0.15s',
            }}
          >
            {f === 'all' ? 'Barchasi' : f === 'online' ? 'Online' : 'Offline'}
          </button>
        ))}
      </div>

      {/* Table */}
      <div style={S.card}>
        <table style={{ width: '100%', borderCollapse: 'collapse' }}>
          <thead>
            <tr style={{ borderBottom: '1px solid #1e293b' }}>
              {['', 'Nomi', 'IP manzil', 'Oxirgi faollik', "Qo'shilgan", ''].map((h, i) => (
                <th key={i} style={{
                  textAlign: 'left', fontSize: '10px', color: '#475569',
                  fontWeight: '500', textTransform: 'uppercase', letterSpacing: '0.5px',
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
                <td colSpan={6} style={{ padding: '48px', textAlign: 'center', color: '#334155', fontFamily: 'monospace', fontSize: '12px' }}>
                  Yuklanmoqda...
                </td>
              </tr>
            ) : filtered.length === 0 ? (
              <tr>
                <td colSpan={6} style={{ padding: '48px', textAlign: 'center', color: '#334155', fontFamily: 'monospace', fontSize: '12px' }}>
                  {agents.length === 0 ? "No Agents Found" : 'Hech narsa topilmadi'}
                </td>
              </tr>
            ) : filtered.map((a, i) => {
              const online = a.status === 'online' // ✅
              return (
                <tr
                  key={a.id}
                  style={{
                    borderBottom: i < filtered.length - 1 ? '1px solid rgba(30,41,59,0.5)' : 'none',
                    transition: 'background 0.15s',
                  }}
                  onMouseEnter={e => e.currentTarget.style.background = 'rgba(30,41,59,0.3)'}
                  onMouseLeave={e => e.currentTarget.style.background = 'transparent'}
                >
                  {/* Status dot */}
                  <td style={{ padding: '14px 12px', width: '32px' }}>
                    <span style={{
                      width: '7px', height: '7px', borderRadius: '50%',
                      background: online ? '#22c55e' : '#334155',
                      display: 'inline-block',
                      boxShadow: online ? '0 0 0 2px rgba(34,197,94,0.15)' : 'none',
                    }}/>
                  </td>

                  {/* Name + status */}
                  <td style={{ padding: '14px 12px' }}>
                    <span style={{ fontSize: '13px', color: '#e2e8f0', fontWeight: '500', display: 'block' }}>
                      {a.name || '—'}
                    </span>
                    <div style={{ marginTop: '3px' }}>
                      <StatusDot status={a.status} />  {/* ✅ status prop */}
                    </div>
                  </td>

                  {/* IP */}
                  <td style={{ padding: '14px 12px' }}>
                    <span style={{ fontSize: '12px', color: '#64748b', fontFamily: 'monospace' }}>
                      {a.ip_address || '—'}
                    </span>
                  </td>

                  {/* Last seen */}
                  <td style={{ padding: '14px 12px' }}>
                    <span style={{ fontSize: '12px', fontFamily: 'monospace', color: online ? '#4ade80' : '#475569' }}>
                      {timeAgo(a.last_seen)}
                    </span>
                  </td>

                  {/* Created at */}
                  <td style={{ padding: '14px 12px' }}>
                    <span style={{ fontSize: '12px', color: '#334155', fontFamily: 'monospace' }}>
                      {formatDate(a.created_at)}
                    </span>
                  </td>

                  {/* Action */}
                  <td style={{ padding: '14px 12px', textAlign: 'right' }}>
                    <a
                      href={`/agents/${a.id}`}
                      style={{
                        fontSize: '11px', color: '#475569', textDecoration: 'none',
                        padding: '4px 10px', borderRadius: '4px',
                        border: '1px solid #1e293b', fontFamily: 'monospace',
                      }}
                      onMouseEnter={e => { e.target.style.color = '#a5b4fc'; e.target.style.borderColor = 'rgba(99,102,241,0.3)' }}
                      onMouseLeave={e => { e.target.style.color = '#475569'; e.target.style.borderColor = '#1e293b' }}
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
  )
}