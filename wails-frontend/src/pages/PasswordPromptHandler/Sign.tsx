import { SyntheticEvent, useRef, useState } from 'react';

import { Mas } from '@massalabs/massa-web3';
import { Button, Password } from '@massalabs/react-ui-kit';
import { walletapp } from '@wailsjs/go/models';
import { SendSignPromptInput } from '@wailsjs/go/walletapp/WalletApp';
import { EventsOnce, WindowSetSize } from '@wailsjs/runtime';
import { FiLock } from 'react-icons/fi';
import { useLocation, useNavigate } from 'react-router-dom';

import { BuySellRoll } from './BuySellRoll/BuySellRoll';
import { CallSc } from './CallSC/CallSc';
import { NetworkName } from './components/NetworkName';
import { ExecuteSC } from './ExecuteSC.tsx/ExecuteSc';
import { PlainText } from './PlainText/PlainText';
import { OperationCost } from './SignComponentUtils/OperationCost';
import { Transaction } from './Transaction/Transaction';
import {
  OPER_BUY_ROLL,
  OPER_CALL_SC,
  OPER_EXECUTE_SC,
  OPER_PLAIN_TEXT,
  OPER_SELL_ROLL,
  OPER_TRANSACTION,
} from '@/const/operations';
import { promptRequest, promptResult } from '@/events';
import Intl from '@/i18n/i18n';
import { SignLayout } from '@/layouts/Layout/SignLayout';
import {
  ErrorCode,
  IErrorObject,
  handleApplyResult,
  handleCancel,
  parseForm,
} from '@/utils';

export interface AssetInfo {
  address: string;
  decimals: number;
  name: string;
  symbol: string;
}

export interface SignBodyProps {
  Description: string;
  Fees: string; // in nanoMassa
  MinFees: string; // in MAS
  OperationType: number;
  Coins: string; // in nanoMassa
  Address: string;
  Function: string;
  MaxCoins: string;
  WalletAddress: string;
  Nickname: string;
  RollCount: number;
  RecipientAddress: string;
  RecipientNickname: string;
  Amount: string;
  PlainText: string;
  AllowFeeEdition: boolean;
  DisplayData: boolean;
  ChainID: number;
  Assets: AssetInfo[];
  Parameters: string; // base64
  DeployedByteCodeSize: number; // for executeSC of type deploySC
  DeployedCoins: string; // for executeSC of type deploySC
  EnabledSignRule: string | null;
  children?: React.ReactNode;
}

export function validate(
  e: SyntheticEvent,
  setError: (error: IErrorObject) => void,
) {
  const formObject = parseForm(e);
  const { password } = formObject;

  if (!password.length) {
    setError({ password: Intl.t('errors.PasswordRequired') });
    return false;
  }

  return true;
}

export function Sign() {
  const navigate = useNavigate();
  const form = useRef(null);
  const { state } = useLocation();
  const req: promptRequest = state.req;
  const signData = req.Data as SignBodyProps;
  const [error, setError] = useState<IErrorObject | null>(null);
  const [errorMessage, setErrorMessage] = useState<string>('');
  const [fees, setFees] = useState<Mas.Mas>(BigInt(signData.Fees || 0));
  const [isEditing, setIsEditing] = useState(false);

  const { EventType } = walletapp;

  function save(e: SyntheticEvent) {
    const form = parseForm(e);
    const { password } = form;

    EventsOnce(EventType.promptResult, handleResult);

    SendSignPromptInput(password, fees.toString());
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
    if (!req.DisablePassword) {
      e.preventDefault();
      if (!validate(e, setError)) return;
    }
    save(e);
  }

  function getTitle(operation: number | undefined) {
    if (operation === OPER_PLAIN_TEXT)
      return Intl.t('password-prompt.title.sign-message');

    return Intl.t('password-prompt.title.sign');
  }

  const winWidth = 460;
  const winHeight = 800;

  WindowSetSize(winWidth, winHeight);

  const operationCostsArgs = {
    fees,
    minFees: Mas.fromString(signData.MinFees || '0'),
    setFees,
    isEditing,
    setIsEditing,
    allowFeeEdition: signData.AllowFeeEdition,
    DeployedByteCodeSize: signData.DeployedByteCodeSize,
    DeployedCoins: signData.DeployedCoins,
  };

  return (
    <SignLayout>
      <form ref={form} onSubmit={handleSubmit}>
        {signData.OperationType !== OPER_PLAIN_TEXT && (
          <NetworkName chainId={signData.ChainID} />
        )}
        <h1 className="mas-title">{getTitle(signData.OperationType)}</h1>
        <div className="mas-body pt-4 break-words">
          {(() => {
            switch (signData.OperationType) {
              case OPER_CALL_SC:
                return (
                  <CallSc {...signData}>
                    <OperationCost
                      coins={signData.Coins}
                      {...operationCostsArgs}
                    />
                  </CallSc>
                );
              case OPER_EXECUTE_SC:
                return (
                  <ExecuteSC {...signData}>
                    <OperationCost {...operationCostsArgs} />
                  </ExecuteSC>
                );
              case OPER_BUY_ROLL:
              case OPER_SELL_ROLL:
                return (
                  <BuySellRoll {...signData}>
                    <OperationCost {...operationCostsArgs} />
                  </BuySellRoll>
                );
              case OPER_TRANSACTION:
                return (
                  <Transaction {...signData}>
                    <OperationCost {...operationCostsArgs} />
                  </Transaction>
                );
              case OPER_PLAIN_TEXT:
                return <PlainText {...signData} />;
            }
          })()}
        </div>
        {!req.DisablePassword && (
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
        )}
        <div className="pt-4 flex gap-4">
          <Button variant={'secondary'} onClick={handleCancel}>
            {Intl.t('password-prompt.buttons.cancel')}
          </Button>
          <Button preIcon={<FiLock />} disabled={isEditing} type="submit">
            {Intl.t('password-prompt.buttons.sign')}
          </Button>
        </div>
      </form>
    </SignLayout>
  );
}
