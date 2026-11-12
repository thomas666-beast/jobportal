<script setup>
import { ref, computed } from 'vue'
import { RouterLink } from 'vue-router'
import { forgotPassword } from '../api/authEmail'
import AppButton from '../components/ui/AppButton.vue'
import AppInput from '../components/ui/AppInput.vue'
import AppCard from '../components/ui/AppCard.vue'

const email = ref('')
const loading = ref(false)
const sent = ref(false)
const errorMsg = ref('')

const canSubmit = computed(() => email.value.trim() && !loading.value)

async function handleSubmit() {
  errorMsg.value = ''
  loading.value = true
  try {
    await forgotPassword(email.value.trim())
    sent.value = true
  } catch (err) {
    errorMsg.value = err.message
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="max-w-md mx-auto py-8">
    <div class="text-center mb-6">
      <h1 class="text-2xl font-semibold text-slate-900">Forgot your password?</h1>
      <p class="mt-1 text-sm text-slate-600">
        Enter your email and we'll send you a reset link.
      </p>
    </div>

    <AppCard>
      <div v-if="sent" class="text-center space-y-3">
        <p class="text-sm text-slate-700">
          If an account exists for <strong>{{ email }}</strong>, we've sent a reset link.
        </p>
        <p class="text-xs text-slate-500">
          Check your inbox (or spam folder). The link expires in 1 hour.
        </p>
        <RouterLink to="/login" class="text-sm text-primary-600 hover:underline">
          Back to login
        </RouterLink>
      </div>

      <form v-else @submit.prevent="handleSubmit" class="space-y-4">
        <AppInput
          v-model="email"
          type="email"
          label="Email"
          placeholder="you@example.com"
          required
          autocomplete="email"
        />

        <p v-if="errorMsg" class="text-sm text-red-600">{{ errorMsg }}</p>

        <AppButton type="submit" :disabled="!canSubmit" :loading="loading" block>
          {{ loading ? 'Sending…' : 'Send reset link' }}
        </AppButton>
      </form>
    </AppCard>

    <p class="mt-6 text-sm text-center text-slate-600">
      Remembered it?
      <RouterLink to="/login" class="text-primary-600 hover:text-primary-700 font-medium">
        Log in
      </RouterLink>
    </p>
  </div>
</template>
