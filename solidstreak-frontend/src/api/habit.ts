import axios from 'axios';
import type { Habit } from '@/models/habit';

export interface Error {
  HTTPCode: number;
	Title: string;
	Detail?: string;
}

export interface GetHabitsResponse {
  data: Habit[];
  errors?: Error[];
}

export type ApiResponse = GetHabitsResponse

export interface RequestResult {
  success: boolean;
  httpCode: number;
  httpError: string | null;
  apiErrors: Error[];
  response: ApiResponse | null;
}

// async function performRequest<T>(method: 'get' | 'post' | 'put', url: string, data?: any): Promise<T> {
async function performRequest(method: 'get' | 'post' | 'put', url: string): Promise<RequestResult> {
  const result: RequestResult = {
    success: true,
    httpCode: 200,
    httpError: null,
    apiErrors: [],
    response: null,
  }

  try {

    console.log(window.Telegram?.WebApp?.initData);
    const response = await axios.request<ApiResponse>({
      method,
      url,
      // data,
      headers: {
        'X-Telegram-InitData': window.Telegram?.WebApp?.initData,
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