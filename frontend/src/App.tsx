import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import { Provider } from 'react-redux';
import { ThemeProvider, CssBaseline } from '@mui/material';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { store } from './store';
import theme from './theme';

// Layout components
import Layout from './components/Layout';
import Navbar from './components/Navbar';

// Page components
import Dashboard from './pages/Dashboard';
import Drivers from './pages/Drivers';
import Teams from './pages/Teams';
import Races from './pages/Races';
import Circuits from './pages/Circuits';
import DriverComparison from './pages/DriverComparison';
import RaceStrategy from './pages/RaceStrategy';
import SeasonPoints from './pages/SeasonPoints';
import TeamPerformance from './pages/TeamPerformance';
import CircuitHistory from './pages/CircuitHistory';

const queryClient = new QueryClient();

function App() {
  return (
    <Provider store={store}>
      <QueryClientProvider client={queryClient}>
        <ThemeProvider theme={theme}>
          <CssBaseline />
          <Router>
            <Layout>
              <Navbar />
              <Routes>
                <Route path="/" element={<Dashboard />} />
                <Route path="/drivers" element={<Drivers />} />
                <Route path="/drivers/:id" element={<DriverComparison />} />
                <Route path="/teams" element={<Teams />} />
                <Route path="/teams/:id" element={<TeamPerformance />} />
                <Route path="/races" element={<Races />} />
                <Route path="/races/:id" element={<RaceStrategy />} />
                <Route path="/circuits" element={<Circuits />} />
                <Route path="/circuits/:id" element={<CircuitHistory />} />
                <Route path="/season-points" element={<SeasonPoints />} />
              </Routes>
            </Layout>
          </Router>
        </ThemeProvider>
      </QueryClientProvider>
    </Provider>
  );
}

export default App;
