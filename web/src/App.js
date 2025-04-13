import React, { useState, useEffect } from 'react';
import './App.css';

function App() {
  const [assignments, setAssignments] = useState([]);
  const [totalHours, setTotalHours] = useState(0);
  const [totalWeeks, setTotalWeeks] = useState(0);
  const [error, setError] = useState(null);
  const [isLoading, setIsLoading] = useState(true);

  const formatHours = (decimalHours) => {
    const hours = Math.floor(decimalHours);
    const minutes = Math.round((decimalHours - hours) * 60);
    if (minutes === 0) {
      return `${hours} hours`;
    }
    if (hours === 0) {
      return `${minutes} minutes`;
    }
    return `${hours} hours ${minutes} minutes`;
  };

  useEffect(() => {
    const controller = new AbortController();
    const signal = controller.signal;

    const fetchAssignments = async () => {
      try {
        setIsLoading(true);
        const response = await fetch(`${process.env.REACT_APP_API_URL}/api/weekly-plan`, {
          signal,
          headers: {
            'Accept': 'application/json',
          }
        });
        const data = await response.json();
        setAssignments(data.assignments);
        setTotalHours(data.totalHours);
        setTotalWeeks(data.totalWeeks);
      } catch (err) {
        if (err.name === 'AbortError') {
          console.log('Fetch aborted');
        } else {
          setError('Error loading assignments. Please try again later.');
          console.error('Error fetching assignments:', err);
        }
      } finally {
        setIsLoading(false);
      }
    };

    fetchAssignments();

    return () => {
      controller.abort();
    };
  }, []);

  if (error) {
    return <div className="error">{error}</div>;
  }

  if (isLoading) {
    return (
      <div className="loading-container">
        <div className="spinner"></div>
        <p>Loading assignments...</p>
      </div>
    );
  }

  return (
    <div className="container">
      <h1>Weekly Assignments</h1>
      <div className="assignments-container">
        {assignments.map((developerAssignments, index) => {
          // Get the first assignment to get the developer info
          const firstAssignment = developerAssignments[0];
          const developerName = firstAssignment.developer.name;
          const productivity = firstAssignment.developer.productivity;
          
          // Group assignments by week
          const assignmentsByWeek = developerAssignments.reduce((acc, assignment) => {
            if (!acc[assignment.week_number]) {
              acc[assignment.week_number] = [];
            }
            acc[assignment.week_number].push(assignment);
            return acc;
          }, {});

          return (
            <div key={index} className="developer-container">
              <div className="developer-header">
                <span className="developer-name">{developerName}</span>
                <span className="productivity">Productivity: {productivity.toFixed(2)}</span>
              </div>
              {Object.entries(assignmentsByWeek).map(([weekNumber, weekAssignments]) => (
                <div key={weekNumber} className="week-container">
                  <div className="week-header">Week {weekNumber}</div>
                  {weekAssignments.map((assignment, assignmentIndex) => (
                    <div key={assignmentIndex} className="assignment">
                      <div className="task-info">
                        {assignment.task_name}
                        <span className="task-details">
                          (Difficulty: {assignment.task.difficulty.toFixed(2)}, Duration: {assignment.task.estimated_duration.toFixed(2)}h)
                        </span>
                      </div>
                      <div className="hours-info">
                        {formatHours(assignment.calculated_hours)}
                      </div>
                    </div>
                  ))}
                </div>
              ))}
            </div>
          );
        })}
      </div>
      <div className="summary">
        <p>Total Time: {formatHours(totalHours)}</p>
        <p>Total Weeks: {totalWeeks}</p>
      </div>
    </div>
  );
}

export default App; 