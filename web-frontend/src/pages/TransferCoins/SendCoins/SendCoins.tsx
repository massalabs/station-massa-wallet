import { useState } from 'react';

export function SendCoins() {
  const [amount, setAmount] = useState<number>(64);
  const [fees, setFees] = useState<number>(42);
  const [recepient, setRecepient] = useState<number>(42);
  return (
    <div className="bg-primary w-full h-screen">
      <h1>Send coins </h1>

      <form className="flex flex-col w-1/2">
        <input type="text" value={amount} placeholder="amount"></input>
        <input type="text" value={fees} placeholder="fees"></input>
        <input
          type="text"
          value={recepient}
          placeholder="recepient address"
        ></input>
      </form>
    </div>
  );
}
