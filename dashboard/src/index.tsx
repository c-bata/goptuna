import * as ReactDOM from "react-dom"
import { jsx } from "@emotion/core"

const appDom = document.getElementById("dashboard")
if (appDom !== null) {
  ReactDOM.render(<h1>Hello World</h1>, appDom)
}
