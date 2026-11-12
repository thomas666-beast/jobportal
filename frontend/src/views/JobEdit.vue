<script setup>
import { ref, onMounted, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ArrowLeft } from 'lucide-vue-next'
import JobForm from '../components/JobForm.vue'
import { getJob, updateJob, submitJob, closeJob, archiveJob } from '../api/jobs'
import { statusLabel, statusVariant } from '../utils/jobStatus'
import AppButton from '../components/ui/AppButton.vue'
import AppBadge from '../components/ui/AppBadge.vue'
import AppCard from '../components/ui/AppCard.vue'
import { useToast } from '../composables/useToast'

const route = useRoute()
const router = useRouter()
const toast = useToast()

const job = ref(null)
const loading = ref(true)
const saving = ref(false)
const acting = ref(false)
const errorMsg = ref('')

const editable = computed(
  () => job.value && (job.value.status === 'draft' || job.value.status === 'pending_review')
)

async function load() {
  loading.value = true
  errorMsg.value = ''
  try {
    job.value = await getJob(route.params.id)
  } catch (err) {
    errorMsg.value = err.message
  } finally {
    loading.value = false
  }
}

async function handleSubmit(payload) {
  errorMsg.value = ''
  saving.value = true
  try {
    job.value = await updateJob(job.value.id, payload)
    toast.success('Saved.')
  } catch (err) {
    errorMsg.value = err.message
  } finally {
    saving.value = false
  }
}

async function doSubmit() {
  acting.value = true
  try {
    job.value = await submitJob(job.value.id)
    toast.success('Submitted for review.')
  } catch (err) {
    toast.error(err.message)
  } finally {
    acting.value = false
  }
}

async function doClose() {
  acting.value = true
  try {
    job.value = await closeJob(job.value.id)
    toast.success('Job closed.')
  } catch (err) {
    toast.error(err.message)
  } finally {
    acting.value = false
  }
}

async function doArchive() {
  if (!confirm('Archive this job?')) return
  acting.value = true
  try {
    job.value = await archiveJob(job.value.id)
    toast.success('Job archived.')
    setTimeout(() => router.push({ name: 'jobs' }), 500)
  } catch (err) {
    toast.error(err.message)
  } finally {
    acting.value = false
  }
}

onMounted(load)
</script>

<template>
  <div class="max-w-2xl mx-auto space-y-6">
    <button
      type="button"
      @click="router.push({ name: 'jobs' })"
      class="inline-flex items-center gap-1 text-sm text-slate-600 hover:text-slate-900"
    >
      <ArrowLeft class="w-4 h-4" /> Back to jobs
    </button>

    <p v-if="loading" class="text-slate-600">Loading…</p>
    <p v-else-if="errorMsg && !job" class="text-red-600">{{ errorMsg }}</p>

    <template v-else-if="job">
      <div class="flex items-start justify-between gap-4">
        <div>
          <h1 class="text-2xl font-semibold text-slate-900">
            {{ job.title || 'Untitled' }}
          </h1>
          <div class="flex items-center gap-2 mt-2">
            <AppBadge :variant="statusVariant(job.status)">
              {{ statusLabel(job.status) }}
            </AppBadge>
            <span v-if="job.rejection_reason" class="text-xs text-red-600">
              Rejected: {{ job.rejection_reason }}
            </span>
          </div>
        </div>

        <div class="flex flex-wrap gap-2">
          <AppButton
            v-if="job.status === 'draft'"
            variant="secondary"
            :disabled="acting"
            @click="doSubmit"
          >
            Submit for review
          </AppButton>
          <AppButton
            v-if="job.status === 'published'"
            variant="secondary"
            :disabled="acting"
            @click="doClose"
          >
            Close
          </AppButton>
          <AppButton
            v-if="job.status !== 'archived'"
            variant="ghost"
            :disabled="acting"
            @click="doArchive"
          >
            Archive
          </AppButton>
        </div>
      </div>

      <p v-if="!editable" class="text-sm text-slate-500">
        This job is {{ statusLabel(job.status).toLowerCase() }} and can no longer be edited.
      </p>

      <p v-if="errorMsg" class="text-sm text-red-600">{{ errorMsg }}</p>

      <AppCard>
        <JobForm
          :initial="job"
          :disabled="!editable"
          :saving="saving"
          submit-label="Save changes"
          @submit="handleSubmit"
        />
      </AppCard>
    </template>
  </div>
</template>
