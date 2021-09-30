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
                }
            });
        return () => subscribed = false;
    }, [sceneId, appId]);

    let selector;
    if (boundaries && boundaries.length > 0) {
        selector = <>
            <>
                <label className="p-mr-2">Boundary:</label>
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
            </>
        </>
    } else {
        selector = <>
            <p>Please use the <strong> gocan create-boundary </strong> command to create selectable boundaries.</p>
        </>
    }

    return <>{selector}</>;
}