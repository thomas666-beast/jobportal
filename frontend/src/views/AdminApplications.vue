<script setup>
import { ref, onMounted, watch } from 'vue'
import { RouterLink, useRoute, useRouter } from 'vue-router'
import { FileText, Search } from 'lucide-vue-next'
import { listAdminApplications } from '../api/admin'
import { applicationStatusLabel, applicationStatusVariant, timeAgo } from '../utils/jobStatus'
import AdminNav from '../components/AdminNav.vue'
import AppBadge from '../components/ui/AppBadge.vue'
import AppCard from '../components/ui/AppCard.vue'
import AppSelect from '../components/ui/AppSelect.vue'
import EmptyState from '../components/ui/EmptyState.vue'

const route = useRoute()
const router = useRouter()

const applications = ref([])
const total = ref(0)
const loading = ref(true)
const errorMsg = ref('')
const statusFilter = ref('')
const seekerQuery = ref('')
const page = ref(1)
const pageSize = 50

let debounceTimer = null

onMounted(() => {
  if (typeof route.query.status === 'string') statusFilter.value = route.query.status
  if (typeof route.query.q === 'string') seekerQuery.value = route.query.q
  if (typeof route.query.page === 'string') { const n = parseInt(route.query.page, 10); if (n > 0) page.value = n }
  load()
})

async function load() {
  loading.value = true
  errorMsg.value = ''
  try {
    const res = await listAdminApplications({ status: statusFilter.value, seeker_q: seekerQuery.value.trim(), limit: pageSize, offset: (page.value - 1) * pageSize })
    applications.value = res.applications || []
    total.value = res.total || 0
  } catch (err) { errorMsg.value = err.message } 
  finally { loading.value = false }
}

function syncUrl() {
  const q = {}
  if (statusFilter.value) q.status = statusFilter.value
  if (seekerQuery.value.trim()) q.q = seekerQuery.value.trim()
  if (page.value > 1) q.page = String(page.value)
  router.replace({ query: q })
}

function debouncedLoad() {
  clearTimeout(debounceTimer)
  debounceTimer = setTimeout(() => { page.value = 1; syncUrl(); load() }, 300)
}

watch(statusFilter, () => { page.value = 1; syncUrl(); load() })
watch(page, () => { syncUrl(); load() })
</script>

<template>
  <div class="space-y-6">
    <AdminNav />
    <h1 class="text-3xl font-bold text-fg">Applications</h1>

    <AppCard>
      <div class="grid md:grid-cols-2 gap-4">
        <div class="relative">
          <Search class="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-fg-subtle" />
          <input v-model="seekerQuery" @input="debouncedLoad" type="text" placeholder="Search seeker..." class="w-full pl-10 pr-3 py-3 text-sm" />
        </div>
        <AppSelect v-model="statusFilter">
          <option value="">All statuses</option>
          <option value="submitted">Submitted</option>
          <option value="viewed">Viewed</option>
          <option value="shortlisted">Shortlisted</option>
          <option value="rejected">Rejected</option>
          <option value="hired">Hired</option>
          <option value="withdrawn">Withdrawn</option>
        </AppSelect>
      </div>
    </AppCard>

    <p v-if="loading" class="text-fg-muted">Loading…</p>
    <p v-else-if="errorMsg" class="text-red-500">{{ errorMsg }}</p>

    <AppCard v-else-if="applications.length === 0" :padded="false">
      <EmptyState :icon="FileText" title="No applications match" description="Try a different filter or search term." />
    </AppCard>

    <div v-else class="space-y-3">
      <AppCard v-for="a in applications" :key="a.id" :padded="false">
        <RouterLink :to="{ name: 'admin-application-detail', params: { id: a.id } }" class="block p-6 hover:bg-surface-muted transition-colors">
          <div class="flex items-start justify-between gap-4">
            <div class="min-w-0 flex-1">
              <p class="text-lg font-bold text-fg">{{ a.seeker.full_name || a.seeker.username }}</p>
              <p v-if="a.seeker.headline" class="text-sm text-fg-muted mt-1">{{ a.seeker.headline }}</p>
              <p class="text-xs text-fg-subtle mt-3 uppercase tracking-wider">Applied to <strong class="text-fg">{{ a.job_title }}</strong> at {{ a.company_name }} · {{ timeAgo(a.created_at) }}</p>
            </div>
            <AppBadge :variant="applicationStatusVariant(a.status)">{{ applicationStatusLabel(a.status) }}</AppBadge>
          </div>
        </RouterLink>
      </AppCard>
    </div>
  </div>
</template>
