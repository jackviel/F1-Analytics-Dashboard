import axios from 'axios';
import type { Driver, Team, Race, Circuit } from '../types/models';

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080/api/v1';

const api = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Driver endpoints
export const getDrivers = () => api.get<Driver[]>('/drivers');
export const getDriver = (id: number) => api.get<Driver>(`/drivers/${id}`);
export const getDriverStats = (id: number) => api.get(`/drivers/${id}/stats`);

// Team endpoints
export const getTeams = () => api.get<Team[]>('/teams');
export const getTeam = (id: number) => api.get<Team>(`/teams/${id}`);
export const getTeamStats = (id: number) => api.get(`/teams/${id}/stats`);

// Race endpoints
export const getRaces = () => api.get<Race[]>('/races');
export const getRace = (id: number) => api.get<Race>(`/races/${id}`);
export const getRaceResults = (id: number) => api.get(`/races/${id}/results`);

// Circuit endpoints
export const getCircuits = () => api.get<Circuit[]>('/circuits');
export const getCircuit = (id: number) => api.get<Circuit>(`/circuits/${id}`);
export const getCircuitHistory = (id: number) => api.get(`/circuits/${id}/history`);

// Add request interceptor for authentication
api.interceptors.request.use((config) => {
  const token = localStorage.getItem('token');
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

// Add response interceptor for error handling
api.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401) {
      // Handle unauthorized access
      localStorage.removeItem('token');
      window.location.href = '/login';
    }
    return Promise.reject(error);
  }
);

export default api; 