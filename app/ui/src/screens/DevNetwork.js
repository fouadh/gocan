import { useEffect, useState } from 'react'
import axios from 'axios'
import { Spinner } from '../components/Spinner'
import { DateSelector } from '../components/DateSelector'
import { Button } from 'primereact/button'
import { Network } from '../components/Network'

export function DevNetwork({ sceneId, appId, date }) {
  const [dateRange, setDateRange] = useState(date)
  const [analyze, setAnalyze] = useState(true)
  const [network, setNetwork] = useState()
  const [loading, setLoading] = useState(false)

  useEffect(() => {
    let subscribed = true
    if (analyze) {
      setLoading(true)
    } else {
      return
    }

    const endpoint = `/api/scenes/${sceneId}/apps/${appId}/developers-network`
    let params = new URLSearchParams()
    if (dateRange.min) {
      params.append('after', dateRange.min)
    }
    if (dateRange.max) {
      params.append('before', dateRange.max)
    }
    axios
      .get(`${endpoint}?${params}`)
      .then((it) => it.data)
      .then((it) => {
        if (subscribed) {
          setNetwork(it)
        }
      })
      .finally(() => {
        setLoading(false)
        setAnalyze(false)
      })

    return () => (subscribed = false)
  }, [sceneId, appId, analyze, date, dateRange])

  let screen

  if (loading) {
    screen = <Spinner />
  } else if (network) {
    screen = (
      <div style={{ display: 'flex', justifyContent: 'center' }}>
        <Network width={975} height={975} data={network} />
      </div>
    )
  } else {
    screen = (
      <>
        <p>Unable to get developers network</p>
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
          <Button label="Submit" onClick={(e) => setAnalyze(true)} />
        </div>
      </div>
      {screen}
    </>
  )
}
