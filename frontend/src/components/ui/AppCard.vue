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
    class="bg-white rounded-xl border border-slate-200"
    :class="[
      padded && 'p-6',
      interactive && 'transition hover:border-slate-300 hover:shadow-sm cursor-pointer',
      hover && 'transition hover:border-primary-300 hover:shadow-sm'
    ]"
  >
    <slot />
  </component>
</template>
