import { defineStore } from 'pinia'
import { ref } from 'vue'
import * as api from '../api/notifications'

export const useNotificationsStore = defineStore('notifications', () => {
  const unread = ref(0)
  const recent = ref([])
  const loading = ref(false)

  let pollTimer = null

  async function refresh() {
    loading.value = true
    try {
      const res = await api.listRecentNotifications()
      recent.value = res.notifications || []
      unread.value = res.unread || 0
    } catch {
      // silent — notifications failing should never block the UI
    } finally {
      loading.value = false
    }
  }

  async function fetchUnreadCount() {
    try {
      const res = await api.getUnreadCount()
      unread.value = res.unread || 0
    } catch {
      // silent
    }
  }

  async function markRead(id) {
    try {
      await api.markRead(id)
      const n = recent.value.find((x) => x.id === id)
      if (n && !n.read_at) {
        n.read_at = new Date().toISOString()
        if (unread.value > 0) unread.value--
      }
    } catch {
      // silent
    }
  }

  async function markAllRead() {
    try {
      await api.markAllRead()
      recent.value.forEach((n) => {
        if (!n.read_at) n.read_at = new Date().toISOString()
      })
      unread.value = 0
    } catch {
      // silent
    }
  }

  function startPolling(intervalMs = 60000) {
    stopPolling()
    fetchUnreadCount()
    pollTimer = setInterval(fetchUnreadCount, intervalMs)
  }

  function stopPolling() {
    if (pollTimer) {
      clearInterval(pollTimer)
      pollTimer = null
    }
  }

  function reset() {
    unread.value = 0
    recent.value = []
    stopPolling()
  }

  return {
    unread,
    recent,
    loading,
    refresh,
    fetchUnreadCount,
    markRead,
    markAllRead,
    startPolling,
    stopPolling,
    reset
  }
})
