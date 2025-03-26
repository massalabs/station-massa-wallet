import { useState, useRef, SyntheticEvent } from 'react';

import { Password, Button, Stepper } from '@massalabs/react-ui-kit';
import { walletapp } from '@wailsjs/go/models';
import {
  SendPKeyPromptInput,
  SendPromptInput,
} from '@wailsjs/go/walletapp/WalletApp';
import { EventsOnce } from '@wailsjs/runtime';
import { FiLock } from 'react-icons/fi';
import { useLocation, useNavigate } from 'react-router-dom';

import { promptRequest } from '@/events';
import { Layout } from '@/layouts/Layout/Layout';
import {
  parseForm,
  handleApplyResult,
  handleCancel,
  IErrorObject,
} from '@/utils';
import { hasMoreThanFiveChars, hasSamePassword } from '@/validation/password';

function NewPassword() {
  const navigate = useNavigate();
  const form = useRef(null);

  const { state } = useLocation();
  const req: promptRequest = state.req;

  const { EventType } = walletapp;

  const { newPassword, import: importReq } = walletapp.PromptRequestAction;
  const isImportAction = req.Action === importReq;

  function getButtonLabel() {
    switch (req.Action) {
      case newPassword:
        return 'Define';
      case importReq:
        return 'Define and import';
      default:
        return 'Apply';
    }
  }

  function getSubtitle() {
    switch (req.Action) {
      case newPassword:
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

    // reset error state
    setError(null);
    setErrorConfirm(null);

    if (!hasMoreThanFiveChars(password)) {
      setError({ password: 'Password must have at least 5 characters' });
      valid = false;
    }

    if (!hasSamePassword(password, passwordConfirm)) {
      setErrorConfirm({ password: "Password doesn't match" });
      valid = false;
    }

    return valid;
  }

  function save(e: SyntheticEvent) {
    const form = parseForm(e);
    const { password } = form;

    EventsOnce(
      EventType.promptResult,
      handleApplyResult(navigate, req, setError, isImportAction),
    );

    return isImportAction
      ? SendPKeyPromptInput(state.privateKey, state.nickname, password)
      : SendPromptInput(password);
  }

  async function handleSubmit(e: SyntheticEvent) {
    e.preventDefault();
    if (!validate(e)) return;

    save(e);
  }

  return (
    <Layout>
      {isImportAction && (
        <div className="pb-6">
          <Stepper
            step={2}
            steps={['Private Key', 'Account name', 'Password']}
          />
        </div>
      )}
      <form ref={form} onSubmit={handleSubmit}>
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
        <div className="pt-4">
          <Password
            defaultValue=""
            name="passwordConfirm"
            placeholder="Confirm your password"
            error={errorConfirm?.password}
          />
        </div>
        <div className="pt-4 flex gap-4">
          <div className="max-w-min">
            <Button variant="secondary" onClick={handleCancel}>
              Cancel
            </Button>
          </div>
          <Button preIcon={<FiLock />} type="submit">
            {getButtonLabel()}
          </Button>
        </div>
      </form>
    </Layout>
  );
}

export default NewPassword;
