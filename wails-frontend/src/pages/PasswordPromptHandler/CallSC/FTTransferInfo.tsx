import { Args } from '@massalabs/massa-web3';
import { formatFTAmount } from '@massalabs/react-ui-kit';

import Intl from '@/i18n/i18n';
import { AmountBox } from '@/pages/PasswordPromptHandler/AmountBox';
import { AssetInfo } from '@/pages/PasswordPromptHandler/Sign';
import { maskAddress } from '@/utils';
import { base64ToArray } from '@/utils/parameters';

interface FTTransferInfoProps {
  asset?: AssetInfo;
  parameters?: string;
}

export function FTTransferInfo(props: FTTransferInfoProps) {
  const { asset, parameters } = props;

  if (!asset || !parameters) {
    return null;
  }

  const args = new Args(base64ToArray(parameters));
  const recipient = args.nextString();
  const amount = args.nextU256();

  return (
    <>
      <AmountBox>
        {formatFTAmount(amount, asset.decimals).amountFormattedFull}{' '}
        {asset.symbol}
      </AmountBox>

      <div className="flex justify-between w-full">
        <div className="flex flex-col h-fit">
          <p>{Intl.t('password-prompt.sign.recipient')}</p>
        </div>
        <div className="flex flex-col items-end h-fit">
          <p>{maskAddress(recipient)}</p>
        </div>
      </div>
    </>
  );
}
