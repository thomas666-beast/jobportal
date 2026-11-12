<script setup>
import { ref, onMounted, computed } from 'vue'
import { RouterLink } from 'vue-router'
import {
  Search,
  MapPin,
  Briefcase,
  Banknote,
  Clock,
  Building2,
  X,
  Bookmark
} from 'lucide-vue-next'
import { browseJobs } from '../api/publicJobs'
import {
  timeAgo,
  formatSalary,
  remoteLabel,
  employmentLabel,
  isNewJob
} from '../utils/jobStatus'
import { useAuthStore } from '../stores/auth'
import AppButton from '../components/ui/AppButton.vue'
import AppBadge from '../components/ui/AppBadge.vue'
import AppCard from '../components/ui/AppCard.vue'
import EmptyState from '../components/ui/EmptyState.vue'
import JobCardSkeleton from '../components/ui/JobCardSkeleton.vue'
import SaveJobButton from '../components/SaveJobButton.vue'
import SaveSearchModal from '../components/SaveSearchModal.vue'
import { useRoute } from 'vue-router'

const route = useRoute()
const auth = useAuthStore()

const q = ref('')
const location = ref('')
const remoteType = ref('')

const jobs = ref([])
const total = ref(0)
const loading = ref(false)
const errorMsg = ref('')

const limit = 20
const offset = ref(0)

const canLoadMore = computed(() => offset.value + jobs.value.length < total.value)
const hasFilters = computed(() => q.value.trim() || location.value.trim() || remoteType.value)
const initialLoading = ref(true)

const showSaveSearch = ref(false)
const canSaveSearch = computed(
  () =>
    auth.isAuthenticated &&
    auth.user?.role === 'job_seeker' &&
    auth.user?.onboardingCompleted === true &&
    hasFilters.value
)

let debounceTimer = null

async function load(reset = true) {
  loading.value = true
  errorMsg.value = ''
  try {
    const res = await browseJobs({
      q: q.value.trim(),
      location: location.value.trim(),
      remote_type: remoteType.value,
      limit,
      offset: reset ? 0 : offset.value
    })
    if (reset) {
      jobs.value = res.jobs || []
      offset.value = 0
    } else {
      jobs.value = [...jobs.value, ...(res.jobs || [])]
    }
    total.value = res.total || 0
  } catch (err) {
    errorMsg.value = err.message
  } finally {
    loading.value = false
    initialLoading.value = false
  }
}

function debouncedLoad() {
  clearTimeout(debounceTimer)
  debounceTimer = setTimeout(() => load(true), 300)
}

function loadMore() {
  offset.value = jobs.value.length
  load(false)
}

function clearFilters() {
  q.value = ''
  location.value = ''
  remoteType.value = ''
  load(true)
}

onMounted(() => {
  if (typeof route.query.q === 'string') q.value = route.query.q
  if (typeof route.query.location === 'string') location.value = route.query.location
  if (typeof route.query.remote_type === 'string') remoteType.value = route.query.remote_type
  load(true)
})
</script>

