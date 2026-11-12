<script setup>
import { ref, computed } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useAuthStore } from '../stores/auth'
import PasswordInput from '../components/PasswordInput.vue'
import AppButton from '../components/ui/AppButton.vue'
import AppInput from '../components/ui/AppInput.vue'
import AppCard from '../components/ui/AppCard.vue'

const router = useRouter()
const route = useRoute()
const auth = useAuthStore()

const email = ref('')
const password = ref('')
const loading = ref(false)
const errorMsg = ref('')

const canSubmit = computed(
  () => email.value.trim() && password.value && !loading.value
)

async function handleSubmit() {
  errorMsg.value = ''
  loading.value = true

  try {
    await auth.login({
      email: email.value.trim(),
      password: password.value
    })

    const redirect = route.query.redirect
    router.push(typeof redirect === 'string' ? redirect : { name: 'dashboard' })
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
      <RouterLink
        to="/browse"
        class="inline-flex items-center gap-2 text-lg font-semibold text-slate-900"
      >
        Job Portal
      </RouterLink>
      <h1 class="mt-6 text-2xl font-semibold text-slate-900">Welcome back</h1>
      <p class="mt-1 text-sm text-slate-600">Sign in to your account</p>
    </div>

    <AppCard>
      <form @submit.prevent="handleSubmit" class="space-y-4">
        <AppInput
          v-model="email"
          type="email"
          label="Email"
          placeholder="you@example.com"
          required
          autocomplete="email"
        />

        <div>
          <label class="block text-sm font-medium text-slate-700 mb-1">
            Password
          </label>
          <PasswordInput
            v-model="password"
            placeholder="Your password"
            required
            :minlength="1"
            :show-generate="false"
          />
        </div>

        <p v-if="errorMsg" class="text-sm text-red-600">{{ errorMsg }}</p>

        <AppButton
          type="submit"
          :disabled="!canSubmit"
          :loading="loading"
          block
          size="lg"
        >
          {{ loading ? 'Logging in…' : 'Log in' }}
        </AppButton>
        <div class="text-right">
          <RouterLink
            to="/forgot-password"
            class="text-xs text-primary-600 hover:text-primary-700"
          >
            Forgot password?
          </RouterLink>
        </div>
      </form>
    </AppCard>

    <p class="mt-6 text-sm text-center text-slate-600">
      Don't have an account?
      <RouterLink
        to="/register"
        class="text-primary-600 hover:text-primary-700 font-medium"
      >
        Register
      </RouterLink>
    </p>
  </div>
</template>
