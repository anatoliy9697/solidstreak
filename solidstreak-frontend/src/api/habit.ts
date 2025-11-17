import axios from 'axios';
import type { Habit, HabitCheck } from '@/models/habit';

export interface Metadata {
  username: string;
}

export interface PostPutHabitRequest {
  data: Habit;
  meta?: Metadata;
}

export interface DeleteHabitRequest {
  meta?: Metadata;
}

export interface PostHabitCheckRequest {
  data: HabitCheck;
  meta?: Metadata;
}

export type ApiRequest = PostPutHabitRequest | DeleteHabitRequest | PostHabitCheckRequest

export interface Error {
  HTTPCode: number;
	Title: string;
	Detail?: string;
}

export interface PutHabitResponse {
  data: Habit;
  errors?: Error[];
}

export interface DeleteHabitResponse {
  data: Habit;
  errors?: Error[];
}

export interface GetHabitsResponse {
  data: Habit[];
  errors?: Error[];
}

export interface PostHabitCheckResponse {
  data: HabitCheck;
  errors?: Error[];
}

export type ApiResponse = PutHabitResponse | DeleteHabitResponse | GetHabitsResponse | PostHabitCheckResponse;

export interface RequestResult {
  success: boolean;
  httpCode: number;
  httpError: string | null;
  apiErrors: Error[];
  response: ApiResponse | null;
}

async function performRequest(method: 'post' | 'put' | 'delete' | 'get', url: string,  data?: ApiRequest): Promise<RequestResult> {
  if (data && !data.meta?.username) data.meta = { username: 'telegram_user' }; // TODO: должно приходить из внешнего контекста
  
  const result: RequestResult = {
    success: true,
    httpCode: 200,
    httpError: null,
    apiErrors: [],
    response: null,
  }

  try {

    const response = await axios.request<ApiResponse>({
      method,
      url,
      data,
      headers: {
        'X-Telegram-InitData': window.Telegram?.WebApp?.initData, // TODO: должно приходить из внешнего контекста
        'X-Request-ID': crypto.randomUUID(),
      },
    });

    result.response = response.data || null;

  } catch (error: unknown) {

    result.success = false; 
    if (typeof error === 'object' && error !== null && 'response' in error) {
      const err = error as { response?: { status?: number; data?: { errors?: Error[] } }; message?: string };
      result.httpCode = err.response?.status || 500;
      result.httpError = err.message || 'Unknown error';
      result.apiErrors = err.response?.data?.errors || [];
    } else {
      result.httpCode = 500;
      result.httpError = String(error);
    }

  }

  return result
}

export async function fetchHabits(userId: number): Promise<RequestResult> {
  return await performRequest('get', `/api/v1/users/${userId}/habits?with_checks=true`);
}

export async function postHabit(userId: number, habit: Habit): Promise<RequestResult> {
  const payload: PostPutHabitRequest = { data: habit };
  return await performRequest('post', `/api/v1/users/${userId}/habits`, payload);
}

export async function putHabit(userId: number, habit: Habit): Promise<RequestResult> {
  const payload: PostPutHabitRequest = { data: habit };
  return await performRequest('put', `/api/v1/users/${userId}/habits/${habit.id}`, payload);
}

export async function deleteHabit(userId: number, habitId: number): Promise<RequestResult> {
  const payload: DeleteHabitRequest = {};
  return await performRequest('delete', `/api/v1/users/${userId}/habits/${habitId}`, payload);
}

export async function postHabitCheck(userId: number, habitId: number, habitCheck: HabitCheck): Promise<RequestResult> {
  const payload: PostHabitCheckRequest = { data: habitCheck };
  return await performRequest('post', `/api/v1/users/${userId}/habits/${habitId}/checks`, payload);
}