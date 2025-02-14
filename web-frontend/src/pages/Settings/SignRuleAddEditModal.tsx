import { useState } from 'react';

import {
  Button,
  Dropdown,
  Input,
  PopupModal,
  PopupModalContent,
  PopupModalHeader,
  Toggle,
} from '@massalabs/react-ui-kit';
import { RuleType, SignRule } from '@massalabs/wallet-provider';

import { useProvider } from '@/custom/useProvider';
import Intl from '@/i18n/i18n';

interface SignRuleModalProps {
  onClose: () => void;
  rule?: SignRule;
  nickname: string;
  onSuccess?: (successMessage: string) => void;
  onError?: (errorMessage: string) => void;
}

export function SignRuleModal(props: SignRuleModalProps) {
  const { onClose, rule, nickname, onSuccess, onError } = props;
  const { wallet } = useProvider();
  const [_isSubmitting, setIsSubmitting] = useState(false);

  const [name, setName] = useState(rule?.name || '');
  const [contract, setContract] = useState(rule?.contract || '');
  const [ruleType, setRuleType] = useState(
    rule?.ruleType || RuleType.DisablePasswordPrompt,
  );
  const [applyToAllContracts, setApplyToAllContracts] = useState(false);

  const isEditMode = !!rule;

  const handleSubmit = async () => {
    try {
      setIsSubmitting(true);
      const signRuleData: SignRule = {
        id: rule?.id || '',
        name,
        contract: applyToAllContracts ? '*' : contract,
        ruleType,
        enabled: rule?.enabled || true,
      };

      if (isEditMode) {
        await wallet?.editSignRule(
          nickname,
          signRuleData,
          `Update sign rule ${name}`,
        );
        onSuccess?.(Intl.t('settings.sign-rules.success.update'));
      } else {
        await wallet?.addSignRule(
          nickname,
          signRuleData,
          `Add sign rule ${name}`,
        );
        onSuccess?.(Intl.t('settings.sign-rules.success.add'));
      }
    } catch (error: any) {
      console.error('Error submitting sign rule:', error);
      let errorMsg = Intl.t(
        `settings.sign-rules.errors.${isEditMode ? 'update' : 'add'}`,
      );

      if (error?.message === 'Rule already exists') {
        errorMsg = Intl.t('settings.sign-rules.errors.already-exist');
      }

      onError?.(errorMsg);
    } finally {
      setIsSubmitting(false);
    }
  };

  const selectedRuleType = Object.values(RuleType).indexOf(
    ruleType as RuleType,
  );

  return (
    <PopupModal
      customClass="w-[580px] h-[400px]"
      fullMode={true}
      onClose={onClose}
    >
      <PopupModalHeader>
        <div className="mb-6">
          <p className="mas-title mb-6">
            {isEditMode
              ? Intl.t('settings.sign-rules.modals.edit-title')
              : Intl.t('settings.sign-rules.modals.add-title')}
          </p>
        </div>
      </PopupModalHeader>
      <PopupModalContent>
        <div className="mas-body2 pb-10 flex flex-col gap-4">
          <Input
            value={name}
            onChange={(e) => setName(e.target.value)}
            name="name"
            placeholder={Intl.t('settings.sign-rules.modals.name-placeholder')}
          />
          <Input
            value={contract}
            onChange={(e) => setContract(e.target.value)}
            name="contract"
            placeholder={
              applyToAllContracts
                ? 'All'
                : Intl.t('settings.sign-rules.modals.contract-placeholder')
            }
            disabled={applyToAllContracts}
          />
          <div className="flex items-center gap-2">
            <Toggle
              checked={applyToAllContracts}
              onChange={(e) => setApplyToAllContracts(e.target.checked)}
            />
            {Intl.t('settings.sign-rules.modals.apply-to-all-contracts')}
          </div>
          <Dropdown
            select={selectedRuleType}
            readOnly={false}
            size="md"
            options={Object.values(RuleType).map((type) => ({
              item:
                type === RuleType.AutoSign
                  ? Intl.t('settings.sign-rules.auto-sign')
                  : Intl.t('settings.sign-rules.disable-password-prompt'),
              onClick: () => setRuleType(type),
            }))}
          />
          <Button customClass="mt-6 self-start" onClick={handleSubmit}>
            {isEditMode
              ? Intl.t('settings.sign-rules.modals.update')
              : Intl.t('settings.sign-rules.modals.add')}
          </Button>
        </div>
      </PopupModalContent>
    </PopupModal>
  );
}
