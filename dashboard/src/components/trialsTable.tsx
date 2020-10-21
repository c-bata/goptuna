import React, { FC, useState } from "react"
import { makeStyles } from "@material-ui/core/styles"
import {
  TableFooter,
  TablePagination,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  Paper,
} from "@material-ui/core"
import TablePaginationActions from "@material-ui/core/TablePagination/TablePaginationActions"

const useStyles = makeStyles({
  table: {
    minWidth: 650,
  },
})

export const TrialsTable: FC<{
  trials: Trial[]
}> = ({ trials = [] }) => {
  const classes = useStyles()
  const [page, setPage] = useState<number>(0)
  const [rowsPerPage, setRowsPerPage] = useState<number>(10)

  const emptyRows =
    rowsPerPage - Math.min(rowsPerPage, trials.length - page * rowsPerPage)
  const handleChangePage = (
    event: React.MouseEvent<HTMLButtonElement> | null,
    newPage: number
  ) => {
    setPage(newPage)
  }

  const handleChangeRowsPerPage = (
    event: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>
  ) => {
    const value = parseInt(event.target.value, 10)
    setRowsPerPage(value)
    setPage(0)
  }

  return (
    <TableContainer component={Paper}>
      <Table className={classes.table} size="small" aria-label="trials table">
        <TableHead>
          <TableRow>
            <TableCell>Trial ID</TableCell>
            <TableCell>Number</TableCell>
            <TableCell>State</TableCell>
            <TableCell>Value</TableCell>
            <TableCell>Params</TableCell>
          </TableRow>
        </TableHead>
        <TableBody>
          {(rowsPerPage > 0
            ? trials.slice(page * rowsPerPage, page * rowsPerPage + rowsPerPage)
            : trials
          ).map((t) => {
            return (
              <TableRow key={t.trial_id}>
                <TableCell component="th" scope="row">
                  {t.trial_id}
                </TableCell>
                <TableCell>{t.number}</TableCell>
                <TableCell>{t.state.toString()}</TableCell>
                <TableCell>{t.value}</TableCell>
                <TableCell>
                  {t.params.map((p) => p.name + ": " + p.value).join(", ")}
                </TableCell>
              </TableRow>
            )
          })}
          {emptyRows > 0 && (
            <TableRow style={{ height: 53 * emptyRows }}>
              <TableCell colSpan={6} />
            </TableRow>
          )}
        </TableBody>
        <TableFooter>
          <TableRow>
            <TablePagination
              rowsPerPageOptions={[10, 50, 100, { label: "All", value: -1 }]}
              colSpan={3}
              count={trials.length}
              rowsPerPage={rowsPerPage}
              page={page}
              SelectProps={{
                inputProps: { "aria-label": "rows per page" },
                native: true,
              }}
              onChangePage={handleChangePage}
              onChangeRowsPerPage={handleChangeRowsPerPage}
              ActionsComponent={TablePaginationActions}
            />
          </TableRow>
        </TableFooter>
      </Table>
    </TableContainer>
  )
}
