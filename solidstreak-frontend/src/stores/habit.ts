import { defineStore } from 'pinia';

import type { Habit, HabitCheck } from '@/models/habit';
import { type RequestResult, fetchHabits, postHabit, putHabit, deleteHabit, postHabitCheck } from '@/api/habit';

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

    async createHabit(userId: number, habit: Habit): Promise<RequestResult> {
      const result = await postHabit(userId, habit);

      if (result.success) {
        const createdHabit = result.response?.data as Habit;
        this.habitsMap.set(createdHabit.id, createdHabit);
        this.habits.push(createdHabit);
      }

      return result;
    },

    async updateHabit(userId: number, habit: Habit): Promise<RequestResult> {
      const result = await putHabit(userId, habit);

      if (result.success) {
        const updatedHabit = result.response?.data as Habit;
        
        habit = this.habitsMap.get(updatedHabit.id)!;
        
        habit.title = updatedHabit.title;
        habit.description = updatedHabit.description;
        habit.color = updatedHabit.color;
        habit.isPublic = updatedHabit.isPublic;
        habit.archived = updatedHabit.archived;
        habit.updatedAt = updatedHabit.updatedAt;

        // this.habitsMap.set(updatedHabit.id, habit);

        // const index = this.habits.findIndex(h => h.id === updatedHabit.id);
        // if (index !== -1) {
        //   this.habits[index] = habit;
        // }
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

      const updatedHabit = {
        id: habit.id,
        title: habit.title,
        description: habit.description,
        color: habit.color,
        isPublic: habit.isPublic,
        archived: archived
      } as Habit;
      
      const result = await this.updateHabit(userId, updatedHabit);

      // if (result.success) this.habits = [...this.habits];

      return result;
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
        if (habit && !habit.checks) {
          habit.checks = [];
        }
        habit?.checks?.push(habitCheck);
      }

      return result;
    },

  },

  getters: {
    
    activeHabits(state): Habit[] {
      return state.habits.filter(habit => !habit.archived);
    },

    archivedHabits(state): Habit[] {
      return state.habits.filter(habit => habit.archived);
    },

    activeHabitsCount(state): number {
      return state.habits.filter(habit => !habit.archived).length;
    },

    activities(state): { date: string; count: number }[] {
      const map = new Map<string, number>();
  
      state.habits
        .filter(habit => !habit.archived)
        .forEach(habit => {
          habit.checks?.forEach(check => {
            if (check.completed) {
              const date = check.checkDate;
              map.set(date, (map.get(date) || 0) + 1);
            }
          });
        });

      return Array.from(map.entries()).map(([date, count]) => ({ date, count }));
    },

    habitById(state): (id: number) => Habit | undefined {
      return (id: number) => state.habitsMap.get(id);
    },

  },

});