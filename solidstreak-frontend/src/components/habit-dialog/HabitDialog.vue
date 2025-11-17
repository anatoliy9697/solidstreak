<script setup lang="ts">

import { computed, ref, watch } from 'vue';
import { useToast } from 'primevue/usetoast';
import 'primeicons/primeicons.css'
import Dialog from 'primevue/dialog';
import InputText from 'primevue/inputtext';
import Textarea from 'primevue/textarea';

import { useHabitStore } from '@/stores/habit';
import type { Habit } from '@/models/habit'
import type { RequestResult } from '@/api/habit'
import { type Color, COLORS, GREEN } from '@/models/color'
import ColorPicker from '@/components/color-picker/ColorPicker.vue';

// ─────────────────────────────────────────────
// Props
// ─────────────────────────────────────────────
const props = defineProps<{
  visible: boolean
  newHabit: boolean
  habit?: Habit
}>();

// ─────────────────────────────────────────────
// Emits
// ─────────────────────────────────────────────
const emit = defineEmits<{
  (e: 'closeHabitDialog'): void
}>();

// ─────────────────────────────────────────────
// Composables & stores
// ─────────────────────────────────────────────
const habitStore = useHabitStore();
const toast = useToast();

// ─────────────────────────────────────────────
// Constants & reactive state
// ─────────────────────────────────────────────
const habitTitle = ref('')
const titleValidationMessage = ref('')
const habitDescription = ref('')
const color = ref(GREEN);

watch(() => props.visible, (visible) => {
  if (visible) {
    habitTitle.value = props.habit?.title || ''
    titleValidationMessage.value = ''
    habitDescription.value = props.habit?.description || ''
    color.value = COLORS[props.habit?.color as keyof typeof COLORS] || GREEN;
  }
}, { immediate: true });

watch(habitTitle, () => {
  titleValidationMessage.value = ''
});

// ─────────────────────────────────────────────
// Computed
// ─────────────────────────────────────────────
const dialogVisible = computed({
  get: () => props.visible,
  set: () => { emit('closeHabitDialog'); }
});

// ─────────────────────────────────────────────
// Methods
// ─────────────────────────────────────────────
async function processHabitSaving(): Promise<void> {
  if (!habitTitle.value) {
    titleValidationMessage.value = 'Title is required';
    return;
  }

  const newHabit: Habit = {} as Habit;

  if (!props.newHabit) newHabit.id = props.habit!.id;
  newHabit.title = habitTitle.value;
  newHabit.description = habitDescription.value;
  newHabit.color = color.value.name;
  newHabit.archived = props.habit?.archived || false;
  newHabit.isPublic = props.habit?.isPublic || false;
  
  let result: RequestResult;
  if (props.newHabit) {
    result = await habitStore.createHabit(3, newHabit); // TODO: передавать не хардкоженный userId
  } else {
    result = await habitStore.updateHabit(3, newHabit); // TODO: передавать не хардкоженный userId
  }
  
  if (!result.success) {
    toast.add({severity:'error', summary: 'Error', detail: 'Failed to save habit', life: 3000});
  } else {
    emit('closeHabitDialog');
  }
}

async function onColorSelected(selectedColor: Color): Promise<void> {
  color.value = selectedColor;
}

</script>

<template>

  <Dialog 
    v-model:visible="dialogVisible" 
    position="bottom" 
    :modal="true" 
    :draggable="false" 
    class="p-4"
    :style="{ width: '95vw', maxWidth: '500px', margin: '0.5rem', marginBottom: 0, borderRadius: '0.375rem 0.375rem 0 0'}"
  >
    <template #container="{ closeCallback }">

      <div class="mb-4 flex justify-between items-start">
        <div>
          <h1>{{ props.newHabit ? 'New habit' : 'Edit habit' }}</h1>
        </div>
        <div>
          <i 
            @click="closeCallback"  
            class="pi pi-times text-gray-400 hover:text-gray-500 cursor-pointer" 
          ></i>
        </div>
      </div>

      <label for="habit-title" class="font-semibold">Title:</label>
      <InputText 
        id="habit-title" 
        v-model="habitTitle"
        :placeholder="titleValidationMessage"
        maxlength="64"
        :class="['w-full', titleValidationMessage ? 'p-invalid' : '']" 
      />
      <p class="text-right text-xs text-gray-400">{{ habitTitle.length }}/64</p>

      <label for="habit-description" class="font-semibold">Description:</label>
      <Textarea 
        id="habit-description" 
        v-model="habitDescription"
        autoResize 
        maxlength="256"
        rows="5" cols="30"
        class="w-full"
      />
      <p class="text-right text-xs text-gray-400">{{ habitDescription.length }}/256</p>

      <div class="flex flex-row items-start mb-4">
        <label class="block font-semibold mr-2">Color:</label>
        <ColorPicker 
          :selectedColor="color" 
          @colorSelected="onColorSelected"
        />
      </div>

      <div class="flex gap-2">
        <button
            @click="closeCallback"
            class="w-1/2 px-4 py-2 rounded-md border border-gray-300 bg-gray-100 text-gray-800 font-medium hover:bg-gray-200"
        >Cancel</button>
        <button
            @click="processHabitSaving"
            class="w-1/2 px-4 py-2 rounded-md border border-green-700 bg-green-600 text-white font-medium hover:bg-green-700 hover:border-green-800"
        >Save</button>
      </div>

    </template>
  </Dialog>
    
</template>