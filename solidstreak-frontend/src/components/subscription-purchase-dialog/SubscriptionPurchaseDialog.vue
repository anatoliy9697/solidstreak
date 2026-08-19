<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import Dialog from 'primevue/dialog'
import RadioButton from 'primevue/radiobutton'
import { useI18n } from 'vue-i18n'

import { useSubscriptionStore } from '@/stores/subscription'
import { type Pricing } from '@/models/subscription'

// ─────────────────────────────────────────────
// Props
// ─────────────────────────────────────────────
const props = defineProps<{
  visible: boolean
}>()

// ─────────────────────────────────────────────
// Emits
// ─────────────────────────────────────────────
const emit = defineEmits<{
  (e: 'closeSubscriptionPurchaseDialog'): void
}>()

// ─────────────────────────────────────────────
// Composables & stores
// ─────────────────────────────────────────────
const { t } = useI18n()
const subscriptionStore = useSubscriptionStore()

// ─────────────────────────────────────────────
// Constants & reactive state
// ─────────────────────────────────────────────
const selectedPricing = ref<Pricing | null>(null)

watch(
  () => props.visible,
  (visible) => {
    if (visible) {
      selectedPricing.value = premiumSubscriptionPricing.value[0] ?? null
    }
  },
  { immediate: true },
)

// ─────────────────────────────────────────────
// Computed
// ─────────────────────────────────────────────
const dialogVisible = computed({
  get: () => props.visible,
  set: () => {
    emit('closeSubscriptionPurchaseDialog')
  },
})

const premiumSubscriptionPricing = computed(() =>
  [...(subscriptionStore.planByCode('premium')?.pricing ?? [])].sort(
    (a, b) => a.displayOrder - b.displayOrder,
  ),
)
</script>

<template>
  <Dialog
    v-model:visible="dialogVisible"
    position="bottom"
    :modal="true"
    :draggable="false"
    class="!m-2 !mb-0 w-[calc(100vw-1rem)] max-w-[31rem] !rounded-t-xl !rounded-b-none !border-b-0 p-4"
  >
    <template #container="{ closeCallback }">
      <div class="mb-4 flex items-start justify-between">
        <div>
          <h1>{{ t('subscriptionPurchaseDialog.premiumSubscription', 'Premium subscription') }}</h1>
        </div>
        <div>
          <i
            @click="closeCallback"
            class="pi pi-times cursor-pointer text-gray-300 hover:text-gray-400 dark:text-gray-500 dark:hover:text-gray-400"
            :title="t('common.close', 'Close')"
          ></i>
        </div>
      </div>

      <div
        v-for="pricingOption in premiumSubscriptionPricing"
        :key="pricingOption.periodCount + pricingOption.periodUnit"
        :class="[
          pricingOption.displayOrder == premiumSubscriptionPricing.length ? 'mb-4' : 'mb-2',
          'flex cursor-pointer items-center justify-between rounded-md border bg-gray-100 px-4 py-2 hover:bg-gray-200 dark:bg-gray-600 dark:hover:bg-gray-500',
          selectedPricing === pricingOption
            ? 'border-orange-600'
            : 'border-gray-200 hover:border-gray-300 dark:border-gray-500 dark:hover:border-gray-400',
        ]"
        @click="selectedPricing = pricingOption"
      >
        <span class="mr-2 text-lg font-semibold whitespace-nowrap">{{
          t('subscriptionPurchaseDialog.' + pricingOption.periodUnit + 'SubscriptionPeriod', {
            subscriptionPeriodCount: pricingOption.periodCount,
          })
        }}</span>
        <div class="flex items-center">
          <span class="mr-2 text-lg font-semibold">{{ pricingOption.price }}</span>
          <span class="mr-2 text-xs">{{
            t('subscriptionPurchaseDialog.' + pricingOption.currency)
          }}</span>
          <RadioButton
            v-model="selectedPricing"
            :inputId="pricingOption.periodCount + pricingOption.periodUnit"
            name="pricing"
            :value="pricingOption"
          />
        </div>
      </div>

      <h2 class="mb-2">
        {{ t('subscriptionPurchaseDialog.premiumBenefits', 'Premium subscription benefits') }}
      </h2>
      <ul class="mb-3 pl-2">
        <li class="mb-1 flex items-start">
          <span class="mr-2">📚</span>
          <span>{{
            t(
              'subscriptionPurchaseDialog.habitLimits',
              {
                premiumSubscriptionActiveHabitsLimit:
                  subscriptionStore.planByCode('premium')!.activeHabitsLimit,
                basicSubscriptionActiveHabitsLimit:
                  subscriptionStore.planByCode('basic')!.activeHabitsLimit,
              },
              'Up to {premiumSubscriptionActiveHabitsLimit} active habits, instead of {basicSubscriptionActiveHabitsLimit} in the basic subscription',
            )
          }}</span>
        </li>
        <li v-if="!subscriptionStore.planByCode('premium')!.showAds" class="mb-1 flex items-start">
          <span class="mr-2">📢</span>
          <span>{{ t('subscriptionPurchaseDialog.noAds', 'No ads') }}</span>
        </li>
        <li class="mb-1 flex items-start">
          <span class="mr-2">✨</span>
          <span>{{
            t('subscriptionPurchaseDialog.futureFeatures', 'Access to some upcoming features')
          }}</span>
        </li>
      </ul>

      <div class="flex gap-2">
        <button
          @click="closeCallback"
          class="w-1/2 rounded-md border border-gray-200 bg-gray-100 px-4 py-2 font-medium text-gray-800 hover:border-gray-300 hover:bg-gray-200 active:bg-gray-300 dark:border-gray-500 dark:bg-gray-600 dark:text-white dark:hover:border-gray-400 dark:hover:bg-gray-500 dark:active:bg-gray-400"
        >
          {{ t('common.upperCancel', 'Cancel') }}
        </button>
        <button
          class="w-1/2 rounded-md border border-green-700 bg-green-600 px-4 py-2 font-medium text-white hover:border-green-800 hover:bg-green-700 active:bg-green-800 dark:border-green-600 dark:bg-green-700 dark:hover:border-green-500 dark:hover:bg-green-600 dark:active:bg-green-500"
        >
          {{ t('subscriptionPurchaseDialog.purchase', 'Purchase') }}
        </button>
      </div>
    </template>
  </Dialog>
</template>
