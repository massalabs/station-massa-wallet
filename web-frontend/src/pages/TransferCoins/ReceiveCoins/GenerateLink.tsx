import { useState } from 'react';
import { formatStandard } from '../../../utils/MassaFormating';
import CopyContent from './CopyContent';
import { AccountObject } from '../../../models/AccountModel';
import Intl from '../../../i18n/i18n';
import {
  Button,
  Identicon,
  Input,
  MassaLogo,
  PopupModal,
  PopupModalContent,
  PopupModalHeader,
  Selector,
} from '@massalabs/react-ui-kit';

interface GenerateLinkProps {
  account: AccountObject;
  presetURL: string;
  url: string;
  setURL: (url: string) => void;
  setModal: (modal: boolean) => void;
}
function GenerateLink(props: GenerateLinkProps) {
  const { account, presetURL, url, setURL, setModal } = props;

  const [amount, setAmount] = useState('');
  const [provider, setProvider] = useState('');
  const recipient = account.nickname;
  const recipientBalance = parseInt(account.candidateBalance) / 10 ** 9;
  const formattedBalance = formatStandard(recipientBalance);

  const handleGenerate = () => {
    const amountArg = amount ? `&amount=${amount}` : '';
    const providerArg = provider ? `&provider=${provider}` : '';
    const newURL = presetURL + amountArg + providerArg;
    setURL(newURL);
  };

  return (
    <>
      <PopupModal
        customClass="!w-1/2 min-w-[775px]"
        fullMode={true}
        onClose={() => setModal(false)}
      >
        <PopupModalHeader>
          <h1 className="mas-banner mb-6">
            {Intl.t('receive.receive-account')}
          </h1>
        </PopupModalHeader>
        <PopupModalContent>
          <div className="pb-10">
            <div className="flex flex-col gap-3 mb-6">
              <p className="mas-body2">{Intl.t('receive.receive-amount')}</p>
              <Input
                placeholder={Intl.t('receive.amount-to-ask')}
                defaultValue=""
                onChange={(e) => setAmount(e.target.value)}
              />
            </div>
            <div className="flex flex-col gap-3 mb-6">
              <p className="mas-body2">{Intl.t('receive.recipient')}</p>
              <Selector
                preIcon={<Identicon username={account.nickname} />}
                content={recipient}
                amount={formattedBalance}
                posIcon={<MassaLogo />}
                variant="secondary"
              />
            </div>
            <div className="flex flex-col gap-3 mb-6">
              <p className="mas-body2">{Intl.t('receive.provider')}</p>
              <Input
                placeholder={Intl.t('receive.provider-description')}
                defaultValue=""
                onChange={(e) => setProvider(e.target.value)}
              />
            </div>
            <div className="flex flex-col gap-3 mb-3">
              <p className="mas-body2">{Intl.t('receive.link-to-share')}</p>
              <CopyContent
                content={url}
                formattedContent={url.slice(0, 50) + '...'}
              />
            </div>
            <div className="pb-3">
              <Button onClick={() => handleGenerate()}>
                {Intl.t('receive.receive-account')}
              </Button>
            </div>
          </div>
        </PopupModalContent>
      </PopupModal>
    </>
  );
}

export default GenerateLink;
