import { useAuth } from '@/context/AuthContext'

export default function Profile() {
  const { user } = useAuth()

  return (
    <div>
      <h1 className="text-3xl font-bold mb-8">Profile</h1>
      
      <div className="bg-white p-6 rounded-lg shadow max-w-2xl">
        <div className="space-y-4">
          <div>
            <label className="text-sm font-medium text-gray-500">Name</label>
            <p className="text-lg">{user?.name}</p>
          </div>
          <div>
            <label className="text-sm font-medium text-gray-500">Email</label>
            <p className="text-lg">{user?.email}</p>
          </div>
          <div>
            <label className="text-sm font-medium text-gray-500">Status</label>
            <p className="text-lg">
              {user?.is_verified ? (
                <span className="text-green-600">Verified</span>
              ) : (
                <span className="text-yellow-600">Not Verified</span>
              )}
            </p>
          </div>
          <div>
            <label className="text-sm font-medium text-gray-500">Member Since</label>
            <p className="text-lg">
              {user?.created_at ? new Date(user.created_at).toLocaleDateString() : 'N/A'}
            </p>
          </div>
          <div>
            <label className="text-sm font-medium text-gray-500">User ID (share to be followed)</label>
            <div className="flex items-center space-x-2">
              <code className="text-sm bg-gray-100 px-2 py-1 rounded break-all">{user?.id}</code>
              <button
                onClick={() => {
                  if (user?.id) {
                    navigator.clipboard.writeText(user.id)
                  }
                }}
                className="text-sm text-primary-600 hover:underline"
              >
                Copy
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
  )
}
