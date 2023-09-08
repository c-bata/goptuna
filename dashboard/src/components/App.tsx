import React, { FC } from "react"
import { RecoilRoot } from "recoil"
import { BrowserRouter as Router, Routes, Route } from "react-router-dom"
import { SnackbarProvider } from "notistack"

import { StudyDetail } from "./StudyDetail"
import { StudyList } from "./StudyList"

export const App: FC<{}> = () => {
  return (
    <RecoilRoot>
      <SnackbarProvider maxSnack={3}>
        <Router>
          <Routes>
            <Route
              path={URL_PREFIX + "/studies/:studyId"}
              element={<StudyDetail />}
            />
            <Route path={URL_PREFIX + "/"} element={<StudyList />} />
          </Routes>
        </Router>
      </SnackbarProvider>
    </RecoilRoot>
  )
}
