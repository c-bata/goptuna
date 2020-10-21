import { jsx } from "@emotion/core"
import * as plotly from "plotly.js-dist"
import { FC, useEffect, useState } from "react"
import {Grid} from "@material-ui/core";
import Typography from "@material-ui/core/Typography";

export const HistoryPlot: FC<{
  trials: Trial[]
}> = ({ trials = [] }) => {
  const [ready, setReady] = useState(false)
  useEffect(() => {
    setReady(true)
  }, [])

  const completedTrials = trials.filter(
    (t: Trial) =>
      t.state === TrialState.Complete || t.state === TrialState.Pruned
  )
  const plotData: Partial<plotly.PlotData>[] = [
    {
      x: completedTrials.map((t: Trial): number => t.number),
      y: completedTrials.map((t: Trial): number => t.value || 0),
      mode: "lines+markers",
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
  }

  if (ready) {
    plotly.react("history-plot", plotData, layout)
  }
  return (
    <Grid
      container
      direction="row"
    >
      <Grid item xs={3}>
        <Typography variant="h5" component="h2">
          Plot settings
        </Typography>
      </Grid>
      <Grid item xs={9}>
        <div id="history-plot" />
      </Grid>
    </Grid>
  )
}
