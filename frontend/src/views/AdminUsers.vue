<script setup>
import { ref, onMounted, watch } from 'vue'
import { RouterLink } from 'vue-router'
import { Users as UsersIcon } from 'lucide-vue-next'
import { listAdminUsers } from '../api/admin'
import { timeAgo } from '../utils/jobStatus'
import AdminNav from '../components/AdminNav.vue'
import AppBadge from '../components/ui/AppBadge.vue'
import AppCard from '../components/ui/AppCard.vue'
import AppSelect from '../components/ui/AppSelect.vue'
import EmptyState from '../components/ui/EmptyState.vue'

const users = ref([])
const total = ref(0)
const loading = ref(false)
const errorMsg = ref('')
const roleFilter = ref('')

async function load() {
  loading.value = true
  errorMsg.value = ''
  try {
    const res = await listAdminUsers({ role: roleFilter.value })
    users.value = res.users || []
    total.value = res.total || 0
  } catch (err) { errorMsg.value = err.message } 
  finally { loading.value = false }
}

watch(roleFilter, load)
onMounted(load)

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
      <h1 class="text-3xl font-bold text-fg">Users</h1>
      <div class="w-48">
        <AppSelect v-model="roleFilter">
          <option value="">All roles</option>
          <option value="job_seeker">Job seekers</option>
          <option value="employer">Employers</option>
          <option value="admin">Admins</option>
        </AppSelect>
      </div>
    </div>

    <p v-if="loading" class="text-fg-muted">Loading…</p>
    <p v-else-if="errorMsg" class="text-red-500">{{ errorMsg }}</p>

    <template v-else>
      <p class="text-xs text-fg-subtle uppercase tracking-wider">{{ total }} total</p>
      <AppCard v-if="users.length === 0" :padded="false"><EmptyState :icon="UsersIcon" title="No users match" /></AppCard>
      <div v-else class="space-y-3">
        <AppCard v-for="u in users" :key="u.id" :padded="false">
          <RouterLink :to="{ name: 'admin-user-detail', params: { id: u.id } }" class="block p-6 hover:bg-surface-muted transition-colors">
            <div class="flex items-center justify-between gap-4">
              <div class="min-w-0 flex-1">
                <p class="text-lg font-bold text-fg">{{ u.username }}</p>
                <p class="text-sm text-fg-muted">{{ u.email }}</p>
                <p v-if="u.company_name" class="text-xs text-fg-subtle mt-2 uppercase tracking-wider">🏢 {{ u.company_name }}</p>
                <p v-else-if="u.seeker_full_name" class="text-xs text-fg-subtle mt-2 uppercase tracking-wider">{{ u.seeker_full_name }} <span v-if="u.seeker_headline">· {{ u.seeker_headline }}</span></p>
              </div>
              <div class="flex items-center gap-4 flex-shrink-0">
                <AppBadge :variant="roleVariant(u.role)">{{ roleLabel(u.role) }}</AppBadge>
                <span class="text-xs text-fg-subtle w-20 text-right">{{ timeAgo(u.created_at) }}</span>
              </div>
            </div>
          </RouterLink>
        </AppCard>
      </div>
    </template>
  </div>
</template>
