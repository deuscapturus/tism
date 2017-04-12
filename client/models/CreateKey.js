import { Keys } from "./Keys.js";
import { Token } from "./Token.js";

//model
var CreateKey = {
	name: "",
	comment: "",
	email: "",
	error: false,
	create: function() {
		return m.request({
			method: "POST",
			url: "key/new",
			data: { "token": Token.current, "name": CreateKey.name, "comment": CreateKey.comment, "email": CreateKey.email },
		})
		.then(function(result){
			CreateKey.error = false
			Keys.getList()
		})
		.catch(function(e) {
			console.log(e.message)
			CreateKey.error = true
		})
	}
}

export { CreateKey }
