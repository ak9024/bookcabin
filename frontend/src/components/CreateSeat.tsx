import { useState, useEffect } from 'react';
import { Button } from '@/components/ui/button';
import { seatService } from '@/services/seat.service';
import { flightService } from '@/services/flight.service';
import { ApiError } from '@/services/api';
import type { CreateSeatRequest } from '@/types/seat.types';
import type { Flight } from '@/types/flight.types';

interface CreateSeatProps {
  onSeatCreated?: () => void;
}

export function CreateSeat({ onSeatCreated }: CreateSeatProps) {
  const [flights, setFlights] = useState<Flight[]>([]);
  const [isLoadingFlights, setIsLoadingFlights] = useState(true);
  const [flightId, setFlightId] = useState<number>(0);
  const [cabin, setCabin] = useState<string>('');
  const [seatLabels, setSeatLabels] = useState<string>('');
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [success, setSuccess] = useState<string | null>(null);

  useEffect(() => {
    const fetchFlights = async () => {
      setIsLoadingFlights(true);
      try {
        const data = await flightService.getAllFlights();
        setFlights(data);
      } catch (err) {
        if (err instanceof ApiError) {
          setError(`Failed to load flights: ${err.message}`);
        } else {
          setError('Failed to load flights. Please refresh the page.');
        }
      } finally {
        setIsLoadingFlights(false);
      }
    };

    fetchFlights();
  }, []);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setIsSubmitting(true);
    setError(null);
    setSuccess(null);

    try {
      const labelsArray = seatLabels
        .split(',')
        .map(label => label.trim())
        .filter(label => label.length > 0);

      if (labelsArray.length === 0) {
        setError('Please enter at least one seat label');
        setIsSubmitting(false);
        return;
      }

      const request: CreateSeatRequest = {
        flight_id: flightId,
        cabin: cabin,
        labels: labelsArray,
      };

      await seatService.createSeats(request);
      setSuccess(`Successfully created ${labelsArray.length} seat(s)!`);
      setFlightId(0);
      setCabin('');
      setSeatLabels('');

      if (onSeatCreated) {
        onSeatCreated();
      }
    } catch (err) {
      if (err instanceof ApiError) {
        setError(err.message);
      } else {
        setError('Failed to create seat(s). Please try again.');
      }
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <div className="w-full max-w-md p-6 bg-white rounded-lg shadow-md">
      <h2 className="text-2xl font-bold text-gray-800 mb-6">Create New Seat(s)</h2>

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
          <label htmlFor="flightId" className="block text-sm font-medium text-gray-700 mb-1">
            Flight
          </label>
          <select
            id="flightId"
            value={flightId || ''}
            onChange={(e) => setFlightId(parseInt(e.target.value))}
            className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
            required
            disabled={isLoadingFlights}
          >
            <option value="">
              {isLoadingFlights ? 'Loading flights...' : 'Select a flight'}
            </option>
            {flights.map((flight) => (
              <option key={flight.id} value={flight.id}>
                {flight.flight_no} - {flight.dep_date}
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
            value={cabin}
            onChange={(e) => setCabin(e.target.value)}
            className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
            required
          >
            <option value="">Select a cabin</option>
            <option value="ECONOMY">Economy</option>
            <option value="BUSINESS">Business</option>
            <option value="FIRST">First</option>
          </select>
        </div>

        <div>
          <label htmlFor="seatLabels" className="block text-sm font-medium text-gray-700 mb-1">
            Seat Labels
          </label>
          <input
            type="text"
            id="seatLabels"
            value={seatLabels}
            onChange={(e) => setSeatLabels(e.target.value)}
            className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
            placeholder="e.g., 1A, 1B, 1C"
            required
          />
          <p className="text-xs text-gray-500 mt-1">
            Separate multiple seat labels with commas
          </p>
        </div>

        <Button
          type="submit"
          disabled={isSubmitting}
          className="w-full py-2"
        >
          {isSubmitting ? 'Creating...' : 'Create Seat(s)'}
        </Button>
      </form>
    </div>
  );
}
