import { Server } from 'miragejs';

import { AppSchema } from '../types';

export function routesForConfig(server: Server) {
  server.get('config', (schema: AppSchema) => {
    const { models: accounts } = schema.all('account');
    const { models: signRules } = schema.all('signRule');

    const accountConfigs = signRules.reduce(
      (acc: Record<string, any>, rule) => {
        const nickname = rule.accountNickname;
        if (!acc[nickname]) {
          acc[nickname] = {
            signRules: [],
          };
        }
        acc[nickname].signRules.push({
          ruleType: rule.ruleType,
          contract: rule.contract,
          enabled: rule.enabled,
          id: rule.id,
          name: rule.name,
        });
        return acc;
      },
      {},
    );

    accounts.forEach((account) => {
      if (!accountConfigs[account.nickname]) {
        accountConfigs[account.nickname] = {
          signRules: [],
        };
      }
    });

    return {
      accounts: accountConfigs,
    };
  });
}
