import { useEffect, useState } from 'react'
import axios from 'axios'
import { CirclePacking } from '../components/CirclePacking'
import { Spinner } from '../components/Spinner'
import { DateSelector } from '../components/DateSelector'
import { Button } from 'primereact/button'
import { Autocomplete } from '../components/Autocomplete'
import { InputNumber } from 'primereact/inputnumber'

export function EntityCoupling({ sceneId, appId, date, entities }) {
  const [minCouplingPercent, setMinCouplingPercent] = useState(40)
  const [minRevsAvg, setMinRevsAvg] = useState(6)
  const [entity, setEntity] = useState('')
  const [dateRange, setDateRange] = useState(date)
  const [analyze, setAnalyze] = useState(true)
  const [error, setError] = useState()
  const [coupling, setCoupling] = useState()
  const [loading, setLoading] = useState(false)

  useEffect(() => {
    let subscribed = true
    if (analyze) {
      if (entity !== '') {
        setError(null)
        setLoading(true)
        let endpoint = `/api/scenes/${sceneId}/apps/${appId}/entity-coupling`
        const params = new URLSearchParams()
        params.append('entity', entity)
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

        axios
          .get(`${endpoint}?${params}`)
          .then((it) => it.data)
          .then((it) => {
            if (subscribed) {
              setCoupling(it)
            }
          })
          .catch(() => setError('Unable to get coupling information'))
          .finally(() => {
            setLoading(false)
            setAnalyze(false)
          })
      } else {
        setCoupling(null)
        setAnalyze(false)
      }
    }

    return () => (subscribed = false)
  }, [
    sceneId,
    appId,
    analyze,
    dateRange,
    entity,
    minCouplingPercent,
    minRevsAvg,
  ])

  let screen

  if (loading) {
    screen = <Spinner />
  } else if (coupling) {
    screen = (
      <div style={{ display: 'flex', justifyContent: 'center' }}>
        <CirclePacking width={975} height={975} data={coupling} />
      </div>
    )
  } else if (error) {
    screen = <p>{error}</p>
  }

  return (
    <>
      <div className="card mt-4">
        <div className="flex align-items-start">
          <div className="p-field p-col-12 p-md-4 mr-4">
            <span className="p-float-label autocomplete">
              <Autocomplete
                suggestions={entities}
                onChange={(e) => setEntity(e.value)}
              />
              <label htmlFor="entity">Entity</label>
            </span>
          </div>
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
          <Button label="Submit" onClick={(e) => setAnalyze(true)} />
        </div>
      </div>
      {screen}
    </>
  )
}
