import React, { FC } from "react"
import { Link } from "react-router-dom"
import { createStyles, makeStyles, Theme } from "@material-ui/core/styles"
import {
  AppBar,
  Toolbar,
  Typography,
  Container,
  Card,
  Grid,
  Box,
} from "@material-ui/core"

import { actionCreator } from "../action"
import { formatDate } from "../dateUtil"
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
      field: "direction",
      label: "Direction",
      sortable: false,
      toCellValue: (i) => studies[i].direction.toString(),
    },
    {
      field: "best_trial",
      label: "Best value",
      sortable: false,
      toCellValue: (i) => studies[i].best_trial?.value || null,
    },
    {
      field: "datetime_start",
      label: "Datetime start",
      sortable: false,
      toCellValue: (i) => formatDate(studies[i].datetime_start),
    },
  ]

  const collapseAttrColumns: DataGridColumn<Attribute>[] = [
    { field: "key", label: "Key", sortable: true },
    { field: "value", label: "Value", sortable: true },
  ]

  const collapseBody = (index: number) => {
    return (
      <Grid container direction="row">
        <Grid item xs={6}>
          <Box margin={1}>
            <Typography variant="h6" gutterBottom component="div">
              Study user attributes
            </Typography>
            <DataGrid<Attribute>
              columns={collapseAttrColumns}
              rows={studies[index].user_attrs}
              keyField={"key"}
              dense={true}
            />
          </Box>
        </Grid>
        <Grid item xs={6}>
          <Box margin={1}>
            <Typography variant="h6" gutterBottom component="div">
              Study system attributes
            </Typography>
            <DataGrid<Attribute>
              columns={collapseAttrColumns}
              rows={studies[index].system_attrs}
              keyField={"key"}
              dense={true}
            />
          </Box>
        </Grid>
      </Grid>
    )
  }

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
            collapseBody={collapseBody}
          />
        </Card>
      </Container>
    </div>
  )
}
