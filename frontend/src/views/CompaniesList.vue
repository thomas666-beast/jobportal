<script setup>
import { ref, onMounted, computed } from 'vue'
import { RouterLink } from 'vue-router'
import { Search, Building2, MapPin, Briefcase } from 'lucide-vue-next'
import { listCompanies } from '../api/publicCompanies'
import AppCard from '../components/ui/AppCard.vue'
import AppBadge from '../components/ui/AppBadge.vue'
import AppButton from '../components/ui/AppButton.vue'
import EmptyState from '../components/ui/EmptyState.vue'

const q = ref('')
const companies = ref([])
const total = ref(0)
const loading = ref(false)
const errorMsg = ref('')

const limit = 20
const offset = ref(0)

const initialLoading = ref(true)
const canLoadMore = computed(() => offset.value + companies.value.length < total.value)

let debounceTimer = null

async function load(reset = true) {
  loading.value = true
  errorMsg.value = ''
  try {
    const res = await listCompanies({
      q: q.value.trim(),
      limit,
      offset: reset ? 0 : offset.value
    })
    if (reset) {
      companies.value = res.companies || []
      offset.value = 0
    } else {
      companies.value = [...companies.value, ...(res.companies || [])]
    }
    total.value = res.total || 0
  } catch (err) {
    errorMsg.value = err.message
  } finally {
    loading.value = false
    initialLoading.value = false
  }
}

function debouncedLoad() {
  clearTimeout(debounceTimer)
  debounceTimer = setTimeout(() => load(true), 300)
}

function loadMore() {
  offset.value = companies.value.length
  load(false)
}

onMounted(() => load(true))
</script>

<template>
  <div class="space-y-6">
    <section class="text-center pt-6 pb-2">
      <h1 class="text-3xl md:text-4xl text-fg">
        Companies hiring now
      </h1>
      <p class="mt-2 text-fg-muted">
        {{ total }} {{ total === 1 ? 'company' : 'companies' }} on Job Portal
      </p>
    </section>

    <AppCard class="!rounded-2xl">
      <div class="relative">
        <Search class="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-fg-subtle pointer-events-none" />
        <input
          v-model="q"
          @input="debouncedLoad"
          type="text"
          placeholder="Search companies by name"
          class="w-full rounded-lg border border-border pl-9 pr-3 py-2.5 text-sm bg-surface text-fg
                 placeholder:text-fg-subtle
                 focus:outline-none focus:ring-2 focus:ring-primary"
        />
      </div>
    </AppCard>

    <p v-if="loading && initialLoading" class="text-fg-muted text-center py-8">
      Loading companies…
    </p>
    <p v-else-if="errorMsg" class="text-danger-fg text-center">{{ errorMsg }}</p>

    <AppCard v-else-if="companies.length === 0" :padded="false">
      <EmptyState
        :icon="Building2"
        title="No companies found"
        :description="q ? 'Try a different search term.' : 'No companies registered yet.'"
      />
    </AppCard>

    <div v-else class="grid md:grid-cols-2 gap-3">
      <RouterLink
        v-for="c in companies"
        :key="c.id"
        :to="{ name: 'company-public', params: { id: c.id } }"
        class="block"
      >
        <AppCard hover class="h-full">
          <div class="flex items-start gap-4">
            <div class="flex-shrink-0">
              <div
                v-if="c.logo_url"
                class="w-14 h-14 rounded-xl border border-border overflow-hidden bg-surface-muted"
              >
                <img
                  :src="c.logo_url"
                  :alt="c.name"
                  class="w-full h-full object-cover"
                />
              </div>
              <div
                v-else
                class="w-14 h-14 rounded-xl bg-surface-muted flex items-center justify-center"
              >
                <Building2 class="w-7 h-7 text-fg-muted" />
              </div>
            </div>

            <div class="flex-1 min-w-0">
              <h2 class="text-base text-fg truncate">
                {{ c.name }}
              </h2>
              <div
                v-if="c.location"
                class="text-sm text-fg-muted mt-0.5 inline-flex items-center gap-1"
              >
                <MapPin class="w-3.5 h-3.5" /> {{ c.location }}
              </div>
              <p
                v-if="c.description"
                class="text-sm text-fg-muted mt-2 line-clamp-2"
              >
                {{ c.description }}
              </p>
              <div class="mt-3">
                <AppBadge v-if="c.published_job_count > 0" variant="primary" size="sm">
                  <Briefcase class="w-3 h-3 mr-1" />
                  {{ c.published_job_count }}
                  {{ c.published_job_count === 1 ? 'open job' : 'open jobs' }}
                </AppBadge>
                <AppBadge v-else variant="neutral" size="sm">
                  No open jobs
                </AppBadge>
              </div>
            </div>
          </div>
        </AppCard>
      </RouterLink>
    </div>

    <div v-if="canLoadMore" class="text-center pt-2">
      <AppButton variant="secondary" :loading="loading" @click="loadMore">
        {{ loading ? 'Loading…' : 'Load more companies' }}
      </AppButton>
    </div>
  </div>
</template>
