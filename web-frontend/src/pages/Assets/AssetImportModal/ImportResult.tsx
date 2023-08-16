import { FiCheck, FiX } from 'react-icons/fi';

import { InputsErrors } from '@/const/assets/assets';
import Intl from '@/i18n/i18n';
import { Asset } from '@/models/AssetModel';
import { maskAddress } from '@/utils';

interface ImportResultProps {
  closeModal: () => void;
  data: Asset;
  importResult: boolean;
  inputError: InputsErrors | null;
}

export function ImportResult(props: ImportResultProps) {
  const { closeModal, data, importResult, inputError } = props;
  setTimeout(() => closeModal(), 2000);

  const bgColor = importResult ? 'bg-brand' : 'bg-s-error';

  return (
    <div className="w-full h-full flex flex-col justify-center items-center px-6">
      <div
        className={`w-12 h-12 flex flex-col justify-center items-center rounded-full mb-6 ${bgColor} `}
      >
        {importResult ? (
          <FiCheck className="w-6 h-6 text-neutral" />
        ) : (
          <FiX className="w-6 h-6 text-neutral" />
        )}
      </div>
      {importResult ? (
        <div className="text-center">
          {Intl.t('assets.import.success-screen.success-message', {
            name: data.name,
            symbol: data.symbol,
          })}
          <br />
          {Intl.t('assets.import.success-screen.success-address', {
            address: maskAddress(data?.address),
          })}
        </div>
      ) : (
        <div className="text-center">{inputError && inputError.address}</div>
      )}
    </div>
  );
}
