<script setup>
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import {
  ArrowLeft,
  MapPin,
  Globe,
  Building2,
  Briefcase,
  Banknote,
  Clock,
  Search
} from 'lucide-vue-next'
import { getPublicCompany } from '../api/publicCompanies'
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
import SaveJobButton from '../components/SaveJobButton.vue'

const route = useRoute()
const router = useRouter()

const company = ref(null)
const loading = ref(true)
const errorMsg = ref('')

async function load() {
  loading.value = true
  errorMsg.value = ''
  try {
    company.value = await getPublicCompany(route.params.id)
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
      @click="router.push({ name: 'companies' })"
      class="inline-flex items-center gap-1 text-sm text-slate-600 hover:text-slate-900"
    >
      <ArrowLeft class="w-4 h-4" /> Back to companies
    </button>

    <p v-if="loading" class="text-slate-600">Loading…</p>
    <p v-else-if="errorMsg" class="text-red-600">{{ errorMsg }}</p>

    <template v-else-if="company">
      <!-- Header card -->
      <AppCard>
        <div class="flex items-start gap-5">
          <div class="flex-shrink-0">
            <div
              v-if="company.logo_url"
              class="w-20 h-20 rounded-2xl border border-slate-200 overflow-hidden bg-white"
            >
              <img
                :src="company.logo_url"
                :alt="company.name"
                class="w-full h-full object-cover"
              />
            </div>
            <div
              v-else
              class="w-20 h-20 rounded-2xl bg-slate-100 flex items-center justify-center"
            >
              <Building2 class="w-10 h-10 text-slate-400" />
            </div>
          </div>

          <div class="flex-1 min-w-0">
            <h1 class="text-2xl font-semibold text-slate-900">
              {{ company.name }}
            </h1>

            <div class="flex flex-wrap gap-3 mt-2 text-sm text-slate-600">
              <span v-if="company.location" class="inline-flex items-center gap-1">
                <MapPin class="w-4 h-4" /> {{ company.location }}
              </span>
              <span v-if="company.website" class="inline-flex items-center gap-1">
                <Globe class="w-4 h-4" />
                <a
                  :href="company.website"
                  target="_blank"
                  rel="noopener"
                  class="text-primary-600 hover:underline"
                >
                  {{ company.website }}
                </a>
              </span>
            </div>

            <div class="mt-3">
              <AppBadge
                :variant="company.published_job_count > 0 ? 'primary' : 'neutral'"
              >
                <Briefcase class="w-3 h-3 mr-1" />
                {{ company.published_job_count }}
                {{ company.published_job_count === 1 ? 'open job' : 'open jobs' }}
              </AppBadge>
            </div>
          </div>
        </div>

        <p
          v-if="company.description"
          class="text-sm text-slate-700 whitespace-pre-wrap mt-5 leading-relaxed"
        >
          {{ company.description }}
        </p>
      </AppCard>

      <!-- Jobs -->
      <div>
        <h2 class="text-lg font-semibold text-slate-900 mb-3">Open positions</h2>

        <AppCard v-if="company.jobs.length === 0" :padded="false">
          <EmptyState
            :icon="Search"
            title="No open positions"
            description="This company has no published jobs right now."
          >
            <template #action>
              <AppButton to="/browse" variant="secondary">Browse all jobs</AppButton>
            </template>
          </EmptyState>
        </AppCard>

        <div v-else class="space-y-3">
          <RouterLink
            v-for="job in company.jobs"
            :key="job.id"
            :to="{ name: 'job-detail', params: { id: job.id } }"
            class="block"
          >
            <AppCard hover>
              <div class="flex items-start justify-between gap-4">
                <div class="min-w-0 flex-1">
                  <h3 class="text-lg font-semibold text-slate-900 truncate">
                    {{ job.title }}
                  </h3>
                  <p class="text-sm text-slate-700 mt-1 line-clamp-2">
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
                    <span v-if="job.published_at" class="inline-flex items-center gap-1 ml-auto">
                      <Clock class="w-3.5 h-3.5" /> {{ timeAgo(job.published_at) }}
                    </span>
                  </div>

                  <div v-if="job.tags?.length" class="flex flex-wrap gap-1.5 mt-3">
                    <AppBadge
                      v-for="t in job.tags"
                      :key="t"
                      variant="neutral"
                      size="sm"
                    >
                      {{ t }}
                    </AppBadge>
                  </div>
                </div>

                <div class="flex flex-col items-end gap-2 flex-shrink-0">
                  <AppBadge variant="neutral" size="sm">
                    {{ remoteLabel(job.remote_type) }}
                  </AppBadge>
                  <SaveJobButton :job-id="job.id" variant="icon" />
                </div>
              </div>
            </AppCard>
          </RouterLink>
        </div>
      </div>
    </template>
  </div>
</template>
