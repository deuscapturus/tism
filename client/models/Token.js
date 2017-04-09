//model
var Token = {
	current: "",
	error: "",
	info: {},
	getInfo: function(token) {
		return m.request({
			method: "POST",
			url: "token/info",
			data: { "token": token },
		})
		.then(function(result){
			Token.info = result
		})
		.catch(function(e) {
			Token.info = {}
			Token.error = e.message
		})
	}
}

export { Token }
