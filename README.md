# URL Health Monitor

A web application that monitors the health and uptime of URLs,
and alerts users when a site goes down or becomes unreachable.

## Tech Stack

**Frontend:** React.js, Vite  
**Backend:** Node.js, Express.js  
**Database:** PostgreSQL, Prisma ORM  

## Features

- Add and manage URLs to monitor
- Real-time health status (UP / DOWN)
- Response time tracking
- Uptime percentage calculation
- Email/alert notifications on failure
- Dashboard with status overview

## Getting Started

### Prerequisites
- Node.js
- PostgreSQL

### Installation

# Clone the repo
git clone https://github.com/your-username/url-health-monitor.git

# Install frontend dependencies
cd frontend
npm install

# Install backend dependencies
cd ../backend
npm install

### Environment Variables

Create a `.env` file in the backend:

DATABASE_URL=postgresql://user:password@localhost:5432/urlmonitor
PORT=5000
CHECK_INTERVAL=60000  # in milliseconds

### Run the App

# Start backend
cd backend
npm run dev

# Start frontend
cd frontend
npm run dev

## How It Works

1. User adds a URL to monitor
2. Backend pings the URL at regular intervals
3. Status (UP/DOWN), response time, and timestamp are logged
4. Dashboard displays live health data
5. Alert is triggered if URL goes down

## Screenshots

(Add screenshots here)

## Author

Truppti — B.Tech CSE, Pimpri Chinchwad University"# Current-Url-Monitor" 
"A URL monitoring tool built with React + Vite" 
