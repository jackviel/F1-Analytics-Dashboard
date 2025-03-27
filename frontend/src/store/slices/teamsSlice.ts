import { createSlice, createAsyncThunk } from '@reduxjs/toolkit';
import type { Team } from '../../types/models';
import { getTeams, getTeam } from '../../services/api';

interface TeamsState {
  teams: Team[];
  selectedTeam: Team | null;
  loading: boolean;
  error: string | null;
}

const initialState: TeamsState = {
  teams: [],
  selectedTeam: null,
  loading: false,
  error: null,
};

export const fetchTeams = createAsyncThunk(
  'teams/fetchTeams',
  async () => {
    const response = await getTeams();
    return response.data;
  }
);

export const fetchTeam = createAsyncThunk(
  'teams/fetchTeam',
  async (id: number) => {
    const response = await getTeam(id);
    return response.data;
  }
);

const teamsSlice = createSlice({
  name: 'teams',
  initialState,
  reducers: {
    clearSelectedTeam: (state) => {
      state.selectedTeam = null;
    },
  },
  extraReducers: (builder) => {
    builder
      // Fetch all teams
      .addCase(fetchTeams.pending, (state) => {
        state.loading = true;
        state.error = null;
      })
      .addCase(fetchTeams.fulfilled, (state, action) => {
        state.loading = false;
        state.teams = action.payload;
      })
      .addCase(fetchTeams.rejected, (state, action) => {
        state.loading = false;
        state.error = action.error.message || 'Failed to fetch teams';
      })
      // Fetch single team
      .addCase(fetchTeam.pending, (state) => {
        state.loading = true;
        state.error = null;
      })
      .addCase(fetchTeam.fulfilled, (state, action) => {
        state.loading = false;
        state.selectedTeam = action.payload;
      })
      .addCase(fetchTeam.rejected, (state, action) => {
        state.loading = false;
        state.error = action.error.message || 'Failed to fetch team';
      });
  },
});

export const { clearSelectedTeam } = teamsSlice.actions;
export default teamsSlice.reducer; 