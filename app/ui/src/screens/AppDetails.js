import axios from "axios";
import {useParams} from 'react-router-dom'
import {useEffect, useState} from "react";
import {TabPanel, TabView} from 'primereact/tabview';
/*import {Revisions} from "./Revisions";
import {CodeChurn} from "./CodeChurn";
import {ModusOperandi} from "./ModusOperandi";
import {RevisionTrends} from "./RevisionTrends";
import {Coupling} from "./Coupling";
import "../components/Timeline.css";
import {ActiveSet} from "./ActiveSet";*/

export function AppDetails() {
  const {sceneId, appId} = useParams();
  const [application, setApplication] = useState();

  useEffect(() => {
    let subscribed = true;
    axios.get(`/api/scenes/${sceneId}/apps/${appId}`)
      .then(it => it.data)
      .then(it => {
        if (subscribed)
          setApplication(it);
      });

    return () => subscribed = false;
  }, [sceneId, appId]);

  return (
    <>
      <div><h3 className="p-ml-4">{application?.name}</h3></div>
      {/*<TabView renderActiveOnly={true}>
        <TabPanel header="Revisions">
          <Revisions sceneId={sceneId} appId={appId}/>
        </TabPanel>
        <TabPanel header="Coupling">
          <Coupling sceneId={sceneId} appId={appId}/>
        </TabPanel>
        <TabPanel header="Code Churn">
          <CodeChurn sceneId={sceneId} appId={appId}/>
        </TabPanel>
        <TabPanel header="Modus Operandi">
          <ModusOperandi sceneId={sceneId} appId={appId}/>
        </TabPanel>
        <TabPanel header="Revisions Trends">
          <RevisionTrends sceneId={sceneId} appId={appId}/>
        </TabPanel>
        <TabPanel header="Entities Active Set">
          <ActiveSet sceneId={sceneId} appId={appId}/>
        </TabPanel>
      </TabView>*/}
    </>
  );
}