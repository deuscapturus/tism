//model
var Encrypt = {
	error: false,
	output: "",
	selectedKey: "",
	encrypt: function(token, decsecret, id) {
		return m.request({
			method: "POST",
			url: "encrypt",
			data: { "token": token, "decsecret": decsecret, "id": id },
			deserialize: function(value) {return value}
		})
		.then(function(result){
			Encrypt.output = result
			Encrypt.error = false
		})
		.catch(function(e) {
			Encrypt.output = e.message
			Encrypt.error = true
		})
	}
}

export { Encrypt }
