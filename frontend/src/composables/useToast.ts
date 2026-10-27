import { ref } from 'vue'

export interface ToastMessage {
  id: number
  type: 'INFO' | 'SUCCESS' | 'WARN' | 'ERROR'
  text: string
}

const toasts = ref<ToastMessage[]>([])
let counter = 0

export function useToast() {
  const show = (text: string, type: ToastMessage['type'] = 'INFO', duration = 3500) => {
    const id = ++counter
    toasts.value.push({ id, type, text })
    setTimeout(() => {
      toasts.value = toasts.value.filter(t => t.id !== id)
    }, duration)
  }

  const remove = (id: number) => {
    toasts.value = toasts.value.filter(t => t.id !== id)
  }

  return {
    toasts,
    show,
    remove,
    info: (msg: string) => show(msg, 'INFO'),
    success: (msg: string) => show(msg, 'SUCCESS'),
    warn: (msg: string) => show(msg, 'WARN'),
    error: (msg: string) => show(msg, 'ERROR')
  }
}
