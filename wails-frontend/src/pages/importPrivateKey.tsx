import { useState } from 'preact/hooks';
import { ImportPrivateKey } from '../../wailsjs/go/walletapp/WalletApp';
import { EventsOnce } from '../../wailsjs/runtime/runtime';
import { h } from 'preact';
import { events, promptRequest } from '../events/events';
import { useLocation, useNavigate } from 'react-router-dom';
import { handleCancel, hideAndReload } from '../utils/utils';

const ImportPrivatekey = () => {
  const navigate = useNavigate();

  const [pkey, setPkey] = useState('');
  const [resultMsg, setResultMsg] = useState('');

  const { state } = useLocation();
  const req: promptRequest = state.req;

  const applyStr = 'Import';

  const title = 'Import';
  const baselineStr = 'Choose an import method';

  const handleApplyResult = (result: any) => {
    if (result.Success) {
      navigate('/success', {
        state: { msg: 'Account imported successfully!' },
      });
      setTimeout(hideAndReload, 2000);
    } else {
      setResultMsg(result.Data);
      if (result.Error === 'timeoutError') {
        setTimeout(hideAndReload, 2000);
      }
    }
  };

  const handleApply = async () => {
    if (pkey.length) {
      EventsOnce(events.passwordResult, handleApplyResult);
      await ImportPrivateKey(pkey);
    }
  };

  const updatePkey = (e: any) => setPkey(e.target.value);

  return (
    <section class="PasswordPrompt">
      <div>{title}</div>
      <div className="baseline">{baselineStr}</div>
      <div id="input" className="input-box">
        <input
          id="name"
          className="input"
          onChange={updatePkey}
          autoComplete="off"
          name="input"
          placeholder="Enter your private key"
        />
      </div>
      <div className="flex">
        <button className="btn" onClick={handleCancel}>
          Cancel
        </button>
        <button className="btn" onClick={handleApply}>
          {applyStr}
        </button>
      </div>
      <div className="baseline">{resultMsg}</div>
    </section>
  );
};

export default ImportPrivatekey;
