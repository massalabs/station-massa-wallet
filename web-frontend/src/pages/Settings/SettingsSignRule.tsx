import { useEffect, useState } from 'react';

import {
  Button,
  ButtonIcon,
  ButtonToggle,
  Clipboard,
  FetchingLine,
  toast,
  Tooltip,
} from '@massalabs/react-ui-kit';
import { maskAddress } from '@massalabs/react-ui-kit/src/lib/massa-react/utils';
import { FiEdit, FiTrash2 } from 'react-icons/fi';

import { SignRuleModal } from './SignRuleAddEditModal';
import { useDelete, usePut, useResource } from '@/custom/api';
import Intl from '@/i18n/i18n';
import { Config, RuleType, SignRule } from '@/models/ConfigModel';

interface SettingsSignRulesProps {
  nickname: string;
}

export default function SettingsSignRules(props: SettingsSignRulesProps) {
  const { nickname } = props;

  const [editingRule, setEditingRule] = useState<SignRule | undefined>(
    undefined,
  );
  const [isAddEditRuleModalOpen, setIsAddEditRuleModalOpen] = useState(false);

  const {
    data: config,
    isLoading: isConfigLoading,
    error: configError,
  } = useResource<Config>('config', true);

  if (configError) {
    console.error('Error fetching config:', configError);
    toast.error('Error fetching config');
  }

  const signRules = config?.accounts?.[nickname]?.signRules ?? [];

  const handleAdd = () => {
    setIsAddEditRuleModalOpen(true);
  };

  const handleEdit = (rule: SignRule) => {
    setEditingRule(rule);
    setIsAddEditRuleModalOpen(true);
  };

  const handleAddEditModalClose = () => {
    setEditingRule(undefined);
    setIsAddEditRuleModalOpen(false);
  };

  const handleAddEditModalSuccess = (message: string) => {
    toast.success(message);
    handleAddEditModalClose();
  };

  const handleAddEditModalError = (message: string) => {
    toast.error(message);
  };

  return (
    <>
      <div className="w-full flex items-center justify-between">
        <p className="mas-body text-f-primary">
          {Intl.t('settings.title-sign-rules')}
        </p>
        <Button customClass="w-fit max-h-[2.5rem]" onClick={handleAdd}>
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
              <tr key={index} className="align-baseline">
                <SignRuleListItem
                  rule={rule}
                  nickname={nickname}
                  setEditingRule={handleEdit}
                />
              </tr>
            ))
          )}
        </tbody>
      </table>
      {isAddEditRuleModalOpen && (
        <SignRuleModal
          rule={editingRule}
          nickname={nickname}
          onClose={handleAddEditModalClose}
          onSuccess={handleAddEditModalSuccess}
          onError={handleAddEditModalError}
        />
      )}
    </>
  );
}

interface SignRuleListItemProps {
  setEditingRule: (rule: SignRule) => void;
  rule: SignRule;
  nickname: string;
}

function SignRuleListItem(props: SignRuleListItemProps) {
  const { rule, nickname, setEditingRule } = props;

  const {
    mutate: updateSignRule,
    isSuccess: updateSignRuleSuccess,
    error: updateSignRuleError,
    reset: resetUpdateError,
  } = usePut(`accounts/${nickname}/signrules/${rule.id}`);

  const {
    mutate: deleteSignRule,
    isSuccess: deleteSignRuleSuccess,
    error: deleteSignRuleError,
    reset: resetDeleteError,
  } = useDelete(`accounts/${nickname}/signrules/${rule.id}`);

  useEffect(() => {
    if (updateSignRuleError) {
      console.error('Error updating sign rule:', updateSignRuleError);
      toast.error(Intl.t('settings.sign-rules.errors.toggle'));
      resetUpdateError();
    }
  }, [updateSignRuleError, resetUpdateError]);

  useEffect(() => {
    if (deleteSignRuleError) {
      console.error('Error deleting sign rule:', deleteSignRuleError);
      toast.error(Intl.t('settings.sign-rules.errors.delete'));
      resetDeleteError();
    }
  }, [deleteSignRuleError, resetDeleteError]);

  useEffect(() => {
    if (updateSignRuleSuccess) {
      toast.success(Intl.t('settings.sign-rules.success.toggle'));
    } else if (deleteSignRuleSuccess) {
      toast.success(Intl.t('settings.sign-rules.success.delete'));
    }
  }, [updateSignRuleSuccess, deleteSignRuleSuccess]);

  const handleEdit = () => {
    setEditingRule(rule);
  };

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
      <td className="max-w-[128px] truncate">
        {rule.name && rule.name.length > 0 && (
          <Tooltip
            body={rule.name}
            placement="right"
            triggerClassName="truncate w-full"
            tooltipClassName="mas-caption"
          >
            {rule.name}
          </Tooltip>
        )}
      </td>
      <td className="max-w-xs truncate whitespace-nowrap">
        {rule.ruleType === RuleType.AutoSign ? (
          <Tooltip
            body={Intl.t('settings.sign-rules.auto-sign-tooltip')}
            placement="top"
            triggerClassName="truncate w-full"
            tooltipClassName="mas-caption"
          >
            <span>{Intl.t('settings.sign-rules.auto-sign')}</span>
          </Tooltip>
        ) : (
          <Tooltip
            body={Intl.t('settings.sign-rules.disable-password-prompt-tooltip')}
            placement="top"
            triggerClassName="truncate w-full"
            tooltipClassName="mas-caption"
          >
            <span>{Intl.t('settings.sign-rules.disable-password-prompt')}</span>
          </Tooltip>
        )}
      </td>
      <td className="max-w-xs truncate whitespace-nowrap">
        <Clipboard
          rawContent={rule.contract}
          displayedContent={maskAddress(rule.contract, 6)}
          className="text-sm p-2 flex"
        />
      </td>
      <td className="text-right">
        <div className="flex justify-end items-center space-x-2">
          <ButtonToggle onClick={handleToggle}>
            {rule.enabled ? 'On' : 'Off'}
          </ButtonToggle>
          <ButtonIcon variant="primary" onClick={handleEdit}>
            <FiEdit />
          </ButtonIcon>
          <ButtonIcon variant="primary" onClick={handleDelete}>
            <FiTrash2 />
          </ButtonIcon>
        </div>
      </td>
    </>
  );
}
