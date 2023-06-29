export function handleKeyDown(e: React.KeyboardEvent<HTMLFormElement>) {
  if (e.key === 'Enter') {
    e.preventDefault();
  }
}
