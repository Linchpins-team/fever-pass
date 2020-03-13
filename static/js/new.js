form = document.getElementById("form")
form.onsubmit = function (ev) {
    newRecord(form)
    ev.preventDefault()
}

async function newRecord(form) {
    let formData = new FormData(form)
    accountID = await getAccount(formData)
    if (accountID == null) {
        return
    }
    formData.append("account_id", accountID)

    let result = await sendNewRecord(formData)
    if (result != null) {
        let table = document.getElementById("table")
        table.insertAdjacentHTML("afterbegin", result)
        while (table.childElementCount > 20) {
            table.lastChild.remove()
        }
        form.reset()
        document.getElementsByName("class")[0].focus()
    }
}

async function getAccount(formData) {
    let response = await fetch(`/api/accounts?class=${formData.get("class")}&number=${formData.get("number")}`, {
        credentials: "include",
        method: "get",
    })
    if (response.status != 200) {
        alert(await response.text())
        return null
    }
    return await response.text()
}

async function sendNewRecord(formData) {
    let response = await fetch("/api/records", {
        credentials: "include",
        method: "post",
        body: formData,
    })
    if (response.status != 200) {
        alert(await response.text())
        return null
    }
    return await response.text()
}
