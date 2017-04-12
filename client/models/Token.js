//model
var Token = {
	current: sessionStorage.getItem("token"),
	error: sessionStorage.getItem("tokenError"),
	info: JSON.parse(sessionStorage.getItem("tokenInfo")) || {},
	getInfo: function(token) {
		return m.request({
			method: "POST",
			url: "token/info",
			data: { "token": token },
		})
		.then(function(result){
			Token.info = result
			sessionStorage.setItem("tokenInfo", JSON.stringify(result))
		})
		.catch(function(e) {
			Token.error = e.message
			sessionStorage.setItem("tokenError", e.message)
			Token.info = {}
			sessionStorage.setItem("tokenInfo", JSON.stringify({}))
		})
	}
}

export { Token }
