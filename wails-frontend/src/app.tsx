import { useNavigate } from 'react-router-dom';

import { events, promptAction, promptRequest } from './events/events';
import { useConfigStore } from './store/store';
import { EventsOn } from '../wailsjs/runtime';

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

  EventsOn(events.promptRequest, handlePromptRequest);

  return <></>;
}

export default App;
