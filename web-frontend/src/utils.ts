import { NavigateFunction } from 'react-router-dom';

export function routeFor(path: string) {
  return `${import.meta.env.VITE_BASE_APP}/${path}`;
}

export function navigateToImportAccount(navigate: NavigateFunction) {
  return () => {
    navigate(routeFor('account-import'));
  };
}

export function goToErrorPage(navigate: NavigateFunction) {
  navigate(routeFor('error'));
}
