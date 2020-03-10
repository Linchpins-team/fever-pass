async function deleteRow(button, type) {
    let response = await fetch(`/api/${type}/${button.value}`, {
        method: "DELETE", 
    })
    switch (response.status) {
        case 200:
            removeRow(button)
            break
        
        default:
            alert(await response.text())
            break
    }
}

function removeRow(button) {
    button.parentNode.parentNode.remove()
}

async function deleteAccount(button) {
    if (confirm("確定刪除帳號？")) {
        deleteRow(button, 'accounts')
    }
}