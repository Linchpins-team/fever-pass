form = document.getElementById("form")
form.onsubmit = function (ev) {
    newRecord(form)
    ev.preventDefault()
}
function newRecord(form) {
    fetch("/api/records", {
        credentials: "include",
        method: "post",
        body: new FormData(form),
    })
    .then(response => {
        switch (response.status) {
            case 200:
                return response.json()
                break
            case 401:
                throw response.text()
                break
            case 415:
                throw response.text()
                break
        }
    })
    .then(record => {
        let table = document.getElementById("table")
        let tr = document.createElement("tr")
        tr.innerHTML = `
        <td>${record.ID}</td>
        <td>${record.Account}</td>
        <td>${record.Temperature}</td>
        <td>${record.Type}</td>
        <td>${record.Fever?'發燒':'正常'}</td>
        <td>${formatTime(record.CreatedAt)}</td>
        <td>${record.RecordedBy}</td>
        <td><button value=${record.ID} onclick="delRecord(this)">刪除</button></td>`
        table.insertBefore(tr, table.firstChild)
        while (table.childElementCount > 20) {
            table.lastChild.remove()
        }
        form.reset()
    })
    .catch(error => {
        alert(error)
    })
}

function formatTime(t) {
    let date = new Date(t)
    return `${two(date.getMonth()+1)}-${two(date.getDate())} ${two(date.getHours())}:${two(date.getMinutes())}`
}

function two(num) {
    return `${num<10?"0":""}${num}`
}