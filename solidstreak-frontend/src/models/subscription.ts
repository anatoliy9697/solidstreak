export interface Pricing {
  tgStarsPerMonth: number
  tgStarsPerYear: number
  tgStarsForever: number
}

export interface Plan {
  code: string
  price: Pricing
  habitsLimit: number
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

export function getDefaultSubscription(): Subscription {
  return {
    active: true,
    planCode: 'basic',
    plan: {
      // TODO: брать планы из бэкенда
      code: 'basic',
      price: {
        tgStarsPerMonth: 0,
        tgStarsPerYear: 0,
        tgStarsForever: 0,
      },
      habitsLimit: 5,
      showAds: true,
    },
  }
}
