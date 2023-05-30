import { useState, useRef, SyntheticEvent } from 'react';
import { useLocation, useNavigate } from 'react-router-dom';
import { ClipboardSetText, EventsOnce } from '../../wailsjs/runtime';
import { events, promptRequest, promptResult } from '../events/events';
import { parseForm } from '../utils/parseForm';
import { SendPromptInput } from '../../wailsjs/go/walletapp/WalletApp';
import { ErrorCode, IErrorObject, getErrorMessage } from '../utils';

import { Password, Button } from '@massalabs/react-ui-kit';
import { Layout } from '../layouts/Layout/Layout';
import { FiCopy, FiArrowRight } from 'react-icons/fi';

function EnterKey() {
  const navigate = useNavigate();

  let { state } = useLocation();
  const req: promptRequest = state.req;

  return (
    <>
      <Button
        variant="secondary"
        onClick={() => navigate('/backup-methods', { state: { req } })}
      >
        Cancel
      </Button>
      <Button posIcon={<FiArrowRight />} type="submit">
        Next
      </Button>
    </>
  );
}

interface CopyProps {
  privateKey: string;
}

function CopyKey({ privateKey }: CopyProps) {
  async function handleCopy() {
    await ClipboardSetText(privateKey);
  }

  return (
    <Button onClick={handleCopy} preIcon={<FiCopy />}>
      Copy to clipboard
    </Button>
  );
}

function BackupKeyPairs() {
  const form = useRef(null);
  const { state } = useLocation();
  const navigate = useNavigate();

  const req: promptRequest = state.req;

  const [privateKey, setPrivateKey] = useState<string | undefined>(undefined);
  const [error, setError] = useState<IErrorObject | null>(null);

  function validate(e: SyntheticEvent) {
    const formObject = parseForm(e);
    const { password } = formObject;

    if (!password || !password.length) {
      setError({ password: 'Password is required' });
      return false;
    }
    return true;
  }

  async function handleResult(result: promptResult) {
    let { Success, CodeMessage } = result;

    if (!Success) {
      if (CodeMessage === ErrorCode.WrongPassword) {
        setError({ password: getErrorMessage(CodeMessage) });
        return;
      } else {
        req.Msg = getErrorMessage(CodeMessage);
        navigate('/failure', {
          state: { req },
        });
      }
    }
    setPrivateKey(result.Data);
  }

  async function handleSubmit(e: SyntheticEvent) {
    e.preventDefault();
    if (!validate(e)) return;

    const formObject = parseForm(e);
    const { password } = formObject;

    EventsOnce(events.promptResult, handleResult);
    await SendPromptInput(password);
  }

  return (
    <Layout>
      <form ref={form} onSubmit={handleSubmit}>
        <h1 className="mas-title">{req.Msg}</h1>
        {!privateKey && (
          <p className="mas-body pt-4">
            Enter your password to show your private key
          </p>
        )}
        <div className="pt-4">
          {privateKey ? (
            <Password
              defaultValue=""
              name="privateKey"
              placeholder="Private key"
              value={privateKey}
            />
          ) : (
            <Password
              defaultValue=""
              name="password"
              placeholder="Password"
              error={error?.password}
            />
          )}
        </div>
        <div className="flex flex-row gap-4 pt-4">
          {privateKey ? <CopyKey privateKey={privateKey} /> : <EnterKey />}
        </div>
      </form>
    </Layout>
  );
}

export default BackupKeyPairs;
