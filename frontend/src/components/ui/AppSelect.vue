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
    <label v-if="label" class="block text-sm font-medium text-slate-700 mb-1">
      {{ label }}
      <span v-if="required" class="text-red-500">*</span>
    </label>

    <select
      :value="modelValue"
      @change="onChange"
      :required="required"
      :disabled="disabled"
      class="w-full rounded-lg border px-3 py-2 text-sm bg-white
             focus:outline-none focus:ring-2
             disabled:bg-slate-100 disabled:cursor-not-allowed"
      :class="error
        ? 'border-red-300 focus:ring-red-500'
        : 'border-slate-300 focus:ring-primary-500'"
    >
      <slot />
    </select>

    <p v-if="error" class="text-xs text-red-600 mt-1">{{ error }}</p>
    <p v-else-if="hint" class="text-xs text-slate-500 mt-1">{{ hint }}</p>
  </div>
</template>
