import { useState, useEffect, useMemo } from 'react';
import { Button } from '@/components/ui/button';
import { voucherService } from '@/services/voucher.service';
import { seatService } from '@/services/seat.service';
import { ApiError } from '@/services/api';
import type { CreateVoucherRequest } from '@/types/voucher.types';
import type { Seat } from '@/types/seat.types';

interface CreateVoucherProps {
  onVoucherCreated?: () => void;
}

export function CreateVoucher({ onVoucherCreated }: CreateVoucherProps) {
  const [formData, setFormData] = useState<CreateVoucherRequest>({
    code: '',
    flight_id: 0,
    cabin: 'ECONOMY',
  });
  const [seats, setSeats] = useState<Seat[]>([]);
  const [isLoadingSeats, setIsLoadingSeats] = useState(true);
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [success, setSuccess] = useState<string | null>(null);

  useEffect(() => {
    const fetchSeats = async () => {
      setIsLoadingSeats(true);
      try {
        const data = await seatService.getAllSeats();
        setSeats(data);
      } catch (err) {
        if (err instanceof ApiError) {
          setError(`Failed to load seats: ${err.message}`);
        } else {
          setError('Failed to load seats. Please refresh the page.');
        }
      } finally {
        setIsLoadingSeats(false);
      }
    };

    fetchSeats();
  }, []);

  const uniqueFlights = useMemo(() => {
    const flightMap = new Map<number, Set<string>>();

    seats.forEach(seat => {
      if (!flightMap.has(seat.flight_id)) {
        flightMap.set(seat.flight_id, new Set());
      }
      flightMap.get(seat.flight_id)?.add(seat.cabin);
    });

    return Array.from(flightMap.entries()).map(([flight_id, cabins]) => ({
      flight_id,
      cabins: Array.from(cabins)
    }));
  }, [seats]);

  const availableCabins = useMemo(() => {
    if (!formData.flight_id) return [];

    const flight = uniqueFlights.find(f => f.flight_id === formData.flight_id);
    return flight?.cabins || [];
  }, [formData.flight_id, uniqueFlights]);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setIsSubmitting(true);
    setError(null);
    setSuccess(null);

    try {
      const voucher = await voucherService.createVoucher(formData);
      setSuccess(`Voucher created successfully!`);
      setFormData({ code: '', flight_id: 0, cabin: '' });

      if (onVoucherCreated) {
        onVoucherCreated();
      }
    } catch (err) {
      if (err instanceof ApiError) {
        setError(err.message);
      } else {
        setError('Failed to create voucher. Please try again.');
      }
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <div className="w-full max-w-md p-6 bg-white rounded-lg shadow-md">
      <h2 className="text-2xl font-bold text-gray-800 mb-6">Create New Voucher</h2>

      {error && (
        <div className="p-3 bg-red-50 border border-red-200 rounded-md mb-4">
          <p className="text-sm text-red-600">{error}</p>
        </div>
      )}

      {success && (
        <div className="p-3 bg-green-50 border border-green-200 rounded-md mb-4">
          <p className="text-sm text-green-600">{success}</p>
        </div>
      )}

      <form onSubmit={handleSubmit} className="space-y-4">
        <div>
          <label htmlFor="code" className="block text-sm font-medium text-gray-700 mb-1">
            Voucher Code
          </label>
          <input
            type="text"
            id="code"
            value={formData.code}
            onChange={(e) => setFormData({ ...formData, code: e.target.value })}
            className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
            placeholder="e.g., VC-XX-XX"
            required
          />
        </div>

        <div>
          <label htmlFor="flight_id" className="block text-sm font-medium text-gray-700 mb-1">
            Flight
          </label>
          <select
            id="flight_id"
            value={formData.flight_id || ''}
            onChange={(e) => {
              const flight_id = parseInt(e.target.value);
              setFormData({
                ...formData,
                flight_id,
                cabin: '' // Reset cabin when flight changes
              });
            }}
            className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
            required
            disabled={isLoadingSeats}
          >
            <option value="">
              {isLoadingSeats ? 'Loading flights...' : 'Select a flight'}
            </option>
            {uniqueFlights.map((flight) => (
              <option key={flight.flight_id} value={flight.flight_id}>
                Flight ID: {flight.flight_id}
              </option>
            ))}
          </select>
        </div>

        <div>
          <label htmlFor="cabin" className="block text-sm font-medium text-gray-700 mb-1">
            Cabin Class
          </label>
          <select
            id="cabin"
            value={formData.cabin}
            onChange={(e) => setFormData({ ...formData, cabin: e.target.value })}
            className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
            required
            disabled={!formData.flight_id}
          >
            <option value="">
              {!formData.flight_id ? 'Select a flight first' : 'Select a cabin'}
            </option>
            {availableCabins.map((cabin) => (
              <option key={cabin} value={cabin}>
                {cabin.charAt(0) + cabin.slice(1).toLowerCase()}
              </option>
            ))}
          </select>
        </div>

        <Button
          type="submit"
          disabled={isSubmitting}
          className="w-full py-2"
        >
          {isSubmitting ? 'Creating...' : 'Create Voucher'}
        </Button>
      </form>
    </div>
  );
}
