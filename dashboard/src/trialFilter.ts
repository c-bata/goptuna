import { useMemo, useState } from "react"

type TargetKind = "objective" | "user_attr" | "params"

export class Target {
  kind: TargetKind
  key: string | null

  constructor(kind: TargetKind, key: string | null) {
    this.kind = kind
    this.key = key
  }

  validate(): boolean {
    if (this.kind === "objective") {
      if (this.key !== null) {
        return false
      }
    } else if (this.kind === "user_attr") {
      if (typeof this.key !== "string") {
        return false
      }
    } else if (this.kind === "params") {
      if (typeof this.key !== "string") {
        return false
      }
    }
    return true
  }

  identifier(): string {
    return `${this.kind}:${this.key}`
  }

  toLabel(): string {
    if (this.kind === "objective") {
      return `Objective`
    } else if (this.kind === "user_attr") {
      return `User Attribute ${this.key}`
    } else {
      return `Param ${this.key}`
    }
  }

  getTargetValue(trial: Trial): number | null {
    if (!this.validate()) {
      return null
    }
    if (this.kind === "objective") {
      if (trial.value === undefined) {
        return null
      }
      return trial.value
    } else if (this.kind === "user_attr") {
      const attr = trial.user_attrs.find((attr) => attr.key === this.key)
      if (attr === undefined) {
        return null
      }
      const value = Number(attr.value)
      if (value === undefined) {
        return null
      }
      return value
    } else if (this.kind === "params") {
      const param = trial.params.find((p) => p.name === this.key)
      if (param === undefined) {
        return null
      }
      return param.param_internal_value
    }
    return null
  }
}

const filterTrials = (
  study: StudyDetail | null,
  targets: Target[],
  filterPruned: boolean
): Trial[] => {
  if (study === null) {
    return []
  }
  return study.trials.filter((t) => {
    if (t.state !== "Complete" && t.state !== "Pruned") {
      return false
    }
    if (t.state === "Pruned" && filterPruned) {
      return false
    }
    return targets.every((target) => target.getTargetValue(t) !== null)
  })
}

export const useFilteredTrials = (
  study: StudyDetail | null,
  targets: Target[],
  filterPruned: boolean
): Trial[] =>
  useMemo<Trial[]>(() => {
    return filterTrials(study, targets, filterPruned)
  }, [study?.trials, targets, filterPruned])

export const useParamTargets = (
  searchSpace: SearchSpaceItem[]
): [Target[], Target | null, (ident: string) => void] => {
  const [selected, setTargetIdent] = useState<string>("")
  const targetList = useMemo<Target[]>(() => {
    const targets = searchSpace.map((s) => new Target("params", s.name))
    if (selected === "" && targets.length > 0)
      setTargetIdent(targets[0].identifier())
    return targets
  }, [searchSpace])
  const selectedTarget = useMemo<Target | null>(
    () => targetList.find((t) => t.identifier() === selected) || null,
    [targetList, selected]
  )
  return [targetList, selectedTarget, setTargetIdent]
}

export const useObjectiveAndUserAttrTargets = (
  study: StudyDetail | null
): [Target[], Target, (ident: string) => void] => {
  const defaultTarget = new Target("objective", null)
  const [selected, setTargetIdent] = useState<string>(
    defaultTarget.identifier()
  )
  const targetList = useMemo<Target[]>(() => {
    if (study !== null) {
      return [
        new Target("objective", null),
        ...study.union_user_attrs
          .filter((attr) => attr.sortable)
          .map((attr) => new Target("user_attr", attr.key)),
      ]
    } else {
      return [defaultTarget]
    }
  }, [study?.union_user_attrs])
  const selectedTarget = useMemo<Target>(
    () => targetList.find((t) => t.identifier() === selected) || defaultTarget,
    [targetList, selected]
  )
  return [targetList, selectedTarget, setTargetIdent]
}
