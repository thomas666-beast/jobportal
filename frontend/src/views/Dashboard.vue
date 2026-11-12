<script setup>
import { onMounted, ref } from 'vue'
import { RouterLink } from 'vue-router'
import { Search, Building2, Plus, UserCircle, FileText } from 'lucide-vue-next'
import { useAuthStore } from '../stores/auth'
import { getProfile } from '../api/seeker'
import { getCompany } from '../api/company'
import { listJobs } from '../api/jobs'
import AppCard from '../components/ui/AppCard.vue'
import AppBadge from '../components/ui/AppBadge.vue'
import AppButton from '../components/ui/AppButton.vue'

const auth = useAuthStore()
const loading = ref(false)
const errorMsg = ref('')
const hasSeekerProfile = ref(null)
const hasCompany = ref(null)
const jobCount = ref(null)

onMounted(async () => {
  loading.value = true
  try {
    await auth.fetchMe()
    if (auth.user?.role === 'job_seeker' && auth.user?.onboardingCompleted) {
      try { await getProfile(); hasSeekerProfile.value = true } catch { hasSeekerProfile.value = false }
    }
    if (auth.user?.role === 'employer' && auth.user?.onboardingCompleted) {
      try {
        await getCompany()
        hasCompany.value = true
        const jobs = await listJobs()
        jobCount.value = jobs.length
      } catch { hasCompany.value = false }
    }
  } catch (err) { errorMsg.value = err.message } finally { loading.value = false }
})

function roleVariant(r) {
  if (r === 'job_seeker') return 'primary'
  if (r === 'employer') return 'success'
  if (r === 'admin') return 'warning'
  return 'neutral'
}
function roleLabel(r) {
  if (r === 'job_seeker') return 'Job seeker'
  if (r === 'employer') return 'Employer'
  if (r === 'admin') return 'Admin'
  return r || ''
}
</script>

<template>
  <div class="space-y-6">
    <AppCard>
      <div class="flex items-start justify-between gap-4">
        <div>
          <h1 class="text-2xl font-bold text-fg">Welcome, {{ auth.user?.username }}</h1>
          <p v-if="loading" class="text-sm text-fg-muted mt-1">Loading…</p>
          <p v-else-if="errorMsg" class="text-sm text-red-500 mt-1">{{ errorMsg }}</p>
          <p v-else class="text-sm text-fg-muted mt-1">Here's your overview</p>
        </div>
        <AppBadge v-if="auth.user?.role" :variant="roleVariant(auth.user.role)">
          {{ roleLabel(auth.user.role) }}
        </AppBadge>
      </div>
    </AppCard>

    <!-- Seeker banner: no profile -->
    <div v-if="auth.user?.role === 'job_seeker' && hasSeekerProfile === false" class="bg-surface border border-warning p-6 flex items-center justify-between gap-4">
      <div class="flex items-start gap-3">
        <UserCircle class="w-5 h-5 text-warning flex-shrink-0 mt-0.5" />
        <div>
          <p class="font-bold text-warning uppercase tracking-wide">Complete your profile</p>
          <p class="text-sm text-fg-muted mt-1">Employers will see this when you apply for jobs.</p>
        </div>
      </div>
      <AppButton to="/profile" size="sm">Create profile</AppButton>
    </div>

    <!-- Seeker banner: has profile -->
    <div v-if="auth.user?.role === 'job_seeker' && hasSeekerProfile === true" class="bg-surface border border-primary p-6 flex items-center justify-between gap-4">
      <div class="flex items-start gap-3">
        <Search class="w-5 h-5 text-primary flex-shrink-0 mt-0.5" />
        <div>
          <p class="font-bold text-primary uppercase tracking-wide">Looking for a job?</p>
          <p class="text-sm text-fg-muted mt-1">Browse published listings and find your next role.</p>
        </div>
      </div>
      <AppButton to="/browse" size="sm">Browse jobs</AppButton>
    </div>

    <!-- Employer banner: no company -->
    <div v-if="auth.user?.role === 'employer' && hasCompany === false" class="bg-surface border border-warning p-6 flex items-center justify-between gap-4">
      <div class="flex items-start gap-3">
        <Building2 class="w-5 h-5 text-warning flex-shrink-0 mt-0.5" />
        <div>
          <p class="font-bold text-warning uppercase tracking-wide">Set up your company</p>
          <p class="text-sm text-fg-muted mt-1">You'll need a company profile before posting jobs.</p>
        </div>
      </div>
      <AppButton to="/company" size="sm">Create company</AppButton>
    </div>

    <!-- Employer banner: has company, no jobs -->
    <div v-if="auth.user?.role === 'employer' && hasCompany === true && jobCount === 0" class="bg-surface border border-success p-6 flex items-center justify-between gap-4">
      <div class="flex items-start gap-3">
        <Plus class="w-5 h-5 text-success flex-shrink-0 mt-0.5" />
        <div>
          <p class="font-bold text-success uppercase tracking-wide">Post your first job</p>
          <p class="text-sm text-fg-muted mt-1">Companies with listings get more visibility.</p>
        </div>
      </div>
      <AppButton to="/jobs/new" size="sm" variant="success">New job</AppButton>
    </div>

    <!-- Employer banner: has company + jobs -->
    <div v-if="auth.user?.role === 'employer' && hasCompany === true && jobCount > 0" class="bg-surface border border-border p-6 flex items-center justify-between gap-4">
      <div class="flex items-start gap-3">
        <FileText class="w-5 h-5 text-fg-muted flex-shrink-0 mt-0.5" />
        <div>
          <p class="font-bold text-fg uppercase tracking-wide">You have {{ jobCount }} {{ jobCount === 1 ? 'job' : 'jobs' }}</p>
          <p class="text-sm text-fg-muted mt-1">Manage listings and review applications.</p>
        </div>
      </div>
      <div class="flex gap-2">
        <AppButton to="/jobs" size="sm" variant="secondary">Manage jobs</AppButton>
        <AppButton to="/employer/applications" size="sm">Applications</AppButton>
      </div>
    </div>

    <!-- Admin hint -->
    <div v-if="auth.user?.role === 'admin'" class="bg-surface border border-border p-6 flex items-center justify-between gap-4">
      <div class="flex items-start gap-3">
        <Building2 class="w-5 h-5 text-fg-muted flex-shrink-0 mt-0.5" />
        <div>
          <p class="font-bold text-fg uppercase tracking-wide">Admin access</p>
          <p class="text-sm text-fg-muted mt-1">Review the moderation queue and platform stats.</p>
        </div>
      </div>
      <AppButton to="/admin" size="sm">Open admin</AppButton>
    </div>
  </div>
</template>
