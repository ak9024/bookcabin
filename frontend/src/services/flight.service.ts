import { apiClient } from './api';
import type { Flight, CreateFlightRequest } from '@/types/flight.types';

export const flightService = {
  getAllFlights: async (): Promise<Flight[]> => {
    const data = await apiClient<Flight[]>('/api/v1/flights', {
      method: 'GET',
    });
    return data ?? [];
  },

  createFlight: async (request: CreateFlightRequest): Promise<string> => {
    return apiClient<string>('/api/v1/flights', {
      method: 'POST',
      body: JSON.stringify(request),
    });
  },
};
