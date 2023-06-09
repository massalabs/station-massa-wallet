import { useState, useRef, SyntheticEvent, ReactNode } from 'react';
import { useLocation, useNavigate } from 'react-router-dom';
import { SendPromptInput } from '../../wailsjs/go/walletapp/WalletApp';
import { EventsOnce } from '../../wailsjs/runtime';
import {
  events,
  promptAction,
  promptRequest,
  promptResult,
} from '../events/events';
import { handleApplyResult, handleCancel } from '../utils/utils';
import { parseForm } from '../utils/parseForm';

import { FiLock, FiTrash2 } from 'react-icons/fi';
import { Password, Button } from '@massalabs/react-ui-kit';
import { ErrorCode, IErrorObject } from '../utils';
import { Layout } from '../layouts/Layout/Layout';
import Intl from '../i18n/i18n';
import { formatStandard } from '../utils/MassaFormating';
import { toMAS } from '@massalabs/massa-web3';

interface PromptRequestDeleteData {
  Nickname: string;
  Balance: string;
}

interface PromptRequestTransferData {
  NicknameFrom: string;
  Amount: string;
  Fee: string;
  RecipientAddress: string;
}

function PasswordPrompt() {
  const navigate = useNavigate();
  const form = useRef(null);

  const { state } = useLocation();
  const req: promptRequest = state.req;

  const data: PromptRequestDeleteData | PromptRequestTransferData = req.Data;

  const { deleteReq, signReq, transferReq } = promptAction;

  function getTitle() {
    switch (req.Action) {
      case deleteReq:
        return Intl.t('password-prompt.title.delete');
      case signReq:
        return Intl.t('password-prompt.title.sign');
      default:
        return Intl.t(`password-prompt.title.${req.CodeMessage}`);
    }
  }

  function getButtonLabel() {
    switch (req.Action) {
      case deleteReq:
        return Intl.t('password-prompt.buttons.delete');
      case signReq:
        return Intl.t('password-prompt.buttons.sign');
      case transferReq:
        return Intl.t('password-prompt.buttons.transfer');
      default:
        return Intl.t('password-prompt.buttons.default');
    }
  }

  function getSubtitle() {
    switch (req.Action) {
      case deleteReq:
        return Intl.t('password-prompt.subtitle.delete');
      case signReq:
        return Intl.t('password-prompt.subtitle.sign');
      case promptAction.transferReq: {
        const transferData = req.Data as PromptRequestTransferData;
        return `Transfer ${formatStandard(
          toMAS(transferData.Amount).toString(),
        )} Massa from ${transferData.NicknameFrom} to
        ${transferData.RecipientAddress},
         with fee(s) ${transferData.Fee} nonaMassa`;
      }
      default:
        return Intl.t('password-prompt.subtitle.default');
    }
  }

  function getButtonIcon(): ReactNode {
    switch (req.Action) {
      case deleteReq:
        return <FiTrash2 />;
      case signReq:
      case transferReq:
        // no icon for sign and transfer requests
        return;
      default:
        return <FiLock />;
    }
  }

  // error state for the input password field
  const [error, setError] = useState<IErrorObject | null>(null);

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

  async function handleResult(result: promptResult) {
    let { Success, CodeMessage } = result;

    if (!Success) {
      if (CodeMessage === ErrorCode.WrongPassword) {
        setError({ password: Intl.t(`errors.${CodeMessage}`) });
        return;
      }
    } else {
      handleApplyResult(navigate, req, setError, false)(result);
    }
  }

  async function handleSubmit(e: SyntheticEvent) {
    e.preventDefault();
    if (!validate(e)) return;

    const form = parseForm(e);
    const { password } = form;

    if (
      req.Action === deleteReq &&
      (data as PromptRequestDeleteData).Balance !== '0'
    ) {
      navigate('/confirm-delete', { state: { req, password } });
    } else {
      save(e);
    }
  }

  return (
    <Layout>
      <form ref={form} onSubmit={handleSubmit}>
        <h1 className="mas-title">{getTitle()}</h1>
        <p className="mas-body pt-4 break-words max-w-xs">{getSubtitle()}</p>
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
          <Button preIcon={getButtonIcon()} type="submit">
            {getButtonLabel()}
          </Button>
        </div>
      </form>
    </Layout>
  );
}

export default PasswordPrompt;
