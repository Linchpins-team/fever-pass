let form = document.getElementById("home-form")
form.onsubmit = ev => {
    login(form)
    ev.preventDefault()
}

async function login(form) {
    let response = await fetch("/api/login", {
        method: "POST",
        body: new FormData(form),
    })
    switch (response.status) {
        case 200:
            window.location = "/"
            break
        
        case 404:
            form.reset()
            errorMessage(await response.text())
            break

        case 403:
            errorMessage(await response.text())
            break
    }
}

function errorMessage(msg) {
    document.getElementById("msg").innerHTML = msg
} 