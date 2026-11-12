<script setup>
const props = defineProps({
  modelValue: { type: [String, Number], default: '' },
  type: { type: String, default: 'text' },
  label: { type: String, default: '' },
  placeholder: { type: String, default: '' },
  required: { type: Boolean, default: false },
  disabled: { type: Boolean, default: false },
  error: { type: String, default: '' },
  hint: { type: String, default: '' },
  maxlength: { type: [String, Number], default: null },
  min: { type: [String, Number], default: null },
  max: { type: [String, Number], default: null },
  autocomplete: { type: String, default: '' }
})

const emit = defineEmits(['update:modelValue'])

function onInput(e) {
  emit('update:modelValue', e.target.value)
}
</script>

<template>
  <div>
    <label v-if="label" class="block text-sm font-medium text-fg-muted mb-1">
      {{ label }}
      <span v-if="required" class="text-danger">*</span>
    </label>

    <input
      :value="modelValue"
      @input="onInput"
      :type="type"
      :placeholder="placeholder"
      :required="required"
      :disabled="disabled"
      :maxlength="maxlength"
      :min="min"
      :max="max"
      :autocomplete="autocomplete || undefined"
      class="w-full rounded-lg border px-3 py-2 text-sm bg-surface text-fg
             placeholder:text-fg-subtle
             focus:outline-none focus:ring-2
             disabled:bg-surface-muted disabled:cursor-not-allowed"
      :class="error
        ? 'border-danger focus:ring-danger'
        : 'border-border-strong focus:ring-primary'"
    />

    <p v-if="error" class="text-xs text-danger-fg mt-1">{{ error }}</p>
    <p v-else-if="hint" class="text-xs text-fg-subtle mt-1">{{ hint }}</p>
  </div>
</template>
