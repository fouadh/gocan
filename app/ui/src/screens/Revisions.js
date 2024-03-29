import { useEffect, useState } from 'react'
import axios from 'axios'
import { Timeline } from '../components/Timeline'
import { Spinner } from '../components/Spinner'
import { DateSelector } from '../components/DateSelector'
import { Button } from 'primereact/button'

export function Revisions({ sceneId, appId, date }) {
  const [dateRange, setDateRange] = useState(date)
  const [analyze, setAnalyze] = useState(true)
  const [revisions, setRevisions] = useState([])
  const [loading, setLoading] = useState(false)

  useEffect(() => {
    let subscribed = true
    if (analyze) {
      setLoading(true)
      let endpoint = `/api/scenes/${sceneId}/apps/${appId}/revisions`
      if (dateRange.min) {
        if (dateRange.max) {
          endpoint += `?after=${dateRange.min}&before=${dateRange.max}`
        } else {
          endpoint += `?after=${dateRange.min}`
        }
      } else if (dateRange.max) {
        endpoint += `?before=${dateRange.max}`
      }
      axios
        .get(endpoint)
        .then((it) => it.data)
        .then((it) => it.revisions)
        .then((it) => {
          const results = []
          it.forEach((each, i) =>
            results.push({ x: i, y: each.numberOfRevisions })
          )
          return results
        })
        .then((it) => {
          if (subscribed) {
            setRevisions(it)
          }
        })
        .finally(() => {
          setLoading(false)
          setAnalyze(false)
        })
    }

    return () => (subscribed = false)
  }, [sceneId, appId, analyze, dateRange])

  let screen
  if (loading) {
    screen = <Spinner />
  } else {
    screen = (
      <Timeline
        data={revisions}
        xAccessor={(d) => {
          return d.x
        }}
        yAccessor={(d) => d.y}
        label="Revisions"
        xFormatter={(d) => d.y}
      />
    )
  }

  return (
    <>
      <div className="card mt-4">
        <div className="flex align-items-center">
          <DateSelector
            min={date.min}
            max={date.max}
            onChange={(e) => setDateRange(e)}
          />
          <Button label="Submit" onClick={(e) => setAnalyze(true)} />
        </div>
      </div>
      {screen}
    </>
  )
}
