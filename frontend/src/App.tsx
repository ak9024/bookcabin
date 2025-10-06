import { useState } from "react"
import { VoucherAssignment } from "@/components/VoucherAssignment"
import { VoucherList } from "@/components/VoucherList"
import { CreateVoucher } from "@/components/CreateVoucher"
import { CreateFlight } from "@/components/CreateFlight"
import { CreateSeat } from "@/components/CreateSeat"

function App() {
  const [refreshTrigger, setRefreshTrigger] = useState(0);

  const handleVoucherCreated = () => {
    setRefreshTrigger(prev => prev + 1);
  };

  const handleFlightCreated = () => {
    setRefreshTrigger(prev => prev + 1);
  };

  const handleSeatCreated = () => {
    setRefreshTrigger(prev => prev + 1);
  };

  return (
    <div className="flex min-h-svh flex-col items-center justify-center gap-8 bg-gray-50 p-8">
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-8 w-full max-w-7xl">
        <CreateFlight onFlightCreated={handleFlightCreated} />
        <CreateSeat onSeatCreated={handleSeatCreated} />
        <CreateVoucher onVoucherCreated={handleVoucherCreated} />
      </div>
      <VoucherAssignment />
      <VoucherList key={refreshTrigger} />
    </div>
  )
}

export default App
