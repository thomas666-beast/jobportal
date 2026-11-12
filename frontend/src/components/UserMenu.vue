<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { RouterLink, useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'
import { LogOut, Settings as SettingsIcon, User, Building2, Briefcase, FileText, Shield, Bookmark, Search } from 'lucide-vue-next'

const auth = useAuthStore()
const router = useRouter()

const open = ref(false)
const rootEl = ref(null)

function toggle() {
  open.value = !open.value
}

function close() {
  open.value = false
}

function handleLogout() {
  auth.logout()
  close()
  router.push({ name: 'login' })
}

function onDocumentClick(e) {
  if (rootEl.value && !rootEl.value.contains(e.target)) {
    close()
  }
}

onMounted(() => document.addEventListener('click', onDocumentClick))
onUnmounted(() => document.removeEventListener('click', onDocumentClick))

const initials = computed(() => {
  const name = auth.user?.username || '?'
  return name.charAt(0).toUpperCase()
})

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
    <button
      type="button"
      @click="toggle"
      class="flex items-center gap-2 rounded-full pl-1 pr-3 py-1
             hover:bg-slate-100 transition-colors"
    >
      <span
        class="w-7 h-7 rounded-full bg-primary-100 text-primary-700 text-xs font-semibold
               flex items-center justify-center"
      >
        {{ initials }}
      </span>
      <span class="text-sm text-slate-700">{{ auth.user?.username }}</span>
    </button>

    <Transition
      enter-active-class="transition duration-100 ease-out"
      enter-from-class="opacity-0 scale-95"
      leave-active-class="transition duration-75 ease-in"
      leave-to-class="opacity-0 scale-95"
    >
      <div
        v-if="open"
        class="absolute right-0 mt-1 w-56 bg-white rounded-lg border border-slate-200
               shadow-lg z-40 py-1 origin-top-right"
      >
        <div class="px-3 py-2 border-b border-slate-100">
          <p class="text-sm font-medium text-slate-800">{{ auth.user?.username }}</p>
          <p class="text-xs text-slate-500">{{ roleLabel }}</p>
        </div>

        <RouterLink
          v-if="auth.user?.role === 'job_seeker' && auth.user?.onboardingCompleted"
          to="/profile"
          @click="close"
          class="flex items-center gap-2 px-3 py-2 text-sm text-slate-700 hover:bg-slate-50"
        >
          <User class="w-4 h-4" /> Profile
        </RouterLink>
        <RouterLink
          v-if="auth.user?.role === 'job_seeker' && auth.user?.onboardingCompleted"
          to="/applications"
          @click="close"
          class="flex items-center gap-2 px-3 py-2 text-sm text-slate-700 hover:bg-slate-50"
        >
          <FileText class="w-4 h-4" /> My applications
        </RouterLink>

        <RouterLink
          v-if="auth.user?.role === 'job_seeker' && auth.user?.onboardingCompleted"
          to="/saved"
          @click="close"
          class="flex items-center gap-2 px-3 py-2 text-sm text-slate-700 hover:bg-slate-50"
        >
          <Bookmark class="w-4 h-4" /> Saved jobs
        </RouterLink>

        <RouterLink
          v-if="auth.user?.role === 'job_seeker' && auth.user?.onboardingCompleted"
          to="/searches"
          @click="close"
          class="flex items-center gap-2 px-3 py-2 text-sm text-slate-700 hover:bg-slate-50"
        >
          <Search class="w-4 h-4" /> Saved searches
        </RouterLink>

        <RouterLink
          v-if="auth.user?.role === 'employer' && auth.user?.onboardingCompleted"
          to="/company"
          @click="close"
          class="flex items-center gap-2 px-3 py-2 text-sm text-slate-700 hover:bg-slate-50"
        >
          <Building2 class="w-4 h-4" /> Company
        </RouterLink>
        <RouterLink
          v-if="auth.user?.role === 'employer' && auth.user?.onboardingCompleted"
          to="/jobs"
          @click="close"
          class="flex items-center gap-2 px-3 py-2 text-sm text-slate-700 hover:bg-slate-50"
        >
          <Briefcase class="w-4 h-4" /> Jobs
        </RouterLink>
        <RouterLink
          v-if="auth.user?.role === 'employer' && auth.user?.onboardingCompleted"
          to="/employer/applications"
          @click="close"
          class="flex items-center gap-2 px-3 py-2 text-sm text-slate-700 hover:bg-slate-50"
        >
          <FileText class="w-4 h-4" /> Applications
        </RouterLink>

        <RouterLink
          v-if="auth.user?.role === 'admin'"
          to="/admin"
          @click="close"
          class="flex items-center gap-2 px-3 py-2 text-sm text-slate-700 hover:bg-slate-50"
        >
          <Shield class="w-4 h-4" /> Admin
        </RouterLink>

        <RouterLink
          to="/settings"
          @click="close"
          class="flex items-center gap-2 px-3 py-2 text-sm text-slate-700 hover:bg-slate-50"
        >
          <SettingsIcon class="w-4 h-4" /> Settings
        </RouterLink>

        <div class="border-t border-slate-100 mt-1">
          <button
            type="button"
            @click="handleLogout"
            class="w-full flex items-center gap-2 px-3 py-2 text-sm text-slate-700
                   hover:bg-slate-50 text-left"
          >
            <LogOut class="w-4 h-4" /> Logout
          </button>
        </div>
      </div>
    </Transition>
  </div>
</template>
