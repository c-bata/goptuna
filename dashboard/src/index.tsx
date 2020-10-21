import * as ReactDOM from "react-dom"
import { jsx } from "@emotion/core"
import { AppContainer } from "./container"
import { RecoilRoot } from "recoil"

const dashboardDOM = document.getElementById("dashboard")
if (dashboardDOM !== null) {
  ReactDOM.render(
    <RecoilRoot>
      <AppContainer />
    </RecoilRoot>,
    dashboardDOM
  )
}
