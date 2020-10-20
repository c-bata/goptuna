import * as ReactDOM from "react-dom"
import { jsx } from "@emotion/core"
import { fetchStudySummariesAction } from "./action"

const appDom = document.getElementById("dashboard")
if (appDom !== null) {
  fetchStudySummariesAction()
  ReactDOM.render(<h1>Hello World</h1>, appDom)
}
