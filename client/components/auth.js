import { Token } from "../models/Token.js"

//component
var auth = {
	view: function() {
		return m("div", [
			m("input[autofocus=yes,id=token][type=text][placeholder=Token]", {
				oninput: m.withAttr("value", function(value) {
						Token.current = value
						Token.getInfo(value)
				}),
				value: Token.current
			}),
			m("p", "Token Information: ", JSON.stringify(Token.info))
		])
	}
}

export { auth }
