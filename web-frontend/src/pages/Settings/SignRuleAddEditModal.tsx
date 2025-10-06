import { useState } from 'react';

import {
  Button,
  Dropdown,
  Input,
  PopupModal,
  PopupModalContent,
  PopupModalHeader,
  Toggle,
  Tooltip,
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
  const [applyToAllContracts, setApplyToAllContracts] = useState(
    rule?.contract === '*',
  );
  const [authorizedOrigin, setAuthorizedOrigin] = useState(
    rule?.authorizedOrigin || '',
  );
  const [applyToAllDomains, setApplyToAllDomains] = useState(
    !rule?.authorizedOrigin,
  );

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
        authorizedOrigin:
          ruleType === RuleType.AutoSign
            ? authorizedOrigin
            : applyToAllDomains
            ? undefined
            : authorizedOrigin,
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
      } else if (error?.message.includes('invalid contract address')) {
        errorMsg = Intl.t(
          'settings.sign-rules.errors.invalid-contract-address',
        );
      } else if (error?.message.includes('action canceled by user')) {
        errorMsg = Intl.t('errors.action-canceled');
      } else if (error?.message.includes('invalid AuthorizedOrigin URL')) {
        errorMsg = Intl.t(
          'settings.sign-rules.errors.invalid-authorized-origin',
          { error: error?.message },
        );
      }

      onError?.(errorMsg);
    } finally {
      setIsSubmitting(false);
    }
  };

  const selectSignRuleType = (type: RuleType) => {
    setRuleType(type);
    if (type === RuleType.AutoSign) {
      if (applyToAllContracts) setApplyToAllContracts(false);
      if (applyToAllDomains && !isEditMode) setApplyToAllDomains(false);
    }
  };

  const selectedRuleType = Object.values(RuleType).indexOf(
    ruleType as RuleType,
  );

  const shouldDisableSubmit =
    !name ||
    (!contract && !applyToAllContracts) ||
    (!authorizedOrigin && !applyToAllDomains);

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
          <div className="flex flex-col gap-1">
            <p className="mas-caption text-s-info">
              {Intl.t('settings.sign-rules.modals.name-placeholder')}
            </p>
            <Input
              value={name}
              onChange={(e) => setName(e.target.value)}
              name="name"
              placeholder={Intl.t(
                'settings.sign-rules.modals.name-placeholder',
              )}
              required
            />
          </div>
          <div className="flex flex-col gap-1">
            <p className="mas-caption text-s-info">
              {Intl.t('settings.sign-rules.modals.contract-placeholder')}
            </p>
            <Input
              value={applyToAllContracts ? 'All' : contract}
              onChange={(e) => setContract(e.target.value)}
              name="contract"
              placeholder={
                applyToAllContracts
                  ? 'All'
                  : Intl.t('settings.sign-rules.modals.contract-placeholder')
              }
              disabled={applyToAllContracts}
              required={!applyToAllContracts}
            />
          </div>
          <div className="flex items-center gap-2">
            {ruleType === RuleType.AutoSign ? (
              <Tooltip
                body={Intl.t(
                  'settings.sign-rules.modals.all-contracts-auto-sign',
                )}
                placement="top"
                triggerClassName="truncate"
                tooltipClassName="mas-caption max-w-96"
              >
                <div className="flex items-center gap-2 opacity-50 cursor-not-allowed">
                  <Toggle
                    checked={applyToAllContracts}
                    onChange={(e) => setApplyToAllContracts(e.target.checked)}
                    disabled
                  />
                  {Intl.t('settings.sign-rules.modals.apply-to-all-contracts')}
                </div>
              </Tooltip>
            ) : (
              <>
                <Toggle
                  checked={applyToAllContracts}
                  onChange={(e) => setApplyToAllContracts(e.target.checked)}
                />
                {Intl.t('settings.sign-rules.modals.apply-to-all-contracts')}
              </>
            )}
          </div>
          <Dropdown
            select={selectedRuleType}
            readOnly={false}
            size="md"
            options={Object.values(RuleType).map((type) => ({
              item:
                type === RuleType.AutoSign ? (
                  <Tooltip
                    body={Intl.t('settings.sign-rules.auto-sign-tooltip')}
                    placement="top"
                    triggerClassName="truncate w-full"
                    tooltipClassName="mas-caption max-w-96"
                  >
                    {Intl.t('settings.sign-rules.auto-sign')}
                  </Tooltip>
                ) : (
                  <Tooltip
                    body={Intl.t(
                      'settings.sign-rules.disable-password-prompt-tooltip',
                    )}
                    placement="top"
                    triggerClassName="truncate w-full"
                    tooltipClassName="mas-caption max-w-96"
                  >
                    {Intl.t('settings.sign-rules.disable-password-prompt')}
                  </Tooltip>
                ),
              onClick: () => selectSignRuleType(type),
            }))}
          />
          <>
            <div className="flex flex-col gap-1">
              <div className="flex items-center gap-1">
                <p className="mas-caption text-s-info">
                  {Intl.t(
                    'settings.sign-rules.modals.authorized-origin-placeholder',
                  )}
                </p>
                <Tooltip
                  body={Intl.t(
                    'settings.sign-rules.modals.authorized-origin-info',
                  )}
                  placement="top"
                  triggerClassName="inline-flex items-center justify-center w-4 h-4
                  rounded-full border border-s-info text-s-info cursor-help"
                  tooltipClassName="mas-caption max-w-96"
                >
                  <span className="mas-caption leading-none">?</span>
                </Tooltip>
              </div>
              <div className={`relative ${isEditMode ? 'opacity-60' : ''}`}>
                <Input
                  value={applyToAllDomains ? 'All domains' : authorizedOrigin}
                  onChange={(e) => setAuthorizedOrigin(e.target.value)}
                  name="authorizedOrigin"
                  placeholder={
                    applyToAllDomains
                      ? 'All domains'
                      : Intl.t(
                          'settings.sign-rules.modals.authorized-origin-placeholder',
                        )
                  }
                  required={!applyToAllDomains}
                  disabled={applyToAllDomains || isEditMode}
                />
                {isEditMode && (
                  <Tooltip
                    body={Intl.t(
                      'settings.sign-rules.modals.authorized-origin-disabled-message',
                    )}
                    placement="top"
                    triggerClassName="absolute right-2 top-1/2 -translate-y-1/2 inline-flex items-center 
                    justify-center w-4 h-4 rounded-full border border-s-info text-s-info cursor-help"
                    tooltipClassName="mas-caption max-w-96"
                  >
                    <span className="mas-caption leading-none">🔒</span>
                  </Tooltip>
                )}
              </div>
            </div>
            <div className="flex items-center gap-2">
              {ruleType === RuleType.AutoSign ? (
                <Tooltip
                  body={Intl.t(
                    'settings.sign-rules.modals.all-authorized-origin-auto-sign',
                  )}
                  placement="top"
                  triggerClassName="truncate"
                  tooltipClassName="mas-caption max-w-96"
                >
                  <div className="flex items-center gap-2 opacity-50 cursor-not-allowed">
                    <Toggle
                      checked={applyToAllDomains}
                      onChange={(e) => setApplyToAllDomains(e.target.checked)}
                      disabled
                    />
                    {Intl.t('settings.sign-rules.modals.apply-to-all-domains')}
                  </div>
                </Tooltip>
              ) : (
                <>
                  <Toggle
                    checked={applyToAllDomains}
                    onChange={(e) => setApplyToAllDomains(e.target.checked)}
                    disabled={isEditMode}
                  />
                  {Intl.t('settings.sign-rules.modals.apply-to-all-domains')}
                </>
              )}
            </div>
            {isEditMode && (
              <p className="mas-caption text-s-info">
                {Intl.t(
                  'settings.sign-rules.modals.authorized-origin-disabled-message',
                )}
              </p>
            )}
          </>
          <Button
            customClass="self-start"
            onClick={handleSubmit}
            disabled={shouldDisableSubmit}
          >
            {isEditMode
              ? Intl.t('settings.sign-rules.modals.update')
              : Intl.t('settings.sign-rules.modals.add')}
          </Button>
        </div>
      </PopupModalContent>
    </PopupModal>
  );
}
