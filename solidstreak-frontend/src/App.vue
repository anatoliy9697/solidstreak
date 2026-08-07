<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { usePrimeVue } from 'primevue/config'
import { useToast } from 'primevue/usetoast'
import Toast from 'primevue/toast'

import { applyTheme } from '@/main'
import { getHeatmapLocale } from '@/i18n'
import { type Theme } from '@/models/theme'
import { dateToISO8601String } from '@/utils/date'
import { ApiFetcher } from '@/api/request'
import { useSubscriptionStore } from '@/stores/subscription'
import { useUserStore } from '@/stores/user'
import { useHabitStore } from '@/stores/habit'
import { type Color, GRAY, ORANGE, generateColorGradient } from '@/models/color'
import ConfirmDialog from '@/components/confirm-dialog/ConfirmDialog.vue'
import TopPanel from '@/components/top-panel/TopPanel.vue'
import CalendarHeatmap from '@/components/calendar-heatmap/CalendarHeatmap.vue'
import DatePicker from '@/components/date-picker/DatePicker.vue'
import HabitCard from '@/components/habit-card/HabitCard.vue'
import HabitDialog from '@/components/habit-dialog/HabitDialog.vue'
import SubscriptionPurchaseDialog from '@/components/subscription-purchase-dialog/SubscriptionPurchaseDialog.vue'

// ─────────────────────────────────────────────
// States & stores
// ─────────────────────────────────────────────
const { t, locale } = useI18n()
const primeVue = usePrimeVue()
const toast = useToast()
const subscriptionStore = useSubscriptionStore()
const userStore = useUserStore()
const habitStore = useHabitStore()
const init = ref<boolean>(true)
const initErrorMsg = ref<string | null>(null)
const view = ref<'active' | 'archived'>('active')
const selectedDate = ref<Date>(new Date())
const mainHeatmapColor = ref<Color>(ORANGE)
const expandedHabitCardId = ref<number | null>(null)
const editingHabitId = ref<number | null>(null)
const isHabitDialogVisible = ref(false)
const isSubscriptionPurchaseDialogVisible = ref(false)

// ─────────────────────────────────────────────
// Methods
// ─────────────────────────────────────────────
function setPrimeVueLocale(t: (key: string, defaultMsg?: string) => string) {
  if (!primeVue?.config?.locale) {
    return
  }
  primeVue.config.locale.dayNamesMin = [
    t('common.daysMin.sun', 'Su'),
    t('common.daysMin.mon', 'Mo'),
    t('common.daysMin.tue', 'Tu'),
    t('common.daysMin.wed', 'We'),
    t('common.daysMin.thu', 'Th'),
    t('common.daysMin.fri', 'Fr'),
    t('common.daysMin.sat', 'Sa'),
  ]
  primeVue.config.locale.monthNames = [
    t('common.months.jan', 'January'),
    t('common.months.feb', 'February'),
    t('common.months.mar', 'March'),
    t('common.months.apr', 'April'),
    t('common.months.may', 'May'),
    t('common.months.jun', 'June'),
    t('common.months.jul', 'July'),
    t('common.months.aug', 'August'),
    t('common.months.sep', 'September'),
    t('common.months.oct', 'October'),
    t('common.months.nov', 'November'),
    t('common.months.dec', 'December'),
  ]
  primeVue.config.locale.monthNamesShort = [
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
  ]
}

async function updateLang(newLang: string): Promise<void> {
  const prevLang = userStore.lang

  if (newLang === prevLang) {
    return
  }

  locale.value = newLang
  setPrimeVueLocale(t)

  const result = await userStore.setLang(newLang)
  if (!result.success) {
    locale.value = prevLang
    setPrimeVueLocale(t)
    toast.add({
      severity: 'error',
      summary: t('common.error', 'Error'),
      detail: t('app.langChangeFailed', 'Failed to change language'),
      life: 3000,
    })
  }
}

function updateTheme(newTheme: Theme): void {
  userStore.setTheme(newTheme)
  applyTheme(newTheme)
}

