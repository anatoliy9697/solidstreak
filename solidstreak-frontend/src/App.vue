<script setup lang="ts">

import { ref, onMounted } from 'vue';
import Toast from 'primevue/toast';

import { useHabitStore } from '@/stores/habit';
import { type Color, PURPLE, generateColorGradient } from '@/models/color'
import ConfirmDialog from '@/components/confirm-dialog/ConfirmDialog.vue'
import CalendarHeatmap from '@/components/calendar-heatmap/CalendarHeatmap.vue'
import HabitCard from '@/components/habit-card/HabitCard.vue'

// ─────────────────────────────────────────────
// States & stores
// ─────────────────────────────────────────────
const habitStore = useHabitStore();

const currentDate = new Date().toISOString().split('T')[0] || '';
const mainHeatmapColor = ref<Color>(PURPLE)

const isHabitsLoading = ref<boolean>(true);
const habitsSuccessfullyLoaded = ref<boolean>(false);
const view = ref<'active' | 'archived'>('active');
const expandedHabitCardId = ref<number | null>(null);

// ─────────────────────────────────────────────
// Lifecycle
// ─────────────────────────────────────────────
onMounted(async () => {
  const result = await habitStore.fetchHabits(3);
  isHabitsLoading.value = false;
  habitsSuccessfullyLoaded.value = result.success || false;
});

</script>

<template>

  <p v-if="isHabitsLoading">Loading...</p>
  <p v-else-if="!habitsSuccessfullyLoaded">Failed to load habits</p>
  <template v-else>

    <CalendarHeatmap
      v-if="!isHabitsLoading && habitsSuccessfullyLoaded"
      :values="habitStore.activities"
      :end-date="new Date().toISOString().split('T')[0] || ''"
      :max="habitStore.activeHabitsCount"
      :range-color="['#ffffff', ...generateColorGradient(habitStore.activeHabitsCount == 2 ? mainHeatmapColor.value400hex : mainHeatmapColor.value200hex, mainHeatmapColor.value800hex, habitStore.activeHabitsCount)]"
      :round="3"
      class="mb-2 px-2"
    />

    <div class="mb-2 px-4">
      <span v-if="view === 'active'" class="text-lg font-semibold text-gray-500">Active</span>
      <a v-else @click="view = 'active'">Active</a>
      <span class="text-gray-500"> / </span>
      <span v-if="view === 'archived'" class="text-lg font-semibold text-gray-500">Archived</span>
      <a v-else @click="view = 'archived'">Archived</a>
    </div>

    <HabitCard
      v-show="view === 'active'"
      v-for="habit in habitStore.activeHabits"
      :key="habit.id"
      :habit="habit"
      :current-date="currentDate"
      :expanded="expandedHabitCardId === habit.id"
      @click="expandedHabitCardId = habit.id"
      class="mb-2"
    />

    <HabitCard
      v-show="view === 'archived'"
      v-for="habit in habitStore.archivedHabits"
      :key="habit.id"
      :habit="habit"
      :current-date="new Date().toISOString().split('T')[0] || ''"
      :expanded="expandedHabitCardId === habit.id"
      @click="expandedHabitCardId = habit.id"
      class="mb-2"
    />

  </template>

  <Toast position="bottom-right"/>
  <ConfirmDialog></ConfirmDialog>

</template>

<style scoped></style>
