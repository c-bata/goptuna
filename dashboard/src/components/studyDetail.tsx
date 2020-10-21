import React, { FC, useEffect, useState } from "react"
import { Link, useParams } from "react-router-dom"
import { useRecoilState } from "recoil"
import { createStyles, makeStyles, Theme } from "@material-ui/core/styles"
import Card from "@material-ui/core/Card"
import CardContent from "@material-ui/core/CardContent"

import { studyDetailsState } from "../state"
import { updateStudyDetail } from "../action"
import { TrialsTable } from "./trialsTable"
import { HistoryPlot } from "./historyPlot"
import { AppBar, Button, Container, Grid, Toolbar } from "@material-ui/core"
import Typography from "@material-ui/core/Typography"
import { ParallelCoordinatePlot } from "./parallelCoordinatePlot"
import { IntermediateValuesPlot } from "./intermediateValuesPlot"

const useStyles = makeStyles((theme: Theme) =>
  createStyles({
    card: {
      margin: theme.spacing(2),
    },
    grow: {
      flexGrow: 1,
    },
  })
)

interface ParamTypes {
  studyId: string
}

export const StudyDetail: FC<{}> = () => {
  const { studyId } = useParams<ParamTypes>()
  const studyIdNumber = parseInt(studyId, 10)
  const [ready, setReady] = useState(false)
  const [studyDetails, setStudyDetails] = useRecoilState<StudyDetails>(
    studyDetailsState
  )
  const classes = useStyles()

  useEffect(() => {
    // fetch immediately
    updateStudyDetail(studyIdNumber, studyDetails, setStudyDetails)
    const intervalId = setInterval(function () {
      updateStudyDetail(studyIdNumber, studyDetails, setStudyDetails)
    }, 1000)
    return () => clearInterval(intervalId)
  }, [])

  useEffect(() => {
    if (!ready && studyDetails[studyIdNumber]) {
      setReady(true)
    }
  }, [studyDetails])

  const studyDetail = studyDetails[studyIdNumber]
  const content = ready ? (
    <div>
      <Card className={classes.card}>
        <CardContent>
          <HistoryPlot trials={studyDetail.trials} />
        </CardContent>
      </Card>
      <Grid container direction="row">
        <Grid item xs={6}>
          <Card className={classes.card}>
            <CardContent>
              <ParallelCoordinatePlot trials={studyDetail.trials} />
            </CardContent>
          </Card>
        </Grid>
        <Grid item xs={6}>
          <Card className={classes.card}>
            <CardContent>
              <IntermediateValuesPlot trials={studyDetail.trials} />
            </CardContent>
          </Card>
        </Grid>
      </Grid>
      <Card className={classes.card}>
        <TrialsTable trials={studyDetail.trials} />
      </Card>
    </div>
  ) : (
    <p>Now loading...</p>
  )

  const title = studyDetail ? studyDetail.name : `Study #${studyId}`
  return (
    <div>
      <AppBar position="static">
        <Container>
          <Toolbar>
            <Typography variant="h6">{title}</Typography>
            <div className={classes.grow} />
            <Button color="inherit" component={Link} to="/">
              Return to Top
            </Button>
          </Toolbar>
        </Container>
      </AppBar>
      <Container>{content}</Container>
    </div>
  )
}
