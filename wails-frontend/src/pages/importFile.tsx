import { useState } from 'preact/hooks';
import {
  ImportWalletFile,
  SelectAccountFile,
} from '../../wailsjs/go/walletapp/WalletApp';
import { h } from 'preact';
import { events, promptRequest } from '../events/events';
import { useLocation, useNavigate } from 'react-router-dom';
import { handleApplyResult, handleCancel } from '../utils/utils';
import { walletapp } from '../../wailsjs/go/models';
import { EventsOnce } from '../../wailsjs/runtime/runtime';

const ImportFile = () => {
  const nav = useNavigate();

  const [account, setAccount] = useState<
    undefined | walletapp.selectFileResult
  >(undefined);
  const [errorMsg, setErrorMsg] = useState('');

  const { state } = useLocation();
  const req: promptRequest = state.req;

  const title = 'Import';
  const successMsg = 'The account has been imported';
  const baselineStr = () =>
    account
      ? `Importing account ${account.nickname}`
      : 'Select an accound file to import';
  const applyStr = () => (account ? 'Import' : 'Select a file');

  const handleApply = async () => {
    if (!account) {
      const res = await SelectAccountFile();
      console.log('res', res);
      if (!Object.keys(res).length) {
        return;
      }
      if (res.err) {
        setErrorMsg(res.err);
        return;
      }
      console.log('file', res.filePath);
      setAccount(res);
    } else {
      EventsOnce(
        events.passwordResult,
        handleApplyResult(nav, successMsg, setErrorMsg),
      );
      await ImportWalletFile(account.filePath);
    }
  };

  return (
    <section class="PasswordPrompt">
      <div>{title}</div>
      <div className="baseline">{baselineStr()}</div>
      <div className="flex">
        <button className="btn" onClick={handleCancel}>
          Cancel
        </button>
        <button className="btn" onClick={handleApply}>
          {applyStr()}
        </button>
      </div>
      <div className="baseline">{errorMsg}</div>
    </section>
  );
};

export default ImportFile;
