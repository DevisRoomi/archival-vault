<script setup lang="ts">
import { ref } from 'vue'
import type { FilterState, RecordItem, RecordPayload } from './types'
import { useAuth } from './composables/useAuth'
import { useRecords } from './composables/useRecords'
import { useToast } from './composables/useToast'
import Navbar from './components/Navbar.vue'
import StatsGrid from './components/StatsGrid.vue'
import RecordCard from './components/RecordCard.vue'
import RecordModal from './components/RecordModal.vue'
import AuthModal from './components/AuthModal.vue'

const { isAuthenticated } = useAuth()
const { toasts, remove: removeToast } = useToast()

const filters = ref<FilterState>({
  search: '',
  category: 'ALL',
  classification: 'ALL',
  status: 'ALL'
})

const {
  records,
  isLoading,
  error,
  stats,
  isStatsLoading,
  refetchRecords,
  createRecord,
  updateRecord,
  deleteRecord,
  isCreating,
  isUpdating
} = useRecords(filters)

const isRecordModalOpen = ref(false)
const isAuthModalOpen = ref(false)
const recordToEdit = ref<RecordItem | null>(null)

const handleOpenCreate = () => {
  recordToEdit.value = null
  isRecordModalOpen.value = true
}

const handleEdit = (rec: RecordItem) => {
  recordToEdit.value = rec
  isRecordModalOpen.value = true
}

const handleDelete = async (id: number) => {
  if (confirm(`Expunge record #${id} from ledger? This action is recorded.`)) {
    await deleteRecord(id)
  }
}

const handleModalSubmit = async (payload: RecordPayload, id?: number) => {
  if (id) {
    await updateRecord({ id, payload })
  } else {
    await createRecord(payload)
  }
}

const resetFilters = () => {
  filters.value = {
    search: '',
    category: 'ALL',
    classification: 'ALL',
    status: 'ALL'
  }
}
</script>

