import * as plotly from "plotly.js-dist"
import React, { ChangeEvent, FC, useEffect, useState } from "react"
import {
  Grid,
  FormControl,
  FormLabel,
  FormControlLabel,
  Checkbox,
  Switch,
  Radio,
  RadioGroup,
} from "@material-ui/core"

export const HistoryPlot: FC<{
  study: StudyDetail
}> = ({
  study = {
    name: "",
    direction: "minimize",
    datetime_start: new Date(),
    trials: [],
  },
}) => {
  const [ready, setReady] = useState<boolean>(false)
  const [xAxis, setXAxis] = useState<string>("number")
  const [logScale, setLogScale] = useState<boolean>(false)
  const [filterCompleteTrial, setFilterCompleteTrial] = useState<boolean>(false)
  const [filterPrunedTrial, setFilterPrunedTrial] = useState<boolean>(false)

  useEffect(() => {
    setReady(true)
  }, [])

  const handleXAxisChange = (e: ChangeEvent<HTMLInputElement>) => {
    console.dir(e.target.value)
    setXAxis(e.target.value)
  }

  const handleLogScaleChange = (e: ChangeEvent<HTMLInputElement>) => {
    e.preventDefault()
    setLogScale(!logScale)
  }

  const handleFilterCompleteChange = (e: ChangeEvent<HTMLInputElement>) => {
    e.preventDefault()
    setFilterCompleteTrial(!filterCompleteTrial)
  }

  const handleFilterPrunedChange = (e: ChangeEvent<HTMLInputElement>) => {
    e.preventDefault()
    setFilterPrunedTrial(!filterPrunedTrial)
  }

  if (ready) {
    let filteredTrials = study.trials.filter(
      (t) => t.state === TrialState.Complete || t.state === TrialState.Pruned
    )
    if (filterCompleteTrial) {
      filteredTrials = filteredTrials.filter(
        (t) => t.state !== TrialState.Complete
      )
    }
    if (filterPrunedTrial) {
      filteredTrials = filteredTrials.filter(
        (t) => t.state !== TrialState.Pruned
      )
    }
    let trialsForLinePlot: Trial[] = []
    let currentBest: number | null = null
    filteredTrials.forEach((item) => {
      if (currentBest === null) {
        currentBest = item.value!
        trialsForLinePlot.push(item)
      } else if (
        study.direction === StudyDirection.Maximize &&
        item.value! > currentBest
      ) {
        currentBest = item.value!
        trialsForLinePlot.push(item)
      } else if (
        study.direction === StudyDirection.Minimize &&
        item.value! < currentBest
      ) {
        currentBest = item.value!
        trialsForLinePlot.push(item)
      }
    })

    const getAxisX = (trial: Trial): number | Date => {
      return xAxis === "number"
        ? trial.number
        : xAxis === "datetime_start"
        ? trial.datetime_start
        : trial.datetime_complete!
    }

    let xForLinePlot = trialsForLinePlot.map(getAxisX)
    xForLinePlot.push(getAxisX(filteredTrials[filteredTrials.length - 1]))
    let yForLinePlot = trialsForLinePlot.map((t: Trial): number => t.value!)
    yForLinePlot.push(yForLinePlot[yForLinePlot.length - 1])

    const plotData: Partial<plotly.PlotData>[] = [
      {
        x: filteredTrials.map(getAxisX),
        y: filteredTrials.map((t: Trial): number => t.value!),
        mode: "markers",
        type: "scatter",
      },
      {
        x: xForLinePlot,
        y: yForLinePlot,
        mode: "lines",
        type: "scatter",
      },
    ]
    const layout: Partial<plotly.Layout> = {
      margin: {
        l: 50,
        t: 0,
        r: 50,
        b: 0,
      },
      yaxis: {
        type: logScale ? "log" : "linear",
      },
      xaxis: {
        type: xAxis === "number" ? "linear" : "date",
      },
      showlegend: false,
    }
    plotly.react("history-plot", plotData, layout)
  }
  return (
    <Grid container direction="row">
      <Grid item xs={3}>
        <Grid container direction="column">
          <FormControl component="fieldset">
            <FormLabel component="legend">Log scale:</FormLabel>
            <Switch
              checked={logScale}
              onChange={handleLogScaleChange}
              value="enable"
            />
          </FormControl>
          <FormControl component="fieldset">
            <FormLabel component="legend">Filter state:</FormLabel>
            <FormControlLabel
              control={
                <Checkbox
                  checked={!filterCompleteTrial}
                  onChange={handleFilterCompleteChange}
                />
              }
              label="Complete"
            />
            <FormControlLabel
              control={
                <Checkbox
                  checked={!filterPrunedTrial}
                  onChange={handleFilterPrunedChange}
                />
              }
              label="Pruned"
            />
          </FormControl>
          <FormControl component="fieldset">
            <FormLabel component="legend">X-axis:</FormLabel>
            <RadioGroup
              aria-label="gender"
              name="gender1"
              value={xAxis}
              onChange={handleXAxisChange}
            >
              <FormControlLabel
                value="number"
                control={<Radio />}
                label="Number"
              />
              <FormControlLabel
                value="datetime_start"
                control={<Radio />}
                label="Datetime start"
              />
              <FormControlLabel
                value="datetime_complete"
                control={<Radio />}
                label="Datetime complete"
              />
            </RadioGroup>
          </FormControl>
        </Grid>
      </Grid>
      <Grid item xs={9}>
        <div id="history-plot" />
      </Grid>
    </Grid>
  )
}
