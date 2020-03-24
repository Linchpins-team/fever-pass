async function loadAccountList(value, class_name, date) {
    response = await fetch(`/api/stats?type=${value}&class=${class_name}&date=${date}`, {
        method: "get",
    })
    if (response.status != 200) {
        alert(await response.text())
        return
    }

    document.getElementById("list").innerHTML = await response.text()
}