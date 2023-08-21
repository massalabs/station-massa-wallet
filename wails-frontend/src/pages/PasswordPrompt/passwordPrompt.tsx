// import { useLocation } from 'react-router-dom';

// import { promptAction, promptRequest } from '@/events/events';

// import Intl from '@/i18n/i18n';

export interface PromptRequestDeleteData {
  Nickname: string;
  Balance: string;
}

function PasswordPrompt() {
  // const { state } = useLocation();
  // const req: promptRequest = state.req;
  // const { deleteReq, signReq, transferReq } = promptAction;

  //maybe directly initialize my useState with the fn return
  // useEffect(() => {
  //   getTitle();
  //   getSubtitle();
  //   getButtonIcon();
  //   getButtonLabel();
  // }, []);

  // function getTitle() {
  //   switch (req.Action) {
  //     case deleteReq:
  //       return setTitle(Intl.t('password-prompt.title.delete'));
  //     case signReq:
  //       return setTitle(Intl.t('password-prompt.title.sign'));
  //     default:
  //       return setTitle(Intl.t(`password-prompt.title.${req.CodeMessage}`));
  //   }
  // }

  // function getButtonLabel() {
  //   switch (req.Action) {
  //     case deleteReq:
  //       return setButtonLabel(Intl.t('password-prompt.buttons.delete'));
  //     case signReq:
  //       return setButtonLabel(Intl.t('password-prompt.buttons.sign'));
  //     case transferReq:
  //       return setButtonLabel(Intl.t('password-prompt.buttons.transfer'));
  //     default:
  //       return setButtonLabel(Intl.t('password-prompt.buttons.default'));
  //   }
  // }

  // function getSubtitle() {
  //   switch (req.Action) {
  //     case deleteReq:
  //       return setSubtitle(Intl.t('password-prompt.subtitle.delete'));
  //     case signReq:
  //       return setSubtitle(Intl.t('password-prompt.subtitle.sign'));
  //     default:
  //       return setSubtitle(Intl.t('password-prompt.subtitle.default'));
  //   }
  // }

  // function getButtonIcon() {
  //   switch (req.Action) {
  //     case deleteReq:
  //       return setButtonIcon(<FiTrash2 />);
  //     case signReq:
  //       return;
  //     case transferReq:
  //       return;
  //     default:
  //       return setButtonIcon(<FiLock />);
  //   }
  // }

  return (
    <>
      <p>Condionally Render</p>
      <p>Correct Prompt</p>
    </>
  );
}

export default PasswordPrompt;
