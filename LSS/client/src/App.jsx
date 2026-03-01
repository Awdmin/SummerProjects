import { useState, useEffect } from 'react'
import './App.css'

const DAYS = ['Ned', 'Pon', 'Tor', 'Sre', 'Čet', 'Pet', 'Sob'];
const MONTHS = ['Jan','Feb','Mar','Apr','Maj','Jun','Jul','Aug','Sep','Okt','Nov','Dec'];

function Modal({ date, onClose, onSubmit }) {
  const [name, setName] = useState('');
  const [loading, setLoading] = useState(false);

  const handleSubmit = async () => {
    if (!name.trim()) return;
    setLoading(true);
    await onSubmit(name, date);
    setLoading(false);
    onClose();
  };

  return (
    <div style={styles.overlay} onClick={onClose}>
      <div style={styles.modal} onClick={e => e.stopPropagation()}>
        <button style={styles.closeBtn} onClick={onClose}>✕</button>
        <div style={styles.modalIcon}>📅</div>
        <h2 style={styles.modalTitle}>Sign Up</h2>
        <p style={styles.modalSubtitle}>
          {new Date(date).toLocaleDateString('sl-SI', { weekday: 'long', month: 'long', day: 'numeric' })}
        </p>
        <input
          style={styles.input}
          type="text"
          placeholder="Ime"
          value={name}
          onChange={e => setName(e.target.value)}
          onKeyDown={e => e.key === 'Enter' && handleSubmit()}
          autoFocus
        />
        <button
          style={{ ...styles.submitBtn, opacity: loading || !name.trim() ? 0.6 : 1 }}
          onClick={handleSubmit}
          disabled={loading || !name.trim()}
        >
          {loading ? 'Signing up...' : 'Potrdi'}
        </button>
      </div>
    </div>
  );
}

function DeletePopover({ card, onClose, onDelete }) {
  const [loading, setLoading] = useState(false);

  const handleDelete = async () => {
    setLoading(true);
    await onDelete(card);
    setLoading(false);
    onClose();
  };

  return (
    <div style={styles.overlay} onClick={onClose}>
      <div style={styles.modal} onClick={e => e.stopPropagation()}>
        <button style={styles.closeBtn} onClick={onClose}>✕</button>
        <div style={styles.modalIcon}>🗑️</div>
        <h2 style={styles.modalTitle}>Odstrani najavo</h2>
        <p style={styles.modalSubtitle}>Odstrani <strong style={{ color: '#f0ebe3' }}>{card.name}</strong> iz dneva?</p>
        <div style={{ display: 'flex', gap: 10, width: '100%', marginTop: 4 }}>
          <button
            style={{ ...styles.submitBtn, background: '#2a2825', color: '#c4bdb5', flex: 1 }}
            onClick={onClose}
          >
            Prekliči
          </button>
          <button
            style={{ ...styles.submitBtn, background: '#8b3a3a', flex: 1, opacity: loading ? 0.6 : 1 }}
            onClick={handleDelete}
            disabled={loading}
          >
            {loading ? 'Removing...' : 'Odstrani'}
          </button>
        </div>
      </div>
    </div>
  );
}

