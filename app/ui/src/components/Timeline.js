import {useChartDimensions} from "../hooks/useChartDimensions";
import * as d3 from "d3";
import {Line} from "../chart/Line";
import Axis from "../chart/Axis";
import {useUniqueId} from "../hooks/useUniqueId";
import Gradient from "../chart/Gradient";
import Chart from "../chart/Chart";
import './Timeline.css';

const gradientColors = ["rgb(226, 222, 243)", "#f8f9fa"]

const Timeline = ({data, xAccessor, yAccessor, label, xFormatter}) => {
  const [ref, dimensions] = useChartDimensions();
  const gradientId = useUniqueId("Timeline-gradient");

  const xScale = d3.scaleTime()
    .domain(d3.extent(data, xAccessor))
    .range([0, dimensions.boundedWidth])

  const yScale = d3.scaleLinear()
    .domain(d3.extent(data, yAccessor))
    .range([dimensions.boundedHeight, 0])
    .nice()

  const xAccessorScaled = d => xScale(xAccessor(d))
  const yAccessorScaled = d => yScale(yAccessor(d))
  const y0AccessorScaled = yScale(yScale.domain()[0])

  return (
    <div className="Timeline js-viz" ref={ref}>
      <Chart dimensions={dimensions}>
        <defs>
          <Gradient
            id={gradientId}
            colors={gradientColors}
            x2="0"
            y2="100%"
          />
        </defs>
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
        <Line
          type="area"
          data={data}
          xAccessor={xAccessorScaled}
          yAccessor={yAccessorScaled}
          y0Accessor={y0AccessorScaled}
          style={{fill: `url(#${gradientId})`}}
        />
        <Line
          data={data}
          xAccessor={xAccessorScaled}
          yAccessor={yAccessorScaled}
        />
      </Chart>
    </div>
  )
}

Timeline.defaultProps = {
  xAccessor: d => d.x,
  yAccessor: d => d.y,
  xFormatter: d3.timeFormat("%-b %-d")
}

export {Timeline};