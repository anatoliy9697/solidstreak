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
    '--a-text-decoration-color': '#1e40af', // blue-800
    '--theme-toggle-p-button-bg': '#e5e7eb', // gray-200
    '--theme-toggle-p-button-border-color': '#d1d5db', // gray-300
    '--p-box-shadow': '0 0 #0000, 0 0 #0000, 0 1px 2px 0 rgba(18, 18, 23, 0.05)',
    '--theme-toggle-p-button-hover-bg': '#d1d5db', // gray-300
    '--theme-toggle-p-button-hover-border-color': '#d1d5db', // gray-300
    '--theme-toggle-p-button-active-border-color': '#9ca3af', // gray-400
    '--theme-toggle-p-button-active-bg': '#9ca3af', // gray-400
    '--theme-toggle-p-button-active-icon-color': '#6b7280', // gray-500
    '--theme-toggle-p-button-icon-color': '#6b7280', // gray-500
    '--lang-select-p-select-bg': '#e5e7eb', // gray-200
    '--lang-select-p-select-border-color': '#d1d5db', // gray-300
    '--lang-select-p-select-hover-bg': '#d1d5db', // gray-300
    '--lang-select-p-select-active-border-color': '#9ca3af', // gray-400
    '--lang-select-p-select-active-bg': '#9ca3af', // gray-400
    '--lang-select-p-select-label-color': '#6b7280', // gray-500
    '--lang-select-p-select-overlay-bg': '#fff', // white
    '--lang-select-p-select-overlay-border-color': '#e5e7eb', // gray-200
    '--lang-select-p-select-dropdown-color': '#6b7280', // gray-500
    '--lang-select-p-select-active-dropdown-color': '#6b7280', // gray-500
    '--lang-select-p-select-option-color': '#6b7280', // gray-500
    '--lang-select-p-select-option-focus-bg': '#eff6ff', // blue-50
    '--lang-select-p-select-option-selected-color': '#fff', // white
    '--lang-select-p-select-option-selected-bg': '#2563eb', // blue-600
    '--p-inputtext-textarea-bg': '#fff', // white
    '--p-inputtext-textarea-color': '#000', // black
    '--p-inputtext-textarea-border-color': '#e5e7eb', // gray-200
    '--p-inputtext-textarea-active-border-color': '#d1d5db', // gray-300
    '--p-datepicker-color': '#000', // black
    '--p-datepicker-bg': '#fff', // white
    '--p-datepicker-border-color': '#e5e7eb', // gray-200
    '--p-datepicker-p-button-text-color': '#193cb8', // blue-800
    '--p-datepicker-item-hover-bg': '#eff6ff', // blue-50
    '--p-datepicker-item-active-bg': '#dbeafe', // blue-100
    '--p-datepicker-day-today-bg': '#bfdbfe', // blue-200
    '--p-datepicker-item-selected-bg': '#2563eb', // blue-600
    '--p-datepicker-item-selected-color': '#fff', // white
  },
  dark: {
    '--a-text-decoration-color': '#60a5fa', // blue-400
    '--theme-toggle-p-button-bg': '#111827', // gray-900
    '--theme-toggle-p-button-border-color': '#374151', // gray-700
    '--p-box-shadow': 'none',
    '--theme-toggle-p-button-hover-bg': '#1f2937', // gray-800
    '--theme-toggle-p-button-hover-border-color': '#374151', // gray-700
    '--theme-toggle-p-button-active-border-color': '#374151', // gray-700
    '--theme-toggle-p-button-active-bg': '#374151', // gray-700
    '--theme-toggle-p-button-active-icon-color': '#9ca3af', // gray-400
    '--theme-toggle-p-button-icon-color': '#9ca3af', // gray-400
    '--lang-select-p-select-bg': '#111827', // gray-900
    '--lang-select-p-select-border-color': '#374151', // gray-700
    '--lang-select-p-select-hover-bg': '#1f2937', // gray-800
    '--lang-select-p-select-active-border-color': '#374151', // gray-700
    '--lang-select-p-select-active-bg': '#374151', // gray-700
    '--lang-select-p-select-label-color': '#9ca3af', // gray-400
    '--lang-select-p-select-overlay-bg': '#374151', // gray-700
    '--lang-select-p-select-overlay-border-color': '#4b5563', // gray-600
    '--lang-select-p-select-dropdown-color': '#9ca3af', // gray-400
    '--lang-select-p-select-active-dropdown-color': '#9ca3af', // gray-400
    '--lang-select-p-select-option-color': '#d1d5db', // gray-300
    '--lang-select-p-select-option-focus-bg': '#4b5563', // gray-600
    '--lang-select-p-select-option-selected-color': '#fff', // white
    '--lang-select-p-select-option-selected-bg': '#2563eb', // blue-600
    '--p-inputtext-textarea-bg': '#374151', // gray-700
    '--p-inputtext-textarea-color': '#d1d5db', // gray-300
    '--p-inputtext-textarea-border-color': '#4b5563', // gray-600
    '--p-inputtext-textarea-active-border-color': '#6b7280', // gray-500
    '--p-datepicker-color': '#fff', // white
    '--p-datepicker-bg': '#374151', // gray-700
    '--p-datepicker-border-color': '#4b5563', // gray-600
    '--p-datepicker-p-button-text-color': '#60a5fa', // blue-400
    '--p-datepicker-item-hover-bg': '#4b5563', // gray-600
    '--p-datepicker-item-active-bg': '#6b7280', // gray-500
    '--p-datepicker-day-today-bg': '#3b82f6', // blue-500
    '--p-datepicker-item-selected-bg': '#2563eb', // blue-600
    '--p-datepicker-item-selected-color': '#fff', // white
  },
}
