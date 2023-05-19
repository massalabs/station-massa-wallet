/* eslint-disable new-cap */
/* eslint-disable @typescript-eslint/no-explicit-any */
import { NavigateFunction } from 'react-router-dom';
import { Hide, AbortAction } from '../../wailsjs/go/walletapp/WalletApp';
import { WindowReloadApp } from '../../wailsjs/runtime';
import { promptResult, promptRequest } from '../events/events';
import { errorsEN } from './errors';

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
      req.Msg = result.Error;
      if (result.Error) {
        errMsgCb(errorsEN[result.Error]);
      } else {
        // we keep this errMsgCb call here until we migrate all the error
        // for retro-compatibility
        errMsgCb(result.Data);
      }
      navigate('/failure', {
        state: { req },
        // TODO: pass errorsEN[result.Error] and display it in Failure component
      });
      if (quitOnError || result.Error === 'timeoutError') {
        setTimeout(hideAndReload, 2000);
        return;
      }
    }
  };
};
