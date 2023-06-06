import { NavigateFunction } from 'react-router-dom';
import { Hide, AbortAction } from '../../wailsjs/go/walletapp/WalletApp';
import { WindowReloadApp } from '../../wailsjs/runtime';
import { promptResult, promptRequest } from '../events/events';
import { IErrorObject } from './errors';
import Intl from '../i18n/i18n';
import { useConfigStore } from '../store/store';

export const handleCancel = () => {
  AbortAction();
  hideAndReload();
};

export const hideAndReload = () => {
  Hide();
  // Reload the wails frontend
  WindowReloadApp();
};

const timeoutDelay = 1000 * 6; // 6 seconds

export const handleApplyResult = (
  navigate: NavigateFunction,
  req: promptRequest,
  errMsgCb:
    | React.Dispatch<React.SetStateAction<IErrorObject | null>>
    | ((msg: string) => void)
    | null = null,
  quitOnError = false,
) => {
  return (result: promptResult) => {
    let timeoutId;
    if (result.Success) {
      navigate('/success', {
        state: { req },
      });
      timeoutId = setTimeout(hideAndReload, timeoutDelay);
    } else {
      req.Msg = Intl.t(`errors.${result.CodeMessage}`);
      if (errMsgCb) errMsgCb(req.Msg);
      navigate('/failure', {
        state: { req },
      });
      if (quitOnError || result.CodeMessage === 'Timeout-0001') {
        timeoutId = setTimeout(hideAndReload, timeoutDelay);
      }
    }
    useConfigStore.setState({ timeoutId: Number(timeoutId) });
  };
};
