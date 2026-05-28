import { useEffect, useState } from 'react'
import toast from 'react-hot-toast'
import {
  BarChart,
  Bar,
  LineChart,
  Line,
  XAxis,
  YAxis,
  Tooltip,
  CartesianGrid,
  ResponsiveContainer,
} from 'recharts'
import { adminAPI, Stats } from '@/api'

function Tile({ label, value }: { label: string; value: string | number }) {
  return (
    <div className="bg-white p-5 rounded-lg shadow">
      <p className="text-xs uppercase text-gray-500">{label}</p>
      <p className="text-3xl font-bold mt-1">{value}</p>
    </div>
  )
}

export default function Dashboard() {
  const [stats, setStats] = useState<Stats | null>(null)
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    adminAPI
      .stats()
      .then((r) => setStats(r.data))
      .catch((e) => toast.error(e?.response?.data?.error || 'Failed to load stats'))
      .finally(() => setLoading(false))
  }, [])

  if (loading) return <p className="text-gray-500">Loading dashboard...</p>
  if (!stats) return <p className="text-gray-500">No data</p>

  const hours = (stats.total_minutes / 60).toFixed(1)

  return (
    <div>
      <h1 className="text-3xl font-bold mb-6">Overview</h1>

      <div className="grid grid-cols-2 md:grid-cols-4 gap-4 mb-6">
        <Tile label="Total users" value={stats.total_users} />
        <Tile label="Verified" value={stats.verified_users} />
        <Tile label="New (7d)" value={stats.new_users_7d} />
        <Tile label="Active (7d)" value={stats.active_users_7d} />
        <Tile label="Skills" value={stats.total_skills} />
        <Tile label="Sessions" value={stats.total_sessions} />
        <Tile label="Total hours" value={hours} />
        <Tile
          label="Avg min/session"
          value={
            stats.total_sessions > 0
              ? Math.round(stats.total_minutes / stats.total_sessions)
              : 0
          }
        />
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        <div className="bg-white p-5 rounded-lg shadow">
          <h2 className="font-semibold mb-4">Sessions per day (last 30)</h2>
          <div style={{ width: '100%', height: 260 }}>
            <ResponsiveContainer>
              <BarChart data={stats.sessions_per_day}>
                <CartesianGrid strokeDasharray="3 3" />
                <XAxis dataKey="date" tick={{ fontSize: 10 }} interval={2} />
                <YAxis allowDecimals={false} />
                <Tooltip />
                <Bar dataKey="value" fill="#3f51b5" />
              </BarChart>
            </ResponsiveContainer>
          </div>
        </div>

        <div className="bg-white p-5 rounded-lg shadow">
          <h2 className="font-semibold mb-4">New users per day (last 30)</h2>
          <div style={{ width: '100%', height: 260 }}>
            <ResponsiveContainer>
              <LineChart data={stats.new_users_per_day}>
                <CartesianGrid strokeDasharray="3 3" />
                <XAxis dataKey="date" tick={{ fontSize: 10 }} interval={2} />
                <YAxis allowDecimals={false} />
                <Tooltip />
                <Line type="monotone" dataKey="value" stroke="#10b981" strokeWidth={2} />
              </LineChart>
            </ResponsiveContainer>
          </div>
        </div>

        <div className="bg-white p-5 rounded-lg shadow lg:col-span-2">
          <h2 className="font-semibold mb-4">Top skills</h2>
          {stats.top_skills && stats.top_skills.length > 0 ? (
            <div style={{ width: '100%', height: 260 }}>
              <ResponsiveContainer>
                <BarChart data={stats.top_skills} layout="vertical">
                  <CartesianGrid strokeDasharray="3 3" />
                  <XAxis type="number" allowDecimals={false} />
                  <YAxis dataKey="name" type="category" width={120} />
                  <Tooltip />
                  <Bar dataKey="count" fill="#3f51b5" />
                </BarChart>
              </ResponsiveContainer>
            </div>
          ) : (
            <p className="text-gray-500 text-sm">No skills logged yet.</p>
          )}
        </div>
      </div>
    </div>
  )
}
