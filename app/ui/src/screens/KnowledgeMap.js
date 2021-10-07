import {useEffect, useState} from "react";
import * as d3 from "d3";
import axios from "axios";
import {CirclePacking} from "../components/CirclePacking";
import {Spinner} from "../components/Spinner";

export function KnowledgeMap({sceneId, appId}) {
  const [knowledgeMap, setKnowledgeMap] = useState();
  const [authors, setAuthors] = useState([]);
  const [loading, setLoading] = useState(false);

  const authorColor = d3.scaleOrdinal(d3.schemeCategory10)
  authorColor.domain(authors.map((it) => it.name))

  function fillColor(d) {
    const color2 = d3
      .scaleLinear()
      .domain([-1, 5])
      .range(['hsl(185,60%,99%)', 'hsl(187,40%,70%)'])
      .interpolate(d3.interpolateHcl)
    return (d.data).weight || 0 > 0.0
      ? authorColor((d.data).mainDeveloper)
      : d.children
        ? color2(d.depth)
        : 'WhiteSmoke'
  }

  useEffect(() => {
    let subscribed = true;
    setLoading(true);

    axios.get(`/api/scenes/${sceneId}/apps/${appId}/developers`)
      .then(it => it.data)
      .then(it => it.authors)
      .then(it => {
        if (subscribed)
          setAuthors(it);
      });

    axios.get(`/api/scenes/${sceneId}/apps/${appId}/knowledge-map`)
      .then(it => it.data)
      .then(it => {
        if (subscribed) {
          console.log(it);
          setKnowledgeMap(it);
        }
      }).finally(() => setLoading(false));


    return (() => subscribed = false);
  }, [sceneId, appId]);

  let screen;

  if (loading) {
    screen = <Spinner/>;
  } else {
    screen = <CirclePacking
        width={975}
        height={975}
        data={knowledgeMap}
        fillColor={fillColor}
        fillOpacity={(d) => 1}
    />
  }

  return <>
    <div className="p-d-flex p-text-center">
      <div className="p-mr-5">
        <ul>
          {authors.map((a) =>
            (
              <li key={a.name}>
                {a.name}{' '}
                <span
                  style={{
                    display: 'inline-block',
                    width: '20px',
                    height: '20px',
                    backgroundColor: authorColor(a.name)
                  }}
                ></span>
              </li>
            ))}
        </ul>
      </div>
      <div className="p-ml-5 p-mr-4">
        {screen}
      </div>
    </div>
  </>;
}