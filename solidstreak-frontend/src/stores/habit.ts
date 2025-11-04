import { defineStore } from 'pinia';
import type { Habit, HabitCheck } from '@/models/habit';
import { fetchHabits, postHabitCheck } from '@/api/habit';
import type { RequestResult } from '@/api/habit';

export const useHabitStore = defineStore('habit', {
  
  state: () => ({
    habits: [] as Habit[],
    habitsMap: new Map<number, Habit>()
  }),

  actions: {

    async fetchHabits(userId: number): Promise<RequestResult> {
      const result = await fetchHabits(userId);

      const data = result.response?.data;
      this.habits = (data ? data as Habit[] : []);

      this.habitsMap.clear();
      for (const habit of this.habits) {
        this.habitsMap.set(habit.id, habit);
      }

      return result;
    },

    async setHabitCheck(userId: number, habitId: number, habitCheck: HabitCheck): Promise<RequestResult> {
      const result = await postHabitCheck(userId, habitId, habitCheck);

      if (result.success) {
        const habit = this.habitsMap.get(habitId);
        for (const check of habit?.checks || []) {
          if (check.checkDate === habitCheck.checkDate) {
            check.completed = habitCheck.completed;
            check.checkedAt = habitCheck.checkedAt;
            return result;
          }
        }
        habit?.checks?.push(habitCheck);
      }

      return result;
    }

  }

});