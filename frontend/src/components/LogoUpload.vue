<script setup>
import { ref, computed } from 'vue'
import { Upload, Trash2, Image as ImageIcon } from 'lucide-vue-next'
import { uploadLogo, deleteLogo } from '../api/uploads'
import AppButton from './ui/AppButton.vue'
import { useToast } from '../composables/useToast'

const props = defineProps({
  logoUrl: { type: String, default: '' }
})

const emit = defineEmits(['updated'])

const toast = useToast()

const fileInput = ref(null)
const uploading = ref(false)
const deleting = ref(false)

const hasLogo = computed(() => !!props.logoUrl)

function pickFile() {
  fileInput.value?.click()
}

async function onFileChange(e) {
  const file = e.target.files?.[0]
  if (!file) return

  const allowed = ['image/png', 'image/jpeg', 'image/webp']
  if (!allowed.includes(file.type)) {
    toast.error('Only PNG, JPEG, or WEBP images are allowed.')
    e.target.value = ''
    return
  }
  if (file.size > 2 * 1024 * 1024) {
    toast.error('Image must be under 2 MB.')
    e.target.value = ''
    return
  }

  uploading.value = true
  try {
    const company = await uploadLogo(file)
    toast.success('Logo uploaded.')
    emit('updated', company)
  } catch (err) {
    toast.error(err.message)
  } finally {
    uploading.value = false
    e.target.value = ''
  }
}

async function onDelete() {
  if (!confirm('Remove the logo?')) return
  deleting.value = true
  try {
    const company = await deleteLogo()
    toast.success('Logo removed.')
    emit('updated', company)
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
      accept="image/png,image/jpeg,image/webp"
      class="hidden"
      @change="onFileChange"
    />

    <div class="flex items-center gap-4">
      <div
        class="w-20 h-20 rounded-xl border-2 border-dashed border-slate-300
               flex items-center justify-center overflow-hidden bg-white flex-shrink-0"
      >
        <img
          v-if="hasLogo"
          :src="logoUrl"
          alt="Logo"
          class="w-full h-full object-cover"
        />
        <ImageIcon v-else class="w-8 h-8 text-slate-300" :stroke-width="1.5" />
      </div>

      <div class="space-y-2">
        <div class="flex gap-2">
          <AppButton
            variant="secondary"
            size="sm"
            :loading="uploading"
            @click="pickFile"
          >
            <Upload class="w-4 h-4" />
            {{ hasLogo ? 'Replace' : 'Upload logo' }}
          </AppButton>
          <AppButton
            v-if="hasLogo"
            variant="ghost"
            size="sm"
            :disabled="deleting"
            @click="onDelete"
          >
            <Trash2 class="w-4 h-4" />
            Remove
          </AppButton>
        </div>
        <p class="text-xs text-slate-500">
          PNG, JPEG, or WEBP. Up to 2 MB. Square images look best.
        </p>
      </div>
    </div>
  </div>
</template>
