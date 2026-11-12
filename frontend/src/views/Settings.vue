<script setup>
import { ref, computed } from 'vue'
import { useAuthStore } from '../stores/auth'
import { changePassword } from '../api/me'
import PasswordInput from '../components/PasswordInput.vue'
import AppButton from '../components/ui/AppButton.vue'
import AppCard from '../components/ui/AppCard.vue'
import AppBadge from '../components/ui/AppBadge.vue'
import { useToast } from '../composables/useToast'

const auth = useAuthStore()
const toast = useToast()

const currentPassword = ref('')
const newPassword = ref('')
const confirmPassword = ref('')
const saving = ref(false)
const errorMsg = ref('')

const passwordsMatch = computed(
  () => !confirmPassword.value || newPassword.value === confirmPassword.value
)

const canSubmit = computed(
  () =>
    currentPassword.value &&
    newPassword.value.length >= 8 &&
    newPassword.value === confirmPassword.value &&
    newPassword.value !== currentPassword.value &&
    !saving.value
)

const roleBadge = computed(() => {
  const r = auth.user?.role
  if (r === 'job_seeker') return { label: 'Job seeker', variant: 'primary' }
  if (r === 'employer') return { label: 'Employer', variant: 'success' }
  if (r === 'admin') return { label: 'Admin', variant: 'warning' }
  return { label: r || '', variant: 'neutral' }
})

async function handleSubmit() {
  errorMsg.value = ''
  saving.value = true
  try {
    await changePassword({
      current_password: currentPassword.value,
      new_password: newPassword.value
    })
    currentPassword.value = ''
    newPassword.value = ''
    confirmPassword.value = ''
    toast.success('Password updated.')
  } catch (err) {
    errorMsg.value = err.message
  } finally {
    saving.value = false
  }
}
</script>

<template>
  <div class="max-w-xl mx-auto space-y-6">
    <div>
      <h1 class="text-2xl font-semibold text-slate-900">Settings</h1>
      <p class="text-sm text-slate-600 mt-1">Manage your account</p>
    </div>

    <AppCard>
      <div class="flex items-start justify-between gap-4">
        <div>
          <p class="text-sm text-slate-500">Signed in as</p>
          <p class="text-base font-medium text-slate-900 mt-0.5">
            {{ auth.user?.username }}
          </p>
        </div>
        <AppBadge :variant="roleBadge.variant">{{ roleBadge.label }}</AppBadge>
      </div>
    </AppCard>

    <AppCard>
      <div class="mb-4">
        <h2 class="text-base font-semibold text-slate-900">Change password</h2>
        <p class="text-sm text-slate-600 mt-1">
          You'll stay signed in on this device.
        </p>
      </div>

      <form @submit.prevent="handleSubmit" class="space-y-4">
        <div>
          <label class="block text-sm font-medium text-slate-700 mb-1">
            Current password
          </label>
          <PasswordInput
            v-model="currentPassword"
            placeholder="Your current password"
            required
            :minlength="1"
            :show-generate="false"
          />
        </div>

        <div>
          <label class="block text-sm font-medium text-slate-700 mb-1">
            New password
          </label>
          <PasswordInput
            v-model="newPassword"
            placeholder="At least 8 characters"
            required
            :minlength="8"
            :show-generate="true"
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

        <AppButton
          type="submit"
          :disabled="!canSubmit"
          :loading="saving"
          block
        >
          {{ saving ? 'Updating…' : 'Update password' }}
        </AppButton>
      </form>
    </AppCard>
  </div>
</template>
