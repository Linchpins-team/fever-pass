let form = document.getElementById("form")
form.onsubmit = ev => {
    login(form)
    ev.preventDefault()
}

function login(form) {
    fetch("/api/login", {
        method: "POST",
        body: new FormData(form),
    })
    .then(response => {
        switch (response.status) {
            case 200:
                window.location = "/"
                break
            
            case 401:
                throw new Error("wrong password")
                break

            case 404:
                throw new Error("user not found")
                break
        }
    })
    .catch(err => {
        alert(err)
    })
}