import { useEffect, useState } from 'react'
import axios from 'axios'
import { Dropdown } from 'primereact/dropdown'
import './BoundarySelector.css'

export function BoundarySelector({ sceneId, appId, onChange }) {
  const [boundary, setBoundary] = useState()
  const [boundaries, setBoundaries] = useState([])

  useEffect(() => {
    let subscribed = true
    axios
      .get(`/api/scenes/${sceneId}/apps/${appId}/boundaries`)
      .then((it) => it.data)
      .then((it) => it.boundaries)
      .then((it) => {
        if (subscribed) {
          setBoundaries(it)
        }
      })
    return () => (subscribed = false)
  }, [sceneId, appId])

  let selector
  if (boundaries && boundaries.length > 0) {
    selector = (
      <>
        <>
          <div className="p-field p-col-12 p-md-4 mr-4">
            <span className="p-float-label">
              <Dropdown
                id="boundaries"
                optionLabel="name"
                optionValue="name"
                options={boundaries}
                value={boundary}
                showClear={true}
                onChange={(e) => {
                  setBoundary(e.value)
                  onChange(e)
                }}
              />
              <label htmlFor="boundaries">Boundary</label>
            </span>
          </div>
        </>
      </>
    )
  } else {
    selector = <></>
  }

  return <>{selector}</>
}
