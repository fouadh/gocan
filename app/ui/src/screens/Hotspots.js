import {useEffect, useState} from "react";
import axios from "axios";
import {CirclePacking} from "../components/CirclePacking";

export function Hotspots({sceneId, appId}) {
  const [hospots, setHotspots] = useState([]);

  useEffect(() => {
    let subscribed = true;

    axios.get(`/api/scenes/${sceneId}/apps/${appId}/hotspots`)
      .then(it => it.data)
      .then(setHotspots)
    ;

    return () => subscribed = false;
  }, [sceneId, appId]);

  return <div>
    <CirclePacking width={975} height={975} data={hospots}/>
  </div>;
}