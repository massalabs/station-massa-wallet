import { EventsOn } from '@wailsjs/runtime/runtime';
import { useNavigate } from 'react-router-dom';

import {
  events,
  promptAction,
  promptRequest,
  promptResult,
} from './events/events';
import Intl from './i18n/i18n';
import { Loading } from './pages/loading';
import { useConfigStore } from './store/store';
import { ErrorCode } from './utils';

export function App() {
  const navigate = useNavigate();

  const handlePromptRequest = (req: promptRequest) => {
    clearTimeout(useConfigStore.getState().timeoutId);

    switch (req.Action) {
      case promptAction.deleteReq:
      case promptAction.signReq:
      case promptAction.tradeRollsReq:
      case promptAction.unprotectReq:
        navigate('/password', { state: { req } });
        return;
      case promptAction.newPasswordReq:
        navigate('/new-password', { state: { req } });
        return;
      case promptAction.importReq:
        navigate('/import-methods', { state: { req } });
        return;
      case promptAction.backupReq:
        navigate('/backup-methods', { state: { req } });
        return;
      default:
    }
  };

  EventsOn(events.promptResult, (result: promptResult) => {
    if (!result.Success && result.CodeMessage === ErrorCode.Timeout) {
      const errorMessage = Intl.t(`errors.${result.CodeMessage}`);
      navigate('/failure', {
        state: { req: { Msg: errorMessage } },
      });
    }
  });

  EventsOn(events.promptRequest, handlePromptRequest);

  return <Loading />;
}

export default App;
