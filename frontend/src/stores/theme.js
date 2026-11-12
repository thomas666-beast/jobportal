import { defineStore } from 'pinia'
import { ref, computed } from 'vue'

export const useThemeStore = defineStore('theme', () => {
  // Default to 'light' if nothing is saved
  const current = ref(localStorage.getItem('theme') || 'light')

  function apply(theme) {
    document.documentElement.setAttribute('data-theme', theme)
  }

  function setTheme(theme) {
    current.value = theme
    localStorage.setItem('theme', theme)
    apply(theme)
  }

  function toggle() {
    setTheme(current.value === 'light' ? 'goth' : 'light')
  }

  function init() {
    apply(current.value)
  }

  const isGoth = computed(() => current.value === 'goth')

  return { current, isGoth, setTheme, toggle, init }
})
