<script setup>
import { ref, onMounted } from 'vue'
import { RouterLink } from 'vue-router'
import {
  Bookmark,
  MapPin,
  Briefcase,
  Banknote,
  Building2,
  Clock
} from 'lucide-vue-next'
import { listSavedJobs, unsaveJob } from '../api/savedJobs'
import {
  timeAgo,
  formatSalary,
  remoteLabel,
  employmentLabel
} from '../utils/jobStatus'
import AppButton from '../components/ui/AppButton.vue'
import AppBadge from '../components/ui/AppBadge.vue'
import AppCard from '../components/ui/AppCard.vue'
import EmptyState from '../components/ui/EmptyState.vue'
import { useToast } from '../composables/useToast'

const toast = useToast()

const jobs = ref([])
const loading = ref(true)
const errorMsg = ref('')
const acting = ref(null)

async function load() {
  loading.value = true
  errorMsg.value = ''
  try {
    jobs.value = await listSavedJobs()
  } catch (err) {
    errorMsg.value = err.message
  } finally {
    loading.value = false
  }
}

async function remove(jobId) {
  acting.value = jobId
  try {
    await unsaveJob(jobId)
    jobs.value = jobs.value.filter((j) => j.job_id !== jobId)
    toast.success('Removed from saved jobs.')
  } catch (err) {
    toast.error(err.message)
  } finally {
    acting.value = null
  }
}

onMounted(load)
</script>

<template>
  <div class="space-y-4">
    <div>
      <h1 class="text-2xl font-semibold text-slate-900">Saved jobs</h1>
      <p class="text-sm text-slate-600 mt-1">
        Jobs you've bookmarked for later
      </p>
    </div>

    <p v-if="loading" class="text-slate-600">Loading…</p>
    <p v-else-if="errorMsg" class="text-red-600">{{ errorMsg }}</p>

    <AppCard v-else-if="jobs.length === 0" :padded="false">
      <EmptyState
        :icon="Bookmark"
        title="No saved jobs yet"
        description="Tap the bookmark icon on any job to save it here."
      >
        <template #action>
          <AppButton to="/browse">Browse jobs</AppButton>
        </template>
      </EmptyState>
    </AppCard>

    <div v-else class="space-y-3">
      <AppCard v-for="job in jobs" :key="job.job_id">
        <div class="flex gap-4">
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

          <div class="flex-1 min-w-0">
            <div class="flex items-start justify-between gap-3">
              <div class="min-w-0">
                <RouterLink
                  :to="{ name: 'job-detail', params: { id: job.job_id } }"
                  class="text-lg font-semibold text-slate-900 hover:text-primary-600"
                >
                  {{ job.title }}
                </RouterLink>
                <p class="text-sm text-slate-600 mt-0.5">
                  {{ job.company.name }}
                  <span v-if="job.company.location"> · {{ job.company.location }}</span>
                </p>
              </div>
              <AppBadge variant="neutral" size="sm" class="flex-shrink-0">
                {{ remoteLabel(job.remote_type) }}
              </AppBadge>
            </div>

            <p class="text-sm text-slate-700 mt-2 line-clamp-2">
              {{ job.description }}
            </p>

            <div class="flex flex-wrap items-center gap-x-4 gap-y-1.5 mt-3 text-xs text-slate-500">
              <span v-if="job.location" class="inline-flex items-center gap-1">
                <MapPin class="w-3.5 h-3.5" /> {{ job.location }}
              </span>
              <span class="inline-flex items-center gap-1">
                <Briefcase class="w-3.5 h-3.5" />
                {{ employmentLabel(job.employment_type) }}
              </span>
              <span
                v-if="formatSalary(job.salary_min, job.salary_max, job.salary_currency)"
                class="inline-flex items-center gap-1"
              >
                <Banknote class="w-3.5 h-3.5" />
                {{ formatSalary(job.salary_min, job.salary_max, job.salary_currency) }}
              </span>
              <span class="inline-flex items-center gap-1 ml-auto">
                <Clock class="w-3.5 h-3.5" /> Saved {{ timeAgo(job.saved_at) }}
              </span>
            </div>

            <div class="flex items-center justify-between mt-3">
              <div class="flex flex-wrap gap-1.5">
                <AppBadge
                  v-for="t in job.tags || []"
                  :key="t"
                  variant="neutral"
                  size="sm"
                >
                  {{ t }}
                </AppBadge>
              </div>
              <AppButton
                size="sm"
                variant="ghost"
                :disabled="acting === job.job_id"
                @click="remove(job.job_id)"
              >
                Remove
              </AppButton>
            </div>
          </div>
        </div>
      </AppCard>
    </div>
  </div>
</template>
