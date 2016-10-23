
var nonJsonErrors = function(xhr) {
	  return xhr.status > 200 ? JSON.stringify(xhr.responseText) : xhr.responseText
}

var tISM = {
	//model

	//controller
	controller: function() {
		return {
			token: "none",
			message: tISM.Key.list(this.token),
		}
	}
};

//view
tISM.view = function(ctrl) {
	return m("div"), [
		m("p", ctrl.token),
		m("h1", "DAMNIT!"),
		m("ol", ctrl.message().map( function(key, index) {
			return m("li", key.Id(), key.Name(), key.CreationTime())
		})),
		m("input", {
			onkeyup: (e) => {
				ctrl.token = e.target.value
			}
		})
	]
}

tISM.Key = function(data) {
	this.Id = m.prop(data.Id);
	this.Name = m.prop(data.Name);
	this.CreationTime = m.prop(data.CreationTime);
};

tISM.Key.list = function(token) {
	return m.request({
		method: "POST",
	       	url: "key/list",
		data: { "token": token },
		type: tISM.Key
	})
}

// ok. if this wasn't ui how would I do this.
//
// 1.  Setup an input that with any change makes an api call and changes the
// keys and scopes list.
// 2.  Create a component with that is a list of key ids from step 1
// 3.  Create a component for each api function.




//initialize
m.mount(document.getElementById("app"), tISM);

