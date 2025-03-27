import { useQuery } from '@tanstack/react-query';
import { getDrivers, getTeams, getRaces, getCircuits } from '../services/api';
import type { Driver, Team, Race, Circuit } from '../types/models';

// Query keys
export const queryKeys = {
  drivers: ['drivers'],
  teams: ['teams'],
  races: ['races'],
  circuits: ['circuits'],
} as const;

// Custom hooks for data fetching
export const useDrivers = () => {
  return useQuery<Driver[]>({
    queryKey: queryKeys.drivers,
    queryFn: async () => {
      const response = await getDrivers();
      return response.data;
    },
    staleTime: 5 * 60 * 1000, // Consider data fresh for 5 minutes
    gcTime: 30 * 60 * 1000, // Keep in cache for 30 minutes
  });
};

export const useTeams = () => {
  return useQuery<Team[]>({
    queryKey: queryKeys.teams,
    queryFn: async () => {
      const response = await getTeams();
      return response.data;
    },
    staleTime: 5 * 60 * 1000,
    gcTime: 30 * 60 * 1000,
  });
};

export const useRaces = () => {
  return useQuery<Race[]>({
    queryKey: queryKeys.races,
    queryFn: async () => {
      const response = await getRaces();
      return response.data;
    },
    staleTime: 5 * 60 * 1000,
    gcTime: 30 * 60 * 1000,
  });
};

export const useCircuits = () => {
  return useQuery<Circuit[]>({
    queryKey: queryKeys.circuits,
    queryFn: async () => {
      const response = await getCircuits();
      return response.data;
    },
    staleTime: 5 * 60 * 1000,
    gcTime: 30 * 60 * 1000,
  });
}; 