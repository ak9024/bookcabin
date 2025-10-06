import type { ApiResponse } from './api.types';

export interface Flight {
  id: number;
  flight_no: string;
  dep_date: string;
}

export type FlightListResponse = ApiResponse<Flight[]>;

export interface CreateFlightRequest {
  flight_numbers: string[];
  dep_date: string;
}
