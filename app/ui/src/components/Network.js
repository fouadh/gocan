import {useEffect, useState} from "react";
import {useChartDimensions} from "../hooks/useChartDimensions";
import Chart from "../chart/Chart";
import * as d3 from "d3";

export function Network({data}) {
    const [nodes, setNodes] = useState([])
    const [ref, dimensions] = useChartDimensions({height: "600px", width: "800px"});

    useEffect(() => {
        const simulation = d3.forceSimulation(data.nodes)
            .force("x", d3.forceX(400))
            .force("y", d3.forceY(300))
           .force('charge', d3.forceManyBody().strength(-60))
        ;

        simulation.on('tick', () => {
            setNodes([...simulation.nodes()])
        });

        simulation.nodes([...data.nodes]);
        simulation.alpha(0.1).restart();

        return () => simulation.stop();
    }, [data])

    return <div className="network" ref={ref}>
        <Chart dimensions={dimensions}>
            { nodes.map((node) => (
                <circle
                    cx={node.x}
                    cy={node.y}
                    r="5"
                    key={node.id}
                    fill="#69b3a2"
                />
            )) }
        </Chart>
    </div>;
}
