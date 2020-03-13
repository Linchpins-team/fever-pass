let form = document.getElementById("form")
form.onsubmit = (event) => {
    setAccountID(event, form)
    event.preventDefault()
}

async function setAccountID(event, form) {
    let formData = new FormData(form)
    let id = await getAccount(formData)
    if (id == null) {
        form.reset()
        return
    }
    window.location = `/list?account_id=${id}&date=${formData.get("date")}`
}