export const INVOICE_STATUS = {
  PENDING: 'pending',
  PAID: 'paid',
  CANCELED: 'canceled',
} as const

export type InvoiceStatus = (typeof INVOICE_STATUS)[keyof typeof INVOICE_STATUS]

export const CURRENCY = {
  XTR: 'XTR',
} as const

export type Currency = (typeof CURRENCY)[keyof typeof CURRENCY]

export interface Invoice {
  active: boolean
  uuid: string
  status: InvoiceStatus
  currency: Currency
  amount: number
  userId: number
  expiresAt: Date
  createdAt: Date
  updatedAt: Date
}
