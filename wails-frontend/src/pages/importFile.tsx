import { useState } from 'react';
import {
  SendPromptInput,
  SelectAccountFile,
} from '../../wailsjs/go/walletapp/WalletApp';
import { events, promptRequest } from '../events/events';
import { useLocation, useNavigate } from 'react-router-dom';
import { handleApplyResult, handleCancel } from '../utils/utils';
import { walletapp } from '../../wailsjs/go/models';
import { EventsOnce } from '../../wailsjs/runtime/runtime';
import { Button } from '@massalabs/react-ui-kit';

const ImportFile = () => {
  const nav = useNavigate();

  const [account, setAccount] = useState<
    undefined | walletapp.selectFileResult
  >(undefined);
  const [errorMsg, setErrorMsg] = useState('');

  const { state } = useLocation();
  const req: promptRequest = state.req;

  const accountStyleSuccess = ' mas-body text-s-success  pb-4';
  const accountStyleNormal = ' mas-body text-neutral  pb-4';

  const baselineStr = () =>
    account
      ? `Selected ${account.nickname}'s account`
      : 'Select an account file to import';

  const importStr = () => (account ? 'Import' : 'Select a file');

  const handleApply = async () => {
    setErrorMsg('');
    if (!account) {
      const res = await SelectAccountFile();
      if (!Object.keys(res).length) {
        return;
      }
      if (res.err) {
        setErrorMsg(res.err);
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
    <div className="bg-primary flex flex-col justify-center items-center h-screen w-full">
      <div className="w-1/4  max-w-sm  min-w-fit">
        <div>
          <p className="mas-title text-neutral pb-4">{req.Msg}</p>
        </div>
        <div>
          <p className={account ? accountStyleSuccess : accountStyleNormal}>
            {baselineStr()}
          </p>
        </div>
        <div className="flex flex-row gap-4  pb-4">
          <div className="min-w-fit">
            <Button variant={'secondary'} onClick={handleCancel}>
              Cancel
            </Button>
          </div>
          <div className="min-w-fit">
            <Button onClick={handleApply}>{importStr()}</Button>
          </div>
        </div>
        <div>
          <p className="mas-body text-s-error">{errorMsg}</p>
        </div>
      </div>
    </div>
  );
};

export default ImportFile;
