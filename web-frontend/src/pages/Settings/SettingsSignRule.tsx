import {
  Button,
  ButtonIcon,
  ButtonToggle,
  Clipboard,
  FetchingLine,
  toast,
} from '@massalabs/react-ui-kit';
import { FiTrash2 } from 'react-icons/fi';

import { useDelete, usePost, usePut, useResource } from '@/custom/api';
import Intl from '@/i18n/i18n';
import { Config, RuleType, SignRule } from '@/models/ConfigModel';

interface SettingsSignRulesProps {
  nickname: string;
}

export default function SettingsSignRules(props: SettingsSignRulesProps) {
  const { nickname } = props;

  const {
    data: config,
    isLoading: isConfigLoading,
    error: configError,
  } = useResource<Config>('config', true);

  const { mutate: addSignRule, error: addSignRuleError } = usePost(
    `accounts/${nickname}/signrules`,
  );

  if (configError) {
    console.error('Error fetching config:', configError);
    toast.error('Error fetching config');
  }

  if (addSignRuleError) {
    console.error('Error adding sign rule:', addSignRuleError);
    toast.error('Error adding sign rule');
  }

  const signRules = config?.accounts?.[nickname]?.signRules ?? [];

  // TODO: move to a proper modal
  const handleAdd = () => {
    addSignRule({
      ruleType: RuleType.AutoSign,
      name: 'New Rule',
      contract: 'AS1BsB34Hq7VGpGG26cudgFr15nshSeJfAkkf6WGY8JcXbG6kzUg',
      enabled: true,
    });
  };

  return (
    <>
      <div className="w-full flex items-center justify-between">
        <p className="mas-body text-f-primary">
          {Intl.t('settings.title-sign-rules')}
        </p>
        <Button customClass="w-fit" onClick={handleAdd}>
          {Intl.t('settings.sign-rules.add')}
        </Button>
      </div>
      <p className="mas-body2 text-f-primary pb-5">
        {Intl.t('settings.sign-rules.description')}
      </p>
      <table className="w-full">
        <thead>
          <tr>
            <th className="text-left">{Intl.t('settings.sign-rules.name')}</th>
            <th className="text-left">
              {Intl.t('settings.sign-rules.rule-type')}
            </th>
            <th className="text-left">
              {Intl.t('settings.sign-rules.contract-address')}
            </th>
            <th className="text-right">
              {Intl.t('settings.sign-rules.actions')}
            </th>
          </tr>
        </thead>
        <tbody>
          {isConfigLoading ? (
            <FetchingLine />
          ) : (
            signRules.map((rule, index) => (
              <tr key={index}>
                <SignRuleListItem rule={rule} nickname={nickname} />
              </tr>
            ))
          )}
        </tbody>
      </table>
    </>
  );
}

interface SignRuleListItemProps {
  rule: SignRule;
  nickname: string;
}

function SignRuleListItem(props: SignRuleListItemProps) {
  const { rule, nickname } = props;

  const { mutate: updateSignRule, error: updateSignRuleError } = usePut(
    `accounts/${nickname}/signrules/${rule.id}`,
  );

  if (updateSignRuleError) {
    toast.error(Intl.t('assets.delete.bad-request'));
  }

  const { mutate: deleteSignRule, error: deleteSignRuleError } = useDelete(
    `accounts/${nickname}/signrules/${rule.id}`,
  );

  if (deleteSignRuleError) {
    toast.error(Intl.t('settings.signRules.deleteModal.error'));
    return;
  }

  const handleToggle = () => {
    updateSignRule({
      ...rule,
      enabled: !rule.enabled,
    });
  };

  const handleDelete = () => {
    deleteSignRule({});
  };

  return (
    <>
      <td className="max-w-xs truncate whitespace-nowrap">{rule.name}</td>
      <td className="max-w-xs truncate whitespace-nowrap">
        <div className="flex items-center gap-2 relative">
          {rule.ruleType === RuleType.AutoSign ? (
            <span>{Intl.t('settings.sign-rules.auto-sign')}</span>
          ) : (
            <span>{Intl.t('settings.sign-rules.disable-password-prompt')}</span>
          )}
        </div>
      </td>
      <td className="max-w-xs truncate whitespace-nowrap">
        <Clipboard rawContent={rule.contract} className="text-sm p-2 flex" />
      </td>
      <td className="text-right">
        <div className="flex justify-end items-center space-x-2">
          <ButtonToggle onClick={handleToggle}>
            {rule.enabled ? 'On' : 'Off'}
          </ButtonToggle>
          <ButtonIcon variant="primary" onClick={handleDelete}>
            <FiTrash2 />
          </ButtonIcon>
        </div>
      </td>
    </>
  );
}
