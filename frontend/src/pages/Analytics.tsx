import { useEffect, useState } from 'react'
import toast from 'react-hot-toast'
import {
  BarChart,
  Bar,
  XAxis,
  YAxis,
  Tooltip,
  ResponsiveContainer,
  CartesianGrid,
  PieChart,
  Pie,
  Cell,
  Legend,
} from 'recharts'
import { analyticsAPI, sessionsAPI, skillsAPI } from '@/services/api'
import type { Session, Skill } from '@/types'

interface SkillStat {
  skill_id: string
  skill_name: string
  total_minutes: number
  total_hours: number
}

interface Overview {
  total_skills: number
  total_minutes: number
  total_hours: number
  skills: SkillStat[]
}

const PIE_COLORS = ['#3B82F6', '#10B981', '#F59E0B', '#EF4444', '#8B5CF6', '#EC4899', '#14B8A6']

export default function Analytics() {
  const [overview, setOverview] = useState<Overview | null>(null)
  const [loading, setLoading] = useState(true)
  const [last30, setLast30] = useState<{ date: string; minutes: number }[]>([])

  useEffect(() => {
    const load = async () => {
      try {
        const [ov, sessRes, skRes] = await Promise.all([
          analyticsAPI.getOverview(),
          sessionsAPI.getAll({ limit: 500 }),
          skillsAPI.getAll(),
        ])
        setOverview(ov.data as Overview)

        const sessions: Session[] = sessRes.data || []
        const skills: Skill[] = skRes.data || []
        void skills

        // Build last 30 days bucket
        const buckets: Record<string, number> = {}
        for (let i = 29; i >= 0; i--) {
          const d = new Date()
          d.setDate(d.getDate() - i)
          const key = d.toISOString().slice(0, 10)
          buckets[key] = 0
        }
        sessions.forEach((s) => {
          const key = s.start_time.slice(0, 10)
          if (key in buckets) buckets[key] += s.duration_minutes
        })
        setLast30(
          Object.entries(buckets).map(([date, minutes]) => ({
            date: date.slice(5),
            minutes,
          }))
        )
      } catch (e: any) {
        toast.error(e?.response?.data?.error || 'Failed to load analytics')
      } finally {
        setLoading(false)
      }
    }
    load()
  }, [])

  if (loading) {
    return (
      <div>
        <h1 className="text-3xl font-bold mb-8">Analytics</h1>
        <div className="text-gray-500">Loading analytics...</div>
      </div>
    )
  }

  const pieData =
    overview?.skills
      .filter((s) => s.total_hours > 0)
      .map((s) => ({ name: s.skill_name, value: Math.round(s.total_hours * 10) / 10 })) || []

  return (
    <div>
      <h1 className="text-3xl font-bold mb-8">Analytics</h1>

      <div className="grid grid-cols-1 md:grid-cols-3 gap-4 mb-6">
        <div className="bg-white p-5 rounded-lg shadow">
          <p className="text-sm text-gray-500">Total skills</p>
          <p className="text-3xl font-bold">{overview?.total_skills ?? 0}</p>
        </div>
        <div className="bg-white p-5 rounded-lg shadow">
          <p className="text-sm text-gray-500">Total hours</p>
          <p className="text-3xl font-bold">
            {(overview?.total_hours ?? 0).toFixed(1)}
          </p>
        </div>
        <div className="bg-white p-5 rounded-lg shadow">
          <p className="text-sm text-gray-500">% toward 10,000h</p>
          <p className="text-3xl font-bold">
            {(((overview?.total_hours ?? 0) / 10000) * 100).toFixed(2)}%
          </p>
        </div>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        <div className="bg-white p-5 rounded-lg shadow">
          <h2 className="font-semibold mb-4">Last 30 days (minutes)</h2>
          <div style={{ width: '100%', height: 280 }}>
            <ResponsiveContainer>
              <BarChart data={last30}>
                <CartesianGrid strokeDasharray="3 3" />
                <XAxis dataKey="date" tick={{ fontSize: 10 }} interval={2} />
                <YAxis />
                <Tooltip />
                <Bar dataKey="minutes" fill="#3B82F6" />
              </BarChart>
            </ResponsiveContainer>
          </div>
        </div>

        <div className="bg-white p-5 rounded-lg shadow">
          <h2 className="font-semibold mb-4">Hours by skill</h2>
          {pieData.length === 0 ? (
            <p className="text-sm text-gray-500">Log sessions to see distribution.</p>
          ) : (
            <div style={{ width: '100%', height: 280 }}>
              <ResponsiveContainer>
                <PieChart>
                  <Pie
                    data={pieData}
                    dataKey="value"
                    nameKey="name"
                    outerRadius={100}
                    label
                  >
                    {pieData.map((_, i) => (
                      <Cell key={i} fill={PIE_COLORS[i % PIE_COLORS.length]} />
                    ))}
                  </Pie>
                  <Tooltip />
                  <Legend />
                </PieChart>
              </ResponsiveContainer>
            </div>
          )}
        </div>
      </div>

      <div className="bg-white p-5 rounded-lg shadow mt-6">
        <h2 className="font-semibold mb-4">Per-skill progress</h2>
        {overview && overview.skills.length > 0 ? (
          <div className="space-y-3">
            {overview.skills.map((s) => {
              const pct = Math.min(100, (s.total_hours / 10000) * 100)
              return (
                <div key={s.skill_id}>
                  <div className="flex justify-between text-sm mb-1">
                    <span>{s.skill_name}</span>
                    <span className="text-gray-500">
                      {s.total_hours.toFixed(1)}h / 10,000h
                    </span>
                  </div>
                  <div className="w-full h-2 bg-gray-100 rounded-full">
                    <div
                      className="h-full bg-primary-600 rounded-full"
                      style={{ width: `${pct}%` }}
                    />
                  </div>
                </div>
              )
            })}
          </div>
        ) : (
          <p className="text-sm text-gray-500">No skills yet.</p>
        )}
      </div>
    </div>
  )
}
