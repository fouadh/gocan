import { useEffect, useState } from 'react'
import axios from 'axios'
import { MultiLineChart } from '../components/MultiLineChart'
import * as d3 from 'd3'

export function ActiveSet({ sceneId, appId }) {
  const [activeSet, setActiveSet] = useState([])

  useEffect(() => {
    let subscribe = true

    axios
      .get(`/api/scenes/${sceneId}/apps/${appId}/active-set`)
      .then((it) => it.data)
      .then((it) => it.activeSet)
      .then((it) => {
        if (subscribe) {
          const opened = []
          const closed = []

          it.forEach((each) => {
            opened.push({ x: each.date, y: each.opened })
            closed.push({ x: each.date, y: each.closed })
          })

          setActiveSet([opened, closed])
        }
      })

    return () => (subscribe = false)
  }, [sceneId, appId])

  return activeSet.flat().length > 0 ? (
    <MultiLineChart
      data={activeSet}
      legend={['Opened Entities', 'Closed Entities']}
      xFormatter={d3.timeFormat('%Y-%m-%d')}
      xAccessor={(d) => d3.utcParse('%Y-%m-%d')(d.x)}
      yAccessor={(d) => (d.y ? d.y : 0)}
      yLabel="Active Set Entities"
    />
  ) : (
    <></>
  )
}
