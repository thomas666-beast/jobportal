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
  neutral: 'bg-surface-muted text-fg-muted',
  primary: 'bg-primary-soft text-primary-soft-fg',
  success: 'bg-success-soft text-success-fg',
  warning: 'bg-warning-soft text-warning-fg',
  danger: 'bg-danger-soft text-danger-fg',
  info: 'bg-info-soft text-info-fg',
  dark: 'bg-fg text-bg'
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
