import {Link, useParams} from 'react-router-dom'
import {useEffect, useState} from 'react'
import axios from "axios";
import {Panel} from 'primereact/panel';
import {DataView} from 'primereact/dataview';
import {TabPanel, TabView} from 'primereact/tabview';
import './Scene.css';

export function Scene() {
    const {sceneId} = useParams();
    const [scene, setScene] = useState();
    const [applications, setApplications] = useState([]);

    useEffect(() => {
        let subscribed = true;

        axios.get(`/api/scenes/${sceneId}`)
            .then(it => it.data)
            .then(it => {
                if (subscribed)
                    setScene(it);
            });

        axios.get(`/api/scenes/${sceneId}/apps`)
            .then(it => it.data)
            .then(it => it.apps)
            .then(it => {
                if (subscribed) {
                    setApplications(it);
                }
            })

        return () => subscribed = false;
    }, [sceneId]);

    const AppSummary = (data) => {
        const headerTemplate = <div className="p-text-center"><Link
            to={`/scenes/${sceneId}/apps/${data.id}?after=${data.dateRange.minDate}&before=${data.dateRange.maxDate}`}>
            <span data-name={data.name}>{data.name}</span>
        </Link>
        </div>;
        return (
            <>
                <Panel header={headerTemplate} className="p-ml-4">
                    <div className="p-d-flex p-text-center">
                        (Between {data.dateRange.minDate} and {data.dateRange.maxDate})
                    </div>
                    <div className="p-d-flex p-text-center">
                        <div className="p-mr-4">
                            <div className="p-d-flex p-flex-column">
                                <div className="p-mb-2 data-title">Commits</div>
                                <div className="p-mb-2 data-value">{data.numberOfCommits}</div>
                            </div>
                        </div>
                        <div className="p-mr-4">
                            <div className="p-d-flex p-flex-column">
                                <div className="p-mb-2 data-title">Entities</div>
                                <div className="p-mb-2 data-value">{data.numberOfEntities}</div>
                            </div>
                        </div>
                        <div className="p-mr-4">
                            <div className="p-d-flex p-flex-column">
                                <div className="p-mb-2 data-title">Changed Entities</div>
                                <div className="p-mb-2 data-value">{data.numberOfEntitiesChanged}</div>
                            </div>
                        </div>
                        <div className="p-d-flex p-flex-column">
                            <div className="p-mb-2 data-title">Authors</div>
                            <div className="p-mb-2 data-value">{data.numberOfAuthors}</div>
                        </div>
                    </div>
                </Panel>
            </>
        );
    }

    const sceneDateRange = {};

    return scene ? <div>
        <h1>{scene.name}</h1>
        <div>
            <TabView renderActiveOnly={true}>
                <TabPanel header="Applications">
                    <div className="apps">
                        <DataView value={applications} layout="grid" itemTemplate={AppSummary}/>
                    </div>
                </TabPanel>
            </TabView>
        </div>
    </div> : <>
        <div>Scene Not Found</div>
    </>;
}