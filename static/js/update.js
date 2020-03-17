async function update(select, id) {
    let form = new FormData()
    form.append("role", select.value)
    let response = await fetch(`/api/accounts/${id}`, {
        method: "PUT",
        body: form,
    })
    if (response.status != 200) {
        alert(await response.text())
    }
}