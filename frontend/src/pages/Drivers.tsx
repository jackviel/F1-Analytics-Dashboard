import { useEffect } from 'react';
import { useDispatch, useSelector } from 'react-redux';
import {
  Grid,
  Card,
  CardContent,
  Typography,
  Box,
  CircularProgress,
  Alert,
} from '@mui/material';
import { AppDispatch, RootState } from '../store';
import { fetchDrivers } from '../store/slices/driversSlice';

const Drivers = () => {
  const dispatch = useDispatch<AppDispatch>();
  const { drivers, loading, error } = useSelector((state: RootState) => state.drivers);

  useEffect(() => {
    dispatch(fetchDrivers());
  }, [dispatch]);

  if (loading) {
    return (
      <Box display="flex" justifyContent="center" alignItems="center" minHeight="60vh">
        <CircularProgress />
      </Box>
    );
  }

  if (error) {
    return <Alert severity="error">{error}</Alert>;
  }

  return (
    <Box>
      <Typography variant="h4" gutterBottom>
        F1 Drivers
      </Typography>

      <Grid container spacing={3}>
        {drivers.map((driver) => (
          <Grid item xs={12} sm={6} md={4} key={driver.id}>
            <Card>
              <CardContent>
                <Typography variant="h6">{driver.name}</Typography>
                <Typography color="textSecondary">
                  Team: {driver.team?.name}
                </Typography>
                <Typography color="textSecondary">
                  Nationality: {driver.nationality}
                </Typography>
                <Typography color="textSecondary">
                  Number: {driver.number}
                </Typography>
              </CardContent>
            </Card>
          </Grid>
        ))}
      </Grid>
    </Box>
  );
};

export default Drivers; 