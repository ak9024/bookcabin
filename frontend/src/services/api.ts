import type { ApiResponse, ApiErrorResponse } from '@/types/api.types';

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080';

export class ApiError extends Error {
  statusCode: number;
  data?: unknown;

  constructor(
    statusCode: number,
    message: string,
    data?: unknown
  ) {
    super(message);
    this.name = 'ApiError';
    this.statusCode = statusCode;
    this.data = data;
  }
}

interface FetchOptions extends RequestInit {
  timeout?: number;
}

export async function apiClient<T>(
  endpoint: string,
  options: FetchOptions = {}
): Promise<T> {
  const { timeout = 30000, ...fetchOptions } = options;

  const controller = new AbortController();
  const timeoutId = setTimeout(() => controller.abort(), timeout);

  try {
    const response = await fetch(`${API_BASE_URL}${endpoint}`, {
      ...fetchOptions,
      signal: controller.signal,
      headers: {
        'Content-Type': 'application/json',
        ...fetchOptions.headers,
      },
    });

    clearTimeout(timeoutId);

    const body = await response.json().catch(() => null);

    if (!response.ok) {
      if (body && typeof body === 'object' && 'status_code' in body && 'data' in body) {
        const errorResponse = body as ApiErrorResponse;
        const errorMessage = errorResponse.data || `HTTP error! status: ${errorResponse.status_code}`;

        throw new ApiError(
          errorResponse.status_code,
          errorMessage,
          errorResponse.data
        );
      }

      throw new ApiError(
        response.status,
        body?.message || `HTTP error! status: ${response.status}`,
        body
      );
    }

    if (body && typeof body === 'object' && 'status_code' in body && 'data' in body) {
      const apiResponse = body as ApiResponse<T>;
      return apiResponse.data;
    }

    return body as T;
  } catch (error) {
    clearTimeout(timeoutId);

    if (error instanceof ApiError) {
      throw error;
    }

    if (error instanceof Error) {
      if (error.name === 'AbortError') {
        throw new ApiError(408, 'Request timeout');
      }
      throw new ApiError(0, error.message);
    }

    throw new ApiError(0, 'An unknown error occurred');
  }
}
