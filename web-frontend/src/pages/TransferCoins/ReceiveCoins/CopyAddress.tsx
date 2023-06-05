import { FiCopy } from 'react-icons/fi';

function CopyAddress({ ...props }) {
  const { address, formattedAddress } = props;

  function handleCopyClick() {
    navigator.clipboard.writeText(address);
  }

  return (
    <>
      <div
        className="flex flex-row items-center justify-between w-full 
        h-12 mb-3.5 px-3 rounded bg-secondary cursor-pointer"
        onClick={() => handleCopyClick()}
      >
        <u>{formattedAddress}</u>
        <FiCopy size={24} />
      </div>
    </>
  );
}

export default CopyAddress;
