<script setup>
import { ref, onMounted, watch } from 'vue'
import { RouterLink, useRoute, useRouter } from 'vue-router'
import { Inbox } from 'lucide-vue-next'
import { listAdminJobs, approveJob, rejectJob } from '../api/admin'
import { statusLabel, statusVariant, timeAgo } from '../utils/jobStatus'
import AdminNav from '../components/AdminNav.vue'
import AppButton from '../components/ui/AppButton.vue'
import AppBadge from '../components/ui/AppBadge.vue'
import AppCard from '../components/ui/AppCard.vue'
import AppSelect from '../components/ui/AppSelect.vue'
import EmptyState from '../components/ui/EmptyState.vue'
import AppSkeleton from '../components/ui/AppSkeleton.vue'
import AppPagination from '../components/ui/AppPagination.vue'
import { useToast } from '../composables/useToast'

const route = useRoute()
const router = useRouter()
const toast = useToast()

const statusFilter = ref('pending_review')
const jobs = ref([])
const total = ref(0)
const loading = ref(false)
const errorMsg = ref('')
const acting = ref(null)
const rejecting = ref(null)
const rejectReason = ref('')
const rejectError = ref('')
const page = ref(1)
const pageSize = 50

onMounted(() => {
  if (typeof route.query.status === 'string') statusFilter.value = route.query.status
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
    const res = await listAdminJobs({
      status: statusFilter.value,
      limit: pageSize,
      offset: (page.value - 1) * pageSize
    })
    jobs.value = res.jobs || []
    total.value = res.total || 0
  } catch (err) {
    errorMsg.value = err.message
  } finally {
    loading.value = false
  }
}

function syncUrl() {
  const q = {}
  if (statusFilter.value && statusFilter.value !== 'pending_review') q.status = statusFilter.value
  if (page.value > 1) q.page = String(page.value)
  router.replace({ query: q })
}

watch(statusFilter, () => {
  page.value = 1
  syncUrl()
  load()
})

watch(page, () => {
  syncUrl()
  load()
})

async function doApprove(id) {
  acting.value = id
  try {
    await approveJob(id)
    toast.success('Job approved and published.')
    await load()
  } catch (err) {
    toast.error(err.message)
  } finally {
    acting.value = null
  }
}

function startReject(id) {
  rejecting.value = id
  rejectReason.value = ''
  rejectError.value = ''
}

function cancelReject() {
  rejecting.value = null
  rejectReason.value = ''
  rejectError.value = ''
}

async function submitReject() {
  rejectError.value = ''
  if (!rejectReason.value.trim()) {
    rejectError.value = 'Please provide a reason.'
    return
  }
  acting.value = rejecting.value
  try {
    await rejectJob(rejecting.value, rejectReason.value.trim())
    toast.success('Job rejected.')
    cancelReject()
    await load()
  } catch (err) {
    rejectError.value = err.message
  } finally {
    acting.value = null
  }
}
</script>

<template>
  <div class="space-y-6">
    <AdminNav />

    <div class="flex items-center justify-between">
      <h1 class="text-2xl font-semibold text-slate-900">Job moderation</h1>
      <div class="w-56">
        <AppSelect v-model="statusFilter">
          <option value="pending_review">Pending review</option>
          <option value="published">Published</option>
          <option value="draft">Draft</option>
          <option value="closed">Closed</option>
          <option value="archived">Archived</option>
        </AppSelect>
      </div>
    </div>

    <div v-if="loading && jobs.length === 0" class="space-y-3">
      <AppSkeleton v-for="i in 3" :key="i" variant="row" />
    </div>
    <p v-else-if="errorMsg" class="text-red-600">{{ errorMsg }}</p>

    <AppCard v-else-if="jobs.length === 0" :padded="false">
      <EmptyState
        :icon="Inbox"
        :title="`No jobs with status &quot;${statusLabel(statusFilter)}&quot;`"
        description="Nothing to moderate right now."
      />
    </AppCard>

    <template v-else>
      <div class="space-y-3">
        <AppCard v-for="job in jobs" :key="job.id" :padded="false">
          <div class="p-5">
            <div class="flex items-start justify-between gap-4">
              <div class="min-w-0 flex-1">
                <RouterLink
                  :to="{ name: 'admin-job-detail', params: { id: job.id } }"
                  class="text-base font-semibold text-slate-900 hover:text-primary-600"
                >
                  {{ job.title }}
                </RouterLink>
                <p class="text-sm text-slate-600 mt-0.5 truncate">
                  {{ job.description }}
                </p>
                <div class="flex items-center gap-3 mt-2 text-xs text-slate-500">
                  <AppBadge :variant="statusVariant(job.status)" size="sm">
                    {{ statusLabel(job.status) }}
                  </AppBadge>
                  <span v-if="job.location">{{ job.location }}</span>
                  <span>{{ timeAgo(job.created_at) }}</span>
                </div>
              </div>

              <div class="flex gap-2 flex-shrink-0">
                <AppButton
                  v-if="job.status === 'pending_review'"
                  size="sm"
                  variant="success"
                  :disabled="acting === job.id"
                  @click="doApprove(job.id)"
                >
                  Approve
                </AppButton>
                <AppButton
                  v-if="job.status === 'pending_review'"
                  size="sm"
                  variant="danger"
                  :disabled="acting === job.id"
                  @click="startReject(job.id)"
                >
                  Reject
                </AppButton>
              </div>
            </div>
          </div>

          <div
            v-if="rejecting === job.id"
            class="border-t border-slate-100 bg-red-50 p-5 space-y-2"
          >
            <label class="block text-sm font-medium text-red-900">
              Rejection reason (visible to employer)
            </label>
            <textarea
              v-model="rejectReason"
              rows="2"
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
                :disabled="acting === job.id"
                @click="submitReject"
              >
                Confirm rejection
              </AppButton>
            </div>
          </div>
        </AppCard>
      </div>

      <AppPagination
        v-model:page="page"
        :page-size="pageSize"
        :total="total"
        :loading="loading"
      />
    </template>
  </div>
</template>
