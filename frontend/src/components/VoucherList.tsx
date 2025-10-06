import { useState, useEffect } from 'react';
import { Button } from '@/components/ui/button';
import { voucherService } from '@/services/voucher.service';
import { ApiError } from '@/services/api';
import type { Voucher } from '@/types/voucher.types';

export function VoucherList() {
  const [vouchers, setVouchers] = useState<Voucher[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const fetchVouchers = async () => {
    setIsLoading(true);
    setError(null);

    try {
      const data = await voucherService.getAllVouchers();
      setVouchers(data);
    } catch (err) {
      if (err instanceof ApiError) {
        setError(err.message);
      } else {
        setError('Failed to fetch vouchers. Please try again.');
      }
    } finally {
      setIsLoading(false);
    }
  };

  useEffect(() => {
    fetchVouchers();
  }, []);

  const formatDate = (dateString: string | undefined) => {
    if (!dateString) return '-';
    try {
      return new Date(dateString).toLocaleString();
    } catch {
      return dateString;
    }
  };

  return (
    <div className="w-full max-w-4xl p-6 bg-white rounded-lg shadow-md">
      <div className="flex items-center justify-between mb-6">
        <h2 className="text-2xl font-bold text-gray-800">All Vouchers</h2>
        <Button
          onClick={fetchVouchers}
          disabled={isLoading}
          className="px-4 py-2"
        >
          {isLoading ? 'Refreshing...' : 'Refresh'}
        </Button>
      </div>

      {error && (
        <div className="p-3 bg-red-50 border border-red-200 rounded-md mb-4">
          <p className="text-sm text-red-600">{error}</p>
        </div>
      )}

      {isLoading && !error ? (
        <div className="flex items-center justify-center py-12">
          <p className="text-gray-500">Loading vouchers...</p>
        </div>
      ) : vouchers.length === 0 ? (
        <div className="flex items-center justify-center py-12">
          <p className="text-gray-500">No vouchers found</p>
        </div>
      ) : (
        <div className="overflow-x-auto">
          <table className="w-full border-collapse">
            <thead>
              <tr className="bg-gray-50 border-b border-gray-200">
                <th className="px-4 py-3 text-left text-sm font-semibold text-gray-700">Code</th>
                <th className="px-4 py-3 text-left text-sm font-semibold text-gray-700">Flight ID</th>
                <th className="px-4 py-3 text-left text-sm font-semibold text-gray-700">Cabin</th>
                <th className="px-4 py-3 text-left text-sm font-semibold text-gray-700">Status</th>
                <th className="px-4 py-3 text-left text-sm font-semibold text-gray-700">Redeemed At</th>
                <th className="px-4 py-3 text-left text-sm font-semibold text-gray-700">Expires At</th>
              </tr>
            </thead>
            <tbody>
              {vouchers.map((voucher) => (
                <tr
                  key={voucher.id}
                  className="border-b border-gray-100 hover:bg-gray-50 transition-colors"
                >
                  <td className="px-4 py-3 text-sm font-medium text-gray-900">
                    {voucher.code}
                  </td>
                  <td className="px-4 py-3 text-sm text-gray-700">
                    {voucher.flight_id}
                  </td>
                  <td className="px-4 py-3 text-sm text-gray-700">
                    <span className="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-blue-100 text-blue-800">
                      {voucher.cabin}
                    </span>
                  </td>
                  <td className="px-4 py-3 text-sm">
                    {voucher.redeemed === 1 ? (
                      <span className="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-green-100 text-green-800">
                        ✓ Redeemed
                      </span>
                    ) : (
                      <span className="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-gray-100 text-gray-800">
                        Available
                      </span>
                    )}
                  </td>
                  <td className="px-4 py-3 text-sm text-gray-700">
                    {formatDate(voucher.redeemed_at)}
                  </td>
                  <td className="px-4 py-3 text-sm text-gray-700">
                    {formatDate(voucher.expires_at)}
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}

      <div className="mt-4 text-sm text-gray-500">
        Total vouchers: {vouchers.length}
      </div>
    </div>
  );
}
