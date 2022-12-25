import { useEffect, useState } from 'react'
import { InputNumber } from 'primereact/inputnumber'
import axios from 'axios'
import { Chord } from '../components/Chord'
import { Spinner } from '../components/Spinner'
import { DateSelector } from '../components/DateSelector'
import { Button } from 'primereact/button'
import { BoundarySelector } from './BoundarySelector'

export function Coupling({ sceneId, appId, date }) {
  const [minCouplingPercent, setMinCouplingPercent] = useState(40)
  const [minRevsAvg, setMinRevsAvg] = useState(6)
  const [boundaryName, setBoundaryName] = useState()
  const [dateRange, setDateRange] = useState({ date })
  const [analyze, setAnalyze] = useState(true)
  const [coupling, setCoupling] = useState()
  const [loading, setLoading] = useState(false)

  useEffect(() => {
    let subscribe = true
    if (analyze) {
      setLoading(true)
      let endpoint = `/api/scenes/${sceneId}/apps/${appId}/coupling-hierarchy`
      const params = new URLSearchParams()
      if (dateRange.min) {
        params.append('after', dateRange.min)
      }
      if (dateRange.max) {
        params.append('before', dateRange.max)
      }
      if (minCouplingPercent) {
        params.append('minCoupling', minCouplingPercent / 100)
      }
      if (minRevsAvg) {
        params.append('minRevisionsAvg', minRevsAvg)
      }
      if (boundaryName) {
        params.append('boundaryName', boundaryName)
      }

      axios
        .get(`${endpoint}?${params}`)
        .then((it) => it.data)
        .then((it) => {
          if (subscribe) {
            setCoupling(it)
          }
        })
        .finally(() => {
          setLoading(false)
          setAnalyze(false)
        })
    }

    return () => (subscribe = false)
  }, [
    sceneId,
    appId,
    analyze,
    dateRange,
    minCouplingPercent,
    minRevsAvg,
    boundaryName,
  ])

  let screen
  if (loading) {
    screen = <Spinner />
  } else if (coupling) {
    screen = (
      <div
        style={{ display: 'flex', justifyContent: 'center', height: '900px' }}
      >
        <Chord data={coupling} />
      </div>
    )
  } else {
    screen = (
      <>
        <p>No coupling found.</p>
      </>
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
          <div className="p-field p-col-12 p-md-4 mr-4">
            <span className="p-float-label">
              <InputNumber
                id="minCouplingPercent"
                value={minCouplingPercent}
                onValueChange={(e) => setMinCouplingPercent(e.value)}
              />
              <label htmlFor="minCouplingPercent">
                Minimal Coupling Percentage
              </label>
            </span>
          </div>
          <div className="p-field p-col-12 p-md-4 mr-4">
            <span className="p-float-label">
              <InputNumber
                id="minRevsAvg"
                value={minRevsAvg}
                onValueChange={(e) => setMinRevsAvg(e.value)}
              />
              <label htmlFor="minRevsAvg">Minimal Revisions Average</label>
            </span>
          </div>
          <BoundarySelector
            appId={appId}
            sceneId={sceneId}
            onChange={(e) => setBoundaryName(e.value)}
          />
          <Button label="Submit" onClick={(e) => setAnalyze(true)} />
        </div>
      </div>
      {screen}
    </>
  )
}
