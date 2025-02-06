import { config } from '@wailsjs/go/models';

// Used in i18n
export enum signRuleActionStr {
  addSignRule = 'addRule',
  deleteSignRule = 'deleteRule',
  updateSignRule = 'updateRule',
}

export interface ruleRequestData {
  WalletAddress: string;
  Nickname: string;
  Description: string;
  SignRule: {
    Name: string;
    Contract: string;
    RuleType: config.RuleType;
    Enabled: boolean;
    ID: string;
  };
}
