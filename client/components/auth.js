import { Token } from "../models/Token.js"

//component
var auth = {
	view: function() {
		return m("div[class=form-group]", [
			m("input[autofocus=yes,id=token][type=text][placeholder=Token][class=form-control]", {
				oninput: m.withAttr("value", function(value) {
						Token.current = value
						sessionStorage.setItem("token", value);
						Token.getInfo(value)
				}),
				value: Token.current
			}),
			m("small[class=form-text text-muted]", Token.adminString, " user - ", Token.numKeys, " Keys Available")
		])
	}
}

export { auth }
