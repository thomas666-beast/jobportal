<script setup>
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ArrowLeft, Mail, Phone, MapPin, FileText, Download } from 'lucide-vue-next'
import { getApplicationDetail, updateApplicationStatus } from '../api/employerApplications'
import {
  applicationStatusLabel,
  applicationStatusVariant,
  timeAgo
} from '../utils/jobStatus'
import AppBadge from '../components/ui/AppBadge.vue'
import AppCard from '../components/ui/AppCard.vue'
import AppButton from '../components/ui/AppButton.vue'
import { useToast } from '../composables/useToast'

const route = useRoute()
const router = useRouter()
const toast = useToast()

const app = ref(null)
const loading = ref(true)
const errorMsg = ref('')
const acting = ref(false)

const allowedStatuses = ['viewed', 'shortlisted', 'rejected', 'hired']

async function load() {
  loading.value = true
  errorMsg.value = ''
  try {
    app.value = await getApplicationDetail(route.params.id)
  } catch (err) {
    errorMsg.value = err.message
  } finally {
    loading.value = false
  }
}

async function setStatus(status) {
  if (app.value.status === status) return
  acting.value = true
  try {
    await updateApplicationStatus(app.value.id, status)
    toast.success(`Marked as ${applicationStatusLabel(status).toLowerCase()}.`)
    await load()
  } catch (err) {
    toast.error(err.message)
  } finally {
    acting.value = false
  }
}

onMounted(load)
</script>

<template>
  <div class="space-y-4">
    <button
      type="button"
      @click="router.push({ name: 'employer-applications' })"
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
            <p class="text-xs text-slate-500 mt-1">
              Applied for <strong>{{ app.job_title }}</strong>
              · {{ timeAgo(app.created_at) }}
            </p>
          </div>
          <AppBadge :variant="applicationStatusVariant(app.status)">
            {{ applicationStatusLabel(app.status) }}
          </AppBadge>
        </div>

        <div class="flex flex-wrap gap-2 mt-5">
          <button
            v-for="s in allowedStatuses"
            :key="s"
            type="button"
            :disabled="acting || app.status === s"
            @click="setStatus(s)"
            class="px-3 py-1.5 text-xs rounded-md border font-medium
                   disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
            :class="{
              'bg-blue-50 border-blue-200 text-blue-700 hover:bg-blue-100':
                s === 'viewed' && app.status !== s,
              'bg-primary-50 border-primary-200 text-primary-700 hover:bg-primary-100':
                s === 'shortlisted' && app.status !== s,
              'bg-red-50 border-red-200 text-red-700 hover:bg-red-100':
                s === 'rejected' && app.status !== s,
              'bg-emerald-50 border-emerald-200 text-emerald-700 hover:bg-emerald-100':
                s === 'hired' && app.status !== s,
              'bg-slate-100 border-slate-200 text-slate-400': app.status === s
            }"
          >
            {{ applicationStatusLabel(s) }}
          </button>
        </div>
      </AppCard>

      <AppCard>
        <h2 class="text-sm font-medium text-slate-500 mb-3">Contact</h2>
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
    </template>
  </div>
</template>
