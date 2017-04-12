import { Token } from "./Token.js";

//model
var Keys = {
	available: [],
	error: false,
	getList: function() {
		return m.request({
			method: "POST",
			url: "key/list",
			data: { "token": Token.current },
		})
		.then(function(result){
			Keys.available = result
			Keys.error = false
		})
		.catch(function(e) {
			console.log(e.message)
			Keys.error = true
		})
	}
}

export { Keys }
