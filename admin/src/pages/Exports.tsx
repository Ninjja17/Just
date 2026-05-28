import { adminAPI, downloadCSV } from '@/api'

export default function Exports() {
  return (
    <div>
      <h1 className="text-3xl font-bold mb-6">Exports</h1>
      <div className="bg-white p-6 rounded-lg shadow max-w-xl space-y-4">
        <div>
          <h2 className="font-semibold">All users (CSV)</h2>
          <p className="text-sm text-gray-500 mb-2">
            Includes email, name, verification, join date, skill/session counts and total minutes.
          </p>
          <button
            onClick={() => downloadCSV(adminAPI.exportUsersCSV(), 'users.csv')}
            className="px-4 py-2 bg-admin-500 text-white rounded-md hover:bg-admin-600"
          >
            Download users.csv
          </button>
        </div>
        <div className="border-t pt-4">
          <h2 className="font-semibold">All sessions (CSV)</h2>
          <p className="text-sm text-gray-500 mb-2">
            Most recent 10,000 sessions joined with user email + skill name.
          </p>
          <button
            onClick={() => downloadCSV(adminAPI.exportSessionsCSV(), 'sessions.csv')}
            className="px-4 py-2 bg-admin-500 text-white rounded-md hover:bg-admin-600"
          >
            Download sessions.csv
          </button>
        </div>
      </div>
    </div>
  )
}
