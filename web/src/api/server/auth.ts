import type { LoginResponse } from '@/models/login-response'
import type { Res } from '@/models/res'
import { useGlobalStore } from '@/stores/global'

export default {
  login: async (email: string, password: string): Promise<void> => {
    const globalStore = useGlobalStore()
    const basePath = import.meta.env.VITE_APP_API_BASE_URL
    return fetch(`${basePath}/auth/login`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ email, password }),
    })
      .then(async (response) => {
        if (!response.ok) {
          throw new Error('Login failed')
        }
        const data: Res<LoginResponse> = await response.json()
        console.log(data)
        globalStore.setToken(data.data?.access_token ?? '')
      })
      .catch((error) => {
        console.error('Error during login:', error)
        throw error
      })
  },
}
