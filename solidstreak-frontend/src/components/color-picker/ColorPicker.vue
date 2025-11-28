<script setup lang="ts">
import { type Color, COLORS } from '@/models/color'

// ─────────────────────────────────────────────
// Props
// ─────────────────────────────────────────────
defineProps<{
  selectedColor: Color
}>()

// ─────────────────────────────────────────────
// Emits
// ─────────────────────────────────────────────
const emit = defineEmits<{
  (e: 'colorSelected', color: Color): void
}>()
</script>

<template>
  <div class="flex flex-wrap gap-2">
    <span
      v-for="color in Object.values(COLORS)"
      :key="color.name"
      @click.stop="emit('colorSelected', color)"
      class="inline-block h-5 w-5 cursor-pointer rounded-full text-center font-semibold text-white hover:ring-2 hover:ring-gray-200"
      :style="{ backgroundColor: color.value500hex }"
      >{{ selectedColor.name === color.name ? '✓' : '' }}</span
    >
  </div>
</template>
