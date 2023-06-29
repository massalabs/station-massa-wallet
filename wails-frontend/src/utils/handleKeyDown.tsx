export function handleKeyDown(
  event: React.KeyboardEvent<HTMLFormElement>,
  handleSubmit: (e: React.SyntheticEvent) => Promise<void>,
) {
  if (event.key === 'Enter') {
    event.preventDefault();
    handleSubmit(event);
  }
}
