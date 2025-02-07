export enum RuleType {
  DisablePasswordPrompt = 'disable_password_prompt',
  AutoSign = 'auto_sign',
}

export interface SignRule {
  id: string;
  ruleType: RuleType;
  name?: string;
  contract: string;
  enabled: boolean;
}

export interface AccountConfig {
  signRules: SignRule[];
}

export interface Config {
  accounts: {
    [key: string]: AccountConfig;
  };
}
