<script setup lang="ts">
import { useI18n } from 'vue-i18n'
import ConfirmDialog from 'primevue/confirmdialog'

// ─────────────────────────────────────────────
// Composables & stores
// ─────────────────────────────────────────────
const { t } = useI18n()
</script>

<template>
  <ConfirmDialog group="headless" class="px-4 pt-2 pb-4">
    <template #container="{ message, acceptCallback, rejectCallback }">
      <h2 v-if="message.header" class="mb-2">{{ message.header }}</h2>
      <p class="mb-2">{{ message.message }}</p>
      <div class="flex gap-2">
        <button
          @click="acceptCallback"
          class="w-1/2 rounded-md border border-gray-300 dark:border-gray-500 bg-gray-100 dark:bg-gray-700 px-4 py-2 font-medium text-gray-800 dark:text-gray-300 hover:bg-gray-200 dark:hover:bg-gray-600 active:bg-gray-300 dark:active:bg-gray-500"
        >
          {{ t('common.yes', 'Yes') }}
        </button>
        <button
          @click="rejectCallback"
          class="w-1/2 rounded-md border border-gray-300 dark:border-gray-500 bg-gray-100 dark:bg-gray-700 px-4 py-2 font-medium text-red-800 dark:text-red-400 hover:bg-gray-200 dark:hover:bg-gray-600 active:bg-gray-300 dark:active:bg-gray-500"
        >
          {{ t('common.no', 'No') }}
        </button>
      </div>
    </template>
  </ConfirmDialog>
</template>
