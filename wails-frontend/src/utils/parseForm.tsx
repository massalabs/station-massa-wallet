import { BaseSyntheticEvent } from 'react';

export function parseForm(e: BaseSyntheticEvent) {
  let form = new FormData(e.target);

  return Object.fromEntries(form.entries());
}
