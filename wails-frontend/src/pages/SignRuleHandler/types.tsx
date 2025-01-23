export enum RuleType {
  disablePasswordPrompt = 'disable_password_prompt',
  autoSign = 'auto_sign',
}

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
    RuleType: RuleType;
    Enabled: boolean;
    ID: string;
  };
}
