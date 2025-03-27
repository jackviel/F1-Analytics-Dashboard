import { createSlice, createAsyncThunk } from '@reduxjs/toolkit';
import type { Race } from '../../types/models';
import { getRaces, getRace, getRaceResults } from '../../services/api';

interface RacesState {
  races: Race[];
  selectedRace: Race | null;
  loading: boolean;
  error: string | null;
}

const initialState: RacesState = {
  races: [],
  selectedRace: null,
  loading: false,
  error: null,
};

export const fetchRaces = createAsyncThunk(
  'races/fetchRaces',
  async () => {
    const response = await getRaces();
    return response.data;
  }
);

export const fetchRace = createAsyncThunk(
  'races/fetchRace',
  async (id: number) => {
    const response = await getRace(id);
    return response.data;
  }
);

export const fetchRaceResults = createAsyncThunk(
  'races/fetchRaceResults',
  async (id: number) => {
    const response = await getRaceResults(id);
    return response.data;
  }
);

const racesSlice = createSlice({
  name: 'races',
  initialState,
  reducers: {
    clearSelectedRace: (state) => {
      state.selectedRace = null;
    },
  },
  extraReducers: (builder) => {
    builder
      // Fetch all races
      .addCase(fetchRaces.pending, (state) => {
        state.loading = true;
        state.error = null;
      })
      .addCase(fetchRaces.fulfilled, (state, action) => {
        state.loading = false;
        state.races = action.payload;
      })
      .addCase(fetchRaces.rejected, (state, action) => {
        state.loading = false;
        state.error = action.error.message || 'Failed to fetch races';
      })
      // Fetch single race
      .addCase(fetchRace.pending, (state) => {
        state.loading = true;
        state.error = null;
      })
      .addCase(fetchRace.fulfilled, (state, action) => {
        state.loading = false;
        state.selectedRace = action.payload;
      })
      .addCase(fetchRace.rejected, (state, action) => {
        state.loading = false;
        state.error = action.error.message || 'Failed to fetch race';
      })
      // Fetch race results
      .addCase(fetchRaceResults.pending, (state) => {
        state.loading = true;
        state.error = null;
      })
      .addCase(fetchRaceResults.fulfilled, (state, action) => {
        state.loading = false;
        if (state.selectedRace) {
          state.selectedRace.results = action.payload;
        }
      })
      .addCase(fetchRaceResults.rejected, (state, action) => {
        state.loading = false;
        state.error = action.error.message || 'Failed to fetch race results';
      });
  },
});

export const { clearSelectedRace } = racesSlice.actions;
export default racesSlice.reducer; 