import { jsx } from "@emotion/core"
import { FC, useEffect, useState } from "react"
import { useParams } from "react-router-dom"
import { useRecoilState } from "recoil"
import { createStyles, makeStyles, Theme } from "@material-ui/core/styles"
import Paper from "@material-ui/core/Paper"
import Grid from "@material-ui/core/Grid"

import { studyDetailsState } from "../state"
import { updateStudyDetail } from "../action"
import { TrialsTable } from "./trialsTable"
import { HistoryPlot } from "./historyPlot"

const useStyles = makeStyles((theme: Theme) =>
  createStyles({
    root: {
      flexGrow: 1,
    },
    paper: {
      padding: theme.spacing(2),
      textAlign: "center",
      color: theme.palette.text.secondary,
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
    <Grid container spacing={3}>
      <Grid item xs={12}>
        <Paper className={classes.paper}>Study {studyId}</Paper>
      </Grid>
      <Grid item xs={12}>
        <HistoryPlot trials={studyDetail.trials} />
      </Grid>
      <Grid item xs={12}>
        <TrialsTable trials={studyDetail.trials} />
      </Grid>
    </Grid>
  ) : (
    <p>Now loading...</p>
  )
  return <div className={classes.root}>{content}</div>
}
