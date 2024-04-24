import Intl from '@/i18n/i18n';

interface AmountBoxProps {
  children?: React.ReactNode;
}

export function AmountBox({ children }: AmountBoxProps) {
  return (
    <div className="bg-secondary text-center w-full rounded-lg p-4">
      <p className="mas-menu-default mb-4">
        {Intl.t('password-prompt.sign.amount')}
      </p>
      <p className="mas-menu-active">{children}</p>
    </div>
  );
}
