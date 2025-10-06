import { apiClient } from './api';
import type { AssignVoucherRequest, VoucherAssignmentData, Voucher, CreateVoucherRequest } from '@/types/voucher.types';

export const voucherService = {
  getAllVouchers: async (): Promise<Voucher[]> => {
    const data = await apiClient<Voucher[]>('/api/v1/vouchers', {
      method: 'GET',
    });
    return data ?? [];
  },

  assignVoucher: async (voucherCode: string): Promise<VoucherAssignmentData> => {
    const requestBody: AssignVoucherRequest = {
      voucher_code: voucherCode,
    };

    return apiClient<VoucherAssignmentData>('/api/v1/vouchers/assigns', {
      method: 'POST',
      body: JSON.stringify(requestBody),
    });
  },

  createVoucher: async (requestBody: CreateVoucherRequest): Promise<Voucher> => {
    return apiClient<Voucher>('/api/v1/vouchers', {
      method: 'POST',
      body: JSON.stringify(requestBody),
    });
  },
};
