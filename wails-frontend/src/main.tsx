import { h, render } from 'preact';
import { App } from './app';
import { BrowserRouter } from 'react-router-dom';

render(
  <BrowserRouter>
    <App />
  </BrowserRouter>,
  document.getElementById('app')!,
);
