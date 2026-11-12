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
    toast.success('You\'re all set!')
    router.push({ name: 'dashboard' })
  } catch (err) {
    errorMsg.value = err.message
  } finally {
    loading.value = null
  }
}
</script>

<template>
  <div class="max-w-3xl mx-auto py-8">
    <div class="text-center mb-8">
      <h1 class="text-3xl font-semibold text-slate-900">
        Welcome, {{ auth.user?.username }}
      </h1>
      <p class="mt-2 text-slate-600">How do you want to use Job Portal?</p>
    </div>

    <div class="grid md:grid-cols-2 gap-4">
      <button
        type="button"
        @click="pick('job_seeker')"
        :disabled="loading !== null"
        class="group text-left p-6 bg-white rounded-xl border-2 border-slate-200
               hover:border-primary-500 hover:shadow-sm
               disabled:opacity-50 disabled:cursor-not-allowed
               transition-all"
      >
        <div
          class="w-12 h-12 rounded-xl bg-primary-50 text-primary-600
                 flex items-center justify-center mb-4
                 group-hover:bg-primary-100 transition-colors"
        >
          <Search class="w-6 h-6" :stroke-width="2" />
        </div>
        <h2 class="text-lg font-semibold text-slate-900 mb-1">
          I'm looking for a job
        </h2>
        <p class="text-sm text-slate-600">
          Browse jobs, save favorites, and apply to positions.
        </p>
        <p v-if="loading === 'job_seeker'" class="text-xs text-primary-600 mt-3">
          Setting up…
        </p>
      </button>

      <button
        type="button"
        @click="pick('employer')"
        :disabled="loading !== null"
        class="group text-left p-6 bg-white rounded-xl border-2 border-slate-200
               hover:border-emerald-500 hover:shadow-sm
               disabled:opacity-50 disabled:cursor-not-allowed
               transition-all"
      >
        <div
          class="w-12 h-12 rounded-xl bg-emerald-50 text-emerald-600
                 flex items-center justify-center mb-4
                 group-hover:bg-emerald-100 transition-colors"
        >
          <Building2 class="w-6 h-6" :stroke-width="2" />
        </div>
        <h2 class="text-lg font-semibold text-slate-900 mb-1">
          I'm hiring
        </h2>
        <p class="text-sm text-slate-600">
          Post jobs, manage listings, and review applications.
        </p>
        <p v-if="loading === 'employer'" class="text-xs text-emerald-600 mt-3">
          Setting up…
        </p>
      </button>
    </div>

    <p v-if="errorMsg" class="mt-4 text-sm text-red-600 text-center">
      {{ errorMsg }}
    </p>

    <p class="mt-6 text-xs text-slate-500 text-center max-w-md mx-auto">
      This choice is permanent for this account. If you need both roles,
      register a second account with a different email.
    </p>
  </div>
</template>
