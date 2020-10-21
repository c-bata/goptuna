import { jsx } from "@emotion/core"
import * as plotly from "plotly.js-dist"
import { FC, useEffect, useState } from "react"

export const IntermediateValuesPlot: FC<{
  trials: Trial[]
}> = ({ trials = [] }) => {
  const [ready, setReady] = useState<boolean>(false)

  useEffect(() => {
    setReady(true)
  }, [])

  let filteredTrials = trials.filter(
    (t) => t.state === TrialState.Complete || t.state === TrialState.Pruned
  )
  const plotData: Partial<plotly.PlotData>[] = [
    {
      x: filteredTrials.map((t: Trial): number => t.number),
      y: filteredTrials.map((t: Trial): number => t.value || 0),
      mode: "lines",
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
    plotly.react("intermediate-values-plot", plotData, layout)
  }
  return (
    <div id="intermediate-values-plot" />
  )
}
