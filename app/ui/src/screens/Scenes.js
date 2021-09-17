import {useEffect, useState} from "react";
import axios from "axios";
import {Link} from 'react-router-dom';

function getScenes() {
  return axios
    .get('/api/scenes')
    .then(it => {
      return it.data;
    })
    .then(it => it.scenes);
}

export function Scenes() {
  const [scenes, setScenes] = useState([]);
  useEffect(() => {
    let subscribed = true;
    getScenes()
      .then(it => {
        if (subscribed)
          setScenes(it);
      });

    return (() => subscribed = false);
  }, []);

  return <div>
    {
      scenes.map((each) => <div key={each.name}>
        <Link to={`/scenes/${each.id}`}>{each.name}</Link>
      </div>)
    }
  </div>;
}