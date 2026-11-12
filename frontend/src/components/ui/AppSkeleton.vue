<script setup>
const props = defineProps({
  variant: { type: String, default: 'text' },
  lines: { type: Number, default: 3 },
  width: { type: String, default: '' }
})

function lineClass(i) {
  if (i === 0) return 'w-3/4'
  if (i === props.lines - 1) return 'w-1/2'
  return 'w-full'
}
</script>

<template>
  <div v-if="variant === 'text'" class="animate-pulse space-y-2">
    <div v-for="i in lines" :key="i" class="h-3 rounded bg-surface-muted" :class="lineClass(i - 1)" />
  </div>

  <div v-else-if="variant === 'card'" class="animate-pulse bg-surface border border-border p-6 space-y-3">
    <div class="flex items-start gap-4">
      <div class="w-12 h-12 bg-surface-muted flex-shrink-0" />
      <div class="flex-1 space-y-2">
        <div class="h-4 bg-surface-muted w-2/3" />
        <div class="h-3 bg-surface-muted w-1/2" />
      </div>
    </div>
    <div class="h-3 bg-surface-muted w-full" />
    <div class="h-3 bg-surface-muted w-5/6" />
  </div>

  <div v-else-if="variant === 'row'" class="animate-pulse flex items-center gap-4 py-3">
    <div class="w-10 h-10 bg-surface-muted flex-shrink-0" />
    <div class="flex-1 space-y-2">
      <div class="h-3 bg-surface-muted w-1/3" />
      <div class="h-3 bg-surface-muted w-1/2" />
    </div>
  </div>
</template>
