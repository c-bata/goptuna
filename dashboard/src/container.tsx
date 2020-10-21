import { jsx, css } from "@emotion/core"
import { FC } from "react"
import { BrowserRouter as Router, Switch, Route } from "react-router-dom"
import { HistoryGraph } from "./components/historyGraph"
import { StudyList } from "./components/studyList"

const style = css``

const AppContainer: FC<{}> = () => {
  return (
    <Router css={style}>
      <Switch>
        <Route path="/studies/:studyId" children={<HistoryGraph />} />
        <Route path="/">
          <StudyList />
        </Route>
      </Switch>
    </Router>
  )
}

export { AppContainer }
