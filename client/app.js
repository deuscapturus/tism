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
	encrypt: function(token, decsecret, key) {
		return m.request({
			method: "POST",
			url: "encrypt",
			data: { "token": token, "decsecret": decsecret, "id": key[0] },
			deserialize: function(value) {return value;}
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
	// m.prop's are defined here with their defaults
	this.token = m.prop("")
	this.keys = m.prop([])
	this.selectedKeys = m.prop([])
	this.task = m.prop("Decrypt")
	this.admin = m.prop(0)
	this.input = m.prop("")
	this.output = m.prop("")


	// update functions are defined here to update m.prop's from the view and
	// run functions that these changes could affect.
	this.updateToken = function(token) {
		this.token(token)
		keys()
		newtoken()
		encrypt()
		decrypt()
	}.bind(this)

	this.updateAdmin = function(admin) {
		this.admin(((admin) ? 1 : 0))
		if (this.task()=="New Token") { newtoken() }
	}.bind(this)

	this.updatedSelectedKeys = function(selectedKeys) {
		// selectedOptions returns a HTMLCollection which
		// we have to cast into a normal array with values.
		var arr = [].slice.call(selectedKeys)
		arr.forEach(function(item, index, array) {
                        arr[index] = item.value; 
                })
		this.selectedKeys(arr)
		newtoken()
		
	}.bind(this)

	this.updateInput = function(input) {
		this.input(input)
		decrypt()
		encrypt()
	}.bind(this)

	this.updateTask = function(task) {
		this.task(task)
		decrypt()
		encrypt()
		newtoken()
	}.bind(this)


	// functions are defined here.  These are actions
	// that are triggered by the above update functions
	var keys = function() {
		tISM.keys(this.token()).then(this.keys)
	}.bind(this)

	var decrypt = function() {
		if (this.task() == "Decrypt") {
			tISM.decrypt(
				this.token(),
				this.input()
			).then(this.output, this.output)
		}
	}.bind(this)

	var encrypt = function() {
		if (this.task() == "Encrypt") {
			tISM.encrypt(
				this.token(),
				this.input(),
				this.selectedKeys()
			).then(this.output, this.output)
		}
	}.bind(this)

	var newtoken = function() {
		if (this.task() == "New Token") {
			tISM.newtoken(
				this.token(),
				this.selectedKeys(),
				this.admin()
			).then(this.output, this.output)
		}
	}.bind(this)

};

//view
tISM.view = function(ctrl) {
	return m("div"), [
		m("label[for=token]", "Token"),
		m("input[autofocus=yes,id=token]", {
			oninput: m.withAttr("value", ctrl.updateToken)
		}),
		m("button", {
				// TODO popup new modal for encryption key creation
			},
			"New Key"
		),
		m("label[for=keys]", "Encryption Keys"),
		m("select[name=keys]", {
				id: "keys",
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
			onchange: m.withAttr("value", ctrl.updateTask)
			}, [
				m("option", "Decrypt"),
				m("option", "Encrypt"),
				m("option", "New Token")
			]
		),
		m("div[id=io]", [
			m("label[for=admin]", "Make Admin"),
			m("input[id=admin]", {
				type: "checkbox",
				value: ctrl.admin(),
				onchange: m.withAttr("checked", ctrl.updateAdmin)
			}),
			m("label[for=input]", "Input"),
			m("textarea[id=input]", {
				value: ctrl.input(),
				oninput: m.withAttr("value", ctrl.updateInput)
			}),
			m("label[for=output]", "Output"),
			m("textarea[id=output]", {
				value: ctrl.output(),
			})
		]),

	]

};


//initialize
m.mount(document.getElementById("app"), tISM);
