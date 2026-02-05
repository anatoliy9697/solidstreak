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
import enCalendarHeatmap from './locales/en/calendar-heatmap.json'
import ruCalendarHeatmap from './locales/ru/calendar-heatmap.json'

const messages = {
  en: {
    common: enCommon,
    app: enApp,
    topPanel: enTopPanel,
    habitCard: enHabitCard,
    habitDialog: enHabitDialog,
    datePicker: enDatePicker,
    calendarHeatmap: enCalendarHeatmap,
  },
  ru: {
    common: ruCommon,
    app: ruApp,
    topPanel: ruTopPanel,
    habitCard: ruHabitCard,
    habitDialog: ruHabitDialog,
    datePicker: ruDatePicker,
    calendarHeatmap: ruCalendarHeatmap,
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

export function getHeatmapLocale(t: (key: string, defaultMsg?: string) => string) {
  return {
    months: [
      t('common.monthsShort.jan', 'Jan'),
      t('common.monthsShort.feb', 'Feb'),
      t('common.monthsShort.mar', 'Mar'),
      t('common.monthsShort.apr', 'Apr'),
      t('common.monthsShort.may', 'May'),
      t('common.monthsShort.jun', 'Jun'),
      t('common.monthsShort.jul', 'Jul'),
      t('common.monthsShort.aug', 'Aug'),
      t('common.monthsShort.sep', 'Sep'),
      t('common.monthsShort.oct', 'Oct'),
      t('common.monthsShort.nov', 'Nov'),
      t('common.monthsShort.dec', 'Dec'),
    ],
    days: [
      t('common.daysShort.sun', 'Sun'),
      t('common.daysShort.mon', 'Mon'),
      t('common.daysShort.tue', 'Tue'),
      t('common.daysShort.wed', 'Wed'),
      t('common.daysShort.thu', 'Thu'),
      t('common.daysShort.fri', 'Fri'),
      t('common.daysShort.sat', 'Sat'),
    ],
    on: t('common.on', 'on'),
    less: t('common.upperLess', 'Less'),
    more: t('common.upperMore', 'More'),
  }
}
