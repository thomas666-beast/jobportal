<script setup>
import { computed } from 'vue'
import { RouterLink } from 'vue-router'

const props = defineProps({
  variant: {
    type: String,
    default: 'primary',
    validator: (v) => ['primary', 'secondary', 'danger', 'ghost', 'success'].includes(v)
  },
  size: {
    type: String,
    default: 'md',
    validator: (v) => ['sm', 'md', 'lg'].includes(v)
  },
  type: { type: String, default: 'button' },
  to: { type: [String, Object], default: null },
  disabled: { type: Boolean, default: false },
  loading: { type: Boolean, default: false },
  block: { type: Boolean, default: false }
})

const baseClasses =
  'inline-flex items-center justify-center gap-2 font-medium rounded-lg ' +
  'transition-colors disabled:opacity-50 disabled:cursor-not-allowed ' +
  'focus:outline-none focus-visible:ring-2 focus-visible:ring-offset-2'

const variantClasses = computed(() => ({
  'bg-primary-600 text-white hover:bg-primary-700 focus-visible:ring-primary-500':
    props.variant === 'primary',
  'bg-white border border-slate-300 text-slate-700 hover:bg-slate-50 focus-visible:ring-slate-400':
    props.variant === 'secondary',
  'bg-red-600 text-white hover:bg-red-700 focus-visible:ring-red-500':
    props.variant === 'danger',
  'bg-emerald-600 text-white hover:bg-emerald-700 focus-visible:ring-emerald-500':
    props.variant === 'success',
  'text-slate-600 hover:text-slate-900 hover:bg-slate-100 focus-visible:ring-slate-400':
    props.variant === 'ghost'
}))

const sizeClasses = computed(() => ({
  'text-xs px-3 py-1.5': props.size === 'sm',
  'text-sm px-4 py-2': props.size === 'md',
  'text-base px-5 py-2.5': props.size === 'lg'
}))

const classes = computed(() => [
  baseClasses,
  variantClasses.value,
  sizeClasses.value,
  props.block && 'w-full'
])
</script>

<template>
  <RouterLink v-if="to && !disabled" :to="to" :class="classes">
    <slot />
  </RouterLink>
  <button
    v-else
    :type="type"
    :disabled="disabled || loading"
    :class="classes"
  >
    <span
      v-if="loading"
      class="inline-block h-3.5 w-3.5 rounded-full border-2 border-current border-t-transparent animate-spin"
    />
    <slot />
  </button>
</template>
