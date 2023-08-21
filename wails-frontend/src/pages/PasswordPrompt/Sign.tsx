// import { maskAddress } from '@/utils';
// import { Button, Password } from '@massalabs/react-ui-kit';
import {
  ErrorCode,
  IErrorObject,
  handleApplyResult,
  handleCancel,
  maskAddress,
  parseForm,
} from '@/utils';
import { Button, Password } from '@massalabs/react-ui-kit';
import { SyntheticEvent, useRef, useState } from 'react';
import { FiAlertCircle } from 'react-icons/fi';

import Intl from '@/i18n/i18n';
import { EventsOnce } from '@wailsjs/runtime/runtime';
import { SendPromptInput } from '@wailsjs/go/walletapp/WalletApp';
import { useLocation, useNavigate } from 'react-router-dom';
import { events, promptRequest, promptResult } from '@/events/events';

export interface PromptRequestCallSCData {
  OperationID: number;
  GasLimit: number;
  Coins: number;
  Address: string;
  Function: string;
  WalletAddress: string;
  OperationType?: string;
  MaxCoins: number;
  MaxGas: number;
}

export function Sign() {
  const navigate = useNavigate();
  const form = useRef(null);
  const { state } = useLocation();
  const req: promptRequest = state.req;
  const [error, setError] = useState<IErrorObject | null>(null);
  const signData = req.Data as PromptRequestCallSCData;

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
    <div className="bg-primary min-h-screen">
      <div
        className="flex flex-col justify-center h-screen
        max-w-sm text-f-primary m-auto"
      >
        <form ref={form} onSubmit={handleSubmit}>
          <h1 className="mas-title">{Intl.t('password-prompt.title.sign')}</h1>
          <h3>{Intl.t('password-prompt.subtitle.sign')}</h3>
          <div className="mas-body pt-4 break-words">
            <div>Operation Type: {signData.OperationType}</div>
            {signData.OperationType === 'Call SC' ? (
              <>
                <div>Gas Limit: {signData.GasLimit}</div>
                <div>Coins: {signData.Coins}</div>
                <div>To: {maskAddress(signData.Address)}</div>
                <div>From: {maskAddress(signData.WalletAddress)}</div>
                <div>Function: {signData.Function}</div>
              </>
            ) : signData.OperationType === 'Execute SC' ? ( // Handle the Execute SC case
              <>
                <div>Max Coins: {signData.MaxCoins}</div>
                <div>Max Gas: {signData.MaxGas}</div>
              </>
            ) : (
              <div>Other Sign Data Content</div>
            )}
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
              {Intl.t('password-prompt.buttons.sign')}
            </Button>
          </div>
        </form>
      </div>
    </div>
  );
}
