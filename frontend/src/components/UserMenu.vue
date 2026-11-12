<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { RouterLink, useRouter } from 'vue-router'
import { LogOut, Settings as SettingsIcon, User, Building2, Briefcase, FileText, Shield, Bookmark, Search } from 'lucide-vue-next'
import { useAuthStore } from '../stores/auth'

const auth = useAuthStore()
const router = useRouter()

const open = ref(false)
const rootEl = ref(null)

function toggle() { open.value = !open.value }
function close() { open.value = false }

function handleLogout() {
  auth.logout()
  close()
  router.push({ name: 'login' })
}

function onDocumentClick(e) {
  if (rootEl.value && !rootEl.value.contains(e.target)) close()
}

onMounted(() => document.addEventListener('click', onDocumentClick))
onUnmounted(() => document.removeEventListener('click', onDocumentClick))

const initials = computed(() => (auth.user?.username || '?').charAt(0).toUpperCase())

const roleLabel = computed(() => {
  const r = auth.user?.role
  if (r === 'job_seeker') return 'Job seeker'
  if (r === 'employer') return 'Employer'
  if (r === 'admin') return 'Admin'
  return r || ''
})
</script>

<template>
  <div ref="rootEl" class="relative">
    <button type="button" @click="toggle" class="flex items-center gap-2 px-2 py-1 hover:bg-surface-muted transition-colors">
      <span class="w-8 h-8 bg-primary text-white text-sm font-bold flex items-center justify-center">
        {{ initials }}
      </span>
      <span class="text-sm font-medium text-fg">{{ auth.user?.username }}</span>
    </button>

    <Transition enter-active-class="transition duration-100 ease-out" enter-from-class="opacity-0 scale-95" leave-active-class="transition duration-75 ease-in" leave-to-class="opacity-0 scale-95">
      <div v-if="open" class="absolute right-0 mt-2 w-56 bg-surface border border-border shadow-2xl z-40 py-1 origin-top-right">
        <div class="px-4 py-3 border-b border-border">
          <p class="text-sm font-bold text-fg truncate">{{ auth.user?.username }}</p>
          <p class="text-xs text-fg-muted">{{ roleLabel }}</p>
        </div>

        <RouterLink v-if="auth.user?.role === 'job_seeker' && auth.user?.onboardingCompleted" to="/profile" @click="close" class="flex items-center gap-3 px-4 py-2.5 text-sm text-fg-muted hover:bg-surface-muted hover:text-fg transition-colors">
          <User class="w-4 h-4" /> Profile
        </RouterLink>
        <RouterLink v-if="auth.user?.role === 'job_seeker' && auth.user?.onboardingCompleted" to="/applications" @click="close" class="flex items-center gap-3 px-4 py-2.5 text-sm text-fg-muted hover:bg-surface-muted hover:text-fg transition-colors">
          <FileText class="w-4 h-4" /> My applications
        </RouterLink>
        <RouterLink v-if="auth.user?.role === 'job_seeker' && auth.user?.onboardingCompleted" to="/saved" @click="close" class="flex items-center gap-3 px-4 py-2.5 text-sm text-fg-muted hover:bg-surface-muted hover:text-fg transition-colors">
          <Bookmark class="w-4 h-4" /> Saved jobs
        </RouterLink>
        <RouterLink v-if="auth.user?.role === 'job_seeker' && auth.user?.onboardingCompleted" to="/searches" @click="close" class="flex items-center gap-3 px-4 py-2.5 text-sm text-fg-muted hover:bg-surface-muted hover:text-fg transition-colors">
          <Search class="w-4 h-4" /> Saved searches
        </RouterLink>

        <RouterLink v-if="auth.user?.role === 'employer' && auth.user?.onboardingCompleted" to="/company" @click="close" class="flex items-center gap-3 px-4 py-2.5 text-sm text-fg-muted hover:bg-surface-muted hover:text-fg transition-colors">
          <Building2 class="w-4 h-4" /> Company
        </RouterLink>
        <RouterLink v-if="auth.user?.role === 'employer' && auth.user?.onboardingCompleted" to="/jobs" @click="close" class="flex items-center gap-3 px-4 py-2.5 text-sm text-fg-muted hover:bg-surface-muted hover:text-fg transition-colors">
          <Briefcase class="w-4 h-4" /> Jobs
        </RouterLink>
        <RouterLink v-if="auth.user?.role === 'employer' && auth.user?.onboardingCompleted" to="/employer/applications" @click="close" class="flex items-center gap-3 px-4 py-2.5 text-sm text-fg-muted hover:bg-surface-muted hover:text-fg transition-colors">
          <FileText class="w-4 h-4" /> Applications
        </RouterLink>

        <RouterLink v-if="auth.user?.role === 'admin'" to="/admin" @click="close" class="flex items-center gap-3 px-4 py-2.5 text-sm text-fg-muted hover:bg-surface-muted hover:text-fg transition-colors">
          <Shield class="w-4 h-4" /> Admin
        </RouterLink>

        <RouterLink to="/settings" @click="close" class="flex items-center gap-3 px-4 py-2.5 text-sm text-fg-muted hover:bg-surface-muted hover:text-fg transition-colors border-t border-border mt-1">
          <SettingsIcon class="w-4 h-4" /> Settings
        </RouterLink>

        <button type="button" @click="handleLogout" class="w-full flex items-center gap-3 px-4 py-2.5 text-sm text-red-500 hover:bg-surface-muted hover:text-red-400 transition-colors text-left border-t border-border">
          <LogOut class="w-4 h-4" /> Logout
        </button>
      </div>
    </Transition>
  </div>
</template>
