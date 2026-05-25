import { IAccount } from '../store/application/types';
import { Link } from 'react-router-dom'
import ContextMenu from './ContextMenu';

interface Props {
  account: IAccount
  workspaceName?: string
}

function Header(props: Props) {
  return (
    <header className="bg-gray-200" role="banner">
      <div className="flex items-center p-1">
        <div className="flex text-lg m-1 w-24">
          <b><Link to="/">Featmap</Link></b>
        </div>
        {props.workspaceName ?
          <div className="text-sm"><Link to={"/" + props.workspaceName}>Projects</Link></div>
          : null
        }
        <div className="flex flex-grow justify-end">
          <div className="flex p-1 rounded items-center">
            <ContextMenu icon="account_circle" aria-label="Account menu">
              <div className="rounded bg-white shadow-md absolute mt-8 top-0 right-0 min-w-full text-sm">
                <ul className="list-reset">
                  {props.workspaceName ?
                    <li className="whitespace-nowrap px-4 py-2 block text-gray-600">
                      Current workspace: <b><Link to={"/" + props.workspaceName}>{props.workspaceName}</Link></b> <Link className="p-1 flex-no-shrink font-bold rounded text-xs bg-gray-300" to="/account/workspaces">Change</Link>
                    </li>
                    : null
                  }
                  {props.workspaceName ?
                    <li className="whitespace-nowrap px-4 py-2 block text-black hover:bg-gray-200">
                      <Link to={"/" + props.workspaceName + "/settings"}>Workspace settings</Link>
                    </li>
                    : null
                  }
                  <li className="whitespace-nowrap px-4 py-2 block text-black">
                    <Link to="/account/workspaces">My workspaces</Link>
                  </li>
                  <li className="whitespace-nowrap block text-black">
                    <hr className="border-b" />
                  </li>
                  <li className="whitespace-nowrap px-4 py-2 block text-gray-600">
                    Logged in as <em>{props.account.email}</em>
                  </li>
                  <li><Link className="whitespace-nowrap px-4 py-2 block text-black hover:bg-gray-200" to="/account/settings">Account settings</Link></li>
                  <li><Link className="whitespace-nowrap px-4 py-2 block text-black hover:bg-gray-200" to="/account/logout">Log out</Link></li>
                </ul>
              </div>
            </ContextMenu>
          </div>
        </div>
      </div>
    </header>
  );
}

export default Header;
