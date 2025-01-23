export enum RuleType {
  disablePasswordPrompt = 'disable_password_prompt',
  autoSign = 'auto_sign',
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