function App() {
  const [data, setData] = useState([]);
  const [modalDate, setModalDate] = useState(null);
  const [deleteTarget, setDeleteTarget] = useState(null);
  const URL = import.meta.env.VITE_API_URL;

  const fetchData = async () => {
    try {
      const response = await fetch(`${URL}/entries`);
      const json = await response.json();

      const result = [];
      const today = new Date();

      for (let i = 0; i < 7; i++) {
        const d = new Date();
        d.setDate(today.getDate() + i);
        const key = d.toISOString().split('T')[0];
        result.push({ date: key, items: [] });
      }

      (json || []).forEach((item) => {
        const itemDate = new Date(item.date).toISOString().split('T')[0];
        const bucket = result.find((r) => r.date === itemDate);
        if (bucket) bucket.items.push(item);
      });

      setData(result);
    } catch (error) {
      console.error(error);
    }
  };

  useEffect(() => {
    fetchData();
  }, [URL]);

  const handleSubmit = async (name, date) => {
    try {
      await fetch(`${URL}/entry`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ name, date }),
      });
      await fetchData();
    } catch (error) {
      console.error(error);
    }
  };

  const handleDelete = async (card) => {
    try {
      await fetch(`${URL}/entry`, {
        method: 'DELETE',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ id: card.id, name: card.name, date: card.date }),
      });
      await fetchData();
    } catch (error) {
      console.error(error);
    }
  };

  const formatDate = (dateStr) => {
    const d = new Date(dateStr);
    return { day: DAYS[d.getDay()], num: d.getDate(), month: MONTHS[d.getMonth()] };
  };

  const isToday = (dateStr) => {
    return dateStr === new Date().toISOString().split('T')[0];
  };

  return (
    <>
      <style>{`
        @import url('https://fonts.googleapis.com/css2?family=DM+Serif+Display:ital@0;1&family=DM+Sans:wght@300;400;500&display=swap');
        * { box-sizing: border-box; margin: 0; padding: 0; }
        body { background: #0f0e0d; font-family: 'DM Sans', sans-serif; min-height: 100vh; }
        @keyframes fadeIn { from { opacity: 0; transform: translateY(8px); } to { opacity: 1; transform: translateY(0); } }
        @keyframes modalIn { from { opacity: 0; transform: scale(0.95) translateY(12px); } to { opacity: 1; transform: scale(1) translateY(0); } }
        @keyframes pulse { 0%,100% { opacity:1; } 50% { opacity:0.5; } }
      `}</style>

      <div style={styles.root}>
        <header style={styles.header}>
          <h1 style={styles.heading}>Prisotnost na kosilu</h1>
          <p style={styles.subheading}>Naslednjih 7 dni</p>
        </header>

        <div style={styles.grid}>
          {data.map((column, i) => {
            const { day, num, month } = formatDate(column.date);
            const today = isToday(column.date);
            return (
              <div
                key={i}
                style={{
                  ...styles.column,
                  ...(today ? styles.columnToday : {}),
                  animationDelay: `${i * 60}ms`,
                }}
              >
                <div style={styles.dateHeader}>
                  <span style={{ ...styles.dayName, ...(today ? styles.todayAccent : {}) }}>{day}</span>
                  <span style={{ ...styles.dayNum, ...(today ? styles.todayAccent : {}) }}>{num}</span>
                  <span style={styles.monthName}>{month}</span>
                </div>

                <button
                  style={styles.signupBtn}
                  onClick={() => setModalDate(column.date)}
                  onMouseEnter={e => e.currentTarget.style.background = '#c8a96e'}
                  onMouseLeave={e => e.currentTarget.style.background = '#b8944f'}
                >
                  + Najava
                </button>

                <div style={styles.entries}>
                  {column.items.length === 0 && (
                    <p style={styles.emptyMsg}>Ni najav</p>
                  )}
                  {column.items.map((card, j) => (
                    <div
                      key={j}
                      style={{ ...styles.card, animationDelay: `${j * 50}ms`, cursor: 'pointer' }}
                      onClick={() => {
                          if(i !== 0) {
                              setDeleteTarget(card)}
                      }
                      }
                      onMouseEnter={e => e.currentTarget.style.background = '#2e2b22'}
                      onMouseLeave={e => e.currentTarget.style.background = '#222120'}
                    >
                      <span style={styles.avatar}>{card.name[0].toUpperCase()}</span>
                      <span style={styles.cardName}>{card.name}</span>
                      <span style={styles.deleteHint}>✕</span>
                    </div>
                  ))}
                </div>
              </div>
            );
          })}
        </div>
      </div>

      {modalDate && (
        <Modal
          date={modalDate}
          onClose={() => setModalDate(null)}
          onSubmit={handleSubmit}
        />
      )}

      {deleteTarget && (
        <DeletePopover
          card={deleteTarget}
          onClose={() => setDeleteTarget(null)}
          onDelete={handleDelete}
        />
      )}
    </>
  );
}

