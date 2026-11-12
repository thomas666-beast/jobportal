<script setup>
import { ref, onMounted, computed } from 'vue'
import { getCompany, saveCompany } from '../api/company'
import LogoUpload from '../components/LogoUpload.vue'
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

const form = ref({
  name: '',
  description: '',
  website: '',
  location: ''
})

const logoUrl = ref('')

const canSubmit = computed(() => form.value.name.trim() && !saving.value)

function applyCompany(c) {
  form.value = {
    name: c.name || '',
    description: c.description || '',
    website: c.website || '',
    location: c.location || ''
  }
  logoUrl.value = c.logo_url || ''
}

onMounted(async () => {
  try {
    const c = await getCompany()
    applyCompany(c)
    isNew.value = false
  } catch (err) {
    if (err.status === 404) {
      isNew.value = true
    } else {
      errorMsg.value = err.message
    }
  } finally {
    loading.value = false
  }
})

async function handleSubmit() {
  errorMsg.value = ''
  saving.value = true
  try {
    const c = await saveCompany({
      ...form.value,
      logo_url: logoUrl.value // preserve existing on save
    })
    applyCompany(c)
    isNew.value = false
    toast.success('Company saved.')
  } catch (err) {
    errorMsg.value = err.message
  } finally {
    saving.value = false
  }
}

function onLogoUpdated(company) {
  logoUrl.value = company.logo_url || ''
}
</script>

<template>
  <div class="max-w-2xl mx-auto space-y-6">
    <div>
      <h1 class="text-2xl font-semibold text-slate-900">
        {{ isNew ? 'Create your company' : 'Your company' }}
      </h1>
      <p class="text-sm text-slate-600 mt-1">
        This is shown alongside every job you post.
      </p>
    </div>

    <p v-if="loading" class="text-slate-600">Loading…</p>

    <template v-else>
      <AppCard>
        <form @submit.prevent="handleSubmit" class="space-y-4">
          <AppInput v-model="form.name" label="Company name" required :maxlength="200" />

          <AppTextarea
            v-model="form.description"
            label="Description"
            :rows="5"
            :maxlength="5000"
            show-count
          />

          <div class="grid md:grid-cols-2 gap-4">
            <AppInput
              v-model="form.website"
              type="url"
              label="Website"
              placeholder="https://example.com"
              :maxlength="255"
            />
            <AppInput
              v-model="form.location"
              label="Location"
              placeholder="City or region"
              :maxlength="120"
            />
          </div>

          <p v-if="errorMsg" class="text-sm text-red-600">{{ errorMsg }}</p>

          <div class="flex justify-end">
            <AppButton type="submit" :disabled="!canSubmit" :loading="saving" size="lg">
              {{ saving ? 'Saving…' : isNew ? 'Create company' : 'Save changes' }}
            </AppButton>
          </div>
        </form>
      </AppCard>

      <AppCard v-if="!isNew">
        <div class="mb-3">
          <h2 class="text-base font-semibold text-slate-900">Company logo</h2>
          <p class="text-sm text-slate-600 mt-1">
            Shown on every job listing and detail page.
          </p>
        </div>

        <LogoUpload :logo-url="logoUrl" @updated="onLogoUpdated" />
      </AppCard>
    </template>
  </div>
</template>
