import {useEffect, useState} from "react";
import axios from "axios";
import {Dropdown} from "primereact/dropdown";

export function BoundarySelector({sceneId, appId, onChange}) {
  const [boundary, setBoundary] = useState();
  const [boundaries, setBoundaries] = useState([]);

  useEffect(() => {
    let subscribed = true;
    axios.get(`/api/scenes/${sceneId}/apps/${appId}/boundaries`)
      .then(it => it.data)
      .then(it => it.boundaries)
      .then((it) => {
        if (subscribed) {
            setBoundaries(it);
            if (it.length > 0) {
                setBoundary(it[0]);
                onChange({value: it[0]});
            }
        }
      });
    return () => subscribed = false;
  }, [sceneId, appId]);

  return (<><label className="p-mr-2">Boundary:</label>
    <Dropdown optionLabel="name"
              optionValue="id"
              options={boundaries}
              placeholder="Select a boundary"
              value={boundary}
              showClear={true}
              onChange={(e) => {
                setBoundary(e.value);
                onChange(e);
              }}/>
  </>);
}