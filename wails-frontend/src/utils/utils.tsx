import { NavigateFunction } from 'react-router-dom';
import { Hide, AbortAction } from '../../wailsjs/go/walletapp/WalletApp';
import { WindowReloadApp } from '../../wailsjs/runtime';
import { promptResult, promptRequest } from '../events/events';
import { IErrorObject } from './errors';
import Intl from '../i18n/i18n';

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
  errMsgCb:
    | React.Dispatch<React.SetStateAction<IErrorObject | null>>
    | ((msg: string) => void),
  quitOnError = false,
) => {
  return (result: promptResult) => {
    if (result.Success) {
      navigate('/success', {
        state: { req },
      });
      setTimeout(hideAndReload, 2000);
    } else {
      req.Msg = Intl.t(`errors.${result.CodeMessage}`);
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
