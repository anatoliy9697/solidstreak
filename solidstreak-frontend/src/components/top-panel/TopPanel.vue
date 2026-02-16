<script setup lang="ts">
import { ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import Button from 'primevue/button'
import Select from 'primevue/select'
import { type Theme, getDefaultTheme, getNextTheme } from '@/models/theme'

import { LANGS, getDefaultLang } from '@/i18n'

// ─────────────────────────────────────────────
// Props
// ─────────────────────────────────────────────
const props = defineProps<{
  lang: string
  theme: Theme
}>()

// ─────────────────────────────────────────────
// Emits
// ─────────────────────────────────────────────
const emit = defineEmits<{
  (e: 'langSelected', lang: string): void
  (e: 'themeToggled', theme: Theme): void
}>()

// ─────────────────────────────────────────────
// Composables & stores
// ─────────────────────────────────────────────
const { t } = useI18n()

// ─────────────────────────────────────────────
// Constants & reactive state
// ─────────────────────────────────────────────
const selectedLang = ref(getDefaultLang())
const selectedTheme = ref(getDefaultTheme())

function toggleTheme() {
  selectedTheme.value = getNextTheme(selectedTheme.value)
  emit('themeToggled', selectedTheme.value)
}

watch(
  () => props.lang,
  (newLang) => {
    selectedLang.value = newLang
  },
  { immediate: true },
)
watch(
  () => props.theme,
  (newTheme) => {
    selectedTheme.value = newTheme
  },
  { immediate: true },
)
</script>

<template>
  <div class="mx-auto w-full max-w-lg px-2 py-2 text-right">
    <Button
      id="theme-toggle"
      :icon="selectedTheme === 'light' ? 'pi pi-moon' : 'pi pi-sun'"
      class="mr-2"
      @click="toggleTheme"
    />
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
