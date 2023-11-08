import { SignBodyProps } from '@/pages/PasswordPromptHandler/Sign';
import { Description } from '@/pages/PasswordPromptHandler/SignComponentUtils/Description';
import { From } from '@/pages/PasswordPromptHandler/SignComponentUtils/From';

export function PlainText(props: SignBodyProps) {
  const {
    PlainText,
    DisplayData,
    WalletAddress,
    Description: description,
    Nickname,
  } = props;

  return (
    <div className="flex flex-col items-center gap-4 mas-menu-default w-[326px]">
      <From nickname={Nickname} walletAddress={WalletAddress} />

      <Description text={description} />

      {DisplayData && (
        <Description text={PlainText} label="password-prompt.sign.message" />
      )}
    </div>
  );
}
