<script setup>
const props = defineProps({
  modelValue: { type: [String, Number], default: '' },
  label: { type: String, default: '' },
  required: { type: Boolean, default: false },
  disabled: { type: Boolean, default: false },
  error: { type: String, default: '' },
  hint: { type: String, default: '' }
})

const emit = defineEmits(['update:modelValue'])

function onChange(e) {
  emit('update:modelValue', e.target.value)
}
</script>

<template>
  <div>
    <label v-if="label" class="block text-sm font-medium text-fg-muted mb-1">
      {{ label }}
      <span v-if="required" class="text-danger">*</span>
    </label>

    <select
      :value="modelValue"
      @change="onChange"
      :required="required"
      :disabled="disabled"
      class="w-full rounded-lg border px-3 py-2 text-sm bg-surface text-fg
             focus:outline-none focus:ring-2
             disabled:bg-surface-muted disabled:cursor-not-allowed"
      :class="error
        ? 'border-danger focus:ring-danger'
        : 'border-border-strong focus:ring-primary'"
    >
      <slot />
    </select>

    <p v-if="error" class="text-xs text-danger-fg mt-1">{{ error }}</p>
    <p v-else-if="hint" class="text-xs text-fg-subtle mt-1">{{ hint }}</p>
  </div>
</template>
