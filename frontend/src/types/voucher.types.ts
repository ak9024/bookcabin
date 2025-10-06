import type { ApiResponse } from './api.types';

export interface AssignVoucherRequest {
  voucher_code: string;
}

export interface VoucherAssignmentData {
  voucher_code: string;
  cabin: string;
  seat_id: number;
  seat_label: string;
}

export type AssignVoucherResponse = ApiResponse<VoucherAssignmentData>;
