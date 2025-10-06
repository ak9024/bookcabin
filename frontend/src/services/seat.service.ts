import { apiClient } from './api';
import type { Seat, CreateSeatRequest } from '@/types/seat.types';

export const seatService = {
  getAllSeats: async (): Promise<Seat[]> => {
    const data = await apiClient<Seat[]>('/api/v1/seats', {
      method: 'GET',
    });
    return data ?? [];
  },

  createSeats: async (request: CreateSeatRequest): Promise<string> => {
    return apiClient<string>('/api/v1/seats', {
      method: 'POST',
      body: JSON.stringify(request),
    });
  },
};
