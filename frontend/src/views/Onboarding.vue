<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { Search, Building2 } from 'lucide-vue-next'
import { useAuthStore } from '../stores/auth'
import { useToast } from '../composables/useToast'

const router = useRouter()
const auth = useAuthStore()
const toast = useToast()

const loading = ref(null)
const errorMsg = ref('')

async function pick(role) {
  errorMsg.value = ''
  loading.value = role
  try {
    await auth.completeOnboarding(role)
    toast.success("You're all set!")
    router.push({ name: 'dashboard' })
  } catch (err) {
    errorMsg.value = err.message
  } finally {
    loading.value = null
  }
}
</script>

<template>
  <div class="max-w-3xl mx-auto py-12">
    <div class="text-center mb-10">
      <h1 class="text-4xl font-bold text-fg">Welcome, {{ auth.user?.username }}</h1>
      <p class="mt-3 text-lg text-fg-muted">How do you want to use Job Portal?</p>
    </div>

    <div class="grid md:grid-cols-2 gap-6">
      <button
        type="button"
        @click="pick('job_seeker')"
        :disabled="loading !== null"
        class="group text-left p-8 bg-surface border border-border-strong
               hover:border-primary hover:-translate-y-1 transition-all duration-200
               disabled:opacity-50 disabled:cursor-not-allowed"
      >
        <div class="w-14 h-14 border border-primary text-primary flex items-center justify-center mb-6 group-hover:bg-primary group-hover:text-white transition-colors">
          <Search class="w-7 h-7" />
        </div>
        <h2 class="text-2xl font-bold text-fg mb-2">I'm looking for a job</h2>
        <p class="text-sm text-fg-muted">Browse jobs, save favorites, and apply to positions.</p>
        <p v-if="loading === 'job_seeker'" class="text-xs text-primary mt-4 font-bold uppercase tracking-widest animate-pulse">Setting up…</p>
      </button>

      <button
        type="button"
        @click="pick('employer')"
        :disabled="loading !== null"
        class="group text-left p-8 bg-surface border border-border-strong
               hover:border-primary hover:-translate-y-1 transition-all duration-200
               disabled:opacity-50 disabled:cursor-not-allowed"
      >
        <div class="w-14 h-14 border border-primary text-primary flex items-center justify-center mb-6 group-hover:bg-primary group-hover:text-white transition-colors">
          <Building2 class="w-7 h-7" />
        </div>
        <h2 class="text-2xl font-bold text-fg mb-2">I'm hiring</h2>
        <p class="text-sm text-fg-muted">Post jobs, manage listings, and review applications.</p>
        <p v-if="loading === 'employer'" class="text-xs text-primary mt-4 font-bold uppercase tracking-widest animate-pulse">Setting up…</p>
      </button>
    </div>

    <p v-if="errorMsg" class="mt-6 text-sm text-red-500 text-center">{{ errorMsg }}</p>

    <p class="mt-10 text-xs text-fg-subtle text-center max-w-md mx-auto uppercase tracking-wider">
      This choice is permanent for this account. If you need both roles, register a second account with a different email.
    </p>
  </div>
</template>
