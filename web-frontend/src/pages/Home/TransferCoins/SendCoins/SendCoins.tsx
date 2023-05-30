import { useState } from 'react';

export function SendCoins() {
  const [amount, setAmount] = useState<number>(64);
  const [fees, setFees] = useState<number>(42);
  const [recipient, setRecepient] = useState<string>('joe');

  function handleAmountChange() {
    console.log('handle amount change');
  }

  function handleFeesChange() {
    console.log('handle fee change');
  }

  function handleRecipientChange() {
    console.log('handle recipient change');
  }

  function handleSubmit() {
    console.log('submit');
  }

  return (
    <div className="bg-primary w-full h-screen">
      <h1>Send coins </h1>

      <form className="flex flex-col w-1/2" onSubmit={handleSubmit}>
        <div className="mb-4">
          <label htmlFor="amount" className="block mb-2">
            Amount:
          </label>
          <input
            type="text"
            id="amount"
            placeholder="Enter amount"
            onChange={handleAmountChange}
            className="border border-gray-300 p-2"
          />
        </div>

        <div className="mb-4">
          <label htmlFor="fees" className="block mb-2">
            Fees:
          </label>
          <input
            type="text"
            id="fees"
            placeholder="Enter fees"
            onChange={handleFeesChange}
            className="border border-gray-300 p-2"
          />
        </div>

        <div className="mb-4">
          <label htmlFor="recipient" className="block mb-2">
            Recipient Address:
          </label>
          <input
            type="text"
            id="recipient"
            placeholder="Enter recipient address"
            onChange={handleRecipientChange}
            className="border border-gray-300 p-2"
          />
        </div>

        <button
          type="submit"
          className="bg-blue-500 text-white py-2 px-4 rounded hover:bg-blue-600"
          onClick={handleSubmit}
        >
          Submit
        </button>
      </form>
    </div>
  );
}
