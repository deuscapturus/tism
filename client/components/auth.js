//component
var auth = {};

var extract = function(xhr) {
  return xhr.status > 200 ? JSON.stringify(xhr.responseText) : xhr.responseText
}

//model
auth = { 
	info:  function(token) {
		return m.request({
			method: "POST",
		       	url: "token/info",
			data: { "token": token },
			extract: extract
		})
	}
};

//controller
auth.controller = function(parentctrl) {
	this.token = m.prop("")

	// UpdateToken update the token prop for this component and passes info to the parent compontent
	// that uses the same function name.
	this.updateToken = function(token) {
		this.token(token)
		auth.info(this.token()).then(parentctrl.udpateToken, parentctrl.udpateToken)
	}.bind(this)
};

//view
auth.view = function(ctrl) {
	return m("div", [
		m("label[for=token]", "Token"),
		m("input[autofocus=yes,id=token]", {
			oninput: m.withAttr("value", ctrl.updateToken)
		}),
	])
};
