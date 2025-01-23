import {
  getNetworkNameByChainId,
  NetworkName as networkNameEnum,
} from '@massalabs/massa-web3';
import { InterrogationPoint, MassaLogo } from '@massalabs/react-ui-kit';

import Intl from '@/i18n/i18n';

export interface NetworkNameProps {
  chainId: number;
}

export function NetworkName({ chainId }: NetworkNameProps) {
  let secondaryColor = undefined;
  let primaryColor = undefined;

  let networkName: undefined | string | networkNameEnum =
    getNetworkNameByChainId(BigInt(chainId));
  if (networkName) {
    if (networkName === networkNameEnum.Mainnet) {
      primaryColor = '#FF0000';
      secondaryColor = '#FFFFFF';
    } else {
      primaryColor = '#FFFFFF';
      secondaryColor = '#151A26';
    }
  } else {
    networkName = `${Intl.t(
      'password-prompt.sign.unknown-network-name',
    )} (${chainId})`;
  }

  return (
    <div
      data-testid="tag"
      className="flex justify-between items-center bg-tertiary mas-caption
        rounded-full w-fit px-3 py-1 text-f-primary mb-4"
    >
      {networkName ? (
        <MassaLogo
          size={16}
          primaryColor={primaryColor}
          secondaryColor={secondaryColor}
        />
      ) : (
        <InterrogationPoint size={16} />
      )}
      <span className="ml-2">{networkName}</span>
    </div>
  );
}