<template>
  <div class="min-h-screen bg-carbon text-raw font-mono flex flex-col">
    <!-- Navbar -->
    <Navbar
      :onOpenCreate="handleOpenCreate"
      :onOpenAuth="() => (isAuthModalOpen = true)"
    />

    <!-- Main Workspace -->
    <main class="flex-1 max-w-7xl w-full mx-auto p-4 md:p-6 space-y-5">
      <!-- ASCII Wireframe Hero Banner -->
      <section class="wire-box p-4 bg-surface text-xs select-none">
        <div class="flex flex-col md:flex-row items-start md:items-center justify-between gap-3 pb-3 border-b border-wire">
          <div>
            <div class="text-[10px] text-accent-cyan tracking-widest uppercase">
              [FULLSTACK ARCHITECTURE DEMO // FIKA RIDAUL MAULAYYA COURSE ALIGNMENT]
            </div>
            <h2 class="text-base md:text-lg font-bold tracking-tight">
              ARCHIVAL RAW &amp; MONOSPACED WIREFRAME LEDGER
            </h2>
          </div>
          <div class="flex flex-wrap gap-1.5 text-[10px]">
            <span class="px-2 py-0.5 border border-wire bg-carbon">[RESTful API]</span>
            <span class="px-2 py-0.5 border border-wire bg-carbon text-accent-amber">[JWT AUTH]</span>
            <span class="px-2 py-0.5 border border-wire bg-carbon text-accent-green">[VUE 3]</span>
            <span class="px-2 py-0.5 border border-wire bg-carbon text-accent-cyan">[TYPESCRIPT]</span>
            <span class="px-2 py-0.5 border border-wire bg-carbon text-accent-red">[COMPOSABLES]</span>
            <span class="px-2 py-0.5 border border-wire bg-carbon text-accent-blue">[TANSTACK QUERY]</span>
          </div>
        </div>

        <p class="text-[11px] text-muted mt-2 leading-relaxed">
          High-assurance technical registry for system blueprints, incident postmortems, and cryptographic claims.
          Backed by Golang RESTful API engine on port 8085, FlyEnv MySQL database, and Vue 3 TypeScript client.
        </p>
      </section>

      <!-- Live Stats Telemetry Grid -->
      <StatsGrid :stats="stats" :isLoading="isStatsLoading" />

      <!-- Filter & Search Toolbar -->
      <section class="wire-box p-3 bg-surface font-mono text-xs space-y-3">
        <div class="flex flex-col md:flex-row items-stretch md:items-center justify-between gap-3">
          <!-- Search input -->
          <div class="flex-1 flex items-center border border-wire bg-carbon px-2 py-1">
            <span class="text-accent-cyan font-bold mr-2">&gt;</span>
            <input
              v-model="filters.search"
              type="text"
              placeholder="Search code, title, tags, or summary..."
              class="w-full bg-transparent border-none text-raw placeholder:text-muted/60 focus:outline-none text-xs"
            />
          </div>

          <!-- Dropdown Filters -->
          <div class="flex flex-wrap items-center gap-2">
            <select v-model="filters.category" class="input-wire py-1">
              <option value="ALL">CAT: ALL</option>
              <option value="SPEC">CAT: SPEC</option>
              <option value="INCIDENT">CAT: INCIDENT</option>
              <option value="SCHEMA">CAT: SCHEMA</option>
              <option value="DISPATCH">CAT: DISPATCH</option>
            </select>

            <select v-model="filters.classification" class="input-wire py-1">
              <option value="ALL">CLASS: ALL</option>
              <option value="PUBLIC">CLASS: PUBLIC</option>
              <option value="INTERNAL">CLASS: INTERNAL</option>
              <option value="RESTRICTED">CLASS: RESTRICTED</option>
              <option value="CONFIDENTIAL">CLASS: CONFIDENTIAL</option>
            </select>

            <select v-model="filters.status" class="input-wire py-1">
              <option value="ALL">STATUS: ALL</option>
              <option value="ACTIVE">STATUS: ACTIVE</option>
              <option value="DRAFT">STATUS: DRAFT</option>
              <option value="ARCHIVED">STATUS: ARCHIVED</option>
            </select>

            <button
              @click="resetFilters"
              class="btn-wire py-1 text-muted hover:text-raw"
              title="Reset query filters"
            >
              [RESET]
            </button>

            <button
              @click="() => refetchRecords()"
              class="btn-wire py-1 text-accent-cyan hover:bg-accent-cyan hover:text-carbon"
              title="TanStack Query Refetch"
            >
              [SYNC]
            </button>
          </div>
        </div>

        <!-- Filter status strip -->
        <div class="flex items-center justify-between text-[11px] text-muted pt-1 border-t border-wire/40">
          <span>
            MATCHED: <strong class="text-raw">{{ records.length }}</strong> ENTRIES
          </span>
          <span class="text-[10px]">
            QUERY_CACHE: TANSTACK_VUE_QUERY_v5
          </span>
        </div>
      </section>

      <!-- Records Grid -->
      <section>
        <!-- Loading State -->
        <div v-if="isLoading" class="wire-box p-8 text-center text-muted text-xs bg-surface animate-pulse">
          [QUERY_PENDING] Streaming archival records from Golang REST endpoint...
        </div>

        <!-- Error State -->
        <div v-else-if="error" class="wire-box p-6 border-accent-red text-accent-red text-xs bg-surface">
          [ERROR] Failed to fetch archival stream: {{ error.message }}
        </div>

        <!-- Empty State -->
        <div v-else-if="records.length === 0" class="wire-box p-8 text-center text-muted text-xs bg-surface space-y-2">
          <p>[LEDGER_EMPTY] No technical artifacts match current filter criteria.</p>
          <button @click="resetFilters" class="btn-wire text-accent-cyan text-xs">
            [CLEAR FILTERS]
          </button>
        </div>

        <!-- Record Cards -->
        <div v-else class="grid grid-cols-1 md:grid-cols-2 gap-4">
          <RecordCard
            v-for="rec in records"
            :key="rec.id"
            :record="rec"
            :canModify="isAuthenticated"
            :onEdit="handleEdit"
            :onDelete="handleDelete"
          />
        </div>
      </section>
    </main>

    <!-- Wireframe Footer -->
    <footer class="border-t border-wire bg-surface text-muted text-[11px] font-mono py-4 px-4 select-none">
      <div class="max-w-7xl mx-auto flex flex-col md:flex-row items-center justify-between gap-2">
        <div>
          ARCHIVAL VAULT v1.0.0 // RAW MONOSPACED WIREFRAME ARCHITECTURE
        </div>
        <div class="flex items-center gap-3">
          <span>PORT: 8085</span>
          <span>•</span>
          <span>HERD: archival-vault.test</span>
          <span>•</span>
          <a href="https://github.com/DevisRoomi/archival-vault" target="_blank" class="text-accent-cyan underline">
            [GITHUB REPOSITORY]
          </a>
        </div>
      </div>
    </footer>

    <!-- Modals -->
    <RecordModal
      :isOpen="isRecordModalOpen"
      :recordToEdit="recordToEdit"
      :isLoading="isCreating || isUpdating"
      :onClose="() => (isRecordModalOpen = false)"
      :onSubmit="handleModalSubmit"
    />

    <AuthModal
      :isOpen="isAuthModalOpen"
      :onClose="() => (isAuthModalOpen = false)"
    />

    <!-- Toast Notifications -->
    <div class="fixed bottom-4 right-4 z-50 space-y-2 max-w-sm w-full select-none">
      <div
        v-for="t in toasts"
        :key="t.id"
        class="wire-box p-2.5 text-xs font-mono shadow-2xl flex items-center justify-between gap-2"
        :class="{
          'bg-carbon text-accent-green border-accent-green': t.type === 'SUCCESS',
          'bg-carbon text-accent-red border-accent-red': t.type === 'ERROR',
          'bg-carbon text-accent-amber border-accent-amber': t.type === 'WARN',
          'bg-carbon text-accent-cyan border-accent-cyan': t.type === 'INFO'
        }"
      >
        <span class="break-all">{{ t.text }}</span>
        <button @click="removeToast(t.id)" class="hover:text-raw px-1">[x]</button>
      </div>
    </div>
  </div>
