import './App.css'
import {ThemeProvider, BaseStyles} from '@primer/react'

function App() {
  return (
    <ThemeProvider>
      <BaseStyles>
        <div>Hello world</div>
      </BaseStyles>
    </ThemeProvider>
  )
}

export default App
