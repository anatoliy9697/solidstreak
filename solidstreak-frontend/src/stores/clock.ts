import { defineStore } from 'pinia'

export const useClockStore = defineStore('clock', {
  state: () => ({
    now: Date.now(),
    timerId: null as ReturnType<typeof setInterval> | null,
  }),

  actions: {
    start() {
      if (this.timerId !== null) return
      this.now = Date.now()
      this.timerId = setInterval(() => {
        this.now = Date.now()
      }, 1000)
    },

    stop() {
      if (this.timerId === null) return
      clearInterval(this.timerId!)
      this.timerId = null
    },
  },
})
