<script setup lang="ts">
import type { Habit } from '@/models/habit'
import { ref, computed } from 'vue';

const props = defineProps<{
  habit: Habit
  currentDate: string
  expanded?: boolean
}>()

const checksMap = ref<Map<string, boolean>>(new Map());
if (props.habit.checks) {
  props.habit.checks.forEach(check => {
    if (check.checkDate && check.completed) {
      checksMap.value.set(check.checkDate, true);
    }
  });
}

const checksArray = computed(() => {
  return Array.from(checksMap.value.keys()).map(date => ({ date, count: 0 }));
});

const currentDateCheck = ref<boolean>(false);
function processCurrentDateCheck() {
  currentDateCheck.value = !currentDateCheck.value;
  if (currentDateCheck.value) {
    checksMap.value.set(props.currentDate, true);
  } else {
    checksMap.value.delete(props.currentDate);
  }
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