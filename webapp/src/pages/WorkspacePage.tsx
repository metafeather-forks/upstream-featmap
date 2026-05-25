import { useEffect, useState } from 'react';
import { Route, Switch, Redirect, useParams, useRouteMatch } from 'react-router-dom';
import { useSelector, useDispatch } from 'react-redux';
import { AppState } from '../store';
import { getWorkspaceByName, application } from '../store/application/selectors';
import NotFound from './NotFound';
import ProjectPage from './ProjectPage';
import Header from '../components/Header';
import { IApplication } from '../store/application/types';
import { loadProjectsAction } from '../store/projects/actions';
import { IProject } from '../store/projects/types';
import ProjectsPage from './ProjectsPage';
import { API_GET_PROJECTS } from '../api'
import WorkspaceSettingsPage from './WorkspaceSettingsPage';

function WorkspacePage() {
  const { workspaceName } = useParams<{ workspaceName: string }>();
  const match = useRouteMatch();
  const dispatch = useDispatch();
  const app = useSelector((state: AppState) => application(state));
  const [loading, setLoading] = useState(true);
  const [notFound, setNotFound] = useState(false);

  useEffect(() => {
    const ws = getWorkspaceByName(app, workspaceName);
    if (!ws) { setNotFound(true); return; }
    API_GET_PROJECTS(ws.id).then(response => {
      if (response.ok) {
        response.json().then((data: IProject[]) => {
          dispatch(loadProjectsAction(data));
          setLoading(false);
        })
      }
    })
  }, [app, workspaceName, dispatch]);

  const ws = getWorkspaceByName(app, workspaceName);

  if (notFound) return <div><Redirect to="/" /></div>;
  if (loading) return <div className="p-2">Loading data...</div>;

  return (
    <div>
      <Header account={app.account!} workspaceName={workspaceName} />
      <Switch>
        <Route exact strict path={match.path} component={ProjectsPage} />
        <Route exact strict path={match.path + "/settings"} component={WorkspaceSettingsPage} />
        <Route strict path={match.path + "/projects/:projectId"} component={ProjectPage} />
        <Route path={match.path} component={NotFound} />
      </Switch>
    </div>
  );
}

export default WorkspacePage;
