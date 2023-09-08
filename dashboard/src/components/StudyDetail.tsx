import React, { FC, useEffect, useState } from "react"
import { useRecoilValue } from "recoil"
import { Link, useParams } from "react-router-dom"
import { styled } from "@mui/system"
import {
  AppBar,
  Card,
  Typography,
  CardContent,
  Container,
  Grid,
  Toolbar,
  Paper,
  Box,
  IconButton,
  Select,
  MenuItem,
  useTheme,
} from "@mui/material"
import { Home, Cached } from "@mui/icons-material"

import { DataGridColumn, DataGrid } from "./DataGrid"
import { GraphParallelCoordinate } from "./GraphParallelCoordinate"
import { GraphIntermediateValues } from "./GraphIntermediateValues"
import { GraphSlice } from "./GraphSlice"
import { GraphHistory } from "./GraphHistory"
import { actionCreator } from "../action"
import { studyDetailsState } from "../state"
const theme = useTheme()

const useStyles = makeStyles((theme: Theme) =>
  createStyles({
    card: {
    },
    reload: {
      position: "relative",
      borderRadius: theme.shape.borderRadius,
      backgroundColor: fade(theme.palette.common.white, 0.15),
      "&:hover": {
        backgroundColor: fade(theme.palette.common.white, 0.25),
      },
      marginLeft: 0,
      width: "100%",
      [theme.breakpoints.up("sm")]: {
        marginLeft: theme.spacing(1),
        width: "auto",
      },
    },
    reloadIcon: {
      padding: theme.spacing(0, 2),
      height: "100%",
      position: "absolute",
      pointerEvents: "none",
      display: "flex",
      alignItems: "center",
      justifyContent: "center",
    },
    reloadSelect: {
      color: "inherit",
      padding: theme.spacing(1, 1, 1, 0),
      // vertical padding + font size from searchIcon
      paddingLeft: `calc(1em + ${theme.spacing(4)}px)`,
      transition: theme.transitions.create("width"),
      width: "100%",
      [theme.breakpoints.up("sm")]: {
        width: "14ch",
        "&:focus": {
          width: "20ch",
        },
      },
    },
  })
)

export const useStudyDetailValue = (studyId: number): StudyDetail | null => {
  const studyDetails = useRecoilValue<StudyDetails>(studyDetailsState)
  return studyDetails[studyId] || null
}

export const StudyDetail: FC<{}> = () => {
  const classes = useStyles()
  const action = actionCreator()
  const { studyId } = useParams<{studyId: string}>()
  const studyIdNumber = parseInt(studyId, 10)
  const studyDetail = useStudyDetailValue(studyIdNumber)
  const [openReloadIntervalSelect, setOpenReloadIntervalSelect] = useState<
    boolean
  >(false)
  const [reloadInterval, setReloadInterval] = useState<number>(10)

  useEffect(() => {
    action.updateStudyDetail(studyIdNumber)
  }, [])

  useEffect(() => {
    if (reloadInterval < 0) {
      return
    }
    const intervalId = setInterval(function () {
      action.updateStudyDetail(studyIdNumber)
    }, reloadInterval * 1000)
    return () => clearInterval(intervalId)
  }, [reloadInterval])

  const title = studyDetail !== null ? studyDetail.name : `Study #${studyId}`
  const trials: Trial[] = studyDetail !== null ? studyDetail.trials : []

  return (
    <div>
      <AppBar position="static">
        <Container>
          <Toolbar>
            <Typography variant="h6">{APP_BAR_TITLE}</Typography>
            <Box sx={{flexGrow: 1}} />
            <div
              className={classes.reload}
              onClick={(e) => {
                setOpenReloadIntervalSelect(!openReloadIntervalSelect)
              }}
            >
              <div className={classes.reloadIcon}>
                <Cached />
              </div>
              <Select
                value={reloadInterval}
                className={classes.reloadSelect}
                open={openReloadIntervalSelect}
                onOpen={() => {
                  setOpenReloadIntervalSelect(true)
                }}
                onClose={() => {
                  setOpenReloadIntervalSelect(false)
                }}
                onChange={(e) => {
                  setReloadInterval(e.target.value as number)
                }}
              >
                <MenuItem value={-1}>stop</MenuItem>
                <MenuItem value={5}>5s</MenuItem>
                <MenuItem value={10}>10s</MenuItem>
                <MenuItem value={30}>30s</MenuItem>
                <MenuItem value={60}>60s</MenuItem>
              </Select>
            </div>
            <IconButton
              aria-controls="menu-appbar"
              aria-haspopup="true"
              component={Link}
              to={URL_PREFIX + "/"}
              color="inherit"
            >
              <Home />
            </IconButton>
          </Toolbar>
        </Container>
      </AppBar>
      <Container>
        <>
          <Paper sx={{
            margin: theme.spacing(2),
            padding: theme.spacing(2),
          }}>
            <Typography variant="h6">{title}</Typography>
          </Paper>
          <Card sx={{margin: theme.spacing(2)}}>
            <CardContent>
              <GraphHistory study={studyDetail} />
            </CardContent>
          </Card>
          <Grid container direction="row">
            <Grid item xs={6}>
              <Card sx={{margin: theme.spacing(2)}}>
                <CardContent>
                  <GraphParallelCoordinate trials={trials} />
                </CardContent>
              </Card>
            </Grid>
            <Grid item xs={6}>
              <Card sx={{margin: theme.spacing(2)}}>
                <CardContent>
                  <GraphIntermediateValues trials={trials} />
                </CardContent>
              </Card>
            </Grid>
          </Grid>
          <Card sx={{margin: theme.spacing(2)}}>
            <CardContent>
              <GraphSlice trials={trials} />
            </CardContent>
          </Card>
          <Card sx={{margin: theme.spacing(2)}}>
            <TrialTable trials={trials} />
          </Card>
        </>
      </Container>
    </div>
  )
}

