<script setup>
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ArrowLeft, Mail, Phone, MapPin, FileText, Download, Building2, Briefcase } from 'lucide-vue-next'
import { getAdminApplication } from '../api/admin'
import {
  applicationStatusLabel,
  applicationStatusVariant,
  statusLabel,
  statusVariant,
  timeAgo
} from '../utils/jobStatus'
import AdminNav from '../components/AdminNav.vue'
import AppBadge from '../components/ui/AppBadge.vue'
import AppCard from '../components/ui/AppCard.vue'

const route = useRoute()
const router = useRouter()

const app = ref(null)
const loading = ref(true)
const errorMsg = ref('')

async function load() {
  loading.value = true
  errorMsg.value = ''
  try {
    app.value = await getAdminApplication(route.params.id)
  } catch (err) {
    errorMsg.value = err.message
  } finally {
    loading.value = false
  }
}

onMounted(load)
</script>

<template>
  <div class="space-y-6">
    <AdminNav />

    <button
      type="button"
      @click="router.push({ name: 'admin-applications' })"
      class="inline-flex items-center gap-1 text-sm text-slate-600 hover:text-slate-900"
    >
      <ArrowLeft class="w-4 h-4" /> Back to applications
    </button>

    <p v-if="loading" class="text-slate-600">Loading…</p>
    <p v-else-if="errorMsg" class="text-red-600">{{ errorMsg }}</p>

    <template v-else-if="app">
      <AppCard>
        <div class="flex items-start justify-between gap-4">
          <div>
            <h1 class="text-2xl font-semibold text-slate-900">
              {{ app.seeker.full_name || app.seeker.username }}
            </h1>
            <p v-if="app.seeker.headline" class="text-sm text-slate-600 mt-1">
              {{ app.seeker.headline }}
            </p>
            <p class="text-xs text-slate-500 mt-2">
              Applied {{ timeAgo(app.created_at) }} · last updated {{ timeAgo(app.updated_at) }}
            </p>
          </div>
          <AppBadge :variant="applicationStatusVariant(app.status)">
            {{ applicationStatusLabel(app.status) }}
          </AppBadge>
        </div>
      </AppCard>

      <AppCard>
        <h2 class="text-sm font-medium text-slate-500 mb-3">Job</h2>
        <div class="flex items-start gap-3">
          <div class="w-10 h-10 rounded-lg bg-slate-100 flex items-center justify-center flex-shrink-0">
            <Briefcase class="w-5 h-5 text-slate-500" />
          </div>
          <div class="min-w-0 flex-1">
            <RouterLink
              :to="{ name: 'admin-job-detail', params: { id: app.job_id } }"
              class="text-base font-semibold text-slate-900 hover:text-primary-600"
            >
              {{ app.job_title }}
            </RouterLink>
            <div class="text-sm text-slate-600 mt-1 inline-flex items-center gap-1">
              <Building2 class="w-3.5 h-3.5" />
              <RouterLink
                :to="{ name: 'company-public', params: { id: app.company_id } }"
                class="hover:text-primary-600"
              >
                {{ app.company_name }}
              </RouterLink>
            </div>
            <div class="mt-2">
              <AppBadge :variant="statusVariant(app.job_status)" size="sm">
                {{ statusLabel(app.job_status) }}
              </AppBadge>
            </div>
          </div>
        </div>
      </AppCard>

      <AppCard>
        <h2 class="text-sm font-medium text-slate-500 mb-3">Applicant</h2>
        <div class="space-y-2 text-sm">
          <div class="flex items-center gap-2 text-slate-700">
            <Mail class="w-4 h-4 text-slate-400" />
            {{ app.seeker.email }}
          </div>
          <div v-if="app.seeker.phone" class="flex items-center gap-2 text-slate-700">
            <Phone class="w-4 h-4 text-slate-400" />
            {{ app.seeker.phone }}
          </div>
          <div v-if="app.seeker.location" class="flex items-center gap-2 text-slate-700">
            <MapPin class="w-4 h-4 text-slate-400" />
            {{ app.seeker.location }}
          </div>
        </div>
      </AppCard>

      <AppCard v-if="app.seeker.resume_url">
        <h2 class="text-sm font-medium text-slate-500 mb-3">Resume</h2>
        <div class="flex items-center gap-3 p-3 rounded-lg bg-slate-50 border border-slate-200">
          <div class="w-10 h-10 rounded-lg bg-red-100 flex items-center justify-center flex-shrink-0">
            <FileText class="w-5 h-5 text-red-600" />
          </div>
          <div class="min-w-0 flex-1">
            <p class="text-sm font-medium text-slate-800 truncate">
              {{ app.seeker.resume_filename || 'resume.pdf' }}
            </p>
          </div>
          <a
            :href="app.seeker.resume_url"
            target="_blank"
            rel="noopener"
            class="p-2 text-slate-500 hover:text-primary-600 rounded-md hover:bg-white"
            title="Download"
          >
            <Download class="w-4 h-4" />
          </a>
        </div>
      </AppCard>

      <AppCard v-if="app.seeker.bio">
        <h2 class="text-sm font-medium text-slate-500 mb-2">Bio</h2>
        <p class="text-slate-800 whitespace-pre-wrap text-sm">{{ app.seeker.bio }}</p>
      </AppCard>

      <AppCard v-if="app.cover_letter">
        <h2 class="text-sm font-medium text-slate-500 mb-2">Cover letter</h2>
        <p class="text-slate-800 whitespace-pre-wrap text-sm">{{ app.cover_letter }}</p>
      </AppCard>

      <p class="text-xs text-slate-400 text-center">
        Admin view is read-only. Employers manage application status.
      </p>
    </template>
  </div>
</template>
