import { jsx, css } from "@emotion/core"
import { FC } from "react"
import {
  useParams
} from "react-router-dom";

interface ParamTypes {
  studyId: string
}

const style = css`

`

export const HistoryGraph: FC<{}> = () => {
  let { studyId } = useParams<ParamTypes>();
  console.log("デバッグ")
  console.log(studyId)
  return (
    <div css={style}>
      {studyId}
    </div>
  )
}
