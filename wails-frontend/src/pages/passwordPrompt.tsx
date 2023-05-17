/* eslint-disable new-cap */
import { useState } from 'react';
import { useLocation, useNavigate } from 'react-router-dom';
import {
  ApplyPassword,
  ImportPrivateKey,
} from '../../wailsjs/go/walletapp/WalletApp';
import { EventsOnce } from '../../wailsjs/runtime';
import { events, promptAction, promptRequest } from '../events/events';
import { handleApplyResult, handleCancel } from '../utils/utils';

import { FiLock } from 'react-icons/fi';
import { Password, Button } from '@massalabs/react-ui-kit';

function PasswordPrompt() {
  const navigate = useNavigate();

  const { state } = useLocation();
  const req: promptRequest = state.req;

  const { newPasswordReq, importReq, deleteReq } = promptAction;
  const isNewPassword = req.Action === newPasswordReq;
  const isImport = req.Action === importReq;

  const [resultMsg, setResultMsg] = useState('');
  const [password, setPassword] = useState('');
  const [passwordConfirm, setPasswordConf] = useState<string | undefined>(
    undefined,
  );

  function getButtonLabel() {
    switch (req.Action) {
      case deleteReq:
        return 'Delete';
      case newPasswordReq:
        return 'Define';
      default:
        return 'Apply';
    }
  }

  function getSubtitle() {
    switch (req.Action) {
      case deleteReq:
        return 'Enter your password to delete your wallet';
      case newPasswordReq:
        return 'Enter a secure password';
      case importReq:
        return 'Define a new password';
      default:
        return 'Enter your password below to validate';
    }
  }

  function isDisabled() {
    if (isNewPassword || isImport) {
      return !(password && passwordConfirm);
    }
    return !password;
  }

  function applyPassword() {
    if (password.length) {
      if (passwordConfirm && password !== passwordConfirm) {
        setResultMsg("Passwords don't match");
        return;
      }
      EventsOnce(
        events.passwordResult,
        handleApplyResult(
          navigate,
          req,
          setResultMsg,
          state.req.Action === promptAction.importReq,
        ),
      );

      if (state.req.Action === promptAction.importReq) {
        return ImportPrivateKey(state.pkey, state.nickname, password);
      }
      return ApplyPassword(password);
    }
  }

  return (
    <div className="bg-primary min-h-screen">
      <div
        className="flex flex-col justify-center h-screen
        max-w-xs min-w-fit text-f-primary m-auto"
      >
        <h1 className="mas-title">{req.Msg}</h1>
        <p className="mas-body pt-4">{getSubtitle()}</p>
        <div className="pt-4">
          <Password
            onChange={(e) => setPassword(e.target.value)}
            value={password}
            name="input"
            placeholder="Enter your password"
          />
        </div>
        {(isNewPassword || isImport) && (
          <div className="pt-4">
            <Password
              onChange={(e) => setPasswordConf(e.target.value)}
              value={passwordConfirm}
              name="input"
              placeholder="Confirm your password"
            />
          </div>
        )}
        <p className="pt-4 mas-body">{resultMsg}</p>
        <div className="pt-4 flex gap-4">
          <div className="max-w-min">
            <Button variant={'secondary'} onClick={handleCancel}>
              Cancel
            </Button>
          </div>
          <Button
            preIcon={<FiLock />}
            onClick={applyPassword}
            disabled={isDisabled()}
          >
            {getButtonLabel()}
          </Button>
        </div>
      </div>
    </div>
  );
}

export default PasswordPrompt;
