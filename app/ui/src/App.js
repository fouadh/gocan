import 'primereact/resources/themes/saga-blue/theme.css'
import 'primereact/resources/primereact.min.css'
import 'primeicons/primeicons.css'
import 'primeflex/primeflex.css'
import { Menu } from './Menu'
import { Scenes } from './screens/Scenes'
import { BrowserRouter, Route, Routes } from 'react-router-dom'
import { Scene } from './screens/Scene'
import { AppDetails } from './screens/AppDetails'

function App() {
  return (
    <BrowserRouter basename="/">
      <div className="App layout-wrapper">
        <div className="layout-topbar">
          <Menu />
          <Routes>
            <Route
              path="/scenes/:sceneId/apps/:appId"
              element={<AppDetails />}
            />
            <Route path="/scenes/:sceneId" element={<Scene />} />
            <Route path="/scenes" element={<Scenes />} />
            <Route path="/" element={<Scenes />} />
          </Routes>
        </div>
      </div>
    </BrowserRouter>
  )
}

export default App
