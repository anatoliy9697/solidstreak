<script setup lang="ts">

import { ref, onMounted } from 'vue';
import Toast from 'primevue/toast';

import { dateToLocalString } from './utils/date';
import { ApiFetcher } from '@/api/request';
import { useUserStore } from '@/stores/user';
import { useHabitStore } from '@/stores/habit';
import { type Color, PURPLE, generateColorGradient } from '@/models/color'
import ConfirmDialog from '@/components/confirm-dialog/ConfirmDialog.vue'
import CalendarHeatmap from '@/components/calendar-heatmap/CalendarHeatmap.vue'
import DatePicker from '@/components/date-picker/DatePicker.vue'
import HabitCard from '@/components/habit-card/HabitCard.vue'
import HabitDialog from '@/components/habit-dialog/HabitDialog.vue'

// ─────────────────────────────────────────────
// States & stores
// ─────────────────────────────────────────────
const userStore = useUserStore();
const habitStore = useHabitStore();

const selectedDate = ref<Date>(new Date());
const mainHeatmapColor = ref<Color>(PURPLE)

const init = ref<boolean>(true);
const initErrorMsg = ref<string | null>(null);
const view = ref<'active' | 'archived'>('active');
const expandedHabitCardId = ref<number | null>(null);
const editingHabitId = ref<number | null>(null);
const isHabitDialogVisible = ref(false);

// ─────────────────────────────────────────────
// Methods
// ─────────────────────────────────────────────
const openHabitDialog = (habitId?: number): void => {
  editingHabitId.value = habitId || null;
  isHabitDialogVisible.value = true;
}

// ─────────────────────────────────────────────
// Lifecycle
// ─────────────────────────────────────────────
onMounted(async (): Promise<void> => {
  const initData = window.Telegram?.WebApp?.initData;
  const user = window.Telegram?.WebApp?.initDataUnsafe?.user;
  const chat = window.Telegram?.WebApp?.initDataUnsafe?.chat;

  if (!initData || !user?.id) {
    initErrorMsg.value = 'Initialization failed';
    init.value = false;
    return;
  }

  const apiFetcher = new ApiFetcher(initData, user.username);
  
  userStore.init(apiFetcher);
  const userInfoResult = await userStore.upsertUserInfo(user, chat || { id: user.id }); // Use personal chat with user if no other chat info
  if (!userInfoResult.success) {
    initErrorMsg.value = 'Initialization failed';
    init.value = false;
    return;
  }

  habitStore.init(apiFetcher);
  const habitsResult = await habitStore.fetchHabits(userStore.id);
  if (!habitsResult.success) {
    initErrorMsg.value = 'Initialization failed';
    init.value = false;
    return;
  }
  
  userStore.setAvatarUrl(user.photo_url || '');

  init.value = false;
});

</script>

<template>

  <p v-if="init">Loading...</p>
  <p v-else-if="initErrorMsg">{{ initErrorMsg }}</p>
  <template v-else>

    <CalendarHeatmap
      v-if="!init && !initErrorMsg"
      :values="habitStore.activities"
      :end-date="dateToLocalString(new Date())"
      :max="habitStore.activeHabitsCount"
      tooltip-unit="checks"
      :range-color="['#ffffff', ...generateColorGradient(habitStore.activeHabitsCount == 2 ? mainHeatmapColor.value400hex : mainHeatmapColor.value200hex, mainHeatmapColor.value800hex, habitStore.activeHabitsCount)]"
      :round="3"
      class="mb-2 px-2"
    />

    <div class="flex justify-between items-center mb-2">

      <div class="h-10 flex items-center px-4">
        <span v-if="view === 'active'" class="text-lg font-semibold text-gray-500">Active</span>
        <a v-else @click="view = 'active'">Active</a>
        <span class="text-gray-500">&nbsp;/&nbsp;</span>
        <span v-if="view === 'archived'" class="text-lg font-semibold text-gray-500">Archived</span>
        <a v-else @click="view = 'archived'">Archived</a>
      </div>

      <div v-show="view === 'active'">
        <button
          @click="openHabitDialog()"
          class="px-4 py-2 rounded-md border border-gray-300 bg-gray-100 text-gray-800 font-medium hover:bg-gray-200"
        >+ New habit</button>
      </div>

    </div>

    <HabitCard
      v-show="view === 'active' && habitStore.activeHabits.length > 0"
      v-for="habit in habitStore.activeHabits"
      :key="habit.id"
      :habit="habit"
      :selectedDate="selectedDate"
      :expanded="expandedHabitCardId === habit.id"
      @click="expandedHabitCardId = habit.id"
      @edit-habit="openHabitDialog"
      class="mb-2"
    />

    <HabitCard
      v-show="view === 'archived' && habitStore.archivedHabits.length > 0"
      v-for="habit in habitStore.archivedHabits"
      :key="habit.id"
      :habit="habit"
      :selectedDate="selectedDate"
      :expanded="expandedHabitCardId === habit.id"
      @click="expandedHabitCardId = habit.id"
      @edit-habit="openHabitDialog"
      class="mb-2"
    />

    <p 
      v-if="view === 'active' && habitStore.activeHabits.length === 0" 
      class="text-gray-500 text-center"
    >No active habits. <a @click="openHabitDialog()">Create one</a>!</p>
    <p 
      v-else-if="view === 'archived' && habitStore.archivedHabits.length === 0" 
      class="text-gray-500 text-center"
    >No archived habits</p>

    <DatePicker 
      v-if="view === 'active'" 
      :date="selectedDate" 
      @dateSelected="selectedDate = $event"
    />

  </template>

  <HabitDialog
    :visible="isHabitDialogVisible"
    :new-habit="editingHabitId === null"
    :habit="editingHabitId !== null ? habitStore.habitById(editingHabitId) : undefined"
    @closeHabitDialog="isHabitDialogVisible = false"
  />
  <ConfirmDialog :style="{ borderRadius: '0.375rem' }"></ConfirmDialog>
  <Toast position="bottom-right"/>

</template>

<style scoped></style>