</template>

<style>
:root {
  --bg-carbon: #090a0c;
  --bg-surface: #111317;
  --border-wire: #262c35;
  --text-raw: #e6edf3;
  --text-muted: #8b949e;
  --accent-cyan: #58a6ff;
  --accent-green: #3fb950;
  --accent-amber: #d29922;
  --accent-red: #f85149;
  --accent-blue: #79c0ff;
}

body {
  margin: 0;
  padding: 0;
  background-color: var(--bg-carbon);
  color: var(--text-raw);
  font-family: 'JetBrains Mono', 'Courier New', Courier, monospace;
  -webkit-font-smoothing: antialiased;
}

.wire-box {
  border: 1px solid var(--border-wire);
  border-radius: 0px !important;
}

.btn-wire {
  border: 1px solid var(--border-wire);
  border-radius: 0px !important;
  background: var(--bg-carbon);
  color: var(--text-raw);
  padding: 4px 10px;
  cursor: pointer;
  font-family: inherit;
  transition: all 0.15s ease-in-out;
}

.btn-wire:hover {
  border-color: var(--text-muted);
}

.input-wire {
  border: 1px solid var(--border-wire);
  border-radius: 0px !important;
  background: var(--bg-carbon);
  color: var(--text-raw);
  padding: 6px 10px;
  font-family: inherit;
  outline: none;
}

.input-wire:focus {
  border-color: var(--accent-cyan);
}

/* Color utility classes */
.bg-carbon { background-color: var(--bg-carbon); }
.bg-surface { background-color: var(--bg-surface); }
.border-wire { border-color: var(--border-wire); }
.text-raw { color: var(--text-raw); }
.text-muted { color: var(--text-muted); }
.text-accent-cyan { color: var(--accent-cyan); }
.text-accent-green { color: var(--accent-green); }
.text-accent-amber { color: var(--accent-amber); }
.text-accent-red { color: var(--accent-red); }
.text-accent-blue { color: var(--accent-blue); }
</style>
