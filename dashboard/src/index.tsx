import * as ReactDOM from "react-dom"
import { jsx } from "@emotion/core"
import { AppContainer } from "./container"
import { RecoilRoot } from "recoil"

const appDom = document.getElementById("dashboard")
if (appDom !== null) {
  ReactDOM.render(
    <RecoilRoot>
      <AppContainer />
    </RecoilRoot>,
    appDom
  )
}
