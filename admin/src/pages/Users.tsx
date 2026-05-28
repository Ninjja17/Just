import { useEffect, useState } from 'react'
import { Link } from 'react-router-dom'
import toast from 'react-hot-toast'
import { adminAPI, UserRow } from '@/api'

const PAGE_SIZE = 25

export default function Users() {
  const [users, setUsers] = useState<UserRow[]>([])
  const [total, setTotal] = useState(0)
  const [search, setSearch] = useState('')
  const [page, setPage] = useState(0)
  const [loading, setLoading] = useState(true)

  const load = async () => {
    setLoading(true)
    try {
      const res = await adminAPI.listUsers({
        search,
        limit: PAGE_SIZE,
        offset: page * PAGE_SIZE,
      })
      setUsers(res.data.users)
      setTotal(res.data.total)
    } catch (e: any) {
      toast.error(e?.response?.data?.error || 'Failed to load users')
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    load()
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [page])

  const onSearch = (e: React.FormEvent) => {
    e.preventDefault()
    setPage(0)
    load()
  }

  const totalPages = Math.max(1, Math.ceil(total / PAGE_SIZE))

  return (
    <div>
      <div className="flex items-center justify-between mb-6">
        <h1 className="text-3xl font-bold">Users</h1>
        <form onSubmit={onSearch} className="flex space-x-2">
          <input
            value={search}
            onChange={(e) => setSearch(e.target.value)}
            placeholder="Search email or name"
            className="px-3 py-2 border rounded-md w-64"
          />
          <button className="px-4 py-2 bg-admin-500 text-white rounded-md hover:bg-admin-600">
            Search
          </button>
        </form>
      </div>

      <div className="bg-white rounded-lg shadow overflow-hidden">
        {loading ? (
          <p className="p-5 text-gray-500">Loading users...</p>
        ) : users.length === 0 ? (
          <p className="p-5 text-gray-500">No users found.</p>
        ) : (
          <table className="min-w-full text-sm">
            <thead className="bg-gray-50 text-left text-gray-600">
              <tr>
                <th className="px-4 py-2">Email</th>
                <th className="px-4 py-2">Name</th>
                <th className="px-4 py-2">Verified</th>
                <th className="px-4 py-2">Joined</th>
                <th className="px-4 py-2 text-right">Skills</th>
                <th className="px-4 py-2 text-right">Sessions</th>
                <th className="px-4 py-2 text-right">Hours</th>
              </tr>
            </thead>
            <tbody>
              {users.map((u) => (
                <tr key={u.id} className="border-t hover:bg-gray-50">
                  <td className="px-4 py-2">
                    <Link to={`/users/${u.id}`} className="text-admin-600 hover:underline">
                      {u.email}
                    </Link>
                  </td>
                  <td className="px-4 py-2">{u.name}</td>
                  <td className="px-4 py-2">
                    {u.is_verified ? (
                      <span className="text-emerald-600">✓</span>
                    ) : (
                      <span className="text-amber-600">pending</span>
                    )}
                  </td>
                  <td className="px-4 py-2">
                    {u.created_at ? new Date(u.created_at).toLocaleDateString() : '-'}
                  </td>
                  <td className="px-4 py-2 text-right">{u.skill_count}</td>
                  <td className="px-4 py-2 text-right">{u.session_count}</td>
                  <td className="px-4 py-2 text-right">{(u.total_minutes / 60).toFixed(1)}</td>
                </tr>
              ))}
            </tbody>
          </table>
        )}
      </div>

      <div className="flex items-center justify-between mt-4 text-sm text-gray-600">
        <span>
          Page {page + 1} of {totalPages} — {total} user{total === 1 ? '' : 's'}
        </span>
        <div className="space-x-2">
          <button
            onClick={() => setPage((p) => Math.max(0, p - 1))}
            disabled={page === 0}
            className="px-3 py-1.5 border rounded-md disabled:opacity-40"
          >
            Prev
          </button>
          <button
            onClick={() => setPage((p) => (p + 1 < totalPages ? p + 1 : p))}
            disabled={page + 1 >= totalPages}
            className="px-3 py-1.5 border rounded-md disabled:opacity-40"
          >
            Next
          </button>
        </div>
      </div>
    </div>
  )
}
