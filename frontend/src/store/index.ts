import { configureStore } from '@reduxjs/toolkit';
import driversReducer from './slices/driversSlice';
import teamsReducer from './slices/teamsSlice';
import racesReducer from './slices/racesSlice';
import circuitsReducer from './slices/circuitsSlice';

export const store = configureStore({
  reducer: {
    drivers: driversReducer,
    teams: teamsReducer,
    races: racesReducer,
    circuits: circuitsReducer,
  },
});

export type RootState = ReturnType<typeof store.getState>;
export type AppDispatch = typeof store.dispatch; 