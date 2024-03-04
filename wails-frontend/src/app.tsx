import { EventsOn } from '@wailsjs/runtime/runtime';
import { useNavigate } from 'react-router-dom';

import {
  events,
  promptAction,
  promptRequest,
  promptResult,
} from './events/events';
import { useConfigStore } from './store/store';
import { hideAndReload } from './utils';

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
    // TODO: redirect to a page that display the error
    if (!result.Success && result.CodeMessage === 'Timeout-0001') {
      hideAndReload();
    }
    // TODO: move the string into constant, use proper error code
    if (
      !result.Success &&
      ['InvalidInputType-0001', 'WrongPromptCorrelationId-0001'].includes(
        result.CodeMessage,
      )
    ) {
      hideAndReload();
    }
  });

  EventsOn(events.promptRequest, handlePromptRequest);

  // TODO: display a prettier loading screen
  return <>Loading...</>;
}

export default App;