<template>
  <div class="space-y-6">
    <section class="text-center pt-6 pb-2">
      <h1 class="text-3xl md:text-4xl font-semibold text-slate-900">
        Find your next role
      </h1>
      <p class="mt-2 text-slate-600">
        {{ total }} {{ total === 1 ? 'job' : 'jobs' }} from companies hiring now
      </p>
    </section>

    <AppCard class="!rounded-2xl">
      <div class="flex flex-col md:flex-row gap-3">
        <div class="flex-1 relative">
          <Search class="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-slate-400 pointer-events-none" />
          <input
            v-model="q"
            @input="debouncedLoad"
            type="text"
            placeholder="Job title, skill, or keyword"
            class="w-full rounded-lg border border-slate-300 pl-9 pr-3 py-2.5 text-sm
                   placeholder:text-slate-400
                   focus:outline-none focus:ring-2 focus:ring-primary-500"
          />
        </div>

        <div class="flex-1 relative">
          <MapPin class="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-slate-400 pointer-events-none" />
          <input
            v-model="location"
            @input="debouncedLoad"
            type="text"
            placeholder="City or region"
            class="w-full rounded-lg border border-slate-300 pl-9 pr-3 py-2.5 text-sm
                   placeholder:text-slate-400
                   focus:outline-none focus:ring-2 focus:ring-primary-500"
          />
        </div>

        <select
          v-model="remoteType"
          @change="load(true)"
          class="rounded-lg border border-slate-300 px-3 py-2.5 text-sm bg-white
                 focus:outline-none focus:ring-2 focus:ring-primary-500"
        >
          <option value="">Any remote type</option>
          <option value="remote">Remote</option>
          <option value="hybrid">Hybrid</option>
          <option value="on_site">On-site</option>
        </select>

        <AppButton v-if="hasFilters" variant="ghost" @click="clearFilters">
          <X class="w-4 h-4" /> Clear
        </AppButton>
      </div>

      <div v-if="canSaveSearch" class="mt-3 pt-3 border-t border-slate-100 flex justify-end">
        <AppButton variant="ghost" size="sm" @click="showSaveSearch = true">
          <Bookmark class="w-4 h-4" /> Save this search
        </AppButton>
      </div>
    </AppCard>

    <div v-if="initialLoading" class="space-y-3">
      <JobCardSkeleton v-for="i in 4" :key="i" />
    </div>
    <p v-else-if="errorMsg" class="text-red-600 text-center">{{ errorMsg }}</p>

    <AppCard v-else-if="jobs.length === 0" :padded="false">
      <EmptyState
        :icon="Search"
        title="No jobs match your search"
        description="Try different keywords or clear your filters."
      >
        <template #action>
          <AppButton variant="secondary" @click="clearFilters">Clear filters</AppButton>
        </template>
      </EmptyState>
    </AppCard>

    <div v-else class="space-y-3">
      <RouterLink
        v-for="job in jobs"
        :key="job.id"
        :to="{ name: 'job-detail', params: { id: job.id } }"
        class="block"
      >
        <AppCard hover>
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
                <div class="min-w-0 flex-1">
                  <div class="flex items-center gap-2">
                    <h2 class="text-lg font-semibold text-slate-900 truncate">
                      {{ job.title }}
                    </h2>
                    <AppBadge
                      v-if="isNewJob(job.published_at)"
                      variant="success"
                      size="sm"
                    >
                      New
                    </AppBadge>
                  </div>
                  <p class="text-sm text-slate-600 mt-0.5">
                    {{ job.company.name }}
                    <span v-if="job.company.location"> · {{ job.company.location }}</span>
                  </p>
                </div>
                <div class="flex items-center gap-2 flex-shrink-0">
                  <AppBadge variant="neutral" size="sm">
                    {{ remoteLabel(job.remote_type) }}
                  </AppBadge>
                  <SaveJobButton :job-id="job.id" variant="icon" />
                </div>
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
                <span v-if="job.published_at" class="inline-flex items-center gap-1 ml-auto">
                  <Clock class="w-3.5 h-3.5" /> {{ timeAgo(job.published_at) }}
                </span>
              </div>

              <div v-if="job.tags?.length" class="flex flex-wrap gap-1.5 mt-3">
                <AppBadge v-for="t in job.tags" :key="t" variant="neutral" size="sm">
                  {{ t }}
                </AppBadge>
              </div>
            </div>
          </div>
        </AppCard>
      </RouterLink>

      <div v-if="canLoadMore" class="text-center pt-2">
        <AppButton variant="secondary" :loading="loading" @click="loadMore">
          {{ loading ? 'Loading…' : 'Load more jobs' }}
        </AppButton>
      </div>
    </div>

    <SaveSearchModal
      :open="showSaveSearch"
      :query="q.trim()"
      :location="location.trim()"
      :remote-type="remoteType"
      @close="showSaveSearch = false"
      @created="showSaveSearch = false"
    />
  </div>
</template>
