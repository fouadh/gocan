import { useEffect, useState } from 'react'
import axios from 'axios'
import { CirclePacking } from '../components/CirclePacking'
import { Spinner } from '../components/Spinner'
import { DateSelector } from '../components/DateSelector'
import { Button } from 'primereact/button'
import * as d3 from 'd3'

export function CoordinationNeeds({ sceneId, appId, date }) {
  const [dateRange, setDateRange] = useState(date)
  const [analyze, setAnalyze] = useState(true)
  const [knowledgeMap, setKnowledgeMap] = useState()
  const [loading, setLoading] = useState(false)

  const color = d3
    .scaleLinear()
    .domain([-1, 5])
    .range(['hsl(185,60%,99%)', 'hsl(187,40%,70%)'])
    .interpolate(d3.interpolateHcl)

  const devDiffusionColor = (diffusion) => {
    switch (diffusion) {
      case 0.25:
        return 'rgb(254, 240, 217)'

      case 0.5:
        return 'rgb(253, 204, 138)'

      case 0.75:
        return 'rgb(252, 141, 89)'

      case 1.0:
        return 'rgb(215, 48, 31)'

      default:
        return 'WhiteSmoke'
    }
  }

  function fillColor(d) {
    return d.children ? color(d.depth) : devDiffusionColor(d.data.devDiffusion)
  }

  function fillOpacity(d) {
    return 1
  }

  useEffect(() => {
    let subscribed = true

    if (analyze) {
      setLoading(true)
      let endpoint = `/api/scenes/${sceneId}/apps/${appId}/knowledge-map`
      const params = new URLSearchParams()
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
          if (subscribed) setKnowledgeMap(it)
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
      <CirclePacking
        width={975}
        height={975}
        data={knowledgeMap}
        fillColor={fillColor}
        fillOpacity={fillOpacity}
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
      <div
        className="p-d-flex p-text-center"
        style={{
          display: 'flex',
          justifyContent: 'space-around',
          alignItems: 'center',
        }}
      >
        <div className="p-ml-5 p-mr-4">{screen}</div>
      </div>
    </>
  )
}
