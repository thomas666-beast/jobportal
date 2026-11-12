<script setup>
import { ref, computed, watch } from 'vue'
import AppButton from './ui/AppButton.vue'
import AppInput from './ui/AppInput.vue'
import AppTextarea from './ui/AppTextarea.vue'
import AppSelect from './ui/AppSelect.vue'

const props = defineProps({
  initial: { type: Object, default: () => ({}) },
  submitLabel: { type: String, default: 'Save' },
  disabled: { type: Boolean, default: false },
  saving: { type: Boolean, default: false }
})

const emit = defineEmits(['submit'])

const form = ref({
  title: '', description: '', requirements: '', location: '', remote_type: 'on_site',
  salary_min: null, salary_max: null, salary_currency: '', employment_type: 'full_time',
  category: '', tags_input: ''
})

watch(
  () => props.initial,
  (val) => {
    form.value = {
      title: val.title || '',
      description: val.description || '',
      requirements: val.requirements || '',
      location: val.location || '',
      remote_type: val.remote_type || 'on_site',
      salary_min: val.salary_min ?? null,
      salary_max: val.salary_max ?? null,
      salary_currency: val.salary_currency || '',
      employment_type: val.employment_type || 'full_time',
      category: val.category || '',
      tags_input: (val.tags || []).join(', ')
    }
  },
  { immediate: true }
)

const canSubmit = computed(() => form.value.title.trim() && form.value.description.trim() && !props.saving && !props.disabled)

const salaryError = computed(() => {
  const { salary_min: min, salary_max: max } = form.value
  if (min == null || max == null) return ''
  if (min > max) return 'Minimum cannot be greater than maximum'
  return ''
})

function handleSubmit() {
  if (salaryError.value) return
  const tags = form.value.tags_input.split(',').map((t) => t.trim()).filter(Boolean).slice(0, 10)
  emit('submit', {
    title: form.value.title.trim(),
    description: form.value.description.trim(),
    requirements: form.value.requirements.trim(),
    location: form.value.location.trim(),
    remote_type: form.value.remote_type,
    salary_min: form.value.salary_min === '' ? null : form.value.salary_min,
    salary_max: form.value.salary_max === '' ? null : form.value.salary_max,
    salary_currency: form.value.salary_currency.trim(),
    employment_type: form.value.employment_type,
    category: form.value.category.trim(),
    tags
  })
}
</script>

<template>
  <form @submit.prevent="handleSubmit" class="space-y-6">
    <AppInput v-model="form.title" label="Title" required :maxlength="200" :disabled="disabled" />
    <AppTextarea v-model="form.description" label="Description" required :rows="6" :maxlength="10000" :disabled="disabled" show-count />
    <AppTextarea v-model="form.requirements" label="Requirements" :rows="4" :maxlength="10000" :disabled="disabled" show-count />

    <div class="grid md:grid-cols-2 gap-6">
      <AppInput v-model="form.location" label="Location" :maxlength="120" :disabled="disabled" />
      <AppSelect v-model="form.remote_type" label="Remote type" :disabled="disabled">
        <option value="on_site">On-site</option>
        <option value="hybrid">Hybrid</option>
        <option value="remote">Remote</option>
      </AppSelect>
    </div>

    <div class="grid md:grid-cols-3 gap-6">
      <AppInput v-model.number="form.salary_min" type="number" label="Salary min" :min="0" :disabled="disabled" />
      <AppInput v-model.number="form.salary_max" type="number" label="Salary max" :min="0" :disabled="disabled" :error="salaryError" />
      <AppInput v-model="form.salary_currency" label="Currency" placeholder="EUR" :maxlength="3" :disabled="disabled" />
    </div>

    <div class="grid md:grid-cols-2 gap-6">
      <AppSelect v-model="form.employment_type" label="Employment type" :disabled="disabled">
        <option value="full_time">Full-time</option>
        <option value="part_time">Part-time</option>
        <option value="contract">Contract</option>
        <option value="internship">Internship</option>
      </AppSelect>
      <AppInput v-model="form.category" label="Category" :maxlength="100" :disabled="disabled" />
    </div>

    <AppInput v-model="form.tags_input" label="Tags" placeholder="golang, postgres, backend" :disabled="disabled" hint="Comma-separated. Max 10 tags, each up to 30 characters." />

    <div class="flex justify-end pt-4">
      <AppButton type="submit" variant="success" :disabled="!canSubmit || !!salaryError" :loading="saving" size="lg">
        {{ saving ? 'Saving…' : submitLabel }}
      </AppButton>
    </div>
  </form>
</template>
