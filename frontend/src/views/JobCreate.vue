<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { ArrowLeft } from 'lucide-vue-next'
import JobForm from '../components/JobForm.vue'
import { createJob } from '../api/jobs'
import AppCard from '../components/ui/AppCard.vue'
import { useToast } from '../composables/useToast'

const router = useRouter()
const toast = useToast()

const saving = ref(false)
const errorMsg = ref('')

async function handleSubmit(payload) {
  errorMsg.value = ''
  saving.value = true
  try {
    await createJob(payload)
    toast.success('Job created as draft.')
    router.push({ name: 'jobs' })
  } catch (err) {
    errorMsg.value = err.message
  } finally {
    saving.value = false
  }
}
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

    <div>
      <h1 class="text-2xl font-semibold text-slate-900">New job</h1>
      <p class="text-sm text-slate-600 mt-1">
        Save as draft first. Submit for review when ready.
      </p>
    </div>

    <p v-if="errorMsg" class="text-sm text-red-600">{{ errorMsg }}</p>

    <AppCard>
      <JobForm :saving="saving" submit-label="Create draft" @submit="handleSubmit" />
    </AppCard>
  </div>
</template>
