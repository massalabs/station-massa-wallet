import { Hide, AbortAction } from '@wailsjs/go/walletapp/WalletApp';
import { WindowReloadApp } from '@wailsjs/runtime/runtime';
import { NavigateFunction } from 'react-router-dom';

import { IErrorObject } from './errors';
import { promptResult, promptRequest } from '@/events';
import Intl from '@/i18n/i18n';
import { useConfigStore } from '@/store/store';

export const handleCancel = () => {
  AbortAction();
  hideAndReload();
};

export const hideAndReload = () => {
  Hide();
  // Reload the wails frontend
  WindowReloadApp();
};

const timeoutDelay = 1000 * 4; // 4 seconds

export const handleApplyResult = (
  navigate: NavigateFunction,
  req: promptRequest,
  errMsgCb:
    | React.Dispatch<React.SetStateAction<IErrorObject | null>>
    | ((msg: string) => void)
    | null = null,
  quitOnError = false,
) => {
  return (result: promptResult<null>) => {
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
      if (quitOnError) {
        timeoutId = setTimeout(hideAndReload, timeoutDelay);
      }
    }
    useConfigStore.setState({ timeoutId: Number(timeoutId) });
  };
};
