import { SyntheticEvent, useRef, useState } from 'react';

import { Button, Password } from '@massalabs/react-ui-kit';
import { walletapp } from '@wailsjs/go/models';
import { SendPromptInput } from '@wailsjs/go/walletapp/WalletApp';
import { EventsOnce } from '@wailsjs/runtime/runtime';
import { FiTrash2 } from 'react-icons/fi';
import { useLocation, useNavigate } from 'react-router-dom';

import { promptRequest, promptResult } from '@/events';
import Intl from '@/i18n/i18n';
import { Layout } from '@/layouts/Layout/Layout';
import { PromptRequestDeleteData } from '@/pages/PasswordPromptHandler/PasswordPromptHandler';
import { validate } from '@/pages/PasswordPromptHandler/Sign';
import {
  ErrorCode,
  IErrorObject,
  handleApplyResult,
  handleCancel,
  parseForm,
} from '@/utils';

export function Delete() {
  const [error, setError] = useState<IErrorObject | null>(null);
  const [errorMessage, setErrorMessage] = useState<string>('');

  const navigate = useNavigate();
  const form = useRef(null);
  const { state } = useLocation();
  const req: promptRequest = state.req;
  const data: PromptRequestDeleteData = req.Data;

  const { EventType } = walletapp;

  function save(e: SyntheticEvent) {
    const form = parseForm(e);
    const { password } = form;

    EventsOnce(EventType.promptResult, handleResult);

    SendPromptInput(password);
  }

  function handleResult(result: promptResult) {
    let { Success, CodeMessage } = result;
    setError(null);
    setErrorMessage('');

    if (!Success) {
      if (CodeMessage === ErrorCode.WrongPassword) {
        setError({ password: Intl.t(`errors.${CodeMessage}`) });
      } else {
        setErrorMessage(Intl.t(`errors.delete`));
      }
      return;
    }
    handleApplyResult(navigate, req, setError, false)(result);
  }

  async function handleSubmit(e: SyntheticEvent) {
    e.preventDefault();
    if (!validate(e, setError)) return;

    const form = parseForm(e);
    const { password } = form;

    if (data.Balance !== '0') {
      navigate('/confirm-delete', { state: { req, password } });
    } else {
      save(e);
    }
  }

  return (
    <Layout>
      <form ref={form} onSubmit={handleSubmit}>
        <h1 className="mas-title">{Intl.t('password-prompt.title.delete')}</h1>
        <div className="mas-body pt-4 break-words">
          {Intl.t('password-prompt.subtitle.delete')}
        </div>
        <div className="pt-4">
          <Password
            defaultValue=""
            name="password"
            placeholder="Password"
            error={error?.password}
          />
          {errorMessage && (
            <p className="mt-2 text-s-error mas-body">{errorMessage}</p>
          )}
        </div>
        <div className="pt-4 flex gap-4">
          <Button variant="secondary" onClick={handleCancel}>
            {Intl.t('password-prompt.buttons.cancel')}
          </Button>
          <Button preIcon={<FiTrash2 />} type="submit">
            {Intl.t('password-prompt.buttons.delete')}
          </Button>
        </div>
      </form>
    </Layout>
  );
}
