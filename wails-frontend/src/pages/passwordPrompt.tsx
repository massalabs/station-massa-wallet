import { useState, useRef, SyntheticEvent } from 'react';
import { useLocation, useNavigate } from 'react-router-dom';
import { SendPromptInput } from '../../wailsjs/go/walletapp/WalletApp';
import { EventsOnce } from '../../wailsjs/runtime';
import {
  events,
  promptAction,
  promptRequest,
  promptResult,
} from '../events/events';
import { handleApplyResult, handleCancel } from '../utils/utils';
import { parseForm } from '../utils/parseForm';
import { hasMoreThanFiveChars } from '../validation/password';

import { FiLock } from 'react-icons/fi';
import { Password, Button } from '@massalabs/react-ui-kit';
import { ErrorCode, IErrorObject, getErrorMessage } from '../utils';
import { Layout } from '../layouts/Layout/Layout';

function PasswordPrompt() {
  const navigate = useNavigate();
  const form = useRef(null);

  const { state } = useLocation();
  const req: promptRequest = state.req;

  const { deleteReq } = promptAction;

  function getButtonLabel() {
    switch (req.Action) {
      case deleteReq:
        return 'Delete';
      default:
        return 'Apply';
    }
  }

  function getSubtitle() {
    switch (req.Action) {
      case deleteReq:
        return 'Enter your password to delete your wallet';
      default:
        return 'Enter your password below to validate';
    }
  }

  // error state for the input password field
  const [error, setError] = useState<IErrorObject | null>(null);

  function validate(e: SyntheticEvent) {
    const form = parseForm(e);
    const { password } = form;

    // reset error state
    setError(null);

    if (!hasMoreThanFiveChars(password)) {
      setError({ password: 'Password must have at least 5 characters' });
      return false;
    }

    return true;
  }

  function save(e: SyntheticEvent) {
    const form = parseForm(e);
    const { password } = form;

    EventsOnce(events.promptResult, handleResult);

    return SendPromptInput(password);
  }

  async function handleResult(result: promptResult) {
    let { Success, CodeMessage } = result;

    if (!Success) {
      if (CodeMessage === ErrorCode.WrongPassword) {
        setError({ password: getErrorMessage(CodeMessage) });
      }
    }

    handleApplyResult(navigate, req, setError, false)(result);
  }

  async function handleSubmit(e: SyntheticEvent) {
    e.preventDefault();
    if (!validate(e)) return;

    save(e);
  }

  return (
    <Layout>
      <form ref={form} onSubmit={handleSubmit}>
        <div
          className="flex flex-col justify-center h-screen
          max-w-xs min-w-fit text-f-primary m-auto"
        >
          <h1 className="mas-title">{req.Msg}</h1>
          <p className="mas-body pt-4">{getSubtitle()}</p>
          <div className="pt-4">
            <Password
              defaultValue=""
              name="password"
              placeholder="Password"
              error={error?.password}
            />
          </div>
          <div className="pt-4 flex gap-4">
            <div className="max-w-min">
              <Button variant={'secondary'} onClick={handleCancel}>
                Cancel
              </Button>
            </div>
            <Button preIcon={<FiLock />} type="submit">
              {getButtonLabel()}
            </Button>
          </div>
        </div>
      </form>
    </Layout>
  );
}

export default PasswordPrompt;
