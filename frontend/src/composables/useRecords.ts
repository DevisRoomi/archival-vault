import { computed, type Ref } from 'vue'
import { useQuery, useMutation, useQueryClient } from '@tanstack/vue-query'
import { request } from '../api/client'
import type { RecordItem, RecordPayload, StatsData, FilterState } from '../types'
import { useToast } from './useToast'

export function useRecords(filter: Ref<FilterState>) {
  const queryClient = useQueryClient()
  const toast = useToast()

  // TanStack Query: Fetch Records
  const recordsQuery = useQuery({
    queryKey: computed(() => ['records', filter.value]),
    queryFn: async () => {
      const p = new URLSearchParams()
      if (filter.value.search) p.set('search', filter.value.search)
      if (filter.value.category && filter.value.category !== 'ALL') p.set('category', filter.value.category)
      if (filter.value.classification && filter.value.classification !== 'ALL') p.set('classification', filter.value.classification)
      if (filter.value.status && filter.value.status !== 'ALL') p.set('status', filter.value.status)

      const qs = p.toString()
      return request<RecordItem[]>(`/records${qs ? `?${qs}` : ''}`)
    },
    staleTime: 1000 * 30 // 30s cache freshness
  })

  // TanStack Query: Fetch Telemetry Stats
  const statsQuery = useQuery({
    queryKey: ['stats'],
    queryFn: () => request<StatsData>('/stats'),
    refetchInterval: 15000 // auto poll every 15s
  })

  // TanStack Mutation: Create Record
  const createMutation = useMutation({
    mutationFn: (payload: RecordPayload) =>
      request<RecordItem>('/records', {
        method: 'POST',
        body: JSON.stringify(payload)
      }),
    onSuccess: (data) => {
      queryClient.invalidateQueries({ queryKey: ['records'] })
      queryClient.invalidateQueries({ queryKey: ['stats'] })
      toast.success(`[COMMIT] Record ${data.code} indexed into vault ledger`)
    },
    onError: (err: any) => {
      toast.error(`[COMMIT: ERR] ${err.message}`)
    }
  })

  // TanStack Mutation: Update Record
  const updateMutation = useMutation({
    mutationFn: ({ id, payload }: { id: number; payload: RecordPayload }) =>
      request<RecordItem>(`/records/${id}`, {
        method: 'PUT',
        body: JSON.stringify(payload)
      }),
    onSuccess: (data) => {
      queryClient.invalidateQueries({ queryKey: ['records'] })
      queryClient.invalidateQueries({ queryKey: ['stats'] })
      toast.success(`[UPDATE] Record ${data.code} patch applied`)
    },
    onError: (err: any) => {
      toast.error(`[UPDATE: ERR] ${err.message}`)
    }
  })

  // TanStack Mutation: Delete Record
  const deleteMutation = useMutation({
    mutationFn: (id: number) =>
      request<{ deleted_id: number }>(`/records/${id}`, {
        method: 'DELETE'
      }),
    onSuccess: (_, id) => {
      queryClient.invalidateQueries({ queryKey: ['records'] })
      queryClient.invalidateQueries({ queryKey: ['stats'] })
      toast.warn(`[PURGE] Record ID #${id} expunged from ledger`)
    },
    onError: (err: any) => {
      toast.error(`[PURGE: ERR] ${err.message}`)
    }
  })

  return {
    records: computed(() => recordsQuery.data.value || []),
    isLoading: recordsQuery.isLoading,
    isError: recordsQuery.isError,
    error: recordsQuery.error,
    refetchRecords: recordsQuery.refetch,

    stats: computed(() => statsQuery.data.value),
    isStatsLoading: statsQuery.isLoading,
    refetchStats: statsQuery.refetch,

    createRecord: createMutation.mutateAsync,
    isCreating: createMutation.isPending,

    updateRecord: updateMutation.mutateAsync,
    isUpdating: updateMutation.isPending,

    deleteRecord: deleteMutation.mutateAsync,
    isDeleting: deleteMutation.isPending
  }
}
