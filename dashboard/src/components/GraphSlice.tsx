import * as plotly from "plotly.js-dist"
import React, { ChangeEvent, FC, useEffect, useState } from "react"
import {
  Grid,
  FormControl,
  InputLabel,
  MenuItem,
  Select,
} from "@material-ui/core"
import { createStyles, makeStyles, Theme } from "@material-ui/core/styles"

const plotDomId = "graph-slice"

const useStyles = makeStyles((theme: Theme) =>
  createStyles({
    formControl: {
      marginBottom: theme.spacing(2),
      marginRight: theme.spacing(5),
      marginTop: theme.spacing(10),
    },
  })
)

const getParamNames = (trials: Trial[]): string[] => {
  const paramSet = new Set<string>(
    ...trials.map<string[]>((t) => t.params.map((p) => p.name))
  )
  return Array.from(paramSet)
}

export const GraphSlice: FC<{
  trials: Trial[]
}> = ({ trials = [] }) => {
  const classes = useStyles()
  const [paramNames, setParamNames] = useState<string[]>([])
  const [selected, setSelected] = useState<string | null>(null)

  useEffect(() => {
    if (trials.length === 0 || paramNames.length !== 0) {
      return
    }

    const p = getParamNames(trials)
    setParamNames(p)
    if (selected === null && p.length !== 0) {
      setSelected(p[0])
    }
  }, [trials])

  useEffect(() => {
    plotSlice(trials, selected)
  }, [trials, selected])

  const handleSelectedParam = (e: ChangeEvent<{ value: unknown }>) => {
    setSelected(e.target.value as string)
  }

  return (
    <Grid container direction="row">
      <Grid item xs={3}>
        <Grid container direction="column">
          {paramNames.length !== 0 && selected !== null ? (
            <FormControl component="fieldset" className={classes.formControl}>
              <InputLabel id="parameter">Parameter</InputLabel>
              <Select value={selected} onChange={handleSelectedParam}>
                {paramNames.map((p, i) => (
                  <MenuItem value={p} key={i}>
                    {p}
                  </MenuItem>
                ))}
              </Select>
            </FormControl>
          ) : null}
        </Grid>
      </Grid>
      <Grid item xs={6}>
        <div id={plotDomId} />
      </Grid>
    </Grid>
  )
}

const plotSlice = (trials: Trial[], selected: string | null) => {
  if (document.getElementById(plotDomId) === null) {
    return
  }

  const layout: Partial<plotly.Layout> = {
    title: "Slice",
    margin: {
      l: 50,
      r: 50,
    },
    xaxis: {
      title: selected || "",
      zerolinecolor: "#f2f5fa",
      zerolinewidth: 1.5,
      linecolor: "#f2f5fa",
      linewidth: 5,
      gridcolor: "#f2f5fa",
      gridwidth: 1,
    },
    yaxis: {
      title: "Objective Values",
      zerolinecolor: "#f2f5fa",
      zerolinewidth: 2,
      linecolor: "#f2f5fa",
      linewidth: 5,
      gridcolor: "#f2f5fa",
      gridwidth: 1,
    },
    plot_bgcolor: "#E5ecf6",
    showlegend: false,
  }

  const filteredTrials = trials.filter(
    (t) =>
      (t.state === "Complete" || t.state === "Pruned") &&
      t.params.find((p) => p.name == selected) !== undefined
  )

  if (filteredTrials.length === 0 || selected === null) {
    plotly.react(plotDomId, [], layout)
    return
  }

  const objectiveValues: number[] = filteredTrials.map((t) => t.value!)
  const valueStrings = filteredTrials.map((t) => {
    return t.params.find((p) => p.name == selected)!.value
  })

  const isnum = valueStrings.every((v) => {
    return !isNaN(parseFloat(v))
  })
  if (isnum) {
    const valuesNum: number[] = valueStrings.map((v) => parseFloat(v))
    const trace: plotly.Data[] = [
      {
        type: "scatter",
        x: valuesNum,
        y: objectiveValues,
        mode: "markers",
        xaxis: selected,
        marker: {
          color: "#185799",
        },
      },
    ]
    layout["xaxis"] = {
      title: selected,
      zerolinecolor: "#f2f5fa",
      zerolinewidth: 1.5,
      linecolor: "#f2f5fa",
      linewidth: 5,
      gridcolor: "#f2f5fa",
      gridwidth: 1,
    }
    plotly.react(plotDomId, trace, layout)
  } else {
    const vocabSet = new Set<string>(valueStrings)
    const vocabArr = Array.from<string>(vocabSet)
    const valuesCategorical: number[] = valueStrings.map((v) =>
      vocabArr.findIndex((vocab) => v === vocab)
    )
    const tickvals: number[] = vocabArr.map((v, i) => i)
    const trace: plotly.Data[] = [
      {
        type: "scatter",
        x: valuesCategorical,
        y: objectiveValues,
        mode: "markers",
        // xaxis: paramName,
        marker: {
          color: "#185799",
        },
      },
    ]
    layout["xaxis"] = {
      title: selected,
      zerolinecolor: "#f2f5fa",
      zerolinewidth: 1.5,
      linecolor: "#f2f5fa",
      linewidth: 5,
      gridcolor: "#f2f5fa",
      gridwidth: 1,
      tickfont: {
        color: "#000000",
      },
      tickvals: tickvals,
      ticktext: vocabArr,
    }
    plotly.react(plotDomId, trace, layout)
  }
}
