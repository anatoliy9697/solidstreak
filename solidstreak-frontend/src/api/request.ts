import axios from 'axios'

import type { User } from '@/models/user'
import type { Habit, HabitCheck } from '@/models/habit'

interface Metadata {
  username: string
}

export interface Error {
  HTTPCode: number
  Title: string
  Detail?: string
}

export interface UserInfoData {
  user: User
  tgChat: {
    tgId: number
  }
}
export interface PostUserInfoRequest {
  data: UserInfoData
  meta?: Metadata
}
export interface PostPutHabitRequest {
  data: Habit
  meta?: Metadata
}
export interface DeleteHabitRequest {
  meta?: Metadata
}
export interface PostHabitCheckRequest {
  data: HabitCheck
  meta?: Metadata
}

type ApiRequest =
  | PostUserInfoRequest
  | PostPutHabitRequest
  | DeleteHabitRequest
  | PostHabitCheckRequest

export interface PostUserInfoResponse {
  data: User
  errors?: Error[]
}
export interface PutHabitResponse {
  data: Habit
  errors?: Error[]
}
export interface DeleteHabitResponse {
  data: Habit
  errors?: Error[]
}
export interface GetHabitsResponse {
  data: Habit[]
  errors?: Error[]
}
export interface PostHabitCheckResponse {
  data: HabitCheck
  errors?: Error[]
}

type ApiResponse =
  | PostUserInfoResponse
  | PutHabitResponse
  | DeleteHabitResponse
  | GetHabitsResponse
  | PostHabitCheckResponse

export interface RequestResult {
  success: boolean
  httpCode: number
  httpError: string | null
  apiErrors: Error[]
  response: ApiResponse | null
}

async function performRequest(
  method: 'post' | 'put' | 'delete' | 'get',
  url: string,
  initData: string,
  data?: ApiRequest,
): Promise<RequestResult> {
  const result: RequestResult = {
    success: true,
    httpCode: 200,
    httpError: null,
    apiErrors: [],
    response: null,
  }

  try {
    const response = await axios.request({
      method,
      url,
      data,
      headers: {
        'X-Telegram-InitData': initData,
        'X-Request-ID': crypto.randomUUID(),
      },
    })

    result.response = response.data || null
  } catch (error: unknown) {
    result.success = false
    if (typeof error === 'object' && error !== null && 'response' in error) {
      const err = error as {
        response?: { status?: number; data?: { errors?: Error[] } }
        message?: string
      }
      result.httpCode = err.response?.status || 500
      result.httpError = err.message || 'Unknown error'
      result.apiErrors = err.response?.data?.errors || []
    } else {
      result.httpCode = 500
      result.httpError = String(error)
    }
  }

  return result
}

export class ApiFetcher {
  initData: string
  username: string | undefined

  constructor(initData: string, username?: string) {
    this.initData = initData
    this.username = username
  }

  async upsertUserInfo(user: User, tgChat: { tgId: number }): Promise<RequestResult> {
    const payload: PostUserInfoRequest = { data: { user, tgChat: tgChat } }
    if (this.username) {
      payload.meta = { username: this.username } as Metadata
    }
    return await performRequest('post', `/api/v1/user-info/upsert`, this.initData, payload)
  }

  async fetchHabits(userId: number): Promise<RequestResult> {
    return await performRequest(
      'get',
      `/api/v1/users/${userId}/habits?with_checks=true`,
      this.initData,
    )
  }

  async postHabit(userId: number, habit: Habit): Promise<RequestResult> {
    const payload: PostPutHabitRequest = { data: habit }
    if (this.username) {
      payload.meta = { username: this.username } as Metadata
    }
    return await performRequest('post', `/api/v1/users/${userId}/habits`, this.initData, payload)
  }

  async putHabit(userId: number, habit: Habit): Promise<RequestResult> {
    const payload: PostPutHabitRequest = { data: habit }
    if (this.username) {
      payload.meta = { username: this.username } as Metadata
    }
    return await performRequest(
      'put',
      `/api/v1/users/${userId}/habits/${habit.id}`,
      this.initData,
      payload,
    )
  }

  async deleteHabit(userId: number, habitId: number): Promise<RequestResult> {
    const payload: DeleteHabitRequest = {}
    if (this.username) {
      payload.meta = { username: this.username } as Metadata
    }
    return await performRequest(
      'delete',
      `/api/v1/users/${userId}/habits/${habitId}`,
      this.initData,
      payload,
    )
  }

  async postHabitCheck(
    userId: number,
    habitId: number,
    habitCheck: HabitCheck,
  ): Promise<RequestResult> {
    const payload: PostHabitCheckRequest = { data: habitCheck }
    if (this.username) {
      payload.meta = { username: this.username } as Metadata
    }
    return await performRequest(
      'post',
      `/api/v1/users/${userId}/habits/${habitId}/checks`,
      this.initData,
      payload,
    )
  }
}