const TrialTable: FC<{ trials: Trial[] }> = ({ trials = [] }) => {
  const columns: DataGridColumn<Trial>[] = [
    { field: "number", label: "Number", sortable: true, padding: "none" },
    {
      field: "state",
      label: "State",
      sortable: true,
      filterable: true,
      padding: "none",
      toCellValue: (i) => trials[i].state.toString(),
    },
    { field: "value", label: "Value", sortable: true },
    {
      field: "datetime_start",
      label: "Duration(sec)",
      toCellValue: (i) => {
        const startMs = trials[i].datetime_start?.getTime()
        const completeMs = trials[i].datetime_complete?.getTime()
        if (startMs !== undefined && completeMs !== undefined) {
          return ((completeMs - startMs) / 1000).toString()
        }
        return null
      },
      sortable: true,
      less: (firstEl, secondEl): number => {
        const firstStartMs = firstEl.datetime_start?.getTime()
        const firstCompleteMs = firstEl.datetime_complete?.getTime()
        const firstDurationMs =
          firstStartMs !== undefined && firstCompleteMs !== undefined
            ? firstCompleteMs - firstStartMs
            : undefined
        const secondStartMs = secondEl.datetime_start?.getTime()
        const secondCompleteMs = secondEl.datetime_complete?.getTime()
        const secondDurationMs =
          secondStartMs !== undefined && secondCompleteMs !== undefined
            ? secondCompleteMs - secondStartMs
            : undefined

        if (firstDurationMs === secondDurationMs) {
          return 0
        } else if (
          firstDurationMs !== undefined &&
          secondDurationMs !== undefined
        ) {
          return firstDurationMs < secondDurationMs ? 1 : -1
        } else if (firstDurationMs !== undefined) {
          return -1
        } else {
          return 1
        }
      },
    },
    {
      field: "params",
      label: "Params",
      toCellValue: (i) =>
        trials[i].params.map((p) => p.name + ": " + p.value).join(", "),
    },
  ]
  const collapseParamColumns: DataGridColumn<TrialParam>[] = [
    { field: "name", label: "Name", sortable: true },
    { field: "value", label: "Value", sortable: true },
  ]
  const collapseIntermediateValueColumns: DataGridColumn<
    TrialIntermediateValue
  >[] = [
    { field: "step", label: "Step", sortable: true },
    { field: "value", label: "Value", sortable: true },
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
              Parameters
            </Typography>
            <DataGrid<TrialParam>
              columns={collapseParamColumns}
              rows={trials[index].params}
              keyField={"name"}
              dense={true}
              rowsPerPageOption={[5, 10, { label: "All", value: -1 }]}
            />
            <Typography variant="h6" gutterBottom component="div">
              Trial user attributes
            </Typography>
            <DataGrid<Attribute>
              columns={collapseAttrColumns}
              rows={trials[index].user_attrs}
              keyField={"key"}
              dense={true}
              rowsPerPageOption={[5, 10, { label: "All", value: -1 }]}
            />
          </Box>
        </Grid>
        <Grid item xs={6}>
          <Box margin={1}>
            <Typography variant="h6" gutterBottom component="div">
              Intermediate values
            </Typography>
            <DataGrid<TrialIntermediateValue>
              columns={collapseIntermediateValueColumns}
              rows={trials[index].intermediate_values}
              keyField={"step"}
              dense={true}
              rowsPerPageOption={[5, 10, { label: "All", value: -1 }]}
            />
            <Typography variant="h6" gutterBottom component="div">
              Trial system attributes
            </Typography>
            <DataGrid<Attribute>
              columns={collapseAttrColumns}
              rows={trials[index].system_attrs}
              keyField={"key"}
              dense={true}
              rowsPerPageOption={[5, 10, { label: "All", value: -1 }]}
            />
          </Box>
        </Grid>
      </Grid>
    )
  }

  return (
    <DataGrid<Trial>
      columns={columns}
      rows={trials}
      keyField={"trial_id"}
      dense={true}
      collapseBody={collapseBody}
    />
  )
}
