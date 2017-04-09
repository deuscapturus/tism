var extract = function(xhr) { 
  return xhr.status > 200 ? JSON.stringify(xhr.responseText) : xhr.responseText 
}

//model
var Token = {
	current: "",
	info: {},
	getInfo: function(token) {
		return m.request({
			method: "POST",
			url: "token/info",
			data: { "token": token },
			extract: extract
		})
		.then(function(result){
			Token.info = result
		})
	}
}

export { Token }
