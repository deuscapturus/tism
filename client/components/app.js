//component
var app = {
	view: function(vnode) {
		if (Object.keys(vnode.attrs.auth).length !== 0) {
			return m("h1", "App Loaded")
		} else {
			return m("h1", "App Not Loaded")
		}
	}
}

export { app }
