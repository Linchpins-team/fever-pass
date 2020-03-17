async function update(id, column, value) {
    let form = new FormData()
    form.append(column, value)
    let response = await fetch(`/api/accounts/${id}`, {
        method: "PUT",
        body: form,
    })
    if (response.status != 200) {
        alert(await response.text())
    }
}