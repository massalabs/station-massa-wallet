import { Link } from 'react-router-dom';

export default function Error() {
  return (
    <div id="error-page">
      <h1 className="mas-banner">Oops!</h1>
      <p className="mas-body">Sorry, an unexpected error has occurred.</p>
      <Link to="index">Go back to the Welcome page!</Link>
    </div>
  );
}
