let form = document.getElementById("form")
form.onsubmit = ev => {
    generateURL(form)
    ev.preventDefault()
}

function generateURL(form) {
    fetch("/api/url", {
        method: "POST",
        body: new FormData(form),
    })
    .then(response => {
        switch (response.status) {
            case 200:
                return response.json()
                break
            
            case 401:
                throw new Error("permission denied")
                break

            case 415:
                throw new Error("invalid data")
        }
    })
    .then(result => {
        document.getElementById("url").value = result.url
        document.getElementById("qrcode").src = result.qrcode
    })
    .catch(err => {
        alert(err)
    })
}

function copy() {
    var copyText = document.getElementById("url")

    copyText.select()
    copyText.setSelectionRange(0, 99999)

    document.execCommand("copy")

    alert("已複製到剪貼簿")
}