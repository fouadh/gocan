import {useEffect, useState} from "react";
import axios from "axios";
import {Dropdown} from "primereact/dropdown";

export function ModuleSelector({appId, sceneId, boundaryName, onChange}) {
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