import { defineStore } from 'pinia'

export const useUserStore = defineStore('user', {
  state: () => ({
    token: '' as string,
  }),

  actions: {
    setToken(t: string) {
      this.token = t
      localStorage.setItem('token', t)
    },

    loadFromStorage() {
      const t = localStorage.getItem('token')
      if (t) this.token = t
    }
  }
})
