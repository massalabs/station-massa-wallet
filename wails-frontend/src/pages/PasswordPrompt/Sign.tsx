import { SyntheticEvent, useRef, useState } from 'react';

import { Button, Password } from '@massalabs/react-ui-kit';
import { SendPromptInput } from '@wailsjs/go/walletapp/WalletApp';
import { EventsOnce } from '@wailsjs/runtime/runtime';
import { FiLock } from 'react-icons/fi';
import { useLocation, useNavigate } from 'react-router-dom';

import { validate } from './Default';
import { events, promptRequest, promptResult } from '@/events/events';
import Intl from '@/i18n/i18n';
import {
  ErrorCode,
  IErrorObject,
  handleApplyResult,
  handleCancel,
  maskAddress,
  parseForm,
} from '@/utils';

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
    }
    handleApplyResult(navigate, req, setError, false)(result);
  }

  async function handleSubmit(e: SyntheticEvent) {
    e.preventDefault();
    if (!validate(e, setError)) return;

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
            {/* components will be returned in switch statement
            right now this is a minimalist refactor
            */}
            <div>Operation Type: {signData.OperationType}</div>
            {(() => {
              switch (signData.OperationType) {
                case 'Call SC':
                  return (
                    <>
                      <div>Gas Limit: {signData.GasLimit}</div>
                      <div>Coins: {signData.Coins}</div>
                      <div>To: {maskAddress(signData.Address)}</div>
                      <div>From: {maskAddress(signData.WalletAddress)}</div>
                      <div>Function: {signData.Function}</div>
                    </>
                  );
                case 'Execute SC':
                  return (
                    <>
                      <div>Max Coins: {signData.MaxCoins}</div>
                      <div>Max Gas: {signData.MaxGas}</div>
                    </>
                  );
                default:
                  return <div>Other Sign Data Content</div>;
              }
            })()}
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
            <Button preIcon={<FiLock />} type="submit">
              {Intl.t('password-prompt.buttons.sign')}
            </Button>
          </div>
        </form>
      </div>
    </div>
  );
}
