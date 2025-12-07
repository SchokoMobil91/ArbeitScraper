import { useEffect, useState } from 'react';
import type { Job } from './models/job';
import { getJobs } from './api';
import './App.css';

function App() {
  const [jobs, setJobs] = useState<Job[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const fetchData = async () => {
      try {
        const jobsData = await getJobs();
        setJobs(jobsData);
      } catch (error) {
        console.error('Error fetching jobs:', error);
      } finally {
        setLoading(false);
      }
    };

    fetchData();
  }, []);

  if (loading)
    return (
      <div className='loading-container'>
        <div className='spin-icon'></div>
        <p>Loading data...</p>
      </div>
    );

  if (jobs.length === 0) {
    return (
      <div className='app-container'>
        <h1>Arbeitsagentur Jobs</h1>
        <p>No job data found.</p>
        <p className='error-message'>
          Please ensure the backend is running and the scraping process completed successfully.
        </p>
      </div>
    );
  }

  const getDescriptionText = (job: Job) => {
    return job.fullDescription.length > 0
      ? job.fullDescription.substring(0, 150) + '...'
      : job.shortDescription;
  };

  return (
    <div className='app-container'>
      <h1>Arbeitsagentur Jobs</h1>

      <div className='table-wrapper'>
        <table className='job-table'>
          <thead>
            <tr>
              <th>Profession</th>
              <th>Company</th>
              <th>Location</th>
              <th>Salary</th>
              <th>Start Date</th>
              <th>Telephone</th>
              <th>Email</th>
              <th>Description</th>
              <th>Links</th>
            </tr>
          </thead>

          <tbody>
            {jobs.map((job, i) => (
              <tr key={i}>
                <td>{job.profession}</td>
                <td>{job.company}</td>
                <td>{job.location}</td>
                <td>{job.salary}</td>
                <td>{job.startDate}</td>
                <td>{job.telephone}</td>
                <td>{job.email}</td>
                <td>{getDescriptionText(job)}</td>
                <td className='links'>
                  {job.applicationLink && (
                    <a href={job.applicationLink} target='_blank'>
                      Apply
                    </a>
                  )}
                  {job.externalLink && (
                    <a href={job.externalLink} target='_blank'>
                      Detail
                    </a>
                  )}
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  );
}

export default App;
