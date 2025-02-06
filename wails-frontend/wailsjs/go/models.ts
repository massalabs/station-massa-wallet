export namespace config {
	
	export enum RuleType {
	    disable_password_prompt = "DISABLE_PASSWORD_PROMPT",
	    auto_sign = "AUTO_SIGN",
	}

}

export namespace walletapp {
	
	export enum PromptRequestAction {
	    delete = "DELETE_ACCOUNT",
	    newPassword = "CREATE_PASSWORD",
	    sign = "SIGN",
	    import = "IMPORT_ACCOUNT",
	    backup = "BACKUP_ACCOUNT",
	    tradeRolls = "TRADE_ROLLS",
	    unprotect = "UNPROTECT",
	    addSignRule = "ADD_SIGN_RULE",
	    deleteSignRule = "DELETE_SIGN_RULE",
	    updateSignRule = "UPDATE_SIGN_RULE",
	}
	export enum EventType {
	    promptResult = "PROMPT_RESULT",
	    promptData = "PROMPT_DATA",
	    promptRequest = "PROMPT_REQUEST",
	}
	export class selectFileResult {
	    err: string;
	    codeMessage: string;
	    filePath: string;
	    nickname: string;
	
	    static createFrom(source: any = {}) {
	        return new selectFileResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.err = source["err"];
	        this.codeMessage = source["codeMessage"];
	        this.filePath = source["filePath"];
	        this.nickname = source["nickname"];
	    }
	}

}

