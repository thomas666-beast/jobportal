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
    <label v-if="label" class="block text-sm font-medium text-slate-700 mb-1">
      {{ label }}
      <span v-if="required" class="text-red-500">*</span>
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
      class="w-full rounded-lg border px-3 py-2 text-sm
             placeholder:text-slate-400
             focus:outline-none focus:ring-2
             disabled:bg-slate-100 disabled:cursor-not-allowed"
      :class="error
        ? 'border-red-300 focus:ring-red-500'
        : 'border-slate-300 focus:ring-primary-500'"
    />

    <p v-if="error" class="text-xs text-red-600 mt-1">{{ error }}</p>
    <p v-else-if="hint" class="text-xs text-slate-500 mt-1">{{ hint }}</p>
  </div>
</template>
