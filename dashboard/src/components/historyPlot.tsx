import { jsx } from "@emotion/core"
import * as plotly from "plotly.js-dist"
import { FC, useEffect, useState } from "react"

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
  return <div id="history-plot" />
}
