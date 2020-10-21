import React from 'react'
import {render} from "react-dom"
import { AppContainer } from "./container"
import { RecoilRoot } from "recoil"

render(
  <React.StrictMode>
    <RecoilRoot>
      <AppContainer />
    </RecoilRoot>
  </React.StrictMode>,
  document.getElementById("dashboard")
)
