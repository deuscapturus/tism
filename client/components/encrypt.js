var encrypt = {};

//view
encrypt.view = function(ctrl) {
	return m("div", [
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

	])

};
