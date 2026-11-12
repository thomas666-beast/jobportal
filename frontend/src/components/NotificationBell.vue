<script setup>
import { ref, onMounted, onUnmounted, computed } from 'vue'
import { RouterLink, useRouter } from 'vue-router'
import { Bell, Check, FileText, Briefcase, XCircle } from 'lucide-vue-next'
import { useNotificationsStore } from '../stores/notifications'
import { timeAgo } from '../utils/jobStatus'

const store = useNotificationsStore()
const router = useRouter()

const open = ref(false)
const rootEl = ref(null)

function toggle() {
  open.value = !open.value
  if (open.value) store.refresh()
}

function close() {
  open.value = false
}

function iconFor(type) {
  if (type === 'application_submitted') return FileText
  if (type === 'application_status_changed') return Briefcase
  if (type === 'job_rejected') return XCircle
  if (type === 'job_approved') return Check
  return Bell
}

function clickNotification(n) {
  store.markRead(n.id)
  if (n.link) {
    router.push(n.link)
  }
  close()
}

function onDocumentClick(e) {
  if (rootEl.value && !rootEl.value.contains(e.target)) {
    close()
  }
}

onMounted(() => {
  document.addEventListener('click', onDocumentClick)
  store.startPolling()
})

onUnmounted(() => {
  document.removeEventListener('click', onDocumentClick)
  store.stopPolling()
})

const hasUnread = computed(() => store.unread > 0)
</script>

<template>
  <div ref="rootEl" class="relative">
    <button
      type="button"
      @click="toggle"
      class="relative p-2 text-slate-600 hover:text-slate-900 hover:bg-slate-100 rounded-full transition-colors"
      aria-label="Notifications"
    >
      <Bell class="w-5 h-5" />
      <span
        v-if="hasUnread"
        class="absolute top-1 right-1 min-w-[16px] h-4 px-1 rounded-full bg-red-500 text-white
               text-[10px] font-semibold flex items-center justify-center"
      >
        {{ store.unread > 99 ? '99+' : store.unread }}
      </span>
    </button>

    <Transition
      enter-active-class="transition duration-100 ease-out"
      enter-from-class="opacity-0 scale-95"
      leave-active-class="transition duration-75 ease-in"
      leave-to-class="opacity-0 scale-95"
    >
      <div
        v-if="open"
        class="absolute right-0 mt-1 w-80 bg-white rounded-lg border border-slate-200
               shadow-lg z-40 origin-top-right overflow-hidden"
      >
        <div class="px-3 py-2 border-b border-slate-100 flex items-center justify-between">
          <p class="text-sm font-medium text-slate-800">Notifications</p>
          <button
            v-if="hasUnread"
            type="button"
            @click="store.markAllRead()"
            class="text-xs text-primary-600 hover:text-primary-700"
          >
            Mark all read
          </button>
        </div>

        <div class="max-h-96 overflow-y-auto">
          <p v-if="store.recent.length === 0" class="px-3 py-6 text-center text-sm text-slate-500">
            No notifications yet
          </p>

          <button
            v-for="n in store.recent"
            :key="n.id"
            type="button"
            @click="clickNotification(n)"
            class="w-full text-left px-3 py-3 hover:bg-slate-50 border-b border-slate-100 last:border-b-0
                   flex gap-3 transition-colors"
            :class="!n.read_at && 'bg-primary-50/40'"
          >
            <component
              :is="iconFor(n.type)"
              class="w-4 h-4 mt-0.5 text-slate-400 flex-shrink-0"
            />
            <div class="min-w-0 flex-1">
              <p class="text-sm font-medium text-slate-800 truncate">
                {{ n.title }}
              </p>
              <p v-if="n.body" class="text-xs text-slate-600 mt-0.5 line-clamp-2">
                {{ n.body }}
              </p>
              <p class="text-[10px] text-slate-400 mt-1">
                {{ timeAgo(n.created_at) }}
              </p>
            </div>
            <span v-if="!n.read_at" class="w-2 h-2 rounded-full bg-primary-500 mt-1.5 flex-shrink-0" />
          </button>
        </div>

        <RouterLink
          to="/notifications"
          @click="close"
          class="block px-3 py-2 text-center text-xs text-primary-600 hover:bg-slate-50 border-t border-slate-100"
        >
          View all notifications
        </RouterLink>
      </div>
    </Transition>
  </div>
</template>
