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
  <div class="max-w-md mx-auto py-12">
    <div class="text-center mb-8">
      <RouterLink to="/" class="inline-block text-2xl font-bold goth-heading text-fg hover:text-primary transition-colors">
        Job Portal
      </RouterLink>
      <h1 class="mt-6 text-3xl font-bold text-fg">Welcome back</h1>
      <p class="mt-2 text-sm text-fg-muted">Sign in to your account</p>
    </div>

    <AppCard>
      <form @submit.prevent="handleSubmit" class="space-y-6">
        <AppInput
          v-model="email"
          type="email"
          label="Email"
          placeholder="you@example.com"
          required
          autocomplete="email"
        />

        <div>
          <div class="flex items-center justify-between mb-2">
            <label class="block text-sm font-bold text-fg-muted uppercase tracking-wide">
              Password
            </label>
            <RouterLink
              to="/forgot-password"
              class="text-xs font-bold text-primary hover:text-primary-hover transition-colors"
            >
              Forgot password?
            </RouterLink>
          </div>
          <PasswordInput
            v-model="password"
            placeholder="Your password"
            required
            :minlength="1"
            :show-generate="false"
          />
        </div>

        <p v-if="errorMsg" class="text-sm text-red-500">{{ errorMsg }}</p>

        <AppButton
          type="submit"
          :disabled="!canSubmit"
          :loading="loading"
          block
          size="lg"
        >
          {{ loading ? 'Logging in…' : 'Log in' }}
        </AppButton>
      </form>
    </AppCard>

    <p class="mt-8 text-sm text-center text-fg-muted">
      Don't have an account?
      <RouterLink
        to="/register"
        class="text-primary hover:text-primary-hover font-bold transition-colors"
      >
        Register
      </RouterLink>
    </p>
  </div>
</template>
