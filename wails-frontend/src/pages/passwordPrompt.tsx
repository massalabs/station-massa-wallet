import { useState, useRef, SyntheticEvent } from 'react';
import { useLocation, useNavigate } from 'react-router-dom';
import {
  ImportPrivateKey,
  SendPromptInput,
} from '../../wailsjs/go/walletapp/WalletApp';
import { EventsOnce } from '../../wailsjs/runtime';
import { events, promptAction, promptRequest } from '../events/events';
import { handleApplyResult, handleCancel } from '../utils/utils';
import { parseForm } from '../utils/parseForm';
import { hasMoreThanFiveChars, hasSamePassword } from '../validation/password';

import { FiLock } from 'react-icons/fi';
import { Password, Button } from '@massalabs/react-ui-kit';
import { IErrorObject } from '../utils';

function PasswordPrompt() {
  const navigate = useNavigate();
  const form = useRef(null);

  const { state } = useLocation();
  const req: promptRequest = state.req;

  const { newPasswordReq, importReq, deleteReq } = promptAction;
  const isNewPasswordAction = req.Action === newPasswordReq;
  const isImportAction = req.Action === importReq;

  function getButtonLabel() {
    switch (req.Action) {
      case deleteReq:
        return 'Delete';
      case newPasswordReq:
        return 'Define';
      default:
        return 'Apply';
    }
  }

  function getSubtitle() {
    switch (req.Action) {
      case deleteReq:
        return 'Enter your password to delete your wallet';
      case newPasswordReq:
        return 'Enter a secure password';
      case importReq:
        return 'Define a new password';
      default:
        return 'Enter your password below to validate';
    }
  }

  // error state for the input password field
  const [error, setError] = useState<IErrorObject | null>(null);
  // error state for the input password confirm field
  const [errorConfirm, setErrorConfirm] = useState<IErrorObject | null>(null);

  function validate(e: SyntheticEvent) {
    let valid = true;
    const form = parseForm(e);
    const { password, passwordConfirm } = form;

    // reset error states
    setError(null);
    setErrorConfirm(null);

    if (!hasMoreThanFiveChars(password)) {
      setError({ password: 'Password must have at least 5 characters' });
      valid = false;
    }

    if (isNewPasswordAction && !hasSamePassword(password, passwordConfirm)) {
      setErrorConfirm({ password: "Password doesn't match" });
      valid = false;
    }

    return valid;
  }

  function save(e: SyntheticEvent) {
    const form = parseForm(e);
    const { password } = form;

    EventsOnce(
      events.promptResult,
      handleApplyResult(navigate, req, setError, isImportAction),
    );

    return isImportAction
      ? ImportPrivateKey(state.pkey, state.nickname, password)
      : SendPromptInput(password);
  }

  async function handleSubmit(e: SyntheticEvent) {
    e.preventDefault();
    if (!validate(e)) return;

    save(e);
  }

  return (
    <div className="bg-primary min-h-screen">
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
          {(isNewPasswordAction || isImportAction) && (
            <div className="pt-4">
              <Password
                defaultValue=""
                name="passwordConfirm"
                placeholder="Confirm your password"
                error={errorConfirm?.password}
              />
            </div>
          )}
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
    </div>
  );
}

export default PasswordPrompt;
