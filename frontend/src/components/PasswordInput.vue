<script setup>
import { ref } from 'vue'
import { generatePassword } from '../utils/password'

const props = defineProps({
  modelValue: { type: String, default: '' },
  placeholder: { type: String, default: '' },
  required: { type: Boolean, default: false },
  minlength: { type: Number, default: 8 },
  showGenerate: { type: Boolean, default: true }
})

const emit = defineEmits(['update:modelValue'])

const visible = ref(false)

function updateValue(e) {
  emit('update:modelValue', e.target.value)
}

function handleGenerate() {
  const pwd = generatePassword(16)
  emit('update:modelValue', pwd)
  visible.value = true
}
</script>

<template>
  <div class="flex gap-2">
    <div class="relative flex-1">
      <input
        :value="modelValue"
        @input="updateValue"
        :type="visible ? 'text' : 'password'"
        :placeholder="placeholder"
        :required="required"
        :minlength="minlength"
        class="w-full rounded-lg border border-slate-300 px-3 py-2 pr-10
               focus:outline-none focus:ring-2 focus:ring-indigo-500"
      />
      <button
        type="button"
        @click="visible = !visible"
        class="absolute right-2 top-1/2 -translate-y-1/2 text-slate-500
               hover:text-slate-700 text-sm"
        :aria-label="visible ? 'Hide password' : 'Show password'"
      >
        {{ visible ? '🙈' : '👁' }}
      </button>
    </div>

    <button
      v-if="showGenerate"
      type="button"
      @click="handleGenerate"
      class="px-3 py-2 rounded-lg bg-slate-200 text-slate-700 text-sm
             hover:bg-slate-300 whitespace-nowrap"
    >
      Generate
    </button>
  </div>
</template>
