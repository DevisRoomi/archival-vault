<script setup lang="ts">
import type { StatsData } from '../types'

defineProps<{
  stats: StatsData | undefined
  isLoading: boolean
}>()
</script>

<template>
  <section class="grid grid-cols-2 lg:grid-cols-4 gap-3 font-mono text-xs select-none">
    <!-- Stat 1: Total Records -->
    <div class="wire-box p-3 bg-surface">
      <div class="text-muted uppercase text-[10px] tracking-wider mb-1 flex justify-between">
        <span>[METRIC:01]</span>
        <span class="text-accent-green">INDEXED</span>
      </div>
      <div class="text-2xl font-bold text-raw font-mono tracking-tight">
        {{ isLoading ? '...' : (stats?.total_records ?? 0).toString().padStart(4, '0') }}
      </div>
      <div class="text-[11px] text-muted mt-1 truncate">
        LATEST: {{ stats?.latest_code || 'ARC-INIT' }}
      </div>
    </div>

    <!-- Stat 2: Confidential & Restricted -->
    <div class="wire-box p-3 bg-surface">
      <div class="text-muted uppercase text-[10px] tracking-wider mb-1 flex justify-between">
        <span>[SECURITY]</span>
        <span class="text-accent-red">GUARDED</span>
      </div>
      <div class="text-lg font-bold flex items-baseline gap-2">
        <span class="text-accent-red">CONF: {{ stats?.by_classification?.CONFIDENTIAL || 0 }}</span>
        <span class="text-muted">/</span>
        <span class="text-accent-amber">REST: {{ stats?.by_classification?.RESTRICTED || 0 }}</span>
      </div>
      <div class="text-[11px] text-muted mt-1">
        HIGH AUDIT CLASSIFICATION
      </div>
    </div>

    <!-- Stat 3: Internal & Public -->
    <div class="wire-box p-3 bg-surface">
      <div class="text-muted uppercase text-[10px] tracking-wider mb-1 flex justify-between">
        <span>[DISTRIB]</span>
        <span class="text-accent-cyan">OPEN/INT</span>
      </div>
      <div class="text-lg font-bold flex items-baseline gap-2">
        <span class="text-accent-cyan">INT: {{ stats?.by_classification?.INTERNAL || 0 }}</span>
        <span class="text-muted">/</span>
        <span class="text-accent-green">PUB: {{ stats?.by_classification?.PUBLIC || 0 }}</span>
      </div>
      <div class="text-[11px] text-muted mt-1">
        GENERAL ARCHIVAL POOL
      </div>
    </div>

    <!-- Stat 4: System Health & Uptime -->
    <div class="wire-box p-3 bg-surface">
      <div class="text-muted uppercase text-[10px] tracking-wider mb-1 flex justify-between">
        <span>[TELEMETRY]</span>
        <span class="text-accent-green">NOMINAL</span>
      </div>
      <div class="text-lg font-bold text-accent-green flex items-center gap-1">
        <span class="inline-block w-2 h-2 bg-accent-green"></span>
        <span>200 OK</span>
        <span class="text-xs text-muted font-normal ml-2">UP: {{ stats?.uptime_seconds ?? 0 }}s</span>
      </div>
      <div class="text-[11px] text-muted mt-1 truncate">
        HOST: FLYENV MYSQL / 3306
      </div>
    </div>
  </section>
</template>
