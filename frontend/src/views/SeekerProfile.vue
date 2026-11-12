<script setup>
import { ref, onMounted, computed } from 'vue'
import { getProfile, saveProfile } from '../api/seeker'
import ResumeUpload from '../components/ResumeUpload.vue'
import AppButton from '../components/ui/AppButton.vue'
import AppInput from '../components/ui/AppInput.vue'
import AppTextarea from '../components/ui/AppTextarea.vue'
import AppCard from '../components/ui/AppCard.vue'
import { useToast } from '../composables/useToast'

const toast = useToast()

const loading = ref(true)
const saving = ref(false)
const errorMsg = ref('')
const isNew = ref(false)

const form = ref({ full_name: '', phone: '', location: '', headline: '', bio: '' })
const resumeUrl = ref('')
const resumeFilename = ref('')
const resumeUploadedAt = ref('')

const canSubmit = computed(() => form.value.full_name.trim() && !saving.value)

function applyProfile(p) {
  form.value = {
    full_name: p.full_name || '',
    phone: p.phone || '',
    location: p.location || '',
    headline: p.headline || '',
    bio: p.bio || ''
  }
  resumeUrl.value = p.resume_url || ''
  resumeFilename.value = p.resume_filename || ''
  resumeUploadedAt.value = p.resume_uploaded_at || ''
}

onMounted(async () => {
  try {
    const p = await getProfile()
    applyProfile(p)
    isNew.value = false
  } catch (err) {
    if (err.status === 404) { isNew.value = true } 
    else { errorMsg.value = err.message }
  } finally {
    loading.value = false
  }
})

async function handleSubmit() {
  errorMsg.value = ''
  saving.value = true
  try {
    const p = await saveProfile(form.value)
    applyProfile(p)
    isNew.value = false
    toast.success('Profile saved.')
  } catch (err) {
    errorMsg.value = err.message
  } finally {
    saving.value = false
  }
}

function onResumeUpdated(p) { applyProfile(p) }
</script>

<template>
  <div class="max-w-2xl mx-auto space-y-6">
    <div>
      <h1 class="text-3xl font-bold text-fg">
        {{ isNew ? 'Create your profile' : 'Your profile' }}
      </h1>
      <p class="text-sm text-fg-muted mt-2">Employers will see this when you apply to their jobs.</p>
    </div>

    <p v-if="loading" class="text-fg-muted">Loading…</p>

    <template v-else>
      <AppCard>
        <form @submit.prevent="handleSubmit" class="space-y-6">
          <AppInput v-model="form.full_name" label="Full name" required :maxlength="120" />
          <AppInput v-model="form.headline" label="Headline" placeholder="e.g. Backend Go developer" :maxlength="120" hint="A short tagline employers see first" />
          
          <div class="grid md:grid-cols-2 gap-6">
            <AppInput v-model="form.phone" type="tel" label="Phone" :maxlength="32" />
            <AppInput v-model="form.location" label="Location" placeholder="City or region" :maxlength="120" />
          </div>
          
          <AppTextarea v-model="form.bio" label="Bio" :rows="6" :maxlength="2000" show-count />

          <p v-if="errorMsg" class="text-sm text-red-500">{{ errorMsg }}</p>

          <div class="flex justify-end">
            <AppButton type="submit" :disabled="!canSubmit" :loading="saving" size="lg">
              {{ saving ? 'Saving…' : isNew ? 'Create profile' : 'Save changes' }}
            </AppButton>
          </div>
        </form>
      </AppCard>

      <AppCard v-if="!isNew">
        <div class="mb-6">
          <h2 class="text-xl font-bold text-fg">Resume</h2>
          <p class="text-sm text-fg-muted mt-1">Employers see this when you apply to jobs.</p>
        </div>
        <ResumeUpload :resume-url="resumeUrl" :resume-filename="resumeFilename" :resume-uploaded-at="resumeUploadedAt" @updated="onResumeUpdated" />
      </AppCard>
    </template>
  </div>
</template>
