<script setup>
import { RouterLink } from 'vue-router'
import { Search, Building2, FileText, Send, CheckCircle2, ArrowRight, Users, Sparkles } from 'lucide-vue-next'
import { APP_NAME } from '../config'
import AppButton from '../components/ui/AppButton.vue'

const HEADLINE = 'Hiring, without the noise.'
const SUBHEADLINE = 'A focused job board for people who care about the work. Find a role you\'ll actually love, or hire someone who\'ll make a difference.'
const FOR_SEEKERS = {
  title: 'For job seekers',
  blurb: 'Browse, save, and apply — no spam, no ghost jobs.',
  features: ['Search across real, verified listings', 'Track every application in one place', 'Save jobs and searches to revisit later'],
}
const FOR_EMPLOYERS = {
  title: 'For employers',
  blurb: 'Post a job in minutes. Review applicants without the noise.',
  features: ['Moderated listings — no spam, no scams', 'Applicant profiles with resumes included', 'Simple status tracking, no ATS complexity'],
}
</script>

<template>
  <div class="min-h-screen bg-bg text-fg">
    <!-- Nav -->
    <header class="border-b border-border bg-bg sticky top-0 z-20">
      <nav class="max-w-6xl mx-auto px-4 h-16 flex items-center justify-between">
        <RouterLink to="/" class="text-xl font-bold goth-heading text-fg">
          {{ APP_NAME }}
        </RouterLink>
        <div class="flex items-center gap-4 text-sm font-medium">
          <RouterLink to="/browse" class="text-fg-muted hover:text-primary transition-colors">Browse jobs</RouterLink>
          <RouterLink to="/companies" class="text-fg-muted hover:text-primary transition-colors">Companies</RouterLink>
          <RouterLink to="/login" class="text-fg-muted hover:text-primary transition-colors">Login</RouterLink>
          <RouterLink to="/register" class="bg-primary text-white border border-primary px-5 py-2 uppercase tracking-wider hover:bg-transparent hover:text-primary transition-colors">Get started</RouterLink>
        </div>
      </nav>
    </header>

    <!-- Hero -->
    <section class="max-w-5xl mx-auto px-4 pt-24 pb-24 text-center">
      <div class="inline-flex items-center gap-2 border border-primary text-primary text-xs font-bold uppercase tracking-widest px-3 py-1 mb-8">
        <Sparkles class="w-3.5 h-3.5" />
        A better way to hire
      </div>

      <h1 class="text-5xl md:text-7xl font-bold text-fg leading-tight">
        {{ HEADLINE }}
      </h1>
      <p class="mt-8 text-lg md:text-xl text-fg-muted max-w-2xl mx-auto leading-relaxed">
        {{ SUBHEADLINE }}
      </p>

      <div class="mt-12 flex flex-col sm:flex-row items-center justify-center gap-4">
        <AppButton to="/browse" size="lg">
          <Search class="w-5 h-5" /> Browse jobs
        </AppButton>
        <AppButton to="/register" variant="secondary" size="lg">
          Post a job <ArrowRight class="w-5 h-5" />
        </AppButton>
      </div>
    </section>

    <!-- For whom -->
    <section class="max-w-5xl mx-auto px-4 py-20 border-t border-border">
      <div class="grid md:grid-cols-2 gap-8">
        <div class="bg-surface border border-border p-10">
          <div class="w-14 h-14 border border-primary text-primary flex items-center justify-center mb-6">
            <Users class="w-7 h-7" />
          </div>
          <h2 class="text-2xl font-bold text-fg mb-3">{{ FOR_SEEKERS.title }}</h2>
          <p class="text-fg-muted text-lg">{{ FOR_SEEKERS.blurb }}</p>
          <ul class="mt-8 space-y-4">
            <li v-for="f in FOR_SEEKERS.features" :key="f" class="flex items-start gap-3 text-fg-muted">
              <CheckCircle2 class="w-5 h-5 text-primary flex-shrink-0 mt-0.5" />
              {{ f }}
            </li>
          </ul>
          <div class="mt-8">
            <RouterLink to="/browse" class="inline-flex items-center gap-2 text-sm font-bold uppercase tracking-widest text-primary hover:text-white transition-colors">
              Find your next role <ArrowRight class="w-4 h-4" />
            </RouterLink>
          </div>
        </div>

        <div class="bg-surface border border-border p-10">
          <div class="w-14 h-14 border border-primary text-primary flex items-center justify-center mb-6">
            <Building2 class="w-7 h-7" />
          </div>
          <h2 class="text-2xl font-bold text-fg mb-3">{{ FOR_EMPLOYERS.title }}</h2>
          <p class="text-fg-muted text-lg">{{ FOR_EMPLOYERS.blurb }}</p>
          <ul class="mt-8 space-y-4">
            <li v-for="f in FOR_EMPLOYERS.features" :key="f" class="flex items-start gap-3 text-fg-muted">
              <CheckCircle2 class="w-5 h-5 text-primary flex-shrink-0 mt-0.5" />
              {{ f }}
            </li>
          </ul>
          <div class="mt-8">
            <RouterLink to="/register" class="inline-flex items-center gap-2 text-sm font-bold uppercase tracking-widest text-primary hover:text-white transition-colors">
              Start hiring <ArrowRight class="w-4 h-4" />
            </RouterLink>
          </div>
        </div>
      </div>
    </section>

    <!-- Footer -->
    <footer class="border-t border-border bg-bg">
      <div class="max-w-6xl mx-auto px-4 py-8 flex flex-col md:flex-row items-center justify-between gap-4 text-sm text-fg-muted">
        <p>© {{ new Date().getFullYear() }} {{ APP_NAME }}</p>
        <div class="flex items-center gap-6">
          <RouterLink to="/browse" class="hover:text-primary transition-colors">Browse jobs</RouterLink>
          <RouterLink to="/companies" class="hover:text-primary transition-colors">Companies</RouterLink>
          <RouterLink to="/login" class="hover:text-primary transition-colors">Login</RouterLink>
        </div>
      </div>
    </footer>
  </div>
</template>
