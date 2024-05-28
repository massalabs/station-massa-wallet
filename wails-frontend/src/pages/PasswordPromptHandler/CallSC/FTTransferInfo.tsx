import { Args } from '@massalabs/massa-web3';
import { formatFTAmount } from '@massalabs/react-ui-kit';
import { maskAddress } from '@massalabs/react-ui-kit/src/lib/massa-react/utils';
import { LogPrint } from '@wailsjs/runtime/runtime';

import Intl from '@/i18n/i18n';
import { AmountBox } from '@/pages/PasswordPromptHandler/AmountBox';
import { AssetInfo } from '@/pages/PasswordPromptHandler/Sign';
import { base64ToArray } from '@/utils/parameters';

interface FTTransferInfoProps {
  targetFunction: string;
  asset?: AssetInfo;
  parameters?: string;
}

export function FTTransferInfo(props: FTTransferInfoProps) {
  const { targetFunction, asset, parameters } = props;

  if (targetFunction !== 'transfer' || !asset || !parameters) {
    return null;
  }

  let amount = 0n;
  let recipient = '';
  try {
    const args = new Args(base64ToArray(parameters));
    recipient = args.nextString();
    amount = args.nextU256();
  } catch (error) {
    LogPrint(`error FTTransferInfo: ${error}`);
  }

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
