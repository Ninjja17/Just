import { useEffect, useState } from 'react'
import toast from 'react-hot-toast'
import { socialAPI } from '@/services/api'
import type { User } from '@/types'
import { UserPlusIcon, UserMinusIcon } from '@heroicons/react/24/outline'

export default function Social() {
  const [followers, setFollowers] = useState<User[]>([])
  const [following, setFollowing] = useState<User[]>([])
  const [leaderboard, setLeaderboard] = useState<any>(null)
  const [feed, setFeed] = useState<any>(null)
  const [loading, setLoading] = useState(true)
  const [followId, setFollowId] = useState('')
  const [busy, setBusy] = useState(false)

  const load = async () => {
    setLoading(true)
    try {
      const [fr, fg, lb, fd] = await Promise.all([
        socialAPI.getFollowers(),
        socialAPI.getFollowing(),
        socialAPI.getLeaderboard(),
        socialAPI.getFeed(),
      ])
      setFollowers((fr.data as User[]) || [])
      setFollowing((fg.data as User[]) || [])
      setLeaderboard(lb.data)
      setFeed(fd.data)
    } catch (e: any) {
      toast.error(e?.response?.data?.error || 'Failed to load social data')
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    load()
  }, [])

  const handleFollow = async () => {
    if (!followId.trim()) return
    setBusy(true)
    try {
      await socialAPI.follow(followId.trim())
      toast.success('Followed')
      setFollowId('')
      load()
    } catch (e: any) {
      toast.error(e?.response?.data?.error || 'Could not follow')
    } finally {
      setBusy(false)
    }
  }

  const handleUnfollow = async (id: string) => {
    if (!confirm('Unfollow this user?')) return
    try {
      await socialAPI.unfollow(id)
      toast.success('Unfollowed')
      setFollowing((prev) => prev.filter((u) => u.id !== id))
    } catch (e: any) {
      toast.error(e?.response?.data?.error || 'Could not unfollow')
    }
  }

  return (
    <div>
      <h1 className="text-3xl font-bold mb-6">Social</h1>

      <div className="bg-white p-5 rounded-lg shadow mb-6">
        <h2 className="font-semibold mb-3">Follow a user</h2>
        <div className="flex space-x-2">
          <input
            value={followId}
            onChange={(e) => setFollowId(e.target.value)}
            placeholder="Paste user ID"
            className="flex-1 px-3 py-2 border rounded-md"
          />
          <button
            onClick={handleFollow}
            disabled={busy}
            className="inline-flex items-center px-4 py-2 bg-primary-600 text-white rounded-md hover:bg-primary-700 disabled:opacity-50"
          >
            <UserPlusIcon className="w-5 h-5 mr-2" />
            Follow
          </button>
        </div>
        <p className="text-xs text-gray-500 mt-2">
          Tip: ask a friend for their account ID (visible on their Profile page once we wire it).
        </p>
      </div>

      {loading ? (
        <div className="text-gray-500">Loading...</div>
      ) : (
        <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
          <div className="bg-white p-5 rounded-lg shadow">
            <h2 className="font-semibold mb-3">Following ({following.length})</h2>
            {following.length === 0 ? (
              <p className="text-sm text-gray-500">Not following anyone yet.</p>
            ) : (
              <ul className="divide-y">
                {following.map((u) => (
                  <li key={u.id} className="py-2 flex items-center justify-between">
                    <div>
                      <p className="font-medium">{u.name}</p>
                      <p className="text-xs text-gray-500">{u.email}</p>
                    </div>
                    <button
                      onClick={() => handleUnfollow(u.id)}
                      className="inline-flex items-center px-2 py-1 text-sm text-red-600 hover:bg-red-50 rounded"
                    >
                      <UserMinusIcon className="w-4 h-4 mr-1" />
                      Unfollow
                    </button>
                  </li>
                ))}
              </ul>
            )}
          </div>

          <div className="bg-white p-5 rounded-lg shadow">
            <h2 className="font-semibold mb-3">Followers ({followers.length})</h2>
            {followers.length === 0 ? (
              <p className="text-sm text-gray-500">No followers yet.</p>
            ) : (
              <ul className="divide-y">
                {followers.map((u) => (
                  <li key={u.id} className="py-2">
                    <p className="font-medium">{u.name}</p>
                    <p className="text-xs text-gray-500">{u.email}</p>
                  </li>
                ))}
              </ul>
            )}
          </div>

          <div className="bg-white p-5 rounded-lg shadow">
            <h2 className="font-semibold mb-3">Leaderboard</h2>
            <pre className="text-xs text-gray-600 whitespace-pre-wrap">
              {leaderboard ? JSON.stringify(leaderboard, null, 2) : 'No data'}
            </pre>
          </div>

          <div className="bg-white p-5 rounded-lg shadow">
            <h2 className="font-semibold mb-3">Activity feed</h2>
            <pre className="text-xs text-gray-600 whitespace-pre-wrap">
              {feed ? JSON.stringify(feed, null, 2) : 'No activity'}
            </pre>
          </div>
        </div>
      )}
    </div>
  )
}
