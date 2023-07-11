import { useState } from 'react';
import {
  SendPromptInput,
  SelectAccountFile,
} from '@wailsjs/go/walletapp/WalletApp';
import { events, promptRequest } from '@/events/events';
import { useLocation, useNavigate } from 'react-router-dom';
import { handleApplyResult, handleCancel } from '@/utils';
import { walletapp } from '@wailsjs/go/models';
import { EventsOnce } from '@wailsjs/runtime/runtime';
import { Button } from '@massalabs/react-ui-kit';
import { Layout } from '@/layouts/Layout/Layout';
import Intl from '@/i18n/i18n';

const ImportFile = () => {
  const nav = useNavigate();

  const [account, setAccount] = useState<
    undefined | walletapp.selectFileResult
  >(undefined);
  const [errorMsg, setErrorMsg] = useState('');

  const { state } = useLocation();
  const req: promptRequest = state.req;

  const accountStyleSuccess = 'mas-body text-s-success pb-4';
  const accountStyleNormal = 'mas-body text-neutral pb-4';

  const baselineStr = () =>
    account
      ? `Selected ${account.nickname}'s account`
      : 'Select an account file to import';

  const getImportLabel = () => (account ? 'Import' : 'Select a file');

  const handleApply = async () => {
    setErrorMsg('');
    if (!account) {
      const res = await SelectAccountFile();
      if (!Object.keys(res).length) {
        return;
      }
      if (res.err) {
        setErrorMsg(Intl.t(`errors.${res.codeMessage}`));
        return;
      }
      setAccount(res);
    } else {
      EventsOnce(
        events.promptResult,
        handleApplyResult(nav, req, setErrorMsg, true),
      );
      await SendPromptInput(account.filePath);
    }
  };

  return (
    <Layout>
      <h1 className="mas-title pb-4">{req.Msg}</h1>
      <p className={account ? accountStyleSuccess : accountStyleNormal}>
        {baselineStr()}
      </p>
      <div className="flex flex-row gap-4 pb-4">
        <Button variant={'secondary'} onClick={handleCancel}>
          Cancel
        </Button>
        <div className="min-w-fit">
          <Button onClick={handleApply}>{getImportLabel()}</Button>
        </div>
      </div>
      {errorMsg && <p className="mas-body text-s-error">{errorMsg}</p>}
    </Layout>
  );
};

export default ImportFile;
