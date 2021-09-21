import {useChartDimensions} from "../hooks/useChartDimensions";
import * as d3 from "d3";
import Chart from "../chart/Chart";
import Axis from "../chart/Axis";
import {Line} from "../chart/Line";
import './MultineChart.css';

export function MultiLineChart({data, xAccessor, yAccessor, xFormatter, label, legend=[]}) {
  const [ref, dimensions] = useChartDimensions();

  const xScale = d3.scaleTime()
    .domain(d3.extent(data.flat(), xAccessor))
    .range([0, dimensions.boundedWidth])

  const yScale = d3.scaleLinear()
    .domain(d3.extent(data.flat(), yAccessor))
    .range([dimensions.boundedHeight, 0])
    .nice()

  const xAccessorScaled = d => xScale(xAccessor(d))
  const yAccessorScaled = d => yScale(yAccessor(d))
  const colorScale = d3.scaleSequential(d3.schemeTableau10);

  return (
    <div className="Timeline" ref={ref}>
      <div className="p-mr-6 p-d-flex p-flex-column">
        {
          legend.map((each, index) => {
            return <div style={{color: colorScale(index)}} className="p-mb-2">
              <span className="legend-box" style={{backgroundColor: colorScale(index)}}></span>
              <span>{each}</span>
            </div>
          })
        }
      </div>
      <Chart dimensions={dimensions}>
        <Axis
          dimension="x"
          scale={xScale}
          formatTick={xFormatter}
        />
        <Axis
          dimension="y"
          scale={yScale}
          label={label}
        />
        {
          data.map((each, index) => {
            return <Line
              stroke={colorScale(index)}
              data={each}
              xAccessor={xAccessorScaled}
              yAccessor={yAccessorScaled}
            />
          })}
      </Chart>
    </div>)
}
