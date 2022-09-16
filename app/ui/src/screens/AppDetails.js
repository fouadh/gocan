import axios from "axios";
import {useParams, useLocation} from 'react-router-dom'
import {useEffect, useState} from "react";
import {TabPanel, TabView} from 'primereact/tabview';
import {Revisions} from "./Revisions";
import {Coupling} from "./Coupling";
import {CodeChurn} from "./CodeChurn";
import {ModusOperandi} from "./ModusOperandi";
import {RevisionTrends} from "./RevisionTrends";
import {KnowledgeMap} from "./KnowledgeMap";
import {Hotspots} from "./Hotspots";
import {ComplexityAnalysis} from "./ComplexityAnalysis";
import {EntityCoupling} from "./EntityCoupling";
import {EntityContributions} from "./EntityContributions";
import {CoordinationNeeds} from "./CoordinationNeeds";
import {CodeAge} from "./CodeAge";

function useQuery() {
    return new URLSearchParams(useLocation().search);
}

export function AppDetails() {
    const query = useQuery();
    const dateRange = {
        min: query.get("after"),
        max: query.get("before")
    }
    const {sceneId, appId} = useParams();
    const [application, setApplication] = useState();
    const [entities, setEntities] = useState([]);

    useEffect(() => {
        let subscribed = true;
        axios.get(`/api/scenes/${sceneId}/apps/${appId}`)
            .then(it => it.data)
            .then(it => {
                if (subscribed)
                    setApplication(it);
            });

        axios.get(`/api/scenes/${sceneId}/apps/${appId}/entities`)
            .then(it => it.data)
            .then(it => it.entities)
            .then(it => {
               if (subscribed)
                   setEntities(it);
            });

        return () => subscribed = false;
    }, [sceneId, appId]);

    return (
        <>
            <div><h3 className="p-ml-4 ml-4">{application?.name}</h3></div>
            <TabView renderActiveOnly={true} id="tabs">
                <TabPanel header="Revisions">
                    <Revisions sceneId={sceneId} appId={appId} date={dateRange}/>
                </TabPanel>
                <TabPanel header="Hotspots">
                    <Hotspots sceneId={sceneId} appId={appId} date={dateRange}/>
                </TabPanel>
                <TabPanel header="Code Age">
                    <CodeAge sceneId={sceneId} appId={appId} date={dateRange}/>
                </TabPanel>
                <TabPanel header="Complexity">
                    <ComplexityAnalysis sceneId={sceneId} appId={appId}/>
                </TabPanel>
                <TabPanel header="Coupling">
                    <Coupling sceneId={sceneId} appId={appId} date={dateRange}/>
                </TabPanel>
                <TabPanel header="Entity Coupling">
                    <EntityCoupling sceneId={sceneId} appId={appId} date={dateRange} entities={entities} />
                </TabPanel>
                <TabPanel header="Code Churn">
                    <CodeChurn sceneId={sceneId} appId={appId} date={dateRange}/>
                </TabPanel>
                <TabPanel header="Modus Operandi">
                    <ModusOperandi sceneId={sceneId} appId={appId}/>
                </TabPanel>
                <TabPanel header="Revisions Trends">
                    <RevisionTrends sceneId={sceneId} appId={appId}/>
                </TabPanel>
                <TabPanel header="Knowledge Map">
                    <KnowledgeMap sceneId={sceneId} appId={appId} date={dateRange}/>
                </TabPanel>
                <TabPanel header="Coordination Needs">
                    <CoordinationNeeds sceneId={sceneId} appId={appId} date={dateRange}/>
                </TabPanel>
                <TabPanel header="Entity Contributions">
                    <EntityContributions sceneId={sceneId} appId={appId} date={dateRange} entities={entities} />
                </TabPanel>
                {/*<TabPanel header="Developers Network">*/}
                {/*    <DevNetwork sceneId={sceneId} appId={appId} date={dateRange}/>*/}
                {/*</TabPanel>*/}
            </TabView>
        </>
    );
}