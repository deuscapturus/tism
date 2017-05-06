import { Encrypt } from "../models/Encrypt.js";
import { Keys } from "../models/Keys.js";

//component
var encrypt = {
  oninit: function(vnode) {
    Keys.getList();
  },
  view: function(vnode) {
    return m("div", [
      m(
        "div",
        Keys.available.map(function(key) {
          return key.Id != "ALL"
            ? m("div[class=radio]", [
                m("label", [
                  m("input[type=radio][name=key][class=radio]", {
                    value: key.Id,
                    id: key.Id,
                    oninput: m.withAttr("value", function(value) {
                      Encrypt.selectedKey = value;
                      if (Encrypt.input != "") {
                        Encrypt.encrypt();
                      }
                    })
                  }),
                  key.Name
                ])
              ])
            : null;
        })
      ),
      m(
        "div",
        m(
          "form",
          m("div[class=radio]", [
            m("label", [
              m("input[type=radio][name=key][class=radio]", {
                value: "armor",
                checked: Encrypt.encoding == "armor" ? true : false,
                oninput: m.withAttr("value", function(value) {
                  Encrypt.encoding = value;
                  if (Encrypt.input != "") {
                    Encrypt.encrypt();
                  }
                })
              }),
              "ASCII Armor Encoding"
            ]),
            m("label", [
              m("input[type=radio][name=key][class=radio]", {
                value: "base64",
                checked: Encrypt.encoding == "base64" ? true : false,
                oninput: m.withAttr("value", function(value) {
                  Encrypt.encoding = value;
                  if (Encrypt.input != "") {
                    Encrypt.encrypt();
                  }
                })
              }),
              "Base64 Encoding"
            ])
          ])
        )
      ),
      m("div[id=io]", [
        m("label[for=input]", "Input"),
        m("textarea[id=input][class=form-control][rows=2]", {
          oninput: m.withAttr("value", function(value) {
            Encrypt.input = value;
            Encrypt.encrypt();
          }),
          value: Encrypt.input
        }),
        m("label[for=output]", "Output"),
        m(
          "textarea[id=output][class=form-control][rows=4][readonly][style=word-break: break-all;]",
          {
            value: Encrypt.output,
            class: Encrypt.error ? "error" : null
          }
        )
      ])
    ]);
  }
};

export { encrypt };
