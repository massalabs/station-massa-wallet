import { walletapp } from '@wailsjs/go/models';
import { useLocation } from 'react-router-dom';

import { Delete } from './Delete';
import { Sign } from './Sign';
import { SignRule } from '../SignRuleHandler/signRulePrompt';
import { promptRequest } from '@/events';

export interface PromptRequestDeleteData {
  Nickname: string;
  Balance: string;
}

export default function PasswordPromptHandler() {
  const { state } = useLocation();
  const req: promptRequest = state.req;

  const { PromptRequestAction } = walletapp;

  return (
    <>
      {(() => {
        switch (req.Action) {
          case PromptRequestAction.delete:
            return <Delete />;
          case PromptRequestAction.sign:
            return <Sign />;
          case PromptRequestAction.addSignRule:
          case PromptRequestAction.updateSignRule:
          case PromptRequestAction.deleteSignRule:
            return <SignRule />;
        }
      })()}
    </>
  );
}
