closeModalOnClickOn("close-button");
closeModalOnClickOn("nicknameCancelBtn");
getWallets();

function addPrefixUrl(relativeURL) {
    return `/${relativeURL}`

//    return `/thyra/plugin/${pluginAuthor}/${pluginName}/${relativeURL}`
}

function openNickNameModal() {
    $("#nicknameModal").modal("show");
}

function closeModal() {
    $("#nicknameModal").modal("hide");
    document.getElementById("nicknameInput").value = "";
}

function closeModalOnClickOn(elementID) {
    document.getElementById(elementID).addEventListener("click", closeModal);
}

let wallets = [];

async function importWallet() {
    axios
        .put(addPrefixUrl(`api/accounts`))
        .then((resp) => {
            tableInsert(resp.data);
            wallets.push(resp.data);
        })
        .catch(handleAPIError);
}

// Create a wallet through POST query
async function getWallets() {
    axios
        .get(addPrefixUrl("api/accounts"))
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
function createAccount() {
    const nicknameCreate = document.getElementById("nicknameCreate").value;

    axios
        .post(addPrefixUrl("api/accounts"), {
            nickname: nicknameCreate,
        })
        .then((resp) => {
            tableInsert(resp.data);
            wallets.push(resp.data);
        })
        .catch(handleAPIError);
}

// Fetch a wallet's pending balance through GET query
async function fetchBalanceOf(address) {
    try {
        const getBalance = await axios.get(
            `/massa/addresses?attributes=balance&addresses=${address}`
        );
        return getBalance.data.addressesAttributes[address].balance.pending;
    } catch (error) {
        console.error(error)
        return '-'
    }
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

    cell2.innerHTML = await fetchBalanceOf(resp.address);
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
        .delete(addPrefixUrl(`api/accounts/${nickname}`))
        .then((_) => {
            wallets = wallets.filter((wallet) => wallet.nickname != nickname);
            document.getElementById("user-wallet-table").deleteRow(rowIndex);
        })
        .catch(handleAPIError);
}
