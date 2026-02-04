import { createI18n } from 'vue-i18n'

import enCommon from './locales/en/common.json'
import ruCommon from './locales/ru/common.json'
import enApp from './locales/en/app.json'
import ruApp from './locales/ru/app.json'
import enTopPanel from './locales/en/top-panel.json'
import ruTopPanel from './locales/ru/top-panel.json'
import enHabitCard from './locales/en/habit-card.json'
import ruHabitCard from './locales/ru/habit-card.json'
import enHabitDialog from './locales/en/habit-dialog.json'
import ruHabitDialog from './locales/ru/habit-dialog.json'
import enDatePicker from './locales/en/date-picker.json'
import ruDatePicker from './locales/ru/date-picker.json'

const messages = {
  en: {
    common: enCommon,
    app: enApp,
    topPanel: enTopPanel,
    habitCard: enHabitCard,
    habitDialog: enHabitDialog,
    datePicker: enDatePicker,
  },
  ru: {
    common: ruCommon,
    app: ruApp,
    topPanel: ruTopPanel,
    habitCard: ruHabitCard,
    habitDialog: ruHabitDialog,
    datePicker: ruDatePicker,
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
