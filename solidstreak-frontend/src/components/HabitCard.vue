<script setup lang="ts">
import { ref, computed } from 'vue';
import type { Habit, HabitCheck } from '@/models/habit'
import { useHabitStore } from '@/stores/habit';

const props = defineProps<{
  habit: Habit
  currentDate: string
  expanded?: boolean
}>()

const habitStore = useHabitStore();

const checksArray = computed(() => {
  return props.habit.checks?.map(check => ({
    date: check.checkDate,
    count: check.completed ? 0 : null
  })) || [];
});

const currentDateCheck = ref<boolean>(props.habit.checks?.some(check => check.checkDate === props.currentDate && check.completed) || false);

async function processCurrentDateCheck(): Promise<void> {

  const habitCheck: HabitCheck = {
    checkDate: props.currentDate,
    completed: !currentDateCheck.value,
    checkedAt: new Date()
  };

  const result = await habitStore.setHabitCheck(
    3,  // TODO: получать из внешнего контекста
    props.habit.id, 
    habitCheck
  );

  if (result.success) currentDateCheck.value = habitCheck.completed;

}

</script>

<template>
  <div class="bg-white rounded-md shadow-sm border border-gray-300 px-4 py-2 cursor-pointer">
  
  <div class="flex items-start justify-between">

      <div :class="[
        'min-h-7 flex',
        expanded ? 'flex-col' : 'items-center'
      ]">
        <h2>{{ habit.title }}</h2>
        <p v-if="expanded" class="mb-2">{{ habit.description }}</p>
      </div>

      <div>
        <button
          @click.stop="processCurrentDateCheck()"
          :class="[
            'w-7 h-7 flex items-center justify-center rounded-lg border cursor-pointer',
            currentDateCheck
              ? 'border-green-600 text-white bg-green-500 hover:border-green-500 hover:bg-green-400'
              : 'border-gray-400 text-gray-400 hover:text-gray-500 hover:border-gray-500'
          ]"
        >
          <svg xmlns="http://www.w3.org/2000/svg" class="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
            <path stroke-linecap="round" stroke-linejoin="round" d="M5 13l4 4L19 7" />
          </svg>
        </button>
      </div>
      
    </div>
    
    <calendar-heatmap
      v-if="expanded"
      :values="checksArray"
      :end-date="currentDate"
      :range-color="['#e5e7eb', '#16a34a']"
      :tooltip="false"
      :round="3"
    />
    
  </div>
</template> 