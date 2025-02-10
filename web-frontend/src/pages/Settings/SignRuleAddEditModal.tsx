import { useEffect, useState } from 'react';

import {
  Button,
  Dropdown,
  Input,
  PopupModal,
  PopupModalContent,
  PopupModalHeader,
} from '@massalabs/react-ui-kit';

import { usePost, usePut } from '@/custom/api';
import Intl from '@/i18n/i18n';
import { RuleType, SignRule } from '@/models/ConfigModel';

interface SignRuleModalProps {
  onClose: () => void;
  rule?: SignRule;
  nickname: string;
  onSuccess?: (successMessage: string) => void;
  onError?: (errorMessage: string) => void;
}

export function SignRuleModal(props: SignRuleModalProps) {
  const { onClose, rule, nickname, onSuccess, onError } = props;

  const [name, setName] = useState(rule?.name || '');
  const [contract, setContract] = useState(rule?.contract || '');
  const [ruleType, setRuleType] = useState(
    rule?.ruleType || RuleType.DisablePasswordPrompt,
  );

  const isEditMode = !!rule;

  const {
    mutate: addSignRule,
    isSuccess: addSignRuleSuccess,
    error: addSignRuleError,
    reset: resetAddSignRuleError,
  } = usePost(`accounts/${nickname}/signrules`);

  const {
    mutate: updateSignRule,
    isSuccess: updateSignRuleSuccess,
    error: updateSignRuleError,
    reset: resetUpdateSignRuleError,
  } = usePut(`accounts/${nickname}/signrules/${rule?.id}`);

  useEffect(() => {
    if (addSignRuleError) {
      console.log('addSignRuleError', addSignRuleError);
      onError?.(Intl.t('settings.sign-rules.errors.add'));
      resetAddSignRuleError();
    } else if (updateSignRuleError) {
      console.log('updateSignRuleError', updateSignRuleError);
      onError?.(Intl.t('settings.sign-rules.errors.update'));
      resetUpdateSignRuleError();
    }
  }, [addSignRuleError, updateSignRuleError]);

  useEffect(() => {
    if (addSignRuleSuccess) {
      onSuccess?.(Intl.t('settings.sign-rules.success.add'));
    } else if (updateSignRuleSuccess) {
      onSuccess?.(Intl.t('settings.sign-rules.success.update'));
    }
  }, [addSignRuleSuccess, updateSignRuleSuccess]);

  const handleSubmit = () => {
    const signRuleData: SignRule = {
      id: rule?.id || '',
      name,
      contract,
      ruleType,
      enabled: rule?.enabled || true,
    };

    if (isEditMode) {
      updateSignRule(signRuleData);
    } else {
      addSignRule(signRuleData);
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
        <div className="mas-body2 pb-10">
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
            placeholder={Intl.t(
              'settings.sign-rules.modals.contract-placeholder',
            )}
          />
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
          <Button customClass="mt-6" onClick={handleSubmit}>
            {isEditMode
              ? Intl.t('settings.sign-rules.modals.update')
              : Intl.t('settings.sign-rules.modals.add')}
          </Button>
        </div>
      </PopupModalContent>
    </PopupModal>
  );
}
