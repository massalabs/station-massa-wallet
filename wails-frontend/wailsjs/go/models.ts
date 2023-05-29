export namespace walletapp {
	
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

