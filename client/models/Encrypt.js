import { Token } from "./Token.js";

//model
var Encrypt = {
	error: false,
	output: "",
	selectedKey: "",
	encrypt: function(decsecret, id) {
		return m.request({
			method: "POST",
			url: "encrypt",
			data: { "token": Token.current, "decsecret": decsecret, "id": id },
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
