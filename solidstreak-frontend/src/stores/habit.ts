import { defineStore } from 'pinia';
import type { Habit, HabitCheck } from '@/models/habit';
import { fetchHabits, putHabit, deleteHabit, postHabitCheck } from '@/api/habit';
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

    async updateHabit(userId: number, habit: Habit): Promise<RequestResult> {
      const result = await putHabit(userId, habit);

      if (result.success) {
        const updatedHabit = result.response?.data as Habit;
        updatedHabit.checks = habit.checks;
        this.habitsMap.set(updatedHabit.id, updatedHabit);
        const index = this.habits.findIndex(h => h.id === updatedHabit.id);
        if (index !== -1) {
          this.habits[index] = updatedHabit;
        }
      }

      return result;
    },

    async setHabitArchived(userId: number, habitId: number, archived: boolean): Promise<RequestResult> {
      const habit = this.habitsMap.get(habitId);

      if (!habit) {
        return {
          success: false,
          httpCode: 404,
          httpError: 'Habit not found',
          apiErrors: [{
            HTTPCode: 404,
            Title: 'not found',
            Detail: `couldn't find habit with specified id`,
          }],
          response: null,
        };
      }

      const updatedHabit = { ...habit, archived };
      
      return await this.updateHabit(userId, updatedHabit);
    },

    async deleteHabit(userId: number, habitId: number): Promise<RequestResult> {
      const result = await deleteHabit(userId, habitId);

      if (result.success) {
        this.habitsMap.delete(habitId);
        this.habits = this.habits.filter(h => h.id !== habitId);
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
    },

  }

});