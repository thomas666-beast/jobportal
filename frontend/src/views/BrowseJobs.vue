<script setup>
import { ref, onMounted, computed } from 'vue'
import { RouterLink } from 'vue-router'
import { Search, MapPin, Briefcase, Banknote, Building2, X } from 'lucide-vue-next'
import { browseJobs } from '../api/publicJobs'
import { timeAgo, formatSalary, remoteLabel, employmentLabel } from '../utils/jobStatus'
import AppButton from '../components/ui/AppButton.vue'
import AppBadge from '../components/ui/AppBadge.vue'
import AppCard from '../components/ui/AppCard.vue'
import EmptyState from '../components/ui/EmptyState.vue'
import JobCardSkeleton from '../components/ui/JobCardSkeleton.vue'

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

let debounceTimer = null

async function load(reset = true) {
  loading.value = true
  errorMsg.value = ''
  try {
    const res = await browseJobs({ q: q.value.trim(), location: location.value.trim(), remote_type: remoteType.value, limit, offset: reset ? 0 : offset.value })
    if (reset) { jobs.value = res.jobs || []; offset.value = 0 } else { jobs.value = [...jobs.value, ...(res.jobs || [])] }
    total.value = res.total || 0
  } catch (err) { errorMsg.value = err.message } finally { loading.value = false; initialLoading.value = false }
}

function debouncedLoad() { clearTimeout(debounceTimer); debounceTimer = setTimeout(() => load(true), 300) }
function loadMore() { offset.value = jobs.value.length; load(false) }
function clearFilters() { q.value = ''; location.value = ''; remoteType.value = ''; load(true) }

onMounted(() => load(true))
</script>

<template>
  <div class="space-y-6">
    <section class="text-center pt-6 pb-2">
      <h1 class="text-4xl font-bold text-fg">Find your next role</h1>
      <p class="mt-2 text-fg-muted">{{ total }} {{ total === 1 ? 'job' : 'jobs' }} available</p>
    </section>

    <AppCard>
      <div class="flex flex-col md:flex-row gap-4">
        <div class="flex-1 relative">
          <Search class="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-fg-subtle" />
          <input v-model="q" @input="debouncedLoad" type="text" placeholder="Search jobs..." class="w-full pl-10 pr-3 py-3 text-sm" />
        </div>
        <div class="flex-1 relative">
          <MapPin class="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-fg-subtle" />
          <input v-model="location" @input="debouncedLoad" type="text" placeholder="Location" class="w-full pl-10 pr-3 py-3 text-sm" />
        </div>
        <select v-model="remoteType" @change="load(true)" class="px-3 py-3 text-sm">
          <option value="">Any remote type</option>
          <option value="remote">Remote</option>
          <option value="hybrid">Hybrid</option>
          <option value="on_site">On-site</option>
        </select>
        <AppButton v-if="hasFilters" variant="ghost" @click="clearFilters"><X class="w-4 h-4" /> Clear</AppButton>
      </div>
    </AppCard>

    <div v-if="initialLoading" class="space-y-4"><JobCardSkeleton v-for="i in 3" :key="i" /></div>
    <p v-else-if="errorMsg" class="text-red-500 text-center">{{ errorMsg }}</p>
    
    <AppCard v-else-if="jobs.length === 0" :padded="false">
      <EmptyState :icon="Search" title="No jobs found" description="Try different keywords." />
    </AppCard>

    <div v-else class="space-y-4">
      <RouterLink v-for="job in jobs" :key="job.id" :to="{ name: 'job-detail', params: { id: job.id } }" class="block">
        <AppCard hover>
          <div class="flex gap-4">
            <div class="flex-shrink-0">
              <div v-if="job.company.logo_url" class="w-12 h-12 border border-border overflow-hidden bg-surface-muted"><img :src="job.company.logo_url" class="w-full h-full object-cover" /></div>
              <div v-else class="w-12 h-12 border border-border bg-surface-muted flex items-center justify-center"><Building2 class="w-6 h-6 text-fg-muted" /></div>
            </div>
            <div class="flex-1 min-w-0">
              <div class="flex items-start justify-between gap-3">
                <div>
                  <h2 class="text-xl font-bold text-fg">{{ job.title }}</h2>
                  <p class="text-sm text-fg-muted mt-1">{{ job.company.name }} <span v-if="job.company.location">· {{ job.company.location }}</span></p>
                </div>
                <AppBadge variant="neutral">{{ remoteLabel(job.remote_type) }}</AppBadge>
              </div>
              <p class="text-sm text-fg-muted mt-3 line-clamp-2">{{ job.description }}</p>
              <div class="flex flex-wrap gap-4 mt-4 text-xs text-fg-subtle">
                <span v-if="job.location" class="inline-flex items-center gap-1"><MapPin class="w-3.5 h-3.5" /> {{ job.location }}</span>
                <span class="inline-flex items-center gap-1"><Briefcase class="w-3.5 h-3.5" /> {{ employmentLabel(job.employment_type) }}</span>
                <span v-if="formatSalary(job.salary_min, job.salary_max, job.salary_currency)" class="inline-flex items-center gap-1"><Banknote class="w-3.5 h-3.5" /> {{ formatSalary(job.salary_min, job.salary_max, job.salary_currency) }}</span>
              </div>
            </div>
          </div>
        </AppCard>
      </RouterLink>

      <div v-if="canLoadMore" class="text-center pt-6">
        <AppButton variant="secondary" :loading="loading" @click="loadMore">
          {{ loading ? 'Loading…' : 'Load more jobs' }}
        </AppButton>
      </div>
    </div>
  </div>
</template>
