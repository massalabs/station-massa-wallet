export function handleKeyDown(e: React.KeyboardEvent) {
  if (e.key === 'Enter') {
    e.preventDefault();
    console.log('Enter key pressed');
    // Perform the desired action for the Enter key press
  }
}
