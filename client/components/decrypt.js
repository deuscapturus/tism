import { Decrypt } from "../models/Decrypt.js";

//component
var decrypt = {
  view: function(vnode) {
    return m("div", [
      m(
        "div",
        m(
          "form",
          m("div[class=radio]", [
            m("label", [
              m("input[type=radio][name=key][class=radio]", {
                value: "armor",
                checked: Decrypt.encoding == "armor" ? true : false,
                oninput: m.withAttr("value", function(value) {
                  Decrypt.encoding = value;
                  if (Decrypt.input != "") {
                    Decrypt.decrypt();
                  }
                })
              }),
              "ASCII Armor Encoding"
            ]),
            m("label", [
              m("input[type=radio][name=key][class=radio]", {
                value: "base64",
                checked: Decrypt.encoding == "base64" ? true : false,
                oninput: m.withAttr("value", function(value) {
                  Decrypt.encoding = value;
                  if (Decrypt.input != "") {
                    Decrypt.decrypt();
                  }
                })
              }),
              "Base64 Encoding"
            ])
          ])
        )
      ),
      m("div[id=io][class=form-group]", [
        m("label[for=input][class=control-label]", "Input"),
        m(
          "textarea[id=input][class=form-control][rows=4][style=word-break: break-all;]",
          {
            oninput: m.withAttr("value", function(value) {
              Decrypt.input = value;
              Decrypt.decrypt();
            }),
            value: Decrypt.input
          }
        ),
        m("label[for=output][class=control-label]", "Output"),
        m("textarea[id=output][class=form-control][rows=2][readonly]", {
          value: Decrypt.output,
          class: Decrypt.error ? "error" : null
        })
      ])
    ]);
  }
};

export { decrypt };
