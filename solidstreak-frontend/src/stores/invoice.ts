import { watch } from 'vue'
import { defineStore } from 'pinia'

import type { ApiFetcher, RequestResult } from '@/api/request'
import type { SubscriptionPeriodUnit } from '@/models/subscription'
import { type Currency, INVOICE_STATUS, type Invoice } from '@/models/invoice'
import { useClockStore } from '@/stores/clock'
import { useUserStore } from '@/stores/user'

export const useInvoiceStore = defineStore('invoice', {
  state: () => ({
    apiFetcher: null as ApiFetcher | null,
    invoice: null as Invoice | null,
    lastPolledAt: 0,
    stopWatcher: null as (() => void) | null,
  }),

  actions: {
    async init(apiFetcher: ApiFetcher) {
      this.apiFetcher = apiFetcher
    },

    async createInvoice(
      userId: number,
      subscriptionPlanCode: string,
      subscriptionPeriodUnit: SubscriptionPeriodUnit,
      subscriptionPeriodCount: number,
      currency: Currency,
    ): Promise<RequestResult> {
      const result = await this.apiFetcher!.postInvoice(
        userId,
        subscriptionPlanCode,
        subscriptionPeriodUnit,
        subscriptionPeriodCount,
        currency,
      )

      if (result.success) {
        this.invoice = result.response?.data as Invoice
        this.invoice.expiresAt = new Date(this.invoice.expiresAt)
      } else {
        this.invoice = null
      }

      return result
    },

    async fetchInvoice(userId: number, uuid: string): Promise<RequestResult> {
      const result = await this.apiFetcher!.fetchInvoice(userId, uuid)

      if (result.success) {
        this.invoice = result.response?.data as Invoice
        this.invoice.expiresAt = new Date(this.invoice.expiresAt)
      }

      return result
    },

    async fetchPendingInvoice(userId: number): Promise<RequestResult> {
      const result = await this.apiFetcher!.fetchInvoice(userId, INVOICE_STATUS.PENDING)
      // const result = await this.apiFetcher!.fetchInvoice(
      //   userId,
      //   'af38eb4b-c82e-4559-94fe-37d26577bcbf',
      // )

      if (result.success) {
        this.invoice = result.response?.data as Invoice
        this.invoice.expiresAt = new Date(this.invoice.expiresAt)
      } else {
        this.invoice = null
      }

      return result
    },

    startPollingInvoiceWhilePending() {
      if (!this.pendingInvoice) {
        this.stopPollingInvoice()
        return
      }

      if (this.stopWatcher) return

      const clockStore = useClockStore()
      this.lastPolledAt = clockStore.now

      this.stopWatcher = watch(
        () => clockStore.now,
        async (now) => {
          if (this.invoice && now - this.lastPolledAt >= 1000) {
            await this.fetchInvoice(this.invoice.userId, this.invoice.uuid)
            this.lastPolledAt = now
          }
          if (
            !this.invoice ||
            !this.invoice.active ||
            this.invoice.status !== INVOICE_STATUS.PENDING
          ) {
            this.stopPollingInvoice()
          }
          if (this.invoice?.status === INVOICE_STATUS.PAID) {
            const userStore = useUserStore()
            userStore.refreshSubscription()
          }
        },
      )
    },

    stopPollingInvoice() {
      if (this.stopWatcher) {
        this.stopWatcher()
        this.stopWatcher = null
      }
    },
  },

  getters: {
    pendingInvoice(state): Invoice | null {
      const clockStore = useClockStore()
      return state.invoice &&
        state.invoice.active &&
        state.invoice.expiresAt.getTime() > clockStore.now &&
        state.invoice.status === INVOICE_STATUS.PENDING
        ? state.invoice
        : null
    },

    paidInvoice(state): Invoice | null {
      return state.invoice && state.invoice.active && state.invoice.status === INVOICE_STATUS.PAID
        ? state.invoice
        : null
    },

    millisecondsLeftUntilExpiration(state): number {
      if (!this.pendingInvoice) return 0

      const clockStore = useClockStore()

      return this.pendingInvoice.expiresAt.getTime() - clockStore.now
    },
  },
})
