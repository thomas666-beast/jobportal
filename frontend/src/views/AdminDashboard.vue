<script setup>
import { ref, onMounted } from 'vue'
import { RouterLink } from 'vue-router'
import { Users, Building2, Briefcase, FileText } from 'lucide-vue-next'
import { getAdminStats } from '../api/admin'
import { statusLabel, statusVariant, applicationStatusLabel, applicationStatusVariant, timeAgo } from '../utils/jobStatus'
import AdminNav from '../components/AdminNav.vue'
import AppCard from '../components/ui/AppCard.vue'
import AppBadge from '../components/ui/AppBadge.vue'

const stats = ref(null)
const loading = ref(true)
const errorMsg = ref('')

async function load() {
  loading.value = true
  errorMsg.value = ''
  try { stats.value = await getAdminStats() } 
  catch (err) { errorMsg.value = err.message } 
  finally { loading.value = false }
}

onMounted(load)

function roleLabel(r) {
  if (r === 'job_seeker') return 'Job seekers'
  if (r === 'employer') return 'Employers'
  if (r === 'admin') return 'Admins'
  return r
}
</script>

<template>
  <div class="space-y-8">
    <AdminNav />
    <h1 class="text-3xl font-bold text-fg">Admin dashboard</h1>

    <p v-if="loading" class="text-fg-muted">Loading…</p>
    <p v-else-if="errorMsg" class="text-red-500">{{ errorMsg }}</p>

    <template v-else-if="stats">
      <div class="grid grid-cols-2 md:grid-cols-4 gap-4">
        <AppCard>
          <div class="flex items-center gap-2 text-fg-muted"><Users class="w-4 h-4" /><span class="text-sm font-bold uppercase tracking-wider">Users</span></div>
          <div class="text-4xl font-bold text-fg mt-3">{{ stats.totals.users }}</div>
        </AppCard>
        <AppCard>
          <div class="flex items-center gap-2 text-fg-muted"><Building2 class="w-4 h-4" /><span class="text-sm font-bold uppercase tracking-wider">Companies</span></div>
          <div class="text-4xl font-bold text-fg mt-3">{{ stats.totals.companies }}</div>
        </AppCard>
        <AppCard>
          <div class="flex items-center gap-2 text-fg-muted"><Briefcase class="w-4 h-4" /><span class="text-sm font-bold uppercase tracking-wider">Jobs</span></div>
          <div class="text-4xl font-bold text-fg mt-3">{{ stats.totals.jobs }}</div>
        </AppCard>
        <AppCard>
          <div class="flex items-center gap-2 text-fg-muted"><FileText class="w-4 h-4" /><span class="text-sm font-bold uppercase tracking-wider">Applications</span></div>
          <div class="text-4xl font-bold text-fg mt-3">{{ stats.totals.applications }}</div>
        </AppCard>
      </div>

      <div class="grid md:grid-cols-3 gap-6">
        <AppCard>
          <h2 class="text-sm font-bold text-fg-muted uppercase tracking-wider mb-4">Users by role</h2>
          <div class="flex flex-wrap gap-3">
            <div v-for="(count, role) in stats.users_by_role" :key="role" class="bg-surface-muted border border-border px-4 py-2 flex-1">
              <div class="text-xs text-fg-muted">{{ roleLabel(role) }}</div>
              <div class="text-xl font-bold text-fg">{{ count }}</div>
            </div>
          </div>
        </AppCard>

        <AppCard>
          <h2 class="text-sm font-bold text-fg-muted uppercase tracking-wider mb-4">Jobs by status</h2>
          <div class="flex flex-wrap gap-2">
            <AppBadge v-for="(count, status) in stats.jobs_by_status" :key="status" :variant="statusVariant(status)">
              {{ statusLabel(status) }} · {{ count }}
            </AppBadge>
          </div>
        </AppCard>

        <AppCard>
          <h2 class="text-sm font-bold text-fg-muted uppercase tracking-wider mb-4">Apps by status</h2>
          <div class="flex flex-wrap gap-2">
            <AppBadge v-for="(count, status) in stats.applications_by_status" :key="status" :variant="applicationStatusVariant(status)">
              {{ applicationStatusLabel(status) }} · {{ count }}
            </AppBadge>
          </div>
        </AppCard>
      </div>
    </template>
  </div>
</template>
