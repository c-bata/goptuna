import { jsx } from "@emotion/core"
import { FC, useEffect, useState } from "react"
import { useParams } from "react-router-dom"
import { useRecoilState } from "recoil"
import { createStyles, makeStyles, Theme } from "@material-ui/core/styles"
import Grid from "@material-ui/core/Grid"
import Card from "@material-ui/core/Card"
import CardContent from "@material-ui/core/CardContent"

import { studyDetailsState } from "../state"
import { updateStudyDetail } from "../action"
import { TrialsTable } from "./trialsTable"
import { HistoryPlot } from "./historyPlot"
import { Container } from "@material-ui/core"
import Typography from "@material-ui/core/Typography"

const useStyles = makeStyles((theme: Theme) =>
  createStyles({
    root: {
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
    <Grid container spacing={3}>
      <Grid item xs={12}>
        <Card>
          <CardContent>
            <Typography variant="h5" component="h2">
              Study {studyId}
            </Typography>
          </CardContent>
        </Card>
      </Grid>
      <Grid item xs={3}>
        <Card>
          <CardContent>
            <Typography variant="h5" component="h2">
              Plot settings
            </Typography>
          </CardContent>
        </Card>
      </Grid>
      <Grid item xs={9}>
        <Card>
          <CardContent>
            <HistoryPlot trials={studyDetail.trials} />
          </CardContent>
        </Card>
      </Grid>
      <Grid item xs={12}>
        <Card>
          <CardContent>
            <TrialsTable trials={studyDetail.trials} />
          </CardContent>
        </Card>
      </Grid>
    </Grid>
  ) : (
    <p>Now loading...</p>
  )
  return <Container className={classes.root}>{content}</Container>
}
