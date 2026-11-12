<script setup>
import { computed } from 'vue'
import { ChevronLeft, ChevronRight } from 'lucide-vue-next'

const props = defineProps({
  page: { type: Number, required: true },       // 1-based
  pageSize: { type: Number, required: true },
  total: { type: Number, required: true },
  loading: { type: Boolean, default: false }
})

const emit = defineEmits(['update:page'])

const totalPages = computed(() => Math.max(1, Math.ceil(props.total / props.pageSize)))
const from = computed(() => (props.total === 0 ? 0 : (props.page - 1) * props.pageSize + 1))
const to = computed(() => Math.min(props.page * props.pageSize, props.total))

const canPrev = computed(() => props.page > 1 && !props.loading)
const canNext = computed(() => props.page < totalPages.value && !props.loading)

// Pages to show: first, last, current ± 2, with "..." gaps.
const pages = computed(() => {
  const tp = totalPages.value
  const cur = props.page
  const set = new Set([1, tp, cur - 1, cur, cur + 1])
  const items = []
  let prev = 0
  for (const n of Array.from(set).filter((x) => x >= 1 && x <= tp).sort((a, b) => a - b)) {
    if (prev && n - prev > 1) items.push('…')
    items.push(n)
    prev = n
  }
  return items
})

function go(p) {
  if (p < 1 || p > totalPages.value || p === props.page || props.loading) return
  emit('update:page', p)
}
</script>

<template>
  <div class="flex flex-col sm:flex-row items-center justify-between gap-3 mt-4 text-sm">
    <p class="text-slate-500">
      <template v-if="total === 0">No results</template>
      <template v-else>
        Showing <span class="font-medium text-slate-700">{{ from }}–{{ to }}</span>
        of <span class="font-medium text-slate-700">{{ total.toLocaleString() }}</span>
      </template>
    </p>

    <div class="flex items-center gap-1">
      <button
        type="button"
        :disabled="!canPrev"
        @click="go(page - 1)"
        class="p-1.5 rounded-md text-slate-600 hover:bg-slate-100 disabled:opacity-40 disabled:cursor-not-allowed"
        aria-label="Previous page"
      >
        <ChevronLeft class="w-4 h-4" />
      </button>

      <template v-for="(p, i) in pages" :key="i">
        <span v-if="p === '…'" class="px-2 text-slate-400">…</span>
        <button
          v-else
          type="button"
          :disabled="loading"
          @click="go(p)"
          class="min-w-[32px] h-8 px-2 rounded-md text-sm font-medium transition-colors disabled:cursor-not-allowed"
          :class="p === page
            ? 'bg-primary-600 text-white'
            : 'text-slate-700 hover:bg-slate-100'"
        >
          {{ p }}
        </button>
      </template>

      <button
        type="button"
        :disabled="!canNext"
        @click="go(page + 1)"
        class="p-1.5 rounded-md text-slate-600 hover:bg-slate-100 disabled:opacity-40 disabled:cursor-not-allowed"
        aria-label="Next page"
      >
        <ChevronRight class="w-4 h-4" />
      </button>
    </div>
  </div>
</template>
