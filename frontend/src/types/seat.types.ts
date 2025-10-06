import type { ApiResponse } from './api.types';

export interface Seat {
  id: number;
  flight_id: number;
  label: string;
  cabin: string;
}

export type SeatListResponse = ApiResponse<Seat[]>;

export interface CreateSeatRequest {
  flight_id: number;
  cabin: string;
  labels: string[];
}
