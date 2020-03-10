async function loadAccountList(select, class_name) {
    response = await fetch(`/api/stats?type=${select.value}&class=${class_name}`, {
        method: "get",
    })
    if (response.status != 200) {
        alert(await response.text())
        return
    }

    document.getElementById("list").innerHTML = await response.text()
}