<script setup lang="ts">
import { ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import Select from 'primevue/select'

import { LANGS } from '@/i18n'

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
// Composables & stores
// ─────────────────────────────────────────────
const { t } = useI18n()

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
      :title="t('topPanel.selectLang', 'Select language')"
    />
  </div>
</template>

<style scoped></style>
