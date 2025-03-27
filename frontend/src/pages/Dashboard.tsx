import { Box, Grid, Card, CardContent, Typography, CircularProgress, Alert } from '@mui/material';
import {
  BarChart,
  Bar,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  Legend,
  ResponsiveContainer,
} from 'recharts';
import { useDrivers, useTeams, useRaces } from '../hooks/useF1Data';
import type { Team } from '../types/models';

interface TeamPointsData {
  name: string;
  points: number;
}

const Dashboard = () => {
  const { data: drivers = [], isLoading: driversLoading, error: driversError } = useDrivers();
  const { data: teams = [], isLoading: teamsLoading, error: teamsError } = useTeams();
  const { data: races = [], isLoading: racesLoading, error: racesError } = useRaces();

  if (driversLoading || teamsLoading || racesLoading) {
    return (
      <Box display="flex" justifyContent="center" alignItems="center" minHeight="60vh">
        <CircularProgress />
      </Box>
    );
  }

  if (driversError || teamsError || racesError) {
    return (
      <Alert severity="error">
        {driversError?.message || teamsError?.message || racesError?.message || 'An error occurred'}
      </Alert>
    );
  }

  // Calculate statistics
  const totalDrivers = drivers.length;
  const totalTeams = teams.length;
  const totalRaces = races.length;
  const activeDrivers = drivers.filter((driver) => driver.active).length;

  // Prepare data for the chart
  const teamPointsData: TeamPointsData[] = teams
    .map((team: Team) => ({
      name: team.name,
      points: team.points,
    }))
    .sort((a: TeamPointsData, b: TeamPointsData) => b.points - a.points)
    .slice(0, 5);

  return (
    <Box>
      <Typography variant="h4" gutterBottom>
        F1 Dashboard
      </Typography>

      <Grid container spacing={3}>
        {/* Statistics Cards */}
        <Grid item xs={12} sm={6} md={3}>
          <Card>
            <CardContent>
              <Typography color="textSecondary" gutterBottom>
                Total Drivers
              </Typography>
              <Typography variant="h4">{totalDrivers}</Typography>
            </CardContent>
          </Card>
        </Grid>
        <Grid item xs={12} sm={6} md={3}>
          <Card>
            <CardContent>
              <Typography color="textSecondary" gutterBottom>
                Active Drivers
              </Typography>
              <Typography variant="h4">{activeDrivers}</Typography>
            </CardContent>
          </Card>
        </Grid>
        <Grid item xs={12} sm={6} md={3}>
          <Card>
            <CardContent>
              <Typography color="textSecondary" gutterBottom>
                Total Teams
              </Typography>
              <Typography variant="h4">{totalTeams}</Typography>
            </CardContent>
          </Card>
        </Grid>
        <Grid item xs={12} sm={6} md={3}>
          <Card>
            <CardContent>
              <Typography color="textSecondary" gutterBottom>
                Total Races
              </Typography>
              <Typography variant="h4">{totalRaces}</Typography>
            </CardContent>
          </Card>
        </Grid>

        {/* Team Points Chart */}
        <Grid item xs={12}>
          <Card>
            <CardContent>
              <Typography variant="h6" gutterBottom>
                Top 5 Teams by Points
              </Typography>
              <Box height={400}>
                <ResponsiveContainer width="100%" height="100%">
                  <BarChart data={teamPointsData}>
                    <CartesianGrid strokeDasharray="3 3" />
                    <XAxis dataKey="name" />
                    <YAxis />
                    <Tooltip />
                    <Legend />
                    <Bar dataKey="points" fill="#8884d8" />
                  </BarChart>
                </ResponsiveContainer>
              </Box>
            </CardContent>
          </Card>
        </Grid>
      </Grid>
    </Box>
  );
};

export default Dashboard; 