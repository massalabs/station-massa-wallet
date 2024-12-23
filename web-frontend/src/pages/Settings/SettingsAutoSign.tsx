import { ButtonIcon, ButtonToggle, Clipboard } from '@massalabs/react-ui-kit';
import { FiTrash2 } from 'react-icons/fi';

import Intl from '@/i18n/i18n';
import { Config } from '@/models/ConfigModel';

interface SettingsAutoSignProps {
  nickname: string;
}

export default function SettingsAutoSign(props: SettingsAutoSignProps) {
  const { nickname } = props;

  // Dummy data
  const config: Config = {
    accounts: {
      [nickname]: {
        signRules: [
          {
            contract: 'AS12sDEt4fFFPmRF87fh4X4SSDRioejDoSppMA23DeJIaDO4QMbn',
            passwordPrompt: false,
            autoSign: true,
          },
          {
            contract: 'AS1EJisvbaE4dDizPoHD5SkOppLM4SCnIEscLPIUeAR9mRSOE4MOS',
            passwordPrompt: true,
            autoSign: false,
          },
        ],
      },
    },
  };

  // Dummy functions
  const handleToggle = (contract: string) => {
    console.log(`Toggled autoSign for contract: ${contract}`);
  };

  // Dummy functions
  const handleDelete = (contract: string) => {
    console.log(`Deleted contract: ${contract}`);
  };

  return (
    <>
      <p className="mas-body text-f-primary pb-5">
        {Intl.t('settings.title-auto-sign')}
      </p>
      <p className="mas-body2 text-f-primary pb-5">
        {Intl.t('settings.auto-sign.description')}
      </p>
      <table className="w-full">
        <thead>
          <tr>
            <th className="text-left">
              {Intl.t('settings.auto-sign.address')}
            </th>
            <th className="text-right">
              {Intl.t('settings.auto-sign.actions')}
            </th>
          </tr>
        </thead>
        <tbody>
          {config.accounts[nickname || 'unknown'].signRules.map(
            (rule, index) => (
              <tr key={index}>
                <td className="max-w-xs truncate whitespace-nowrap">
                  <Clipboard
                    rawContent={rule.contract}
                    className="text-sm p-2 flex"
                  />
                </td>
                <td className="text-right">
                  <div className="flex justify-end items-center space-x-2">
                    <ButtonToggle onClick={() => handleToggle(rule.contract)}>
                      {rule.autoSign ? 'On' : 'Off'}
                    </ButtonToggle>
                    <ButtonIcon
                      variant="primary"
                      onClick={() => handleDelete(rule.contract)}
                    >
                      <FiTrash2 />
                    </ButtonIcon>
                  </div>
                </td>
              </tr>
            ),
          )}
        </tbody>
      </table>
    </>
  );
}
