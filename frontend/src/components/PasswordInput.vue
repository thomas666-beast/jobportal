<script setup>
import { ref } from 'vue'
import { generatePassword } from '../utils/password'
import { Eye, EyeOff } from 'lucide-vue-next'

const props = defineProps({
  modelValue: { type: String, default: '' },
  placeholder: { type: String, default: '' },
  required: { type: Boolean, default: false },
  minlength: { type: Number, default: 8 },
  showGenerate: { type: Boolean, default: true }
})

const emit = defineEmits(['update:modelValue'])
const visible = ref(false)

function updateValue(e) { emit('update:modelValue', e.target.value) }
function handleGenerate() {
  emit('update:modelValue', generatePassword(16))
  visible.value = true
}
</script>

<template>
  <div class="flex gap-3">
    <div class="relative flex-1">
      <input
        :value="modelValue"
        @input="updateValue"
        :type="visible ? 'text' : 'password'"
        :placeholder="placeholder"
        :required="required"
        :minlength="minlength"
        class="w-full rounded-none border border-border bg-bg px-3 py-2.5 pr-10 text-sm text-fg
               placeholder:text-fg-subtle focus:outline-none focus:border-primary transition-colors"
      />
      <button
        type="button"
        @click="visible = !visible"
        class="absolute right-3 top-1/2 -translate-y-1/2 text-fg-muted hover:text-fg transition-colors"
      >
        <EyeOff v-if="visible" class="w-4 h-4" />
        <Eye v-else class="w-4 h-4" />
      </button>
    </div>

    <button
      v-if="showGenerate"
      type="button"
      @click="handleGenerate"
      class="px-5 py-2.5 bg-surface border border-border-strong text-fg-muted font-bold uppercase tracking-wider text-xs
             hover:border-primary hover:text-primary transition-colors whitespace-nowrap"
    >
      Generate
    </button>
  </div>
</template>
