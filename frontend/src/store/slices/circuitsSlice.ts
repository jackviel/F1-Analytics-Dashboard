import { createSlice, createAsyncThunk } from '@reduxjs/toolkit';
import type { Circuit } from '../../types/models';
import { getCircuits, getCircuit } from '../../services/api';

interface CircuitsState {
  circuits: Circuit[];
  selectedCircuit: Circuit | null;
  loading: boolean;
  error: string | null;
}

const initialState: CircuitsState = {
  circuits: [],
  selectedCircuit: null,
  loading: false,
  error: null,
};

export const fetchCircuits = createAsyncThunk(
  'circuits/fetchCircuits',
  async () => {
    const response = await getCircuits();
    return response.data;
  }
);

export const fetchCircuit = createAsyncThunk(
  'circuits/fetchCircuit',
  async (id: number) => {
    const response = await getCircuit(id);
    return response.data;
  }
);

const circuitsSlice = createSlice({
  name: 'circuits',
  initialState,
  reducers: {
    clearSelectedCircuit: (state) => {
      state.selectedCircuit = null;
    },
  },
  extraReducers: (builder) => {
    builder
      // Fetch all circuits
      .addCase(fetchCircuits.pending, (state) => {
        state.loading = true;
        state.error = null;
      })
      .addCase(fetchCircuits.fulfilled, (state, action) => {
        state.loading = false;
        state.circuits = action.payload;
      })
      .addCase(fetchCircuits.rejected, (state, action) => {
        state.loading = false;
        state.error = action.error.message || 'Failed to fetch circuits';
      })
      // Fetch single circuit
      .addCase(fetchCircuit.pending, (state) => {
        state.loading = true;
        state.error = null;
      })
      .addCase(fetchCircuit.fulfilled, (state, action) => {
        state.loading = false;
        state.selectedCircuit = action.payload;
      })
      .addCase(fetchCircuit.rejected, (state, action) => {
        state.loading = false;
        state.error = action.error.message || 'Failed to fetch circuit';
      });
  },
});

export const { clearSelectedCircuit } = circuitsSlice.actions;
export default circuitsSlice.reducer; 