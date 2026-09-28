<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { useToast } from 'primevue/usetoast'
import Dialog from 'primevue/dialog'
import RadioButton from 'primevue/radiobutton'
import ProgressSpinner from 'primevue/progressspinner'

import type { RequestResult } from '@/api/request'
import { type Pricing } from '@/models/subscription'
import { useSubscriptionStore } from '@/stores/subscription'
import { useInvoiceStore } from '@/stores/invoice'
import { useUserStore } from '@/stores/user'

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
const toast = useToast()
const subscriptionStore = useSubscriptionStore()
const userStore = useUserStore()
const invoiceStore = useInvoiceStore()
const init = ref<boolean>(true)

// ─────────────────────────────────────────────
// Constants & reactive state
// ─────────────────────────────────────────────
const selectedPricing = ref<Pricing | null>(null)

watch(
  () => props.visible,
  async (visible) => {
    if (visible) {

      init.value = true
      selectedPricing.value = premiumSubscriptionPricing.value[0] ?? null
      
      if (!invoiceStore.pendingInvoice) {

        const fetchPendingInvoiceResult = await invoiceStore.fetchPendingInvoice(userStore.id)
        if (!fetchPendingInvoiceResult.success) {
          toast.add({
            severity: 'error',
            summary: t('common.error', 'Error'),
            detail: t('subscriptionPurchaseDialog.dialogInitFailMsg', 'Failed to initialize the subscription purchase dialog'),
            life: 3000,
          })
          emit('closeSubscriptionPurchaseDialog')
        }

        if (invoiceStore.pendingInvoice) invoiceStore.startPollingInvoiceWhilePending()

      }

      init.value = false

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

// ─────────────────────────────────────────────
// Methods
// ─────────────────────────────────────────────
async function processSubscriptionPurchase(): Promise<void> {
  if (!selectedPricing.value) {
    return
  }

  let result: RequestResult
  result = await invoiceStore.createInvoice(
    userStore.id,
    'premium',
    selectedPricing.value.periodUnit,
    selectedPricing.value.periodCount,
    selectedPricing.value.currency,
  )

  if (result.success) {
    invoiceStore.startPollingInvoiceWhilePending()
  } else {
    toast.add({
      severity: 'error',
      summary: t('common.error', 'Error'),
      detail: `${t('common.failedTo', 'Failed to')} ${t('subscriptionPurchaseDialog.createInvoice', 'create invoice')}`,
      life: 3000,
    })
  }
}

function millisecondsToMMSSString(time: number): string {
  if (time <= 0) return '00:00'

  const minutes = Math.floor(time / 60000)
  const seconds = Math.floor((time % 60000) / 1000)

  return `${minutes.toString().padStart(2, '0')}:${seconds.toString().padStart(2, '0')}`
}
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
      <!-- Form header -->

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

      <!-- Form content -->

      <template v-if="init">
        <div class="mb-4 flex justify-center">
          <ProgressSpinner />
        </div>
      </template>

      <template v-else-if="invoiceStore.pendingInvoice">
        <span class="mb-2">{{
          t(
            'subscriptionPurchaseDialog.sentInvoiceMsg',
            'An invoice has been sent to your private messages',
          )
        }}</span>
        <span class="mb-2">{{
          t(
            'subscriptionPurchaseDialog.pendingInvoiceMsg',
            'It must be paid before the timer below expires, otherwise the invoice will be automatically canceled',
          )
        }}</span>
        <span class="mb-4 text-center text-4xl">
          {{ millisecondsToMMSSString(invoiceStore.millisecondsLeftUntilExpiration) }}
        </span>
      </template>

      <template v-else-if="invoiceStore.paidInvoice">
        <span class="mb-2">{{
          t(
            'subscriptionPurchaseDialog.purchasedSubscriptionMsg',
            'Congratulations! Your subscription has been successfully purchased',
          )
        }}</span>
        <span class="mb-6 text-center text-9xl">🎉</span>
      </template>

      <template v-else>
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
          <li
            v-if="!subscriptionStore.planByCode('premium')!.showAds"
            class="mb-1 flex items-start"
          >
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
      </template>

      <!-- Buttons -->

      <template v-if="init || invoiceStore.pendingInvoice || invoiceStore.paidInvoice">
        <button
          @click="closeCallback"
          class="rounded-md border border-gray-200 bg-gray-100 px-4 py-2 font-medium text-gray-800 hover:border-gray-300 hover:bg-gray-200 active:bg-gray-300 dark:border-gray-500 dark:bg-gray-600 dark:text-white dark:hover:border-gray-400 dark:hover:bg-gray-500 dark:active:bg-gray-400"
        >
          {{ t('common.close', 'Close') }}
        </button>
      </template>

      <template v-else>
        <div class="flex gap-2">
          <button
            @click="closeCallback"
            class="w-1/2 rounded-md border border-gray-200 bg-gray-100 px-4 py-2 font-medium text-gray-800 hover:border-gray-300 hover:bg-gray-200 active:bg-gray-300 dark:border-gray-500 dark:bg-gray-600 dark:text-white dark:hover:border-gray-400 dark:hover:bg-gray-500 dark:active:bg-gray-400"
          >
            {{ t('common.upperCancel', 'Cancel') }}
          </button>
          <button
            @click="processSubscriptionPurchase"
            class="w-1/2 rounded-md border border-green-700 bg-green-600 px-4 py-2 font-medium text-white hover:border-green-800 hover:bg-green-700 active:bg-green-800 dark:border-green-600 dark:bg-green-700 dark:hover:border-green-500 dark:hover:bg-green-600 dark:active:bg-green-500"
          >
            {{ t('subscriptionPurchaseDialog.purchase', 'Purchase') }}
          </button>
        </div>
      </template>
    </template>
  </Dialog>
</template>
