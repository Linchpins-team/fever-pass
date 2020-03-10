async function loadAccountList(select) {
    response = await fetch(`/api/stats?type=${select.value}`, {
        method: "get",
    })
    if (response.status != 200) {
        alert(await response.text())
        return
    }

    document.getElementById("list").innerHTML = await response.text()
}