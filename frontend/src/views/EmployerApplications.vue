<script setup>
import { ref, onMounted, computed } from 'vue'
import { RouterLink } from 'vue-router'
import { Inbox } from 'lucide-vue-next'
import { listCompanyApplications } from '../api/employerApplications'
import { applicationStatusLabel, applicationStatusVariant, timeAgo } from '../utils/jobStatus'
import AppBadge from '../components/ui/AppBadge.vue'
import AppCard from '../components/ui/AppCard.vue'
import AppSelect from '../components/ui/AppSelect.vue'
import EmptyState from '../components/ui/EmptyState.vue'

const apps = ref([])
const loading = ref(true)
const errorMsg = ref('')
const statusFilter = ref('')
const jobFilter = ref('')

const jobOptions = computed(() => {
  const set = new Map()
  for (const a of apps.value) { set.set(a.job_id, a.job.title) }
  return Array.from(set, ([id, title]) => ({ id, title }))
})

const filtered = computed(() => {
  return apps.value.filter((a) => {
    if (statusFilter.value && a.status !== statusFilter.value) return false
    if (jobFilter.value && a.job_id !== jobFilter.value) return false
    return true
  })
})

async function load() {
  loading.value = true
  errorMsg.value = ''
  try { apps.value = await listCompanyApplications() } 
  catch (err) { errorMsg.value = err.message } 
  finally { loading.value = false }
}

onMounted(load)
</script>

<template>
  <div class="space-y-6">
    <div>
      <h1 class="text-3xl font-bold text-fg">Applications</h1>
      <p class="text-sm text-fg-muted mt-2">Review candidates who applied to your jobs</p>
    </div>

    <AppCard>
      <div class="grid md:grid-cols-2 gap-4">
        <AppSelect v-model="statusFilter">
          <option value="">All statuses</option>
          <option value="submitted">Submitted</option>
          <option value="viewed">Viewed</option>
          <option value="shortlisted">Shortlisted</option>
          <option value="rejected">Rejected</option>
          <option value="hired">Hired</option>
          <option value="withdrawn">Withdrawn</option>
        </AppSelect>
        <AppSelect v-model="jobFilter">
          <option value="">All jobs</option>
          <option v-for="j in jobOptions" :key="j.id" :value="j.id">{{ j.title }}</option>
        </AppSelect>
      </div>
    </AppCard>

    <p v-if="loading" class="text-fg-muted">Loading…</p>
    <p v-else-if="errorMsg" class="text-red-500">{{ errorMsg }}</p>

    <AppCard v-else-if="filtered.length === 0" :padded="false">
      <EmptyState :icon="Inbox" title="No applications match" :description="apps.length === 0 ? 'Applications will appear here once candidates apply.' : 'Try a different filter.'" />
    </AppCard>

    <div v-else class="space-y-4">
      <AppCard v-for="a in filtered" :key="a.id" :padded="false">
        <RouterLink :to="{ name: 'employer-application-detail', params: { id: a.id } }" class="block p-6 hover:bg-surface-muted transition-colors">
          <div class="flex items-start justify-between gap-4">
            <div class="min-w-0 flex-1">
              <p class="text-lg font-bold text-fg">{{ a.seeker.full_name || a.seeker.username }}</p>
              <p v-if="a.seeker.headline" class="text-sm text-fg-muted mt-1">{{ a.seeker.headline }}</p>
              <p class="text-xs text-fg-subtle mt-3 uppercase tracking-wider">
                Applied for <strong class="text-fg">{{ a.job.title }}</strong> · {{ timeAgo(a.created_at) }}
              </p>
            </div>
            <AppBadge :variant="applicationStatusVariant(a.status)">{{ applicationStatusLabel(a.status) }}</AppBadge>
          </div>
        </RouterLink>
      </AppCard>
    </div>
  </div>
</template>
