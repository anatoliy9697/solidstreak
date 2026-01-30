import { createI18n } from 'vue-i18n'

import enApp from './locales/en/app.json'
import ruApp from './locales/ru/app.json'

const messages = {
  en: {
    app: enApp,
  },
  ru: {
    app: ruApp,
  },
}

export const i18n = createI18n({
  legacy: false,
  locale: 'en',
  fallbackLocale: 'en',
  messages,
})

export const LANGS = [
  { code: 'en', name: 'Eng' },
  { code: 'ru', name: 'Рус' },
]
