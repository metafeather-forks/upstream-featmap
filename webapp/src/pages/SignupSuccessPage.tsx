import { useHistory } from 'react-router-dom';
import { Button } from '../components/elements'

function SignupSuccessPage() {
  const history = useHistory();
  return (
    <div className="flex p-2 w-full justify-center items-center flex-col">
      <div className="flex p-3 max-w-xl w-full items-center flex-col">
        <div className="flex p-2 flex-col items-baseline">
          <div className="p-1 text-2xl font-bold">Welcome to Featmap!</div>
          <div className="p-1 text-center"><Button title="Get started" primary handleOnClick={() => history.push("/")} /></div>
        </div>
      </div>
    </div>
  );
}

export default SignupSuccessPage;
