<script setup>
const props = defineProps({
  padded: { type: Boolean, default: true },
  interactive: { type: Boolean, default: false },
  hover: { type: Boolean, default: false },
  as: { type: String, default: 'div' }
})
</script>

<template>
  <component
    :is="props.as"
    class="bg-surface border border-border"
    :class="[
      padded && 'p-6',
      (interactive || hover) && 'transition-all duration-200 hover:border-primary'
    ]"
  >
    <slot />
  </component>
</template>
