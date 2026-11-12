<script setup>
import { RouterLink } from 'vue-router'
</script>

<template>
  <nav class="flex gap-1 border-b border-slate-200 overflow-x-auto">
    <RouterLink
      to="/admin"
      class="px-4 py-2 text-sm font-medium border-b-2 -mb-px transition-colors whitespace-nowrap"
      :class="$route.path === '/admin'
        ? 'border-primary-600 text-primary-700'
        : 'border-transparent text-slate-600 hover:text-slate-900'"
    >
      Dashboard
    </RouterLink>
    <RouterLink
      to="/admin/jobs"
      class="px-4 py-2 text-sm font-medium border-b-2 -mb-px transition-colors whitespace-nowrap"
      :class="$route.path.startsWith('/admin/jobs')
        ? 'border-primary-600 text-primary-700'
        : 'border-transparent text-slate-600 hover:text-slate-900'"
    >
      Jobs
    </RouterLink>
    <RouterLink
      to="/admin/users"
      class="px-4 py-2 text-sm font-medium border-b-2 -mb-px transition-colors whitespace-nowrap"
      :class="$route.path.startsWith('/admin/users')
        ? 'border-primary-600 text-primary-700'
        : 'border-transparent text-slate-600 hover:text-slate-900'"
    >
      Users
    </RouterLink>
    <RouterLink
      to="/admin/applications"
      class="px-4 py-2 text-sm font-medium border-b-2 -mb-px transition-colors whitespace-nowrap"
      :class="$route.path.startsWith('/admin/applications')
        ? 'border-primary-600 text-primary-700'
        : 'border-transparent text-slate-600 hover:text-slate-900'"
    >
      Applications
    </RouterLink>
  </nav>
</template>
