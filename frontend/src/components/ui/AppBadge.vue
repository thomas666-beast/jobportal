<script setup>
const props = defineProps({
  variant: {
    type: String,
    default: 'neutral',
    validator: (v) =>
      ['neutral', 'primary', 'success', 'warning', 'danger', 'info', 'dark'].includes(v)
  },
  size: {
    type: String,
    default: 'md',
    validator: (v) => ['sm', 'md'].includes(v)
  }
})

const variantClasses = {
  neutral: 'bg-slate-100 text-slate-700',
  primary: 'bg-primary-100 text-primary-700',
  success: 'bg-emerald-100 text-emerald-700',
  warning: 'bg-amber-100 text-amber-700',
  danger: 'bg-red-100 text-red-700',
  info: 'bg-blue-100 text-blue-700',
  dark: 'bg-slate-800 text-white'
}

const sizeClasses = {
  sm: 'text-[10px] px-1.5 py-0.5',
  md: 'text-xs px-2 py-1'
}
</script>

<template>
  <span
    class="inline-flex items-center font-medium rounded-full whitespace-nowrap"
    :class="[variantClasses[props.variant], sizeClasses[props.size]]"
  >
    <slot />
  </span>
</template>
