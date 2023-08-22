import { useLocation } from 'react-router-dom';

import { Default } from './Default';
import { Delete } from './Delete';
import { Sign } from './Sign';
import { Transfer } from './Transfer';
import { promptAction, promptRequest } from '@/events/events';

export interface PromptRequestDeleteData {
  Nickname: string;
  Balance: string;
}

function PasswordPromptHandler() {
  const { state } = useLocation();
  const req: promptRequest = state.req;
  const { deleteReq, signReq, transferReq } = promptAction;

  return (
    <>
      {(() => {
        switch (req.Action) {
          case deleteReq:
            return <Delete />;
          case signReq:
            return <Sign />;
          case transferReq:
            return <Transfer />;
          default:
            return <Default />;
        }
      })()}
    </>
  );
}
export default PasswordPromptHandler;
