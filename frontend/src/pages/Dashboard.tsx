import { useEffect, useState } from 'react'
import { analyticsAPI } from '@/services/api'
import toast from 'react-hot-toast'

export default function Dashboard() {
  const [overview, setOverview] = useState<any>(null)
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    fetchOverview()
  }, [])

  const fetchOverview = async () => {
    try {
      const response = await analyticsAPI.getOverview()
      setOverview(response.data)
    } catch (error) {
      toast.error('Failed to load overview')
    } finally {
      setLoading(false)
    }
  }

  if (loading) {
    return <div>Loading...</div>
  }

  return (
    <div>
      <h1 className="text-3xl font-bold mb-8">Dashboard</h1>
      
      <div className="grid grid-cols-1 md:grid-cols-3 gap-6 mb-8">
        <div className="bg-white p-6 rounded-lg shadow">
          <h3 className="text-gray-500 text-sm font-medium">Total Skills</h3>
          <p className="text-3xl font-bold mt-2">{overview?.total_skills || 0}</p>
        </div>
        <div className="bg-white p-6 rounded-lg shadow">
          <h3 className="text-gray-500 text-sm font-medium">Total Hours</h3>
          <p className="text-3xl font-bold mt-2">{(overview?.total_hours || 0).toFixed(1)}</p>
        </div>
        <div className="bg-white p-6 rounded-lg shadow">
          <h3 className="text-gray-500 text-sm font-medium">Progress</h3>
          <p className="text-3xl font-bold mt-2">
            {overview?.total_hours ? ((overview.total_hours / 10000) * 100).toFixed(1) : 0}%
          </p>
        </div>
      </div>

      <div className="bg-white p-6 rounded-lg shadow">
        <h2 className="text-xl font-bold mb-4">Your Skills</h2>
        {overview?.skills && overview.skills.length > 0 ? (
          <div className="space-y-4">
            {overview.skills.map((skill: any) => (
              <div key={skill.skill_id} className="flex items-center justify-between">
                <div>
                  <h3 className="font-medium">{skill.skill_name}</h3>
                  <p className="text-sm text-gray-500">{skill.total_hours.toFixed(1)} hours</p>
                </div>
                <div className="w-48 bg-gray-200 rounded-full h-2">
                  <div
                    className="bg-primary-600 h-2 rounded-full"
                    style={{ width: `${Math.min((skill.total_hours / 10000) * 100, 100)}%` }}
                  ></div>
                </div>
              </div>
            ))}
          </div>
        ) : (
          <p className="text-gray-500">No skills yet. Create your first skill to get started!</p>
        )}
      </div>
    </div>
  )
}
