export interface SignRule {
  contract: string;
  passwordPrompt: boolean;
  autoSign: boolean;
}

export interface AccountConfig {
  signRules: SignRule[];
}

export interface Config {
  accounts: {
    [key: string]: AccountConfig;
  };
}
