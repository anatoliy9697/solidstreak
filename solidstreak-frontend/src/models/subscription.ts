export const SUBSCRIPTION_PERIOD_UNIT = {
  MONTH: 'month',
  YEAR: 'year',
  LIFETIME: 'lifetime',
} as const

export type SubscriptionPeriodUnit =
  (typeof SUBSCRIPTION_PERIOD_UNIT)[keyof typeof SUBSCRIPTION_PERIOD_UNIT]

export const CURRENCY = {
  XTR: 'XTR',
} as const

export type Currency = (typeof CURRENCY)[keyof typeof CURRENCY]

export interface Pricing {
  periodUnit: SubscriptionPeriodUnit
  periodCount: number
  price: number
  currency: Currency
  displayOrder: number
}

export interface Plan {
  code: string
  pricing: Pricing[]
  activeHabitsLimit: number
  showAds: boolean
}

export interface Subscription {
  id?: number
  active: boolean
  planCode: string
  plan: Plan
  startDt?: Date
  finishDt?: Date
  createdAt?: Date
}
