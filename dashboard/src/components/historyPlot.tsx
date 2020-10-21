import * as plotly from "plotly.js-dist"
import React, { ChangeEvent, FC, useEffect, useState } from "react"
import { Checkbox, Grid, Switch } from "@material-ui/core"
import FormControl from "@material-ui/core/FormControl"
import FormLabel from "@material-ui/core/FormLabel"
import { FormControlLabel, Radio, RadioGroup } from "@material-ui/core"

export const HistoryPlot: FC<{
  trials: Trial[]
}> = ({ trials = [] }) => {
  const [ready, setReady] = useState<boolean>(false)
  const [xAxis, setXAxis] = useState<string>("number")
  const [showPoints, setShowPoints] = useState<boolean>(true)
  const [logScale, setLogScale] = useState<boolean>(false)
  const [filterCompleteTrial, setFilterCompleteTrial] = useState<boolean>(false)
  const [filterPrunedTrial, setFilterPrunedTrial] = useState<boolean>(false)

  useEffect(() => {
    setReady(true)
  }, [])

  let filteredTrials = trials.filter(
    (t) => t.state === TrialState.Complete || t.state === TrialState.Pruned
  )
  if (filterCompleteTrial) {
    filteredTrials = filteredTrials.filter(
      (t) => t.state !== TrialState.Complete
    )
  }
  if (filterPrunedTrial) {
    filteredTrials = filteredTrials.filter((t) => t.state !== TrialState.Pruned)
  }
  const plotMode = showPoints ? "lines+markers" : "lines"
  const dataX =
    xAxis === "number"
      ? filteredTrials.map((t: Trial): number => t.number)
      : xAxis === "datetime_start"
      ? filteredTrials.map((t: Trial): Date => t.datetime_start)
      : filteredTrials.map((t: Trial): Date => t.datetime_complete!)
  const plotData: Partial<plotly.PlotData>[] = [
    {
      x: dataX,
      y: filteredTrials.map((t: Trial): number => t.value || 0),
      mode: plotMode,
      type: "scatter",
      name: "history",
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
      type: logScale ? 'log' : 'linear'
    },
    xaxis: {
      type: xAxis === "number" ? "linear" : "date"
    }
  }

  const handleXAxisChange = (e: ChangeEvent<HTMLInputElement>) => {
    console.dir(e.target.value)
    setXAxis(e.target.value)
  }

  const handleShowPointChange = (e: ChangeEvent<HTMLInputElement>) => {
    e.preventDefault()
    setShowPoints(!showPoints)
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
    plotly.react("history-plot", plotData, layout)
  }
  return (
    <Grid container direction="row">
      <Grid item xs={3}>
        <Grid container direction="column">
          <FormControl component="fieldset">
            <FormLabel component="legend">Show points:</FormLabel>
            <Switch checked={showPoints} onChange={handleShowPointChange} />
          </FormControl>
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
