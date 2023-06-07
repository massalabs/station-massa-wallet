const errorCodes = new Map([
    ["Wallet-0001", "Wrong password. Try again"],
    ["Wallet-0002", "Error while retrieving that wallet. Try again"],
    ["Wallet-0003", "Action stopped"],
    ["Wallet-0004", "Enter a wallet nickname"],
    ["Wallet-0005", "Enter a wallet password"],
   
    ["Wallet-0006", "Error while creating your wallet. Try again"],

    ["Wallet-0007", "Select a wallet to delete"],
    [
        "Wallet-0008",
        "Error while deleting that wallet. Close all programs using this wallet and try again",
    ],
    [
        "Wallet-0009",
        "Error while connecting all your wallets. Reconnect all your wallets and try again",
    ],
    ["Wallet-0010", "Please enter a wallet name"],
    ["Wallet-0011", "Error while signing transaction"],
    ["Wallet-0012", "Action stopped. No wallet added to MassaStation"],
    ["Wallet-0013", "That nickname is taken. Try Another"],
    ["Wallet-0014", "Error while importing your wallet. Please try again"],
    ["Unknown-0001", "An unknown error occurred. Please try again"],
]);

function getErrorMessage(errorCode) {
    return errorCodes.get(errorCode) || errorCodes.get("Unknown-0001");
}

// If the error is from MassaStation, we display the error to the user and log the details in the console.
// Otherwise, we simply display the details in the console.
function handleAPIError(error) {
    if (error.response && error.response.data) {
        if (error.response.data.code) {
            errorAlert(getErrorMessage(error.response.data.code));
        }
        console.error("MassaStation error:", error.response.data);
    } else {
        errorAlert(getErrorMessage("Unknown-0001"));
        console.error(error);
    }
}
