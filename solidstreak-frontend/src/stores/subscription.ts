import { defineStore } from 'pinia'

import type { Plan, Subscription } from '@/models/subscription'
import type { ApiFetcher, RequestResult } from '@/api/request'

export const useSubscriptionStore = defineStore('subscription', {
  state: () => ({
    apiFetcher: null as ApiFetcher | null,
    subPlansMap: new Map<string, Plan>(),
  }),

  actions: {
    async init(apiFetcher: ApiFetcher) {
      this.apiFetcher = apiFetcher
    },

    async fetchSubscriptionPlans(): Promise<RequestResult> {
      const result = await this.apiFetcher!.fetchSubscriptionPlans()

      const data = result.response?.data
      const plans = data ? (data as Plan[]) : []

      this.subPlansMap.clear()

      for (const plan of plans) {
        this.subPlansMap.set(plan.code, plan)
      }

      return result
    },
  },

  getters: {
    planByCode: (state) => (code: string): Plan | undefined => {
      return state.subPlansMap.get(code)
    },

    default(state): Subscription {
      return {
        active: true,
        planCode: 'basic',
        plan: state.subPlansMap.get('basic'),
      } as Subscription
    },
  },
})
