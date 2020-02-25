function update(button, role, password) {
    let form = new FormData()
    if (role != 0) {
        form.append("role", role)
    }
    if (password != "") {
        form.append("password", password)
    }
    fetch(`/api/accounts/${button.value}`, {
        method: "PUT",
        body: form,
    })
    .then(response => {
        switch (response.status) {
        case 200:
            button.disabled = true
            break

        case 404:
            throw "user not found"

        case 415:
            throw "invalid data"

        case 401:
            throw "permission denied"
        }
    })
    .catch(err => {
        alert(err)
    })
}

function resetPassword(button) {
    let password = prompt("設定新密碼")
    update(button, 0, password)
    alert("密碼已更新")
}