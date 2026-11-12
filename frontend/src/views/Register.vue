<script setup>
import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'
import PasswordInput from '../components/PasswordInput.vue'
import AppButton from '../components/ui/AppButton.vue'
import AppInput from '../components/ui/AppInput.vue'
import AppCard from '../components/ui/AppCard.vue'
import { useToast } from '../composables/useToast'

const router = useRouter()
const auth = useAuthStore()
const toast = useToast()

const username = ref('')
const email = ref('')
const password = ref('')
const loading = ref(false)
const errorMsg = ref('')

const passwordStrength = computed(() => {
  const p = password.value
  if (!p) return null
  let score = 0
  if (p.length >= 8) score++
  if (p.length >= 12) score++
  if (/[A-Z]/.test(p)) score++
  if (/[0-9]/.test(p)) score++
  if (/[^A-Za-z0-9]/.test(p)) score++
  if (score <= 2) return { label: 'Weak', classes: 'bg-red-500', width: 'w-1/3' }
  if (score === 3) return { label: 'Fair', classes: 'bg-yellow-500', width: 'w-2/3' }
  return { label: 'Strong', classes: 'bg-green-500', width: 'w-full' }
})

const canSubmit = computed(() => username.value.trim() && email.value.trim() && password.value.length >= 8 && !loading.value)

async function handleSubmit() {
  errorMsg.value = ''
  loading.value = true
  try {
    await auth.register({ username: username.value.trim(), email: email.value.trim(), password: password.value })
    toast.success('Account created. Redirecting to login…')
    setTimeout(() => router.push({ name: 'login' }), 800)
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
      <h1 class="mt-6 text-3xl font-bold text-fg">Create your account</h1>
      <p class="mt-2 text-sm text-fg-muted">Join to find jobs or hire talent</p>
    </div>

    <AppCard>
      <form @submit.prevent="handleSubmit" class="space-y-6">
        <AppInput v-model="username" label="Username" placeholder="yourname" required autocomplete="username" />
        <AppInput v-model="email" type="email" label="Email" placeholder="you@example.com" required autocomplete="email" />
        
        <div>
          <label class="block text-sm font-bold text-fg-muted mb-2 uppercase tracking-wide">
            Password <span class="text-primary">*</span>
          </label>
          <PasswordInput v-model="password" placeholder="At least 8 characters" required :minlength="8" />
          
          <div v-if="passwordStrength" class="mt-3">
            <div class="h-1 w-full bg-surface-muted overflow-hidden">
              <div class="h-full transition-all" :class="[passwordStrength.classes, passwordStrength.width]" />
            </div>
            <p class="text-xs text-fg-subtle mt-2">Strength: {{ passwordStrength.label }}</p>
          </div>
        </div>

        <p v-if="errorMsg" class="text-sm text-red-500">{{ errorMsg }}</p>

        <AppButton type="submit" :disabled="!canSubmit" :loading="loading" block size="lg">
          {{ loading ? 'Creating…' : 'Create account' }}
        </AppButton>
      </form>
    </AppCard>

    <p class="mt-8 text-sm text-center text-fg-muted">
      Already have an account?
      <RouterLink to="/login" class="text-primary hover:text-white transition-colors font-bold">Log in</RouterLink>
    </p>
  </div>
</template>
