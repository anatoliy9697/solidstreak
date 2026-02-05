<script setup lang="ts">
import { ref, watch, computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { useConfirm } from 'primevue/useconfirm'
import { useToast } from 'primevue/usetoast'
import { SquarePen, Package, PackageOpen, Trash2 } from 'lucide-vue-next'

import { getHeatmapLocale } from '@/i18n'
import { dateToLocalString } from '@/utils/date'
import { useUserStore } from '@/stores/user'
import { useHabitStore } from '@/stores/habit'
import { COLORS, GREEN } from '@/models/color'
import type { Habit, HabitCheck } from '@/models/habit'
import CalendarHeatmap from '@/components/calendar-heatmap/CalendarHeatmap.vue'

// ─────────────────────────────────────────────
// Props
// ─────────────────────────────────────────────
const props = defineProps<{
  habit: Habit
  selectedDate: Date
  expanded?: boolean
}>()

// ─────────────────────────────────────────────
// Emits
// ─────────────────────────────────────────────
const emit = defineEmits<{
  (e: 'editHabit', habitId: number): void
  (e: 'expandHabitCard', habitId: number): void
  (e: 'collapseHabitCard', habitId: number): void
}>()

// ─────────────────────────────────────────────
// Composables & stores
// ─────────────────────────────────────────────
const { t } = useI18n()
const confirm = useConfirm()
const toast = useToast()
const userStore = useUserStore()
const habitStore = useHabitStore()

// ─────────────────────────────────────────────
// Constants & reactive state
// ─────────────────────────────────────────────
const selectedDateStr = ref<string>(dateToLocalString(props.selectedDate))
const isCheckButtonHovered = ref<boolean>(false)
const selectedDateChecked = ref<boolean>(
  props.habit.checks?.some(
    (check) => check.checkDate === dateToLocalString(props.selectedDate) && check.completed,
  ) || false,
)

watch(
  () => props.selectedDate,
  (newDate) => {
    selectedDateStr.value = dateToLocalString(newDate)
    selectedDateChecked.value =
      props.habit.checks?.some(
        (check) => check.checkDate === selectedDateStr.value && check.completed,
      ) || false
  },
)

// ─────────────────────────────────────────────
// Computed
// ─────────────────────────────────────────────
const checksArray = computed(() => {
  return (
    props.habit.checks
      ?.filter((check) => check.completed)
      .map((check) => ({
        date: check.checkDate,
        count: 1,
      })) || []
  )
})

const color = computed(() => COLORS[props.habit.color as keyof typeof COLORS] || GREEN)

// ─────────────────────────────────────────────
// Methods
// ─────────────────────────────────────────────
async function processCurrentDateCheck(): Promise<void> {
  const check = !selectedDateChecked.value

  const habitCheck: HabitCheck = {
    checkDate: selectedDateStr.value,
    completed: check,
    checkedAt: new Date(),
  }

  const result = await habitStore.setHabitCheck(userStore.id, props.habit.id, habitCheck)

  if (result.success) {
    selectedDateChecked.value = check
  } else {
    toast.add({
      severity: 'error',
      summary: t('common.error', 'Error'),
      detail: `${t('common.failedTo', 'Failed to')} ${check ? t('habitCard.check', 'check') : t('habitCard.uncheck', 'uncheck')} ${t('common.gcHabit', 'habit')}`,
      life: 3000,
    })
  }
}

async function processHabitArchiving(): Promise<void> {
  confirm.require({
    group: 'headless',
    message: `${t('common.areYouSure', 'Are you sure you want to')} ${props.habit.archived ? t('habitCard.unarchive', 'unarchive') : t('habitCard.archive', 'archive')} ${t('common.acHabit', 'habit')}?`,
    header: t('common.confirmation', 'Confirmation'),
    accept: async () => {
      const archive = !props.habit.archived
      const result = await habitStore.setHabitArchived(userStore.id, props.habit.id, archive)
      if (!result.success) {
        toast.add({
          severity: 'error',
          summary: t('common.error', 'Error'),
          detail: `${t('common.failedTo', 'Failed to')} ${archive ? t('habitCard.archive', 'archive') : t('habitCard.unarchive', 'unarchive')} ${t('common.acHabit', 'habit')}`,
          life: 3000,
        })
      }
    },
  })
}

async function processHabitDeletion(): Promise<void> {
  confirm.require({
    group: 'headless',
    message: `${t('common.areYouSure', 'Are you sure you want to')} ${t('habitCard.delete', 'delete')} ${t('common.acHabit', 'habit')}?`,
    header: t('common.confirmation', 'Confirmation'),
    accept: async () => {
      const result = await habitStore.deleteHabit(userStore.id, props.habit.id)
      if (!result.success) {
        toast.add({
          severity: 'error',
          summary: t('common.error', 'Error'),
          detail: `${t('common.failedTo', 'Failed to')} ${t('habitCard.delete', 'delete')} ${t('common.acHabit', 'habit')}`,
          life: 3000,
        })
      }
    },
  })
}
</script>

<template>
  <div
    :class="[
      'cursor-pointer rounded-md border border-gray-200 bg-white px-4 py-2 shadow-sm',
      habit.archived ? 'opacity-50' : '',
    ]"
  >
    <div
      :class="[
        'flex justify-between',
        expanded && !habit.archived ? 'mb-2' : '',
        expanded ? 'items-start' : '',
      ]"
    >
      <div
        class="mr-4 flex min-h-7 flex-1 items-center"
        @click.stop="
          props.expanded
            ? emit('collapseHabitCard', props.habit.id)
            : emit('expandHabitCard', props.habit.id)
        "
      >
        <h2 class="leading-none" style="word-break: break-word">{{ habit.title }}</h2>
      </div>

      <div class="flex items-center">
        <div v-if="expanded" class="flex items-center">
          <span :title="t('habitCard.upperDelete', 'Delete')">
            <Trash2
              @click.stop="processHabitDeletion()"
              class="mr-2 h-5 w-5 cursor-pointer text-gray-300 hover:text-gray-400"
            />
          </span>
          <span :title="t('habitCard.upperArchive', 'Archive')" v-if="!habit.archived">
            <Package
              @click.stop="processHabitArchiving()"
              class="mr-2 h-5 w-5 cursor-pointer text-gray-300 hover:text-gray-400"
            />
          </span>
          <span :title="t('habitCard.upperUnarchive', 'Unarchive')" v-else>
            <PackageOpen
              @click.stop="processHabitArchiving()"
              class="mr-2 h-5 w-5 cursor-pointer text-gray-300 hover:text-gray-400"
            />
          </span>
          <span :title="t('habitCard.upperEdit', 'Edit')">
            <SquarePen
              @click.stop="emit('editHabit', habit.id)"
              class="h-5 w-5 cursor-pointer text-gray-300 hover:text-gray-400"
            />
          </span>
        </div>

        <div v-if="!habit.archived" class="ml-4">
          <button
            @click.stop="processCurrentDateCheck()"
            @mouseover="isCheckButtonHovered = true"
            @mouseleave="isCheckButtonHovered = false"
            :style="
              selectedDateChecked
                ? {
                    borderColor: isCheckButtonHovered ? color.value500hex : color.value600hex,
                    backgroundColor: isCheckButtonHovered ? color.value400hex : color.value500hex,
                    color: '#fff',
                  }
                : {}
            "
            :class="[
              'flex h-7 w-7 cursor-pointer items-center justify-center rounded-lg border',
              selectedDateChecked
                ? ''
                : 'border-gray-400 text-gray-400 hover:border-gray-500 hover:text-gray-500',
            ]"
            :title="`${selectedDateChecked ? t('habitCard.upperUncheck', 'Uncheck') : t('habitCard.upperCheck', 'Check')} ${t('common.gcHabit', 'habit')} ${t('habitCard.forSelectedDate', 'for selected date')}`"
          >
            <svg
              xmlns="http://www.w3.org/2000/svg"
              class="h-6 w-6"
              fill="none"
              viewBox="0 0 24 24"
              stroke="currentColor"
              stroke-width="2"
            >
              <path stroke-linecap="round" stroke-linejoin="round" d="M5 13l4 4L19 7" />
            </svg>
          </button>
        </div>
      </div>
    </div>

    <div v-if="expanded && habit.description" class="mb-2">
      <p style="word-break: break-word; white-space: pre-wrap">{{ habit.description }}</p>
    </div>

    <CalendarHeatmap
      v-if="expanded && !habit.archived"
      :values="checksArray"
      :endDate="dateToLocalString(new Date())"
      :max="1"
      :rangeColor="[color.value100hex, color.value600hex]"
      :tooltip="false"
      :locale="getHeatmapLocale(t)"
      :round="3"
    />
  </div>
</template>
