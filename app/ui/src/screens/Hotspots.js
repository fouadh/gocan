import {useEffect, useState} from "react";
import axios from "axios";
import {CirclePacking} from "../components/CirclePacking";
import {Spinner} from "../components/Spinner";
import {DateSelector} from "../components/DateSelector";
import {Button} from 'primereact/button';
import {BoundarySelector} from "./BoundarySelector";
import {ModuleSelector} from "./ModuleSelector";
import * as _ from 'lodash'

function revisionsToHotspots(appName, revisions) {
  const hierarchies = {}
  const filteredRevisions = revisions.filter(rev => rev.code > 0)
  const maxNumberOfRevs = _.maxBy(filteredRevisions, rev => rev.numberOfRevisions).numberOfRevisions

  filteredRevisions.filter(rev => rev.code > 0).forEach(rev => {
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
          .then(revisions => revisionsToHotspots(appName, revisions.revisions))
          .then(it => {
            if (subscribed) {
              setHotspots(it);
            }
          }).finally(() => {
        setLoading(false);
        setAnalyze(false);
      });
    }

    return () => subscribed = false;
  }, [sceneId, appId, appName, analyze, dateRange, boundaryName, moduleName]);

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
    {screen}
  </>
}