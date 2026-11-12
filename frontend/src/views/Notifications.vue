<script setup>
import { ref, onMounted } from 'vue'
import { RouterLink, useRouter } from 'vue-router'
import { Bell, Check, FileText, Briefcase, XCircle } from 'lucide-vue-next'
import { listNotifications, markRead, markAllRead } from '../api/notifications'
import { useNotificationsStore } from '../stores/notifications'
import { timeAgo } from '../utils/jobStatus'
import AppButton from '../components/ui/AppButton.vue'
import AppCard from '../components/ui/AppCard.vue'
import EmptyState from '../components/ui/EmptyState.vue'
import AppSkeleton from '../components/ui/AppSkeleton.vue'
import { useToast } from '../composables/useToast'

const router = useRouter()
const store = useNotificationsStore()
const toast = useToast()

const notifications = ref([])
const unread = ref(0)
const loading = ref(true)
const errorMsg = ref('')

async function load() {
  loading.value = true
  errorMsg.value = ''
  try {
    const res = await listNotifications({ limit: 50 })
    notifications.value = res.notifications || []
    unread.value = res.unread || 0
    // Sync the store's recent dropdown counter.
    store.fetchUnreadCount()
  } catch (err) {
    errorMsg.value = err.message
  } finally {
    loading.value = false
  }
}

function iconFor(type) {
  if (type === 'application_submitted') return FileText
  if (type === 'application_status_changed') return Briefcase
  if (type === 'job_rejected') return XCircle
  if (type === 'job_approved') return Check
  return Bell
}

async function clickNotification(n) {
  if (!n.read_at) {
    try {
      await markRead(n.id)
      n.read_at = new Date().toISOString()
      if (unread.value > 0) unread.value--
      store.fetchUnreadCount()
    } catch {
      // silent
    }
  }
  if (n.link) {
    router.push(n.link)
  }
}

async function handleMarkAll() {
  try {
    await markAllRead()
    notifications.value.forEach((n) => {
      if (!n.read_at) n.read_at = new Date().toISOString()
    })
    unread.value = 0
    store.fetchUnreadCount()
    toast.success('All notifications marked as read.')
  } catch (err) {
    toast.error(err.message)
  }
}

onMounted(load)
</script>

<template>
  <div class="max-w-2xl mx-auto space-y-4">
    <div class="flex items-center justify-between">
      <div>
        <h1 class="text-2xl font-semibold text-slate-900">Notifications</h1>
        <p class="text-sm text-slate-600 mt-1">
          {{ unread }} unread
        </p>
      </div>
      <AppButton
        v-if="unread > 0"
        variant="secondary"
        size="sm"
        @click="handleMarkAll"
      >
        Mark all read
      </AppButton>
    </div>

    <div v-if="loading" class="space-y-2">
      <AppSkeleton v-for="i in 5" :key="i" variant="row" />
    </div>
    <p v-else-if="errorMsg" class="text-red-600">{{ errorMsg }}</p>

    <AppCard v-else-if="notifications.length === 0" :padded="false">
      <EmptyState
        :icon="Bell"
        title="No notifications"
        description="You'll see updates here when something happens."
      />
    </AppCard>

    <div v-else class="space-y-2">
      <AppCard
        v-for="n in notifications"
        :key="n.id"
        :padded="false"
        :class="!n.read_at && 'border-primary-200 bg-primary-50/30'"
      >
        <button
          type="button"
          @click="clickNotification(n)"
          class="w-full text-left p-4 flex gap-3"
        >
          <component
            :is="iconFor(n.type)"
            class="w-5 h-5 mt-0.5 text-slate-400 flex-shrink-0"
          />
          <div class="min-w-0 flex-1">
            <div class="flex items-start justify-between gap-2">
              <p class="text-sm font-medium text-slate-800">{{ n.title }}</p>
              <span
                v-if="!n.read_at"
                class="w-2 h-2 rounded-full bg-primary-500 flex-shrink-0 mt-1.5"
              />
            </div>
            <p v-if="n.body" class="text-sm text-slate-600 mt-0.5">
              {{ n.body }}
            </p>
            <p class="text-xs text-slate-400 mt-1">{{ timeAgo(n.created_at) }}</p>
          </div>
        </button>
      </AppCard>
    </div>
  </div>
</template>
