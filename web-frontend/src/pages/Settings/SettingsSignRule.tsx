import { useEffect, useState } from 'react';

import {
  Button,
  ButtonIcon,
  ButtonToggle,
  Clipboard,
  FetchingLine,
  maskAddress,
  toast,
  Tooltip,
} from '@massalabs/react-ui-kit';
import { Config, RuleType, SignRule } from '@massalabs/wallet-provider';
import { FiEdit, FiTrash2 } from 'react-icons/fi';

import { SignRuleModal } from './SignRuleAddEditModal';
import { useProvider } from '@/custom/useProvider';
import Intl from '@/i18n/i18n';

interface SettingsSignRulesProps {
  nickname: string;
}

export default function SettingsSignRules(props: SettingsSignRulesProps) {
  const { nickname } = props;

  const [editingRule, setEditingRule] = useState<SignRule | undefined>(
    undefined,
  );
  const [isAddEditRuleModalOpen, setIsAddEditRuleModalOpen] = useState(false);

  const [config, setConfig] = useState<Config | undefined>(undefined);
  const [isLoading, setIsLoading] = useState(true);

  const { wallet } = useProvider();

  const fetchConfig = async () => {
    try {
      const walletConfig = await wallet?.getConfig();
      setConfig(walletConfig);
    } catch (error) {
      console.error('Error fetching config:', error);
      toast.error('Error fetching config');
    } finally {
      setIsLoading(false);
    }
  };

  useEffect(() => {
    if (!isAddEditRuleModalOpen) {
      fetchConfig();
    }
  }, [wallet, isAddEditRuleModalOpen]);

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
          {isLoading ? (
            <FetchingLine />
          ) : (
            signRules.map((rule, index) => (
              <tr key={index} className="align-baseline">
                <SignRuleListItem
                  rule={rule}
                  nickname={nickname}
                  setEditingRule={handleEdit}
                  refreshConfig={fetchConfig}
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
  refreshConfig: () => void;
}

function SignRuleListItem(props: SignRuleListItemProps) {
  const { rule, nickname, setEditingRule, refreshConfig } = props;
  const { wallet } = useProvider();
  const [_isUpdating, setIsUpdating] = useState(false);
  const [_isDeleting, setIsDeleting] = useState(false);

  const isAllContract = rule.contract === '*';

  const handleToggle = async () => {
    try {
      setIsUpdating(true);
      if (!rule.id) {
        throw new Error('Rule ID is required');
      }
      await wallet?.editSignRule(
        nickname,
        {
          ...rule,
          enabled: !rule.enabled,
        },
        `Turn ${rule.enabled ? 'off' : 'on'} sign rule ${rule.name}`,
      );
      toast.success(Intl.t('settings.sign-rules.success.toggle'));
    } catch (error) {
      console.error('Error updating sign rule:', error);
      toast.error(Intl.t('settings.sign-rules.errors.toggle'));
    } finally {
      setIsUpdating(false);
      refreshConfig();
    }
  };

  const handleDelete = async () => {
    try {
      setIsDeleting(true);
      if (!rule.id) {
        throw new Error('Rule ID is required');
      }
      await wallet?.deleteSignRule(nickname, rule.id);
      toast.success(Intl.t('settings.sign-rules.success.delete'));
    } catch (error) {
      console.error('Error deleting sign rule:', error);
      toast.error(Intl.t('settings.sign-rules.errors.delete'));
    } finally {
      setIsDeleting(false);
      refreshConfig();
    }
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
        {isAllContract ? (
          <div className="flex justify-center">
            <span>All</span>
          </div>
        ) : (
          <Clipboard
            rawContent={rule.contract}
            displayedContent={maskAddress(rule.contract, 6)}
            className="text-sm p-2 flex"
          />
        )}
      </td>
      <td className="text-right">
        <div className="flex justify-end items-center space-x-2">
          <ButtonToggle onClick={handleToggle}>
            {rule.enabled ? 'On' : 'Off'}
          </ButtonToggle>
          <ButtonIcon variant="primary" onClick={() => setEditingRule(rule)}>
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
