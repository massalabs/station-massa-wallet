import { useState } from 'preact/hooks';
import {
  AbortAction,
  ApplyPassword,
  Hide,
} from '../../wailsjs/go/walletapp/WalletApp';
import { EventsOnce, WindowReloadApp } from '../../wailsjs/runtime';
import { h } from 'preact';
import { events, promptAction, promptRequest } from '../events/events';
import { useLocation, useNavigate } from 'react-router-dom';

const PasswordPrompt = () => {
  const navigate = useNavigate();

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
      default:
        return 'Enter your password below to validate';
    }
  };

  const [resultMsg, setResultMsg] = useState('');
  const [password, setPassword] = useState('');
  const [passwordConfirm, setPasswordConf] = useState(undefined);

  const hideAndReload = () => {
    Hide();
    // Reload the wails frontend
    WindowReloadApp();
  };

  const handleCancel = () => {
    AbortAction();
    hideAndReload();
  };

  const handleApplyResult = (result: any) => {
    console.log('result', result);
    if (result.Success) {
      navigate('/success', {
        state: { msg: 'Password applied successfully!' },
      });
      setTimeout(hideAndReload, 2000);
    } else {
      setResultMsg(result.Data);
      if (result.Error === 'timeoutError') {
        setTimeout(hideAndReload, 2000);
      }
    }
  };

  const applyPassword = (event: any) => {
    console.log('event', event);
    if (password.length) {
      if (passwordConfirm && password !== passwordConfirm) {
        setResultMsg("Passwords don't match");
        return;
      }
      EventsOnce(events.passwordResult, handleApplyResult);
      ApplyPassword(password);
    }
  };

  const updatePassword = (e: any) => setPassword(e.target.value);
  const updatePasswordConfirm = (e: any) => setPasswordConf(e.target.value);

  return (
    <section class="PasswordPrompt">
      <div>{req.Msg}</div>
      <div id="result" className="result">
        {baselineStr(req)}
      </div>
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
      {req.Action === promptAction.newPasswordReq && (
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
      <div id="result" className="result">
        {resultMsg}
      </div>
    </section>
  );
};

export default PasswordPrompt;
