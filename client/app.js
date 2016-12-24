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
		m("h5", ctrl.tokeninfo()),
		m("h5", ctrl.auth()),
		(function() {
			if (ctrl.auth()) {
				return [
					m.component(menu),
					(function() {
						console.log(m.route.param("task"))
						switch(m.route.param("task")) {
							case "encrypt":
								console.log("in encrypt")
								return m.component(encrypt)
						}
					})()
				]
			} else {
				return [
					m("p[class=error]", ctrl.tokeninfo())
				]
			}
		})(),
	]
};


//initialize
//m.mount(document.getElementById("app"), tISM);
m.route.mode = "hash";
//routes
m.route(document.getElementById("app"), "/", {
    "/": tISM,
    "/:task": tISM
});
