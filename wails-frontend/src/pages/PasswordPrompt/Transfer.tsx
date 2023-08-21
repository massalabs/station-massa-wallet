import {
  ErrorCode,
  IErrorObject,
  formatStandard,
  handleApplyResult,
  handleCancel,
  maskAddress,
  parseForm,
} from '@/utils';
import { toMAS } from '@massalabs/massa-web3';
import { Balance, Button, Password } from '@massalabs/react-ui-kit';
import Intl from '@/i18n/i18n';
import { Layout } from '@/layouts/Layout/Layout';
import { SyntheticEvent, useRef, useState } from 'react';
import { FiAlertCircle } from 'react-icons/fi';
import { events, promptRequest, promptResult } from '@/events/events';
import { useLocation, useNavigate } from 'react-router-dom';
import { EventsOnce } from '@wailsjs/runtime/runtime';
import { SendPromptInput } from '@wailsjs/go/walletapp/WalletApp';

export interface PromptRequestTransferData {
  NicknameFrom: string;
  Amount: string;
  Fee: string;
  RecipientAddress: string;
}

export function Transfer() {
  const [error, setError] = useState<IErrorObject | null>(null);
  const navigate = useNavigate();
  const form = useRef(null);
  const { state } = useLocation();
  const req: promptRequest = state.req;
  const data: PromptRequestTransferData = req.Data;

  function validate(e: SyntheticEvent) {
    const formObject = parseForm(e);
    const { password } = formObject;

    if (!password || !password.length) {
      setError({ password: Intl.t('errors.PasswordRequired') });
      return false;
    }

    return true;
  }

  function save(e: SyntheticEvent) {
    const form = parseForm(e);
    const { password } = form;

    EventsOnce(events.promptResult, handleResult);

    SendPromptInput(password);
  }

  function handleResult(result: promptResult) {
    let { Success, CodeMessage } = result;

    if (!Success) {
      if (CodeMessage === ErrorCode.WrongPassword) {
        setError({ password: Intl.t(`errors.${CodeMessage}`) });
        return;
      }
      handleApplyResult(navigate, req, setError, false)(result);
    } else {
      handleApplyResult(navigate, req, setError, false)(result);
    }
  }

  async function handleSubmit(e: SyntheticEvent) {
    e.preventDefault();
    if (!validate(e)) return;

    save(e);
  }

  return (
    <Layout>
      <form ref={form} onSubmit={handleSubmit}>
        <h1 className="mas-title">Title</h1>
        <div className="mas-body pt-4 break-words">
          <div className="p-4 mb-2 bg-secondary rounded-lg w-full">
            <Balance amount={formatStandard(Number(toMAS(data.Amount)))} />
            <div className="mb-4 mt-2 mas-caption">
              {Intl.t('password-prompt.transfer.fee', { fee: data.Fee })}
            </div>
            <div className="flex items-center gap-2">
              {Intl.t('password-prompt.transfer.from')}
              <b>{data.NicknameFrom}</b>
            </div>
            <div className="flex items-center gap-2">
              {Intl.t('password-prompt.transfer.to')}
              <p>{maskAddress(data.RecipientAddress)}</p>
            </div>
          </div>
        </div>
        <div className="pt-4">
          <Password
            defaultValue=""
            name="password"
            placeholder="Password"
            error={error?.password}
          />
        </div>
        <div className="pt-4 flex gap-4">
          <Button variant={'secondary'} onClick={handleCancel}>
            {Intl.t('password-prompt.buttons.cancel')}
          </Button>
          <Button preIcon={<FiAlertCircle />} type="submit">
            {Intl.t('password-prompt.buttons.transfer')}
          </Button>
        </div>
      </form>
    </Layout>
  );
}
