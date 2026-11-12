<script setup>
import { RouterLink } from 'vue-router'
import { Search, Building2, FileText, Send, CheckCircle2, ArrowRight, Users, Briefcase, Sparkles } from 'lucide-vue-next'
import { APP_NAME } from '../config'
import AppButton from '../components/ui/AppButton.vue'

// Editable copy — change these strings to rebrand the whole page.
const HEADLINE = 'Hiring, without the noise.'
const SUBHEADLINE =
  'A focused job board for people who care about the work. Find a role you\'ll actually love, or hire someone who\'ll make a difference.'
const FOR_SEEKERS = {
  title: 'For job seekers',
  blurb: 'Browse, save, and apply — no spam, no ghost jobs.',
  features: [
    'Search across real, verified listings',
    'Track every application in one place',
    'Save jobs and searches to revisit later',
  ],
}
const FOR_EMPLOYERS = {
  title: 'For employers',
  blurb: 'Post a job in minutes. Review applicants without the noise.',
  features: [
    'Moderated listings — no spam, no scams',
    'Applicant profiles with resumes included',
    'Simple status tracking, no ATS complexity',
  ],
}
const STATS = [
  { value: '100%', label: 'Verified employers' },
  { value: '1-click', label: 'Apply to jobs' },
  { value: '0', label: 'Hidden fees' },
]
</script>

