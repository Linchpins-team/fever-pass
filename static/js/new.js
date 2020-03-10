form = document.getElementById("form")
form.onsubmit = function (ev) {
    newRecord(form)
    ev.preventDefault()
}
async function newRecord(form) {
    let response = await fetch("/api/records", {
        credentials: "include",
        method: "post",
        body: new FormData(form),
    })
    switch (response.status) {
        case 200:
            let text = await response.text()
            let table = document.getElementById("table")
            table.insertAdjacentHTML("afterbegin", text)
            while (table.childElementCount > 20) {
                table.lastChild.remove()
            }
            form.reset()
            document.getElementsByName("class")[0].focus()
            break

        case 401:
            alert("需要登入")
            break

        case 415:
            alert("資料無效")
            break

        case 403:
            alert("沒有權限進行此操作")
            break

        case 404:
            alert("找不到此帳號")
            break
    }
}

function formatTime(t) {
    let date = new Date(t)
    return `${two(date.getMonth()+1)}-${two(date.getDate())} ${two(date.getHours())}:${two(date.getMinutes())}`
}

function two(num) {
    return `${num<10?"0":""}${num}`
}