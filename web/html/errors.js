const errorCodes = new Map([
    ["Wallet-01", "That nickname is taken. Try Another"],
    ["Wallet-02", "Wrong guiModal. Try again"],
    ["Wallet-03", "Error while retrieving that wallet. Try again"],
    ["Wallet-06", "Enter a wallet nickname"],
    ["Wallet-07", "Enter a wallet password"],
    ["Wallet-1003", "Error while creating your wallet. Try again"],
    ["Wallet-15", "Action stopped. No wallet added to Thyra"],
    ["Wallet-09", "Select a wallet to delete"],
    [
        "Wallet-10",
        "Error while deleting that wallet. Close all programs using this wallet and try again",
    ],
    ["Wallet-11", "Error while importing this wallet. Try again"],
    [
        "Wallet-4001",
        "Error while connecting all your wallets. Reconnect all your wallets and try again",
    ],
    ["Wallet-16" , "A wallet with the same nickname already exists."],
    ["Unknown-0001", "An unknown error occured. Please try again"],
]);

function getErrorMessage(errorCode) {
    return errorCodes.get(errorCode) || errorCodes.get("Unknown-0001");
}

// If the error is from Thyra, we display the error to the user and log the details in the console.
// Otherwise, we simply display the details in the console.
function handleAPIError(error) {
    if (error.response && error.response.data) {
        if (error.response.data.code) {
            errorAlert(getErrorMessage(error.response.data.code));
        }
        console.error("Thyra error:", error.response.data);
    } else {
        errorAlert(getErrorMessage("Unknown-0001"));
        console.error(error);
    }
}
