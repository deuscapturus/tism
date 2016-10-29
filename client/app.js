//component
var tISM = {};

//model
tISM.Key = function(data) {
	this.Id = m.prop(data.Id);
	this.Name = m.prop(data.Name);
	this.CreationTime = m.prop(data.CreationTime);
}
tISM = { 

	keys:  function(token) {
		return m.request({
			method: "POST",
		       	url: "key/list",
			data: { "token": token },
			type: tISM.Key
		})
	},
	decrypt: function(token, encsecret) {
		return m.request({
			method: "POST",
			url: "decrypt",
			data: { "token": token, "encsecret": encsecret },
			deserialize: function(value) {return value;}
		})
	},
	newtoken: function(token, keys, admin) {
		alert(admin)
		return m.request({
			method: "POST",
			url: "token/new",
			data: { "token": token, "keys": keys, "admin": admin },
			deserialize: function(value) {return value;}
		})
	},
};

//controller
tISM.controller = function() {
	this.token = m.prop("")
	this.keys = m.prop([])
	this.selectedKeys = m.prop([])
	this.task = m.prop("Decrypt")
	this.admin = m.prop(0)
	this.input = m.prop("")
	this.output = m.prop("")

	this.updateToken = function(token) {
		this.token(token)
		tISM.keys(this.token()).then(this.keys)
	}.bind(this)

	this.updateAdmin = function(admin) {
		this.admin(((admin) ? 1 : 0))
		this.updatedSelectedKeys()
	}.bind(this)

	this.updatedSelectedKeys = function(selectedKeys) {
		// selectedOptions returns a HTMLCollection which
		// we have to cast into a normal array with values.
		var arr = [].slice.call(selectedKeys)
		arr.forEach(function(item, index, array) {
                        arr[index] = item.value; 
                })
		this.selectedKeys(arr)
		
		switch(this.task()) {
			//TODO add case for encrypt secret
			case "New Token":
				tISM.newtoken(this.token(), this.selectedKeys(), this.admin()).then(this.output, this.output)
				break;
		}
	}.bind(this)

	this.updateInput = function(input) {
		this.input(input)
		switch(this.task()) {
			//TODO add case for encrypt
			case "Decrypt":
				tISM.decrypt(this.token(), this.input()).then(this.output, this.output)
				break;
		}
	}.bind(this)

};

//view
tISM.view = function(ctrl) {
	return m("div"), [
		m("input[autofocus=yes]", {
			oninput: m.withAttr("value", ctrl.updateToken)
		}),
		m("button", {
				// TODO popup new modal for encryption key creation
			},
			"New Key"
		),
		m("select[name=key]", {
				id: "test",
				multiple: "multiple",
				onchange: m.withAttr("selectedOptions", ctrl.updatedSelectedKeys)
			},
			ctrl.keys().map( function(key, index) {
				return m("option", {
						value: key.Id,
					},
					key.Id,
					key.Name,
					key.CreationTime
				)
		})),
		m("select[name=task]", {
			value: ctrl.task(),
			onchange: m.withAttr("value", ctrl.task)
			}, [
				m("option", "Decrypt"),
				m("option", "Encrypt"),
				m("option", "New Token")
			]
		),
		m("div[id=io]", [
			m("input[name=isAdmin]", {
				type: "checkbox",
				value: ctrl.admin(),
				onchange: m.withAttr("checked", ctrl.updateAdmin)
			}),
			m("textarea[name=input]", {
				value: ctrl.input(),
				oninput: m.withAttr("value", ctrl.updateInput)
			}),
			m("textarea[name=output]", {
				value: ctrl.output(),
			})
		]),

	]

};


//initialize
m.mount(document.getElementById("app"), tISM);
