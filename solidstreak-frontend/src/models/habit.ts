export interface HabitCheck {
    checkDate: string
	completed: boolean
	checkedAt: Date
}

export interface Habit {
    id: number
    archived: boolean
    title: string
    description?: string
    color?: string
    isPublic: boolean
    createdAt?: Date
    updatedAt?: Date
    checks?: HabitCheck[]
}