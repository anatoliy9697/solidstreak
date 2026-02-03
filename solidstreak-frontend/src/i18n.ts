import { createI18n } from 'vue-i18n'

import enCommon from './locales/en/common.json'
import ruCommon from './locales/ru/common.json'
import enApp from './locales/en/app.json'
import ruApp from './locales/ru/app.json'
import enHabitCard from './locales/en/habit-card.json'
import ruHabitCard from './locales/ru/habit-card.json'
import enHabitDialog from './locales/en/habit-dialog.json'
import ruHabitDialog from './locales/ru/habit-dialog.json'

const messages = {
  en: {
    common: enCommon,
    app: enApp,
    habitCard: enHabitCard,
    habitDialog: enHabitDialog,
  },
  ru: {
    common: ruCommon,
    app: ruApp,
    habitCard: ruHabitCard,
    habitDialog: ruHabitDialog,
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
