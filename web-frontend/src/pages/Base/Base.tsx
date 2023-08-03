import { Toast } from '@massalabs/react-ui-kit';
import { FiSun, FiMoon } from 'react-icons/fi';
import { Outlet } from 'react-router-dom';

import { useLocalStorage } from '@/custom/useLocalStorage';

type ThemeSettings = {
  [key: string]: {
    icon: JSX.Element;
    label: string;
  };
};

export const themeSettings: ThemeSettings = {
  'theme-dark': {
    icon: <FiSun />,
    label: 'light theme',
  },
  'theme-light': {
    icon: <FiMoon />,
    label: 'dark theme',
  },
};

function Base() {
  const [theme, setTheme] = useLocalStorage<string>(
    'massa-station-theme',
    'theme-dark',
  );
  const themeIcon = themeSettings[theme].icon;
  const themeLabel = themeSettings[theme].label;
  const context = { themeLabel, themeIcon, theme, handleSetTheme };

  function handleSetTheme() {
    setTheme(theme === 'theme-dark' ? 'theme-light' : 'theme-dark');
  }

  return (
    // TODO
    // remove theme-dark
    // this needs to be removed as soon we fix the steps to create an account
    <div className={`${theme} theme-dark`}>
      <Outlet context={context} />
      <Toast theme={theme} />
    </div>
  );
}

export default Base;
