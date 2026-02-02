<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import Toast from 'primevue/toast'

import { dateToLocalString } from './utils/date'
import { ApiFetcher } from '@/api/request'
import { useUserStore } from '@/stores/user'
import { useHabitStore } from '@/stores/habit'
import { type Color, ORANGE, generateColorGradient } from '@/models/color'
import ConfirmDialog from '@/components/confirm-dialog/ConfirmDialog.vue'
import TopPanel from '@/components/top-panel/TopPanel.vue'
import CalendarHeatmap from '@/components/calendar-heatmap/CalendarHeatmap.vue'
import DatePicker from '@/components/date-picker/DatePicker.vue'
import HabitCard from '@/components/habit-card/HabitCard.vue'
import HabitDialog from '@/components/habit-dialog/HabitDialog.vue'

// ─────────────────────────────────────────────
// States & stores
// ─────────────────────────────────────────────
const { t, locale } = useI18n()
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

// ─────────────────────────────────────────────
// Methods
// ─────────────────────────────────────────────
function updateLocale(newLocale: string): void {
  userStore.setLang(newLocale)
  locale.value = newLocale
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

  const initData = window.Telegram?.WebApp?.initData
  const user = window.Telegram?.WebApp?.initDataUnsafe?.user
  const chat = window.Telegram?.WebApp?.initDataUnsafe?.chat

  if (!initData || !user?.id) {
    finishInitialization('Initialization failed')
    return
  }

  const apiFetcher = new ApiFetcher(initData, user.username)

  userStore.init(apiFetcher)
  const userInfoResult = await userStore.upsertUserInfo(user, chat || { id: user.id }) // Use personal chat with user if no other chat info
  if (!userInfoResult.success) {
    finishInitialization('Initialization failed')
    return
  }

  locale.value = userStore.lang

  habitStore.init(apiFetcher)
  const habitsResult = await habitStore.fetchHabits(userStore.id)
  if (!habitsResult.success) {
    finishInitialization('Initialization failed')
    return
  }

  userStore.setAvatarUrl(user.photo_url || '')

  finishInitialization()
})
</script>

<template>
  <p v-if="init">Loading...</p>
  <p v-else-if="initErrorMsg">{{ initErrorMsg }}</p>
  <template v-else>
    <div class="mb-2 border-b border-gray-300 bg-gray-200">
      <TopPanel :lang="userStore.lang" @langSelected="updateLocale" />
    </div>

    <div id="content" style="flex: 1 0 auto" class="mx-auto w-full max-w-lg px-2">
      <CalendarHeatmap
        v-if="!init && !initErrorMsg"
        :values="habitStore.activities"
        :endDate="dateToLocalString(new Date())"
        :max="habitStore.activeHabitsCount"
        tooltipUnit="checks"
        :rangeColor="[
          '#ffffff',
          ...generateColorGradient(
            habitStore.activeHabitsCount == 2
              ? mainHeatmapColor.value400hex
              : mainHeatmapColor.value200hex,
            mainHeatmapColor.value600hex,
            habitStore.activeHabitsCount,
          ),
        ]"
        :round="3"
        class="mb-2 px-2"
      /><!-- TODO: добавить "checks" в переводы -->

      <div class="mb-2 flex items-center justify-between">
        <div class="flex h-10 items-center px-4">
          <span v-if="view === 'active'" class="text-lg font-semibold text-gray-500">{{
            t('app.active', 'Active')
          }}</span>
          <a
            v-else
            @click="view = 'active'"
            :title="t('app.showActiveHabits', 'Show active habits')"
            >{{ t('app.active', 'Active') }}</a
          >
          <span class="text-gray-500">&nbsp;/&nbsp;</span>
          <span v-if="view === 'archived'" class="text-lg font-semibold text-gray-500">{{
            t('app.archived', 'Archived')
          }}</span>
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
            class="rounded-md border border-gray-300 bg-gray-100 px-4 py-2 font-medium text-blue-800 hover:border-blue-100 hover:bg-blue-100 active:border-blue-200 active:bg-blue-200"
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
        @editHabit="openHabitDialog"
        @expandHabitCard="expandedHabitCardId = $event"
        @collapseHabitCard="expandedHabitCardId = null"
        class="mb-2"
      />

      <p
        v-if="view === 'active' && habitStore.activeHabits.length === 0"
        class="text-center text-gray-500"
      >
        {{ t('app.noActiveHabits', 'No active habits') }}.
        <a @click="openHabitDialog()" :title="t('app.createHabit', 'Create a new habit')">{{
          t('app.createOne', 'Create one')
        }}</a
        >!
      </p>
      <p
        v-else-if="view === 'archived' && habitStore.archivedHabits.length === 0"
        class="text-center text-gray-500"
      >
        {{ t('app.noArchivedHabits', 'No archived habits') }}
      </p>

      <DatePicker
        v-if="view === 'active'"
        :date="selectedDate"
        @dateSelected="selectedDate = $event"
      />
    </div>

    <div id="footer" class="mb-2 w-full text-center text-xs text-gray-500 opacity-50">
      <span
        >{{ t('app.madeBy', 'Made by') }} <a href="https://t.me/avasin_dev">@avasin_dev</a></span
      >
    </div>
  </template>

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
