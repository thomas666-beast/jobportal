<script setup>
import { ref, onMounted, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import {
  MapPin,
  Briefcase,
  Tag,
  Banknote,
  ArrowLeft,
  Globe,
  Building2,
  CheckCircle2,
  Clock
} from 'lucide-vue-next'
import { getPublicJob } from '../api/publicJobs'
import { createApplication, checkApplicationStatus } from '../api/seekerApplications'
import { useAuthStore } from '../stores/auth'
import { useToast } from '../composables/useToast'
import {
  formatSalary,
  remoteLabel,
  employmentLabel,
  applicationStatusLabel,
  applicationStatusVariant,
  timeAgo
} from '../utils/jobStatus'
import ApplyModal from '../components/ApplyModal.vue'
import SaveJobButton from '../components/SaveJobButton.vue'
import AppButton from '../components/ui/AppButton.vue'
import AppCard from '../components/ui/AppCard.vue'
import AppBadge from '../components/ui/AppBadge.vue'

const route = useRoute()
const router = useRouter()
const auth = useAuthStore()
const toast = useToast()

const job = ref(null)
const loading = ref(true)
const errorMsg = ref('')

const applied = ref(false)
const appliedStatus = ref('')

const showModal = ref(false)
const submitting = ref(false)

const isSeeker = computed(
  () =>
    auth.isAuthenticated &&
    auth.user?.role === 'job_seeker' &&
    auth.user?.onboardingCompleted === true
)

const applyState = computed(() => {
  if (!auth.isAuthenticated) {
    return { disabled: false, label: 'Log in as a job seeker to apply', action: 'login' }
  }
  if (auth.user?.role !== 'job_seeker') {
    return { disabled: true, label: 'Only job seekers can apply', action: null }
  }
  if (!auth.user?.onboardingCompleted) {
    return { disabled: false, label: 'Complete onboarding first', action: 'onboarding' }
  }
  if (applied.value) {
    return { disabled: true, label: 'Applied', action: null }
  }
  return { disabled: false, label: 'Apply now', action: 'apply' }
})

function handleApply() {
  const a = applyState.value.action
  if (a === 'login') {
    router.push({ name: 'login', query: { redirect: route.fullPath } })
  } else if (a === 'onboarding') {
    router.push({ name: 'onboarding' })
  } else if (a === 'apply') {
    showModal.value = true
  }
}

async function handleSubmitApplication(payload) {
  submitting.value = true
  try {
    await createApplication(payload)
    showModal.value = false
    applied.value = true
    appliedStatus.value = 'submitted'
    toast.success('Application submitted.')
  } catch (err) {
    if (err.status === 409) {
      applied.value = true
      toast.info('You have already applied to this job.')
    } else {
      toast.error(err.message)
    }
  } finally {
    submitting.value = false
  }
}

async function load() {
  loading.value = true
  errorMsg.value = ''
  try {
    job.value = await getPublicJob(route.params.id)
    if (isSeeker.value) {
      try {
        const st = await checkApplicationStatus(route.params.id)
        applied.value = !!st.applied
        appliedStatus.value = st.status || ''
      } catch {
        // ignore — user can still click Apply
      }
    }
  } catch (err) {
    errorMsg.value = err.message
  } finally {
    loading.value = false
  }
}

onMounted(load)
</script>

<template>
  <div class="space-y-4">
    <button
      type="button"
      @click="router.push({ name: 'browse' })"
      class="inline-flex items-center gap-1 text-sm text-slate-600 hover:text-slate-900"
    >
      <ArrowLeft class="w-4 h-4" /> Back to jobs
    </button>

    <p v-if="loading" class="text-slate-600">Loading…</p>
    <p v-else-if="errorMsg" class="text-red-600">{{ errorMsg }}</p>

    <template v-else-if="job">
      <div class="grid lg:grid-cols-3 gap-6 items-start">
        <!-- Main column -->
        <div class="lg:col-span-2 space-y-4">
          <AppCard>
            <div class="flex items-start gap-4">
              <div class="flex-shrink-0">
                <div
                  v-if="job.company.logo_url"
                  class="w-14 h-14 rounded-xl border border-slate-200 overflow-hidden bg-white"
                >
                  <img
                    :src="job.company.logo_url"
                    :alt="job.company.name"
                    class="w-full h-full object-cover"
                  />
                </div>
                <div
                  v-else
                  class="w-14 h-14 rounded-xl bg-slate-100 flex items-center justify-center"
                >
                  <Building2 class="w-7 h-7 text-slate-400" />
                </div>
              </div>

              <div class="flex-1 min-w-0">
                <h1 class="text-2xl font-semibold text-slate-900">
                  {{ job.title }}
                </h1>
                <p class="text-sm text-slate-600 mt-1">
                  <RouterLink
                    :to="{ name: 'company-public', params: { id: job.company.id } }"
                    class="hover:text-primary-600 font-medium"
                  >
                    {{ job.company.name }}
                  </RouterLink>
                  <span v-if="job.company.location">
                    · {{ job.company.location }}
                  </span>
                </p>

                <div class="flex flex-wrap gap-3 mt-3 text-sm text-slate-600">
                  <span v-if="job.location" class="inline-flex items-center gap-1">
                    <MapPin class="w-4 h-4" /> {{ job.location }}
                  </span>
                  <span class="inline-flex items-center gap-1">
                    <Briefcase class="w-4 h-4" />
                    {{ employmentLabel(job.employment_type) }}
                  </span>
                  <span
                    v-if="formatSalary(job.salary_min, job.salary_max, job.salary_currency)"
                    class="inline-flex items-center gap-1"
                  >
                    <Banknote class="w-4 h-4" />
                    {{ formatSalary(job.salary_min, job.salary_max, job.salary_currency) }}
                  </span>
                </div>

                <div class="flex flex-wrap items-center gap-2 mt-3">
                  <AppBadge variant="neutral">
                    {{ remoteLabel(job.remote_type) }}
                  </AppBadge>
                  <AppBadge v-if="job.category" variant="neutral">
                    <Tag class="w-3 h-3 mr-1" /> {{ job.category }}
                  </AppBadge>
                  <span
                    v-if="job.published_at"
                    class="text-xs text-slate-500 inline-flex items-center gap-1 ml-1"
                  >
                    <Clock class="w-3 h-3" /> Posted {{ timeAgo(job.published_at) }}
                  </span>
                </div>
              </div>
            </div>

            <div v-if="job.tags?.length" class="flex flex-wrap gap-2 mt-4">
              <AppBadge
                v-for="t in job.tags"
                :key="t"
                variant="neutral"
                size="sm"
              >
                {{ t }}
              </AppBadge>
            </div>
          </AppCard>

          <AppCard>
            <h2 class="text-sm font-medium text-slate-500 mb-2">Description</h2>
            <p class="text-slate-800 whitespace-pre-wrap leading-relaxed">
              {{ job.description }}
            </p>
          </AppCard>

          <AppCard v-if="job.requirements">
            <h2 class="text-sm font-medium text-slate-500 mb-2">Requirements</h2>
            <p class="text-slate-800 whitespace-pre-wrap leading-relaxed">
              {{ job.requirements }}
            </p>
          </AppCard>

          <AppCard>
            <h2 class="text-sm font-medium text-slate-500 mb-3">About the company</h2>
            <div class="flex items-start gap-3">
              <div class="flex-shrink-0">
                <div
                  v-if="job.company.logo_url"
                  class="w-12 h-12 rounded-lg border border-slate-200 overflow-hidden bg-white"
                >
                  <img
                    :src="job.company.logo_url"
                    :alt="job.company.name"
                    class="w-full h-full object-cover"
                  />
                </div>
                <div
                  v-else
                  class="w-12 h-12 rounded-lg bg-slate-100 flex items-center justify-center"
                >
                  <Building2 class="w-6 h-6 text-slate-400" />
                </div>
              </div>
              <div class="min-w-0 flex-1">
                <RouterLink
                  :to="{ name: 'company-public', params: { id: job.company.id } }"
                  class="text-base font-semibold text-slate-900 hover:text-primary-600"
                >
                  {{ job.company.name }}
                </RouterLink>
                <div
                  v-if="job.company.location"
                  class="text-sm text-slate-600 mt-0.5"
                >
                  {{ job.company.location }}
                </div>
                <div v-if="job.company.website" class="text-sm mt-1">
                  <a
                    :href="job.company.website"
                    target="_blank"
                    rel="noopener"
                    class="text-primary-600 hover:underline inline-flex items-center gap-1"
                  >
                    <Globe class="w-3.5 h-3.5" /> {{ job.company.website }}
                  </a>
                </div>
                <p
                  v-if="job.company.description"
                  class="text-sm text-slate-700 whitespace-pre-wrap mt-3 leading-relaxed"
                >
                  {{ job.company.description }}
                </p>
              </div>
            </div>
          </AppCard>
        </div>

        <!-- Sidebar — sticky apply + save -->
        <aside class="lg:col-span-1">
          <div class="lg:sticky lg:top-20 space-y-3">
            <AppCard>
              <h3 class="text-sm font-medium text-slate-500 mb-3">Interested?</h3>

              <AppButton
                block
                size="lg"
                :disabled="applyState.disabled"
                @click="handleApply"
              >
                <CheckCircle2 v-if="applied" class="w-4 h-4" />
                {{ applyState.label }}
              </AppButton>

              <div class="mt-3">
                <SaveJobButton :job-id="job.id" variant="full" />
              </div>

              <div
                v-if="applied && appliedStatus"
                class="mt-3 flex items-center gap-2"
              >
                <span class="text-xs text-slate-500">Your status:</span>
                <AppBadge :variant="applicationStatusVariant(appliedStatus)">
                  {{ applicationStatusLabel(appliedStatus) }}
                </AppBadge>
              </div>

              <div class="mt-4 pt-4 border-t border-slate-100 space-y-2 text-sm">
                <div class="flex justify-between gap-2">
                  <span class="text-slate-500">Job type</span>
                  <span class="text-slate-800 text-right">
                    {{ employmentLabel(job.employment_type) }}
                  </span>
                </div>
                <div v-if="job.remote_type" class="flex justify-between gap-2">
                  <span class="text-slate-500">Workplace</span>
                  <span class="text-slate-800 text-right">
                    {{ remoteLabel(job.remote_type) }}
                  </span>
                </div>
                <div v-if="job.location" class="flex justify-between gap-2">
                  <span class="text-slate-500">Location</span>
                  <span class="text-slate-800 text-right">{{ job.location }}</span>
                </div>
                <div
                  v-if="formatSalary(job.salary_min, job.salary_max, job.salary_currency)"
                  class="flex justify-between gap-2"
                >
                  <span class="text-slate-500">Salary</span>
                  <span class="text-slate-800 text-right">
                    {{ formatSalary(job.salary_min, job.salary_max, job.salary_currency) }}
                  </span>
                </div>
              </div>
            </AppCard>
          </div>
        </aside>
      </div>

      <ApplyModal
        :job-id="job.id"
        :job-title="job.title"
        :open="showModal"
        :saving="submitting"
        @close="showModal = false"
        @submit="handleSubmitApplication"
      />
    </template>
  </div>
</template>
