<script setup>
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import {
  ArrowLeft,
  MapPin,
  Briefcase,
  Tag,
  Banknote,
  Globe,
  Building2
} from 'lucide-vue-next'
import { getAdminJob, approveJob, rejectJob } from '../api/admin'
import {
  statusLabel,
  statusVariant,
  remoteLabel,
  employmentLabel,
  formatSalary,
  timeAgo
} from '../utils/jobStatus'
import AdminNav from '../components/AdminNav.vue'
import AppButton from '../components/ui/AppButton.vue'
import AppBadge from '../components/ui/AppBadge.vue'
import AppCard from '../components/ui/AppCard.vue'
import { useToast } from '../composables/useToast'

const route = useRoute()
const router = useRouter()
const toast = useToast()

const job = ref(null)
const loading = ref(true)
const errorMsg = ref('')
const acting = ref(false)
const rejecting = ref(false)
const rejectReason = ref('')
const rejectError = ref('')

async function load() {
  loading.value = true
  errorMsg.value = ''
  try {
    job.value = await getAdminJob(route.params.id)
  } catch (err) {
    errorMsg.value = err.message
  } finally {
    loading.value = false
  }
}

async function doApprove() {
  acting.value = true
  try {
    await approveJob(job.value.id)
    toast.success('Job approved and published.')
    await load()
  } catch (err) {
    toast.error(err.message)
  } finally {
    acting.value = false
  }
}

function startReject() {
  rejecting.value = true
  rejectReason.value = ''
  rejectError.value = ''
}

function cancelReject() {
  rejecting.value = false
  rejectReason.value = ''
  rejectError.value = ''
}

async function submitReject() {
  rejectError.value = ''
  if (!rejectReason.value.trim()) {
    rejectError.value = 'Please provide a reason.'
    return
  }
  acting.value = true
  try {
    await rejectJob(job.value.id, rejectReason.value.trim())
    toast.success('Job rejected.')
    cancelReject()
    await load()
  } catch (err) {
    rejectError.value = err.message
  } finally {
    acting.value = false
  }
}

onMounted(load)
</script>

<template>
  <div class="space-y-6">
    <AdminNav />

    <button
      type="button"
      @click="router.push({ name: 'admin-jobs' })"
      class="inline-flex items-center gap-1 text-sm text-slate-600 hover:text-slate-900"
    >
      <ArrowLeft class="w-4 h-4" /> Back to moderation list
    </button>

    <p v-if="loading" class="text-slate-600">Loading…</p>
    <p v-else-if="errorMsg" class="text-red-600">{{ errorMsg }}</p>

    <template v-else-if="job">
      <AppCard>
        <div class="flex items-start justify-between gap-4">
          <div>
            <h1 class="text-2xl font-semibold text-slate-900">{{ job.title }}</h1>
            <div class="flex items-center gap-2 mt-2">
              <AppBadge :variant="statusVariant(job.status)">
                {{ statusLabel(job.status) }}
              </AppBadge>
              <span class="text-xs text-slate-500">
                Posted {{ timeAgo(job.created_at) }}
              </span>
            </div>
          </div>

          <div v-if="job.status === 'pending_review'" class="flex gap-2">
            <AppButton variant="success" :disabled="acting" @click="doApprove">
              Approve
            </AppButton>
            <AppButton variant="danger" :disabled="acting" @click="startReject">
              Reject
            </AppButton>
          </div>
        </div>

        <div v-if="rejecting" class="bg-red-50 border border-red-200 rounded-lg p-4 mt-4 space-y-2">
          <label class="block text-sm font-medium text-red-900">
            Rejection reason (visible to employer)
          </label>
          <textarea
            v-model="rejectReason"
            rows="3"
            maxlength="1000"
            class="w-full rounded-lg border border-red-300 px-3 py-2 text-sm
                   focus:outline-none focus:ring-2 focus:ring-red-500"
          ></textarea>
          <p v-if="rejectError" class="text-sm text-red-700">{{ rejectError }}</p>
          <div class="flex gap-2 justify-end">
            <AppButton variant="secondary" size="sm" @click="cancelReject">
              Cancel
            </AppButton>
            <AppButton
              variant="danger"
              size="sm"
              :disabled="acting"
              @click="submitReject"
            >
              Confirm rejection
            </AppButton>
          </div>
        </div>

        <p v-if="job.rejection_reason" class="text-sm text-red-700 mt-4">
          Previously rejected: {{ job.rejection_reason }}
        </p>

        <div class="flex flex-wrap gap-4 mt-5 text-sm text-slate-600">
          <span v-if="job.location" class="inline-flex items-center gap-1">
            <MapPin class="w-4 h-4" /> {{ job.location }}
          </span>
          <span class="inline-flex items-center gap-1">
            <Briefcase class="w-4 h-4" /> {{ employmentLabel(job.employment_type) }}
          </span>
          <span class="inline-flex items-center gap-1">
            {{ remoteLabel(job.remote_type) }}
          </span>
          <span v-if="job.category" class="inline-flex items-center gap-1">
            <Tag class="w-4 h-4" /> {{ job.category }}
          </span>
          <span
            v-if="formatSalary(job.salary_min, job.salary_max, job.salary_currency)"
            class="inline-flex items-center gap-1"
          >
            <Banknote class="w-4 h-4" />
            {{ formatSalary(job.salary_min, job.salary_max, job.salary_currency) }}
          </span>
        </div>

        <div v-if="job.tags?.length" class="flex flex-wrap gap-2 mt-4">
          <AppBadge v-for="t in job.tags" :key="t" variant="neutral" size="sm">
            {{ t }}
          </AppBadge>
        </div>
      </AppCard>

      <AppCard>
        <h2 class="text-sm font-medium text-slate-500 mb-2">Description</h2>
        <p class="text-slate-800 whitespace-pre-wrap">{{ job.description }}</p>
      </AppCard>

      <AppCard v-if="job.requirements">
        <h2 class="text-sm font-medium text-slate-500 mb-2">Requirements</h2>
        <p class="text-slate-800 whitespace-pre-wrap">{{ job.requirements }}</p>
      </AppCard>

      <AppCard>
        <div class="flex items-start gap-3">
          <div class="w-10 h-10 rounded-lg bg-slate-100 flex items-center justify-center flex-shrink-0">
            <Building2 class="w-5 h-5 text-slate-500" />
          </div>
          <div class="min-w-0 flex-1">
            <h2 class="text-base font-semibold text-slate-900">{{ job.company.name }}</h2>
            <div v-if="job.company.location" class="text-sm text-slate-600 mt-1">
              {{ job.company.location }}
            </div>
            <div v-if="job.company.website" class="text-sm mt-1">
              <a
                :href="job.company.website"
                target="_blank"
                rel="noopener"
                class="text-primary-600 hover:underline inline-flex items-center gap-1"
              >
                <Globe class="w-3.5 h-3.5" /> {{ job.company.website }}
              </a>
            </div>
            <p v-if="job.company.description" class="text-sm text-slate-700 whitespace-pre-wrap mt-3">
              {{ job.company.description }}
            </p>
          </div>
        </div>

        <div class="pt-4 mt-4 border-t border-slate-100">
          <h3 class="text-sm font-medium text-slate-500 mb-1">Posted by</h3>
          <div class="text-sm text-slate-800">
            {{ job.employer.username || '(deleted user)' }}
            <span v-if="job.employer.email" class="text-slate-500">
              — {{ job.employer.email }}
            </span>
          </div>
        </div>
      </AppCard>
    </template>
  </div>
</template>
