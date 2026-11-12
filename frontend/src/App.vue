<script setup>
import { ref, watch, computed } from 'vue'
import { RouterView, RouterLink, useRoute, useRouter } from 'vue-router'
import { Menu, X } from 'lucide-vue-next'
import { useAuthStore } from './stores/auth'
import { useNotificationsStore } from './stores/notifications'
import { APP_NAME } from './config'
import UserMenu from './components/UserMenu.vue'
import NotificationBell from './components/NotificationBell.vue'
import ToastContainer from './components/ui/ToastContainer.vue'
import EmailVerificationBanner from './components/EmailVerificationBanner.vue'

const auth = useAuthStore()
const notifications = useNotificationsStore()
const route = useRoute()
const router = useRouter()

const menuOpen = ref(false)

const isLanding = computed(() => route.name === 'landing')

watch(() => route.fullPath, () => {
  menuOpen.value = false
})

watch(
  () => auth.isAuthenticated,
  (authed) => {
    if (authed) {
      notifications.startPolling()
    } else {
      notifications.stopPolling()
    }
  },
  { immediate: true }
)

function handleLogout() {
  auth.logout()
  notifications.reset()
  menuOpen.value = false
  router.push({ name: 'login' })
}
</script>

<template>
  <!-- Landing page renders without app chrome -->
  <RouterView v-if="isLanding" />

  <!-- Everything else uses the app shell -->
  <div v-else class="min-h-screen flex flex-col bg-slate-50">
    <header class="bg-white border-b border-slate-200 sticky top-0 z-30">
      <nav class="max-w-6xl mx-auto px-4 h-14 flex items-center justify-between">
        <div class="flex items-center gap-6">
          <RouterLink
            to="/dashboard"
            class="text-base font-semibold text-slate-900"
          >
            {{ APP_NAME }}
          </RouterLink>
          <RouterLink
            to="/browse"
            class="hidden md:inline text-sm text-slate-600 hover:text-slate-900 transition-colors"
          >
            Browse jobs
          </RouterLink>
          <RouterLink
            to="/companies"
            class="hidden md:inline text-sm text-slate-600 hover:text-slate-900 transition-colors"
          >
            Companies
          </RouterLink>
        </div>

        <div class="hidden md:flex items-center gap-3 text-sm">
          <template v-if="auth.isAuthenticated">
            <NotificationBell />
            <UserMenu />
          </template>
          <template v-else>
            <RouterLink
              to="/login"
              class="text-sm text-slate-600 hover:text-slate-900 px-3 py-1.5"
            >
              Login
            </RouterLink>
            <RouterLink
              to="/register"
              class="text-sm text-white bg-primary-600 hover:bg-primary-700
                     font-medium px-4 py-1.5 rounded-lg transition-colors"
            >
              Register
            </RouterLink>
          </template>
        </div>

        <button
          type="button"
          class="md:hidden p-2 text-slate-600 hover:text-slate-900"
          @click="menuOpen = !menuOpen"
          :aria-label="menuOpen ? 'Close menu' : 'Open menu'"
        >
          <component :is="menuOpen ? X : Menu" class="w-5 h-5" />
        </button>
      </nav>

      <!-- Mobile menu -->
      <Transition
        enter-active-class="transition duration-150 ease-out"
        enter-from-class="opacity-0 -translate-y-2"
        leave-active-class="transition duration-100 ease-in"
        leave-to-class="opacity-0 -translate-y-2"
      >
        <div v-if="menuOpen" class="md:hidden border-t border-slate-200 bg-white">
          <div class="px-4 py-3 space-y-1">
            <RouterLink
              to="/browse"
              class="block px-3 py-2 text-sm text-slate-700 hover:bg-slate-50 rounded-lg"
            >
              Browse jobs
            </RouterLink>
            <RouterLink
              to="/companies"
              class="block px-3 py-2 text-sm text-slate-700 hover:bg-slate-50 rounded-lg"
            >
              Companies
            </RouterLink>

            <template v-if="auth.isAuthenticated">
              <div class="border-t border-slate-100 my-2"></div>
              <div class="px-3 py-1 text-xs text-slate-500">
                {{ auth.user?.username }}
              </div>

              <RouterLink
                to="/notifications"
                class="flex items-center justify-between px-3 py-2 text-sm text-slate-700 hover:bg-slate-50 rounded-lg"
              >
                <span>Notifications</span>
                <span
                  v-if="notifications.unread > 0"
                  class="min-w-[20px] h-5 px-1.5 rounded-full bg-red-500 text-white
                         text-[10px] font-semibold flex items-center justify-center"
                >
                  {{ notifications.unread > 99 ? '99+' : notifications.unread }}
                </span>
              </RouterLink>

              <RouterLink
                v-if="auth.user?.role === 'job_seeker' && auth.user?.onboardingCompleted"
                to="/profile"
                class="block px-3 py-2 text-sm text-slate-700 hover:bg-slate-50 rounded-lg"
              >
                Profile
              </RouterLink>
              <RouterLink
                v-if="auth.user?.role === 'job_seeker' && auth.user?.onboardingCompleted"
                to="/applications"
                class="block px-3 py-2 text-sm text-slate-700 hover:bg-slate-50 rounded-lg"
              >
                My applications
              </RouterLink>
              <RouterLink
                v-if="auth.user?.role === 'job_seeker' && auth.user?.onboardingCompleted"
                to="/saved"
                class="block px-3 py-2 text-sm text-slate-700 hover:bg-slate-50 rounded-lg"
              >
                Saved jobs
              </RouterLink>
              <RouterLink
                v-if="auth.user?.role === 'job_seeker' && auth.user?.onboardingCompleted"
                to="/searches"
                class="block px-3 py-2 text-sm text-slate-700 hover:bg-slate-50 rounded-lg"
              >
                Saved searches
              </RouterLink>

              <RouterLink
                v-if="auth.user?.role === 'employer' && auth.user?.onboardingCompleted"
                to="/company"
                class="block px-3 py-2 text-sm text-slate-700 hover:bg-slate-50 rounded-lg"
              >
                Company
              </RouterLink>
              <RouterLink
                v-if="auth.user?.role === 'employer' && auth.user?.onboardingCompleted"
                to="/jobs"
                class="block px-3 py-2 text-sm text-slate-700 hover:bg-slate-50 rounded-lg"
              >
                Jobs
              </RouterLink>
              <RouterLink
                v-if="auth.user?.role === 'employer' && auth.user?.onboardingCompleted"
                to="/employer/applications"
                class="block px-3 py-2 text-sm text-slate-700 hover:bg-slate-50 rounded-lg"
              >
                Applications
              </RouterLink>

              <RouterLink
                v-if="auth.user?.role === 'admin'"
                to="/admin"
                class="block px-3 py-2 text-sm text-slate-700 hover:bg-slate-50 rounded-lg"
              >
                Admin
              </RouterLink>

              <RouterLink
                to="/settings"
                class="block px-3 py-2 text-sm text-slate-700 hover:bg-slate-50 rounded-lg"
              >
                Settings
              </RouterLink>

              <button
                type="button"
                @click="handleLogout"
                class="w-full text-left px-3 py-2 text-sm text-slate-700 hover:bg-slate-50 rounded-lg"
              >
                Logout
              </button>
            </template>

            <template v-else>
              <div class="border-t border-slate-100 my-2"></div>
              <RouterLink
                to="/login"
                class="block px-3 py-2 text-sm text-slate-700 hover:bg-slate-50 rounded-lg"
              >
                Login
              </RouterLink>
              <RouterLink
                to="/register"
                class="block px-3 py-2 text-sm font-medium text-primary-600 hover:bg-primary-50 rounded-lg"
              >
                Register
              </RouterLink>
            </template>
          </div>
        </div>
      </Transition>
    </header>

    <main class="flex-1 max-w-6xl mx-auto w-full px-4 py-8 space-y-4">
      <EmailVerificationBanner />
      <RouterView />
    </main>

    <footer class="border-t border-slate-200 bg-white mt-12">
      <div class="max-w-6xl mx-auto px-4 py-6 flex flex-col md:flex-row items-center justify-between gap-3 text-sm text-slate-500">
        <p>© {{ new Date().getFullYear() }} {{ APP_NAME }}</p>
        <div class="flex items-center gap-4">
          <RouterLink to="/browse" class="hover:text-slate-900">Browse jobs</RouterLink>
          <RouterLink to="/companies" class="hover:text-slate-900">Companies</RouterLink>
          <RouterLink
            v-if="!auth.isAuthenticated"
            to="/register"
            class="hover:text-slate-900"
          >
            Create account
          </RouterLink>
        </div>
      </div>
    </footer>

    <ToastContainer />
  </div>
</template>
