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
    
    <div class="h-14 p-2 fixed bottom-8 left-1/2 -translate-x-1/2 z-50 bg-white rounded-full shadow-sm border border-gray-200 flex items-center">
      
      <ChevronLeft 
        :class="['w-10 h-10 py-2 rounded-full', isBeforeDay(minDate, selectedDate) ? 'text-blue-800  hover:bg-blue-50 active:bg-blue-100 cursor-pointer' : 'text-gray-400 opacity-50 cursor-not-allowed']"
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
            class="px-2 py-1 p-button-text p-datepicker-today-button rounded-md"
          >
            <span class="p-button-label justify-center">Today</span>
          </button>
        </template>
      </DatePicker>

      <ChevronRight 
        :class="['w-10 h-10 py-2 text-blue-800 rounded-full', isBeforeDay(selectedDate, maxDate) ? 'text-blue-800  hover:bg-blue-50 active:bg-blue-100 cursor-pointer' : 'text-gray-400 opacity-50 cursor-not-allowed']"
        @click="processDateShift(1)"
      />

    </div>

    <div 
      class="fixed bottom-0 left-0 h-28 w-dvw z-49 pointer-events-none"
      style="background: linear-gradient(to top, rgba(0,0,0,0.4), rgba(0,0,0,0));"
    ></div>

</template>