import React, { ReactNode, useState } from "react"
import {
  Dialog,
  DialogTitle,
  DialogContent,
  DialogContentText,
  Button,
  DialogActions,
} from "@mui/material"
import { actionCreator } from "../action"
import { DebouncedInputTextField } from "./Debounce"
import { useRecoilValue } from "recoil"
import { studySummariesState } from "../state"

export const useCreateStudyDialog = (): [() => void, () => ReactNode] => {
  const action = actionCreator()

  const [newStudyName, setNewStudyName] = useState("")
  const [openNewStudyDialog, setOpenNewStudyDialog] = useState(false)
  const [direction, setDirection] = useState<StudyDirection>("minimize")
  const studies = useRecoilValue<StudySummary[]>(studySummariesState)
  const newStudyNameAlreadyUsed = studies.some(
    (v) => v.study_name === newStudyName
  )

  const handleCloseNewStudyDialog = () => {
    setOpenNewStudyDialog(false)
    setNewStudyName("")
    setDirection("minimize")
  }

  const handleCreateNewStudy = () => {
    action.createNewStudy(newStudyName, direction)
    setOpenNewStudyDialog(false)
    setNewStudyName("")
    setDirection("minimize")
  }

  const openDialog = () => {
    setOpenNewStudyDialog(true)
  }

  const renderCreateNewStudyDialog = () => {
    return (
      <Dialog
        open={openNewStudyDialog}
        onClose={() => {
          handleCloseNewStudyDialog()
        }}
        aria-labelledby="create-study-dialog-title"
      >
        <DialogTitle id="create-study-dialog-title">New Study</DialogTitle>
        <DialogContent>
          <DialogContentText>
            Please enter the study name and directions here.
          </DialogContentText>
          <DebouncedInputTextField
            onChange={(s) => {
              setNewStudyName(s)
            }}
            delay={500}
            textFieldProps={{
              autoFocus: true,
              fullWidth: true,
              error: newStudyNameAlreadyUsed,
              helperText: newStudyNameAlreadyUsed
                ? `"${newStudyName}" is already used`
                : "",
              label: "Study name",
              type: "text",
            }}
          />
        </DialogContent>
        <DialogActions>
          <Button onClick={handleCloseNewStudyDialog} color="primary">
            Cancel
          </Button>
          <Button
            onClick={handleCreateNewStudy}
            color="primary"
            disabled={newStudyName === "" || newStudyNameAlreadyUsed}
          >
            Create
          </Button>
        </DialogActions>
      </Dialog>
    )
  }
  return [openDialog, renderCreateNewStudyDialog]
}
