<script setup>
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ArrowLeft } from 'lucide-vue-next'
import { getAdminUser } from '../api/admin'
import { timeAgo } from '../utils/jobStatus'
import AdminNav from '../components/AdminNav.vue'
import AppBadge from '../components/ui/AppBadge.vue'
import AppCard from '../components/ui/AppCard.vue'

const route = useRoute()
const router = useRouter()

const user = ref(null)
const loading = ref(true)
const errorMsg = ref('')

async function load() {
  loading.value = true
  errorMsg.value = ''
  try {
    user.value = await getAdminUser(route.params.id)
  } catch (err) {
    errorMsg.value = err.message
  } finally {
    loading.value = false
  }
}

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

onMounted(load)
</script>

<template>
  <div class="space-y-6">
    <AdminNav />

    <button
      type="button"
      @click="router.push({ name: 'admin-users' })"
      class="inline-flex items-center gap-1 text-sm text-slate-600 hover:text-slate-900"
    >
      <ArrowLeft class="w-4 h-4" /> Back to users
    </button>

    <p v-if="loading" class="text-slate-600">Loading…</p>
    <p v-else-if="errorMsg" class="text-red-600">{{ errorMsg }}</p>

    <template v-else-if="user">
      <AppCard>
        <div class="flex items-start justify-between gap-4">
          <div>
            <h1 class="text-2xl font-semibold text-slate-900">{{ user.username }}</h1>
            <p class="text-sm text-slate-600 mt-1">{{ user.email }}</p>
          </div>
          <AppBadge :variant="roleVariant(user.role)">
            {{ roleLabel(user.role) }}
          </AppBadge>
        </div>

        <div class="grid md:grid-cols-2 gap-4 mt-5 text-sm">
          <div>
            <div class="text-slate-500">Onboarding</div>
            <div class="text-slate-800 mt-0.5">
              {{ user.onboarding_completed ? 'Completed' : 'Not completed' }}
            </div>
          </div>
          <div>
            <div class="text-slate-500">Joined</div>
            <div class="text-slate-800 mt-0.5">
              {{ timeAgo(user.created_at) }}
            </div>
          </div>
        </div>
      </AppCard>

      <AppCard v-if="user.seeker">
        <h2 class="text-base font-semibold text-slate-900 mb-3">Seeker profile</h2>
        <div class="grid md:grid-cols-2 gap-4 text-sm">
          <div>
            <div class="text-slate-500">Full name</div>
            <div class="text-slate-800 mt-0.5">{{ user.seeker.full_name || '—' }}</div>
          </div>
          <div>
            <div class="text-slate-500">Headline</div>
            <div class="text-slate-800 mt-0.5">{{ user.seeker.headline || '—' }}</div>
          </div>
          <div>
            <div class="text-slate-500">Phone</div>
            <div class="text-slate-800 mt-0.5">{{ user.seeker.phone || '—' }}</div>
          </div>
          <div>
            <div class="text-slate-500">Location</div>
            <div class="text-slate-800 mt-0.5">{{ user.seeker.location || '—' }}</div>
          </div>
        </div>
        <div v-if="user.seeker.bio" class="mt-4">
          <div class="text-slate-500 text-sm">Bio</div>
          <p class="text-slate-800 whitespace-pre-wrap text-sm mt-1">
            {{ user.seeker.bio }}
          </p>
        </div>
      </AppCard>

      <AppCard v-if="user.company">
        <h2 class="text-base font-semibold text-slate-900 mb-3">Company</h2>
        <div class="grid md:grid-cols-2 gap-4 text-sm">
          <div>
            <div class="text-slate-500">Name</div>
            <div class="text-slate-800 mt-0.5">{{ user.company.name }}</div>
          </div>
          <div>
            <div class="text-slate-500">Location</div>
            <div class="text-slate-800 mt-0.5">{{ user.company.location || '—' }}</div>
          </div>
          <div class="md:col-span-2">
            <div class="text-slate-500">Website</div>
            <div class="text-slate-800 mt-0.5">
              <a
                v-if="user.company.website"
                :href="user.company.website"
                target="_blank"
                rel="noopener"
                class="text-primary-600 hover:underline"
              >
                {{ user.company.website }}
              </a>
              <span v-else>—</span>
            </div>
          </div>
        </div>
        <div v-if="user.company.description" class="mt-4">
          <div class="text-slate-500 text-sm">Description</div>
          <p class="text-slate-800 whitespace-pre-wrap text-sm mt-1">
            {{ user.company.description }}
          </p>
        </div>
      </AppCard>
    </template>
  </div>
</template>
