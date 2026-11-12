<script setup>
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { CheckCircle2, XCircle, Loader2 } from 'lucide-vue-next'
import { verifyEmail } from '../api/authEmail'
import AppButton from '../components/ui/AppButton.vue'
import AppCard from '../components/ui/AppCard.vue'

const route = useRoute()
const router = useRouter()

const status = ref('loading') // 'loading' | 'success' | 'error'
const errorMsg = ref('')

onMounted(async () => {
  const token = route.query.token
  if (typeof token !== 'string' || !token) {
    status.value = 'error'
    errorMsg.value = 'Missing verification token.'
    return
  }
  try {
    await verifyEmail(token)
    status.value = 'success'
  } catch (err) {
    status.value = 'error'
    errorMsg.value = err.message
  }
})
</script>

<template>
  <div class="max-w-md mx-auto py-8">
    <AppCard>
      <div v-if="status === 'loading'" class="text-center py-8 space-y-3">
        <Loader2 class="w-10 h-10 text-primary-600 animate-spin mx-auto" />
        <p class="text-sm text-slate-700">Verifying your email…</p>
      </div>

      <div v-else-if="status === 'success'" class="text-center py-8 space-y-3">
        <CheckCircle2 class="w-12 h-12 text-emerald-500 mx-auto" />
        <h1 class="text-lg font-semibold text-slate-900">Email verified</h1>
        <p class="text-sm text-slate-600">
          Your account is now fully activated.
        </p>
        <div class="pt-2">
          <AppButton @click="router.push({ name: 'dashboard' })">
            Go to dashboard
          </AppButton>
        </div>
      </div>

      <div v-else class="text-center py-8 space-y-3">
        <XCircle class="w-12 h-12 text-red-500 mx-auto" />
        <h1 class="text-lg font-semibold text-slate-900">Verification failed</h1>
        <p class="text-sm text-slate-600">{{ errorMsg }}</p>
        <div class="pt-2 flex gap-2 justify-center">
          <AppButton variant="secondary" @click="router.push({ name: 'dashboard' })">
            Go to dashboard
          </AppButton>
          <AppButton @click="router.push({ name: 'login' })">Log in</AppButton>
        </div>
      </div>
    </AppCard>
  </div>
</template>
