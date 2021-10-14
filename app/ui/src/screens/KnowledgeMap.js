import {useEffect, useState} from "react";
import * as d3 from "d3";
import axios from "axios";
import {CirclePacking} from "../components/CirclePacking";
import {Spinner} from "../components/Spinner";
import {DateSelector} from "../components/DateSelector";

export function KnowledgeMap({sceneId, appId}) {
  const [dateRange, setDateRange] = useState({});
  const [analyze, setAnalyze] = useState(true);
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

    if (analyze) {
      setLoading(true);
      let endpoint1 = `/api/scenes/${sceneId}/apps/${appId}/developers`;
      if (dateRange.min) {
        if (dateRange.max) {
          endpoint1 += `?after=${dateRange.min}&before=${dateRange.max}`
        } else {
          endpoint1 += `?after=${dateRange.min}`
        }
      } else if (dateRange.max) {
        endpoint1 += `?before=${dateRange.max}`
      }
      axios.get(endpoint1)
          .then(it => it.data)
          .then(it => it.authors)
          .then(it => {
            if (subscribed)
              setAuthors(it);
          });

      let endpoint2 = `/api/scenes/${sceneId}/apps/${appId}/knowledge-map`;
      if (dateRange.min) {
        if (dateRange.max) {
          endpoint2 += `?after=${dateRange.min}&before=${dateRange.max}`
        } else {
          endpoint2 += `?after=${dateRange.min}`
        }
      } else if (dateRange.max) {
        endpoint2 += `?before=${dateRange.max}`
      }
      axios.get(endpoint2)
          .then(it => it.data)
          .then(it => {
            if (subscribed) {
              console.log(it);
              setKnowledgeMap(it);
            }
          }).finally(() => {
        setLoading(false);
        setAnalyze(false);
      });
    }

    return (() => subscribed = false);
  }, [sceneId, appId, analyze]);

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
    <DateSelector onChange={e => setDateRange(e)}/>
    <button onClick={e => setAnalyze(true)}>Submit</button>
    <div className="p-d-flex p-text-center" style={{display: "flex", justifyContent: "space-around", alignItems: "center"}}>
      <div className="p-mr-5">
        <ul>
          {authors.map((a) =>
            (
              <li key={a.name} style={{listStyleType: 'none'}}>
                <span
                    style={{
                      display: 'inline-block',
                      width: '20px',
                      height: '20px',
                      marginRight: '.5em',
                      backgroundColor: authorColor(a.name)
                    }}
                ></span>
                <span>{a.name}</span>
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