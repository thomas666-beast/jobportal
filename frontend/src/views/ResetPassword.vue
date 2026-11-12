<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { resetPassword } from '../api/authEmail'
import PasswordInput from '../components/PasswordInput.vue'
import AppButton from '../components/ui/AppButton.vue'
import AppCard from '../components/ui/AppCard.vue'
import { useToast } from '../composables/useToast'

const route = useRoute()
const router = useRouter()
const toast = useToast()

const token = ref('')
const newPassword = ref('')
const confirmPassword = ref('')
const loading = ref(false)
const errorMsg = ref('')
const done = ref(false)

const passwordsMatch = computed(
  () => !confirmPassword.value || newPassword.value === confirmPassword.value
)

const canSubmit = computed(
  () =>
    token.value &&
    newPassword.value.length >= 8 &&
    newPassword.value === confirmPassword.value &&
    !loading.value
)

onMounted(() => {
  const t = route.query.token
  if (typeof t === 'string') {
    token.value = t
  } else {
    errorMsg.value = 'Invalid or missing reset token.'
  }
})

async function handleSubmit() {
  errorMsg.value = ''
  loading.value = true
  try {
    await resetPassword(token.value, newPassword.value)
    done.value = true
    toast.success('Password updated.')
    setTimeout(() => router.push({ name: 'login' }), 1200)
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
      <h1 class="text-2xl font-semibold text-slate-900">Set a new password</h1>
    </div>

    <AppCard>
      <div v-if="done" class="text-center space-y-3">
        <p class="text-sm text-slate-700">
          Your password has been updated. Redirecting to login…
        </p>
      </div>

      <form v-else @submit.prevent="handleSubmit" class="space-y-4">
        <div>
          <label class="block text-sm font-medium text-slate-700 mb-1">
            New password
          </label>
          <PasswordInput
            v-model="newPassword"
            placeholder="At least 8 characters"
            required
            :minlength="8"
          />
        </div>

        <div>
          <label class="block text-sm font-medium text-slate-700 mb-1">
            Confirm new password
          </label>
          <PasswordInput
            v-model="confirmPassword"
            placeholder="Repeat new password"
            required
            :minlength="8"
            :show-generate="false"
          />
          <p v-if="!passwordsMatch" class="text-xs text-red-600 mt-1">
            Passwords do not match.
          </p>
        </div>

        <p v-if="errorMsg" class="text-sm text-red-600">{{ errorMsg }}</p>

        <AppButton type="submit" :disabled="!canSubmit" :loading="loading" block>
          {{ loading ? 'Updating…' : 'Set new password' }}
        </AppButton>
      </form>
    </AppCard>

    <p class="mt-6 text-sm text-center text-slate-600">
      <RouterLink to="/login" class="text-primary-600 hover:text-primary-700 font-medium">
        Back to login
      </RouterLink>
    </p>
  </div>
</template>