const styles = {
  root: {
    maxWidth: 1400,
    margin: '0 auto',
    padding: '48px 24px',
  },
  header: {
    marginBottom: 40,
  },
  heading: {
    fontFamily: "'DM Serif Display', serif",
    fontSize: 42,
    fontWeight: 400,
    color: '#f0ebe3',
    letterSpacing: '-0.5px',
  },
  subheading: {
    fontSize: 14,
    color: '#6b6560',
    marginTop: 6,
    fontWeight: 300,
    letterSpacing: '0.05em',
    textTransform: 'uppercase',
  },
  grid: {
    display: 'grid',
    gridTemplateColumns: 'repeat(7, 1fr)',
    gap: 12,
  },
  column: {
    background: '#1a1916',
    border: '1px solid #2a2825',
    borderRadius: 12,
    padding: '20px 14px',
    display: 'flex',
    flexDirection: 'column',
    gap: 12,
    animation: 'fadeIn 0.4s ease both',
    transition: 'border-color 0.2s',
  },
  columnToday: {
    border: '1px solid #b8944f',
    background: '#1d1b14',
  },
  dateHeader: {
    display: 'flex',
    flexDirection: 'column',
    alignItems: 'center',
    gap: 2,
    paddingBottom: 12,
    borderBottom: '1px solid #2a2825',
  },
  dayName: {
    fontSize: 11,
    fontWeight: 500,
    color: '#6b6560',
    letterSpacing: '0.1em',
    textTransform: 'uppercase',
  },
  dayNum: {
    fontFamily: "'DM Serif Display', serif",
    fontSize: 32,
    color: '#c4bdb5',
    lineHeight: 1,
  },
  monthName: {
    fontSize: 11,
    color: '#4a4845',
    letterSpacing: '0.08em',
    textTransform: 'uppercase',
  },
  todayAccent: {
    color: '#c8a96e',
  },
  signupBtn: {
    background: '#b8944f',
    color: '#0f0e0d',
    border: 'none',
    borderRadius: 7,
    padding: '8px 16px',
    fontSize: 13,
    fontWeight: 500,
    cursor: 'pointer',
    transition: 'background 0.2s',
    letterSpacing: '0.02em',
    fontFamily: "'DM Sans', sans-serif",
  },
  entries: {
    display: 'flex',
    flexDirection: 'column',
    gap: 6,
    flex: 1,
  },
  emptyMsg: {
    fontSize: 12,
    color: '#3a3835',
    textAlign: 'center',
    marginTop: 8,
    fontStyle: 'italic',
  },
  card: {
    display: 'flex',
    alignItems: 'center',
    gap: 8,
    background: '#222120',
    borderRadius: 8,
    padding: '8px 10px',
    animation: 'fadeIn 0.3s ease both',
  },
  avatar: {
    width: 26,
    height: 26,
    borderRadius: '50%',
    background: '#2e2b22',
    border: '1px solid #b8944f44',
    color: '#b8944f',
    fontSize: 12,
    fontWeight: 500,
    display: 'flex',
    alignItems: 'center',
    justifyContent: 'center',
    flexShrink: 0,
  },
  cardName: {
    fontSize: 13,
    color: '#c4bdb5',
    fontWeight: 400,
    overflow: 'hidden',
    textOverflow: 'ellipsis',
    whiteSpace: 'nowrap',
  },
  deleteHint: {
    marginLeft: 'auto',
    fontSize: 11,
    color: '#4a4845',
    flexShrink: 0,
    transition: 'color 0.2s',
  },
  // Modal
  overlay: {
    position: 'fixed',
    inset: 0,
    background: 'rgba(0,0,0,0.75)',
    backdropFilter: 'blur(4px)',
    display: 'flex',
    alignItems: 'center',
    justifyContent: 'center',
    zIndex: 1000,
  },
  modal: {
    background: '#1a1916',
    border: '1px solid #2e2b22',
    borderRadius: 16,
    padding: '36px 32px',
    width: 340,
    display: 'flex',
    flexDirection: 'column',
    alignItems: 'center',
    gap: 12,
    position: 'relative',
    animation: 'modalIn 0.25s ease',
  },
  closeBtn: {
    position: 'absolute',
    top: 14,
    right: 16,
    background: 'none',
    border: 'none',
    color: '#6b6560',
    fontSize: 16,
    cursor: 'pointer',
    padding: 4,
  },
  modalIcon: {
    fontSize: 32,
    marginBottom: 4,
  },
  modalTitle: {
    fontFamily: "'DM Serif Display', serif",
    fontSize: 26,
    color: '#f0ebe3',
    fontWeight: 400,
  },
  modalSubtitle: {
    fontSize: 13,
    color: '#6b6560',
    marginBottom: 8,
  },
  input: {
    width: '100%',
    background: '#111110',
    border: '1px solid #2e2b22',
    borderRadius: 8,
    padding: '12px 14px',
    color: '#f0ebe3',
    fontSize: 15,
    fontFamily: "'DM Sans', sans-serif",
    outline: 'none',
    transition: 'border-color 0.2s',
  },
  submitBtn: {
    width: '100%',
    background: '#b8944f',
    color: '#0f0e0d',
    border: 'none',
    borderRadius: 8,
    padding: '13px 0',
    fontSize: 15,
    fontWeight: 500,
    cursor: 'pointer',
    marginTop: 4,
    fontFamily: "'DM Sans', sans-serif",
    transition: 'opacity 0.2s, background 0.2s',
  },
};

export default App;
