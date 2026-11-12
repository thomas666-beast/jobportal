<script setup>
import { computed } from 'vue'
import { RouterLink } from 'vue-router'

const props = defineProps({
  variant: { type: String, default: 'primary' },
  size: { type: String, default: 'md' },
  type: { type: String, default: 'button' },
  to: { type: [String, Object], default: null },
  disabled: { type: Boolean, default: false },
  loading: { type: Boolean, default: false },
  block: { type: Boolean, default: false }
})

const baseClasses = 'inline-flex items-center justify-center gap-2 font-semibold uppercase tracking-wider transition-all duration-200 disabled:opacity-50 disabled:cursor-not-allowed border'

const variantClasses = computed(() => ({
  'bg-primary text-white border-primary hover:bg-transparent hover:text-primary': props.variant === 'primary',
  'bg-transparent text-fg border-border-strong hover:border-primary hover:text-primary': props.variant === 'secondary',
  'bg-transparent text-red-500 border-red-500 hover:bg-red-500 hover:text-white': props.variant === 'danger',
  'bg-transparent text-fg-muted border-transparent hover:text-fg': props.variant === 'ghost',
  'bg-green-600 text-white border-green-600 hover:bg-transparent hover:text-green-500': props.variant === 'success'
}))

const sizeClasses = computed(() => ({
  'text-xs px-4 py-1.5': props.size === 'sm',
  'text-sm px-6 py-2.5': props.size === 'md',
  'text-base px-8 py-3': props.size === 'lg'
}))

const classes = computed(() => [baseClasses, variantClasses.value, sizeClasses.value, props.block && 'w-full'])
</script>

<template>
  <RouterLink v-if="to && !disabled" :to="to" :class="classes"><slot /></RouterLink>
  <button v-else :type="type" :disabled="disabled || loading" :class="classes">
    <span v-if="loading" class="inline-block h-3.5 w-3.5 rounded-full border-2 border-current border-t-transparent animate-spin" />
    <slot />
  </button>
</template>
