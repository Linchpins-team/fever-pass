function delRecord(button) {
    fetch(`/api/records/${button.value}`, {
        method: "DELETE", 
    })
    .then(response => {
        switch (response.status) {
            case 200:
                removeRow(button)
                break
            
            default:
                throw response.text()
        }
    })
    .catch(err => {
        alert(err)
    })
}

function removeRow(button) {
    button.parentNode.parentNode.remove()
}