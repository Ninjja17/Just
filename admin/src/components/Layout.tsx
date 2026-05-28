import { NavLink, Outlet, useNavigate } from 'react-router-dom'

export default function Layout() {
  const navigate = useNavigate()
  const me = (() => {
    try {
      return JSON.parse(localStorage.getItem('admin_user') || 'null')
    } catch {
      return null
    }
  })()

  const logout = () => {
    localStorage.removeItem('admin_token')
    localStorage.removeItem('admin_user')
    navigate('/login')
  }

  const linkClass = ({ isActive }: { isActive: boolean }) =>
    `block px-4 py-2 rounded-md text-sm ${
      isActive ? 'bg-admin-500 text-white' : 'text-gray-700 hover:bg-admin-100'
    }`

  return (
    <div className="min-h-screen flex">
      <aside className="w-60 bg-white border-r shadow-sm flex flex-col">
        <div className="px-5 py-4 border-b">
          <p className="text-xs uppercase text-gray-500">10,000 HR</p>
          <h1 className="font-bold text-lg">Admin Console</h1>
        </div>
        <nav className="p-3 space-y-1 flex-1">
          <NavLink to="/" end className={linkClass}>
            Dashboard
          </NavLink>
          <NavLink to="/users" className={linkClass}>
            Users
          </NavLink>
          <NavLink to="/exports" className={linkClass}>
            Exports
          </NavLink>
        </nav>
        <div className="p-4 border-t text-xs text-gray-600">
          <p className="font-medium">{me?.name}</p>
          <p className="text-gray-500">{me?.email}</p>
          <p className="text-gray-400 mt-1">Role: {me?.role}</p>
          <button
            onClick={logout}
            className="mt-3 w-full px-3 py-1.5 text-sm rounded-md bg-gray-100 hover:bg-gray-200"
          >
            Sign out
          </button>
        </div>
      </aside>
      <main className="flex-1 p-8 overflow-auto">
        <Outlet />
      </main>
    </div>
  )
}
