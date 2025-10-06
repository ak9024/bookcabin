import { apiClient } from './api';
import type { AssignVoucherRequest, VoucherAssignmentData } from '@/types/voucher.types';

export const voucherService = {
  assignVoucher: async (voucherCode: string): Promise<VoucherAssignmentData> => {
    const requestBody: AssignVoucherRequest = {
      voucher_code: voucherCode,
    };

    return apiClient<VoucherAssignmentData>('/api/v1/vouchers/assigns', {
      method: 'POST',
      body: JSON.stringify(requestBody),
    });
  },
};
