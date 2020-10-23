import React, { FC } from "react"
import { Link } from "react-router-dom"
import { createStyles, makeStyles, Theme } from "@material-ui/core/styles"
import { AppBar, Toolbar, Typography, Container, Card } from "@material-ui/core"

import { actionCreator } from "../action"
import { formatDate } from "../utils/date"
import { useSnackbar } from "notistack"
import { useStudySummaries } from "../hook"
import { DataGrid, DataGridColumn } from "./dataGrid"

const useStyles = makeStyles((theme: Theme) =>
  createStyles({
    card: {
      margin: theme.spacing(2),
    },
  })
)

export const StudyList: FC<{}> = () => {
  const classes = useStyles()

  const { enqueueSnackbar } = useSnackbar()
  const action = actionCreator(enqueueSnackbar)
  const studies = useStudySummaries(action)

  const columns: DataGridColumn<StudySummary>[] = [
    {
      field: "study_id",
      label: "Study ID",
      sortable: true,
    },
    {
      field: "study_name",
      label: "Name",
      sortable: true,
      toCellValue: (i) => (
        <Link to={`/studies/${studies[i].study_id}`}>
          {studies[i].study_name}
        </Link>
      ),
    },
    {
      field: "datetime_start",
      label: "Datetime start",
      sortable: false,
      toCellValue: (i) => formatDate(studies[i].datetime_start),
    },
  ]

  return (
    <div>
      <AppBar position="static">
        <Container>
          <Toolbar>
            <Typography variant="h6">Goptuna dashboard</Typography>
          </Toolbar>
        </Container>
      </AppBar>
      <Container>
        <Card className={classes.card}>
          <DataGrid<StudySummary>
            columns={columns}
            rows={studies}
            keyField={"study_id"}
          />
        </Card>
      </Container>
    </div>
  )
}
