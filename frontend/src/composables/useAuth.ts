import { ref, computed } from 'vue'
import { request } from '../api/client'
import type { User, AuthResponse } from '../types'
import { useToast } from './useToast'

const token = ref<string | null>(localStorage.getItem('vault_token'))
const user = ref<User | null>(
  localStorage.getItem('vault_user') ? JSON.parse(localStorage.getItem('vault_user')!) : null
)
const isLoading = ref<boolean>(false)

export function useAuth() {
  const toast = useToast()

  const isAuthenticated = computed(() => !!token.value && !!user.value)

  const login = async (username: string, password: string): Promise<boolean> => {
    isLoading.value = true
    try {
      const res = await request<AuthResponse>('/auth/login', {
        method: 'POST',
        body: JSON.stringify({ username, password })
      })

      token.value = res.token
      user.value = res.user
      localStorage.setItem('vault_token', res.token)
      localStorage.setItem('vault_user', JSON.stringify(res.user))
      toast.success(`[AUTH] Session established: @${res.user.username} [${res.user.role}]`)
      return true
    } catch (err: any) {
      toast.error(`[AUTH: FAIL] ${err.message}`)
      return false
    } finally {
      isLoading.value = false
    }
  }

  const register = async (username: string, email: string, password: string): Promise<boolean> => {
    isLoading.value = true
    try {
      const res = await request<AuthResponse>('/auth/register', {
        method: 'POST',
        body: JSON.stringify({ username, email, password })
      })

      token.value = res.token
      user.value = res.user
      localStorage.setItem('vault_token', res.token)
      localStorage.setItem('vault_user', JSON.stringify(res.user))
      toast.success(`[AUTH] Identity registered & token generated: @${res.user.username}`)
      return true
    } catch (err: any) {
      toast.error(`[AUTH: FAIL] ${err.message}`)
      return false
    } finally {
      isLoading.value = false
    }
  }

  const logout = () => {
    const prev = user.value?.username || 'user'
    token.value = null
    user.value = null
    localStorage.removeItem('vault_token')
    localStorage.removeItem('vault_user')
    toast.info(`[AUTH] Token revoked. Session ended for @${prev}`)
  }

  const fetchProfile = async () => {
    if (!token.value) return
    try {
      const profile = await request<User>('/auth/me')
      user.value = profile
      localStorage.setItem('vault_user', JSON.stringify(profile))
    } catch {
      logout()
    }
  }

  return {
    token,
    user,
    isLoading,
    isAuthenticated,
    login,
    register,
    logout,
    fetchProfile
  }
}
