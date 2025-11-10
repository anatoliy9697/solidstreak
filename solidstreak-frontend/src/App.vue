<script setup lang="ts">

import { ref, onMounted, computed } from 'vue'
import Toast from 'primevue/toast';
import { type Color, PURPLE, generateColorGradient } from '@/models/color'
import { useHabitStore } from '@/stores/habit';
import HabitCard from '@/components/habit-card/HabitCard.vue'
import CalendarHeatmap from '@/components/calendar-heatmap/CalendarHeatmap.vue'

const habitStore = useHabitStore();
const isHabitsLoading = ref<boolean>(true);
const habitsSuccessfullyLoaded = ref<boolean>(false);
const currentDate = new Date().toISOString().split('T')[0] || '';
const mainHeatmapColor = ref<Color>(PURPLE)

onMounted(async () => {
  const result = await habitStore.fetchHabits(3);
  isHabitsLoading.value = false;
  habitsSuccessfullyLoaded.value = result.success || false;
});

const activeHabitsCount = computed(() => {
  return habitStore.habits.filter(habit => !habit.archived).length;
});
const activitiesMap = computed(() => {
  const map = new Map<string, number>();
  
  habitStore.habits
    .filter(habit => !habit.archived)
    .forEach(habit => {
      habit.checks?.forEach(check => {
        if (check.completed) {
          const date = check.checkDate;
          map.set(date, (map.get(date) || 0) + 1);
        }
      });
    });

  return map;
});
const expandedHabitCardId = ref<number | null>(null);

</script>

<template>

  <p v-if="isHabitsLoading">Loading...</p>
  <p v-else-if="!habitsSuccessfullyLoaded">Failed to load habits</p>
  <template v-else>

    <calendar-heatmap
      v-if="!isHabitsLoading && habitsSuccessfullyLoaded"
      :values="Array.from(activitiesMap.entries()).map(([date, count]) => ({ date, count }))"
      :end-date="currentDate"
      :max="activeHabitsCount"
      :range-color="['#ffffff', ...generateColorGradient(mainHeatmapColor.value200hex, mainHeatmapColor.value800hex, activeHabitsCount)]"
      :round="3"
      class="mb-4 px-2"
    />
    
    <habit-card
      v-for="habit in habitStore.habits"
      :key="habit.id"
      :habit="habit"
      :current-date="new Date().toISOString().split('T')[0] || ''"
      :expanded="expandedHabitCardId === habit.id"
      @click="expandedHabitCardId = habit.id"
      class="mb-2"
    />

  </template>

  <Toast position="bottom-right"/>

</template>

<style scoped></style>
