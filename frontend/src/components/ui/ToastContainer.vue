<script setup>
import { CheckCircle2, XCircle, Info, X } from 'lucide-vue-next'
import { useToast } from '../../composables/useToast'

const { toasts, dismiss } = useToast()

function iconFor(type) {
  if (type === 'success') return CheckCircle2
  if (type === 'error') return XCircle
  return Info
}

function classesFor(type) {
  if (type === 'success') return 'bg-emerald-50 border-emerald-200 text-emerald-900'
  if (type === 'error') return 'bg-red-50 border-red-200 text-red-900'
  return 'bg-white border-slate-200 text-slate-900'
}
</script>

<template>
  <div class="fixed bottom-4 right-4 z-50 flex flex-col gap-2 pointer-events-none">
    <TransitionGroup
      enter-active-class="transition duration-200 ease-out"
      enter-from-class="opacity-0 translate-y-2"
      leave-active-class="transition duration-150 ease-in"
      leave-to-class="opacity-0 translate-y-2"
    >
      <div
        v-for="t in toasts"
        :key="t.id"
        class="pointer-events-auto min-w-64 max-w-sm rounded-lg border shadow-sm px-4 py-3 flex items-start gap-3"
        :class="classesFor(t.type)"
      >
        <component
          :is="iconFor(t.type)"
          class="w-5 h-5 mt-0.5 flex-shrink-0"
          :stroke-width="2"
        />
        <p class="text-sm flex-1">{{ t.message }}</p>
        <button
          type="button"
          @click="dismiss(t.id)"
          class="text-current opacity-60 hover:opacity-100"
        >
          <X class="w-4 h-4" />
        </button>
      </div>
    </TransitionGroup>
  </div>
</template>
