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
