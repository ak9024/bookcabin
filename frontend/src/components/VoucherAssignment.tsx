import { useState, type FormEvent } from 'react';
import { Button } from '@/components/ui/button';
import { voucherService } from '@/services/voucher.service';
import { ApiError } from '@/services/api';
import type { VoucherAssignmentData } from '@/types/voucher.types';

export function VoucherAssignment() {
  const [voucherCode, setVoucherCode] = useState('');
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [assignmentResult, setAssignmentResult] = useState<VoucherAssignmentData | null>(null);

  const handleSubmit = async (e: FormEvent<HTMLFormElement>) => {
    e.preventDefault();

    if (!voucherCode.trim()) {
      setError('Please enter a voucher code');
      return;
    }

    setIsLoading(true);
    setError(null);
    setAssignmentResult(null);

    try {
      const response = await voucherService.assignVoucher(voucherCode);
      setAssignmentResult(response);
      setVoucherCode('');
    } catch (err) {
      if (err instanceof ApiError) {
        setError(err.message);
      } else {
        setError('Failed to assign voucher. Please try again.');
      }
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="w-full max-w-md p-6 bg-white rounded-lg shadow-md">
      <h2 className="text-2xl font-bold mb-6 text-gray-800">Assign Voucher</h2>

      <form onSubmit={handleSubmit} className="space-y-4">
        <div>
          <label
            htmlFor="voucher-code"
            className="block text-sm font-medium text-gray-700 mb-2"
          >
            Voucher Code
          </label>
          <input
            id="voucher-code"
            type="text"
            value={voucherCode}
            onChange={(e) => setVoucherCode(e.target.value)}
            placeholder="e.g., VC-XX-XX"
            className="w-full px-4 py-2 border border-gray-300 rounded-md focus:ring-2 focus:ring-blue-500 focus:border-transparent outline-none"
            disabled={isLoading}
          />
        </div>

        {error && (
          <div className="p-3 bg-red-50 border border-red-200 rounded-md">
            <p className="text-sm text-red-600">{error}</p>
          </div>
        )}

        {assignmentResult && (
          <div className="p-4 bg-green-50 border border-green-200 rounded-md space-y-2">
            <p className="text-sm font-semibold text-green-800">✓ Voucher Assigned Successfully!</p>
            <div className="text-sm text-gray-700 space-y-1">
              <p><span className="font-medium">Voucher:</span> {assignmentResult.voucher_code}</p>
              <p><span className="font-medium">Cabin:</span> {assignmentResult.cabin}</p>
              <p><span className="font-medium">Seat:</span> {assignmentResult.seat_label}</p>
            </div>
          </div>
        )}

        <Button
          type="submit"
          className="w-full"
          disabled={isLoading}
        >
          {isLoading ? 'Assigning...' : 'Assign Voucher'}
        </Button>
      </form>
    </div>
  );
}
