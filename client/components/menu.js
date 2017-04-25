import { Token } from "../models/Token.js";

//component
var menu = {
  view: function() {
    return m(
      "div",
      m("ul[class=nav nav-tabs]", [
        m(
          "li[class=nav-item]",
          m(
            "a[href=/decrypt][class=nav-link]",
            {
              oncreate: m.route.link,
              class: m.route.param("task") == "decrypt" ? "active" : null
            },
            "Decrypt"
          )
        ),
        m(
          "li[class=nav-item]",
          m(
            "a[href=/encrypt][class=nav-link]",
            {
              oncreate: m.route.link,
              class: m.route.param("task") == "encrypt" ? "active" : null
            },
            "Encrypt"
          )
        ),
        m(
          "li[class=nav-item]",
          m(
            "a[href=/keys][class=nav-link]",
            {
              oncreate: m.route.link,
              class: m.route.param("task") == "keys" ? "active" : null
            },
            "Keys"
          )
        ),
        m(
          "li[class=nav-item]",
          m(
            "a[href=/tokens][class=nav-link]",
            {
              oncreate: m.route.link,
              class: m.route.param("task") == "tokens" ? "active" : null
            },
            "Tokens"
          )
        )
        //		(Token.info.admin == 1) ? m("li[class=dropdown]",
        //			{class: (m.route.param("task") == "tokens"|| m.route.param("task") == "keys") ? "active" : null}, [
        //			m("a[href=#][class=dropdown-toggle][data-toggle=dropdown][role=button][aria-haspopup=true][aria-expanded=false]",
        //				"Admin ",
        //				m("span[class=caret]"),
        //			),
        //			m("ul[class=navbar-nav navbar-right]", m("ul[class=dropdown-menu]", [
        //				m("li", m("a[href=/tokens]", {oncreate: m.route.link}, "Tokens")),
        //				m("li", m("a[href=/keys]", {oncreate: m.route.link}, "Keys"))
        //			]))
        //		]) : null
      ])
    );
  }
};

export { menu };
