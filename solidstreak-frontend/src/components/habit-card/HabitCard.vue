<script setup lang="ts">
import { ref, computed } from 'vue';
import { useToast } from 'primevue/usetoast';
import type { Habit, HabitCheck } from '@/models/habit'
import { COLORS, GREEN } from '@/models/color'
import { useHabitStore } from '@/stores/habit';
import { SquarePen, Package, PackageOpen, Trash2 } from 'lucide-vue-next';
import CalendarHeatmap from '@/components/calendar-heatmap/CalendarHeatmap.vue'

const props = defineProps<{
  habit: Habit
  currentDate: string
  expanded?: boolean
}>();

const toast = useToast();

const color = COLORS[props.habit.color as keyof typeof COLORS] || GREEN;

const userId = 3; // TODO: получать из внешнего контекста

const habitStore = useHabitStore();

const checksArray = computed(() => {
  return props.habit.checks?.filter(check => check.completed)
    .map(check => ({
      date: check.checkDate,
      count: 1
    })) || [];
});

const currentDateCheck = ref<boolean>(props.habit.checks?.some(check => check.checkDate === props.currentDate && check.completed) || false);

const isCheckButtonHovered = ref<boolean>(false);

async function processCurrentDateCheck(): Promise<void> {

  const check = !currentDateCheck.value;

  const habitCheck: HabitCheck = {
    checkDate: props.currentDate,
    completed: check,
    checkedAt: new Date()
  };

  const result = await habitStore.setHabitCheck(
    userId,
    props.habit.id, 
    habitCheck
  );

  if (result.success) {
    currentDateCheck.value = check;
  } else {
    toast.add({severity:'error', summary: 'Error', detail: 'Failed to ' + (check ? 'check' : 'uncheck') + ' habit', life: 3000});
  }

}

async function processHabitArchiving(): Promise<void> {
  const archive = !props.habit.archived;
  const result = await habitStore.setHabitArchived(userId, props.habit.id, archive);
  if (!result.success) {
    toast.add({severity:'error', summary: 'Error', detail: 'Failed to ' + (archive ? 'archive' : 'unarchive') + ' habit', life: 3000});
  }
}

async function processHabitDeletion(): Promise<void> {
  const result = await habitStore.deleteHabit(userId, props.habit.id);
  if (!result.success) {
    toast.add({severity:'error', summary: 'Error', detail: 'Failed to delete habit', life: 3000});
  }
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
            <Trash2 
              @click.stop="processHabitDeletion()"
              class="w-5 h-5 mr-2 text-gray-300 hover:text-gray-400 cursor-pointer"
            />
          </span>
          <span title="Archive" v-if="!habit.archived">
            <Package 
              @click.stop="processHabitArchiving()"
              class="w-5 h-5 mr-2 text-gray-300 hover:text-gray-400 cursor-pointer"
            />
          </span>
          <span title="Unarchive" v-else>
            <PackageOpen 
              @click.stop="processHabitArchiving()"
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
            @mouseover="isCheckButtonHovered = true"
            @mouseleave="isCheckButtonHovered = false"
            :style="currentDateCheck
              ? {
                  borderColor: isCheckButtonHovered ? color.value500hex : color.value600hex,
                  backgroundColor: isCheckButtonHovered ? color.value400hex : color.value500hex,
                  color: '#fff'
                }
              : {}"
            :class="[
              'w-7 h-7 flex items-center justify-center rounded-lg border cursor-pointer',
              currentDateCheck
                ? ''
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
      :max="1"
      :range-color="[color.value100hex, color.value600hex]"
      :tooltip="false" 
      :round="3"
    />
    
  </div>

</template> 