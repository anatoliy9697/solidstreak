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

export function getHeatmapLocale(t: (key: string) => string) {
  return {
    months: [
      t('common.months.jan'),
      t('common.months.feb'),
      t('common.months.mar'),
      t('common.months.apr'),
      t('common.months.may'),
      t('common.months.jun'),
      t('common.months.jul'),
      t('common.months.aug'),
      t('common.months.sep'),
      t('common.months.oct'),
      t('common.months.nov'),
      t('common.months.dec'),
    ],
    days: [
      t('common.days.sun'),
      t('common.days.mon'),
      t('common.days.tue'),
      t('common.days.wed'),
      t('common.days.thu'),
      t('common.days.fri'),
      t('common.days.sat'),
    ],
    on: t('common.on'),
    less: t('common.upperLess'),
    more: t('common.upperMore'),
  }
}
