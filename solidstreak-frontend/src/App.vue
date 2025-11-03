<script setup lang="ts">

import { ref, onMounted } from 'vue'
import HabitCard from '@/components/HabitCard.vue'
import { useHabitStore } from '@/stores/habit';

const habitStore = useHabitStore();
const isHabitsLoading = ref<boolean>(true);
const habitsSuccessfullyLoaded = ref<boolean>(false);

onMounted(async () => {
  await habitStore.fetchHabits(3);
  isHabitsLoading.value = false;
  habitsSuccessfullyLoaded.value = habitStore.requestResult?.success || false;
});

const expandedHabitCardId = ref<number | null>(null);

</script>

<template>
  <div class="p-2">

  <p v-if="isHabitsLoading">Loading...</p>
  <p v-else-if="!habitsSuccessfullyLoaded">Failed to load habits</p>
  <template v-else>
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

  </div>
</template>

<style scoped></style>
