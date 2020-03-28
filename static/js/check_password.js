function checkPassword() { 
    let form = document.getElementById("form")
    let password1 = form.new_password.value
    let password2 = form.match_password.value

    if (password1 != password2) 
        document.getElementById("msg").innerHTML = "密碼不相符"

    else {
        document.getElementById("msg").innerHTML = ""
        document.getElementById("reset").disabled = false
    }
}