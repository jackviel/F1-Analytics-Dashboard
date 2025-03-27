import { Box, Grid, Card, CardContent, Typography, CircularProgress, Alert } from '@mui/material';
import { useTeams } from '../hooks/useF1Data';

const Teams = () => {
  const { data: teams = [], isLoading, error } = useTeams();

  if (isLoading) {
    return (
      <Box display="flex" justifyContent="center" alignItems="center" minHeight="60vh">
        <CircularProgress />
      </Box>
    );
  }

  if (error) {
    return <Alert severity="error">Failed to load teams</Alert>;
  }

  return (
    <Box>
      <Typography variant="h4" gutterBottom>
        F1 Teams
      </Typography>

      <Grid container spacing={3}>
        {teams.map((team) => (
          <Grid item xs={12} sm={6} md={4} key={team.id}>
            <Card>
              <CardContent>
                <Typography variant="h6">{team.name}</Typography>
                <Typography color="textSecondary">
                  Base: {team.baseLocation}
                </Typography>
                <Typography color="textSecondary">
                  Principal: {team.teamPrincipal}
                </Typography>
                <Typography color="textSecondary">
                  Championships: {team.worldTitles}
                </Typography>
              </CardContent>
            </Card>
          </Grid>
        ))}
      </Grid>
    </Box>
  );
};

export default Teams; 