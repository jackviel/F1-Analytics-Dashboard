import { Box, Grid, Card, CardContent, Typography, CircularProgress, Alert, Chip } from '@mui/material';
import { useRaces } from '../hooks/useF1Data';
import type { Race } from '../types/models';

const Races = () => {
  const { data: races = [], isLoading, error } = useRaces();

  if (isLoading) {
    return (
      <Box display="flex" justifyContent="center" alignItems="center" minHeight="60vh">
        <CircularProgress />
      </Box>
    );
  }

  if (error) {
    return <Alert severity="error">Failed to load races</Alert>;
  }

  return (
    <Box>
      <Typography variant="h4" gutterBottom>
        F1 Races
      </Typography>

      <Grid container spacing={3}>
        {races.map((race: Race) => (
          <Grid item xs={12} sm={6} md={4} key={race.id}>
            <Card>
              <CardContent>
                <Typography variant="h6">{race.name}</Typography>
                <Typography color="textSecondary" gutterBottom>
                  Round {race.round} - Season {race.season}
                </Typography>
                <Typography color="textSecondary">
                  Circuit: {race.circuit.name}
                </Typography>
                <Typography color="textSecondary">
                  Date: {new Date(race.date).toLocaleDateString()}
                </Typography>
                <Box mt={1}>
                  <Chip 
                    label={race.status}
                    color={
                      race.status === 'Completed' ? 'success' :
                      race.status === 'Cancelled' ? 'error' : 'primary'
                    }
                    size="small"
                  />
                </Box>
              </CardContent>
            </Card>
          </Grid>
        ))}
      </Grid>
    </Box>
  );
};

export default Races; 