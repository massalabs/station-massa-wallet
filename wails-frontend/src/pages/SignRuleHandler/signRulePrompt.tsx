import { SyntheticEvent, useRef, useState } from 'react';

import {
  Button,
  Clipboard,
  Password,
  maskAddress,
} from '@massalabs/react-ui-kit';
import { walletapp } from '@wailsjs/go/models';
import { SendPromptInput } from '@wailsjs/go/walletapp/WalletApp';
import { EventsOnce, WindowSetSize } from '@wailsjs/runtime/runtime';
import { FiAlertTriangle } from 'react-icons/fi';
import { useLocation, useNavigate } from 'react-router-dom';

import { ruleRequestData, signRuleActionStr } from './types';
import { Account } from '../PasswordPromptHandler/components/account';
import { promptRequest, promptResult } from '@/events';
import Intl from '@/i18n/i18n';
import { Layout } from '@/layouts/Layout/Layout';
import { validate } from '@/pages/PasswordPromptHandler/Sign';
import {
  ErrorCode,
  IErrorObject,
  handleApplyResult,
  handleCancel,
  parseForm,
} from '@/utils';

export function SignRule() {
  const [error, setError] = useState<IErrorObject | null>(null);
  const [errorMessage, setErrorMessage] = useState<string>('');

  const navigate = useNavigate();
  const form = useRef(null);
  const { state } = useLocation();
  const req: promptRequest = state.req;
  const data: ruleRequestData = req.Data;

  let action = '';
  let showWarning = false;

  const { addSignRule, updateSignRule, deleteSignRule } =
    walletapp.PromptRequestAction;

  const { EventType } = walletapp;

  switch (req.Action) {
    case addSignRule:
      action = signRuleActionStr.addSignRule;
      showWarning = data.SignRule.Enabled;
      break;
    case updateSignRule:
      action = signRuleActionStr.updateSignRule;
      showWarning = data.SignRule.Enabled;
      break;
    case deleteSignRule:
      action = signRuleActionStr.deleteSignRule;
      break;
    default:
      setErrorMessage(Intl.t(`errors.unknownRuleRequest`));
  }

  const winWidth = 460;
  const winHeight = showWarning ? 650 : 540;

  WindowSetSize(winWidth, winHeight);

  function handleResult(result: promptResult) {
    let { Success, CodeMessage } = result;
    setError(null);
    setErrorMessage('');

    if (!Success) {
      if (CodeMessage === ErrorCode.WrongPassword) {
        setError({ password: Intl.t(`errors.${CodeMessage}`) });
      } else {
        setErrorMessage(Intl.t(`errors.${action}`));
      }
      return;
    }
    handleApplyResult(navigate, req, setError, false)(result);
  }

  function submitPassword(e: SyntheticEvent) {
    const form = parseForm(e);
    const { password } = form;

    EventsOnce(EventType.promptResult, handleResult);

    SendPromptInput(password);
  }

  async function handleSubmitPassword(e: SyntheticEvent) {
    e.preventDefault();
    if (!validate(e, setError)) return;
    submitPassword(e);
  }

  const isAllContract = data.SignRule.Contract === '*';

  return (
    <Layout>
      <form ref={form} onSubmit={handleSubmitPassword}>
        <h1 className="mas-title">{Intl.t(`signRule.${action}.title`)}</h1>
        {data.Description && (
          <>
            <div className="mas-body pt-4 break-words">
              {Intl.t('signRule.description')}:
            </div>
            <div className="text-sm break-words text-f-disabled-1">
              {data.Description}
            </div>
          </>
        )}

        <Account nickname={data.Nickname} walletAddress={data.WalletAddress} />

        {showWarning && (
          <div className="p-4 flex items-center">
            <FiAlertTriangle size={42} className="text-s-warning " />
            <div className="ml-2 text-s-warning text-sm">
              <p>{Intl.t(`signRule.behavior.${data.SignRule.RuleType}`)}</p>
            </div>
          </div>
        )}

        <div className="bg-secondary p-2 mt-4 rounded-md">
          <div className="flex-row">
            <p className="flex justify-between">
              <strong>{Intl.t('signRule.name')}:</strong>
              <div className="flex flex-row justify-right">
                {' '}
                {data.SignRule.Name}
              </div>
            </p>
          </div>
          <div className="flex-row items-center">
            <p className="flex justify-between">
              <strong>{Intl.t('signRule.contract')}:</strong>
              {isAllContract ? (
                'All'
              ) : (
                <Clipboard
                  customClass="flex h-6 pl-20 ml-20"
                  rawContent={data.SignRule.Contract}
                  displayedContent={maskAddress(data.SignRule.Contract, 6)}
                />
              )}
            </p>
          </div>

          <p className="flex justify-between">
            <strong>{Intl.t('signRule.ruleType')}:</strong>{' '}
            {Intl.t(`signRule.rulesTypeDesc.${data.SignRule.RuleType}`)}
          </p>
          <p className="flex justify-between">
            <strong>{Intl.t('signRule.enabled')}:</strong>{' '}
            {data.SignRule.Enabled ? 'Yes' : 'No'}
          </p>
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
          <Button type="submit">
            {Intl.t('password-prompt.buttons.default')}
          </Button>
        </div>
      </form>
    </Layout>
  );
}
