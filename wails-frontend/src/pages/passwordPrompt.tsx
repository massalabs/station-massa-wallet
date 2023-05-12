/* eslint-disable new-cap */
import { useState } from 'react';
import {
  ApplyPassword,
  ImportPrivateKey,
} from '../../wailsjs/go/walletapp/WalletApp';
import { EventsOnce } from '../../wailsjs/runtime';
import { events, promptAction, promptRequest } from '../events/events';
import { useLocation, useNavigate } from 'react-router-dom';
import { handleApplyResult, handleCancel } from '../utils/utils';

const PasswordPrompt = () => {
  const nav = useNavigate();

  const { state } = useLocation();
  const req: promptRequest = state.req;

  const applyStr = (req: promptRequest) => {
    switch (req.Action) {
      case promptAction.deleteReq:
        return 'Delete';
      case promptAction.newPasswordReq:
        return 'Define a password';
      default:
        return 'Apply';
    }
  };

  const baselineStr = (req: promptRequest) => {
    switch (req.Action) {
      case promptAction.deleteReq:
        return 'Enter your password to delete your wallet';
      case promptAction.newPasswordReq:
        return 'Enter a secure password';
      case promptAction.importReq:
        return 'Define a new password';
      default:
        return 'Enter your password below to validate';
    }
  };

  const [resultMsg, setResultMsg] = useState<string | []>('');
  const [password, setPassword] = useState<string>('');
  const [passwordConfirm, setPasswordConf] = useState<string | undefined>(
    undefined,
  );

  const applyPassword = () => {
    if (password.length) {
      if (passwordConfirm && password !== passwordConfirm) {
        setResultMsg("Passwords don't match");
        return;
      }
      EventsOnce(
        events.passwordResult,
        handleApplyResult(
          nav,
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
  };

  const updatePassword = (e: React.ChangeEvent<HTMLInputElement>) =>
    setPassword(e.target.value);
  const updatePasswordConfirm = (e: React.ChangeEvent<HTMLInputElement>) =>
    setPasswordConf(e.target.value);

  return (
    <section className="PasswordPrompt">
      <div>{req.Msg}</div>
      <div className="baseline">{baselineStr(req)}</div>
      <div id="input" className="input-box">
        <input
          id="name"
          className="input"
          onChange={updatePassword}
          autoComplete="off"
          name="input"
          type="password"
          placeholder="Enter your password"
        />
      </div>
      {(req.Action === promptAction.newPasswordReq ||
        req.Action === promptAction.importReq) && (
        <div id="input" className="input-box">
          <input
            id="name"
            className="input"
            onChange={updatePasswordConfirm}
            autoComplete="off"
            name="input"
            type="password"
            placeholder="Confirm your password"
          />
        </div>
      )}
      <div>
        <button className="btn" onClick={handleCancel}>
          Cancel
        </button>
        <button className="btn" onClick={applyPassword}>
          {applyStr(req)}
        </button>
      </div>
      <div className="baseline">{resultMsg}</div>
    </section>
  );
};

export default PasswordPrompt;
