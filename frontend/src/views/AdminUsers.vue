<script setup>
import { ref, onMounted, watch, computed } from 'vue'
import { RouterLink, useRoute, useRouter } from 'vue-router'
import { Users as UsersIcon } from 'lucide-vue-next'
import { listAdminUsers } from '../api/admin'
import { timeAgo } from '../utils/jobStatus'
import AdminNav from '../components/AdminNav.vue'
import AppBadge from '../components/ui/AppBadge.vue'
import AppCard from '../components/ui/AppCard.vue'
import AppSelect from '../components/ui/AppSelect.vue'
import EmptyState from '../components/ui/EmptyState.vue'
import AppSkeleton from '../components/ui/AppSkeleton.vue'
import AppPagination from '../components/ui/AppPagination.vue'

const route = useRoute()
const router = useRouter()

const users = ref([])
const total = ref(0)
const loading = ref(false)
const errorMsg = ref('')
const roleFilter = ref('')
const page = ref(1)
const pageSize = 50

// Init from URL
onMounted(() => {
  if (typeof route.query.role === 'string') roleFilter.value = route.query.role
  if (typeof route.query.page === 'string') {
    const n = parseInt(route.query.page, 10)
    if (n > 0) page.value = n
  }
  load()
})

async function load() {
  loading.value = true
  errorMsg.value = ''
  try {
    const res = await listAdminUsers({
      role: roleFilter.value,
      limit: pageSize,
      offset: (page.value - 1) * pageSize
    })
    users.value = res.users || []
    total.value = res.total || 0
  } catch (err) {
    errorMsg.value = err.message
  } finally {
    loading.value = false
  }
}

function syncUrl() {
  const q = {}
  if (roleFilter.value) q.role = roleFilter.value
  if (page.value > 1) q.page = String(page.value)
  router.replace({ query: q })
}

watch(roleFilter, () => {
  page.value = 1
  syncUrl()
  load()
})

watch(page, () => {
  syncUrl()
  load()
})

function roleLabel(r) {
  if (r === 'job_seeker') return 'Job seeker'
  if (r === 'employer') return 'Employer'
  if (r === 'admin') return 'Admin'
  return r
}

function roleVariant(r) {
  if (r === 'job_seeker') return 'primary'
  if (r === 'employer') return 'success'
  if (r === 'admin') return 'warning'
  return 'neutral'
}
</script>

<template>
  <div class="space-y-6">
    <AdminNav />

    <div class="flex items-center justify-between">
      <h1 class="text-2xl font-semibold text-slate-900">Users</h1>
      <div class="w-48">
        <AppSelect v-model="roleFilter">
          <option value="">All roles</option>
          <option value="job_seeker">Job seekers</option>
          <option value="employer">Employers</option>
          <option value="admin">Admins</option>
        </AppSelect>
      </div>
    </div>

    <div v-if="loading && users.length === 0" class="space-y-2">
      <AppSkeleton v-for="i in 5" :key="i" variant="row" />
    </div>
    <p v-else-if="errorMsg" class="text-red-600">{{ errorMsg }}</p>

    <template v-else>
      <AppCard v-if="users.length === 0" :padded="false">
        <EmptyState :icon="UsersIcon" title="No users match" />
      </AppCard>

      <div v-else class="space-y-2">
        <AppCard v-for="u in users" :key="u.id" :padded="false">
          <RouterLink
            :to="{ name: 'admin-user-detail', params: { id: u.id } }"
            class="block p-4"
          >
            <div class="flex items-center justify-between gap-4">
              <div class="min-w-0 flex-1">
                <p class="font-medium text-slate-900">{{ u.username }}</p>
                <p class="text-xs text-slate-500">{{ u.email }}</p>
                <p v-if="u.company_name" class="text-xs text-slate-600 mt-1">
                  🏢 {{ u.company_name }}
                </p>
                <p v-else-if="u.seeker_full_name" class="text-xs text-slate-600 mt-1">
                  {{ u.seeker_full_name }}
                  <span v-if="u.seeker_headline"> · {{ u.seeker_headline }}</span>
                </p>
              </div>
              <div class="flex items-center gap-3 flex-shrink-0">
                <AppBadge :variant="roleVariant(u.role)" size="sm">
                  {{ roleLabel(u.role) }}
                </AppBadge>
                <span class="text-xs text-slate-400 w-16 text-right">
                  {{ timeAgo(u.created_at) }}
                </span>
              </div>
            </div>
          </RouterLink>
        </AppCard>

        <AppPagination
          v-model:page="page"
          :page-size="pageSize"
          :total="total"
          :loading="loading"
        />
      </div>
    </template>
  </div>
</template>
