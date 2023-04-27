import { NavigateFunction, useNavigate } from 'react-router-dom';
import { Hide, AbortAction } from '../../wailsjs/go/walletapp/WalletApp';
import { WindowReloadApp } from '../../wailsjs/runtime';

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
  successMsg: string,
  errMsgCb: (msg: any) => void,
) => {
  return (result: any) => {
    if (result.Success) {
      navigate('/success', {
        state: { msg: successMsg },
      });
      setTimeout(hideAndReload, 2000);
    } else {
      errMsgCb(result.Data);
      if (result.Error === 'timeoutError') {
        setTimeout(hideAndReload, 2000);
      }
    }
  };
};
