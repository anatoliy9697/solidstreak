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

export const THEMES: Record<Theme, { [key: string]: string }> = {
  light: {
    // '--background': '#fff',
    // '--text': '#222',
    // ...другие CSS-переменные
  },
  dark: {
    // '--background': '#181818',
    // '--text': '#eee',
    // ...другие CSS-переменные
  },
}
