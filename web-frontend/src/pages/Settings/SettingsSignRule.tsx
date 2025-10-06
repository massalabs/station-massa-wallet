import { useCallback, useEffect, useState } from 'react';

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
import { RuleType } from '@massalabs/wallet-provider';
import { FcExpired } from 'react-icons/fc';
import { FiEdit, FiTrash2 } from 'react-icons/fi';

import { SignRuleModal } from './SignRuleAddEditModal';
import { useProvider } from '@/custom/useProvider';
import Intl from '@/i18n/i18n';
import { SignRule, Config } from '@/models/ConfigModel';

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

  const fetchConfig = useCallback(async () => {
    try {
      const walletConfig = await wallet?.getConfig();
      setConfig(walletConfig as Config);
    } catch (error) {
      console.error('Error fetching config:', error);
      toast.error('Error fetching config');
    } finally {
      setIsLoading(false);
    }
  }, [wallet]);

  useEffect(() => {
    if (!isAddEditRuleModalOpen) {
      fetchConfig();
    }
  }, [fetchConfig, isAddEditRuleModalOpen]);

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
      <div className="w-full flex items-center justify-between mb-4">
        <div>
          <h2 className="mas-h2 text-f-primary mb-2">
            {Intl.t('settings.title-sign-rules')}
          </h2>
          <p className="mas-body2 text-f-primary">
            {Intl.t('settings.sign-rules.description')}
          </p>
        </div>
        <Button
          customClass="w-fit max-h-[2.5rem] min-h-[2.5rem]"
          onClick={handleAdd}
        >
          {Intl.t('settings.sign-rules.add')}
        </Button>
      </div>
      <div className="w-full overflow-x-auto">
        <table className="w-full border-collapse min-w-[800px]">
          <thead>
            <tr className="border-b border-neutral-100 dark:border-neutral-400">
              <th className="text-left py-3 px-4">
                {Intl.t('settings.sign-rules.name')}
              </th>
              <th className="text-left py-3 px-4">
                {Intl.t('settings.sign-rules.rule-type')}
              </th>
              <th className="text-left py-3 px-4">
                {Intl.t('settings.sign-rules.contract-address')}
              </th>
              <th className="text-left py-3 px-4">
                {Intl.t('settings.sign-rules.authorized-origin')}
              </th>
              <th className="text-center py-3 px-4">
                {Intl.t('settings.sign-rules.actions')}
              </th>
            </tr>
          </thead>
          <tbody>
            {isLoading ? (
              <FetchingLine />
            ) : (
              signRules.map((rule, index) => (
                <tr
                  key={index}
                  className={`
                    align-baseline 
                    hover:bg-neutral-50 
                    dark:hover:bg-neutral-800 
                    transition-colors 
                    border-b-[1px]
                    border-neutral-100 
                    dark:border-neutral-400 
                    last:border-b-0
                  `}
                >
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
      </div>
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

  // Check if the rule has expired
  const isExpired = rule.expireAfter
    ? new Date(rule.expireAfter) < new Date()
    : false;

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
      <td className="max-w-[128px] truncate py-3 px-4">
        <div className="flex items-center gap-1 truncate">
          {isExpired && (
            <Tooltip
              body={Intl.t('settings.sign-rules.expired-tooltip')}
              placement="top"
              tooltipClassName="mas-caption max-w-96"
            >
              <FcExpired className="mr-1" />
            </Tooltip>
          )}
          {rule.name && rule.name.length > 0 ? (
            <Tooltip
              body={rule.name}
              placement="right"
              triggerClassName="truncate w-full"
              tooltipClassName="mas-caption"
            >
              {rule.name}
            </Tooltip>
          ) : (
            <div className="flex justify-center">
              <span className="text-f-primary">-</span>
            </div>
          )}
        </div>
      </td>
      <td className="max-w-xs truncate whitespace-nowrap py-3 px-4">
        {rule.ruleType === RuleType.AutoSign ? (
          <Tooltip
            body={Intl.t('settings.sign-rules.auto-sign-tooltip')}
            placement="top"
            triggerClassName="truncate w-full"
            tooltipClassName="mas-caption max-w-96 cursor-default"
          >
            <span
              className={`
              inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium
              bg-primary-100 text-primary-800
              dark:bg-primary-800 dark:text-primary-100
              border border-primary-200 dark:border-primary-700
            `}
            >
              {Intl.t('settings.sign-rules.auto-sign-short')}
            </span>
          </Tooltip>
        ) : (
          <Tooltip
            body={Intl.t('settings.sign-rules.disable-password-prompt-tooltip')}
            placement="top"
            triggerClassName="truncate w-full"
            tooltipClassName="mas-caption max-w-96"
          >
            <span
              className={`
              inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium
              bg-warning-100 text-warning-800
              dark:bg-warning-800 dark:text-warning-100
              border border-warning-200 dark:border-warning-700
              cursor-default
            `}
            >
              {Intl.t('settings.sign-rules.disable-password-prompt-short')}
            </span>
          </Tooltip>
        )}
      </td>
      <td className="max-w-xs truncate whitespace-nowrap py-3 px-4">
        {isAllContract ? (
          <div className="flex justify-center">
            <span className="text-f-primary">All</span>
          </div>
        ) : (
          <Clipboard
            rawContent={rule.contract}
            displayedContent={maskAddress(rule.contract, 10)}
            customClass="text-sm p-2 flex cursor-pointer"
          />
        )}
      </td>
      <td className="max-w-[200px] truncate whitespace-nowrap py-3 px-4">
        {rule.authorizedOrigin ? (
          <Tooltip
            body={rule.authorizedOrigin}
            placement="top"
            triggerClassName="truncate w-full"
            tooltipClassName="mas-caption max-w-96"
          >
            <span className="truncate">{rule.authorizedOrigin}</span>
          </Tooltip>
        ) : (
          <div className="flex justify-center">
            <span className="text-f-primary">All domains</span>
          </div>
        )}
      </td>
      <td className="text-right py-3 px-4">
        <div className="flex justify-end items-center space-x-2">
          <ButtonToggle
            onClick={handleToggle}
            customClass={`min-w-[60px] ${
              rule.enabled
                ? 'bg-primary-100 text-primary-800'
                : 'bg-warning-100 text-warning-800'
            }`}
          >
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