function getCalendarHeatmapColorRange(
  activeHabitsCount: number,
  color: Color,
  theme: Theme,
): string[] {
  if (theme === 'dark') {
    return [
      GRAY.value700hex,
      ...generateColorGradient(
        generateColorGradient(GRAY.value700hex, color.value700hex, 5)[1]!,
        color.value600hex,
        activeHabitsCount,
      ),
    ]
  } else {
    return [
      '#fff',
      ...generateColorGradient(
        activeHabitsCount === 2 ? color.value400hex : color.value200hex,
        color.value600hex,
        activeHabitsCount,
      ),
    ]
  }
}

const openHabitDialog = (habitId?: number): void => {
  editingHabitId.value = habitId || null
  isHabitDialogVisible.value = true
}

// ─────────────────────────────────────────────
// Lifecycle
// ─────────────────────────────────────────────
onMounted(async (): Promise<void> => {
  function finishInitialization(errorMsg: string | null = null): void {
    initErrorMsg.value = errorMsg
    init.value = false
    window.Telegram?.WebApp?.ready()
  }

  locale.value = userStore.lang
  setPrimeVueLocale(t)

  const initData = window.Telegram?.WebApp?.initData
  const user = window.Telegram?.WebApp?.initDataUnsafe?.user
  const chat = window.Telegram?.WebApp?.initDataUnsafe?.chat

  if (!initData || !user?.id) {
    finishInitialization(t('app.initFailed', 'Initialization failed'))
    return
  }

  const apiFetcher = new ApiFetcher(initData, user.username)

  subscriptionStore.init(apiFetcher)
  const subscriptionPlansResult = await subscriptionStore.fetchSubscriptionPlans()
  if (!subscriptionPlansResult.success || !subscriptionStore.planByCode('basic')) {
    finishInitialization(t('app.initFailed', 'Initialization failed'))
    return
  }

  console.log('subscriptionStore.planByCode("basic")', subscriptionStore.planByCode('basic'))
  console.log('subscriptionStore.planByCode("premium")', subscriptionStore.planByCode('premium'))

  userStore.init(apiFetcher, subscriptionStore.default)
  const upsertUserInfoResult = await userStore.upsertUserInfo(user, chat || { id: user.id }) // Use personal chat with user if no other chat info
  if (!upsertUserInfoResult.success) {
    finishInitialization(t('app.initFailed', 'Initialization failed'))
    return
  }

  applyTheme(userStore.theme)

  locale.value = userStore.lang
  setPrimeVueLocale(t)

  habitStore.init(apiFetcher)
  const habitsResult = await habitStore.fetchHabits(userStore.id)
  if (!habitsResult.success) {
    finishInitialization(t('app.initFailed', 'Initialization failed'))
    return
  }

  userStore.setAvatarUrl(user.photo_url || '')

  finishInitialization()
})
</script>

