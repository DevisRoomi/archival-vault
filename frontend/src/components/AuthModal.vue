<script setup lang="ts">
import { ref } from 'vue'
import { useAuth } from '../composables/useAuth'

const props = defineProps<{
  isOpen: boolean
  onClose: () => void
}>()

const { login, register, isLoading } = useAuth()
const mode = ref<'LOGIN' | 'REGISTER'>('LOGIN')

const username = ref<string>('')
const email = ref<string>('')
const password = ref<string>('')

const fillDemo = () => {
  username.value = 'archivist'
  password.value = 'password123'
}

const handleSubmit = async () => {
  if (mode.value === 'LOGIN') {
    const ok = await login(username.value, password.value)
    if (ok) props.onClose()
  } else {
    const ok = await register(username.value, email.value, password.value)
    if (ok) props.onClose()
  }
}
</script>

<template>
  <div v-if="isOpen" class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/80 font-mono select-none">
    <div class="wire-box bg-surface w-full max-w-md border-2 border-wire text-raw p-5">
      <!-- Header -->
      <div class="flex items-center justify-between pb-3 mb-4 border-b border-wire">
        <div>
          <span class="text-[10px] text-accent-cyan tracking-widest uppercase">
            [ACCESS_CONTROL_GATEWAY]
          </span>
          <h2 class="text-base font-bold uppercase">
            {{ mode === 'LOGIN' ? 'OPERATOR AUTHENTICATION' : 'IDENTITY REGISTRATION' }}
          </h2>
        </div>
        <button @click="onClose" class="btn-wire text-accent-red font-bold">
          [X]
        </button>
      </div>

      <!-- Mode Toggle -->
      <div class="flex gap-2 mb-4 border-b border-wire pb-2 text-xs">
        <button
          type="button"
          @click="mode = 'LOGIN'"
          class="btn-wire flex-1 py-1"
          :class="mode === 'LOGIN' ? 'bg-wire text-raw font-bold' : 'text-muted'"
        >
          [SIGN IN]
        </button>
        <button
          type="button"
          @click="mode = 'REGISTER'"
          class="btn-wire flex-1 py-1"
          :class="mode === 'REGISTER' ? 'bg-wire text-raw font-bold' : 'text-muted'"
        >
          [REGISTER IDENTITY]
        </button>
      </div>

      <!-- Form -->
      <form @submit.prevent="handleSubmit" class="space-y-3 text-xs">
        <div>
          <label class="block text-muted mb-1 uppercase font-bold">
            Username / Identifier
          </label>
          <input
            v-model="username"
            required
            type="text"
            placeholder="archivist"
            class="input-wire w-full"
          />
        </div>

        <div v-if="mode === 'REGISTER'">
          <label class="block text-muted mb-1 uppercase font-bold">
            Email Address
          </label>
          <input
            v-model="email"
            required
            type="email"
            placeholder="operator@archive.local"
            class="input-wire w-full"
          />
        </div>

        <div>
          <label class="block text-muted mb-1 uppercase font-bold">
            Keyphrase / Password
          </label>
          <input
            v-model="password"
            required
            type="password"
            placeholder="••••••••"
            class="input-wire w-full"
          />
        </div>

        <!-- Quick Demo Fill -->
        <div class="pt-2 flex items-center justify-between">
          <button
            type="button"
            @click="fillDemo"
            class="text-[10px] text-accent-amber underline hover:text-raw"
          >
            [AUTO-FILL DEMO ARCHIVIST]
          </button>
        </div>

        <!-- Submit Button -->
        <div class="pt-3 border-t border-wire flex items-center justify-end gap-2">
          <button type="button" @click="onClose" class="btn-wire text-muted">
            [CANCEL]
          </button>
          <button
            type="submit"
            :disabled="isLoading"
            class="btn-wire font-bold bg-accent-cyan/20 text-accent-cyan hover:bg-accent-cyan hover:text-carbon"
          >
            {{ isLoading ? '[AUTHENTICATING...]' : (mode === 'LOGIN' ? '[VERIFY JWT]' : '[ISSUE IDENTITY]') }}
          </button>
        </div>
      </form>
    </div>
  </div>
</template>
