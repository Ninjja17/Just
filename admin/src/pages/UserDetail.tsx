import { useEffect, useState } from 'react'
import { Link, useNavigate, useParams } from 'react-router-dom'
import toast from 'react-hot-toast'
import { adminAPI } from '@/api'

export default function UserDetail() {
  const { id } = useParams<{ id: string }>()
  const navigate = useNavigate()
  const [data, setData] = useState<any>(null)
  const [loading, setLoading] = useState(true)

  const me = (() => {
    try {
      return JSON.parse(localStorage.getItem('admin_user') || 'null')
    } catch {
      return null
    }
  })()

  useEffect(() => {
    if (!id) return
    adminAPI
      .getUser(id)
      .then((r) => setData(r.data))
      .catch((e) => toast.error(e?.response?.data?.error || 'Failed to load user'))
      .finally(() => setLoading(false))
  }, [id])

  const remove = async () => {
    if (!id) return
    if (!confirm('Permanently delete this user and ALL their data? This cannot be undone.')) return
    try {
      await adminAPI.deleteUser(id)
      toast.success('User deleted')
      navigate('/users')
    } catch (e: any) {
      toast.error(e?.response?.data?.error || 'Delete failed')
    }
  }

  if (loading) return <p className="text-gray-500">Loading...</p>
  if (!data) return <p className="text-gray-500">Not found</p>

  const u = data.user
  return (
    <div>
      <Link to="/users" className="text-sm text-admin-600 hover:underline">
        ← Back to users
      </Link>
      <div className="flex items-start justify-between mt-2 mb-6">
        <div>
          <h1 className="text-3xl font-bold">{u.name}</h1>
          <p className="text-gray-500">{u.email}</p>
          <p className="text-xs text-gray-400 mt-1">ID: {u.id}</p>
        </div>
        {me?.role === 'superadmin' && (
          <button
            onClick={remove}
            className="px-4 py-2 bg-red-600 text-white rounded-md hover:bg-red-700 text-sm"
          >
            Delete user
          </button>
        )}
      </div>

      <div className="grid grid-cols-2 md:grid-cols-4 gap-4 mb-6">
        <div className="bg-white p-4 rounded-lg shadow">
          <p className="text-xs uppercase text-gray-500">Status</p>
          <p className="text-lg font-semibold">
            {u.is_verified ? 'Verified' : 'Pending'}
          </p>
        </div>
        <div className="bg-white p-4 rounded-lg shadow">
          <p className="text-xs uppercase text-gray-500">Joined</p>
          <p className="text-lg font-semibold">
            {u.created_at ? new Date(u.created_at).toLocaleDateString() : '-'}
          </p>
        </div>
        <div className="bg-white p-4 rounded-lg shadow">
          <p className="text-xs uppercase text-gray-500">Sessions</p>
          <p className="text-lg font-semibold">{data.session_count}</p>
        </div>
        <div className="bg-white p-4 rounded-lg shadow">
          <p className="text-xs uppercase text-gray-500">Total hours</p>
          <p className="text-lg font-semibold">{(data.total_minutes / 60).toFixed(1)}</p>
        </div>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        <div className="bg-white p-5 rounded-lg shadow">
          <h2 className="font-semibold mb-3">Skills</h2>
          {data.skills?.length ? (
            <table className="w-full text-sm">
              <thead className="text-left text-gray-500">
                <tr>
                  <th className="py-1">Name</th>
                  <th className="py-1">Category</th>
                  <th className="py-1 text-right">Hours / Target</th>
                </tr>
              </thead>
              <tbody>
                {data.skills.map((s: any) => (
                  <tr key={s.id} className="border-t">
                    <td className="py-2">{s.name}</td>
                    <td className="py-2 text-gray-600">{s.category}</td>
                    <td className="py-2 text-right">
                      {(s.total_minutes / 60).toFixed(1)} / {s.target_hours}
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          ) : (
            <p className="text-sm text-gray-500">No skills</p>
          )}
        </div>

        <div className="bg-white p-5 rounded-lg shadow">
          <h2 className="font-semibold mb-3">Recent sessions</h2>
          {data.recent_sessions?.length ? (
            <table className="w-full text-sm">
              <thead className="text-left text-gray-500">
                <tr>
                  <th className="py-1">Skill</th>
                  <th className="py-1">Start</th>
                  <th className="py-1 text-right">Min</th>
                </tr>
              </thead>
              <tbody>
                {data.recent_sessions.map((s: any) => (
                  <tr key={s.id} className="border-t">
                    <td className="py-2">{s.skill_name}</td>
                    <td className="py-2 text-gray-600">
                      {s.start_time ? new Date(s.start_time).toLocaleString() : '-'}
                    </td>
                    <td className="py-2 text-right">{s.duration_minutes}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          ) : (
            <p className="text-sm text-gray-500">No sessions</p>
          )}
        </div>
      </div>
    </div>
  )
}
