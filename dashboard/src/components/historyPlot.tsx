import * as plotly from "plotly.js-dist"
import React, { ChangeEvent, FC, useEffect, useState } from "react"
import { Checkbox, Grid, Switch } from "@material-ui/core"
import FormControl from "@material-ui/core/FormControl"
import FormLabel from "@material-ui/core/FormLabel"
import { FormControlLabel, Radio, RadioGroup } from "@material-ui/core"

export const HistoryPlot: FC<{
  study: StudyDetail
}> = ({ study = {
  name: "",
  direction: "minimize",
  datetime_start: new Date(),
  trials: [],
} }) => {
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
      t => t.state === TrialState.Complete || t.state === TrialState.Pruned
    )
    if (filterCompleteTrial) {
      filteredTrials = filteredTrials.filter(
        t => t.state !== TrialState.Complete
      )
    }
    if (filterPrunedTrial) {
      filteredTrials = filteredTrials.filter((t) => t.state !== TrialState.Pruned)
    }
    let trialsForLinePlot: Trial[] = []
    let currentBest: number | null = null
    filteredTrials.forEach(item => {
      if (currentBest === null) {
        currentBest = item.value!
        trialsForLinePlot.push(item)
      } else if (study.direction === StudyDirection.Maximize && item.value! > currentBest) {
        currentBest = item.value!
        trialsForLinePlot.push(item)
      } else if (study.direction === StudyDirection.Minimize && item.value! < currentBest) {
        currentBest = item.value!
        trialsForLinePlot.push(item)
      }
    })

    const getAxisXList = (trials: Trial[]): number[] | Date[] => {
      return xAxis === "number"
        ? filteredTrials.map((t: Trial): number => t.number)
        : xAxis === "datetime_start"
          ? filteredTrials.map((t: Trial): Date => t.datetime_start)
          : filteredTrials.map((t: Trial): Date => t.datetime_complete!)
    }
    const plotData: Partial<plotly.PlotData>[] = [
      {
        x: getAxisXList(filteredTrials),
        y: filteredTrials.map((t: Trial): number => t.value!),
        mode: "markers",
        type: "scatter",
      },
      {
        x: getAxisXList(trialsForLinePlot),
        y: trialsForLinePlot.map((t: Trial): number => t.value || 0),
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
