<script setup>
import { ref, onMounted } from 'vue'
import { RouterLink } from 'vue-router'
import { Plus, Briefcase } from 'lucide-vue-next'
import { listJobs, submitJob, closeJob, archiveJob } from '../api/jobs'
import { statusLabel, statusVariant } from '../utils/jobStatus'
import AppButton from '../components/ui/AppButton.vue'
import AppBadge from '../components/ui/AppBadge.vue'
import AppCard from '../components/ui/AppCard.vue'
import EmptyState from '../components/ui/EmptyState.vue'
import AppSkeleton from '../components/ui/AppSkeleton.vue'
import { useToast } from '../composables/useToast'

const toast = useToast()
const jobs = ref([])
const loading = ref(true)
const errorMsg = ref('')
const acting = ref(null)

async function load() {
  loading.value = true
  errorMsg.value = ''
  try { jobs.value = await listJobs() } 
  catch (err) { errorMsg.value = err.message } 
  finally { loading.value = false }
}

function replaceJob(updated) {
  const i = jobs.value.findIndex((j) => j.id === updated.id)
  if (i !== -1) jobs.value[i] = updated
}

async function doSubmit(id) {
  acting.value = id
  try { replaceJob(await submitJob(id)); toast.success('Submitted for review.') } 
  catch (err) { toast.error(err.message) } 
  finally { acting.value = null }
}

async function doClose(id) {
  acting.value = id
  try { replaceJob(await closeJob(id)); toast.success('Job closed.') } 
  catch (err) { toast.error(err.message) } 
  finally { acting.value = null }
}

async function doArchive(id) {
  if (!confirm('Archive this job?')) return
  acting.value = id
  try { replaceJob(await archiveJob(id)); toast.success('Job archived.') } 
  catch (err) { toast.error(err.message) } 
  finally { acting.value = null }
}

onMounted(load)
</script>

<template>
  <div class="space-y-6">
    <div class="flex items-center justify-between">
      <div>
        <h1 class="text-3xl font-bold text-fg">Your jobs</h1>
        <p class="text-sm text-fg-muted mt-2">Manage your job listings</p>
      </div>
      <AppButton to="/jobs/new" variant="success"><Plus class="w-4 h-4" /> New job</AppButton>
    </div>

    <div v-if="loading" class="space-y-4"><AppSkeleton v-for="i in 3" :key="i" variant="row" /></div>
    <p v-else-if="errorMsg" class="text-red-500">{{ errorMsg }}</p>

    <AppCard v-else-if="jobs.length === 0" :padded="false">
      <EmptyState :icon="Briefcase" title="No jobs yet" description="Post your first job to start receiving applications.">
        <template #action><AppButton to="/jobs/new" variant="success">Create your first job</AppButton></template>
      </EmptyState>
    </AppCard>

    <div v-else class="space-y-4">
      <AppCard v-for="job in jobs" :key="job.id">
        <div class="flex flex-col md:flex-row md:items-start md:justify-between gap-4">
          <div class="min-w-0 flex-1">
            <RouterLink :to="{ name: 'job-edit', params: { id: job.id } }" class="text-xl font-bold text-fg hover:text-primary transition-colors">
              {{ job.title }}
            </RouterLink>
            <p class="text-sm text-fg-muted mt-1">{{ job.location || 'No location' }}</p>
          </div>

          <div class="flex flex-wrap items-center gap-3 flex-shrink-0">
            <AppBadge :variant="statusVariant(job.status)">{{ statusLabel(job.status) }}</AppBadge>
            <AppButton v-if="job.status === 'draft'" size="sm" variant="secondary" :disabled="acting === job.id" @click="doSubmit(job.id)">Submit for review</AppButton>
            <AppButton v-if="job.status === 'published'" size="sm" variant="secondary" :disabled="acting === job.id" @click="doClose(job.id)">Close</AppButton>
            <AppButton v-if="job.status !== 'archived'" size="sm" variant="ghost" :disabled="acting === job.id" @click="doArchive(job.id)">Archive</AppButton>
          </div>
        </div>
      </AppCard>
    </div>
  </div>
</template>
