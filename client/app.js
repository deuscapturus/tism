//component
var tISM = {};

//controller
tISM.controller = function() {

        this.tokeninfo = m.prop({})
	this.auth = m.prop(false)

	// updateToken is called from the child auth component with token info
	this.updateToken = function(tokeninfo) {
		this.tokeninfo(tokeninfo)
		if ( Object.keys(this.tokeninfo()).length != 0 && typeof this.tokeninfo() == "object" ) {
			this.auth(true)
		} else {
			this.auth(false)
		}
	}.bind(this)

};

//view
tISM.view = function(ctrl) {
	return [
		m.component(auth, {udpateToken: ctrl.updateToken}),
		(function() {
			if (ctrl.auth()) {
				return m("div", [
					m("h2", "Authenticated")
				])
			} else {
				return m("div", [
					m("h2", "Not Authenticated"),
					m("p", ctrl.tokeninfo())
				])
			}
		})(),
	]
};


//initialize
m.mount(document.getElementById("app"), tISM);
