export enum RuleType {
  DisablePasswordPrompt = 'DISABLE_PASSWORD_PROMPT',
  AutoSign = 'AUTO_SIGN',
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
