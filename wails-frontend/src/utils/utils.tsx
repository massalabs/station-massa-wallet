/* eslint-disable new-cap */
/* eslint-disable @typescript-eslint/no-explicit-any */
import { NavigateFunction } from 'react-router-dom';
import { Hide, AbortAction } from '../../wailsjs/go/walletapp/WalletApp';
import { WindowReloadApp } from '../../wailsjs/runtime';
import { promptResult, promptRequest } from '../events/events';
import { getErrorMessage } from './errors';

export const handleCancel = () => {
  AbortAction();
  hideAndReload();
};

export const hideAndReload = () => {
  Hide();
  // Reload the wails frontend
  WindowReloadApp();
};

export const handleApplyResult = (
  navigate: NavigateFunction,
  req: promptRequest,
  errMsgCb: (msg: any) => void,
  quitOnError = false,
) => {
  return (result: promptResult) => {
    if (result.Success) {
      navigate('/success', {
        state: { req },
      });
      setTimeout(hideAndReload, 2000);
    } else {
      req.Msg = getErrorMessage(result.CodeMessage);
      errMsgCb(req.Msg);
      navigate('/failure', {
        state: { req },
      });
      if (quitOnError || result.CodeMessage === 'Timeout-0001') {
        setTimeout(hideAndReload, 2000);
        return;
      }
    }
  };
};
