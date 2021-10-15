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
            <TabView renderActiveOnly={true}>
                <TabPanel header="Revisions">
                    <Revisions sceneId={sceneId} appId={appId} date={dateRange}/>
                </TabPanel>
                <TabPanel header="Hotspots">
                    <Hotspots sceneId={sceneId} appId={appId}/>
                </TabPanel>
                <TabPanel header="Complexity">
                    <ComplexityAnalysis sceneId={sceneId} appId={appId}/>
                </TabPanel>
                <TabPanel header="Coupling">
                    <Coupling sceneId={sceneId} appId={appId}/>
                </TabPanel>
                <TabPanel header="Entity Coupling">
                    <EntityCoupling sceneId={sceneId} appId={appId}/>
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
                <TabPanel header="Knowledge Map">
                    <KnowledgeMap sceneId={sceneId} appId={appId}/>
                </TabPanel>
            </TabView>
        </>
    );
}