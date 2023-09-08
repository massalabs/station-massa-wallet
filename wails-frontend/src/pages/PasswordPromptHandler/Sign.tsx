import { SyntheticEvent, useRef, useState } from 'react';

import { Button, Password } from '@massalabs/react-ui-kit';
import { SendPromptInput } from '@wailsjs/go/walletapp/WalletApp';
import { EventsOnce } from '@wailsjs/runtime/runtime';
import { FiLock } from 'react-icons/fi';
import { useLocation, useNavigate } from 'react-router-dom';

import { BuySellRoll } from './BuySellRoll/BuySellRoll';
import { validate } from './Default';
import { ExecuteSC } from './ExecuteSC.tsx/ExecuteSc';
import { PlainText } from './PlainText/PlainText';
import { CallSc } from './SignSC/CallSc';
import { Transaction } from './Transaction/Transaction';
import {
  OPER_BUY_ROLL,
  OPER_CALL_SC,
  OPER_EXECUTE_SC,
  OPER_PLAIN_TEXT,
  OPER_SELL_ROLL,
  OPER_TRANSACTION,
} from '@/const/operations';
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
  Fees: number;
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
  RecipientAddress: string;
  Amount: string;
  PlainText: string;
  DisplayData: boolean;
  Nickname: string;
}

export function Sign() {
  const navigate = useNavigate();
  const form = useRef(null);
  const { state } = useLocation();
  const req: promptRequest = state.req;
  const [error, setError] = useState<IErrorObject | null>(null);
  const [errorMessage, setErrorMessage] = useState<string>('');
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
      } else {
        setErrorMessage(Intl.t(`errors.${CodeMessage}`));
      }
      return;
    }
    handleApplyResult(navigate, req, setError, false)(result);
  }

  async function handleSubmit(e: SyntheticEvent) {
    e.preventDefault();
    if (!validate(e, setError)) return;

    save(e);
  }

  function _getTitle(operation: string | undefined) {
    if (operation === OPER_PLAIN_TEXT)
      return Intl.t('password-prompt.title.sign-message');

    return Intl.t('password-prompt.title.sign');
  }

  return (
    <SignLayout>
      <form ref={form} onSubmit={handleSubmit}>
        <h1 className="mas-title">{_getTitle(signData.OperationType)}</h1>
        <div className="mas-body pt-4 break-words">
          {(() => {
            switch (signData.OperationType) {
              case OPER_CALL_SC:
                return (
                  <>
                    <CallSc {...signData} />
                  </>
                );
              case OPER_EXECUTE_SC:
                return (
                  <>
                    <ExecuteSC {...signData} />
                  </>
                );
              case OPER_BUY_ROLL:
              case OPER_SELL_ROLL:
                return (
                  <>
                    <BuySellRoll {...signData} />
                  </>
                );
              case OPER_TRANSACTION:
                return (
                  <>
                    <Transaction {...signData} />;
                  </>
                );
              case OPER_PLAIN_TEXT:
                return (
                  <>
                    <PlainText {...signData} />
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
          {errorMessage && (
            <p className="mt-2 text-s-error mas-body">{errorMessage}</p>
          )}
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
