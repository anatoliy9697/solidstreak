<script setup lang="ts">
import { ref, computed } from 'vue';
import type { Habit, HabitCheck } from '@/models/habit'
import { useHabitStore } from '@/stores/habit';
import { SquarePen, Package, PackageOpen, Trash2 } from 'lucide-vue-next';

const props = defineProps<{
  habit: Habit
  currentDate: string
  expanded?: boolean
}>()

const userId = 3; // TODO: получать из внешнего контекста

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
    userId,
    props.habit.id, 
    habitCheck
  );

  if (result.success) currentDateCheck.value = habitCheck.completed;

}

</script>

<template>
  <div :class="['bg-white rounded-md shadow-sm border border-gray-300 px-4 py-2 cursor-pointer', habit.archived ? 'opacity-50' : '']">

  <div :class="['flex items-start justify-between', expanded && !habit.archived ? 'mb-2' : '']">

      <div :class="['min-h-7 flex mr-4', expanded ? 'flex-col' : 'items-center']">
        <h2 :class="['leading-none', expanded && habit.description ? ' mb-2' : '']">{{ habit.title }}</h2>
        <p v-if="expanded && habit.description">{{ habit.description }}</p>
      </div>

      <div class="flex items-center">

        <div class="flex items-center" v-if="expanded">
          <span title="Delete">
            <!-- TODO: сделать обработку ошибки запроса -->
            <Trash2 
              @click.stop="habitStore.deleteHabit(userId, habit.id)"
              class="w-5 h-5 mr-2 text-gray-300 hover:text-gray-400 cursor-pointer"
            />
          </span>
          <span title="Archive" v-if="!habit.archived">
            <!-- TODO: сделать обработку ошибки запроса -->
            <Package 
              @click.stop="habitStore.setHabitArchived(userId, habit.id, true)"
              class="w-5 h-5 mr-2 text-gray-300 hover:text-gray-400 cursor-pointer"
            />
          </span>
          <span title="Unarchive" v-else>
            <!-- TODO: сделать обработку ошибки запроса -->
            <PackageOpen 
              @click.stop="habitStore.setHabitArchived(userId, habit.id, false)"
              class="w-5 h-5 mr-2 text-gray-300 hover:text-gray-400 cursor-pointer"
            />
          </span>
          <span title="Edit">
            <SquarePen 
              class="w-5 h-5 text-gray-300 hover:text-gray-400 cursor-pointer"
            />
          </span>
        </div>

        <div v-if="!habit.archived" class="ml-4">
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
      
    </div>
    
    <calendar-heatmap
      v-if="expanded && !habit.archived"
      :values="checksArray"
      :end-date="currentDate"
      :range-color="['#e5e7eb', '#16a34a']"
      :tooltip="false"
      :round="3"
    />
    
  </div>
</template> 