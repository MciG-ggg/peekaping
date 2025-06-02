import React, { createContext, useContext, useEffect, useState } from 'react';

const TIMEZONE_STORAGE_KEY = 'selectedTimezone';

// Get browser default timezone
const getDefaultTimezone = () =>
  Intl.DateTimeFormat().resolvedOptions().timeZone || 'UTC';

interface TimezoneContextType {
  timezone: string;
  setTimezone: (tz: string) => void;
}

const TimezoneContext = createContext<TimezoneContextType | undefined>(undefined);

export const TimezoneProvider: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const [timezone, setTimezoneState] = useState<string>(() => {
    return localStorage.getItem(TIMEZONE_STORAGE_KEY) || getDefaultTimezone();
  });

  useEffect(() => {
    localStorage.setItem(TIMEZONE_STORAGE_KEY, timezone);
  }, [timezone]);

  const setTimezone = (tz: string) => {
    setTimezoneState(tz);
  };

  return (
    <TimezoneContext.Provider value={{ timezone, setTimezone }}>
      {children}
    </TimezoneContext.Provider>
  );
};

export const useTimezone = () => {
  const ctx = useContext(TimezoneContext);
  if (!ctx) throw new Error('useTimezone must be used within a TimezoneProvider');
  return ctx;
};
