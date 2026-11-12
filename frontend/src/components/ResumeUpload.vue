<script setup>
import { ref, computed } from 'vue'
import { FileText, Upload, Trash2, Download, Loader2 } from 'lucide-vue-next'
import { uploadResume, deleteResume } from '../api/uploads'
import AppButton from './ui/AppButton.vue'
import { useToast } from '../composables/useToast'

const props = defineProps({
  resumeUrl: { type: String, default: '' },
  resumeFilename: { type: String, default: '' },
  resumeUploadedAt: { type: String, default: '' }
})

const emit = defineEmits(['updated'])

const toast = useToast()

const fileInput = ref(null)
const uploading = ref(false)
const deleting = ref(false)

const hasResume = computed(() => !!props.resumeUrl)

function pickFile() {
  fileInput.value?.click()
}

async function onFileChange(e) {
  const file = e.target.files?.[0]
  if (!file) return

  if (file.type !== 'application/pdf') {
    toast.error('Only PDF files are allowed.')
    e.target.value = ''
    return
  }
  if (file.size > 5 * 1024 * 1024) {
    toast.error('File must be under 5 MB.')
    e.target.value = ''
    return
  }

  uploading.value = true
  try {
    const profile = await uploadResume(file)
    toast.success('Resume uploaded.')
    emit('updated', profile)
  } catch (err) {
    toast.error(err.message)
  } finally {
    uploading.value = false
    e.target.value = ''
  }
}

async function onDelete() {
  if (!confirm('Remove your resume?')) return
  deleting.value = true
  try {
    const profile = await deleteResume()
    toast.success('Resume removed.')
    emit('updated', profile)
  } catch (err) {
    toast.error(err.message)
  } finally {
    deleting.value = false
  }
}
</script>

<template>
  <div>
    <input
      ref="fileInput"
      type="file"
      accept="application/pdf,.pdf"
      class="hidden"
      @change="onFileChange"
    />

    <div v-if="!hasResume" class="space-y-2">
      <div
        class="border-2 border-dashed border-slate-300 rounded-lg p-6 text-center
               hover:border-primary-400 transition-colors"
      >
        <FileText class="w-8 h-8 mx-auto text-slate-400" :stroke-width="1.5" />
        <p class="text-sm text-slate-600 mt-2">No resume uploaded</p>
        <p class="text-xs text-slate-500 mt-1">PDF, up to 5 MB</p>
        <div class="mt-3">
          <AppButton
            variant="secondary"
            size="sm"
            :loading="uploading"
            @click="pickFile"
          >
            <Upload class="w-4 h-4" />
            {{ uploading ? 'Uploading…' : 'Upload resume' }}
          </AppButton>
        </div>
      </div>
    </div>

    <div v-else class="flex items-center gap-3 p-3 rounded-lg bg-slate-50 border border-slate-200">
      <div class="w-10 h-10 rounded-lg bg-red-100 flex items-center justify-center flex-shrink-0">
        <FileText class="w-5 h-5 text-red-600" />
      </div>
      <div class="min-w-0 flex-1">
        <p class="text-sm font-medium text-slate-800 truncate">
          {{ resumeFilename || 'resume.pdf' }}
        </p>
        <p v-if="resumeUploadedAt" class="text-xs text-slate-500">
          Uploaded {{ new Date(resumeUploadedAt).toLocaleDateString() }}
        </p>
      </div>
      <div class="flex items-center gap-1">
        <a
          :href="resumeUrl"
          target="_blank"
          rel="noopener"
          class="p-2 text-slate-500 hover:text-primary-600 rounded-md hover:bg-white"
          title="Download"
        >
          <Download class="w-4 h-4" />
        </a>
        <button
          type="button"
          :disabled="uploading"
          @click="pickFile"
          class="p-2 text-slate-500 hover:text-primary-600 rounded-md hover:bg-white disabled:opacity-50"
          title="Replace"
        >
          <Loader2 v-if="uploading" class="w-4 h-4 animate-spin" />
          <Upload v-else class="w-4 h-4" />
        </button>
        <button
          type="button"
          :disabled="deleting"
          @click="onDelete"
          class="p-2 text-slate-500 hover:text-red-600 rounded-md hover:bg-white disabled:opacity-50"
          title="Delete"
        >
          <Trash2 class="w-4 h-4" />
        </button>
      </div>
    </div>
  </div>
</template>
