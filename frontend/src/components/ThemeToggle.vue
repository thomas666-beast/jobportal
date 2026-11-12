<script setup>
import { Sun, Moon } from 'lucide-vue-next'
import { useThemeStore } from '../stores/theme'

const theme = useThemeStore()
</script>

<template>
  <button
    type="button"
    @click="theme.toggle()"
    class="p-2 transition-colors text-fg-muted hover:text-fg hover:bg-surface-muted"
    :title="theme.current === 'light' ? 'Switch to Goth' : 'Switch to Light'"
  >
    <Sun v-if="theme.isGoth" class="w-5 h-5" />
    <Moon v-else class="w-5 h-5" />
  </button>
</template>
