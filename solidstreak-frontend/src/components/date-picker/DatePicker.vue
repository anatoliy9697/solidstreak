<script setup lang="ts">

import { ref, watch } from 'vue';
import DatePicker from 'primevue/datepicker';
import { ChevronLeft, ChevronRight } from 'lucide-vue-next';

import { isSameDay, isBeforeDay, isAfterDay } from '@/utils/date';

// ─────────────────────────────────────────────
// Props
// ─────────────────────────────────────────────
const props = defineProps<{
  date: Date
}>();

// ─────────────────────────────────────────────
// Emits
// ─────────────────────────────────────────────
const emit = defineEmits<{
  (e: 'dateSelected', date: Date): void
}>();

// ─────────────────────────────────────────────
// Constants & reactive state
// ─────────────────────────────────────────────
const selectedDate = ref<Date>(props.date);
const minDate = ((): Date => {
  const curD = new Date();
  curD.setDate(curD.getDate() - 364 - curD.getDay());
  return curD;
})();
const maxDate = new Date();

watch(() => props.date, (newDate) => {
  selectedDate.value = newDate;
}, { immediate: true });

// ─────────────────────────────────────────────
// Methods
// ─────────────────────────────────────────────
function processDateShift(diff: number): void {
  let newDate = new Date(selectedDate.value)
  newDate.setDate(newDate.getDate() + diff);

  if (isBeforeDay(newDate, minDate)) newDate = minDate;
  if (isAfterDay(newDate, maxDate)) newDate = maxDate;
  
  if (!isSameDay(selectedDate.value, newDate)) {
    selectedDate.value = newDate;
    emit('dateSelected', newDate);
  }
}

</script>

<template>
    
    <div class="h-10 fixed bottom-4 left-1/2 -translate-x-1/2 z-50 bg-white rounded-full shadow-sm border border-gray-300 flex items-center">
      
      <ChevronLeft 
        :class="['w-10 h-10 py-2 text-gray-400', isBeforeDay(minDate, selectedDate) ? 'hover:text-gray-500 cursor-pointer' : 'opacity-50 cursor-not-allowed']"
        @click="processDateShift(-1)"
      />

      <DatePicker 
        v-model="selectedDate"
        @dateSelect="emit('dateSelected', $event)" 
        dateFormat="dd/mm/yy"
        :minDate="minDate"
        :maxDate="maxDate"
        :manualInput="false"
        showButtonBar
        placeholder="Customized"
      >
        <template #buttonbar="{ todayCallback }">
          <button 
            @click="todayCallback"
            type="button"
            class="px-2 py-1 p-button p-component p-button-text p-datepicker-today-button"
          >
            <span class="p-button-label justify-center">Today</span>
          </button>
        </template>
      </DatePicker>

      <ChevronRight 
        :class="['w-10 h-10 py-2 text-gray-400', isBeforeDay(selectedDate, maxDate) ? 'hover:text-gray-500 cursor-pointer' : 'opacity-50 cursor-not-allowed']"
        @click="processDateShift(1)"
      />

    </div>

    <div 
      class="fixed bottom-0 left-0 h-20 w-dvw z-49 pointer-events-none"
      style="background: linear-gradient(to top, rgba(0,0,0,0.4), rgba(0,0,0,0));"
    ></div>

</template>