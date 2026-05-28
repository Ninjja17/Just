import { useEffect, useMemo, useState } from 'react'
import toast from 'react-hot-toast'
import { sessionsAPI, skillsAPI } from '@/services/api'
import type { Session, Skill, CreateSessionRequest } from '@/types'
import Modal from '@/components/Modal'
import {
  PlusIcon,
  PencilSquareIcon,
  TrashIcon,
  PlayIcon,
  StopIcon,
} from '@heroicons/react/24/outline'

const SESSION_TYPES = ['focused', 'casual', 'deliberate', 'performance']

function toLocalInput(d: Date) {
  const pad = (n: number) => String(n).padStart(2, '0')
  return `${d.getFullYear()}-${pad(d.getMonth() + 1)}-${pad(d.getDate())}T${pad(
    d.getHours()
  )}:${pad(d.getMinutes())}`
}

function fromLocalInput(s: string) {
  return new Date(s).toISOString()
}

export default function Sessions() {
  const [sessions, setSessions] = useState<Session[]>([])
  const [skills, setSkills] = useState<Skill[]>([])
  const [loading, setLoading] = useState(true)
  const [filterSkill, setFilterSkill] = useState<string>('')
  const [modalOpen, setModalOpen] = useState(false)
  const [editing, setEditing] = useState<Session | null>(null)
  const [saving, setSaving] = useState(false)

  const [form, setForm] = useState({
    skill_id: '',
    start_time: toLocalInput(new Date(Date.now() - 60 * 60 * 1000)),
    end_time: toLocalInput(new Date()),
    notes: '',
    session_type: 'focused',
    quality_rating: 3,
  })

  const skillMap = useMemo(() => {
    const m: Record<string, Skill> = {}
    skills.forEach((s) => (m[s.id] = s))
    return m
  }, [skills])

  const load = async () => {
    setLoading(true)
    try {
      const [sRes, kRes] = await Promise.all([
        sessionsAPI.getAll(filterSkill ? { skill_id: filterSkill } : undefined),
        skillsAPI.getAll(),
      ])
      setSessions(sRes.data || [])
      setSkills(kRes.data || [])
    } catch (e: any) {
      toast.error(e?.response?.data?.error || 'Failed to load sessions')
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    load()
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [filterSkill])

  const openCreate = () => {
    if (skills.length === 0) {
      toast.error('Create a skill first')
      return
    }
    setEditing(null)
    setForm({
      skill_id: skills[0].id,
      start_time: toLocalInput(new Date(Date.now() - 60 * 60 * 1000)),
      end_time: toLocalInput(new Date()),
      notes: '',
      session_type: 'focused',
      quality_rating: 3,
    })
    setModalOpen(true)
  }

  const openEdit = (s: Session) => {
    setEditing(s)
    setForm({
      skill_id: s.skill_id,
      start_time: toLocalInput(new Date(s.start_time)),
      end_time: s.end_time ? toLocalInput(new Date(s.end_time)) : toLocalInput(new Date()),
      notes: s.notes || '',
      session_type: s.session_type || 'focused',
      quality_rating: s.quality_rating ?? 3,
    })
    setModalOpen(true)
  }

  const submit = async (e: React.FormEvent) => {
    e.preventDefault()
    if (!form.skill_id) {
      toast.error('Pick a skill')
      return
    }
    const start = new Date(form.start_time)
    const end = new Date(form.end_time)
    if (end <= start) {
      toast.error('End time must be after start time')
      return
    }
    const duration = Math.max(1, Math.round((end.getTime() - start.getTime()) / 60000))

    const payload: CreateSessionRequest = {
      skill_id: form.skill_id,
      start_time: fromLocalInput(form.start_time),
      end_time: fromLocalInput(form.end_time),
      duration_minutes: duration,
      notes: form.notes || undefined,
      session_type: form.session_type,
      quality_rating: form.quality_rating,
    }

    setSaving(true)
    try {
      if (editing) {
        await sessionsAPI.update(editing.id, payload)
        toast.success('Session updated')
      } else {
        await sessionsAPI.create(payload)
        toast.success('Session logged')
      }
      setModalOpen(false)
      load()
    } catch (e: any) {
      toast.error(e?.response?.data?.error || 'Save failed')
    } finally {
      setSaving(false)
    }
  }

  const remove = async (s: Session) => {
    if (!confirm('Delete this session?')) return
    try {
      await sessionsAPI.delete(s.id)
      toast.success('Session deleted')
      setSessions((prev) => prev.filter((x) => x.id !== s.id))
    } catch (e: any) {
      toast.error(e?.response?.data?.error || 'Delete failed')
    }
  }

  const startTimer = async (skill_id: string) => {
    try {
      await sessionsAPI.startTimer(skill_id)
      toast.success('Timer started')
      load()
    } catch (e: any) {
      toast.error(e?.response?.data?.error || 'Could not start timer')
    }
  }

  const stopTimer = async (session_id: string) => {
    try {
      await sessionsAPI.stopTimer(session_id)
      toast.success('Timer stopped')
      load()
    } catch (e: any) {
      toast.error(e?.response?.data?.error || 'Could not stop timer')
    }
  }

  const running = sessions.filter((s) => !s.end_time)

  return (
    <div>
      <div className="flex items-center justify-between mb-6">
        <h1 className="text-3xl font-bold">Practice Sessions</h1>
        <div className="flex items-center space-x-2">
          <select
            value={filterSkill}
            onChange={(e) => setFilterSkill(e.target.value)}
            className="px-3 py-2 border rounded-md text-sm"
          >
            <option value="">All skills</option>
            {skills.map((s) => (
              <option key={s.id} value={s.id}>
                {s.name}
              </option>
            ))}
          </select>
          <button
            onClick={openCreate}
            className="inline-flex items-center px-4 py-2 bg-primary-600 text-white rounded-md hover:bg-primary-700"
          >
            <PlusIcon className="w-5 h-5 mr-2" />
            Log session
          </button>
        </div>
      </div>

      {/* Quick timer */}
      <div className="bg-white rounded-lg shadow p-5 mb-6">
        <h2 className="font-semibold mb-3">Quick timer</h2>
        {skills.length === 0 ? (
          <p className="text-sm text-gray-500">Create a skill to start a timer.</p>
        ) : (
          <div className="flex flex-wrap gap-2">
            {skills.map((s) => {
              const r = running.find((x) => x.skill_id === s.id)
              return r ? (
                <button
                  key={s.id}
                  onClick={() => stopTimer(r.id)}
                  className="inline-flex items-center px-3 py-2 rounded-md bg-red-600 text-white text-sm hover:bg-red-700"
                >
                  <StopIcon className="w-4 h-4 mr-1" /> Stop {s.name}
                </button>
              ) : (
                <button
                  key={s.id}
                  onClick={() => startTimer(s.id)}
                  className="inline-flex items-center px-3 py-2 rounded-md bg-emerald-600 text-white text-sm hover:bg-emerald-700"
                >
                  <PlayIcon className="w-4 h-4 mr-1" /> Start {s.name}
                </button>
              )
            })}
          </div>
        )}
      </div>

      {loading ? (
        <div className="text-gray-500">Loading sessions...</div>
      ) : sessions.length === 0 ? (
        <div className="bg-white p-8 rounded-lg shadow text-center text-gray-500">
          No sessions yet.
        </div>
      ) : (
        <div className="bg-white rounded-lg shadow overflow-hidden">
          <table className="min-w-full text-sm">
            <thead className="bg-gray-50 text-left text-gray-600">
              <tr>
                <th className="px-4 py-2">Skill</th>
                <th className="px-4 py-2">Start</th>
                <th className="px-4 py-2">End</th>
                <th className="px-4 py-2">Duration</th>
                <th className="px-4 py-2">Type</th>
                <th className="px-4 py-2">Quality</th>
                <th className="px-4 py-2">Notes</th>
                <th className="px-4 py-2"></th>
              </tr>
            </thead>
            <tbody>
              {sessions.map((s) => (
                <tr key={s.id} className="border-t">
                  <td className="px-4 py-2">{skillMap[s.skill_id]?.name || s.skill_id.slice(0, 8)}</td>
                  <td className="px-4 py-2">{new Date(s.start_time).toLocaleString()}</td>
                  <td className="px-4 py-2">
                    {s.end_time ? new Date(s.end_time).toLocaleString() : <span className="text-emerald-600">running</span>}
                  </td>
                  <td className="px-4 py-2">{s.duration_minutes} min</td>
                  <td className="px-4 py-2 capitalize">{s.session_type}</td>
                  <td className="px-4 py-2">{s.quality_rating ?? '-'}</td>
                  <td className="px-4 py-2 max-w-[240px] truncate" title={s.notes || ''}>
                    {s.notes || '-'}
                  </td>
                  <td className="px-4 py-2 text-right whitespace-nowrap">
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
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}

      <Modal
        open={modalOpen}
        title={editing ? 'Edit session' : 'Log session'}
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
              form="session-form"
              disabled={saving}
              className="px-4 py-2 text-sm rounded-md bg-primary-600 text-white hover:bg-primary-700 disabled:opacity-50"
            >
              {saving ? 'Saving...' : editing ? 'Save changes' : 'Create session'}
            </button>
          </div>
        }
      >
        <form id="session-form" onSubmit={submit} className="space-y-4">
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">Skill</label>
            <select
              required
              value={form.skill_id}
              onChange={(e) => setForm({ ...form, skill_id: e.target.value })}
              className="w-full px-3 py-2 border rounded-md"
            >
              <option value="">Select a skill</option>
              {skills.map((s) => (
                <option key={s.id} value={s.id}>
                  {s.name}
                </option>
              ))}
            </select>
          </div>
          <div className="grid grid-cols-2 gap-4">
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">Start</label>
              <input
                required
                type="datetime-local"
                value={form.start_time}
                onChange={(e) => setForm({ ...form, start_time: e.target.value })}
                className="w-full px-3 py-2 border rounded-md"
              />
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">End</label>
              <input
                required
                type="datetime-local"
                value={form.end_time}
                onChange={(e) => setForm({ ...form, end_time: e.target.value })}
                className="w-full px-3 py-2 border rounded-md"
              />
            </div>
          </div>
          <div className="grid grid-cols-2 gap-4">
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">Type</label>
              <select
                value={form.session_type}
                onChange={(e) => setForm({ ...form, session_type: e.target.value })}
                className="w-full px-3 py-2 border rounded-md capitalize"
              >
                {SESSION_TYPES.map((t) => (
                  <option key={t} value={t} className="capitalize">
                    {t}
                  </option>
                ))}
              </select>
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                Quality (1-5)
              </label>
              <input
                type="number"
                min={1}
                max={5}
                value={form.quality_rating}
                onChange={(e) =>
                  setForm({ ...form, quality_rating: parseInt(e.target.value || '3', 10) })
                }
                className="w-full px-3 py-2 border rounded-md"
              />
            </div>
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">Notes</label>
            <textarea
              value={form.notes}
              onChange={(e) => setForm({ ...form, notes: e.target.value })}
              className="w-full px-3 py-2 border rounded-md min-h-[80px]"
              placeholder="Optional notes"
            />
          </div>
        </form>
      </Modal>
    </div>
  )
}
