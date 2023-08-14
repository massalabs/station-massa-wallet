import { FiCheck } from 'react-icons/fi';

export function ImportSuccess({ ...props }) {
  const { closeModal } = props;
  setTimeout(() => closeModal(), 100000);
  return (
    <div className="w-full h-full flex flex-col justify-center items-center">
      <div className="w-12 h-12 bg-brand flex flex-col justify-center items-center rounded-full mb-6">
        <FiCheck className="w-6 h-6 text-neutral" />
      </div>
      <div>Fungible Token [symbol] has been added</div>
    </div>
  );
}
