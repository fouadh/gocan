import React, {createContext, useContext} from "react"

import "./Chart.css"

const ChartContext = createContext()
export const useDimensionsContext = () => useContext(ChartContext)

const Chart = ({ dimensions, children }) => (
  <ChartContext.Provider value={dimensions}>
    <svg className="Chart" width={dimensions.width} height={dimensions.height}>
      <g transform={`translate(${dimensions.marginLeft}, ${dimensions.marginTop})`}>
        { children }
      </g>
    </svg>
  </ChartContext.Provider>
)

Chart.defaultProps = {
  dimensions: {}
}

export default Chart
