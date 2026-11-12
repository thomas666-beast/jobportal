<script setup>
import { ref, onMounted } from 'vue'
import { RouterLink } from 'vue-router'
import { Inbox, MapPin } from 'lucide-vue-next'
import { listMyApplications, withdrawApplication } from '../api/seekerApplications'
import { applicationStatusLabel, applicationStatusVariant, timeAgo, remoteLabel } from '../utils/jobStatus'
import AppButton from '../components/ui/AppButton.vue'
import AppBadge from '../components/ui/AppBadge.vue'
import AppCard from '../components/ui/AppCard.vue'
import EmptyState from '../components/ui/EmptyState.vue'
import AppSkeleton from '../components/ui/AppSkeleton.vue'
import { useToast } from '../composables/useToast'

const toast = useToast()
const apps = ref([])
const loading = ref(true)
const errorMsg = ref('')
const acting = ref(null)

async function load() {
  loading.value = true
  errorMsg.value = ''
  try { apps.value = await listMyApplications() } 
  catch (err) { errorMsg.value = err.message } 
  finally { loading.value = false }
}

async function doWithdraw(id) {
  if (!confirm('Withdraw this application?')) return
  acting.value = id
  try {
    await withdrawApplication(id)
    toast.success('Application withdrawn.')
    await load()
  } catch (err) { toast.error(err.message) } 
  finally { acting.value = null }
}

onMounted(load)
</script>

<template>
  <div class="space-y-6">
    <div>
      <h1 class="text-3xl font-bold text-fg">My applications</h1>
      <p class="text-sm text-fg-muted mt-2">Track the jobs you've applied to</p>
    </div>

    <div v-if="loading" class="space-y-4"><AppSkeleton v-for="i in 3" :key="i" variant="row" /></div>
    <p v-else-if="errorMsg" class="text-red-500">{{ errorMsg }}</p>

    <AppCard v-else-if="apps.length === 0" :padded="false">
      <EmptyState :icon="Inbox" title="No applications yet" description="Browse jobs and apply to get started.">
        <template #action><AppButton to="/browse">Browse jobs</AppButton></template>
      </EmptyState>
    </AppCard>

    <div v-else class="space-y-4">
      <AppCard v-for="a in apps" :key="a.id">
        <div class="flex flex-col md:flex-row md:items-start md:justify-between gap-4">
          <div class="min-w-0 flex-1">
            <RouterLink :to="{ name: 'job-detail', params: { id: a.job_id } }" class="text-xl font-bold text-fg hover:text-primary transition-colors">
              {{ a.job.title }}
            </RouterLink>
            <p class="text-sm text-fg-muted mt-1">{{ a.company.name }}</p>

            <div class="flex flex-wrap gap-4 mt-3 text-xs text-fg-subtle">
              <span v-if="a.job.location" class="inline-flex items-center gap-1"><MapPin class="w-3.5 h-3.5" /> {{ a.job.location }}</span>
              <span>{{ remoteLabel(a.job.remote_type) }}</span>
              <span>Applied {{ timeAgo(a.created_at) }}</span>
            </div>
          </div>

          <div class="flex md:flex-col items-center md:items-end gap-3 md:flex-shrink-0">
            <AppBadge :variant="applicationStatusVariant(a.status)">{{ applicationStatusLabel(a.status) }}</AppBadge>
            <AppButton v-if="a.status === 'submitted' || a.status === 'viewed'" size="sm" variant="ghost" :disabled="acting === a.id" @click="doWithdraw(a.id)">Withdraw</AppButton>
          </div>
        </div>
      </AppCard>
    </div>
  </div>
</template>
