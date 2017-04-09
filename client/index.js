import { index } from "./components/index.js";

m.route(document.body, "/decrypt", {
	"/": index,
    "/:task": index
})
