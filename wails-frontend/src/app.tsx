import { EventsOn } from '../wailsjs/runtime';
import { events, promptAction, promptRequest } from './events/events';
import { useNavigate } from 'react-router-dom';

export function App() {
  const navigate = useNavigate();

  const handlePromptRequest = (req: promptRequest) => {
    switch (req.Action) {
      case promptAction.deleteReq:
      case promptAction.signReq:
      case promptAction.transferReq:
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
