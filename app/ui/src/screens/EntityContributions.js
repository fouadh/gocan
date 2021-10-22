import {useChartDimensions} from "../hooks/useChartDimensions";
import Chart from "../chart/Chart";
import * as d3 from "d3"

export function EntityContributions() {
    const [ref, dimensions] = useChartDimensions();

    const contributions = [
        {dev: "dev1", contributions: 5},
        {dev: "dev2", contributions: 2},
        {dev: "dev3", contributions: 2},
        {dev: "dev4", contributions: 1},
        {dev: "dev5", contributions: 1},
        {dev: "dev6", contributions: 1},
        {dev: "dev7", contributions: 1},
        {dev: "dev8", contributions: 1}
    ];

    const buildRectangles = (contributions) => {
        const colorScale = d3.scaleOrdinal(d3.schemeSet1);
        const totalContributions = contributions.map(c => c.contributions).reduce((a, c) => a + c, 0);

        const initialRectangle = {x: 10, y: 10, height: 600, width: 400, fill: "#69b3a2"};
        const totalArea = initialRectangle.height * initialRectangle.width;
        const rectangle0 = {
            x: initialRectangle.x,
            y: initialRectangle.y,
            height: initialRectangle.height,
            width: (contributions[0].contributions / totalContributions) * totalArea / initialRectangle.height,
            fill: colorScale(0)
        };

        const rectangles = [initialRectangle, rectangle0];

        for (let i = 2; i <= contributions.length; i++) {
            if (i % 2 === 0) {
                rectangles.push({
                    x: rectangles[i-1].x + rectangles[i-1].width,
                    y: rectangles[i-1].y,
                    height: (contributions[i-1].contributions / totalContributions) * totalArea / (rectangles[i-2].width - rectangles[i-1].width),
                    width: rectangles[i-2].width - rectangles[i-1].width,
                    fill: colorScale(i)
                });
            } else {
                rectangles.push({
                    x: rectangles[i-1].x,
                    y: rectangles[i-1].y + rectangles[i-1].height,
                    height: rectangles[i-2].height - rectangles[i-1].height,
                    width: (contributions[i-1].contributions / totalContributions) * totalArea / (rectangles[i-2].height - rectangles[i-1].height),
                    fill: colorScale(i)
                });
            }
        }

        return rectangles;
    }
    const rectangles = buildRectangles(contributions);
    return <>
        <div className="js-viz" ref={ref}>
            <Chart dimensions={dimensions}>
                {
                    rectangles.map(r => <rect stroke={"#eee"} x={r.x} y={r.y} height={r.height} width={r.width}
                                              fill={r.fill}/>)
                }
            </Chart>


        </div>
    </>
}