import type { Job } from './models/job';

export async function getJobs(): Promise<Job[]> {
  const res = await fetch('http://localhost:8080/jobs');
  return res.json();
}
