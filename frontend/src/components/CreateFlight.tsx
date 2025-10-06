import { useState } from 'react';
import { Button } from '@/components/ui/button';
import { flightService } from '@/services/flight.service';
import { ApiError } from '@/services/api';
import type { CreateFlightRequest } from '@/types/flight.types';

interface CreateFlightProps {
  onFlightCreated?: () => void;
}

export function CreateFlight({ onFlightCreated }: CreateFlightProps) {
  const [flightNumbers, setFlightNumbers] = useState<string>('');
  const [deptDate, setDeptDate] = useState<string>('');
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [success, setSuccess] = useState<string | null>(null);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setIsSubmitting(true);
    setError(null);
    setSuccess(null);

    try {
      const flightNumbersArray = flightNumbers
        .split(',')
        .map(fn => fn.trim())
        .filter(fn => fn.length > 0);

      if (flightNumbersArray.length === 0) {
        setError('Please enter at least one flight number');
        setIsSubmitting(false);
        return;
      }

      const request: CreateFlightRequest = {
        flight_numbers: flightNumbersArray,
        dep_date: deptDate,
      };

      await flightService.createFlight(request);
      setSuccess(`Successfully created ${flightNumbersArray.length} flight(s)!`);
      setFlightNumbers('');
      setDeptDate('');

      if (onFlightCreated) {
        onFlightCreated();
      }
    } catch (err) {
      if (err instanceof ApiError) {
        setError(err.message);
      } else {
        setError('Failed to create flight(s). Please try again.');
      }
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <div className="w-full max-w-md p-6 bg-white rounded-lg shadow-md">
      <h2 className="text-2xl font-bold text-gray-800 mb-6">Create New Flight(s)</h2>

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
          <label htmlFor="flightNumbers" className="block text-sm font-medium text-gray-700 mb-1">
            Flight Numbers
          </label>
          <input
            type="text"
            id="flightNumbers"
            value={flightNumbers}
            onChange={(e) => setFlightNumbers(e.target.value)}
            className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
            placeholder="e.g., XX1, XX2, XX3"
            required
          />
          <p className="text-xs text-gray-500 mt-1">
            Separate multiple flight numbers with commas
          </p>
        </div>

        <div>
          <label htmlFor="deptDate" className="block text-sm font-medium text-gray-700 mb-1">
            Departure Date
          </label>
          <input
            type="date"
            id="deptDate"
            value={deptDate}
            onChange={(e) => setDeptDate(e.target.value)}
            className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
            required
          />
        </div>

        <Button
          type="submit"
          disabled={isSubmitting}
          className="w-full py-2"
        >
          {isSubmitting ? 'Creating...' : 'Create Flight(s)'}
        </Button>
      </form>
    </div>
  );
}
