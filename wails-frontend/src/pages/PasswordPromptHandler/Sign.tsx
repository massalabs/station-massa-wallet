import { SyntheticEvent, useRef, useState } from 'react';

import { Button, Password } from '@massalabs/react-ui-kit';
import { SendPromptInput } from '@wailsjs/go/walletapp/WalletApp';
import { EventsOnce } from '@wailsjs/runtime/runtime';
import { FiLock } from 'react-icons/fi';
import { useLocation, useNavigate } from 'react-router-dom';

import { validate } from './Default';
import { CallSc } from './SignSC/CallSc';
import { events, promptRequest, promptResult } from '@/events/events';
import Intl from '@/i18n/i18n';
import { SignLayout } from '@/layouts/Layout/SignLayout';
import {
  ErrorCode,
  IErrorObject,
  handleApplyResult,
  handleCancel,
  maskAddress,
  parseForm,
} from '@/utils';

export interface PromptRequestData {
  Description: string;
  OperationID: number;
  GasLimit: number;
  Coins: number;
  Address: string;
  Function: string;
  WalletAddress: string;
  OperationType?: string;
  MaxCoins: number;
  MaxGas: number;
  RollCount: number;
}

export function Sign() {
  const navigate = useNavigate();
  const form = useRef(null);
  const { state } = useLocation();
  const req: promptRequest = state.req;
  const [error, setError] = useState<IErrorObject | null>(null);
  const signData = req.Data as PromptRequestData;

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
    <SignLayout>
      <form ref={form} onSubmit={handleSubmit}>
        <h1 className="mas-title">{Intl.t('password-prompt.title.sign')}</h1>
        <div className="mas-body pt-4 break-words">
          {/* components will be returned in switch statement
            right now this is a minimalist refactor
            */}
          {(() => {
            switch (signData.OperationType) {
              case 'Call SC':
                return (
                  <>
                    <CallSc {...signData} />
                  </>
                );
              case 'Execute SC':
                return (
                  <>
                    <div>Description: {signData.Description}</div>
                    <div>Max Coins: {signData.MaxCoins}</div>
                    <div>Max Gas: {signData.MaxGas}</div>
                    <div>From: {maskAddress(signData.WalletAddress)}</div>
                  </>
                );
              case 'Buy Roll':
              case 'Sell Roll':
                return (
                  <>
                    <div>Description: {signData.Description}</div>
                    <div>Rolls Count: {signData.RollCount}</div>
                    <div>From: {maskAddress(signData.WalletAddress)}</div>
                  </>
                );

              default:
                return (
                  <>
                    <div>Description: {signData.Description}</div>
                    <div>From: {maskAddress(signData.WalletAddress)}</div>
                  </>
                );
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
    </SignLayout>
  );
}
