<script setup>
const props = defineProps({
  modelValue: { type: String, default: '' },
  label: { type: String, default: '' },
  placeholder: { type: String, default: '' },
  required: { type: Boolean, default: false },
  disabled: { type: Boolean, default: false },
  error: { type: String, default: '' },
  hint: { type: String, default: '' },
  rows: { type: Number, default: 5 },
  maxlength: { type: [String, Number], default: null },
  showCount: { type: Boolean, default: false }
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

    <textarea
      :value="modelValue"
      @input="onInput"
      :placeholder="placeholder"
      :required="required"
      :disabled="disabled"
      :rows="rows"
      :maxlength="maxlength"
      class="w-full rounded-lg border px-3 py-2 text-sm bg-surface text-fg
             placeholder:text-fg-subtle
             focus:outline-none focus:ring-2
             disabled:bg-surface-muted disabled:cursor-not-allowed"
      :class="error
        ? 'border-danger focus:ring-danger'
        : 'border-border-strong focus:ring-primary'"
    />

    <div class="flex justify-between items-center mt-1">
      <p v-if="error" class="text-xs text-danger-fg">{{ error }}</p>
      <p v-else-if="hint" class="text-xs text-fg-subtle">{{ hint }}</p>
      <p v-if="showCount && maxlength" class="text-xs text-fg-subtle ml-auto">
        {{ modelValue.length }} / {{ maxlength }}
      </p>
    </div>
  </div>
</template>