<template>
  <div class="min-h-screen bg-slate-50">
    <!-- Nav -->
    <header class="border-b border-slate-200 bg-white/80 backdrop-blur sticky top-0 z-20">
      <nav class="max-w-6xl mx-auto px-4 h-14 flex items-center justify-between">
        <RouterLink to="/" class="text-base font-semibold text-slate-900">
          {{ APP_NAME }}
        </RouterLink>
        <div class="flex items-center gap-3 text-sm">
          <RouterLink
            to="/browse"
            class="hidden sm:inline text-slate-600 hover:text-slate-900 px-3 py-1.5"
          >
            Browse jobs
          </RouterLink>
          <RouterLink
            to="/companies"
            class="hidden sm:inline text-slate-600 hover:text-slate-900 px-3 py-1.5"
          >
            Companies
          </RouterLink>
          <RouterLink
            to="/login"
            class="text-slate-600 hover:text-slate-900 px-3 py-1.5"
          >
            Login
          </RouterLink>
          <RouterLink
            to="/register"
            class="text-sm text-white bg-primary-600 hover:bg-primary-700
                   font-medium px-4 py-1.5 rounded-lg transition-colors"
          >
            Get started
          </RouterLink>
        </div>
      </nav>
    </header>

    <!-- Hero -->
    <section class="max-w-5xl mx-auto px-4 pt-16 pb-20 text-center">
      <div class="inline-flex items-center gap-2 bg-primary-50 text-primary-700 text-xs font-medium px-3 py-1 rounded-full mb-6">
        <Sparkles class="w-3.5 h-3.5" />
        A better way to hire and be hired
      </div>

      <h1 class="text-4xl md:text-5xl lg:text-6xl font-semibold text-slate-900 leading-tight tracking-tight">
        {{ HEADLINE }}
      </h1>
      <p class="mt-6 text-lg md:text-xl text-slate-600 max-w-2xl mx-auto leading-relaxed">
        {{ SUBHEADLINE }}
      </p>

      <div class="mt-10 flex flex-col sm:flex-row items-center justify-center gap-3">
        <AppButton
          to="/browse"
          size="lg"
          class="w-full sm:w-auto"
        >
          <Search class="w-4 h-4" />
          Browse jobs
        </AppButton>
        <AppButton
          to="/register"
          variant="secondary"
          size="lg"
          class="w-full sm:w-auto"
        >
          Post a job
          <ArrowRight class="w-4 h-4" />
        </AppButton>
      </div>

      <div class="mt-14 grid grid-cols-3 max-w-lg mx-auto gap-4">
        <div v-for="s in STATS" :key="s.label" class="text-center">
          <div class="text-2xl font-semibold text-slate-900">{{ s.value }}</div>
          <div class="text-xs text-slate-500 mt-1">{{ s.label }}</div>
        </div>
      </div>
    </section>

    <!-- For whom -->
    <section class="max-w-5xl mx-auto px-4 py-16">
      <div class="grid md:grid-cols-2 gap-6">
        <div class="bg-white rounded-2xl border border-slate-200 p-8">
          <div class="w-12 h-12 rounded-xl bg-primary-50 text-primary-600 flex items-center justify-center mb-4">
            <Users class="w-6 h-6" />
          </div>
          <h2 class="text-xl font-semibold text-slate-900">{{ FOR_SEEKERS.title }}</h2>
          <p class="text-slate-600 mt-2">{{ FOR_SEEKERS.blurb }}</p>
          <ul class="mt-5 space-y-3">
            <li
              v-for="f in FOR_SEEKERS.features"
              :key="f"
              class="flex items-start gap-2 text-sm text-slate-700"
            >
              <CheckCircle2 class="w-4 h-4 text-emerald-500 flex-shrink-0 mt-0.5" />
              {{ f }}
            </li>
          </ul>
          <div class="mt-6">
            <RouterLink
              to="/browse"
              class="inline-flex items-center gap-1 text-sm font-medium text-primary-600 hover:text-primary-700"
            >
              Find your next role
              <ArrowRight class="w-3.5 h-3.5" />
            </RouterLink>
          </div>
        </div>

        <div class="bg-white rounded-2xl border border-slate-200 p-8">
          <div class="w-12 h-12 rounded-xl bg-emerald-50 text-emerald-600 flex items-center justify-center mb-4">
            <Building2 class="w-6 h-6" />
          </div>
          <h2 class="text-xl font-semibold text-slate-900">{{ FOR_EMPLOYERS.title }}</h2>
          <p class="text-slate-600 mt-2">{{ FOR_EMPLOYERS.blurb }}</p>
          <ul class="mt-5 space-y-3">
            <li
              v-for="f in FOR_EMPLOYERS.features"
              :key="f"
              class="flex items-start gap-2 text-sm text-slate-700"
            >
              <CheckCircle2 class="w-4 h-4 text-emerald-500 flex-shrink-0 mt-0.5" />
              {{ f }}
            </li>
          </ul>
          <div class="mt-6">
            <RouterLink
              to="/register"
              class="inline-flex items-center gap-1 text-sm font-medium text-primary-600 hover:text-primary-700"
            >
              Start hiring
              <ArrowRight class="w-3.5 h-3.5" />
            </RouterLink>
          </div>
        </div>
      </div>
    </section>

    <!-- How it works -->
    <section class="max-w-5xl mx-auto px-4 py-16">
      <div class="text-center mb-10">
        <h2 class="text-2xl md:text-3xl font-semibold text-slate-900">
          How it works
        </h2>
        <p class="text-slate-600 mt-2 max-w-xl mx-auto">
          Simple on both sides. No hidden steps, no unnecessary friction.
        </p>
      </div>

      <div class="grid md:grid-cols-3 gap-6">
        <div class="text-center">
          <div class="w-14 h-14 rounded-2xl bg-primary-50 text-primary-600 flex items-center justify-center mx-auto mb-4">
            <FileText class="w-7 h-7" />
          </div>
          <h3 class="font-semibold text-slate-900">Build your profile</h3>
          <p class="text-sm text-slate-600 mt-1">
            Add your resume and a short bio. Takes two minutes.
          </p>
        </div>
        <div class="text-center">
          <div class="w-14 h-14 rounded-2xl bg-primary-50 text-primary-600 flex items-center justify-center mx-auto mb-4">
            <Search class="w-7 h-7" />
          </div>
          <h3 class="font-semibold text-slate-900">Find the right job</h3>
          <p class="text-sm text-slate-600 mt-1">
            Search, save, and compare. No spam, no scam listings.
          </p>
        </div>
        <div class="text-center">
          <div class="w-14 h-14 rounded-2xl bg-primary-50 text-primary-600 flex items-center justify-center mx-auto mb-4">
            <Send class="w-7 h-7" />
          </div>
          <h3 class="font-semibold text-slate-900">Apply with one click</h3>
          <p class="text-sm text-slate-600 mt-1">
            Attach a short cover letter if you want. Track every application.
          </p>
        </div>
      </div>
    </section>

    <!-- Final CTA -->
    <section class="max-w-4xl mx-auto px-4 pb-20">
      <div class="bg-gradient-to-br from-primary-600 to-indigo-700 rounded-3xl px-8 py-12 text-center text-white">
        <h2 class="text-2xl md:text-3xl font-semibold">
          Ready to start?
        </h2>
        <p class="mt-3 text-primary-100 max-w-xl mx-auto">
          Whether you're hiring or looking — it takes a minute to get going.
        </p>
        <div class="mt-8 flex flex-col sm:flex-row items-center justify-center gap-3">
          <RouterLink
            to="/register"
            class="w-full sm:w-auto px-6 py-3 rounded-lg bg-white text-primary-700 font-medium hover:bg-primary-50 transition-colors"
          >
            Create your account
          </RouterLink>
          <RouterLink
            to="/browse"
            class="w-full sm:w-auto px-6 py-3 rounded-lg border border-white/30 text-white font-medium hover:bg-white/10 transition-colors"
          >
            Browse jobs
          </RouterLink>
        </div>
      </div>
    </section>

    <!-- Footer -->
    <footer class="border-t border-slate-200 bg-white">
      <div class="max-w-6xl mx-auto px-4 py-6 flex flex-col md:flex-row items-center justify-between gap-3 text-sm text-slate-500">
        <p>© {{ new Date().getFullYear() }} {{ APP_NAME }}</p>
        <div class="flex items-center gap-4">
          <RouterLink to="/browse" class="hover:text-slate-900">Browse jobs</RouterLink>
          <RouterLink to="/companies" class="hover:text-slate-900">Companies</RouterLink>
          <RouterLink to="/login" class="hover:text-slate-900">Login</RouterLink>
        </div>
      </div>
    </footer>
  </div>
</template>
