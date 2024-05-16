import { Toast } from '@massalabs/react-ui-kit';
import { useLocalStorage } from '@massalabs/react-ui-kit/src/lib/util/hooks/useLocalStorage';
import { FiSun, FiMoon } from 'react-icons/fi';
import { Outlet } from 'react-router-dom';

type ThemeSettings = {
  [key: string]: {
    icon: JSX.Element;
    label: string;
  };
};

const THEME_STORAGE_KEY = 'massa-wallet-theme';

const themeSettings: ThemeSettings = {
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
    THEME_STORAGE_KEY,
    'theme-dark',
  );
  const themeIcon = themeSettings[theme].icon;
  const themeLabel = themeSettings[theme].label;
  const context = { themeLabel, themeIcon, theme, handleSetTheme };

  function handleSetTheme() {
    setTheme(theme === 'theme-dark' ? 'theme-light' : 'theme-dark');
  }

  return (
    <div className={theme}>
      <Outlet context={context} />
      <Toast />
    </div>
  );
}

export default Base;
