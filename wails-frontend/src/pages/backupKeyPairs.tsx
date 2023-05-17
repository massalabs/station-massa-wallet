import { useState } from 'react';
import { promptRequest } from '../events/events';
import { useLocation } from 'react-router-dom';
import { handleCancel } from '../utils/utils';
import { Password, Button } from '@massalabs/react-ui-kit';
import { FiCopy } from 'react-icons/fi';

function BackupKeyPairs() {
  const [inputPassword, setInputPassword] = useState<string | undefined>(
    undefined,
  );
  const [show, setShow] = useState<boolean>(false);
  const [privateKey, setPkey] = useState<string | undefined>(undefined);

  let { state } = useLocation();
  const req: promptRequest = state.req;

  const baselineStr = 'Enter your password to show your private key';

  function handleCopy() {
    console.log('copy');
  }

  function handleShow() {
    if (!inputPassword || !inputPassword.length) {
      return;
    }

    if (!privateKey) {
      // TODO: validate private key
      setPkey(inputPassword);
      setShow(!show);
      setInputPassword('');
      return;
    }
  }

  function updateInput(e: React.ChangeEvent<HTMLInputElement>) {
    setInputPassword(e.target.value);
  }

  function enterPassword() {
    return (
      <>
        <div>
          <p className="mas-title text-neutral pb-4">{req.Msg}</p>
        </div>
        <div>
          <p className="mas-body text-neutral pb-4">{baselineStr}</p>
        </div>
        <div className="pb-4">
          <Password
            id="input"
            onChange={updateInput}
            autoComplete="off"
            name="input"
            placeholder={'Password'}
            value={inputPassword}
          />
        </div>
        <div className="flex flex-row gap-4 w-full pb-4">
          <div className="min-w-fit w-full">
            <Button variant={'secondary'} onClick={handleCancel}>
              Cancel
            </Button>
          </div>
          <div className="min-w-fit w-full">
            <Button onClick={handleShow}>Show</Button>
          </div>
        </div>
      </>
    );
  }

  function copyPrivateKey() {
    return (
      <>
        <div>
          <div>
            <p className="mas-title text-neutral pb-4">{req.Msg}</p>
          </div>
        </div>
        <div className="pb-4">
          <Password
            id="input"
            onChange={updateInput}
            autoComplete="off"
            name="input"
            placeholder={''}
            value={inputPassword}
          />
        </div>
        <div className="flex flex-row">
          <div className="pb-4 w-full">
            <Button
              onClick={handleCopy}
              preIcon={<FiCopy className="w-6 h-6" />}
            >
              Copy to clipboard
            </Button>
          </div>
        </div>
      </>
    );
  }

  return (
    <div className="bg-primary flex flex-col justify-center items-center h-screen w-full">
      <div className="w-1/4 max-w-sm min-w-fit">
        {!show ? enterPassword() : copyPrivateKey()}
      </div>
    </div>
  );
}

export default BackupKeyPairs;
