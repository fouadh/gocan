import {useEffect, useState} from "react";
import axios from "axios";
import {CirclePacking} from "../components/CirclePacking";
import {Spinner} from "../components/Spinner";
import {DateSelector} from "../components/DateSelector";
import {Button} from 'primereact/button';
import {BoundarySelector} from "./BoundarySelector";
import {ModuleSelector} from "./ModuleSelector";
import {Slider} from 'primereact/slider'
import './Hotspots.css'
import * as _ from 'lodash'

function revisionsToHotspots(appName, revisions, minRevs) {
  const hierarchies = {}
  const filteredRevisions = revisions.filter(rev => rev.code > 0 && rev.numberOfRevisions >= minRevs)
  const revsRange = [_.minBy(revisions, rev => rev.numberOfRevisions).numberOfRevisions, _.maxBy(revisions, rev => rev.numberOfRevisions).numberOfRevisions]
  const maxNumberOfRevs = revsRange[1]
  const codeRange = [_.minBy(revisions, rev => rev.code).code, _.maxBy(revisions, rev => rev.code).code]

  filteredRevisions.forEach(rev => {
    const filename = rev.entity
    const pathElements = filename.split('/')

    let currentFolder = null
    pathElements.forEach((element, index) => {
      if (index === pathElements.length - 1) {
        // last element is a file
        const fileInfo = {
          name: element,
          size: rev.code,
          weight: rev.numberOfRevisions / maxNumberOfRevs
        }
        if (currentFolder) {
          currentFolder.children.push(fileInfo)
        } else {
          hierarchies[element] = fileInfo
        }
      } else {
        // still a folder
        if (currentFolder) {
          const existingElement = _.find(currentFolder.children, folder => folder.name === element)
          if (existingElement) {
            currentFolder = existingElement
          } else {
            const folder = {
              name: element,
              children: []
            }
            currentFolder.children.push(folder)
            currentFolder = folder
          }

        } else {
          // we are at the top folder
          if (hierarchies[element]) {
            currentFolder = hierarchies[element]
          } else {
            currentFolder = {
              name: element,
              children: []
            }
            hierarchies[element] = currentFolder
          }
        }
      }
    })
  })

  return {
    name: appName,
    maxRevisions: maxNumberOfRevs,
    codeRange: codeRange,
    revsRange: revsRange,
    children: Object.values(hierarchies)
  }
}

export function Hotspots({appName, sceneId, appId, date}) {
  const [dateRange, setDateRange] = useState(date);
  const [boundaryName, setBoundaryName] = useState();
  const [moduleName, setModuleName] = useState();
  const [analyze, setAnalyze] = useState(true);
  const [hospots, setHotspots] = useState();
  const [loading, setLoading] = useState(false);
  const [minRevs, setMinRevs] = useState(0);
  const [revisions, setRevisions] = useState();
  const [revsRange, setRevsRange] = useState([0, 1]);

  useEffect(() => {
    let subscribed = true;
    if (analyze) {
      setLoading(true);
      let endpoint;
      if (appId) {
        endpoint = `/api/scenes/${sceneId}/apps/${appId}/revisions`;
      } else {
        endpoint = `/api/scenes/${sceneId}/revisions`;
      }
      let params = new URLSearchParams();
      if (dateRange.min) {
        params.append("after", dateRange.min);
      }
      if (dateRange.max) {
        params.append("before", dateRange.max);
      }
      if (boundaryName) {
        params.append("boundaryName", boundaryName);
      }
      if (moduleName) {
        params.append("moduleName", moduleName);
      }
      axios.get(`${endpoint}?${params}`)
          .then(it => it.data)
          .then(revisions => setRevisions(revisions.revisions))
          .finally(() => {
            setLoading(false);
            setAnalyze(false);
          });
    }

    return () => subscribed = false;
  }, [sceneId, appId, appName, analyze, dateRange, boundaryName, moduleName]);

  useEffect(() => {
    if (revisions) {
      const hotspots = revisionsToHotspots(appName, revisions, minRevs)
      setRevsRange(hotspots.revsRange)
      setHotspots(hotspots)
    }
  }, [revisions, minRevs])

  let screen;

  if (loading) {
    screen = <Spinner/>;
  } else if (hospots) {
    screen = <div style={{display: "flex", justifyContent: "center"}}>
      <CirclePacking width={975} height={975} data={hospots}/>
    </div>;
  } else {
    screen = <><p>Unable to get hotspots</p></>
  }

  return <>
    <div className="card mt-4">
      <div className="flex align-items-center">
        <DateSelector min={date.min} max={date.max} onChange={e => setDateRange(e)}/>
        {appId && <BoundarySelector sceneId={sceneId} appId={appId} onChange={e => setBoundaryName(e.value)}/>}
        {boundaryName && <ModuleSelector sceneId={sceneId} appId={appId} boundaryName={boundaryName}
                                         onChange={e => setModuleName(e.value)}/>}
        <Button label="Submit" onClick={e => setAnalyze(true)}/>
      </div>
    </div>
    <div className="controls card mt-4">
      <div className="slider">
        <h5>Minimal number of revisions ({minRevs})</h5>
        <Slider value={minRevs} max={revsRange[1]} min={revsRange[0]} onChange={e => setMinRevs(e.value)}/>
      </div>
    </div>
    {screen}
  </>
}