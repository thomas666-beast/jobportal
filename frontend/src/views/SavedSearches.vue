<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { Bookmark, Search, MapPin, Globe, Play, Trash2 } from 'lucide-vue-next'
import { listSavedSearches, deleteSavedSearch } from '../api/savedSearches'
import { timeAgo } from '../utils/jobStatus'
import AppButton from '../components/ui/AppButton.vue'
import AppBadge from '../components/ui/AppBadge.vue'
import AppCard from '../components/ui/AppCard.vue'
import EmptyState from '../components/ui/EmptyState.vue'
import AppSkeleton from '../components/ui/AppSkeleton.vue'
import { useToast } from '../composables/useToast'

const router = useRouter()
const toast = useToast()

const searches = ref([])
const loading = ref(true)
const errorMsg = ref('')
const acting = ref(null)

async function load() {
  loading.value = true
  errorMsg.value = ''
  try {
    searches.value = await listSavedSearches()
  } catch (err) {
    errorMsg.value = err.message
  } finally {
    loading.value = false
  }
}

async function remove(id) {
  if (!confirm('Delete this saved search?')) return
  acting.value = id
  try {
    await deleteSavedSearch(id)
    searches.value = searches.value.filter((s) => s.id !== id)
    toast.success('Saved search deleted.')
  } catch (err) {
    toast.error(err.message)
  } finally {
    acting.value = null
  }
}

function runSearch(s) {
  const query = {}
  if (s.query) query.q = s.query
  if (s.location) query.location = s.location
  if (s.remote_type) query.remote_type = s.remote_type
  router.push({ name: 'browse', query })
}

function remoteLabel(rt) {
  if (rt === 'remote') return 'Remote'
  if (rt === 'hybrid') return 'Hybrid'
  if (rt === 'on_site') return 'On-site'
  return rt
}

onMounted(load)
</script>

<template>
  <div class="space-y-4">
    <div>
      <h1 class="text-2xl font-semibold text-slate-900">Saved searches</h1>
      <p class="text-sm text-slate-600 mt-1">
        Re-run your favorite job searches with one click
      </p>
    </div>

    <div v-if="loading" class="space-y-3">
      <AppSkeleton v-for="i in 3" :key="i" variant="row" />
    </div>
    <p v-else-if="errorMsg" class="text-red-600">{{ errorMsg }}</p>

    <AppCard v-else-if="searches.length === 0" :padded="false">
      <EmptyState
        :icon="Bookmark"
        title="No saved searches"
        description="Save your browse filters to find matching jobs later."
      >
        <template #action>
          <AppButton to="/browse">Browse jobs</AppButton>
        </template>
      </EmptyState>
    </AppCard>

    <div v-else class="space-y-3">
      <AppCard v-for="s in searches" :key="s.id">
        <div class="flex flex-col md:flex-row md:items-start md:justify-between gap-4">
          <div class="min-w-0 flex-1">
            <div class="flex items-center gap-2">
              <h2 class="text-base font-semibold text-slate-900">{{ s.name }}</h2>
              <AppBadge
                :variant="s.match_count > 0 ? 'primary' : 'neutral'"
                size="sm"
              >
                {{ s.match_count }} {{ s.match_count === 1 ? 'match' : 'matches' }}
              </AppBadge>
            </div>

            <div class="flex flex-wrap items-center gap-3 mt-2 text-xs text-slate-500">
              <span v-if="s.query" class="inline-flex items-center gap-1">
                <Search class="w-3 h-3" /> {{ s.query }}
              </span>
              <span v-if="s.location" class="inline-flex items-center gap-1">
                <MapPin class="w-3 h-3" /> {{ s.location }}
              </span>
              <span v-if="s.remote_type" class="inline-flex items-center gap-1">
                <Globe class="w-3 h-3" /> {{ remoteLabel(s.remote_type) }}
              </span>
              <span class="ml-auto">Saved {{ timeAgo(s.created_at) }}</span>
            </div>
          </div>

          <div class="flex items-center gap-2 flex-shrink-0">
            <AppButton size="sm" @click="runSearch(s)">
              <Play class="w-4 h-4" /> Run
            </AppButton>
            <AppButton
              size="sm"
              variant="ghost"
              :disabled="acting === s.id"
              @click="remove(s.id)"
            >
              <Trash2 class="w-4 h-4" />
            </AppButton>
          </div>
        </div>
      </AppCard>
    </div>
  </div>
</template>
