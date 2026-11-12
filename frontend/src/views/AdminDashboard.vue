<script setup>
import { ref, onMounted } from 'vue'
import { RouterLink } from 'vue-router'
import { Users, Building2, Briefcase, FileText } from 'lucide-vue-next'
import { getAdminStats } from '../api/admin'
import {
  statusLabel,
  statusVariant,
  applicationStatusLabel,
  applicationStatusVariant,
  timeAgo
} from '../utils/jobStatus'
import AdminNav from '../components/AdminNav.vue'
import AppCard from '../components/ui/AppCard.vue'
import AppBadge from '../components/ui/AppBadge.vue'

const stats = ref(null)
const loading = ref(true)
const errorMsg = ref('')

async function load() {
  loading.value = true
  errorMsg.value = ''
  try {
    stats.value = await getAdminStats()
  } catch (err) {
    errorMsg.value = err.message
  } finally {
    loading.value = false
  }
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
  <div class="space-y-6">
    <AdminNav />

    <h1 class="text-2xl font-semibold text-slate-900">Admin dashboard</h1>

    <p v-if="loading" class="text-slate-600">Loading…</p>
    <p v-else-if="errorMsg" class="text-red-600">{{ errorMsg }}</p>

    <template v-else-if="stats">
      <div class="grid grid-cols-2 md:grid-cols-4 gap-4">
        <AppCard>
          <div class="flex items-center gap-2 text-slate-500">
            <Users class="w-4 h-4" />
            <span class="text-sm">Users</span>
          </div>
          <div class="text-3xl font-semibold text-slate-900 mt-2">
            {{ stats.totals.users }}
          </div>
        </AppCard>
        <AppCard>
          <div class="flex items-center gap-2 text-slate-500">
            <Building2 class="w-4 h-4" />
            <span class="text-sm">Companies</span>
          </div>
          <div class="text-3xl font-semibold text-slate-900 mt-2">
            {{ stats.totals.companies }}
          </div>
        </AppCard>
        <AppCard>
          <div class="flex items-center gap-2 text-slate-500">
            <Briefcase class="w-4 h-4" />
            <span class="text-sm">Jobs</span>
          </div>
          <div class="text-3xl font-semibold text-slate-900 mt-2">
            {{ stats.totals.jobs }}
          </div>
        </AppCard>
        <AppCard>
          <div class="flex items-center gap-2 text-slate-500">
            <FileText class="w-4 h-4" />
            <span class="text-sm">Applications</span>
          </div>
          <div class="text-3xl font-semibold text-slate-900 mt-2">
            {{ stats.totals.applications }}
          </div>
        </AppCard>
      </div>

      <AppCard>
        <h2 class="text-sm font-medium text-slate-500 mb-3">Users by role</h2>
        <div class="flex flex-wrap gap-3">
          <div
            v-for="(count, role) in stats.users_by_role"
            :key="role"
            class="bg-slate-50 rounded-lg px-4 py-3"
          >
            <div class="text-xs text-slate-500">{{ roleLabel(role) }}</div>
            <div class="text-xl font-semibold text-slate-900">{{ count }}</div>
          </div>
        </div>
      </AppCard>

      <AppCard>
        <h2 class="text-sm font-medium text-slate-500 mb-3">Jobs by status</h2>
        <div class="flex flex-wrap gap-3">
          <div
            v-for="(count, status) in stats.jobs_by_status"
            :key="status"
            class="rounded-lg px-4 py-3"
          >
            <AppBadge :variant="statusVariant(status)">
              {{ statusLabel(status) }} · {{ count }}
            </AppBadge>
          </div>
        </div>
      </AppCard>

      <AppCard>
        <h2 class="text-sm font-medium text-slate-500 mb-3">Applications by status</h2>
        <div class="flex flex-wrap gap-3">
          <div
            v-for="(count, status) in stats.applications_by_status"
            :key="status"
            class="rounded-lg px-4 py-3"
          >
            <AppBadge :variant="applicationStatusVariant(status)">
              {{ applicationStatusLabel(status) }} · {{ count }}
            </AppBadge>
          </div>
        </div>
      </AppCard>

      <AppCard>
        <div class="flex items-center justify-between mb-3">
          <h2 class="text-sm font-medium text-slate-500">Recent signups</h2>
          <RouterLink
            to="/admin/users"
            class="text-xs text-primary-600 hover:underline"
          >
            View all →
          </RouterLink>
        </div>
        <p v-if="!stats.recent_users?.length" class="text-sm text-slate-500">
          No users yet.
        </p>
        <table v-else class="w-full text-sm">
          <tbody>
            <tr
              v-for="u in stats.recent_users"
              :key="u.id"
              class="border-t border-slate-100 first:border-t-0"
            >
              <td class="py-2">
                <RouterLink
                  :to="{ name: 'admin-user-detail', params: { id: u.id } }"
                  class="font-medium text-slate-900 hover:text-primary-600"
                >
                  {{ u.username }}
                </RouterLink>
                <div class="text-xs text-slate-500">{{ u.email }}</div>
              </td>
              <td class="py-2 text-right">
                <AppBadge variant="neutral" size="sm">{{ roleLabel(u.role) }}</AppBadge>
              </td>
              <td class="py-2 text-right text-xs text-slate-500 w-24">
                {{ timeAgo(u.created_at) }}
              </td>
            </tr>
          </tbody>
        </table>
      </AppCard>

      <AppCard>
        <h2 class="text-sm font-medium text-slate-500 mb-3">Recent applications</h2>
        <p v-if="!stats.recent_applications?.length" class="text-sm text-slate-500">
          No applications yet.
        </p>
        <table v-else class="w-full text-sm">
          <tbody>
            <tr
              v-for="a in stats.recent_applications"
              :key="a.id"
              class="border-t border-slate-100 first:border-t-0"
            >
              <td class="py-2">
                <span class="font-medium text-slate-900">{{ a.seeker_username }}</span>
                <span class="text-slate-500"> → </span>
                <span class="text-slate-700">{{ a.job_title }}</span>
                <div class="text-xs text-slate-500">at {{ a.company_name }}</div>
              </td>
              <td class="py-2 text-right">
                <AppBadge :variant="applicationStatusVariant(a.status)" size="sm">
                  {{ applicationStatusLabel(a.status) }}
                </AppBadge>
              </td>
              <td class="py-2 text-right text-xs text-slate-500 w-24">
                {{ timeAgo(a.created_at) }}
              </td>
            </tr>
          </tbody>
        </table>
      </AppCard>
    </template>
  </div>
</template>
