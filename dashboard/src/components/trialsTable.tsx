import { jsx } from "@emotion/core"
import { FC } from "react"
import Table from "@material-ui/core/Table"
import TableBody from "@material-ui/core/TableBody"
import TableCell from "@material-ui/core/TableCell"
import TableContainer from "@material-ui/core/TableContainer"
import TableHead from "@material-ui/core/TableHead"
import TableRow from "@material-ui/core/TableRow"
import Paper from "@material-ui/core/Paper"
import { makeStyles } from "@material-ui/core/styles"

const useStyles = makeStyles({
  table: {
    minWidth: 650,
  },
})

export const TrialsTable: FC<{
  trials: Trial[]
}> = ({ trials = [] }) => {
  const classes = useStyles()
  return (
    <TableContainer component={Paper}>
      <Table className={classes.table} size="small" aria-label="trials table">
        <TableHead>
          <TableRow>
            <TableCell>Trial ID</TableCell>
            <TableCell>Number</TableCell>
            <TableCell>Value</TableCell>
          </TableRow>
        </TableHead>
        <TableBody>
          {trials.map((t) => {
            return (
              <TableRow key={t.trial_id}>
                <TableCell component="th" scope="row">
                  {t.trial_id}
                </TableCell>
                <TableCell>{t.number}</TableCell>
                <TableCell>{t.value}</TableCell>
              </TableRow>
            )
          })}
        </TableBody>
      </Table>
    </TableContainer>
  )
}
