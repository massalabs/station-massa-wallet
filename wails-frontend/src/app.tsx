import { EventsOn } from '@wailsjs/runtime/runtime';
import { useNavigate } from 'react-router-dom';

import { promptRequest, promptResult } from './events';
import Intl from './i18n/i18n';
import { Loading } from './pages/loading';
import { useConfigStore } from './store/store';
import { ErrorCode } from './utils';
import { walletapp } from '../wailsjs/go/models';

export function App() {
  const navigate = useNavigate();

  const { PromptRequestAction, EventType } = walletapp;

  const handlePromptRequest = (req: promptRequest) => {
    clearTimeout(useConfigStore.getState().timeoutId);

    switch (req.Action) {
      case PromptRequestAction.delete:
      case PromptRequestAction.sign:
      case PromptRequestAction.tradeRolls:
      case PromptRequestAction.unprotect:
      case PromptRequestAction.addSignRule:
      case PromptRequestAction.deleteSignRule:
      case PromptRequestAction.updateSignRule:
        navigate('/password', { state: { req } });
        return;
      case PromptRequestAction.newPassword:
        navigate('/new-password', { state: { req } });
        return;
      case PromptRequestAction.import:
        navigate('/import-methods', { state: { req } });
        return;
      case PromptRequestAction.backup:
        navigate('/backup-methods', { state: { req } });
        return;
      default:
    }
  };

  EventsOn(EventType.promptResult, (result: promptResult) => {
    if (!result.Success && result.CodeMessage === ErrorCode.Timeout) {
      const errorMessage = Intl.t(`errors.${result.CodeMessage}`);
      navigate('/failure', {
        state: { req: { Msg: errorMessage } },
      });
    }
  });

  EventsOn(EventType.promptRequest, handlePromptRequest);

  return <Loading />;
}

export default App;
