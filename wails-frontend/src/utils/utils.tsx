import { NavigateFunction } from 'react-router-dom';
import { Hide, AbortAction } from '../../wailsjs/go/walletapp/WalletApp';
import { WindowReloadApp } from '../../wailsjs/runtime';
import { promptResult, promptRequest } from '../events/events';

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
      errMsgCb(result.Data);
      if (quitOnError || result.Error === 'timeoutError') {
        setTimeout(hideAndReload, 2000);
        return;
      }
    }
  };
};
