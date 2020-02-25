function deleteRow(button, type) {
    fetch(`/api/${type}/${button.value}`, {
        method: "DELETE", 
    })
    .then(response => {
        switch (response.status) {
            case 200:
                removeRow(button)
                break
            
            case 404:
                throw new Error("user not found")

            case 401:
                throw new Error("permission denied")

            case 415:
                throw new Error("invalid id")
        }
    })
    .catch(err => {
        alert(err)
    })
}

function removeRow(button) {
    button.parentNode.parentNode.remove()
}