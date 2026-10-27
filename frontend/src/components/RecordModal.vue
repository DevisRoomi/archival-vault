<script setup lang="ts">
import { ref, watch } from 'vue'
import type { RecordItem, RecordPayload, Classification, Category, Status } from '../types'

const props = defineProps<{
  isOpen: boolean
  recordToEdit: RecordItem | null
  isLoading: boolean
  onClose: () => void
  onSubmit: (payload: RecordPayload, id?: number) => Promise<void>
}>()

const form = ref<RecordPayload>({
  title: '',
  classification: 'INTERNAL',
  category: 'SPEC',
  status: 'ACTIVE',
  tags: '',
  summary: '',
  payload: ''
})

watch(() => props.recordToEdit, (rec) => {
  if (rec) {
    form.value = {
      title: rec.title,
      classification: rec.classification,
      category: rec.category,
      status: rec.status,
      tags: rec.tags,
      summary: rec.summary,
      payload: rec.payload
    }
  } else {
    form.value = {
      title: '',
      classification: 'INTERNAL',
      category: 'SPEC',
      status: 'ACTIVE',
      tags: '',
      summary: '',
      payload: ''
    }
  }
}, { immediate: true })

const handleSubmit = async () => {
  if (!form.value.title.trim()) return
  await props.onSubmit(form.value, props.recordToEdit?.id)
  props.onClose()
}
</script>

<template>
  <div v-if="isOpen" class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/80 font-mono">
    <div class="wire-box bg-surface w-full max-w-2xl border-2 border-wire text-raw p-5 max-h-[90vh] overflow-y-auto">
      <!-- Modal Header -->
      <div class="flex items-center justify-between pb-3 mb-4 border-b border-wire">
        <div>
          <span class="text-[10px] text-accent-cyan tracking-widest uppercase">
            {{ recordToEdit ? `[PATCH RECORD: ${recordToEdit.code}]` : '[INGEST NEW RECORD]' }}
          </span>
          <h2 class="text-base font-bold uppercase">
            {{ recordToEdit ? 'UPDATE ARCHIVAL ENTRY' : 'REGISTER TECHNICAL ARTIFACT' }}
          </h2>
        </div>
        <button @click="onClose" class="btn-wire text-accent-red font-bold">
          [ESC / CLOSE]
        </button>
      </div>

      <!-- Form Body -->
      <form @submit.prevent="handleSubmit" class="space-y-4 text-xs">
        <div>
          <label class="block text-muted mb-1 uppercase tracking-wider font-bold">
            Record Title *
          </label>
          <input
            v-model="form.title"
            required
            type="text"
            placeholder="e.g. Distributed Consensus Timing Specification"
            class="input-wire w-full"
          />
        </div>

        <div class="grid grid-cols-1 sm:grid-cols-3 gap-3">
          <div>
            <label class="block text-muted mb-1 uppercase tracking-wider">Classification</label>
            <select v-model="form.classification" class="input-wire w-full">
              <option value="PUBLIC">PUBLIC</option>
              <option value="INTERNAL">INTERNAL</option>
              <option value="RESTRICTED">RESTRICTED</option>
              <option value="CONFIDENTIAL">CONFIDENTIAL</option>
            </select>
          </div>

          <div>
            <label class="block text-muted mb-1 uppercase tracking-wider">Category</label>
            <select v-model="form.category" class="input-wire w-full">
              <option value="SPEC">SPEC</option>
              <option value="INCIDENT">INCIDENT</option>
              <option value="SCHEMA">SCHEMA</option>
              <option value="DISPATCH">DISPATCH</option>
            </select>
          </div>

          <div>
            <label class="block text-muted mb-1 uppercase tracking-wider">Status</label>
            <select v-model="form.status" class="input-wire w-full">
              <option value="ACTIVE">ACTIVE</option>
              <option value="DRAFT">DRAFT</option>
              <option value="ARCHIVED">ARCHIVED</option>
            </select>
          </div>
        </div>

        <div>
          <label class="block text-muted mb-1 uppercase tracking-wider">Tags (comma-separated)</label>
          <input
            v-model="form.tags"
            type="text"
            placeholder="e.g. golang, vue3, security, network"
            class="input-wire w-full"
          />
        </div>

        <div>
          <label class="block text-muted mb-1 uppercase tracking-wider">Summary / Abstract</label>
          <textarea
            v-model="form.summary"
            rows="2"
            placeholder="Brief technical abstract of this entry..."
            class="input-wire w-full"
          ></textarea>
        </div>

        <div>
          <label class="block text-muted mb-1 uppercase tracking-wider">
            Raw Payload (JSON / Markdown / Plaintext)
          </label>
          <textarea
            v-model="form.payload"
            rows="6"
            placeholder="Paste raw config, schema, or system logs..."
            class="input-wire w-full font-mono text-[11px]"
          ></textarea>
        </div>

        <!-- Form Actions -->
        <div class="pt-3 border-t border-wire flex items-center justify-end gap-2">
          <button type="button" @click="onClose" class="btn-wire text-muted">
            [CANCEL]
          </button>
          <button
            type="submit"
            :disabled="isLoading"
            class="btn-wire font-bold bg-accent-green/20 text-accent-green hover:bg-accent-green hover:text-carbon"
          >
            {{ isLoading ? '[WRITING...]' : (recordToEdit ? '[APPLY PATCH]' : '[COMMIT TO LEDGER]') }}
          </button>
        </div>
      </form>
    </div>
  </div>
</template>
