<script setup lang="ts">

import { type Color, COLORS } from '@/models/color'

// ─────────────────────────────────────────────
// Props
// ─────────────────────────────────────────────
defineProps<{
  selectedColor: Color
}>();

// ─────────────────────────────────────────────
// Emits
// ─────────────────────────────────────────────
const emit = defineEmits<{
  (e: 'colorSelected', color: Color): void
}>();

</script>

<template>

    <div class="flex flex-wrap gap-2">
        <span
            v-for="color in Object.values(COLORS)"
            :key="color.name"
            @click.stop="emit('colorSelected', color)"
            class="inline-block w-5 h-5 rounded-full cursor-pointer hover:ring-2 hover:ring-gray-200 text-center font-semibold text-white"
            :style="{ backgroundColor: color.value500hex }"
        >{{ selectedColor.name === color.name ? '✓' : '' }}</span>
    </div>

</template>