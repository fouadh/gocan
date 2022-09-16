import {useEffect, useState} from "react";
import axios from "axios";
import {CirclePacking} from "../components/CirclePacking";
import {Spinner} from "../components/Spinner";
import {DateSelector} from "../components/DateSelector";
import {Button} from 'primereact/button';
import {BoundarySelector} from "./BoundarySelector";
import {Dropdown} from "primereact/dropdown";

function ModuleSelector({appId, sceneId, boundaryName, onChange}) {
    const [module, setModule] = useState()
    const [modules, setModules] = useState([]);

    useEffect(() => {
        let subscribed = true;
        axios.get(`/api/scenes/${sceneId}/apps/${appId}/boundaries/${boundaryName}/modules`)
            .then(it => it.data)
            .then(it => it.modules)
            .then((it) => {
                if (subscribed) {
                    setModules(it);
                }
            });
        return () => subscribed = false;
    }, [sceneId, appId, boundaryName]);


    let selector;
    if (modules && modules.length > 0) {
        selector = <>
            <>
                <div className="p-field p-col-12 p-md-4 mr-4">
                    <span className="p-float-label">
                        <Dropdown id="modules"
                                  optionLabel="name"
                                  optionValue="name"
                                  options={modules}
                                  value={module}
                                  showClear={true}
                                  onChange={(e) => {
                                      setModule(e.value);
                                      onChange(e);
                                  }}/>
                        <label htmlFor="modules">Module</label>
                    </span>
                </div>
            </>
        </>
    } else {
        selector = <>
        </>
    }

    return <>{selector}</>;
}

export function Hotspots({sceneId, appId, date}) {
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
                endpoint = `/api/scenes/${sceneId}/apps/${appId}/hotspots`;
            } else {
                endpoint = `/api/scenes/${sceneId}/hotspots`;
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
    }, [sceneId, appId, analyze, dateRange, boundaryName, moduleName]);

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
                { appId && <BoundarySelector sceneId={sceneId} appId={appId} onChange={e => setBoundaryName(e.value)} /> }
                { boundaryName && <ModuleSelector sceneId={sceneId} appId={appId} boundaryName={boundaryName} onChange={e => setModuleName(e.value)} /> }
                <Button label="Submit" onClick={e => setAnalyze(true)}/>
            </div>
        </div>
        {screen}
    </>
}