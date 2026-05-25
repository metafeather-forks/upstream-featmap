import { useEffect, useState } from 'react';
import { useParams, Link } from 'react-router-dom';
import { API_VERIFY_EMAIL } from '../api'

function VerifyEmailPage() {
  const { key } = useParams<{ key: string }>();
  const [ok, setOk] = useState(false);

  useEffect(() => {
    API_VERIFY_EMAIL(key).then(response => {
      if (response.ok) setOk(true)
    })
  }, [key]);

  return ok ?
    <div className="p-2">
      <em>Email address verified!</em>
      <br /><br />
      <Link className="link" to="/">Back to Featmap</Link>
    </div>
    : <>Something went wrong</>;
}

export default VerifyEmailPage;
