import React, { FC } from "react"
import { Link, useParams } from "react-router-dom"
import { createStyles, makeStyles, Theme } from "@material-ui/core/styles"
import {
  AppBar,
  Card,
  Typography,
  CardContent,
  Button,
  Container,
  Grid,
  Toolbar,
} from "@material-ui/core"

import { GraphParallelCoordinate } from "./graphParallelCoordinate"
import { GraphIntermediateValues } from "./graphIntermediateValues"
import { GraphHistory } from "./graphHistory"
import { actionCreator } from "../action"
import { useStudyDetail } from "../hook"
import { useSnackbar } from "notistack"
import { DataGridColumn, TrialsTable } from "./trialsTable"

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
  const classes = useStyles()
  const { enqueueSnackbar } = useSnackbar()
  const action = actionCreator(enqueueSnackbar)
  const { studyId } = useParams<ParamTypes>()
  const studyIdNumber = parseInt(studyId, 10)
  const studyDetail = useStudyDetail(action, studyIdNumber)

  const title = studyDetail !== null ? studyDetail.name : `Study #${studyId}`
  const trials: Trial[] = studyDetail !== null ? studyDetail.trials : []

  const columns: DataGridColumn<Trial>[] = [
    { field: "trial_id", label: "Trial ID", sortable: true },
    { field: "number", label: "Number", sortable: true },
    {
      field: "state",
      label: "State",
      sortable: false,
      toCellValue: (i) => trials[i].state.toString(),
    },
    { field: "value", label: "Value", sortable: true },
    {
      field: "params",
      label: "Params",
      sortable: false,
      toCellValue: (i) =>
        trials[i].params.map((p) => p.name + ": " + p.value).join(", "),
    },
  ]
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
      <Container>
        <div>
          <Card className={classes.card}>
            <CardContent>
              <GraphHistory study={studyDetail} />
            </CardContent>
          </Card>
          <Grid container direction="row">
            <Grid item xs={6}>
              <Card className={classes.card}>
                <CardContent>
                  <GraphParallelCoordinate trials={trials} />
                </CardContent>
              </Card>
            </Grid>
            <Grid item xs={6}>
              <Card className={classes.card}>
                <CardContent>
                  <GraphIntermediateValues trials={trials} />
                </CardContent>
              </Card>
            </Grid>
          </Grid>
          <Card className={classes.card}>
            <TrialsTable<Trial>
              columns={columns}
              rows={trials}
              keyField={"trial_id"}
            />
          </Card>
        </div>
      </Container>
    </div>
  )
}
