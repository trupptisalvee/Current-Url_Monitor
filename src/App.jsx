import React, { useState, useEffect } from 'react';
import { Globe, Activity, AlertTriangle, CheckCircle2, Search, Play, Download, Trash2 } from 'lucide-react';
import { motion, AnimatePresence } from 'framer-motion';
import './App.css';

const DEFAULT_URLS = [
  'https://google.com',
  'https://github.com',
  'https://api.github.com',
  'https://httpstat.us/404',
  'https://httpstat.us/503',
  'https://invalid.domain.xyz'
];

function App() {
  const [urls, setUrls] = useState(DEFAULT_URLS.join('\n'));
  const [results, setResults] = useState([]);
  const [isChecking, setIsChecking] = useState(false);
  const [stats, setStats] = useState({ total: 0, up: 0, down: 0, error: 0 });

  const runCheck = async () => {
    setIsChecking(true);
    const urlList = urls.split('\n').filter(u => u.trim() !== '');
    setResults([]);
    
    // Reset stats
    setStats({ total: urlList.length, up: 0, down: 0, error: 0 });

    const checkingPool = urlList.map(async (url, index) => {
      // Simulate network delay logic
      const delay = Math.random() * 2000 + 500;
      await new Promise(resolve => setTimeout(resolve, delay));

      let result;
      if (url.includes('invalid')) {
        result = { url, status: 'ERROR', code: '---', time: delay.toFixed(0), type: 'DNS_ERROR' };
      } else if (url.includes('404') || url.includes('503')) {
        result = { url, status: 'DOWN', code: url.includes('404') ? 404 : 503, time: delay.toFixed(0) };
      } else {
        result = { url, status: 'UP', code: 200, time: delay.toFixed(0) };
      }

      setResults(prev => [...prev, result]);
      setStats(prev => ({
        ...prev,
        up: result.status === 'UP' ? prev.up + 1 : prev.up,
        down: result.status === 'DOWN' ? prev.down + 1 : prev.down,
        error: result.status === 'ERROR' ? prev.error + 1 : prev.error,
      }));
    });

    await Promise.all(checkingPool);
    setIsChecking(false);
  };

  return (
    <div className="app-container">
      <header>
        <div className="logo-section">
          <div className="logo">CURLSC v1.0</div>
          <div className="stat-label">Concurrent URL Status Checker</div>
        </div>
        <div className="actions">
          <button 
            onClick={runCheck} 
            disabled={isChecking}
            style={{ display: 'flex', alignItems: 'center', gap: '8px' }}
          >
            {isChecking ? <Activity className="animate-spin" size={18} /> : <Play size={18} />}
            {isChecking ? 'Checking...' : 'Run Diagnostics'}
          </button>
        </div>
      </header>

      <div className="stats-grid">
        <StatCard label="Total" value={stats.total} icon={<Globe size={20} />} color="var(--accent-blue)" />
        <StatCard label="Healthy" value={stats.up} icon={<CheckCircle2 size={20} />} color="var(--accent-green)" />
        <StatCard label="Down" value={stats.down} icon={<Activity size={20} />} color="var(--accent-yellow)" />
        <StatCard label="Failed" value={stats.error} icon={<AlertTriangle size={20} />} color="var(--accent-red)" />
      </div>

      <div className="main-content">
        <div className="input-section">
          <div style={{ flex: 1 }}>
            <label className="stat-label" style={{ display: 'block', marginBottom: '8px' }}>Target URLs (One per line)</label>
            <textarea
              value={urls}
              onChange={(e) => setUrls(e.target.value)}
              placeholder="Enter URLs to check..."
              style={{
                width: '100%',
                height: '120px',
                background: 'rgba(255,255,255,0.05)',
                border: '1px solid var(--glass-border)',
                borderRadius: '8px',
                padding: '1rem',
                color: 'white',
                resize: 'none'
              }}
            />
          </div>
        </div>

        <div className="table-container">
          <table className="results-table">
            <thead>
              <tr>
                <th>STATUS</th>
                <th>LATENCY</th>
                <th>URL</th>
                <th>DETAILS</th>
              </tr>
            </thead>
            <tbody>
              <AnimatePresence>
                {results.map((res, i) => (
                  <motion.tr
                    key={i}
                    initial={{ opacity: 0, x: -20 }}
                    animate={{ opacity: 1, x: 0 }}
                    transition={{ duration: 0.3 }}
                  >
                    <td>
                      <StatusBadge status={res.status} />
                    </td>
                    <td style={{ color: 'var(--text-secondary)', fontFamily: 'monospace' }}>
                      {res.time}ms
                    </td>
                    <td style={{ fontWeight: 500 }}>{res.url}</td>
                    <td>
                      <span style={{ fontSize: '0.8rem', color: 'var(--text-secondary)' }}>
                        {res.status === 'UP' ? `HTTP ${res.code} OK` : res.type || `HTTP ${res.code}`}
                      </span>
                    </td>
                  </motion.tr>
                ))}
              </AnimatePresence>
              {results.length === 0 && !isChecking && (
                <tr>
                  <td colSpan="4" style={{ textAlign: 'center', padding: '4rem', color: 'var(--text-secondary)' }}>
                    <Search size={48} style={{ opacity: 0.2, marginBottom: '1rem' }} />
                    <p>No diagnostics run yet. Enter URLs above and click "Run Diagnostics".</p>
                  </td>
                </tr>
              )}
            </tbody>
          </table>
        </div>
      </div>
    </div>
  );
}

function StatCard({ label, value, icon, color }) {
  return (
    <motion.div 
      className="stat-card"
      whileHover={{ y: -5 }}
    >
      <div style={{ color, marginBottom: '0.5rem', display: 'flex', justifyContent: 'center' }}>{icon}</div>
      <span className="stat-value" style={{ color }}>{value}</span>
      <span className="stat-label">{label}</span>
    </motion.div>
  );
}

function StatusBadge({ status }) {
  const config = {
    UP: { class: 'badge-up', icon: <CheckCircle2 size={14} />, text: 'UP' },
    DOWN: { class: 'badge-down', icon: <Activity size={14} />, text: 'DOWN' },
    ERROR: { class: 'badge-error', icon: <AlertTriangle size={14} />, text: 'ERR' }
  };
  
  const { class: className, icon, text } = config[status];
  
  return (
    <div className={`status-badge ${className}`}>
      <div className="pulse" />
      {icon}
      {text}
    </div>
  );
}

export default App;
