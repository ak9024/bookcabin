/**
 * Standard API response wrapper for all endpoints
 */
export interface ApiResponse<T> {
  status_code: number;
  data: T;
}

export interface ApiErrorResponse {
  status_code: number;
  data: string;
}
