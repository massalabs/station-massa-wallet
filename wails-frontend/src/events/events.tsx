
export type promptResult = {
    Success: boolean,
    Data: any,
}

export enum promptAction {
    deleteReq = 0,
    newPasswordReq = 1,
    signReq = 2,
    importReq = 3,
    exportReq = 4,
}

export type promptRequest = {
    Action: promptAction,
    Msg: string,
    Data: any
}

export const events = {
    passwordResult: "passwordResult",
    promptRequest :"promptRequest"
}

