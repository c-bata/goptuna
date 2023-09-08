import * as plotly from "plotly.js-dist-min"
import React, { FC, useEffect } from "react"
import { Typography, useTheme, Box } from "@mui/material"
import { plotlyDarkTemplate } from "./PlotlyDarkMode"
import { Target, useFilteredTrials } from "../trialFilter"

const domId = "graph-edf"

export const GraphEdf: FC<{ study: StudyDetail | null }> = ({ study }) => {
  const theme = useTheme()
  const target = new Target("objective", null)
  const trials = useFilteredTrials(study, [target], false)

  useEffect(() => {
    if (study !== null) {
      plotEdf(trials, target, domId, theme.palette.mode)
    }
  }, [trials, target, domId, theme.palette.mode])
  return (
    <Box>
      <Typography
        variant="h6"
        sx={{ margin: "1em 0", fontWeight: theme.typography.fontWeightBold }}
      >
        EDF
      </Typography>
      <Box id={domId} sx={{ height: "450px" }} />
    </Box>
  )
}

const plotEdf = (
  trials: Trial[],
  target: Target,
  domId: string,
  mode: string
) => {
  if (document.getElementById(domId) === null) {
    return
  }
  if (trials.length === 0) {
    plotly.react(domId, [], {
      template: mode === "dark" ? plotlyDarkTemplate : {},
    })
    return
  }

  const target_name = "Objective Value"
  const layout: Partial<plotly.Layout> = {
    xaxis: {
      title: target_name,
    },
    yaxis: {
      title: "Cumulative Probability",
    },
    margin: {
      l: 50,
      t: 0,
      r: 50,
      b: 50,
    },
    uirevision: "true",
    template: mode === "dark" ? plotlyDarkTemplate : {},
  }

  const values = trials.map((t) => target.getTargetValue(t) as number)
  const numValues = values.length
  const minX = Math.min(...values)
  const maxX = Math.max(...values)
  const numStep = 100
  const _step = (maxX - minX) / (numStep - 1)

  const xValues = []
  const yValues = []
  for (let i = 0; i < numStep; i++) {
    const boundary_right = minX + _step * i
    xValues.push(boundary_right)
    yValues.push(values.filter((v) => v <= boundary_right).length / numValues)
  }

  const plotData: Partial<plotly.PlotData>[] = [
    {
      type: "scatter",
      x: xValues,
      y: yValues,
    },
  ]
  plotly.react(domId, plotData, layout)
}
