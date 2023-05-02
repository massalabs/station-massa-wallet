import { useState } from 'preact/hooks';
import { h } from 'preact';
import { promptRequest } from '../events/events';
import { useLocation, useNavigate } from 'react-router-dom';
import { handleCancel } from '../utils/utils';

const ImportPrivatekey = () => {
  const navigate = useNavigate();

  const [input, setInput] = useState<string | undefined>(undefined);
  const [privateKey, setPkey] = useState<string | undefined>(undefined);

  let { state } = useLocation();
  const req: promptRequest = state.req;

  const applyStr = 'Next';

  const baselineStr = () => {
    if (privateKey) {
      return 'Choose a username';
    }
    return 'To import enter your private key';
  };

  const placeholder = () => {
    if (privateKey) {
      return 'Username';
    }
    return 'Private key';
  };

  const handleApply = () => {
    if (!input || !input.length) {
      return;
    }

    if (!privateKey) {
      // TODO: validate private key
      setPkey(input);
      setInput('');
      return;
    }

    state = { ...state, pkey: privateKey, nickname: input };

    navigate('/password', { state });
  };

  const updateInput = (e: any) => setInput(e.target.value);

  return (
    <section>
      <div>{req.Msg}</div>
      <div className="baseline">{baselineStr()}</div>
      <div id="input" className="input-box">
        <input
          id="input"
          className="input"
          onChange={updateInput}
          autoComplete="off"
          name="input"
          type={privateKey ? 'text' : 'password'}
          placeholder={placeholder()}
          value={input}
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
    </section>
  );
};

export default ImportPrivatekey;
