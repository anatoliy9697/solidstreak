<script setup lang="ts">
import { ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import Button from 'primevue/button'
import Select from 'primevue/select'
import { type Theme, getDefaultTheme, getNextTheme } from '@/models/theme'
import { useUserStore } from '@/stores/user'

import { LANGS, getDefaultLang } from '@/i18n'
import { dateToDDMMYYYYString } from '@/utils/date'

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
const userStore = useUserStore()

// ─────────────────────────────────────────────
// Constants & reactive state
// ─────────────────────────────────────────────
const selectedLang = ref(getDefaultLang())
const selectedTheme = ref(getDefaultTheme())

// ─────────────────────────────────────────────
// Methods
// ─────────────────────────────────────────────
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
  <div class="mx-auto flex w-full max-w-lg items-center justify-between px-2 py-2">
    <div class="px-2 text-left">
      <span
        v-if="userStore._subscription && userStore._subscription!.active"
        class="text-gray-600 dark:text-gray-400"
      >
        {{
          userStore._subscription!.plan.code === 'basic'
            ? t('common.basicSubscriptionUpper', 'Basic subscription')
            : t('common.premiumSubscriptionUpper', 'Premium subscription')
        }}
        {{
          userStore._subscription!.finishDt
            ? t('common.until', 'until') +
              ' ' +
              dateToDDMMYYYYString(userStore._subscription!.finishDt)
            : ''
        }}
        <template
          v-if="
            userStore._subscription!.plan.code === 'basic' ||
            (userStore._subscription!.plan.code === 'premium' && userStore._subscription!.finishDt)
          "
        >
          (<a
            href="#"
            :title="
              userStore._subscription!.plan.code === 'basic'
                ? t('topPanel.buyPremiumSubscription', 'Buy premium subscription')
                : t('topPanel.renewPremiumSubscription', 'Renew premium subscription')
            "
            >{{
              userStore._subscription!.plan.code === 'basic'
                ? t('topPanel.buyPremium', 'buy premium')
                : t('topPanel.renew', 'renew')
            }}</a
          >)</template
        >
      </span>
    </div>

    <div class="flex items-center">
      <Button
        id="theme-toggle"
        :icon="selectedTheme !== 'dark' ? 'pi pi-moon' : 'pi pi-sun'"
        class="mr-2"
        @click="toggleTheme"
        :title="
          selectedTheme !== 'dark'
            ? t('topPanel.darkTheme', 'Turn on dark theme')
            : t('topPanel.lightTheme', 'Turn on light theme')
        "
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
  </div>
</template>

<style scoped></style>
