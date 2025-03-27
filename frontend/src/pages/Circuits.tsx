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
import { fetchCircuits } from '../store/slices/circuitsSlice';
import type { Circuit } from '../types/models';

const Circuits = () => {
  const dispatch = useDispatch<AppDispatch>();
  const { circuits, loading, error } = useSelector((state: RootState) => state.circuits);

  useEffect(() => {
    dispatch(fetchCircuits());
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
        F1 Circuits
      </Typography>

      <Grid container spacing={3}>
        {circuits.map((circuit: Circuit) => (
          <Grid item xs={12} sm={6} md={4} key={circuit.id}>
            <Card>
              <CardContent>
                <Typography variant="h6">{circuit.name}</Typography>
                <Typography color="textSecondary">
                  Location: {circuit.location}, {circuit.country}
                </Typography>
                <Typography color="textSecondary">
                  Length: {circuit.length}km
                </Typography>
                {circuit.lapRecord && (
                  <Typography color="textSecondary">
                    Lap Record: {circuit.lapRecord} ({circuit.lapRecordHolder}, {circuit.lapRecordYear})
                  </Typography>
                )}
              </CardContent>
            </Card>
          </Grid>
        ))}
      </Grid>
    </Box>
  );
};

export default Circuits; 