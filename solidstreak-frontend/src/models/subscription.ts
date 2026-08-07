export const SUBSCRIPTION_PERIOD = {
  MONTH: 'month',
  YEAR: 'year',
  LIFETIME: 'lifetime',
} as const

export type SubscriptionPeriod = (typeof SUBSCRIPTION_PERIOD)[keyof typeof SUBSCRIPTION_PERIOD]

export const CURRENCY = {
  XTR: 'XTR',
} as const

export type Currency = (typeof CURRENCY)[keyof typeof CURRENCY]

export interface Pricing {
  period: SubscriptionPeriod
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
