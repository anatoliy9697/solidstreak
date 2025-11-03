import { defineStore } from 'pinia';
import type { Habit } from '@/models/habit';
import { fetchHabits } from '@/api/habit';
import type { RequestResult } from '@/api/habit';

export const useHabitStore = defineStore('habit', {
  
  state: () => ({
    habits: [] as Habit[],
    requestResult: null as RequestResult | null,
  }),

  actions: {
    async fetchHabits(userId: number): Promise<void> {
      this.requestResult = await fetchHabits(userId);
      this.habits = this.requestResult.response?.data || [];
    }
  }
  
});