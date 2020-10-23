import React from "react"
import { createStyles, makeStyles, Theme } from "@material-ui/core/styles"
import Table from "@material-ui/core/Table"
import TableBody from "@material-ui/core/TableBody"
import TableCell from "@material-ui/core/TableCell"
import TableContainer from "@material-ui/core/TableContainer"
import TableHead from "@material-ui/core/TableHead"
import TablePagination from "@material-ui/core/TablePagination"
import TableRow from "@material-ui/core/TableRow"
import TableSortLabel from "@material-ui/core/TableSortLabel"

const useStyles = makeStyles((theme: Theme) =>
  createStyles({
    root: {
      width: "100%",
    },
    table: {
      minWidth: 750,
    },
    visuallyHidden: {
      border: 0,
      clip: "rect(0 0 0 0)",
      height: 1,
      margin: -1,
      overflow: "hidden",
      padding: 0,
      position: "absolute",
      top: 20,
      width: 1,
    },
  })
)

function descendingComparator<T>(a: T, b: T, orderBy: keyof T) {
  if (b[orderBy] < a[orderBy]) {
    return -1
  }
  if (b[orderBy] > a[orderBy]) {
    return 1
  }
  return 0
}

type Order = "asc" | "desc"

function getComparator<T>(
  order: Order,
  orderBy: keyof T
): (a: T, b: T) => number {
  return order === "desc"
    ? (a, b) => descendingComparator<T>(a, b, orderBy)
    : (a, b) => -descendingComparator<T>(a, b, orderBy)
}

function stableSort<T>(array: T[], comparator: (a: T, b: T) => number) {
  const stabilizedThis = array.map((el, index) => [el, index] as [T, number])
  stabilizedThis.sort((a, b) => {
    const order = comparator(a[0], b[0])
    if (order !== 0) return order
    return a[1] - b[1]
  })
  return stabilizedThis.map((el) => el[0])
}

export interface DataGridColumn<T> {
  field: keyof T
  label: string
  sortable: boolean
  toCellValue?: (dataIndex: number) => string | React.ReactNode
}

export function DataGrid<T>(props: {
  columns: DataGridColumn<T>[]
  rows: T[]
  keyField: keyof T
  dense?: boolean
}) {
  const classes = useStyles()
  const { columns, rows, keyField, dense } = props
  const [order, setOrder] = React.useState<Order>("asc")
  const [orderBy, setOrderBy] = React.useState<keyof T>(keyField)
  const [page, setPage] = React.useState(0)
  const [rowsPerPage, setRowsPerPage] = React.useState(10)

  const handleRequestSort = (
    event: React.MouseEvent<unknown>,
    property: keyof T
  ) => {
    const isAsc = orderBy === property && order === "asc"
    setOrder(isAsc ? "desc" : "asc")
    setOrderBy(property)
  }
  const createSortHandler = (property: keyof T) => (
    event: React.MouseEvent<unknown>
  ) => {
    handleRequestSort(event, property)
  }

  const handleChangePage = (event: unknown, newPage: number) => {
    setPage(newPage)
  }

  const handleChangeRowsPerPage = (
    event: React.ChangeEvent<HTMLInputElement>
  ) => {
    setRowsPerPage(parseInt(event.target.value, 10))
    setPage(0)
  }

  const emptyRows =
    rowsPerPage - Math.min(rowsPerPage, rows.length - page * rowsPerPage)

  const sortedRows = stableSort<T>(rows, getComparator(order, orderBy))
  const paginateRows =
    rowsPerPage > 0
      ? sortedRows.slice(page * rowsPerPage, page * rowsPerPage + rowsPerPage)
      : sortedRows

  return (
    <div className={classes.root}>
      <TableContainer>
        <Table
          className={classes.table}
          aria-labelledby="tableTitle"
          size={dense ? "small" : "medium"}
          aria-label="enhanced table"
        >
          <TableHead>
            <TableRow>
              {columns.map((column, index) => (
                <TableCell
                  key={index}
                  sortDirection={orderBy === column.field ? order : false}
                >
                  {column.sortable ? (
                    <TableSortLabel
                      active={orderBy === column.field}
                      direction={orderBy === column.field ? order : "asc"}
                      onClick={createSortHandler(column.field)}
                    >
                      {column.label}
                      {orderBy === column.field ? (
                        <span className={classes.visuallyHidden}>
                          {order === "desc"
                            ? "sorted descending"
                            : "sorted ascending"}
                        </span>
                      ) : null}
                    </TableSortLabel>
                  ) : (
                    column.label
                  )}
                </TableCell>
              ))}
            </TableRow>
          </TableHead>
          <TableBody>
            {paginateRows.map((row, index) => {
              return (
                <TableRow hover role="checkbox" tabIndex={-1} key={index}>
                  {columns.map((column) => (
                    <TableCell key={`${row[keyField]}:${column.field}`}>
                      {column.toCellValue
                        ? column.toCellValue(index)
                        : row[column.field]}
                    </TableCell>
                  ))}
                </TableRow>
              )
            })}
            {emptyRows > 0 && (
              <TableRow style={{ height: (dense ? 33 : 53) * emptyRows }}>
                <TableCell colSpan={6} />
              </TableRow>
            )}
          </TableBody>
        </Table>
      </TableContainer>
      <TablePagination
        rowsPerPageOptions={[10, 50, 100, { label: "All", value: -1 }]}
        component="div"
        count={rows.length}
        rowsPerPage={rowsPerPage}
        page={page}
        onChangePage={handleChangePage}
        onChangeRowsPerPage={handleChangeRowsPerPage}
      />
    </div>
  )
}
