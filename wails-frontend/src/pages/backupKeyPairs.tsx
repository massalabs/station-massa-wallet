/* eslint-disable new-cap */
import { useState, useRef, SyntheticEvent } from 'react';
import { useLocation, useNavigate } from 'react-router-dom';
import { EventsOnce } from '../../wailsjs/runtime';
import { events, promptAction, promptRequest } from '../events/events';
import { parseForm } from '../utils/parseForm';
import { hasMoreThanFiveChars } from '../validation/password';
import { handleApplyResult } from '../utils/utils';
import { ApplyPassword } from '../../wailsjs/go/walletapp/WalletApp';

import { Layout } from '../layouts/Layout/Layout';
import { Password, Button } from '@massalabs/react-ui-kit';
import { FiCopy } from 'react-icons/fi';

interface IErrorObject {
  password: string;
}

function EnterKey() {
  const navigate = useNavigate();

  let { state } = useLocation();
  const req: promptRequest = state.req;

  return (
    <>
      <div className="min-w-fit w-full">
        <Button
          variant="secondary"
          onClick={() => navigate('/backup-methods', { state: { req } })}
        >
          Cancel
        </Button>
      </div>
      <div className="min-w-fit w-full">
        <Button type="submit">Next</Button>
      </div>
    </>
  );
}

function CopyKey() {
  function handleCopy() {
    // TODO: copy to clipboard
    console.log('Copy to clipboard');
  }

  return (
    <Button onClick={handleCopy} preIcon={<FiCopy />}>
      Copy to clipboard
    </Button>
  );
}

function BackupKeyPairs() {
  const navigate = useNavigate();
  const form = useRef(null);
  const { state } = useLocation();

  const req: promptRequest = state.req;
  const { backupReq } = promptAction;

  const isExportAction = req.Action === backupReq;

  const [show, setShow] = useState<boolean>(false);
  const [error, setError] = useState<IErrorObject | null>(null);

  function validate(e: SyntheticEvent) {
    const formObject = parseForm(e);
    const { password } = formObject;

    if (!hasMoreThanFiveChars(password)) {
      setError({ password: 'Password must have at least 5 characters' });
      return false;
    } else if (!password || !password.length) {
      setError({ password: 'Password is required' });
      return false;
    }
    return true;
  }

  function save(e: SyntheticEvent) {
    const formObject = parseForm(e);
    const { password } = formObject;

    EventsOnce(
      events.promptResult,
      handleApplyResult(navigate, req, setError, isExportAction),
    );

    return ApplyPassword(password);
  }

  async function handleSubmit(e: SyntheticEvent) {
    e.preventDefault();
    if (!validate(e)) return;

    save(e);
    // setShow(true); when backend is ready
    setShow(false);
  }

  return (
    <Layout>
      <form ref={form} onSubmit={handleSubmit}>
        <h1 className="mas-title text-neutral">{req.Msg}</h1>
        {!show && (
          <p className="mas-body text-neutral pt-4">
            Enter your password to show your private key
          </p>
        )}
        <div className="pt-4">
          <Password
            defaultValue=""
            name="password"
            placeholder={!show ? 'Password' : 'Private key'}
            error={error?.password}
          />
        </div>
        <div className="flex flex-row gap-4 pt-4">
          {!show ? <EnterKey /> : <CopyKey />}
        </div>
      </form>
    </Layout>
  );
}

export default BackupKeyPairs;
