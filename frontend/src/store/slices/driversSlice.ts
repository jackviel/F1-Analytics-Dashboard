import { createSlice, createAsyncThunk } from '@reduxjs/toolkit';
import type { Driver } from '../../types/models';
import { getDrivers, getDriver } from '../../services/api';

interface DriversState {
  drivers: Driver[];
  selectedDriver: Driver | null;
  loading: boolean;
  error: string | null;
}

const initialState: DriversState = {
  drivers: [],
  selectedDriver: null,
  loading: false,
  error: null,
};

export const fetchDrivers = createAsyncThunk(
  'drivers/fetchDrivers',
  async () => {
    const response = await getDrivers();
    return response.data;
  }
);

export const fetchDriver = createAsyncThunk(
  'drivers/fetchDriver',
  async (id: number) => {
    const response = await getDriver(id);
    return response.data;
  }
);

const driversSlice = createSlice({
  name: 'drivers',
  initialState,
  reducers: {
    clearSelectedDriver: (state) => {
      state.selectedDriver = null;
    },
  },
  extraReducers: (builder) => {
    builder
      // Fetch all drivers
      .addCase(fetchDrivers.pending, (state) => {
        state.loading = true;
        state.error = null;
      })
      .addCase(fetchDrivers.fulfilled, (state, action) => {
        state.loading = false;
        state.drivers = action.payload;
      })
      .addCase(fetchDrivers.rejected, (state, action) => {
        state.loading = false;
        state.error = action.error.message || 'Failed to fetch drivers';
      })
      // Fetch single driver
      .addCase(fetchDriver.pending, (state) => {
        state.loading = true;
        state.error = null;
      })
      .addCase(fetchDriver.fulfilled, (state, action) => {
        state.loading = false;
        state.selectedDriver = action.payload;
      })
      .addCase(fetchDriver.rejected, (state, action) => {
        state.loading = false;
        state.error = action.error.message || 'Failed to fetch driver';
      });
  },
});

export const { clearSelectedDriver } = driversSlice.actions;
export default driversSlice.reducer; 