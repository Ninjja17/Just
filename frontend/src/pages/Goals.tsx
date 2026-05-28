import { useEffect, useMemo, useState } from 'react'
import toast from 'react-hot-toast'
import { goalsAPI, skillsAPI } from '@/services/api'
import type { Goal, Skill, CreateGoalRequest } from '@/types'
import Modal from '@/components/Modal'
import { PlusIcon, PencilSquareIcon, TrashIcon, CheckCircleIcon } from '@heroicons/react/24/outline'

const TYPES: Array<Goal['type']> = ['daily', 'weekly', 'milestone']

export default function Goals() {
  const [goals, setGoals] = useState<Goal[]>([])
  const [skills, setSkills] = useState<Skill[]>([])
  const [loading, setLoading] = useState(true)
  const [modalOpen, setModalOpen] = useState(false)
  const [editing, setEditing] = useState<Goal | null>(null)
  const [saving, setSaving] = useState(false)
  const [form, setForm] = useState<{
    skill_id: string
    type: Goal['type']
    target_hours: number
    deadline: string
  }>({
    skill_id: '',
    type: 'daily',
    target_hours: 1,
    deadline: '',
  })

  const skillMap = useMemo(() => {
    const m: Record<string, Skill> = {}
    skills.forEach((s) => (m[s.id] = s))
    return m
  }, [skills])

  const load = async () => {
    setLoading(true)
    try {
      const [g, s] = await Promise.all([goalsAPI.getAll(), skillsAPI.getAll()])
      setGoals(g.data || [])
      setSkills(s.data || [])
    } catch (e: any) {
      toast.error(e?.response?.data?.error || 'Failed to load goals')
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    load()
  }, [])

  const openCreate = () => {
    setEditing(null)
    setForm({ skill_id: '', type: 'daily', target_hours: 1, deadline: '' })
    setModalOpen(true)
  }

  const openEdit = (g: Goal) => {
    setEditing(g)
    setForm({
      skill_id: g.skill_id || '',
      type: g.type,
      target_hours: g.target_hours,
      deadline: g.deadline ? g.deadline.slice(0, 10) : '',
    })
    setModalOpen(true)
  }

  const submit = async (e: React.FormEvent) => {
    e.preventDefault()
    if (!form.target_hours || form.target_hours <= 0) {
      toast.error('Target hours must be greater than 0')
      return
    }
    const payload: CreateGoalRequest = {
      type: form.type,
      target_hours: form.target_hours,
      skill_id: form.skill_id || undefined,
      deadline: form.deadline ? new Date(form.deadline).toISOString() : undefined,
    }
    setSaving(true)
    try {
      if (editing) {
        await goalsAPI.update(editing.id, payload)
        toast.success('Goal updated')
      } else {
        await goalsAPI.create(payload)
        toast.success('Goal created')
      }
      setModalOpen(false)
      load()
    } catch (e: any) {
      toast.error(e?.response?.data?.error || 'Save failed')
    } finally {
      setSaving(false)
    }
  }

  const remove = async (g: Goal) => {
    if (!confirm('Delete this goal?')) return
    try {
      await goalsAPI.delete(g.id)
      toast.success('Goal deleted')
      setGoals((prev) => prev.filter((x) => x.id !== g.id))
    } catch (e: any) {
      toast.error(e?.response?.data?.error || 'Delete failed')
    }
  }

  return (
    <div>
      <div className="flex items-center justify-between mb-6">
        <h1 className="text-3xl font-bold">Goals</h1>
        <button
          onClick={openCreate}
          className="inline-flex items-center px-4 py-2 bg-primary-600 text-white rounded-md hover:bg-primary-700"
        >
          <PlusIcon className="w-5 h-5 mr-2" />
          New Goal
        </button>
      </div>

      {loading ? (
        <div className="text-gray-500">Loading goals...</div>
      ) : goals.length === 0 ? (
        <div className="bg-white p-8 rounded-lg shadow text-center text-gray-500">
          No goals yet.
        </div>
      ) : (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
          {goals.map((g) => (
            <div key={g.id} className="bg-white p-5 rounded-lg shadow">
              <div className="flex justify-between items-start">
                <div>
                  <p className="text-xs uppercase tracking-wide text-gray-500">
                    {g.type}
                  </p>
                  <h3 className="font-semibold text-lg">
                    {g.target_hours} hour{g.target_hours === 1 ? '' : 's'}
                  </h3>
                  <p className="text-sm text-gray-600 mt-1">
                    {g.skill_id ? skillMap[g.skill_id]?.name || 'Skill' : 'All skills'}
                  </p>
                  {g.deadline && (
                    <p className="text-xs text-gray-500 mt-1">
                      Due {new Date(g.deadline).toLocaleDateString()}
                    </p>
                  )}
                </div>
                <div className="flex flex-col items-end">
                  {g.is_completed && (
                    <CheckCircleIcon className="w-6 h-6 text-emerald-500" />
                  )}
                  <div className="flex space-x-1 mt-2">
                    <button
                      onClick={() => openEdit(g)}
                      className="p-1 text-gray-500 hover:text-primary-600"
                      title="Edit"
                    >
                      <PencilSquareIcon className="w-5 h-5" />
                    </button>
                    <button
                      onClick={() => remove(g)}
                      className="p-1 text-gray-500 hover:text-red-600"
                      title="Delete"
                    >
                      <TrashIcon className="w-5 h-5" />
                    </button>
                  </div>
                </div>
              </div>
            </div>
          ))}
        </div>
      )}

      <Modal
        open={modalOpen}
        title={editing ? 'Edit goal' : 'New goal'}
        onClose={() => setModalOpen(false)}
        footer={
          <div className="flex justify-end space-x-2">
            <button
              type="button"
              onClick={() => setModalOpen(false)}
              className="px-4 py-2 text-sm rounded-md border"
            >
              Cancel
            </button>
            <button
              type="submit"
              form="goal-form"
              disabled={saving}
              className="px-4 py-2 text-sm rounded-md bg-primary-600 text-white hover:bg-primary-700 disabled:opacity-50"
            >
              {saving ? 'Saving...' : editing ? 'Save changes' : 'Create goal'}
            </button>
          </div>
        }
      >
        <form id="goal-form" onSubmit={submit} className="space-y-4">
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">Type</label>
            <select
              value={form.type}
              onChange={(e) => setForm({ ...form, type: e.target.value as Goal['type'] })}
              className="w-full px-3 py-2 border rounded-md capitalize"
            >
              {TYPES.map((t) => (
                <option key={t} value={t} className="capitalize">
                  {t}
                </option>
              ))}
            </select>
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              Skill (optional)
            </label>
            <select
              value={form.skill_id}
              onChange={(e) => setForm({ ...form, skill_id: e.target.value })}
              className="w-full px-3 py-2 border rounded-md"
            >
              <option value="">All skills</option>
              {skills.map((s) => (
                <option key={s.id} value={s.id}>
                  {s.name}
                </option>
              ))}
            </select>
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              Target hours
            </label>
            <input
              type="number"
              min={0.1}
              step={0.1}
              value={form.target_hours}
              onChange={(e) =>
                setForm({ ...form, target_hours: parseFloat(e.target.value || '0') })
              }
              className="w-full px-3 py-2 border rounded-md"
            />
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              Deadline (optional)
            </label>
            <input
              type="date"
              value={form.deadline}
              onChange={(e) => setForm({ ...form, deadline: e.target.value })}
              className="w-full px-3 py-2 border rounded-md"
            />
          </div>
        </form>
      </Modal>
    </div>
  )
}
