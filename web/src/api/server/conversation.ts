import { useGlobalStore } from '@/stores/global'

export default {
  getConversations: async (): Promise<void> => {
    const basePath = import.meta.env.VITE_APP_API_BASE_URL + '/backoffice'
    const globalStore = useGlobalStore()
    return fetch(`${basePath}/conversations`, {
      method: 'GET',
      headers: {
        Authorization: `Bearer ${globalStore.accessToken}`,
        'Content-Type': 'application/json',
      },
    })
      .then(async (response) => {
        if (!response.ok) {
          throw new Error('Failed to fetch conversations')
        }
        const data = await response.json()
        console.log(data)
      })
      .catch((error) => {
        console.error('Error fetching conversations:', error)
        throw error
      })
  },
}