<template>
  <p v-if="init">{{ t('app.loading', 'Loading') }}...</p>
  <p v-else-if="initErrorMsg">{{ initErrorMsg }}</p>
  <template v-else>
    <div class="mb-2 border-b border-gray-300 bg-gray-200 dark:border-gray-700 dark:bg-gray-900">
      <TopPanel
        :lang="userStore.lang"
        :theme="userStore.theme"
        @langSelected="updateLang"
        @themeToggled="updateTheme"
        @subscriptionPurchaseRequested="isSubscriptionPurchaseDialogVisible = true"
      />
    </div>

    <div id="content" style="flex: 1 0 auto" class="mx-auto w-full max-w-lg px-2">
      <CalendarHeatmap
        id="main-calendar-heatmap"
        v-if="!init && !initErrorMsg"
        :values="habitStore.activities"
        :endDate="dateToISO8601String(new Date())"
        :max="habitStore.activeHabitsCount"
        :tooltipUnit="t('calendarHeatmap.checks', 'checks')"
        :rangeColor="
          getCalendarHeatmapColorRange(
            habitStore.activeHabitsCount,
            mainHeatmapColor,
            userStore.theme,
          )
        "
        :locale="getHeatmapLocale(t)"
        :round="3"
        class="mb-2 px-2"
      />

      <div class="mb-2 flex items-center justify-between">
        <div class="flex h-10 items-center px-4">
          <span
            v-if="view === 'active'"
            class="text-lg font-semibold text-gray-500 dark:text-gray-400"
            >{{ t('app.active', 'Active') }}</span
          >
          <a
            v-else
            @click="view = 'active'"
            :title="t('app.showActiveHabits', 'Show active habits')"
            >{{ t('app.active', 'Active') }}</a
          >
          <span class="text-lg font-semibold text-gray-500 dark:text-gray-400">&nbsp;/&nbsp;</span>
          <span
            v-if="view === 'archived'"
            class="text-lg font-semibold text-gray-500 dark:text-gray-400"
            >{{ t('app.archived', 'Archived') }}</span
          >
          <a
            v-else
            @click="view = 'archived'"
            :title="t('app.showArchivedHabits', 'Show archived habits')"
            >{{ t('app.archived', 'Archived') }}</a
          >
        </div>

        <div v-show="view === 'active'">
          <button
            @click="openHabitDialog()"
            class="cursor-pointer rounded-md border border-gray-300 bg-gray-100 px-4 py-2 font-medium text-blue-800 hover:border-blue-100 hover:bg-blue-100 active:border-blue-200 active:bg-blue-200 dark:border-gray-700 dark:bg-gray-800 dark:text-blue-400 dark:hover:border-gray-700 dark:hover:bg-gray-700 dark:active:border-gray-600 dark:active:bg-gray-600"
            :title="t('app.createHabit', 'Create a new habit')"
          >
            + {{ t('app.newHabit', 'New habit') }}
          </button>
        </div>
      </div>

      <HabitCard
        v-show="view === 'active' && habitStore.activeHabits.length > 0"
        v-for="habit in habitStore.activeHabits"
        :key="habit.id"
        :habit="habit"
        :selectedDate="selectedDate"
        :expanded="expandedHabitCardId === habit.id"
        :theme="userStore.theme"
        @editHabit="openHabitDialog"
        @expandHabitCard="expandedHabitCardId = $event"
        @collapseHabitCard="expandedHabitCardId = null"
        class="mb-2"
      />

      <HabitCard
        v-show="view === 'archived' && habitStore.archivedHabits.length > 0"
        v-for="habit in habitStore.archivedHabits"
        :key="habit.id"
        :habit="habit"
        :selectedDate="selectedDate"
        :expanded="expandedHabitCardId === habit.id"
        :theme="userStore.theme"
        @editHabit="openHabitDialog"
        @expandHabitCard="expandedHabitCardId = $event"
        @collapseHabitCard="expandedHabitCardId = null"
        class="mb-2"
      />

      <p
        v-if="view === 'active' && habitStore.activeHabits.length === 0"
        class="text-center text-gray-500 dark:text-gray-400"
      >
        {{ t('app.noActiveHabits', 'No active habits') }}.
        <a @click="openHabitDialog()" :title="t('app.createHabit', 'Create a new habit')">{{
          t('app.createOne', 'Create one')
        }}</a
        >!
      </p>
      <p
        v-else-if="view === 'archived' && habitStore.archivedHabits.length === 0"
        class="text-center text-gray-500 dark:text-gray-400"
      >
        {{ t('app.noArchivedHabits', 'No archived habits') }}
      </p>

      <DatePicker
        v-if="view === 'active'"
        :date="selectedDate"
        @dateSelected="selectedDate = $event"
      />
    </div>

    <div
      id="footer"
      class="mb-2 w-full text-center text-xs text-gray-500 opacity-50 dark:text-gray-400"
    >
      <span
        >{{ t('app.madeBy', 'Made by') }} <a href="https://t.me/avasin_dev">@avasin_dev</a></span
      >
    </div>
  </template>

  <SubscriptionPurchaseDialog
    :visible="isSubscriptionPurchaseDialogVisible"
    @closeSubscriptionPurchaseDialog="isSubscriptionPurchaseDialogVisible = false"
  />
  <HabitDialog
    :visible="isHabitDialogVisible"
    :newHabit="editingHabitId === null"
    :habit="editingHabitId !== null ? habitStore.habitById(editingHabitId) : undefined"
    @closeHabitDialog="isHabitDialogVisible = false"
  />
  <ConfirmDialog :style="{ borderRadius: '0.375rem' }"></ConfirmDialog>
  <Toast position="bottom-right" />
</template>

<style scoped></style>
