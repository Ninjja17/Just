import { useEffect, useState } from 'react'
import toast from 'react-hot-toast'
import { skillsAPI } from '@/services/api'
import type { Skill, CreateSkillRequest } from '@/types'
import Modal from '@/components/Modal'
import { PlusIcon, PencilSquareIcon, TrashIcon } from '@heroicons/react/24/outline'

const CATEGORY_OPTIONS = ['Music', 'Sports', 'Coding', 'Art', 'Language', 'Cooking', 'Other']
const COLOR_OPTIONS = ['#3B82F6', '#10B981', '#F59E0B', '#EF4444', '#8B5CF6', '#EC4899', '#14B8A6']
const ICON_OPTIONS = ['🎯', '🎵', '💻', '🎨', '🏀', '📚', '🍳', '✍️']

const emptyForm: CreateSkillRequest = {
  name: '',
  category: 'Other',
  color: COLOR_OPTIONS[0],
  icon: ICON_OPTIONS[0],
  target_hours: 10000,
}

export default function Skills() {
  const [skills, setSkills] = useState<Skill[]>([])
  const [loading, setLoading] = useState(true)
  const [modalOpen, setModalOpen] = useState(false)
  const [editing, setEditing] = useState<Skill | null>(null)
  const [form, setForm] = useState<CreateSkillRequest>(emptyForm)
  const [saving, setSaving] = useState(false)

  const load = async () => {
    setLoading(true)
    try {
      const res = await skillsAPI.getAll()
      setSkills(res.data || [])
    } catch (e: any) {
      toast.error(e?.response?.data?.error || 'Failed to load skills')
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    load()
  }, [])

  const openCreate = () => {
    setEditing(null)
    setForm(emptyForm)
    setModalOpen(true)
  }

  const openEdit = (s: Skill) => {
    setEditing(s)
    setForm({
      name: s.name,
      category: s.category,
      color: s.color,
      icon: s.icon,
      target_hours: s.target_hours,
    })
    setModalOpen(true)
  }

  const submit = async (e: React.FormEvent) => {
    e.preventDefault()
    if (!form.name.trim()) {
      toast.error('Name is required')
      return
    }
    setSaving(true)
    try {
      if (editing) {
        await skillsAPI.update(editing.id, form)
        toast.success('Skill updated')
      } else {
        await skillsAPI.create(form)
        toast.success('Skill created')
      }
      setModalOpen(false)
      load()
    } catch (e: any) {
      toast.error(e?.response?.data?.error || 'Save failed')
    } finally {
      setSaving(false)
    }
  }

  const remove = async (s: Skill) => {
    if (!confirm(`Delete skill "${s.name}"? This also removes its sessions.`)) return
    try {
      await skillsAPI.delete(s.id)
      toast.success('Skill deleted')
      setSkills((prev) => prev.filter((x) => x.id !== s.id))
    } catch (e: any) {
      toast.error(e?.response?.data?.error || 'Delete failed')
    }
  }

  return (
    <div>
      <div className="flex items-center justify-between mb-6">
        <h1 className="text-3xl font-bold">Skills</h1>
        <button
          onClick={openCreate}
          className="inline-flex items-center px-4 py-2 bg-primary-600 text-white rounded-md hover:bg-primary-700"
        >
          <PlusIcon className="w-5 h-5 mr-2" />
          New Skill
        </button>
      </div>

      {loading ? (
        <div className="text-gray-500">Loading skills...</div>
      ) : skills.length === 0 ? (
        <div className="bg-white p-8 rounded-lg shadow text-center text-gray-500">
          No skills yet. Click <strong>New Skill</strong> to create one.
        </div>
      ) : (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
          {skills.map((s) => (
            <div key={s.id} className="bg-white p-5 rounded-lg shadow flex flex-col">
              <div className="flex items-start justify-between">
                <div className="flex items-center">
                  <div
                    className="w-10 h-10 rounded-md flex items-center justify-center text-xl mr-3"
                    style={{ backgroundColor: s.color + '22', color: s.color }}
                  >
                    {s.icon}
                  </div>
                  <div>
                    <h3 className="font-semibold">{s.name}</h3>
                    <p className="text-xs text-gray-500">{s.category}</p>
                  </div>
                </div>
                <div className="flex space-x-1">
                  <button
                    onClick={() => openEdit(s)}
                    className="p-1 text-gray-500 hover:text-primary-600"
                    title="Edit"
                  >
                    <PencilSquareIcon className="w-5 h-5" />
                  </button>
                  <button
                    onClick={() => remove(s)}
                    className="p-1 text-gray-500 hover:text-red-600"
                    title="Delete"
                  >
                    <TrashIcon className="w-5 h-5" />
                  </button>
                </div>
              </div>
              <div className="mt-4 text-sm text-gray-600">
                Target: <strong>{s.target_hours.toLocaleString()}</strong> hours
              </div>
            </div>
          ))}
        </div>
      )}

      <Modal
        open={modalOpen}
        title={editing ? 'Edit skill' : 'New skill'}
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
              form="skill-form"
              disabled={saving}
              className="px-4 py-2 text-sm rounded-md bg-primary-600 text-white hover:bg-primary-700 disabled:opacity-50"
            >
              {saving ? 'Saving...' : editing ? 'Save changes' : 'Create skill'}
            </button>
          </div>
        }
      >
        <form id="skill-form" onSubmit={submit} className="space-y-4">
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">Name</label>
            <input
              required
              value={form.name}
              onChange={(e) => setForm({ ...form, name: e.target.value })}
              className="w-full px-3 py-2 border rounded-md"
              placeholder="e.g. Guitar"
            />
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">Category</label>
            <select
              value={form.category}
              onChange={(e) => setForm({ ...form, category: e.target.value })}
              className="w-full px-3 py-2 border rounded-md"
            >
              {CATEGORY_OPTIONS.map((c) => (
                <option key={c} value={c}>
                  {c}
                </option>
              ))}
            </select>
          </div>
          <div className="grid grid-cols-2 gap-4">
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">Color</label>
              <div className="flex flex-wrap gap-2">
                {COLOR_OPTIONS.map((c) => (
                  <button
                    type="button"
                    key={c}
                    onClick={() => setForm({ ...form, color: c })}
                    className={`w-7 h-7 rounded-full border-2 ${
                      form.color === c ? 'border-black' : 'border-transparent'
                    }`}
                    style={{ backgroundColor: c }}
                  />
                ))}
              </div>
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">Icon</label>
              <div className="flex flex-wrap gap-2">
                {ICON_OPTIONS.map((i) => (
                  <button
                    type="button"
                    key={i}
                    onClick={() => setForm({ ...form, icon: i })}
                    className={`w-8 h-8 rounded border text-lg ${
                      form.icon === i ? 'border-primary-600' : 'border-gray-200'
                    }`}
                  >
                    {i}
                  </button>
                ))}
              </div>
            </div>
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              Target hours
            </label>
            <input
              type="number"
              min={1}
              value={form.target_hours ?? 0}
              onChange={(e) =>
                setForm({ ...form, target_hours: parseInt(e.target.value || '0', 10) })
              }
              className="w-full px-3 py-2 border rounded-md"
            />
          </div>
        </form>
      </Modal>
    </div>
  )
}
