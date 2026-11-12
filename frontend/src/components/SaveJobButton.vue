<script setup>
import { ref, onMounted, computed } from 'vue'
import { Bookmark, BookmarkCheck } from 'lucide-vue-next'
import { saveJob, unsaveJob, checkSaved } from '../api/savedJobs'
import { useAuthStore } from '../stores/auth'
import { useToast } from '../composables/useToast'

const props = defineProps({
  jobId: { type: String, required: true },
  variant: { type: String, default: 'icon' } // 'icon' | 'full'
})

const auth = useAuthStore()
const toast = useToast()

const saved = ref(false)
const loading = ref(false)
const checking = ref(true)

const visible = computed(
  () =>
    auth.isAuthenticated &&
    auth.user?.role === 'job_seeker' &&
    auth.user?.onboardingCompleted === true
)

async function check() {
  if (!visible.value) {
    checking.value = false
    return
  }
  try {
    const res = await checkSaved(props.jobId)
    saved.value = !!res.saved
  } catch {
    saved.value = false
  } finally {
    checking.value = false
  }
}

async function toggle() {
  if (loading.value) return
  loading.value = true
  try {
    if (saved.value) {
      await unsaveJob(props.jobId)
      saved.value = false
      toast.success('Removed from saved jobs.')
    } else {
      await saveJob(props.jobId)
      saved.value = true
      toast.success('Saved to your list.')
    }
  } catch (err) {
    toast.error(err.message)
  } finally {
    loading.value = false
  }
}

onMounted(check)
</script>

<template>
  <button
    v-if="visible"
    type="button"
    :disabled="loading || checking"
    @click.stop.prevent="toggle"
    :title="saved ? 'Remove from saved' : 'Save job'"
    class="inline-flex items-center gap-2 rounded-lg transition-colors
           disabled:opacity-50 disabled:cursor-not-allowed"
    :class="[
      variant === 'icon'
        ? 'p-2 text-slate-500 hover:text-primary-600 hover:bg-slate-100'
        : 'px-4 py-2 border text-sm font-medium',
      variant === 'full' && saved
        ? 'bg-primary-50 border-primary-200 text-primary-700 hover:bg-primary-100'
        : variant === 'full'
          ? 'bg-white border-slate-300 text-slate-700 hover:bg-slate-50'
          : ''
    ]"
  >
    <component
      :is="saved ? BookmarkCheck : Bookmark"
      class="w-4 h-4"
    />
    <span v-if="variant === 'full'">
      {{ saved ? 'Saved' : 'Save job' }}
    </span>
  </button>
</template>
