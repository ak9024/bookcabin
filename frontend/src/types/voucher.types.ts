import type { ApiResponse } from './api.types';

export interface AssignVoucherRequest {
  voucher_code: string;
}

export interface CreateVoucherRequest {
  code: string;
  flight_id: number;
  cabin: string;
}

export interface VoucherAssignmentData {
  voucher_code: string;
  cabin: string;
  seat_id: number;
  seat_label: string;
}

export type AssignVoucherResponse = ApiResponse<VoucherAssignmentData>;

export interface Voucher {
  id: number;
  code: string;
  flight_id: number;
  cabin: string;
  expires_at: string;
  redeemed: number;
  redeemed_at?: string;
}

export type VoucherListResponse = ApiResponse<Voucher[]>;

export type CreateVoucherResponse = ApiResponse<Voucher>;
