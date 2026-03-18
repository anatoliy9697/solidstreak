import type { Subscription } from '@/models/subscription'

export interface User {
  id?: number
  tgId: number
  tgUsername?: string
  tgFirstName: string
  tgLastName?: string
  tgLangCode?: string
  langCode?: string
  tgIsBot?: boolean
  subscription?: Subscription
}
