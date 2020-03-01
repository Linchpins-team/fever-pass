let form = document.getElementById("home-form")
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
                throw new Error("密碼錯誤")

            case 404:
                throw new Error("找不到此帳號")
        }
    })
    .catch(err => {
        document.getElementById("msg").innerHTML = err
    })
}