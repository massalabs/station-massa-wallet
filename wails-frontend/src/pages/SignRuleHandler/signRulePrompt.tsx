import { SyntheticEvent, useRef, useState } from 'react';

import {
  Button,
  Clipboard,
  Password,
  RadioButton,
  Tooltip,
  maskAddress,
} from '@massalabs/react-ui-kit';
import { walletapp } from '@wailsjs/go/models';
import {
  SendPromptInput,
  SendExpiredSignRulePromptInput,
} from '@wailsjs/go/walletapp/WalletApp';
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

enum ExpiredRuleAction {
  Refresh = 'refresh',
  Delete = 'delete',
  None = 'none',
}

export function SignRule() {
  const [error, setError] = useState<IErrorObject | null>(null);
  const [errorMessage, setErrorMessage] = useState<string>('');
  const [
    refreshDeleteExpiredSignRuleError,
    setRefreshDeleteExpiredSignRuleError,
  ] = useState<string>('');
  const [expiredRuleAction, setExpiredRuleAction] = useState<ExpiredRuleAction>(
    ExpiredRuleAction.None,
  );

  const navigate = useNavigate();
  const form = useRef(null);
  const { state } = useLocation();
  const req: promptRequest = state.req;
  const data: ruleRequestData = req.Data;

  let action = '';
  let showWarning = false;
  let winHeight = 0;
  const { addSignRule, updateSignRule, deleteSignRule, expiredSignRule } =
    walletapp.PromptRequestAction;

  const { EventType } = walletapp;

  switch (req.Action) {
    case addSignRule:
      action = signRuleActionStr.addSignRule;
      showWarning = data.SignRule.Enabled;
      winHeight = 680;
      break;
    case updateSignRule:
      action = signRuleActionStr.updateSignRule;
      showWarning = data.SignRule.Enabled;
      winHeight = 670;
      break;
    case deleteSignRule:
      action = signRuleActionStr.deleteSignRule;
      winHeight = 510;
      break;
    case expiredSignRule:
      action = signRuleActionStr.expiredSignRule;
      winHeight = 630;
      break;
    default:
      setErrorMessage(Intl.t(`errors.unknownRuleRequest`));
  }

  const winWidth = 460;
  if (data.SignRule.AuthorizedOrigin) {
    winHeight += 50;
  }
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

    if (password === undefined) {
      setErrorMessage(Intl.t('errors.PasswordRequired'));
      return;
    }

    EventsOnce(EventType.promptResult, handleResult);

    if (req.Action === expiredSignRule) {
      SendExpiredSignRulePromptInput(
        password,
        expiredRuleAction === ExpiredRuleAction.Delete,
      );
      return;
    }
    SendPromptInput(password);
  }

  async function handleSubmitPassword(e: SyntheticEvent) {
    e.preventDefault();
    if (!validate(e, setError)) return;

    // Additional validation for expired sign rule
    if (
      req.Action === expiredSignRule &&
      expiredRuleAction === ExpiredRuleAction.None
    ) {
      setRefreshDeleteExpiredSignRuleError(
        Intl.t('signRule.expiredSignRule.radioFormNotChecked'),
      );
      return;
    }

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
            <div className="text-sm break-words text-f-disabled-1 mb-4">
              {data.Description}
            </div>
          </>
        )}

        <div className="mb-4">
          <Account
            nickname={data.Nickname}
            walletAddress={data.WalletAddress}
          />
        </div>

        {showWarning && (
          <div className="p-4 flex items-center mb-4">
            <FiAlertTriangle size={42} className="text-s-warning" />
            <div className="ml-2 text-s-warning text-sm">
              <p>{Intl.t(`signRule.behavior.${data.SignRule.RuleType}`)}</p>
            </div>
          </div>
        )}

        <div className="bg-secondary p-4 mb-4 rounded-md space-y-4">
          <div className="flex items-center space-x-2 min-h-8 justify-between">
            <strong className="min-w-0 flex-shrink-0">
              {Intl.t('signRule.name')}:
            </strong>
            <span className="text-f-primary-1 truncate">
              {data.SignRule.Name}
            </span>
          </div>

          <div className="flex items-center space-x-2 justify-between min-h-8">
            <strong className="min-w-0 flex-shrink-0">
              {Intl.t('signRule.contract')}:
            </strong>
            {isAllContract ? (
              <span className="text-f-primary-1">All</span>
            ) : (
              <Clipboard
                customClass="flex h-6 px-0"
                rawContent={data.SignRule.Contract}
                displayedContent={maskAddress(data.SignRule.Contract, 12)}
              />
            )}
          </div>

          {data.SignRule.AuthorizedOrigin && (
            <div className="flex items-center space-x-2 justify-between">
              <strong className="min-w-0 flex-shrink-0">
                <Tooltip
                  body={Intl.t('signRule.authorizedOrigin')}
                  placement="top"
                  tooltipClassName="mas-caption py-0"
                  triggerClassName="flex items-center py-0"
                >
                  Origin :
                </Tooltip>
              </strong>
              <Tooltip
                body={data.SignRule.AuthorizedOrigin}
                placement="top"
                triggerClassName="flex items-center py-0"
                tooltipClassName="mas-caption max-w-80 py-0"
              >
                <Clipboard
                  customClass="ml-1 truncate max-w-48 px-0 py-0"
                  rawContent={data.SignRule.AuthorizedOrigin}
                  displayedContent={
                    data.SignRule.AuthorizedOrigin.length > 28
                      ? data.SignRule.AuthorizedOrigin.slice(0, 26) + '...'
                      : data.SignRule.AuthorizedOrigin
                  }
                />
              </Tooltip>
            </div>
          )}

          <div className="flex items-center space-x-2 min-h-8 justify-between">
            <strong className="min-w-0 flex-shrink-0">
              {Intl.t('signRule.ruleType')}:
            </strong>
            <span className="text-f-primary-1">
              {Intl.t(`signRule.rulesTypeDesc.${data.SignRule.RuleType}`)}
            </span>
          </div>

          <div className="flex items-center space-x-2 min-h-8 justify-between">
            <strong className="min-w-0 flex-shrink-0">
              {Intl.t('signRule.enabled')}:
            </strong>
            <span className="text-f-primary-1">
              {data.SignRule.Enabled ? 'Yes' : 'No'}
            </span>
          </div>
        </div>

        {req.Action === expiredSignRule && (
          <div className="mb-4">
            <div className="flex items-center space-x-6">
              <div className="flex items-center">
                <RadioButton
                  name="expiredRuleAction"
                  value={ExpiredRuleAction.Refresh}
                  checked={expiredRuleAction === 'refresh'}
                  onChange={() =>
                    setExpiredRuleAction(ExpiredRuleAction.Refresh)
                  }
                />
                <p
                  className="h-full ml-3 pb-1 cursor-pointer"
                  onClick={() =>
                    setExpiredRuleAction(ExpiredRuleAction.Refresh)
                  }
                >
                  {Intl.t(
                    'signRule.expiredSignRule.expiredRuleOptions.refreshRule',
                  )}
                </p>
              </div>
              <div className="flex items-center">
                <RadioButton
                  name="expiredRuleAction"
                  value={ExpiredRuleAction.Delete}
                  checked={expiredRuleAction === 'delete'}
                  onChange={() =>
                    setExpiredRuleAction(ExpiredRuleAction.Delete)
                  }
                />
                <p
                  className="h-full ml-3 pb-1 cursor-pointer"
                  onClick={() => setExpiredRuleAction(ExpiredRuleAction.Delete)}
                >
                  {Intl.t(
                    'signRule.expiredSignRule.expiredRuleOptions.deleteRule',
                  )}
                </p>
              </div>
            </div>
            {refreshDeleteExpiredSignRuleError && (
              <p className="mt-2 text-s-error mas-body">
                {refreshDeleteExpiredSignRuleError}
              </p>
            )}
          </div>
        )}

        <div className="mb-4">
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
        <div className="flex gap-4">
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
