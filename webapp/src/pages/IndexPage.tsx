import { useEffect, useState } from 'react';
import AccountPage from './AccountPage';
import WorkspacePage from './WorkspacePage';
import WorkspacesPage from './WorkspacesPage';
import NotFound from './NotFound';
import Footer from '../components/Footer';
import { Route, Switch, Redirect, useRouteMatch } from 'react-router-dom'
import { useSelector, useDispatch } from 'react-redux'
import { AppState } from '../store'
import { IApplication } from '../store/application/types'
import { receiveAppAction } from '../store/application/actions';
import { API_FETCH_APP, API_FETCH_APP_RESP } from '../api';
import { useHistory } from 'react-router-dom';

function Index() {
  const match = useRouteMatch();
  const history = useHistory();
  const dispatch = useDispatch();
  const application = useSelector((state: AppState) => state.application.application);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    API_FETCH_APP().then(response => {
      if (response.ok) {
        response.json().then((data: API_FETCH_APP_RESP) => {
          dispatch(receiveAppAction(data))
          setLoading(false)
        })
      } else {
        history.push("/account/login")
      }
    })
  }, [history, dispatch]);

  if (loading) {
    return <div className="p-2">Loading data...</div>;
  }
  if (application.account) {
    return (
      <div>
        <Switch>
          <Route exact path={match.path + ""} render={() => {
            if (application.workspaces && application.workspaces.length === 1) {
              return <Redirect to={application.workspaces[0].name} />;
            }
            return <Redirect to={match.path + "account/workspaces"} />;
          }} />
          <Route exact path={match.path + "account/workspaces"} component={WorkspacesPage} />
          <Route exact path={match.path + "account/settings"} component={AccountPage} />
          <Route path={match.path + "account/"} component={NotFound} />
          <Route path={match.path + ":workspaceName"} component={WorkspacePage} />
          <Route path={match.path} component={NotFound} />
        </Switch>
        <Footer />
      </div>
    );
  }
  return <Redirect to="/account/login" />;
}

export default Index;
