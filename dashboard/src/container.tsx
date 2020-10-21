import React, { FC } from "react"
import { BrowserRouter as Router, Switch, Route } from "react-router-dom"
import { StudyDetail } from "./components/studyDetail"
import { StudyList } from "./components/studyList"
import { SnackbarProvider } from "notistack"

const AppContainer: FC<{}> = () => {
  return (
    <SnackbarProvider maxSnack={3}>
      <Router>
        <Switch>
          <Route path="/studies/:studyId" children={<StudyDetail />} />
          <Route path="/">
            <StudyList />
          </Route>
        </Switch>
      </Router>
    </SnackbarProvider>
  )
}

export { AppContainer }
