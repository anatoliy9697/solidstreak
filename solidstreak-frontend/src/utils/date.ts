export function dateToISO8601String(date: Date): string {
  const yyyy = date.getFullYear()
  const mm = String(date.getMonth() + 1).padStart(2, '0')
  const dd = String(date.getDate()).padStart(2, '0')
  return `${yyyy}-${mm}-${dd}`
}

export function dateToDDMMYYYYString(date: Date): string {
  const yyyy = date.getFullYear()
  const mm = String(date.getMonth() + 1).padStart(2, '0')
  const dd = String(date.getDate()).padStart(2, '0')
  return `${dd}.${mm}.${yyyy}`
}

export function dateToDDMMYYYYUTCString(date: Date): string {
  const yyyy = date.getUTCFullYear()
  const mm = String(date.getUTCMonth() + 1).padStart(2, '0')
  const dd = String(date.getUTCDate()).padStart(2, '0')
  return `${dd}.${mm}.${yyyy}`
}

export function isSameDay(a: Date, b: Date): boolean {
  return (
    a.getFullYear() === b.getFullYear() &&
    a.getMonth() === b.getMonth() &&
    a.getDate() === b.getDate()
  )
}

export function isBeforeDay(a: Date, b: Date): boolean {
  const aDay = new Date(a.getFullYear(), a.getMonth(), a.getDate())
  const bDay = new Date(b.getFullYear(), b.getMonth(), b.getDate())
  return aDay.getTime() < bDay.getTime()
}

export function isAfterDay(a: Date, b: Date): boolean {
  const aDay = new Date(a.getFullYear(), a.getMonth(), a.getDate())
  const bDay = new Date(b.getFullYear(), b.getMonth(), b.getDate())
  return aDay.getTime() > bDay.getTime()
}

export function parseDate(entry: Date | string): Date {
  if (entry instanceof Date) return entry
  const [year, month, day] = entry.split('-').map(Number)
  return new Date(year!, month! - 1, day)
}

export function parseUTCDate(entry: Date | string): Date {
  if (entry instanceof Date) return entry
  const [year, month, day] = entry.split('-').map(Number)
  return new Date(Date.UTC(year!, month! - 1, day))
}