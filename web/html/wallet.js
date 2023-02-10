document
    .getElementById("import-wallet")
    .addEventListener("click", openNickNameModal);

closeModalOnClickOn("close-button");
closeModalOnClickOn("nicknameCancelBtn");
getWallets();

function openNickNameModal() {
    $("#nicknameModal").modal("show");
}

function closeModal() {
    $("#nicknameModal").modal("hide");
    document.getElementById("nicknameInput").value = ""
}

function closeModalOnClickOn(elementID) {
    document.getElementById(elementID).addEventListener("click", closeModal);
}

let wallets = [];

// Import a wallet through PUT query
async function importWallet() {
    let nickname = document.getElementById("nicknameInput").value;
    axios
        .post(`/rest/wallet/import/${nickname}`)
        .then((resp) => {
            tableInsert(resp.data);
            wallets.push(resp.data);
        })
        .catch(handleAPIError);
    closeModal()
}

// Create a wallet through POST query
async function getWallets() {
    axios
        .get("/rest/wallet")
        .then((resp) => {
            if (resp) {
                const data = resp.data;
                for (const wallet of data) {
                    tableInsert(wallet);
                }
                wallets = data;
            }
        })
        .catch(handleAPIError);
}

// Create a wallet through POST query
function createWallet() {
    const nicknameCreate = document.getElementById("nicknameCreate").value;
    const password = document.getElementById("password").value;

    axios
        .post("/rest/wallet", {
            nickname: nicknameCreate,
            password: password,
        })
        .then((resp) => {
            tableInsert(resp.data);
            wallets.push(resp.data);
        })
        .catch(handleAPIError);
}

// Fetch a wallet's balance through POST query
async function fetchBalanceOf(addresses) {
    const getBalance = await axios.post("http://my.massa/massa/addresses", {
        addresses: addresses,
        options: ["balances"],
    });

    return getBalance.data.pendingBalances;
}

async function tableInsert(resp) {
    const tBody = document
        .getElementById("user-wallet-table")
        .getElementsByTagName("tbody")[0];
    const row = tBody.insertRow(-1);

    const cell0 = row.insertCell();
    const cell1 = row.insertCell();
    const cell2 = row.insertCell();
    const cell3 = row.insertCell();

    cell0.innerHTML = addressInnerHTML(resp.address);

    cell1.innerHTML = resp.nickname;

    const balance = await fetchBalanceOf([resp.address]);
    cell2.innerHTML = parseFloat(balance[0]);
    cell3.innerHTML =
        '<svg class="quit-button" onclick="deleteRow(this)" xmlns="http://www.w3.org/2000/svg" width="24" height="22" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="feather feather-x"><line x1="18" y1="6" x2="6" y2="18"></line> <line x1="6" y1="6" x2="18" y2="18"></line></svg>';
}

function deleteRow(element) {
    const rowIndex = element.parentNode.parentNode.rowIndex;

    const tBody = document
        .getElementById("user-wallet-table")
        .getElementsByTagName("tbody")[0];
    const nickname = tBody.rows[rowIndex - 1].cells[1].innerHTML;

    axios
        .delete("/rest/wallet/" + nickname)
        .then((_) => {
            wallets = wallets.filter((wallet) => wallet.nickname != nickname);
        })
        .catch(handleAPIError);

    document.getElementById("user-wallet-table").deleteRow(rowIndex);
}
