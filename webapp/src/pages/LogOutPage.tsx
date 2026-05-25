import { useEffect } from 'react';
import { useHistory } from 'react-router-dom';
import { API_LOG_OUT } from '../api';
import { resetAppAction } from '../store/application/actions'
import { useDispatch } from 'react-redux'

function LogoutPage() {
  const history = useHistory();
  const dispatch = useDispatch();

  useEffect(() => {
    API_LOG_OUT().then(resp => {
      if (resp.ok) {
        dispatch(resetAppAction())
        history.push("/")
      }
    })
  }, [history, dispatch]);

  return <div />;
}

export default LogoutPage;
