<script setup lang="ts">
import { ref, watch } from 'vue'
import { LANGS } from '@/i18n'
import Select from 'primevue/select'

// ─────────────────────────────────────────────
// Props
// ─────────────────────────────────────────────
const props = defineProps<{
  lang: string
}>()

// ─────────────────────────────────────────────
// Emits
// ─────────────────────────────────────────────
const emit = defineEmits<{
  (e: 'langSelected', lang: string): void
}>()

// ─────────────────────────────────────────────
// Constants & reactive state
// ─────────────────────────────────────────────
const selectedLang = ref('ru')

watch(
  () => props.lang,
  (newLang) => {
    selectedLang.value = newLang
  },
  { immediate: true },
)
</script>

<template>
  <div class="mx-auto w-full max-w-lg px-2 py-2 text-right">
    <Select
      v-model="selectedLang"
      @update:modelValue="(value) => emit('langSelected', value)"
      :options="LANGS"
      optionLabel="name"
      optionValue="code"
    />
  </div>
</template>

<style scoped></style>
