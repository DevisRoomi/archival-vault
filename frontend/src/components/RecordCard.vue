<script setup lang="ts">
import { ref } from 'vue'
import type { RecordItem } from '../types'
import { useToast } from '../composables/useToast'

const props = defineProps<{
  record: RecordItem
  canModify: boolean
  onEdit: (rec: RecordItem) => void
  onDelete: (id: number) => void
}>()

const showPayload = ref<boolean>(false)
const toast = useToast()

const copyCurl = () => {
  const curlCmd = `curl -s http://127.0.0.1:8085/api/records/${props.record.id}`
  navigator.clipboard.writeText(curlCmd)
  toast.info(`[CLIPBOARD] Copied curl command for ${props.record.code}`)
}

const copyPayload = () => {
  navigator.clipboard.writeText(props.record.payload || '')
  toast.info(`[CLIPBOARD] Copied raw payload`)
}

const getClassColor = (c: string) => {
  switch (c) {
    case 'CONFIDENTIAL':
      return 'border-accent-red text-accent-red'
    case 'RESTRICTED':
      return 'border-accent-amber text-accent-amber'
    case 'INTERNAL':
      return 'border-accent-cyan text-accent-cyan'
    case 'PUBLIC':
      return 'border-accent-green text-accent-green'
    default:
      return 'border-wire text-raw'
  }
}
</script>

<template>
  <article class="wire-box bg-surface p-4 font-mono transition-colors hover:border-muted flex flex-col justify-between">
    <div>
      <!-- Header row: Code, Classification, Status -->
      <div class="flex items-start justify-between gap-2 mb-2 pb-2 border-b border-wire/60">
        <div class="flex items-center gap-2 flex-wrap">
          <span class="font-bold text-sm tracking-wide text-accent-cyan">
            [{{ record.code }}]
          </span>
          <span class="text-[11px] px-1.5 py-0.5 border" :class="getClassColor(record.classification)">
            [{{ record.classification }}]
          </span>
          <span class="text-[11px] px-1.5 py-0.5 border border-wire text-muted">
            {{ record.category }}
          </span>
        </div>
        <span class="text-[10px] uppercase tracking-wider px-1.5 py-0.5 border border-wire"
              :class="record.status === 'ACTIVE' ? 'text-accent-green' : record.status === 'DRAFT' ? 'text-accent-amber' : 'text-muted'">
          ● {{ record.status }}
        </span>
      </div>

      <!-- Title -->
      <h3 class="text-base font-bold text-raw leading-snug mb-2 hover:text-accent-cyan transition-colors">
        {{ record.title }}
      </h3>

      <!-- Summary -->
      <p class="text-xs text-muted/90 mb-3 leading-relaxed">
        {{ record.summary || 'No archival summary registered.' }}
      </p>

      <!-- Tags list -->
      <div v-if="record.tags" class="flex flex-wrap gap-1 mb-3">
        <span
          v-for="tag in record.tags.split(',').map(t => t.trim()).filter(Boolean)"
          :key="tag"
          class="text-[10px] px-1 bg-carbon border border-wire text-muted hover:text-raw"
        >
          #{{ tag }}
        </span>
      </div>
    </div>

    <!-- Collapsible Raw Payload Inspector -->
    <div v-if="showPayload" class="my-2 border border-wire bg-carbon p-2 text-xs">
      <div class="flex items-center justify-between pb-1 mb-1 border-b border-wire/40 text-[10px] text-muted">
        <span>RAW_PAYLOAD_VIEWER :: {{ record.code }}</span>
        <button @click="copyPayload" class="hover:text-raw underline">[COPY RAW]</button>
      </div>
      <pre class="font-mono text-[11px] text-accent-green overflow-x-auto p-1 max-h-48 whitespace-pre-wrap select-text">{{ record.payload || '// Empty payload register' }}</pre>
    </div>

    <!-- Footer metadata & action buttons -->
    <div class="pt-2 border-t border-wire/60 flex flex-col sm:flex-row items-start sm:items-center justify-between gap-2 text-xs text-muted">
      <div class="text-[11px] truncate">
        BY: <span class="text-raw">@{{ record.created_by }}</span> |
        <span>{{ new Date(record.created_at).toLocaleDateString() }}</span>
      </div>

      <div class="flex items-center gap-1.5 flex-wrap">
        <button
          @click="showPayload = !showPayload"
          class="btn-wire text-[11px] py-0.5"
          :class="showPayload ? 'bg-wire text-raw' : ''"
        >
          {{ showPayload ? '[HIDE RAW]' : '[VIEW RAW]' }}
        </button>

        <button
          @click="copyCurl"
          class="btn-wire text-[11px] py-0.5 hover:text-accent-cyan"
          title="Copy curl command"
        >
          [CURL]
        </button>

        <button
          v-if="canModify"
          @click="onEdit(record)"
          class="btn-wire text-[11px] py-0.5 text-accent-amber hover:bg-accent-amber hover:text-carbon"
        >
          [EDIT]
        </button>

        <button
          v-if="canModify"
          @click="onDelete(record.id)"
          class="btn-wire text-[11px] py-0.5 text-accent-red hover:bg-accent-red hover:text-carbon"
        >
          [DEL]
        </button>
      </div>
    </div>
  </article>
</template>
