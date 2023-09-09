import * as plotly from "plotly.js-dist-min"
import React, { ChangeEvent, FC, useEffect, useState } from "react"
import {
  Grid,
  FormControl,
  FormLabel,
  FormControlLabel,
  MenuItem,
  Select,
  Radio,
  RadioGroup,
  Typography,
  SelectChangeEvent,
  useTheme,
} from "@mui/material"
import { plotlyDarkTemplate } from "./PlotlyDarkMode"
import {
  useFilteredTrials,
  Target,
  useObjectiveAndUserAttrTargets,
} from "../trialFilter"

const plotDomId = "graph-history"

export const GraphHistory: FC<{
  study: StudyDetail | null
  logScale: boolean
  includePruned: boolean
}> = ({ study, logScale, includePruned }) => {
  const theme = useTheme()
  const [xAxis, setXAxis] = useState<
    "number" | "datetime_start" | "datetime_complete"
  >("number")
  const [targets, selected, setTarget] = useObjectiveAndUserAttrTargets(study)
  const trials = useFilteredTrials(study, [selected], !includePruned)

  useEffect(() => {
    if (study !== null) {
      plotHistory(
        trials,
        study.direction,
        selected,
        xAxis,
        logScale,
        theme.palette.mode
      )
    }
  }, [trials, study?.direction, selected, logScale, xAxis, theme.palette.mode])

  const handleObjectiveChange = (event: SelectChangeEvent<string>) => {
    setTarget(event.target.value)
  }

  const handleXAxisChange = (e: ChangeEvent<HTMLInputElement>) => {
    if (e.target.value === "number") {
      setXAxis("number")
    } else if (e.target.value === "datetime_start") {
      setXAxis("datetime_start")
    } else if (e.target.value === "datetime_complete") {
      setXAxis("datetime_complete")
    }
  }

  return (
    <Grid container direction="row">
      <Grid
        item
        xs={3}
        container
        direction="column"
        sx={{ paddingRight: theme.spacing(2) }}
      >
        <Typography
          variant="h6"
          sx={{ margin: "1em 0", fontWeight: theme.typography.fontWeightBold }}
        >
          History
        </Typography>
        {targets.length >= 2 ? (
          <FormControl
            component="fieldset"
            sx={{ marginBottom: theme.spacing(2) }}
          >
            <FormLabel component="legend">y Axis</FormLabel>
            <Select
              value={selected.identifier()}
              onChange={handleObjectiveChange}
            >
              {targets.map((t, i) => (
                <MenuItem value={t.identifier()} key={i}>
                  {t.toLabel()}
                </MenuItem>
              ))}
            </Select>
          </FormControl>
        ) : null}
        <FormControl
          component="fieldset"
          sx={{ marginBottom: theme.spacing(2) }}
        >
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
      <Grid item xs={9}>
        <div id={plotDomId} />
      </Grid>
    </Grid>
  )
}

const plotHistory = (
  trials: Trial[],
  direction: StudyDirection,
  target: Target,
  xAxis: "number" | "datetime_start" | "datetime_complete",
  logScale: boolean,
  mode: string
) => {
  if (document.getElementById(plotDomId) === null) {
    return
  }

  const layout: Partial<plotly.Layout> = {
    margin: {
      l: 50,
      t: 0,
      r: 50,
      b: 0,
    },
    yaxis: {
      title: target.toLabel(),
      type: logScale ? "log" : "linear",
    },
    xaxis: {
      title: xAxis === "number" ? "Trial" : "Time",
      type: xAxis === "number" ? "linear" : "date",
    },
    showlegend: true,
    uirevision: "true",
    template: mode === "dark" ? plotlyDarkTemplate : {},
  }
  if (trials.length === 0) {
    plotly.react(plotDomId, [], layout)
    return
  }

  const getAxisX = (trial: Trial): number | Date => {
    return xAxis === "number"
      ? trial.number
      : xAxis === "datetime_start"
      ? trial.datetime_start!
      : trial.datetime_complete!
  }

  const plotData: Partial<plotly.PlotData>[] = [
    {
      x: trials.map(getAxisX),
      y: trials.map((t: Trial): number => target.getTargetValue(t) as number),
      name: target.toLabel(),
      mode: "markers",
      type: "scatter",
    },
  ]

  if (target.kind === "objective") {
    const xForLinePlot: (number | Date)[] = []
    const yForLinePlot: number[] = []
    let currentBest: number | null = null
    for (let i = 0; i < trials.length; i++) {
      const t = trials[i]
      if (currentBest === null) {
        currentBest = t.value as number
        xForLinePlot.push(getAxisX(t))
        yForLinePlot.push(t.value as number)
      } else if (
        direction === "maximize" &&
        (t.value as number) > currentBest
      ) {
        const p = trials[i - 1]
        if (!xForLinePlot.includes(getAxisX(p))) {
          xForLinePlot.push(getAxisX(p))
          yForLinePlot.push(currentBest)
        }
        currentBest = t.value as number
        xForLinePlot.push(getAxisX(t))
        yForLinePlot.push(t.value as number)
      } else if (
        direction === "minimize" &&
        (t.value as number) < currentBest
      ) {
        const p = trials[i - 1]
        if (!xForLinePlot.includes(getAxisX(p))) {
          xForLinePlot.push(getAxisX(p))
          yForLinePlot.push(currentBest)
        }
        currentBest = t.value as number
        xForLinePlot.push(getAxisX(t))
        yForLinePlot.push(t.value as number)
      }
    }
    xForLinePlot.push(getAxisX(trials[trials.length - 1]))
    yForLinePlot.push(yForLinePlot[yForLinePlot.length - 1])
    plotData.push({
      x: xForLinePlot,
      y: yForLinePlot,
      name: "Best Value",
      mode: "lines",
      type: "scatter",
    })
  }
  plotly.react(plotDomId, plotData, layout)
}
