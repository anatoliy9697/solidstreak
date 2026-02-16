export type Theme = 'light' | 'dark'

export function getDefaultTheme(): Theme {
  return 'light'
}

export function toLocalTheme(theme: string | null | undefined): Theme {
  return theme == 'dark' ? 'dark' : 'light'
}

export function getNextTheme(theme: Theme): Theme {
  return theme == 'light' ? 'dark' : 'light'
}
