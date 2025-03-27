export interface Driver {
  id: number;
  name: string;
  nationality: string;
  dateOfBirth: string;
  number: number;
  teamId: number;
  team: Team;
  careerPoints: number;
  careerWins: number;
  careerPoles: number;
  careerFastLaps: number;
  careerPodiums: number;
  active: boolean;
  profileImageUrl: string;
  biography: string;
}

export interface Team {
  id: number;
  name: string;
  nationality: string;
  founded: string;
  baseLocation: string;
  teamPrincipal: string;
  technicalChief: string;
  chassis: string;
  engine: string;
  firstEntry: string;
  worldTitles: number;
  raceWins: number;
  poles: number;
  fastestLaps: number;
  podiums: number;
  points: number;
  active: boolean;
  logoUrl: string;
  website: string;
}

export interface Circuit {
  id: number;
  name: string;
  location: string;
  country: string;
  length: number;
  firstGrandPrix: string;
  lapRecord: string;
  lapRecordHolder: string;
  lapRecordYear: number;
  imageUrl: string;
  description: string;
}

export interface Race {
  id: number;
  name: string;
  season: number;
  round: number;
  circuitId: number;
  circuit: Circuit;
  date: string;
  raceTime: string;
  qualifyingTime: string;
  practice1Time: string;
  practice2Time: string;
  practice3Time: string;
  sprintTime?: string;
  status: 'Scheduled' | 'Completed' | 'Cancelled';
  laps: number;
  raceDistance: number;
  weather: string;
  temperature: number;
  trackCondition: string;
  drivers: Driver[];
  teams: Team[];
  results: RaceDriver[];
  teamResults: RaceTeam[];
}

export interface RaceDriver {
  driverId: number;
  raceId: number;
  position: number;
  points: number;
  grid: number;
  fastestLap: string;
  raceTime: string;
  status: string;
}

export interface RaceTeam {
  teamId: number;
  raceId: number;
  points: number;
  position: number;
}

export interface Lap {
  id: number;
  raceId: number;
  driverId: number;
  lapNumber: number;
  lapTime: string;
  position: number;
  isFastest: boolean;
  pitStop: boolean;
  pitStopTime?: string;
} 