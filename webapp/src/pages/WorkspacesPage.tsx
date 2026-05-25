import { useState } from 'react';
import Header from '../components/Header';
import { useSelector } from 'react-redux'
import { Link } from 'react-router-dom'
import { Button, CardLayout } from '../components/elements';
import CreateWorkspaceModal from '../components/CreateWorkspaceModal';
import TimeAgo from 'react-timeago'
import { AppState } from '../store';

function WorkspacesPage() {
  const application = useSelector((state: AppState) => state.application.application);
  const [showCreate, setShowCreate] = useState(false);

  return (
    <div>
      <Header account={application.account!} />
      {showCreate ? <CreateWorkspaceModal close={() => setShowCreate(false)} /> : null}
      <div>
        <div className="p-2 flex flex-row mb-2 items-center">
          <div><h3>Workspaces</h3></div>
          <div className="ml-2"><Button title="New workspace" primary icon="add" handleOnClick={() => setShowCreate(true)} /></div>
        </div>
        <CardLayout>
          <div>
            {(application.workspaces && application.workspaces.length > 0) ?
              <div className="flex flex-col max-w-lg">
                <div className="p-2">{application.workspaces.length} workspace(s)</div>
                <div>
                  {application.workspaces.map(x =>
                    <div className="p-2" key={x.id}>
                      <div className="mb-1">
                        <p><b><Link to={"/" + x.name}>{x.name}</Link></b></p>
                        <p className="text-xs">Created <TimeAgo date={x.createdAt} />.</p>
                      </div>
                    </div>
                  )}
                </div>
              </div>
              : "No workspaces"
            }
          </div>
        </CardLayout>
      </div>
    </div>
  );
}

export default WorkspacesPage;
