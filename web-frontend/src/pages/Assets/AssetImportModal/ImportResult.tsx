import { FiCheck, FiX } from 'react-icons/fi';

import Intl from '@/i18n/i18n';

export function ImportResult({ ...props }) {
  const { closeModal, token, importResult, inputError } = props;
  setTimeout(() => closeModal(), 2000);

  return (
    <div className="w-full h-full flex flex-col justify-center items-center px-6">
      <div
        className={`w-12 h-12 flex flex-col justify-center items-center rounded-full mb-6 ${
          importResult ? 'bg-brand' : 'bg-s-error'
        }`}
      >
        {importResult ? (
          <FiCheck className="w-6 h-6 text-neutral" />
        ) : (
          <FiX className="w-6 h-6  text-neutral" />
        )}
      </div>
      {importResult ? (
        <div className="text-center">
          {Intl.t('assets.import.success-screen.success-message', {
            name: token.name,
            symbol: token.symbol,
          })}
          <br />
          {Intl.t('assets.import.success-screen.success-address', {
            name: token?.assetAddress,
          })}
        </div>
      ) : (
        <div className="text-center">{inputError.address}</div>
      )}
    </div>
  );
}
