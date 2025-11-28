<script setup lang="ts">
import { ref, watch } from 'vue'
import DatePicker from 'primevue/datepicker'
import { ChevronLeft, ChevronRight } from 'lucide-vue-next'

import { isSameDay, isBeforeDay, isAfterDay } from '@/utils/date'

// ─────────────────────────────────────────────
// Props
// ─────────────────────────────────────────────
const props = defineProps<{
  date: Date
}>()

// ─────────────────────────────────────────────
// Emits
// ─────────────────────────────────────────────
const emit = defineEmits<{
  (e: 'dateSelected', date: Date): void
}>()

// ─────────────────────────────────────────────
// Constants & reactive state
// ─────────────────────────────────────────────
const selectedDate = ref<Date>(props.date)
const minDate = ((): Date => {
  const curD = new Date()
  curD.setDate(curD.getDate() - 364 - curD.getDay())
  return curD
})()
const maxDate = new Date()

watch(
  () => props.date,
  (newDate) => {
    selectedDate.value = newDate
  },
  { immediate: true },
)

// ─────────────────────────────────────────────
// Methods
// ─────────────────────────────────────────────
function processDateShift(diff: number): void {
  let newDate = new Date(selectedDate.value)
  newDate.setDate(newDate.getDate() + diff)

  if (isBeforeDay(newDate, minDate)) newDate = minDate
  if (isAfterDay(newDate, maxDate)) newDate = maxDate

  if (!isSameDay(selectedDate.value, newDate)) {
    selectedDate.value = newDate
    emit('dateSelected', newDate)
  }
}
</script>

<template>
  <div
    class="fixed bottom-8 left-1/2 z-50 flex h-14 -translate-x-1/2 items-center rounded-full border border-gray-200 bg-white p-2 shadow-sm"
  >
    <ChevronLeft
      :class="[
        'h-10 w-10 rounded-full py-2',
        isBeforeDay(minDate, selectedDate)
          ? 'cursor-pointer text-blue-800 hover:bg-blue-50 active:bg-blue-100'
          : 'cursor-not-allowed text-gray-400 opacity-50',
      ]"
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
          class="p-button-text p-datepicker-today-button rounded-md px-2 py-1"
        >
          <span class="p-button-label justify-center">Today</span>
        </button>
      </template>
    </DatePicker>

    <ChevronRight
      :class="[
        'h-10 w-10 rounded-full py-2 text-blue-800',
        isBeforeDay(selectedDate, maxDate)
          ? 'cursor-pointer text-blue-800 hover:bg-blue-50 active:bg-blue-100'
          : 'cursor-not-allowed text-gray-400 opacity-50',
      ]"
      @click="processDateShift(1)"
    />
  </div>

  <div
    class="pointer-events-none fixed bottom-0 left-0 z-49 h-28 w-dvw"
    style="background: linear-gradient(to top, rgba(0, 0, 0, 0.4), rgba(0, 0, 0, 0))"
  ></div>
</template>
