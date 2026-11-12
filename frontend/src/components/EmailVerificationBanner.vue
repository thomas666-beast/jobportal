<script setup>
import { ref } from 'vue'
import { MailWarning, X } from 'lucide-vue-next'
import { useAuthStore } from '../stores/auth'
import { resendVerification } from '../api/authEmail'
import AppButton from './ui/AppButton.vue'
import { useToast } from '../composables/useToast'

const auth = useAuthStore()
const toast = useToast()

const sending = ref(false)
const dismissed = ref(false)

async function resend() {
  sending.value = true
  try {
    await resendVerification()
    toast.success('Verification email sent. Check your inbox.')
  } catch (err) {
    toast.error(err.message)
  } finally {
    sending.value = false
  }
}
</script>

<template>
  <div
    v-if="auth.isAuthenticated && !auth.emailVerified && !dismissed"
    class="bg-amber-50 border border-amber-200 rounded-xl p-4 flex items-start gap-3"
  >
    <MailWarning class="w-5 h-5 text-amber-700 flex-shrink-0 mt-0.5" />
    <div class="flex-1">
      <p class="text-sm font-medium text-amber-900">Please verify your email</p>
      <p class="text-sm text-amber-800 mt-0.5">
        We sent a link to your inbox. Some features — like applying to jobs — require verification.
      </p>
    </div>
    <div class="flex items-center gap-2 flex-shrink-0">
      <AppButton size="sm" variant="primary" :loading="sending" @click="resend">
        {{ sending ? 'Sending…' : 'Resend' }}
      </AppButton>
      <button
        type="button"
        @click="dismissed = true"
        class="text-amber-700 hover:text-amber-900 p-1"
        aria-label="Dismiss"
      >
        <X class="w-4 h-4" />
      </button>
    </div>
  </div>
</template>
