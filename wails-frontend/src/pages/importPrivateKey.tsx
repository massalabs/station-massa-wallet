import { useState } from 'react';
import { promptRequest } from '../events/events';
import { useLocation, useNavigate } from 'react-router-dom';
import { handleCancel } from '../utils/utils';
import { Password, Input, Button } from '@massalabs/react-ui-kit';

const ImportPrivatekey = () => {
  const navigate = useNavigate();

  const [input, setInput] = useState<string | undefined>(undefined);
  const [privateKey, setPkey] = useState<string | undefined>(undefined);

  let { state } = useLocation();
  const req: promptRequest = state.req;

  const applyStr = 'Next';

  const baselineStr = () => {
    if (privateKey) {
      return 'Choose an account name';
    }
    return 'Enter your Private Key';
  };

  const placeholder = () => {
    if (privateKey) {
      return 'Account name';
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

    navigate('/new-password', { state });
  };

  const updateInput = (e: React.ChangeEvent<HTMLInputElement>) => {
    setInput(e.target.value);
  };

  return (
    <div className="bg-primary flex flex-col justify-center items-center h-screen w-full">
      <div className="w-1/4  max-w-sm  min-w-fit">
        <div>
          <p className="mas-title text-neutral pb-4">{req.Msg}</p>
        </div>
        <div>
          <p className="mas-body text-neutral pb-4">{baselineStr()}</p>
        </div>
        {privateKey ? (
          <div className="pb-4">
            <Input
              id="input"
              onChange={updateInput}
              autoComplete="off"
              name="input"
              placeholder={placeholder()}
              value={input}
            />
          </div>
        ) : (
          <div className="pb-4">
            <Password
              id="input"
              onChange={updateInput}
              autoComplete="off"
              name="input"
              placeholder={placeholder()}
              value={input}
            />
          </div>
        )}
        <div className="flex flex-row gap-4 w-full pb-4">
          <div className="min-w-fit w-full">
            <Button variant={'secondary'} onClick={handleCancel}>
              Cancel
            </Button>
          </div>
          <div className="min-w-fit w-full">
            <Button onClick={handleApply}>{applyStr}</Button>
          </div>
        </div>
      </div>
    </div>
  );
};

export default ImportPrivatekey;
