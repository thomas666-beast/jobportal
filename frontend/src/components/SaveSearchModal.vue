<script setup>
import { ref, watch, computed } from 'vue'
import { X } from 'lucide-vue-next'
import { createSavedSearch } from '../api/savedSearches'
import AppButton from './ui/AppButton.vue'
import AppInput from './ui/AppInput.vue'
import { useToast } from '../composables/useToast'

const props = defineProps({
  open: { type: Boolean, default: false },
  query: { type: String, default: '' },
  location: { type: String, default: '' },
  remoteType: { type: String, default: '' }
})

const emit = defineEmits(['close', 'created'])

const toast = useToast()

const name = ref('')
const saving = ref(false)
const errorMsg = ref('')

watch(
  () => props.open,
  (v) => {
    if (v) {
      name.value = ''
      errorMsg.value = ''
    }
  }
)

const canSubmit = computed(() => name.value.trim() && !saving.value)

async function submit() {
  errorMsg.value = ''
  saving.value = true
  try {
    const created = await createSavedSearch({
      name: name.value.trim(),
      query: props.query,
      location: props.location,
      remote_type: props.remoteType
    })
    toast.success('Search saved.')
    emit('created', created)
  } catch (err) {
    errorMsg.value = err.message
  } finally {
    saving.value = false
  }
}
</script>

<template>
  <div
    v-if="open"
    class="fixed inset-0 bg-black/40 flex items-center justify-center z-50 p-4"
    @click.self="emit('close')"
  >
    <div class="bg-white rounded-xl shadow-lg max-w-md w-full p-6 space-y-4">
      <div class="flex items-start justify-between gap-4">
        <div>
          <h2 class="text-lg font-semibold text-slate-900">Save this search</h2>
          <p class="text-sm text-slate-600 mt-1">
            You'll find it under Saved searches.
          </p>
        </div>
        <button
          type="button"
          @click="emit('close')"
          class="text-slate-400 hover:text-slate-600 -mt-1 -mr-1 p-1"
        >
          <X class="w-5 h-5" />
        </button>
      </div>

      <div class="text-xs text-slate-500 space-y-1 p-3 bg-slate-50 rounded-lg">
        <p><strong>Query:</strong> {{ query || '—' }}</p>
        <p><strong>Location:</strong> {{ location || '—' }}</p>
        <p><strong>Remote:</strong> {{ remoteType || 'any' }}</p>
      </div>

      <AppInput
        v-model="name"
        label="Name"
        placeholder="e.g. Go remote in Helsinki"
        required
        :maxlength="100"
        :error="errorMsg"
      />

      <div class="flex justify-end gap-2">
        <AppButton variant="secondary" :disabled="saving" @click="emit('close')">
          Cancel
        </AppButton>
        <AppButton :disabled="!canSubmit" :loading="saving" @click="submit">
          {{ saving ? 'Saving…' : 'Save search' }}
        </AppButton>
      </div>
    </div>
  </div>
</template>
