<script setup lang="ts">

import { ref, onMounted } from 'vue';
import Toast from 'primevue/toast';

import { useHabitStore } from '@/stores/habit';
import { type Color, PURPLE, generateColorGradient } from '@/models/color'
import ConfirmDialog from '@/components/confirm-dialog/ConfirmDialog.vue'
import CalendarHeatmap from '@/components/calendar-heatmap/CalendarHeatmap.vue'
import HabitCard from '@/components/habit-card/HabitCard.vue'
import HabitDialog from '@/components/habit-dialog/HabitDialog.vue'

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
      :current-date="currentDate"
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
      :current-date="new Date().toISOString().split('T')[0] || ''"
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
