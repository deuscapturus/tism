var encrypt = {};

//view
encrypt.view = function(vnode) {
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
				onchange: m.withAttr("selectedOptions", vnode.state.updatedSelectedKeys)
			},
			vnode.state.keys().map( function(key, index) {
				return m("option", {
						value: key.Id,
					},
					key.Id,
					key.Name,
					key.CreationTime
				)
		})),
		m("select[name=task]", {
			value: vnode.state.task(),
			onchange: m.withAttr("value", vnode.state.updateTask)
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
				value: vnode.state.admin(),
				onchange: m.withAttr("checked", vnode.state.updateAdmin)
			}),
			m("label[for=input]", "Input"),
			m("textarea[id=input]", {
				value: vnode.state.input(),
				oninput: m.withAttr("value", vnode.state.updateInput)
			}),
			m("label[for=output]", "Output"),
			m("textarea[id=output]", {
				value: vnode.state.output(),
			})
		]),

	])

};
