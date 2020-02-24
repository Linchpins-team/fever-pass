let form = document.getElementById("form")
form.onsubmit = ev => {
    register(form)
    ev.preventDefault()
}

function register(form) {
    fetch("/api/register", {
        method: "POST",
        body: new FormData(form),
    })
    .then(response => {
        switch (response.status) {
            case 200:
                window.location = "/admin/new"
                break
            
            case 401:
                throw new Error("invalid key")
                break
        }
    })
    .catch(err => {
        alert(err)
    })
}