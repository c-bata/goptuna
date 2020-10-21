import * as plotly from "plotly.js-dist"
import React, { FC, useEffect, useState } from "react"

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
  const plotData: Partial<plotly.PlotData>[] = filteredTrials.map(trial => {
    return {
      x: trial.intermediate_values.map(iv => iv.step),
      y: trial.intermediate_values.map(iv => iv.value),
      mode: "lines+markers",
      type: "scatter",
      name: `trial #${trial.number}`,
    }
  })
  const layout: Partial<plotly.Layout> = {
    title: "Intermediate values",
    margin: {
      l: 50,
      r: 50,
      b: 0,
    },
  }

  if (ready) {
    plotly.react("intermediate-values-plot", plotData, layout)
  }
  return <div id="intermediate-values-plot" />
}
