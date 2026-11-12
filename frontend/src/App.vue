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

watch(() => route.fullPath, () => { menuOpen.value = false })
watch(() => auth.isAuthenticated, (authed) => {
  if (authed) notifications.startPolling()
  else notifications.stopPolling()
}, { immediate: true })

function handleLogout() {
  auth.logout()
  notifications.reset()
  menuOpen.value = false
  router.push({ name: 'login' })
}
</script>

<template>
  <RouterView v-if="isLanding" />

  <div v-else class="min-h-screen flex flex-col bg-bg text-fg">
    <!-- HEADER -->
    <header class="bg-surface border-b border-border sticky top-0 z-30">
      <nav class="max-w-6xl mx-auto px-4 h-16 flex items-center justify-between">
        <div class="flex items-center gap-6">
          <RouterLink to="/dashboard" class="text-xl font-bold goth-heading text-fg hover:text-primary transition-colors">
            {{ APP_NAME }}
          </RouterLink>
          <RouterLink to="/browse" class="hidden md:inline text-sm text-fg-muted hover:text-fg transition-colors">Browse jobs</RouterLink>
          <RouterLink to="/companies" class="hidden md:inline text-sm text-fg-muted hover:text-fg transition-colors">Companies</RouterLink>
        </div>

        <div class="hidden md:flex items-center gap-4 text-sm">
          <template v-if="auth.isAuthenticated">
            <NotificationBell />
            <UserMenu />
          </template>
          <template v-else>
            <RouterLink to="/login" class="text-sm text-fg-muted hover:text-fg px-3 py-1.5 transition-colors">Login</RouterLink>
            <RouterLink to="/register" class="text-sm text-white bg-primary hover:bg-primary-hover font-bold uppercase tracking-wider px-5 py-2 transition-colors">Register</RouterLink>
          </template>
        </div>

        <button type="button" class="md:hidden p-2 text-fg-muted hover:text-fg" @click="menuOpen = !menuOpen">
          <component :is="menuOpen ? X : Menu" class="w-6 h-6" />
        </button>
      </nav>

      <!-- MOBILE MENU -->
      <Transition enter-active-class="transition duration-150 ease-out" enter-from-class="opacity-0 -translate-y-2" leave-active-class="transition duration-100 ease-in" leave-to-class="opacity-0 -translate-y-2">
        <div v-if="menuOpen" class="md:hidden border-t border-border bg-surface">
          <div class="px-4 py-3 space-y-1">
            <RouterLink to="/browse" class="block px-3 py-2 text-sm text-fg-muted hover:bg-surface-muted hover:text-fg">Browse jobs</RouterLink>
            <RouterLink to="/companies" class="block px-3 py-2 text-sm text-fg-muted hover:bg-surface-muted hover:text-fg">Companies</RouterLink>
            <template v-if="auth.isAuthenticated">
              <div class="border-t border-border my-2"></div>
              <div class="px-3 py-1 text-xs text-fg-subtle font-bold uppercase">{{ auth.user?.username }}</div>
              <RouterLink to="/settings" class="block px-3 py-2 text-sm text-fg-muted hover:bg-surface-muted hover:text-fg">Settings</RouterLink>
              <button type="button" @click="handleLogout" class="w-full text-left px-3 py-2 text-sm text-fg-muted hover:bg-surface-muted hover:text-fg">Logout</button>
            </template>
            <template v-else>
              <div class="border-t border-border my-2"></div>
              <RouterLink to="/login" class="block px-3 py-2 text-sm text-fg-muted hover:bg-surface-muted hover:text-fg">Login</RouterLink>
              <RouterLink to="/register" class="block px-3 py-2 text-sm font-bold text-primary hover:bg-surface-muted">Register</RouterLink>
            </template>
          </div>
        </div>
      </Transition>
    </header>

    <!-- MAIN CONTENT -->
    <main class="flex-1 max-w-6xl mx-auto w-full px-4 py-8 space-y-6">
      <EmailVerificationBanner />
      <RouterView />
    </main>

    <!-- FOOTER -->
    <footer class="border-t border-border bg-surface mt-12">
      <div class="max-w-6xl mx-auto px-4 py-8 flex flex-col md:flex-row items-center justify-between gap-4 text-sm text-fg-muted">
        <p>© {{ new Date().getFullYear() }} {{ APP_NAME }}</p>
        <div class="flex items-center gap-6">
          <RouterLink to="/browse" class="hover:text-fg transition-colors">Browse jobs</RouterLink>
          <RouterLink to="/companies" class="hover:text-fg transition-colors">Companies</RouterLink>
          <RouterLink v-if="!auth.isAuthenticated" to="/register" class="hover:text-fg transition-colors">Create account</RouterLink>
        </div>
      </div>
    </footer>

    <ToastContainer />
  </div>
</template>
