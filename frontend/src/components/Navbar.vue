<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import { useAuth } from '../composables/useAuth'

defineProps<{
  onOpenCreate: () => void
  onOpenAuth: () => void
}>()

const { user, isAuthenticated, logout } = useAuth()
const timeStr = ref<string>('')

let timer: number
onMounted(() => {
  const update = () => {
    const d = new Date()
    timeStr.value = d.toISOString().replace('T', ' ').substring(0, 19) + ' UTC'
  }
  update()
  timer = window.setInterval(update, 1000)
})

onUnmounted(() => {
  clearInterval(timer)
})
</script>

<template>
  <header class="border-b border-wire bg-surface text-raw font-mono px-4 py-3 select-none">
    <div class="max-w-7xl mx-auto flex flex-col md:flex-row items-start md:items-center justify-between gap-3">
      <!-- Title & Archival System Brand -->
      <div class="flex items-center gap-3">
        <span class="inline-block px-2 py-0.5 border border-wire bg-carbon font-bold text-xs tracking-widest text-accent-cyan">
          [REGISTRY:ARC-SYS]
        </span>
        <div>
          <h1 class="text-sm md:text-base font-bold tracking-tight uppercase flex items-center gap-2">
            <span>ARCHIVAL VAULT</span>
            <span class="text-xs text-muted font-normal">:: RAW MONOSPACED LEDGER</span>
          </h1>
          <div class="text-[11px] text-muted flex items-center gap-2 font-mono">
            <span>{{ timeStr }}</span>
            <span>|</span>
            <span class="text-accent-green">[GOLANG REST]</span>
            <span>|</span>
            <span class="text-accent-blue">[VUE 3 TS]</span>
            <span>|</span>
            <span class="text-accent-amber">[TANSTACK QUERY]</span>
          </div>
        </div>
      </div>

      <!-- Controls & Auth Actions -->
      <div class="flex flex-wrap items-center gap-2 text-xs">
        <button
          v-if="isAuthenticated"
          @click="onOpenCreate"
          class="btn-wire font-bold text-accent-green hover:bg-accent-green hover:text-carbon"
        >
          [+ COMMIT RECORD]
        </button>

        <div v-if="isAuthenticated" class="flex items-center gap-2">
          <span class="px-2 py-1 border border-wire bg-carbon text-accent-cyan">
            @{{ user?.username }} [{{ user?.role }}]
          </span>
          <button
            @click="logout"
            class="btn-wire text-accent-red hover:bg-accent-red hover:text-carbon"
            title="Revoke session JWT"
          >
            [LOGOUT]
          </button>
        </div>

        <button
          v-else
          @click="onOpenAuth"
          class="btn-wire text-accent-amber hover:bg-accent-amber hover:text-carbon font-bold"
        >
          [SIGN IN / AUTHENTICATE]
        </button>
      </div>
    </div>
  </header>
</template>
