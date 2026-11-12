<script setup>
import { ref, watch } from 'vue'
import { X } from 'lucide-vue-next'
import AppButton from './ui/AppButton.vue'
import AppTextarea from './ui/AppTextarea.vue'

const props = defineProps({
  jobId: { type: String, required: true },
  jobTitle: { type: String, default: '' },
  open: { type: Boolean, default: false },
  saving: { type: Boolean, default: false }
})

const emit = defineEmits(['close', 'submit'])

const coverLetter = ref('')
const errorMsg = ref('')

watch(
  () => props.open,
  (v) => {
    if (v) {
      coverLetter.value = ''
      errorMsg.value = ''
    }
  }
)

function handleSubmit() {
  if (coverLetter.value.length > 5000) {
    errorMsg.value = 'Cover letter must be 5000 characters or fewer.'
    return
  }
  emit('submit', { job_id: props.jobId, cover_letter: coverLetter.value.trim() })
}
</script>

<template>
  <div
    v-if="open"
    class="fixed inset-0 bg-black/40 flex items-center justify-center z-50 p-4"
    @click.self="emit('close')"
  >
    <div class="bg-white rounded-xl shadow-lg max-w-lg w-full p-6 space-y-4">
      <div class="flex items-start justify-between gap-4">
        <div>
          <h2 class="text-lg font-semibold text-slate-900">Apply to {{ jobTitle }}</h2>
          <p class="text-sm text-slate-600 mt-1">
            You can optionally include a short cover letter.
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

      <AppTextarea
        v-model="coverLetter"
        :rows="6"
        :maxlength="5000"
        show-count
        placeholder="Optional cover letter…"
        :error="errorMsg"
      />

      <div class="flex justify-end gap-2">
        <AppButton variant="secondary" :disabled="saving" @click="emit('close')">
          Cancel
        </AppButton>
        <AppButton :loading="saving" @click="handleSubmit">
          {{ saving ? 'Submitting…' : 'Submit application' }}
        </AppButton>
      </div>
    </div>
  </div>
</template>
