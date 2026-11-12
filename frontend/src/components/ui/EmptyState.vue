<script setup>
import { Inbox } from 'lucide-vue-next'

const props = defineProps({
  icon: { type: [Object, Function], default: null },
  title: { type: String, required: true },
  description: { type: String, default: '' }
})
</script>

<template>
  <div class="text-center py-12 px-4">
    <component
      v-if="props.icon"
      :is="props.icon"
      class="w-12 h-12 mx-auto text-slate-300"
      :stroke-width="1.5"
    />
    <Inbox
      v-else
      class="w-12 h-12 mx-auto text-slate-300"
      :stroke-width="1.5"
    />

    <h3 class="mt-4 text-base font-medium text-slate-800">{{ title }}</h3>
    <p v-if="description" class="mt-1 text-sm text-slate-500">
      {{ description }}
    </p>

    <div v-if="$slots.action" class="mt-4">
      <slot name="action" />
    </div>
  </div>
</template>
