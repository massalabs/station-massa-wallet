/* eslint-disable new-cap */
import { EventsOn } from '../wailsjs/runtime';
import { events, promptAction, promptRequest } from './events/events';
import { useNavigate } from 'react-router-dom';

export function App() {
  const navigate = useNavigate();

  const handlePromptRequest = (req: promptRequest) => {
    switch (req.Action) {
      case promptAction.deleteReq:
      case promptAction.newPasswordReq:
      case promptAction.signReq:
      case promptAction.exportReq:
      case promptAction.transferReq:
      case promptAction.tradeRollsReq:
        navigate('/password', { state: { req } });
        return;
      case promptAction.importReq:
        navigate('/import-methods', { state: { req } });
        return;
      default:
    }
  };

  EventsOn(events.promptRequest, handlePromptRequest);

  return <></>;
}

export default App;
