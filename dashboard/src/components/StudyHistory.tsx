import React, { FC, useState } from "react"
import {
  Box,
  Card,
  CardContent,
  FormControl,
  Switch,
  Typography,
  useTheme,
} from "@mui/material"
import { GraphHistory } from "./GraphHistory"
import { GraphIntermediateValues } from "./GraphIntermediateValues"
import Grid2 from "@mui/material/Unstable_Grid2"
import { DataGrid, DataGridColumn } from "./DataGrid"
import { useStudyDetailValue, useStudySummaryValue } from "../state"
import FormControlLabel from "@mui/material/FormControlLabel"

export const StudyHistory: FC<{ studyId: number }> = ({ studyId }) => {
  const theme = useTheme()
  const studySummary = useStudySummaryValue(studyId)
  const studyDetail = useStudyDetailValue(studyId)
  const [logScale, setLogScale] = useState<boolean>(false)
  const [includePruned, setIncludePruned] = useState<boolean>(true)

  const handleLogScaleChange = () => {
    setLogScale(!logScale)
  }

  const handleIncludePrunedChange = () => {
    setIncludePruned(!includePruned)
  }

  const userAttrs = studySummary?.user_attrs || studyDetail?.user_attrs || []
  const userAttrColumns: DataGridColumn<Attribute>[] = [
    { field: "key", label: "Key", sortable: true },
    { field: "value", label: "Value", sortable: true },
  ]
  const trials: Trial[] = studyDetail?.trials || []
  return (
    <Box sx={{ display: "flex", width: "100%", flexDirection: "column" }}>
      <FormControl
        component="fieldset"
        sx={{
          display: "flex",
          flexDirection: "row",
          justifyContent: "flex-end",
          padding: theme.spacing(2),
        }}
      >
        <FormControlLabel
          control={
            <Switch
              checked={logScale}
              onChange={handleLogScaleChange}
              value="enable"
            />
          }
          label="Log y scale"
        />
        <FormControlLabel
          control={
            <Switch
              checked={
                studyDetail
                  ? studyDetail.has_intermediate_values && includePruned
                  : false
              }
              onChange={handleIncludePrunedChange}
              disabled={!studyDetail?.has_intermediate_values}
              value="enable"
            />
          }
          label="Include PRUNED trials"
        />
      </FormControl>
      <Card
        sx={{
          margin: theme.spacing(2),
        }}
      >
        <CardContent>
          <GraphHistory
            study={studyDetail}
            includePruned={includePruned}
            logScale={logScale}
          />
        </CardContent>
      </Card>
      <Grid2 container spacing={2} sx={{ padding: theme.spacing(0, 2) }}>
        {studyDetail !== null && studyDetail.has_intermediate_values ? (
          <Grid2 xs={6}>
            <GraphIntermediateValues
              trials={trials}
              includePruned={includePruned}
              logScale={logScale}
            />
          </Grid2>
        ) : null}
        <Grid2 xs={6}>
          <Card>
            <CardContent
              sx={{
                display: "flex",
                flexDirection: "column",
              }}
            >
              <Typography
                variant="h6"
                sx={{
                  margin: "1em 0",
                  fontWeight: theme.typography.fontWeightBold,
                }}
              >
                Study User Attributes
              </Typography>
              <DataGrid<Attribute>
                columns={userAttrColumns}
                rows={userAttrs}
                keyField={"key"}
                dense={true}
                initialRowsPerPage={5}
                rowsPerPageOption={[5, 10, { label: "All", value: -1 }]}
              />
            </CardContent>
          </Card>
        </Grid2>
      </Grid2>
    </Box>
  )
}
